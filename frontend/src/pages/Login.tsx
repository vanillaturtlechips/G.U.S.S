import React from 'react';
import { Shield, User, Lock } from 'lucide-react';

export default function Login() {
  return (
    <div className="min-h-screen bg-black text-white flex items-center justify-center relative overflow-hidden">
      {/* 배경 효과 */}
      <div className="absolute inset-0 opacity-10" style={{ backgroundImage: `radial-gradient(#10b981 1px, transparent 1px)`, backgroundSize: '20px 20px' }} />
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-emerald-500/10 rounded-full blur-[120px]" />

      <div className="relative z-10 w-full max-w-md p-8">
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-10 shadow-[0_0_50px_rgba(0,0,0,0.5)] backdrop-blur-xl">
          <div className="text-center mb-10">
            <Shield className="w-16 h-16 text-emerald-500 mx-auto mb-4 animate-pulse" />
            <h1 className="text-4xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400" style={{ fontFamily: 'Orbitron' }}>GUSS LOGIN</h1>
            <p className="text-zinc-500 mt-2">Access the Trinity System</p>
          </div>

          <div className="space-y-6">
            <div className="space-y-2">
              <label className="text-sm font-bold text-emerald-500 ml-1 flex items-center gap-2"><User className="w-4 h-4"/> USER ID</label>
              <input type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 text-white placeholder-zinc-700 outline-none transition-all" placeholder="ID를 입력하세요" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-bold text-emerald-500 ml-1 flex items-center gap-2"><Lock className="w-4 h-4"/> PASSWORD</label>
              <input type="password" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 text-white placeholder-zinc-700 outline-none transition-all" placeholder="••••••••" />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 mt-10">
            <button className="py-3 bg-zinc-900 hover:bg-zinc-800 border border-emerald-500/30 rounded-xl text-white font-bold transition-all">
              회원가입
            </button>
            <button className="py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold shadow-lg shadow-emerald-500/30 hover:scale-105 transition-all">
              로그인
            </button>
          </div>
        </div>
      </div>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;900&display=swap');
      `}</style>
    </div>
  );
}