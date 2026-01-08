// src/pages/Login.tsx
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [id, setId] = useState('');
  const [pw, setPw] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await api.post('/api/login', { user_id: id, user_pw: pw });
      const { token } = response.data;
      
      // [핵심] 받은 토큰을 로컬스토리지에 저장 (이후 모든 API 요청에 자동으로 붙음)
      localStorage.setItem('token', token);
      localStorage.setItem('isLoggedIn', 'true');
      
      alert('GUSS 시스템 접속 성공');
      navigate('/'); // 대시보드로 이동
    } catch (error: any) {
      alert('로그인 실패: 아이디 또는 비밀번호를 확인하세요.');
    }
  };

  return (
    <div className="min-h-screen bg-black flex items-center justify-center p-6">
      <form onSubmit={handleLogin} className="bg-zinc-950 border-2 border-emerald-500/30 p-8 rounded-3xl w-full max-w-md">
        <h2 className="text-3xl font-black text-emerald-400 mb-8 text-center" style={{ fontFamily: 'Orbitron' }}>GUSS LOGIN</h2>
        <div className="space-y-6">
          <input 
            type="text" placeholder="ID" 
            className="w-full bg-black border-2 border-zinc-800 p-4 rounded-xl text-white outline-none focus:border-emerald-500"
            value={id} onChange={(e) => setId(e.target.value)}
          />
          <input 
            type="password" placeholder="PASSWORD" 
            className="w-full bg-black border-2 border-zinc-800 p-4 rounded-xl text-white outline-none focus:border-emerald-500"
            value={pw} onChange={(e) => setPw(e.target.value)}
          />
          <button type="submit" className="w-full py-4 bg-emerald-500 text-black font-black rounded-xl hover:bg-emerald-400">
            시스템 접속
          </button>
          <p className="text-center text-zinc-500 text-sm cursor-pointer" onClick={() => navigate('/register')}>
            계정이 없으신가요? <span className="text-emerald-400">회원가입</span>
          </p>
        </div>
      </form>
    </div>
  );
};

export default Login;