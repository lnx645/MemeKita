import { Card, Field, TextInput, Button } from '../../components/Astryx'
export default function SettingsAdmin(){ return <div><h1>Pengaturan Aplikasi</h1><Card><Field label="Batas Ukuran Unggahan (MB)"><TextInput defaultValue="10" /></Field><Field label="Mode Pemeliharaan"><input type="checkbox" /> Aktif</Field><Button>Simpan</Button></Card></div> }
