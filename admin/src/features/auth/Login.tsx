import { useState } from 'react'
import { Button, Field, TextInput, Card } from '../../components/Astryx'
import api from '../../services/api'

export default function Login(){
  const [email, setEmail] = useState('admin@memekita.id')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const login = async ()=>{
    setLoading(true)
    setError('')
    try{
      const r = await api.post('/auth/login', { email, password })
      const token = r.data.access_token || r.data.accessToken || 'dev-token'
      localStorage.setItem('admin_token', token)
      location.reload()
    }catch(e:any){
      setError(e.response?.data?.message || 'Login gagal — periksa email/password')
    }finally{
      setLoading(false)
    }
  }

  return <div style={{ display:'flex', alignItems:'center', justifyContent:'center', minHeight:'100vh', background:'#F6F5FA', padding:16 }}>
    <Card style={{ width:400, maxWidth:'100%' }}>
      <div style={{ display:'flex', alignItems:'center', gap:12, marginBottom:16 }}>
        <img src="/logo.svg" alt="MemeKita" width={32} height={32} />
        <h2 style={{ margin:0 }}>MemeKita Admin</h2>
      </div>
      <p style={{ color:'#6B7280', fontSize:14, marginTop:0 }}>Masuk untuk mengelola seluruh siklus aplikasi</p>
      <Field label="Email">
        <TextInput value={email} onChange={(e:any)=>setEmail(e.target.value)} placeholder="admin@memekita.id" />
      </Field>
      <Field label="Password">
        <TextInput type="password" value={password} onChange={(e:any)=>setPassword(e.target.value)} placeholder="••••••••" />
      </Field>
      {error && <div style={{ color:'#EF4444', fontSize:13, marginBottom:12 }}>{error}</div>}
      <Button onClick={login} style={{ width:'100%' }} disabled={loading}>{loading ? 'Memproses...' : 'Masuk'}</Button>
      <p style={{ fontSize:12, color:'#6B7280', marginTop:12, textAlign:'center' }}>Demo: admin@memekita.id / admin12345</p>
    </Card>
  </div>
}
