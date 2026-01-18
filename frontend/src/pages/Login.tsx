import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';
import StatusModal from './StatusModal';
import { requestFCMToken } from '../firebase/firebaseConfig'; 

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [id, setId] = useState('');
  const [pw, setPw] = useState('');
  const [statusModal, setStatusModal] = useState({ isOpen: false, type: 'SUCCESS' as any, title: '', message: '' });

  const handleLogin = async (e: React.FormEvent) => {
  e.preventDefault();
  try {
    const response = await api.post('/api/login', { user_id: id, user_pw: pw });
    const { token, user_name, role } = response.data;

    localStorage.setItem('token', token);
    localStorage.setItem('isLoggedIn', 'true');
    localStorage.setItem('userRole', role || 'USER');

    // 🔥 FCM 토큰 발급 (에러 무시)
    try {
      const fcmToken = await requestFCMToken();
      if (fcmToken) {
        await api.post('/api/login', {
          user_id: id,
          user_pw: pw,
          fcm_token: fcmToken
        });
      }
    } catch (fcmError) {
      console.log('푸시 알림 설정 실패 (무시됨):', fcmError);
      // 로그인은 계속 진행
    }

    setStatusModal({
      isOpen: true,
      type: 'SUCCESS',
      title: 'ACCESS GRANTED',
      message: `${user_name || id} 요원님, GUSS 시스템 접속을 환영합니다.`
    });
  } catch (error: any) {
    setStatusModal({
      isOpen: true,
      type: 'ERROR',
      title: 'AUTH FAILED',
      message: '아이디 또는 비밀번호가 일치하지 않습니다.'
    });
  }
};


  return (
    <div className="min-h-screen bg-black flex items-center justify-center p-6">
      <form onSubmit={handleLogin} className="bg-zinc-950 border-2 border-emerald-500/30 p-8 rounded-3xl w-full max-w-md shadow-[0_0_50px_rgba(16,185,129,0.1)]">
        <h2 className="text-3xl font-black text-emerald-400 mb-8 text-center" style={{ fontFamily: 'Orbitron' }}>GUSS LOGIN</h2>
        <div className="space-y-6">
          <input 
            type="text" placeholder="ID" 
            className="w-full bg-black border-2 border-zinc-800 p-4 rounded-xl text-white outline-none focus:border-emerald-500 transition-all"
            value={id} onChange={(e) => setId(e.target.value)}
          />
          <input 
            type="password" placeholder="PASSWORD" 
            className="w-full bg-black border-2 border-zinc-800 p-4 rounded-xl text-white outline-none focus:border-emerald-500 transition-all"
            value={pw} onChange={(e) => setPw(e.target.value)}
          />
          <button type="submit" className="w-full py-4 bg-emerald-500 text-black font-black rounded-xl hover:bg-emerald-400 transition-all active:scale-95 shadow-lg shadow-emerald-500/20">
            시스템 접속
          </button>
          <p className="text-center text-zinc-500 text-sm cursor-pointer" onClick={() => navigate('/register')}>
            계정이 없으신가요? <span className="text-emerald-400">회원가입</span>
          </p>
        </div>
      </form>

      <StatusModal 
        isOpen={statusModal.isOpen}
        type={statusModal.type}
        title={statusModal.title}
        message={statusModal.message}
        onClose={() => setStatusModal({ ...statusModal, isOpen: false })}
        onConfirm={() => {
          if (statusModal.type === 'SUCCESS') navigate('/');
          else setStatusModal({ ...statusModal, isOpen: false });
        }}
      />
    </div>
  );
};

export default Login;
