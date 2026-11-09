import { Table, Button, Dialog, Field, TextInput, Card } from '../../components/Astryx'
import { useState } from 'react'

export default function Coins(){
  const [open, setOpen] = useState(false)
  const packages = [
    { id:'1', coins:100, price_idr:15000, play_sku:'coins_100', is_active:true },
    { id:'2', coins:500, price_idr:49000, play_sku:'coins_500', is_active:true },
  ]
  const ledger = [
    { id:'1', user:'budi123', type:'topup', amount:500, balance_after:1500, created_at:'2026-11-09' },
  ]
  return <div><h1>Koin & Sawer</h1>
    <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:16 }}>
      <Card><h3>Paket Koin</h3>
        <Table columns={[
          {key:'coins',title:'Koin'},{key:'price_idr',title:'Harga'},{key:'play_sku',title:'SKU'},{key:'aksi',title:'',render:()=><Button>Edit</Button>}
        ]} data={packages} />
        <Button onClick={()=>setOpen(true)} style={{ marginTop:12 }}>Tambah Paket</Button>
      </Card>
      <Card><h3>Leaderboard Mingguan</h3><Table columns={[{key:'user',title:'Kreator'},{key:'total',title:'Koin Diterima'}]} data={[{user:'kreator1', total:3200},{user:'kreator2', total:1800}]} /></Card>
    </div>
    <Card style={{ marginTop:16 }}><h3>Ledger Transaksi</h3><Table columns={[{key:'user',title:'User'},{key:'type',title:'Tipe'},{key:'amount',title:'Jumlah'},{key:'balance_after',title:'Saldo Akhir'}]} data={ledger} /></Card>
    <Dialog open={open} onClose={()=>setOpen(false)} title="Tambah Paket">
      <Field label="Jumlah Koin"><TextInput type="number" /></Field>
      <Field label="Harga (IDR)"><TextInput type="number" /></Field>
      <Field label="Play SKU"><TextInput /></Field>
      <Button>Simpan</Button>
    </Dialog>
  </div>
}
