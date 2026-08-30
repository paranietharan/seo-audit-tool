import logging
import os
from contextlib import asynccontextmanager
from typing import AsyncIterator

import httpx
from fastapi import FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from gemini_client import generate_report_async
from models import ReportRequest, ReportResponse

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("report-service")

AUDIT_SERVICE_URL = os.getenv("AUDIT_SERVICE_URL", "http://localhost:8081").rstrip("/")
http_client: httpx.AsyncClient | None = None


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    global http_client
    http_client = httpx.AsyncClient(timeout=60.0)
    logger.info(f"Report service initialized with AUDIT_SERVICE_URL={AUDIT_SERVICE_URL}")
    yield
    if http_client:
        await http_client.aclose()


app = FastAPI(
    title="SEO Report Service",
    description="Microservice providing AI-powered SEO analysis using Google Gemini and rule heuristics",
    version="1.0.0",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.post("/report", response_model=ReportResponse)
async def create_report(req: ReportRequest):
    global http_client
    if http_client is None:
        http_client = httpx.AsyncClient(timeout=60.0)

    # 1. Fetch raw audit data from audit-service
    try:
        response = await http_client.post(
            f"{AUDIT_SERVICE_URL}/audit",
            json={"url": req.url},
        )
        if response.status_code != 200:
            err_detail = "Failed to audit website"
            try:
                err_json = response.json()
                err_detail = err_json.get("error") or err_json.get("detail") or err_detail
            except Exception:
                err_detail = response.text or err_detail
            raise HTTPException(
                status_code=response.status_code,
                detail=f"Audit service error: {err_detail}",
            )
        audit_json = response.json()
    except httpx.HTTPError as exc:
        logger.error(f"HTTP communication with audit-service failed: {exc}")
        raise HTTPException(
            status_code=status.HTTP_502_BAD_GATEWAY,
            detail=f"Failed to communicate with audit-service at {AUDIT_SERVICE_URL}: {str(exc)}",
        )

    # 2. Asynchronously generate SEO report via Gemini / Fallback
    try:
        gemini_report = await generate_report_async(audit_json)
    except Exception as exc:
        logger.error(f"Report generation error: {exc}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Report generation failed: {str(exc)}",
        )

    return ReportResponse(audit=audit_json, report=gemini_report)


@app.get("/health")
async def health():
    return {"status": "ok", "service": "report-service"}
