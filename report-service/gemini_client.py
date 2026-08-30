import asyncio
import json
import logging
import os
import re
from typing import Any, Dict, List
import google.generativeai as genai
from models import GeminiReport, Issue

logger = logging.getLogger("report-service.gemini")


def _clean_json_markdown(text: str) -> str:
    """Strip markdown code fence blocks if returned by the LLM."""
    text = text.strip()
    match = re.search(r"```(?:json)?\s*([\s\S]*?)\s*```", text)
    if match:
        return match.group(1).strip()
    if text.startswith("```"):
        text = re.sub(r"^```(?:json)?", "", text)
        text = re.sub(r"```$", "", text)
    return text.strip()


def _fallback_rule_based_report(audit_data: Dict[str, Any]) -> GeminiReport:
    """Generate a deterministic rule-based report when Gemini API is unavailable."""
    score = 100
    issues: List[Issue] = []
    quick_wins: List[str] = []
    strengths: List[str] = []

    title = audit_data.get("title", "")
    if not title:
        score -= 20
        issues.append(Issue(
            severity="critical",
            category="On-Page SEO",
            description="Page is missing a <title> tag.",
            recommendation="Add a descriptive <title> tag (50-60 characters) containing primary keywords."
        ))
        quick_wins.append("Add a concise, keyword-rich title tag to the page header.")
    else:
        strengths.append(f"Title tag is present ({len(title)} characters).")

    meta_desc = audit_data.get("meta_description", "")
    if not meta_desc:
        score -= 15
        issues.append(Issue(
            severity="warning",
            category="Meta Tags",
            description="Missing meta description tag.",
            recommendation="Add a compelling meta description (120-160 characters) summarizing page content."
        ))
        quick_wins.append("Write a relevant meta description to improve search engine click-through rates.")
    else:
        strengths.append("Meta description is properly configured.")

    h1_count = audit_data.get("h1_count", 0)
    if h1_count == 0:
        score -= 15
        issues.append(Issue(
            severity="critical",
            category="Headings Structure",
            description="No <h1> heading found on the page.",
            recommendation="Ensure the page has exactly one main <h1> heading representing the primary topic."
        ))
        quick_wins.append("Add a single, prominent <h1> tag at the top of your main content.")
    elif h1_count > 1:
        score -= 5
        issues.append(Issue(
            severity="warning",
            category="Headings Structure",
            description=f"Multiple <h1> headings found ({h1_count}).",
            recommendation="Structure headings with a single <h1> and subordinate <h2>/<h3> elements."
        ))
    else:
        strengths.append("Single <h1> heading structure correctly implemented.")

    missing_alt = audit_data.get("images_missing_alt", [])
    if missing_alt:
        score -= min(15, len(missing_alt) * 3)
        issues.append(Issue(
            severity="warning",
            category="Accessibility & Media",
            description=f"{len(missing_alt)} image(s) missing alt text attributes.",
            recommendation="Add descriptive 'alt' attributes to all content images for SEO and screen-reader accessibility."
        ))
        quick_wins.append(f"Add alt attributes to {len(missing_alt)} missing images.")
    else:
        strengths.append("All discovered images have descriptive alt attributes.")

    if not audit_data.get("has_structured_data", False):
        score -= 10
        issues.append(Issue(
            severity="info",
            category="Structured Data",
            description="No JSON-LD schema structured data found.",
            recommendation="Implement Schema.org structured data (e.g., WebPage, Article, Organization) to enhance search snippets."
        ))
    else:
        strengths.append("Schema.org JSON-LD structured data detected.")

    canonical = audit_data.get("canonical_url", "")
    if not canonical:
        score -= 5
        issues.append(Issue(
            severity="info",
            category="Indexing",
            description="Canonical URL tag is missing.",
            recommendation="Specify a rel='canonical' link element to prevent duplicate content issues."
        ))
    else:
        strengths.append("Canonical URL link tag is configured.")

    load_ms = audit_data.get("page_load_ms", 0)
    if load_ms > 3000:
        score -= 15
        issues.append(Issue(
            severity="warning",
            category="Performance",
            description=f"Page response time is slow ({load_ms} ms).",
            recommendation="Optimize asset payloads, minify scripts, and enable server-side caching."
        ))
    elif load_ms > 0:
        strengths.append(f"Fast page load time ({load_ms} ms).")

    score = max(10, min(100, score))
    summary = f"Rule-based SEO audit evaluated {audit_data.get('url', 'the website')} with an overall score of {score}/100 based on core technical SEO metrics."

    return GeminiReport(
        score=score,
        summary=summary,
        issues=issues,
        quick_wins=quick_wins,
        strengths=strengths
    )


def _generate_report_sync(audit_data: dict) -> GeminiReport:
    """Synchronous report generation calling the Gemini API."""
    api_key = os.getenv("GEMINI_API_KEY", "").strip()
    if not api_key:
        logger.warning("GEMINI_API_KEY is not configured. Falling back to rule-based analysis.")
        return _fallback_rule_based_report(audit_data)

    model_name = os.getenv("GEMINI_MODEL", "gemini-2.0-flash").strip()
    genai.configure(api_key=api_key)

    system_prompt = (
        "You are a professional technical SEO specialist and web auditor. "
        "Analyze the provided website audit data thoroughly and return ONLY a valid JSON object matching the requested schema. "
        "Do not include markdown code block markers or any commentary outside the JSON."
    )

    try:
        model = genai.GenerativeModel(model_name, system_instruction=system_prompt)

        user_prompt = (
            f"Audit Data:\n{json.dumps(audit_data, indent=2)}\n\n"
            "Return a JSON object with this exact structure:\n"
            "{\n"
            '  "score": <integer from 0 to 100>,\n'
            '  "summary": "<2-3 sentence executive summary of SEO performance>",\n'
            '  "issues": [\n'
            '    {\n'
            '      "severity": "critical" | "warning" | "info",\n'
            '      "category": "<e.g., On-Page, Performance, Accessibility, Technical>",\n'
            '      "description": "<concise explanation of the issue>",\n'
            '      "recommendation": "<practical, actionable fix>"\n'
            "    }\n"
            "  ],\n"
            '  "quick_wins": ["<high-impact low-effort action 1>", "<action 2>"],\n'
            '  "strengths": ["<positive SEO signal 1>", "<signal 2>"]\n'
            "}"
        )

        response = model.generate_content(
            contents=user_prompt,
            generation_config={
                "response_mime_type": "application/json",
                "temperature": 0.2
            }
        )

        raw_text = response.text or ""
        cleaned_json = _clean_json_markdown(raw_text)
        parsed = json.loads(cleaned_json)

        issues_list = []
        for iss in parsed.get("issues", []):
            issues_list.append(Issue(
                severity=str(iss.get("severity", "info")).lower(),
                category=str(iss.get("category", "General")),
                description=str(iss.get("description", "")),
                recommendation=str(iss.get("recommendation", ""))
            ))

        return GeminiReport(
            score=int(parsed.get("score", 0)),
            summary=str(parsed.get("summary", "")),
            issues=issues_list,
            quick_wins=[str(w) for w in parsed.get("quick_wins", [])],
            strengths=[str(s) for s in parsed.get("strengths", [])]
        )
    except Exception as exc:
        logger.error(f"Gemini API generation failed ({exc}). Falling back to rule-based report.")
        fallback = _fallback_rule_based_report(audit_data)
        fallback.summary += f" (Note: AI summary fallback due to: {str(exc)})"
        return fallback


async def generate_report_async(audit_data: dict) -> GeminiReport:
    """Asynchronous wrapper ensuring the event loop is never blocked."""
    return await asyncio.to_thread(_generate_report_sync, audit_data)
