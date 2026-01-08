import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Dashboard from './pages/dashboard';
import GussPage from './pages/guss';
import LoginPage from './pages/login';
import RegisterPage from './pages/register';
import AdminPage from './pages/admin';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* 1. 서비스 탐색 (비로그인 가능) */}
        <Route path="/" element={<Dashboard />} />
        <Route path="/guss" element={<GussPage />} />
        
        {/* 2. 인증 프로세스 */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        
        {/* 3. 관리자 전용 */}
        <Route path="/admin" element={<AdminPage />} />
      </Routes>
    </BrowserRouter>
  );
}