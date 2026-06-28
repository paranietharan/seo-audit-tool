import os
import json
import google.generativeai as genai
from models import GeminiReport, Issue

# Configure Gemini API
api_key = os.getenv("GEMINI_API_KEY", "")
if api_key:
    genai.configure(api_key=api_key)

def generate_report(audit_data: dict) -> GeminiReport:
    if not api_key:
        raise ValueError("GEMINI_API_KEY environment variable is not set")

    system_prompt = "You are an SEO expert. Analyze the provided website audit data and return ONLY valid JSON with no markdown, no explanation."
    model = genai.GenerativeModel("gemini-2.5-flash", system_instruction=system_prompt)

    # Build prompt instructing Gemini
    user_prompt = (
        f"Audit data: {json.dumps(audit_data)}\n\n"
        "Return JSON format:\n"
        "{\n"
        '  "score": number,\n'
        '  "summary": "string",\n'
        '  "issues": [{"severity": "critical|warning|info", "category": "string", "description": "string", "recommendation": "string"}],\n'
        '  "quick_wins": ["string"],\n'
        '  "strengths": ["string"]\n'
        "}"
    )

    response = model.generate_content(
        contents=user_prompt,
        generation_config={
            "response_mime_type": "application/json"
        }
    )

    text = response.text.strip()

    # Clean markdown code block fences if any
    if text.startswith("```"):
        # Strip start
        if text.startswith("```json"):
            text = text[7:]
        elif text.startswith("```"):
            text = text[3:]
        # Strip end
        if text.endswith("```"):
            text = text[:-3]
        text = text.strip()

    parsed = json.loads(text)
    
    # Map raw dictionary to Pydantic models
    issues_list = []
    for iss in parsed.get("issues", []):
        issues_list.append(Issue(
            severity=iss.get("severity", "info"),
            category=iss.get("category", "General"),
            description=iss.get("description", ""),
            recommendation=iss.get("recommendation", "")
        ))

    return GeminiReport(
        score=int(parsed.get("score", 0)),
        summary=str(parsed.get("summary", "")),
        issues=issues_list,
        quick_wins=list(parsed.get("quick_wins", [])),
        strengths=list(parsed.get("strengths", []))
    )
