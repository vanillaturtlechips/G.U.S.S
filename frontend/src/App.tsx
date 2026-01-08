import React from 'react';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom';
import { Search, MapPin, LogIn, Activity, Shield, User, Lock } from 'lucide-react';

// --- 1. 대시보드(지도) 컴포넌트 ---
const Dashboard = () => {
  const navigate = useNavigate(); // 페이지 이동을 위한 함수

  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden font-sans">
      <div className="absolute inset-0 opacity-20" style={{ backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`, backgroundSize: '40px 40px' }} />
      <div className="relative z-10 p-6 max-w-7xl mx-auto">
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center">
          <h1 className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400" style={{ fontFamily: 'Orbitron' }}>GUSS MAP</h1>
          {/* 로그인 버튼 클릭 시 /login 경로로 이동 */}
          <button onClick={() => navigate('/login')} className="px-6 py-2 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all flex items-center gap-2">
            <LogIn className="w-5 h-5" /> LOGIN
          </button>
        </div>
        {/* ... 나머지 지도 코드 (생략, 기존과 동일) ... */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-4 h-[550px] relative">
                <div className="w-full h-full bg-zinc-900 rounded-xl flex items-center justify-center overflow-hidden border border-zinc-800">
                    <img src="/seoul_map.png" alt="서울 지도" className="opacity-40 w-full h-full object-cover" />
                    <MapPin className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-12 h-12 text-emerald-500 animate-bounce" />
                </div>
            </div>
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col">
                <h2 className="text-2xl font-bold text-emerald-400 mb-8 flex items-center gap-3" style={{ fontFamily: 'Orbitron' }}><Shield className="w-6 h-6" /> INFO_PANEL</h2>
                <div className="space-y-8 flex-1">
                    <div className="border-l-4 border-emerald-500 pl-4"><p className="text-zinc-500 text-xs">Center Name</p><p className="text-xl font-bold">Trinity Fitness</p></div>
                    <div className="border-l-4 border-zinc-800 pl-4"><p className="text-zinc-500 text-xs">Location</p><p className="text-white">서울시 강남구</p></div>
                </div>
                <button className="w-full mt-10 py-4 bg-zinc-900 border border-emerald-500/30 rounded-xl text-white font-bold">DETAIL DATA</button>
            </div>
        </div>
      </div>
    </div>
  );
};

// --- 2. 로그인 컴포넌트 (스케치 반영) ---
const LoginPage = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-black text-white flex items-center justify-center relative overflow-hidden font-sans">
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
              <input type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none" placeholder="Username" />
            </div>
            <div className="space-y-2">
              <label className="text-xs font-bold text-emerald-500 flex items-center gap-2"><Lock className="w-4 h-4"/> PWD</label>
              <input type="password" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none" placeholder="Password" />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <button className="py-3 bg-zinc-900 border border-emerald-500/30 rounded-xl font-bold">SIGN UP</button>
            <button onClick={() => navigate('/')} className="py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold">LOGIN</button>
          </div>
        </div>
      </div>
    </div>
  );
};

// --- 3. 메인 App (경로 설정) ---
export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/login" element={<LoginPage />} />
      </Routes>
    </BrowserRouter>
  );
}