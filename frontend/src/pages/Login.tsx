import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Shield, User, Lock } from 'lucide-react';

export default function LoginPage() {
  const navigate = useNavigate();
  const [id, setId] = useState('');

  const handleLogin = () => {
    // 임시 로그인 로직
    localStorage.setItem('isLoggedIn', 'true');
    
    if (id === 'admin') {
      navigate('/admin'); // 관리자 이동
    } else {
      navigate('/'); // 메인으로 복귀
    }
  };

  return (
    <div className="min-h-screen bg-black text-white flex items-center justify-center relative overflow-hidden">
      <div className="absolute inset-0 opacity-10" style={{ backgroundImage: `radial-gradient(#10b981 1px, transparent 1px)`, backgroundSize: '20px 20px' }} />
      
      <div className="relative z-10 w-full max-w-md p-8">
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-10 shadow-2xl">
          <div className="text-center mb-10">
            <Shield className="w-16 h-16 text-emerald-500 mx-auto mb-4 animate-pulse" />
            <h1 className="text-3xl font-black text-emerald-400" style={{ fontFamily: 'Orbitron' }}>GUSS ACCESS</h1>
          </div>
          
          <div className="space-y-6 mb-10">
            <div className="space-y-2">
              <label className="text-xs font-bold text-emerald-500 flex items-center gap-2"><User className="w-4 h-4"/> ID</label>
              <input 
                type="text" 
                value={id}
                onChange={(e) => setId(e.target.value)}
                className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none transition-all" 
                placeholder="admin 입력 시 관리자 페이지" 
              />
            </div>
            <div className="space-y-2">
              <label className="text-xs font-bold text-emerald-500 flex items-center gap-2"><Lock className="w-4 h-4"/> PWD</label>
              <input type="password" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none transition-all" placeholder="Password" />
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <button onClick={() => navigate('/register')} className="py-3 bg-zinc-900 border border-emerald-500/30 rounded-xl font-bold hover:bg-zinc-800 transition-all">SIGN UP</button>
            <button onClick={handleLogin} className="py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all">LOGIN</button>
          </div>
        </div>
      </div>
    </div>
  );
}