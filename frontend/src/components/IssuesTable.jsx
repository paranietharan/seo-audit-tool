import React from 'react'

export default function IssuesTable({ issues }) {
  if (!issues || issues.length === 0) {
    return (
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 text-center text-slate-400">
        No issues detected! Excellent job.
      </div>
    )
  }

  const sortedIssues = [...issues].sort((a, b) => {
    const severityOrder = { critical: 0, warning: 1, info: 2 }
    const aVal = severityOrder[a.severity.toLowerCase()] ?? 99
    const bVal = severityOrder[b.severity.toLowerCase()] ?? 99
    return aVal - bVal
  })

  const getBadgeClass = (sev) => {
    switch (sev.toLowerCase()) {
      case 'critical':
        return 'bg-rose-500/20 text-rose-400 border-rose-500/30'
      case 'warning':
        return 'bg-amber-500/20 text-amber-400 border-amber-500/30'
      default:
        return 'bg-sky-500/20 text-sky-400 border-sky-500/30'
    }
  }

  return (
    <div className="bg-slate-800 border border-slate-700 rounded-xl overflow-hidden shadow-lg">
      <div className="p-5 border-b border-slate-700">
        <h3 className="text-lg font-semibold text-slate-100 flex items-center gap-2">
          <span>Detailed SEO Issues</span>
          <span className="bg-slate-700 text-xs px-2.5 py-1 rounded-full text-slate-300">
            {issues.length}
          </span>
        </h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm text-slate-300">
          <thead className="bg-slate-800/50 text-slate-400 uppercase text-xs tracking-wider border-b border-slate-700">
            <tr>
              <th className="px-6 py-3.5 font-semibold">Severity</th>
              <th className="px-6 py-3.5 font-semibold">Category</th>
              <th className="px-6 py-3.5 font-semibold">Description</th>
              <th className="px-6 py-3.5 font-semibold">Recommendation</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-700/50">
            {sortedIssues.map((issue, idx) => (
              <tr key={idx} className="hover:bg-slate-700/30 transition-colors">
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`inline-flex items-center px-2.5 py-1 rounded-md text-xs font-semibold border ${getBadgeClass(issue.severity)}`}>
                    {issue.severity.toUpperCase()}
                  </span>
                </td>
                <td className="px-6 py-4 font-semibold text-slate-200">
                  {issue.category}
                </td>
                <td className="px-6 py-4 max-w-xs break-words">
                  {issue.description}
                </td>
                <td className="px-6 py-4 text-emerald-450 max-w-sm break-words">
                  {issue.recommendation}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
