import { Card, Table, Button } from '../../components/Astryx'
export default function Comments(){ return <div><h1>Komentar</h1><Card><Table columns={[{key:'body',title:'Isi'},{key:'aksi',title:'',render:()=><Button variant="danger">Hapus</Button>}]} data={[{body:'Komentar contoh'}]} /></Card></div> }
