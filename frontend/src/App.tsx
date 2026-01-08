import { useState } from 'react';
import RegisterPage from './pages/register';
import GussPage from './pages/guss';
import AdminPage from './pages/admin';

export default function App() {
  const [page, setPage] = useState('guss');

  return (
    <div className="relative">
      {/* 네비게이션 - 우측 상단 (작고 투명하게) */}
      <div style={{ 
        position: 'fixed', 
        top: 20, 
        right: 20, 
        zIndex: 9999,
        display: 'flex',
        gap: '8px',
        background: 'rgba(0, 0, 0, 0.5)',
        padding: '8px',
        borderRadius: '12px',
        backdropFilter: 'blur(10px)'
      }}>
        <button 
          onClick={() => setPage('register')}
          style={{ 
            padding: '8px 16px', 
            background: page === 'register' ? '#10b981' : 'transparent', 
            color: '#fff', 
            border: '1px solid #10b981', 
            borderRadius: '8px', 
            cursor: 'pointer', 
            fontWeight: 'bold',
            fontSize: '12px',
            transition: 'all 0.3s'
          }}
        >
          회원가입
        </button>
        <button 
          onClick={() => setPage('guss')}
          style={{ 
            padding: '8px 16px', 
            background: page === 'guss' ? '#10b981' : 'transparent', 
            color: '#fff', 
            border: '1px solid #10b981', 
            borderRadius: '8px', 
            cursor: 'pointer', 
            fontWeight: 'bold',
            fontSize: '12px',
            transition: 'all 0.3s'
          }}
        >
          상세정보
        </button>
        <button 
          onClick={() => setPage('admin')}
          style={{ 
            padding: '8px 16px', 
            background: page === 'admin' ? '#10b981' : 'transparent', 
            color: '#fff', 
            border: '1px solid #10b981', 
            borderRadius: '8px', 
            cursor: 'pointer', 
            fontWeight: 'bold',
            fontSize: '12px',
            transition: 'all 0.3s'
          }}
        >
          관리자
        </button>
      </div>

      {/* 페이지 표시 */}
      {page === 'register' && <RegisterPage />}
      {page === 'guss' && <GussPage />}
      {page === 'admin' && <AdminPage />}
    </div>
  );
}