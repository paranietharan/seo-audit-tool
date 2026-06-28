import React, { useState } from 'react'
import UrlForm from './components/UrlForm'
import ScoreGauge from './components/ScoreGauge'
import IssuesTable from './components/IssuesTable'
import ReportCards from './components/ReportCards'
import AuditAccordion from './components/AuditAccordion'

export default function App() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [data, setData] = useState(null)

  const handleAuditSubmit = async (url) => {
    setLoading(true)
    setError(null)
    setData(null)

    try {
      const response = await fetch('http://localhost:8080/api/report', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      })

      if (!response.ok) {
        const errText = await response.text()
        throw new Error(errText || 'Failed to generate report')
      }

      const resData = await response.json()
      setData(resData)
    } catch (err) {
      setError(err.message || 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-12">
      <header className="text-center mb-12">
        <h1 className="text-4xl font-extrabold tracking-tight bg-gradient-to-r from-emerald-405 to-sky-405 bg-clip-text text-transparent mb-3">
          SEO Auditor & Analyzer
        </h1>
        <p className="text-slate-400 max-w-xl mx-auto text-sm sm:text-base">
          Submit any website URL to perform a full technical SEO audit and receive AI-generated recommendations instantly.
        </p>
      </header>

      <main className="space-y-10">
        <UrlForm onSubmit={handleAuditSubmit} loading={loading} />

        {loading && (
          <div className="flex flex-col items-center justify-center py-16 space-y-4">
            <svg className="animate-spin h-10 w-10 text-emerald-500" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            <p className="text-slate-400 font-medium text-sm animate-pulse">
              Scraping URL and generating report, please wait... (up to 60s)
            </p>
          </div>
        )}

        {error && (
          <div className="w-full max-w-3xl mx-auto bg-rose-500/10 border border-rose-500/25 p-5 rounded-xl flex items-start gap-3">
            <svg className="w-6 h-6 text-rose-500 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div className="space-y-1">
              <h4 className="font-bold text-rose-400 text-sm">Audit Execution Failed</h4>
              <p className="text-xs text-rose-300 break-words">{error}</p>
            </div>
          </div>
        )}

        {data && data.report && (
          <div className="space-y-10">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 items-stretch">
              <ScoreGauge score={data.report.score} />
              
              <div className="md:col-span-2 bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg flex flex-col justify-center">
                <h3 className="text-lg font-bold text-slate-100 mb-2">Executive Summary</h3>
                <p className="text-sm leading-relaxed text-slate-350">
                  {data.report.summary}
                </p>
              </div>
            </div>

            <ReportCards
              quickWins={data.report.quick_wins}
              strengths={data.report.strengths}
            />

            <IssuesTable issues={data.report.issues} />

            <AuditAccordion audit={data.audit} />
          </div>
        )}
      </main>

      <footer className="mt-20 pt-8 border-t border-slate-800 text-center text-xs text-slate-500">
        &copy; {new Date().getFullYear()} SEO Audit Tool. Powered by Gemini API & Go-Rod.
      </footer>
    </div>
  )
}
