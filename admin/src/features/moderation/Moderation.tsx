import { Table, Button, Dialog, Field, Selector, Badge, Tabs } from '../../components/Astryx'
import { useState } from 'react'

export default function Moderation(){
  const [tab, setTab] = useState('pending')
  const [selected, setSelected] = useState<any>(null)
  const data = tab==='pending' ? [
    { id:'1', reporter:'andi', target:'post #123', reason:'spam', status:'pending', preview:'Meme spam' }
  ] : []
  const cols = [
    { key:'reporter', title:'Pelapor' },
    { key:'target', title:'Konten' },
    { key:'reason', title:'Alasan', render:(v:any)=><Badge>{v}</Badge> },
    { key:'status', title:'Status' },
    { key:'aksi', title:'Aksi', render:(_:any,row:any)=><Button onClick={()=>setSelected(row)}>Tinjau</Button> },
  ]
  return <div><h1>Moderasi</h1>
    <Tabs tabs={[{key:'pending',label:'Menunggu'},{key:'handled',label:'Selesai'}]} active={tab} onChange={setTab} />
    <Table columns={cols} data={data} />
    <Dialog open={!!selected} onClose={()=>setSelected(null)} title="Putusan">
      {selected && <>
        <p>Konten: {selected.preview}</p>
        <Field label="Tindakan"><Selector><option>Abaikan</option><option>Sembunyikan</option><option>Hapus</option><option>Hapus + Blokir</option></Selector></Field>
        <Field label="Catatan"><textarea style={{ width:'100%', border:'1px solid #E8E6F0', borderRadius:8, padding:8 }} rows={3} /></Field>
        <Button onClick={()=>setSelected(null)}>Simpan Keputusan</Button>
      </>}
    </Dialog>
  </div>
}
