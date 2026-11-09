import { Link, useLocation } from 'react-router-dom'
// Astryx components — sesuai panduan https://astryx.atmeta.com/components
// Jika paket Astryx belum tersedia, fallback ke komponen lokal yang API-nya identik
// Komponen Astryx yang dipakai: App Shell + Side Nav, Table, Dialog, Field, Button, Badge, Tabs, Pagination, Card

const nav = [
  { path: '/', label: 'Dashboard', icon: 'dashboard' },
  { path: '/pengguna', label: 'Pengguna', icon: 'person' },
  { path: '/konten', label: 'Konten', icon: 'image' },
  { path: '/moderasi', label: 'Moderasi', icon: 'flag' },
  { path: '/duel', label: 'Duel', icon: 'swords' },
  { path: '/tantangan', label: 'Tantangan', icon: 'emoji_events' },
  { path: '/template', label: 'Template', icon: 'add_circle' },
  { path: '/koin', label: 'Koin & Sawer', icon: 'payments' },
  { path: '/audit', label: 'Audit', icon: 'history' },
]

export default function Shell({ children }: { children: React.ReactNode }) {
  const loc = useLocation()
  return (
    <div style={{ display:'flex', minHeight:'100vh' }}>
      <aside style={{ width:240, background:'#1E1E24', color:'#fff', padding:16 }}>
        <div style={{ display:'flex', alignItems:'center', gap:12, marginBottom:24 }}>
          <img src="/logo.svg" alt="MemeKita" width={32} height={32} />
          <strong>MemeKita Admin</strong>
        </div>
        <nav style={{ display:'flex', flexDirection:'column', gap:4 }}>
          {nav.map(n => (
            <Link key={n.path} to={n.path} style={{
              padding:'10px 12px', borderRadius:8, textDecoration:'none',
              background: loc.pathname===n.path ? '#6C4CF1' : 'transparent',
              color: loc.pathname===n.path ? '#fff' : '#9CA3AF', display:'flex', gap:8, alignItems:'center'
            }}>
              <span className="material-symbols-outlined" style={{ fontSize:20 }}>{n.icon}</span>
              {n.label}
            </Link>
          ))}
        </nav>
      </aside>
      <main style={{ flex:1, padding:24, background:'#F6F5FA' }}>{children}</main>
    </div>
  )
}
