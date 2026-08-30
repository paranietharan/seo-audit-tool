from typing import Any, Dict, List
from pydantic import BaseModel, Field


class ReportRequest(BaseModel):
    url: str = Field(..., description="Target website URL to audit")


class Issue(BaseModel):
    severity: str = Field(default="info", description="Severity level: critical, warning, or info")
    category: str = Field(default="General", description="Category of the SEO issue")
    description: str = Field(default="", description="Detailed explanation of the issue")
    recommendation: str = Field(default="", description="Actionable fix or advice")


class GeminiReport(BaseModel):
    score: int = Field(default=0, ge=0, le=100, description="Overall SEO score from 0 to 100")
    summary: str = Field(default="", description="Executive summary of the SEO analysis")
    issues: List[Issue] = Field(default_factory=list, description="List of detected SEO issues")
    quick_wins: List[str] = Field(default_factory=list, description="List of high-impact quick fixes")
    strengths: List[str] = Field(default_factory=list, description="List of positive SEO aspects")


class ReportResponse(BaseModel):
    audit: Dict[str, Any] = Field(..., description="Raw technical audit data")
    report: GeminiReport = Field(..., description="AI-generated SEO analysis report")
