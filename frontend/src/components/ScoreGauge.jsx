import React from 'react'

export default function ScoreGauge({ score }) {
  const normalizedScore = Math.max(0, Math.min(100, score))
  const radius = 50
  const circumference = 2 * Math.PI * radius
  const offset = circumference - (normalizedScore / 100) * circumference

  const getColor = (s) => {
    if (s >= 70) return 'stroke-emerald-500'
    if (s >= 40) return 'stroke-amber-500'
    return 'stroke-rose-500'
  }

  const getTextColor = (s) => {
    if (s >= 70) return 'text-emerald-500'
    if (s >= 40) return 'text-amber-500'
    return 'text-rose-500'
  }

  const getBgColor = (s) => {
    if (s >= 70) return 'bg-emerald-500/10'
    if (s >= 40) return 'bg-amber-500/10'
    return 'bg-rose-500/10'
  }

  return (
    <div className="flex flex-col items-center justify-center p-6 bg-slate-800 rounded-xl border border-slate-700 shadow-lg">
      <div className={`relative flex items-center justify-center w-36 h-36 rounded-full ${getBgColor(normalizedScore)}`}>
        <svg className="w-full h-full transform -rotate-90" viewBox="0 0 120 120">
          <circle
            cx="60"
            cy="60"
            r={radius}
            className="stroke-slate-700"
            strokeWidth="8"
            fill="transparent"
          />
          <circle
            cx="60"
            cy="60"
            r={radius}
            className={`${getColor(normalizedScore)} transition-all duration-1000 ease-out`}
            strokeWidth="8"
            fill="transparent"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            strokeLinecap="round"
          />
        </svg>
        <div className="absolute text-center">
          <span className={`text-4xl font-extrabold ${getTextColor(normalizedScore)}`}>
            {normalizedScore}
          </span>
          <span className="block text-xs uppercase tracking-wider text-slate-400 font-medium mt-1">SEO Score</span>
        </div>
      </div>
    </div>
  )
}
