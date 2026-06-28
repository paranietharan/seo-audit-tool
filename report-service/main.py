import os
import httpx
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from models import ReportRequest, ReportResponse
from gemini_client import generate_report

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

AUDIT_SERVICE_URL = os.getenv("AUDIT_SERVICE_URL", "http://localhost:8081")

@app.post("/report", response_model=ReportResponse)
async def create_report(req: ReportRequest):
    async with httpx.AsyncClient(timeout=60.0) as client:
        try:
            response = await client.post(
                f"{AUDIT_SERVICE_URL}/audit",
                json={"url": req.url}
            )
            response.raise_for_status()
            audit_json = response.json()
        except httpx.HTTPError as exc:
            raise HTTPException(
                status_code=502,
                detail=f"Failed to communicate with audit-service: {str(exc)}"
            )

    try:
        gemini_report = generate_report(audit_json)
    except Exception as exc:
        raise HTTPException(
            status_code=500,
            detail=f"Failed to generate report from Gemini API: {str(exc)}"
        )

    return ReportResponse(audit=audit_json, report=gemini_report)

@app.get("/health")
def health():
    return {"status": "ok"}
