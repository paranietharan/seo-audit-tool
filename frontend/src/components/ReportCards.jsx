import React, { useState } from 'react'

export default function ReportCards({ quickWins, strengths, onExportJSON }) {
  const [copied, setCopied] = useState(false)

  const handleCopySummary = (items, title) => {
    if (!items || items.length === 0) return
    const text = `${title}:\n` + items.map(it => `• ${it}`).join('\n')
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      {/* Quick Wins */}
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg flex flex-col justify-between">
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-bold text-amber-400 flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              High-Impact Quick Wins
            </h3>
            {quickWins && quickWins.length > 0 && (
              <button
                onClick={() => handleCopySummary(quickWins, 'Quick Wins')}
                className="text-xs text-slate-400 hover:text-slate-200 transition-colors"
                title="Copy Quick Wins"
              >
                {copied ? 'Copied!' : 'Copy'}
              </button>
            )}
          </div>

          {(!quickWins || quickWins.length === 0) ? (
            <p className="text-slate-400 text-sm">No immediate quick wins identified.</p>
          ) : (
            <ul className="space-y-3">
              {quickWins.map((win, idx) => (
                <li key={idx} className="flex items-start gap-2.5 text-sm text-slate-300">
                  <svg className="w-4 h-4 text-amber-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M9 5l7 7-7 7" />
                  </svg>
                  <span>{win}</span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      {/* Strengths */}
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 shadow-lg flex flex-col justify-between">
        <div>
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-bold text-emerald-400 flex items-center gap-2">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              SEO Strengths & Best Practices
            </h3>
          </div>

          {(!strengths || strengths.length === 0) ? (
            <p className="text-slate-400 text-sm">No specific strengths highlighted.</p>
          ) : (
            <ul className="space-y-3">
              {strengths.map((strength, idx) => (
                <li key={idx} className="flex items-start gap-2.5 text-sm text-slate-300">
                  <svg className="w-4 h-4 text-emerald-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M5 13l4 4L19 7" />
                  </svg>
                  <span>{strength}</span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}
