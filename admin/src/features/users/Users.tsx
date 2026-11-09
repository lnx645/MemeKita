import { Table, Button, Dialog, Field, TextInput, Selector, Badge } from '../../components/Astryx'
import { useState } from 'react'

export default function Users(){
  const [filter, setFilter] = useState('')
  const data = [
    { id:'1', username:'budi123', email:'budi@memekita.id', role:'user', status:'active' },
    { id:'2', username:'admin', email:'admin@memekita.id', role:'admin', status:'active' },
  ]
  const cols = [
    { key:'username', title:'Username' },
    { key:'email', title:'Email' },
    { key:'role', title:'Role', render:(v:any)=><Badge>{v}</Badge> },
    { key:'status', title:'Status' },
    { key:'aksi', title:'Aksi', render:(_:any,row:any)=><Button onClick={()=>alert('Blokir '+row.username)}>Blokir</Button> },
  ]
  return <div><h1>Pengguna</h1>
    <div style={{ marginBottom:12, display:'flex', gap:8 }}><TextInput placeholder="Cari user..." value={filter} onChange={(e:any)=>setFilter(e.target.value)} /><Button>Cari</Button></div>
    <Table columns={cols} data={data.filter(d=>d.username.includes(filter))} />
  </div>
}
