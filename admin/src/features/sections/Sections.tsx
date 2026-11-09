import { Card, Table, Button } from '../../components/Astryx'
export default function Sections(){ return <div><h1>Section & Kanal Daerah</h1><Card><Table columns={[{key:'name',title:'Nama'},{key:'slug',title:'Slug'},{key:'aksi',title:'',render:()=><Button>Edit</Button>}]} data={[{name:'Funny', slug:'funny'},{name:'Sunda', slug:'sunda'}]} /></Card></div> }
