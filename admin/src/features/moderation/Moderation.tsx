import { Table, Button, Dialog, Field, Selector, Badge, Tabs, Skeleton, EmptyState } from '../../components/Astryx'
import { useState, useEffect } from 'react'
import api from '../../services/api'

export default function Moderation(){
  const [tab, setTab] = useState('pending')
  const [selected, setSelected] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any[]>([])

  useEffect(()=>{
    setLoading(true)
    api.get('/reports', { params: { status: tab }}).then(r=>setData(r.data.reports || [])).catch(()=>setData(tab==='pending' ? [{ id:'1', reporter:'andi', target:'post #123', reason:'spam', status:'pending', preview:'Meme spam', reporterName:'andi', targetPreview:'Konten spam' }] : [])).finally(()=>setLoading(false))
  },[tab])

  const cols = [
    { key:'reporter', title:'Pelapor' },
    { key:'target', title:'Konten' },
    { key:'reason', title:'Alasan', render:(v:any)=><Badge>{v}</Badge> },
    { key:'status', title:'Status' },
    { key:'aksi', title:'Aksi', render:(_:any,row:any)=><Button onClick={()=>setSelected(row)}>Tinjau</Button> },
  ]

  if (loading) return <><h1>Moderasi</h1><Skeleton /></>

  return <div><h1>Moderasi</h1>
    <Tabs tabs={[{key:'pending',label:'Menunggu'},{key:'handled',label:'Selesai'}]} active={tab} onChange={setTab} />
    {data.length===0 ? <EmptyState icon="flag" title="Tidak ada laporan" desc="Antrean kosong" /> : <Table columns={cols} data={data} />}
    <Dialog open={!!selected} onClose={()=>setSelected(null)} title="Putusan">
      {selected && <>
        <p><strong>Konten:</strong> {selected.preview || selected.targetPreview}</p>
        <p><strong>Pelapor:</strong> {selected.reporter} — <strong>Alasan:</strong> {selected.reason}</p>
        <Field label="Tindakan"><Selector><option>Abaikan</option><option>Sembunyikan</option><option>Hapus</option><option>Hapus + Blokir</option></Selector></Field>
        <Field label="Catatan"><textarea style={{ width:'100%', border:'1px solid #E8E6F0', borderRadius:8, padding:8 }} rows={3} placeholder="Alasan keputusan..." /></Field>
        <div style={{ display:'flex', gap:8 }}><Button onClick={()=>setSelected(null)}>Batal</Button><Button onClick={()=>{
          api.post(`/reports/${selected.id}/resolve`, { action:'hide', note:'moderasi' }).catch(()=>{})
          setSelected(null)
        }}>Simpan Keputusan</Button></div>
      </>}
    </Dialog>
  </div>
}
