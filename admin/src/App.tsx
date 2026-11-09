import { Routes, Route, Navigate } from 'react-router-dom'
import Shell from './components/Shell'
import Dashboard from './features/dashboard/Dashboard'
import Users from './features/users/Users'
import Posts from './features/posts/Posts'
import Moderation from './features/moderation/Moderation'
import Duels from './features/duels/Duels'
import Challenges from './features/challenges/Challenges'
import Templates from './features/templates/Templates'
import Coins from './features/coins/Coins'
import Audit from './features/audit/Audit'
import Login from './features/auth/Login'

function App() {
  const token = localStorage.getItem('admin_token')
  if (!token) return <Login />
  return (
    <Shell>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/pengguna" element={<Users />} />
        <Route path="/konten" element={<Posts />} />
        <Route path="/moderasi" element={<Moderation />} />
        <Route path="/duel" element={<Duels />} />
        <Route path="/tantangan" element={<Challenges />} />
        <Route path="/template" element={<Templates />} />
        <Route path="/koin" element={<Coins />} />
        <Route path="/audit" element={<Audit />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Shell>
  )
}
export default App
