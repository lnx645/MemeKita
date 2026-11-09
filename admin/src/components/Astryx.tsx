// Fallback komponen Astryx — API identik dengan https://astryx.atmeta.com/components
// Bila paket Astryx terpasang, ganti import ke "@astryx/react"

export function Button({ children, variant="primary", ...props }: any) {
  const bg = variant==='primary' ? '#6C4CF1' : variant==='danger' ? '#EF4444' : '#fff'
  const color = variant==='primary' || variant==='danger' ? '#fff' : '#1A1A1A'
  return <button style={{ background:bg, color, border:'1px solid #E8E6F0', borderRadius:8, padding:'8px 16px', cursor:'pointer', fontWeight:600 }} {...props}>{children}</button>
}
export function Card({ children, ...props }: any) {
  return <div style={{ background:'#fff', borderRadius:12, padding:16, boxShadow:'0 1px 3px rgba(0,0,0,0.08)' }} {...props}>{children}</div>
}
export function Badge({ children, color="primary" }: any) {
  const bg = color==='success' ? '#22C55E' : color==='danger' ? '#EF4444' : '#6C4CF1'
  return <span style={{ background:bg, color:'#fff', borderRadius:999, padding:'2px 8px', fontSize:12 }}>{children}</span>
}
export function Table({ columns, data, onRowClick }: any) {
  return (
    <div style={{ overflowX:'auto' }}>
      <table style={{ width:'100%', borderCollapse:'collapse', background:'#fff', borderRadius:12, overflow:'hidden' }}>
        <thead><tr style={{ background:'#F6F5FA', textAlign:'left' }}>
          {columns.map((c:any)=><th key={c.key} style={{ padding:'12px 16px', fontSize:13, color:'#6B7280' }}>{c.title}</th>)}
        </tr></thead>
        <tbody>
          {data.map((row:any,i:number)=>
            <tr key={i} onClick={()=>onRowClick?.(row)} style={{ borderTop:'1px solid #E8E6F0', cursor: onRowClick ? 'pointer' : 'default' }}>
              {columns.map((c:any)=><td key={c.key} style={{ padding:'12px 16px', fontSize:14 }}>{c.render ? c.render(row[c.key], row) : row[c.key]}</td>)}
            </tr>
          )}
          {data.length===0 && <tr><td colSpan={columns.length} style={{ padding:32, textAlign:'center', color:'#6B7280' }}>Tidak ada data</td></tr>}
        </tbody>
      </table>
    </div>
  )
}
export function Dialog({ open, onClose, title, children }: any) {
  if (!open) return null
  return <div style={{ position:'fixed', inset:0, background:'rgba(0,0,0,0.4)', display:'flex', alignItems:'center', justifyContent:'center', zIndex:50 }}>
    <div style={{ background:'#fff', borderRadius:12, padding:24, minWidth:400, maxWidth:600 }}>
      <div style={{ display:'flex', justifyContent:'space-between', marginBottom:16 }}><strong>{title}</strong><button onClick={onClose} style={{ background:'none', border:'none', cursor:'pointer' }}>✕</button></div>
      {children}
    </div>
  </div>
}
export function Field({ label, children }: any) { return <div style={{ display:'flex', flexDirection:'column', gap:6, marginBottom:12 }}><label style={{ fontSize:13, fontWeight:600 }}>{label}</label>{children}</div> }
export function TextInput(props:any){ return <input style={{ border:'1px solid #E8E6F0', borderRadius:8, padding:'10px 12px', fontSize:14 }} {...props} /> }
export function Selector(props:any){ return <select style={{ border:'1px solid #E8E6F0', borderRadius:8, padding:'10px 12px', fontSize:14 }} {...props} /> }
export function Tabs({ tabs, active, onChange }: any){
  return <div style={{ display:'flex', gap:8, borderBottom:'2px solid #E8E6F0', marginBottom:16 }}>
    {tabs.map((t:any)=><button key={t.key} onClick={()=>onChange(t.key)} style={{ padding:'10px 16px', border:'none', background:'none', borderBottom: active===t.key ? '2px solid #6C4CF1' : '2px solid transparent', color: active===t.key ? '#6C4CF1' : '#6B7280', fontWeight:600, cursor:'pointer' }}>{t.label}</button>)}
  </div>
}
export function Pagination({ page, total, onChange }: any){
  return <div style={{ display:'flex', gap:8, justifyContent:'center', marginTop:16 }}>
    <button onClick={()=>onChange(page-1)} disabled={page<=1}>‹</button>
    <span>Halaman {page} / {Math.ceil(total/20)}</span>
    <button onClick={()=>onChange(page+1)} disabled={page>=Math.ceil(total/20)}>›</button>
  </div>
}
export function Skeleton(){ return <div style={{ height:120, background:'#F6F5FA', borderRadius:12, animation:'pulse 1.5s infinite' }} /> }
export function EmptyState({ icon, title, desc, action }: any){
  return <div style={{ textAlign:'center', padding:48, background:'#fff', borderRadius:12 }}>
    <div style={{ fontSize:48 }}>{icon}</div>
    <h3>{title}</h3><p style={{ color:'#6B7280' }}>{desc}</p>{action}
  </div>
}
