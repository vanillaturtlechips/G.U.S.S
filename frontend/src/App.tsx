import { useState } from 'react';
import RegisterPage from './pages/register';
import GussPage from './pages/guss';
import AdminPage from './pages/admin';

export default function App() {
  // 여기에 'guss', 'register', 'admin' 중 보고 싶은 페이지 이름을 넣으세요.
  const [page] = useState('guss'); 

  return (
    <div className="relative">
      {/* 버튼이 있던 자리가 삭제되었습니다. */}

      {/* 페이지 표시 */}
      {page === 'register' && <RegisterPage />}
      {page === 'guss' && <GussPage />}
      {page === 'admin' && <AdminPage />}
    </div>
  );
}