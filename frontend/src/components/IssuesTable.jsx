import React, { useState } from 'react'

export default function IssuesTable({ issues }) {
  const [filter, setFilter] = useState('all')

  if (!issues || issues.length === 0) {
    return (
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-8 text-center text-slate-400 shadow-lg">
        <svg className="w-12 h-12 text-emerald-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h4 className="text-base font-bold text-slate-200">No SEO Issues Detected</h4>
        <p className="text-sm text-slate-400 mt-1">The audited page passed all baseline automated checks.</p>
      </div>
    )
  }

  const severityOrder = { critical: 0, warning: 1, info: 2 }

  const sortedIssues = [...issues].sort((a, b) => {
    const aVal = severityOrder[(a?.severity || 'info').toLowerCase()] ?? 99
    const bVal = severityOrder[(b?.severity || 'info').toLowerCase()] ?? 99
    return aVal - bVal
  })

  const filteredIssues = sortedIssues.filter((issue) => {
    if (filter === 'all') return true
    return (issue?.severity || 'info').toLowerCase() === filter
  })

  const getBadgeClass = (sev) => {
    const normalized = (sev || 'info').toLowerCase()
    switch (normalized) {
      case 'critical':
        return 'bg-rose-500/20 text-rose-400 border-rose-500/30'
      case 'warning':
        return 'bg-amber-500/20 text-amber-400 border-amber-500/30'
      default:
        return 'bg-sky-500/20 text-sky-400 border-sky-500/30'
    }
  }

  const criticalCount = issues.filter(i => (i?.severity || '').toLowerCase() === 'critical').length
  const warningCount = issues.filter(i => (i?.severity || '').toLowerCase() === 'warning').length
  const infoCount = issues.filter(i => (i?.severity || '').toLowerCase() === 'info').length

  return (
    <div className="bg-slate-800 border border-slate-700 rounded-xl overflow-hidden shadow-lg">
      <div className="p-5 border-b border-slate-700 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h3 className="text-lg font-semibold text-slate-100 flex items-center gap-2">
          <span>Detailed SEO Issues & Recommendations</span>
          <span className="bg-slate-700 text-xs px-2.5 py-1 rounded-full text-slate-300">
            {issues.length}
          </span>
        </h3>

        <div className="flex items-center gap-1.5 bg-slate-900/60 p-1 rounded-lg border border-slate-700/50 self-start sm:self-auto">
          <button
            onClick={() => setFilter('all')}
            className={`px-3 py-1 text-xs font-semibold rounded-md transition-colors ${
              filter === 'all' ? 'bg-slate-700 text-white' : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            All ({issues.length})
          </button>
          {criticalCount > 0 && (
            <button
              onClick={() => setFilter('critical')}
              className={`px-3 py-1 text-xs font-semibold rounded-md transition-colors ${
                filter === 'critical' ? 'bg-rose-600 text-white' : 'text-rose-400 hover:bg-rose-950/40'
              }`}
            >
              Critical ({criticalCount})
            </button>
          )}
          {warningCount > 0 && (
            <button
              onClick={() => setFilter('warning')}
              className={`px-3 py-1 text-xs font-semibold rounded-md transition-colors ${
                filter === 'warning' ? 'bg-amber-600 text-white' : 'text-amber-400 hover:bg-amber-950/40'
              }`}
            >
              Warnings ({warningCount})
            </button>
          )}
          {infoCount > 0 && (
            <button
              onClick={() => setFilter('info')}
              className={`px-3 py-1 text-xs font-semibold rounded-md transition-colors ${
                filter === 'info' ? 'bg-sky-600 text-white' : 'text-sky-400 hover:bg-sky-950/40'
              }`}
            >
              Info ({infoCount})
            </button>
          )}
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm text-slate-300">
          <thead className="bg-slate-800/80 text-slate-400 uppercase text-xs tracking-wider border-b border-slate-700">
            <tr>
              <th className="px-6 py-3.5 font-semibold">Severity</th>
              <th className="px-6 py-3.5 font-semibold">Category</th>
              <th className="px-6 py-3.5 font-semibold">Issue Description</th>
              <th className="px-6 py-3.5 font-semibold">Actionable Fix</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700/50">
            {filteredIssues.map((issue, idx) => (
              <tr key={idx} className="hover:bg-slate-700/30 transition-colors">
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`inline-flex items-center px-2.5 py-1 rounded-md text-xs font-semibold border ${getBadgeClass(issue?.severity)}`}>
                    {(issue?.severity || 'info').toUpperCase()}
                  </span>
                </td>
                <td className="px-6 py-4 font-semibold text-slate-200">
                  {issue?.category || 'General'}
                </td>
                <td className="px-6 py-4 max-w-xs break-words text-slate-300">
                  {issue?.description || 'N/A'}
                </td>
                <td className="px-6 py-4 text-emerald-400 max-w-sm break-words">
                  {issue?.recommendation || 'N/A'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
