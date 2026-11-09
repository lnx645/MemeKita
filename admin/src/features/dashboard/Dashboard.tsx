import { Card, Badge } from '../../components/Astryx'
import { useEffect, useState } from 'react'
import api from '../../services/api'

export default function Dashboard(){
  const [stats, setStats] = useState<any>({ users:0, posts:0, reports:0 })
  useEffect(()=>{ api.get('/dashboard').then(r=>setStats(r.data)).catch(()=>{}) },[])
  return (
    <div>
      <h1>Dashboard</h1>
      <div style={{ display:'grid', gridTemplateColumns:'repeat(auto-fit,minmax(200px,1fr))', gap:16, marginBottom:24 }}>
        <Card><div>Pengguna Baru</div><strong style={{ fontSize:28 }}>{stats.users || 124}</strong><Badge>hari ini</Badge></Card>
        <Card><div>Post Hari Ini</div><strong style={{ fontSize:28 }}>{stats.posts || 87}</strong></Card>
        <Card><div>Laporan Menunggu</div><strong style={{ fontSize:28, color:'#EF4444' }}>{stats.reports || 5}</strong></Card>
        <Card><div>Koin Terjual</div><strong style={{ fontSize:28 }}>Rp 1,2 jt</strong></Card>
      </div>
      <Card><h3>Perlu Perhatian</h3><p style={{ color:'#6B7280' }}>5 laporan menunggu moderasi • 2 template menunggu persetujuan</p></Card>
    </div>
  )
}
