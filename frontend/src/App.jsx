import React, { useState, useEffect } from 'react'
import UrlForm from './components/UrlForm'
import ScoreGauge from './components/ScoreGauge'
import IssuesTable from './components/IssuesTable'
import ReportCards from './components/ReportCards'
import AuditAccordion from './components/AuditAccordion'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'
const STORAGE_KEY = 'seo_audit_history'

export default function App() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [data, setData] = useState(null)
  const [history, setHistory] = useState([])

  useEffect(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        setHistory(JSON.parse(saved))
      }
    } catch (e) {
      console.error('Failed to load history from localStorage', e)
    }
  }, [])

  const saveToHistory = (auditResult) => {
    try {
      const item = {
        url: auditResult.audit?.url,
        score: auditResult.report?.score,
        date: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        data: auditResult,
      }
      const updated = [item, ...history.filter(h => h.url !== item.url)].slice(0, 5)
      setHistory(updated)
      localStorage.setItem(STORAGE_KEY, JSON.stringify(updated))
    } catch (e) {
      console.error('Failed to save to localStorage', e)
    }
  }

  const handleAuditSubmit = async (url) => {
    setLoading(true)
    setError(null)
    setData(null)

    try {
      const response = await fetch(`${API_BASE_URL}/report`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      })

      if (!response.ok) {
        let errMessage = `HTTP error ${response.status}`
        try {
          const errJson = await response.json()
          errMessage = errJson.detail || errJson.error || errMessage
        } catch {
          const errText = await response.text()
          if (errText) errMessage = errText
        }
        throw new Error(errMessage)
      }

      const resData = await response.json()
      setData(resData)
      saveToHistory(resData)
    } catch (err) {
      setError(err.message || 'Failed to connect to the audit service')
    } finally {
      setLoading(false)
    }
  }

  const handleExportJSON = () => {
    if (!data) return
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    const hostname = data.audit?.url ? new URL(data.audit.url).hostname : 'seo-audit'
    a.href = url
    a.download = `${hostname}-seo-report.json`
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-12">
      {/* Header */}
      <header className="text-center mb-10">
        <div className="inline-flex items-center gap-2 px-3 py-1 bg-emerald-500/10 border border-emerald-500/20 rounded-full text-xs font-semibold text-emerald-400 mb-4">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          AI-Powered Technical SEO Engine
        </div>
        <h1 className="text-4xl sm:text-5xl font-extrabold tracking-tight bg-gradient-to-r from-emerald-400 to-sky-400 bg-clip-text text-transparent mb-3">
          SEO Auditor & Analyzer
        </h1>
        <p className="text-slate-400 max-w-xl mx-auto text-sm sm:text-base">
          Submit any website URL to perform a full technical DOM audit, extract core web signals, and generate AI-driven action plans.
        </p>
      </header>

      <main className="space-y-10">
        <UrlForm onSubmit={handleAuditSubmit} loading={loading} />

        {/* Recent History Quick Bar */}
        {history.length > 0 && !loading && (
          <div className="flex items-center gap-2 flex-wrap justify-center text-xs text-slate-400">
            <span className="font-semibold text-slate-500">Recent:</span>
            {history.map((h, idx) => (
              <button
                key={idx}
                onClick={() => setData(h.data)}
                className="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded-md text-slate-300 transition-colors flex items-center gap-1.5"
              >
                <span className="truncate max-w-[140px]">{h.url}</span>
                <span className={`font-bold ${h.score >= 70 ? 'text-emerald-400' : h.score >= 50 ? 'text-amber-400' : 'text-rose-400'}`}>
                  {h.score}
                </span>
              </button>
            ))}
          </div>
        )}

        {/* Loading Spinner */}
        {loading && (
          <div className="flex flex-col items-center justify-center py-16 space-y-4 bg-slate-800/30 rounded-2xl border border-slate-800">
            <svg className="animate-spin h-10 w-10 text-emerald-400" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            <div className="text-center space-y-1">
              <p className="text-slate-200 font-semibold text-sm">
                Scraping DOM & Generating AI SEO Report...
              </p>
              <p className="text-slate-500 text-xs">
                Rendering page in headless Chromium and analyzing with Google Gemini
              </p>
            </div>
          </div>
        )}

        {/* Error Alert */}
        {error && (
          <div className="w-full max-w-3xl mx-auto bg-rose-500/10 border border-rose-500/30 p-5 rounded-xl flex items-start gap-3 shadow-lg">
            <svg className="w-6 h-6 text-rose-400 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div className="space-y-1">
              <h4 className="font-bold text-rose-400 text-sm">Audit Request Failed</h4>
              <p className="text-xs text-rose-300 break-words leading-relaxed">{error}</p>
            </div>
          </div>
        )}

        {/* Report Results */}
        {data && data.report && (
          <div className="space-y-8 animate-fadeIn">
            {/* Header / Actions Bar */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-slate-800/80 border border-slate-700/60 p-4 rounded-xl">
              <div>
                <span className="text-xs text-slate-400 block font-medium">Audited Target:</span>
                <a
                  href={data.audit?.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm font-bold text-emerald-400 hover:underline break-all"
                >
                  {data.audit?.url}
                </a>
              </div>
              <button
                onClick={handleExportJSON}
                className="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-200 text-xs font-semibold rounded-lg transition-colors flex items-center justify-center gap-1.5 self-start sm:self-auto"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                Export JSON Report
              </button>
            </div>

            {/* Score & Summary */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 items-stretch">
              <ScoreGauge score={data.report.score} />

              <div className="md:col-span-2 bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg flex flex-col justify-center">
                <h3 className="text-base font-bold text-slate-100 mb-2 flex items-center gap-2">
                  <svg className="w-5 h-5 text-sky-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  Executive SEO Summary
                </h3>
                <p className="text-sm leading-relaxed text-slate-300">
                  {data.report.summary}
                </p>
              </div>
            </div>

            {/* Quick Wins & Strengths */}
            <ReportCards
              quickWins={data.report.quick_wins}
              strengths={data.report.strengths}
            />

            {/* Detailed Issues Table */}
            <IssuesTable issues={data.report.issues} />

            {/* Raw Technical Accordion */}
            <AuditAccordion audit={data.audit} />
          </div>
        )}
      </main>

      <footer className="mt-20 pt-8 border-t border-slate-800 text-center text-xs text-slate-500">
        &copy; {new Date().getFullYear()} SEO Audit Tool &bull; Powered by Google Gemini & Go-Rod Headless Engine
      </footer>
    </div>
  )
}
