import { Table, Button, Field, TextInput, Badge, Skeleton, EmptyState, Dialog, Selector } from '../../components/Astryx'
import { useState, useEffect } from 'react'
import api from '../../services/api'

export default function Users(){
  const [filter, setFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<any[]>([])
  const [selected, setSelected] = useState<any>(null)

  useEffect(()=>{
    setLoading(true)
    api.get('/users', { params: { search: filter }}).then(r=>setData(r.data.users || [
      { id:'1', username:'budi123', email:'budi@memekita.id', role:'user', status:'active' },
      { id:'2', username:'admin', email:'admin@memekita.id', role:'admin', status:'active' },
    ])).catch(()=>setData([
      { id:'1', username:'budi123', email:'budi@memekita.id', role:'user', status:'active' },
    ])).finally(()=>setLoading(false))
  },[filter])

  const cols = [
    { key:'username', title:'Username' },
    { key:'email', title:'Email' },
    { key:'role', title:'Role', render:(v:any)=><Badge>{v}</Badge> },
    { key:'status', title:'Status' },
    { key:'aksi', title:'Aksi', render:(_:any,row:any)=><><Button onClick={()=>setSelected(row)}>Lihat</Button> <Button variant="danger" onClick={()=>alert('Blokir '+row.username)}>Blokir</Button></> },
  ]

  if (loading) return <><h1>Pengguna</h1><Skeleton /></>
  if (data.length===0) return <EmptyState icon="person" title="Tidak ada pengguna" desc="Coba ubah filter" action={<Button onClick={()=>setFilter('')}>Reset</Button>} />

  return <div><h1>Pengguna</h1>
    <div style={{ marginBottom:12, display:'flex', gap:8 }}><TextInput placeholder="Cari user..." value={filter} onChange={(e:any)=>setFilter(e.target.value)} /><Button>Cari</Button></div>
    <Table columns={cols} data={data.filter(d=>d.username.includes(filter) || d.email.includes(filter))} onRowClick={setSelected} />
    <Dialog open={!!selected} onClose={()=>setSelected(null)} title="Detail Pengguna">
      {selected && <>
        <p><strong>{selected.username}</strong> — {selected.email}</p>
        <Field label="Role"><Selector defaultValue={selected.role}><option value="user">user</option><option value="moderator">moderator</option><option value="admin">admin</option></Selector></Field>
        <Field label="Tindakan"><div style={{ display:'flex', gap:8 }}><Button>Simpan Role</Button><Button variant="danger">Hapus Akun</Button></div></Field>
      </>}
    </Dialog>
  </div>
}
