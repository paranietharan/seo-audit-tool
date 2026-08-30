import React, { useState } from 'react'

export default function AuditAccordion({ audit }) {
  const [isOpen, setIsOpen] = useState(false)
  const [showInternal, setShowInternal] = useState(false)
  const [showExternal, setShowExternal] = useState(false)

  if (!audit) return null

  return (
    <div className="bg-slate-800 border border-slate-700 rounded-xl shadow-lg overflow-hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full px-6 py-4 flex items-center justify-between text-left font-bold text-slate-100 bg-slate-800/50 hover:bg-slate-700/50 transition-colors focus:outline-none"
      >
        <span className="flex items-center gap-2">
          <svg className="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
          </svg>
          Raw Technical Audit & DOM Extraction Data
        </span>
        <svg
          className={`w-5 h-5 text-slate-400 transform transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {isOpen && (
        <div className="p-6 border-t border-slate-700 bg-slate-900/40 space-y-6">
          {/* Core Metadata */}
          <div>
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Core Page Meta</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Target URL</span>
                <p className="text-sm font-semibold text-slate-200 break-all">{audit.url}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Page Title</span>
                <p className="text-sm font-semibold text-slate-200">{audit.title || '(missing)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Canonical URL</span>
                <p className="text-sm font-semibold text-slate-200 break-all">{audit.canonical_url || '(missing)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Meta Description</span>
                <p className="text-sm font-semibold text-slate-200">{audit.meta_description || '(missing)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Robots Meta</span>
                <p className="text-sm font-semibold text-slate-200">{audit.robots_meta || '(default index/follow)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Meta Viewport</span>
                <p className="text-sm font-semibold text-slate-200">{audit.viewport_meta || '(missing mobile viewport)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Structured Data</span>
                <p className="text-sm font-semibold text-slate-200">
                  {audit.has_structured_data ? (
                    <span className="text-emerald-400 flex items-center gap-1">
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M5 13l4 4L19 7" />
                      </svg>
                      Detected (JSON-LD)
                    </span>
                  ) : (
                    <span className="text-slate-400">None detected</span>
                  )}
                </p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Status Code</span>
                <p className="text-sm font-semibold text-slate-200">{audit.status_code || 200}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Render Time</span>
                <p className="text-sm font-semibold text-slate-200">{audit.page_load_ms || 0} ms</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs uppercase tracking-wider text-slate-400 font-medium">Word Count</span>
                <p className="text-sm font-semibold text-slate-200">{audit.word_count || 0} words</p>
              </div>
            </div>
          </div>

          <hr className="border-slate-700/60" />

          {/* Headings */}
          <div>
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Headings Distribution</h4>
            <div className="grid grid-cols-3 sm:grid-cols-6 gap-3">
              {['H1', 'H2', 'H3', 'H4', 'H5', 'H6'].map((tag, i) => {
                const count = audit[`h${i+1}_count`] ?? 0
                return (
                  <div key={tag} className="bg-slate-800 border border-slate-700/50 px-4 py-2.5 rounded-lg text-center">
                    <span className="text-xs text-slate-400 block font-semibold">{tag}</span>
                    <span className="text-lg font-bold text-slate-200">{count}</span>
                  </div>
                )
              })}
            </div>
            {audit.heading_issues && audit.heading_issues.length > 0 && (
              <div className="mt-3 bg-amber-500/10 border border-amber-500/30 p-3 rounded-lg text-xs text-amber-300">
                {audit.heading_issues.join(', ')}
              </div>
            )}
          </div>

          <hr className="border-slate-700/60" />

          {/* Social Graph Tags */}
          <div>
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Social Graph Meta (Open Graph & Twitter)</h4>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
              <div className="space-y-1">
                <span className="text-xs text-slate-400 block font-semibold">og:title</span>
                <p className="text-sm font-semibold text-slate-300">{audit.og_title || '(none)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs text-slate-400 block font-semibold">og:description</span>
                <p className="text-sm font-semibold text-slate-300">{audit.og_description || '(none)'}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs text-slate-400 block font-semibold">og:image</span>
                <p className="text-sm font-semibold text-slate-300 break-all">{audit.og_image || '(none)'}</p>
              </div>
              {audit.twitter_card && (
                <div className="space-y-1">
                  <span className="text-xs text-slate-400 block font-semibold">twitter:card</span>
                  <p className="text-sm font-semibold text-slate-300">{audit.twitter_card}</p>
                </div>
              )}
            </div>
          </div>

          <hr className="border-slate-700/60" />

          {/* Images Missing Alt */}
          <div>
            <h4 className="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">
              Images Missing Alt Attributes ({audit.images_missing_alt?.length || 0})
            </h4>
            {(!audit.images_missing_alt || audit.images_missing_alt.length === 0) ? (
              <p className="text-sm text-emerald-400 flex items-center gap-1.5">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                </svg>
                All images have alternative descriptions!
              </p>
            ) : (
              <div className="max-h-40 overflow-y-auto bg-slate-800 border border-slate-700/50 rounded-lg p-3 space-y-2">
                {audit.images_missing_alt.map((src, i) => (
                  <div key={i} className="text-xs text-rose-400 font-mono truncate hover:bg-slate-700 p-1 rounded transition-colors" title={src}>
                    {src}
                  </div>
                ))}
              </div>
            )}
          </div>

          <hr className="border-slate-700/60" />

          {/* Links Section */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <button
                onClick={() => setShowInternal(!showInternal)}
                className="w-full flex items-center justify-between text-left text-sm font-bold text-slate-300 bg-slate-800 border border-slate-700 px-4 py-2.5 rounded-lg hover:bg-slate-700 transition-colors focus:outline-none"
              >
                <span>Internal Links ({audit.internal_links?.length || 0})</span>
                <svg className={`w-4 h-4 transform transition-transform ${showInternal ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              {showInternal && audit.internal_links && (
                <div className="max-h-60 overflow-y-auto bg-slate-900 border border-slate-700/50 border-t-0 rounded-b-lg p-3 space-y-3">
                  {audit.internal_links.map((link, i) => (
                    <div key={i} className="text-xs space-y-0.5 border-b border-slate-800 pb-2 last:border-b-0">
                      <span className="text-slate-400 block font-semibold">Anchor: "{link.text || '(empty text)'}"</span>
                      <a href={link.href} target="_blank" rel="noopener noreferrer" className="text-emerald-400 hover:underline break-all block">
                        {link.href}
                      </a>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div>
              <button
                onClick={() => setShowExternal(!showExternal)}
                className="w-full flex items-center justify-between text-left text-sm font-bold text-slate-300 bg-slate-800 border border-slate-700 px-4 py-2.5 rounded-lg hover:bg-slate-700 transition-colors focus:outline-none"
              >
                <span>External Outbound Links ({audit.external_links?.length || 0})</span>
                <svg className={`w-4 h-4 transform transition-transform ${showExternal ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              {showExternal && audit.external_links && (
                <div className="max-h-60 overflow-y-auto bg-slate-900 border border-slate-700/50 border-t-0 rounded-b-lg p-3 space-y-3">
                  {audit.external_links.map((link, i) => (
                    <div key={i} className="text-xs space-y-0.5 border-b border-slate-800 pb-2 last:border-b-0">
                      <span className="text-slate-400 block font-semibold">Anchor: "{link.text || '(empty text)'}"</span>
                      <a href={link.href} target="_blank" rel="noopener noreferrer" className="text-sky-400 hover:underline break-all block">
                        {link.href}
                      </a>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
