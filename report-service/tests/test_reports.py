import json
import os
from unittest.mock import MagicMock, patch
import pytest
from fastapi.testclient import TestClient
from gemini_client import _generate_report_sync
from main import app


def test_missing_key_uses_deterministic_fallback():
    with patch.dict(os.environ, {}, clear=True):
        report = _generate_report_sync({"url": "https://example.com"})
    assert report.score == 35
    assert len(report.issues) == 5


def test_model_must_be_explicit():
    with patch.dict(os.environ, {"GEMINI_API_KEY": "test-key"}, clear=True), patch('gemini_client.genai.Client') as client:
        assert 'Rule-based' in _generate_report_sync({}).summary
        client.assert_not_called()


@pytest.mark.parametrize('response', ['not json', '{"score": 200}', None])
def test_invalid_ai_responses_fall_back(response):
    with patch.dict(os.environ, {"GEMINI_API_KEY": "test-key", "GEMINI_MODEL": "test-model"}), patch('gemini_client.genai.Client') as factory:
        factory.return_value.__enter__.return_value.models.generate_content.return_value.text = response
        assert 'Rule-based' in _generate_report_sync({}).summary


def test_new_sdk_returns_report_and_closes_client():
    with patch.dict(os.environ, {"GEMINI_API_KEY": "test-key", "GEMINI_MODEL": "test-model"}), patch('gemini_client.genai.Client') as factory:
        client = factory.return_value.__enter__.return_value
        client.models.generate_content.return_value.text = json.dumps({"score": 90, "summary": "Good title", "issues": []})
        report = _generate_report_sync({"title": "Example"})
        assert report.score == 90
        assert client.models.generate_content.call_args.kwargs['model'] == 'test-model'
        factory.return_value.__exit__.assert_called_once()


def test_provider_error_does_not_expose_details():
    with patch.dict(os.environ, {"GEMINI_API_KEY": "test-key", "GEMINI_MODEL": "test-model"}), patch('gemini_client.genai.Client', side_effect=RuntimeError('private-provider-detail')):
        assert 'private-provider-detail' not in _generate_report_sync({}).summary


@pytest.mark.parametrize('target', ['not-a-url', 'ftp://example.com', 'file:///etc/passwd'])
def test_rejects_invalid_audit_targets(target):
    with TestClient(app) as client:
        assert client.post('/report', json={'url': target}).status_code == 422
