import { Card, Table, Button } from '../../components/Astryx'
export default function Challenges(){ return <div><h1>Tantangan</h1><Card><Table columns={[{key:'theme',title:'Tema'},{key:'status',title:'Status'},{key:'aksi',title:'',render:()=><Button>Edit</Button>}]} data={[{theme:'macet', status:'open'}]} /><Button style={{marginTop:12}}>Buat Tema</Button></Card></div> }
