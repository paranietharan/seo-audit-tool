from pydantic import BaseModel
from typing import List, Dict, Any

class ReportRequest(BaseModel):
    url: str

class Issue(BaseModel):
    severity: str
    category: str
    description: str
    recommendation: str

class GeminiReport(BaseModel):
    score: int
    summary: str
    issues: List[Issue]
    quick_wins: List[str]
    strengths: List[str]

class ReportResponse(BaseModel):
    audit: Dict[str, Any]
    report: GeminiReport
