import { Card, Table } from '../../components/Astryx'
export default function Audit(){ return <div><h1>Audit Log</h1><Card><Table columns={[{key:'actor',title:'Aktor'},{key:'action',title:'Aksi'},{key:'target',title:'Target'},{key:'created_at',title:'Waktu'}]} data={[{actor:'admin', action:'moderation:hide', target:'post #123', created_at:'2026-11-09'}]} /></Card></div> }
