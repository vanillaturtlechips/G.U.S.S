import React, { useState } from 'react';
import { User, Mail, Lock, Phone, Shield, CheckCircle, Dumbbell } from 'lucide-react';

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    name: '',
    id: '',
    password: '',
    phone: ''
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('회원가입:', formData);
  };

  return (
    <div className="min-h-screen relative overflow-hidden bg-black">
      {/* Animated Grid Background */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute inset-0" style={{
          backgroundImage: `
            linear-gradient(to right, #10b981 1px, transparent 1px),
            linear-gradient(to bottom, #10b981 1px, transparent 1px)
          `,
          backgroundSize: '40px 40px',
          animation: 'gridMove 20s linear infinite'
        }} />
      </div>

      {/* Gradient Overlays */}
      <div className="absolute top-0 left-0 w-96 h-96 bg-emerald-500/20 rounded-full blur-3xl animate-pulse" 
           style={{ animationDuration: '3s' }} />
      <div className="absolute bottom-0 right-0 w-96 h-96 bg-lime-500/20 rounded-full blur-3xl animate-pulse" 
           style={{ animationDuration: '4s', animationDelay: '1s' }} />

      <div className="relative z-10 min-h-screen flex items-center justify-center p-6">
        <div className="w-full max-w-6xl grid lg:grid-cols-2 gap-8 items-center">
          {/* Left Side - Info Panel */}
          <div className="hidden lg:block space-y-8">
            <div className="space-y-4">
              <div className="inline-block">
                <div className="flex items-center gap-3 px-6 py-3 bg-emerald-500/10 border border-emerald-500/30 rounded-2xl backdrop-blur-xl">
                  <Dumbbell className="w-6 h-6 text-emerald-400" />
                  <span className="text-emerald-400 font-bold text-lg tracking-wider" style={{ fontFamily: 'Orbitron, sans-serif' }}>GUSS SYSTEM</span>
                </div>
              </div>
              
              <h1 className="text-6xl font-black leading-tight" style={{ fontFamily: 'Orbitron, sans-serif' }}>
                <span className="text-white">GYM</span>
                <br />
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400">
                  REVOLUTION
                </span>
              </h1>
              
              <p className="text-xl text-emerald-300">
                차세대 헬스장 관리 시스템에 오신 것을 환영합니다
              </p>
            </div>

            {/* Features */}
            <div className="space-y-4">
              {[
                { icon: '⚡', title: '실시간 모니터링', desc: 'IoT 센서 기반 실시간 현황 파악' },
                { icon: '🎯', title: '스마트 예약', desc: 'AI 추천 시스템으로 최적 시간 제안' },
                { icon: '📊', title: '데이터 분석', desc: '운동 패턴 분석 및 개인 맞춤 리포트' }
              ].map((feature, idx) => (
                <div 
                  key={idx}
                  className="flex items-center gap-4 p-4 bg-zinc-900/50 border border-emerald-500/20 rounded-2xl hover:border-emerald-500/40 hover:bg-zinc-900/70 transition-all duration-300"
                >
                  <div className="text-4xl">{feature.icon}</div>
                  <div>
                    <h3 className="font-bold text-white">{feature.title}</h3>
                    <p className="text-sm text-emerald-400/70">{feature.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Right Side - Registration Form */}
          <div className="w-full">
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 shadow-2xl shadow-emerald-500/10">
              {/* Glowing Top Border */}
              <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-t-3xl mb-8 animate-pulse" />
              
              {/* Header */}
              <div className="text-center mb-8">
                <div className="inline-block p-4 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-2xl mb-4 shadow-lg shadow-emerald-500/50">
                  <Shield className="w-10 h-10 text-black" strokeWidth={2.5} />
                </div>
                <h2 className="text-3xl font-black text-white mb-2" style={{ fontFamily: 'Orbitron, sans-serif' }}>회원가입</h2>
                <p className="text-emerald-400">정보를 입력하고 시스템에 등록하세요</p>
              </div>

              {/* Form */}
              <form onSubmit={handleSubmit} className="space-y-5">
                {/* 성명 */}
                <div className="space-y-2">
                  <label className="text-sm font-bold text-emerald-400 uppercase tracking-wider">
                    성명 (Name)
                  </label>
                  <div className="relative group">
                    <User className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-emerald-500 transition-all group-focus-within:scale-110 group-focus-within:text-lime-400" />
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl pl-12 pr-4 py-4 text-white placeholder-zinc-600 focus:outline-none transition-all font-mono"
                      placeholder="홍길동"
                      required
                    />
                  </div>
                </div>

                {/* ID */}
                <div className="space-y-2">
                  <label className="text-sm font-bold text-emerald-400 uppercase tracking-wider">
                    아이디 (ID)
                  </label>
                  <div className="relative group">
                    <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-emerald-500 transition-all group-focus-within:scale-110 group-focus-within:text-lime-400" />
                    <input
                      type="text"
                      value={formData.id}
                      onChange={(e) => setFormData({ ...formData, id: e.target.value })}
                      className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl pl-12 pr-4 py-4 text-white placeholder-zinc-600 focus:outline-none transition-all font-mono"
                      placeholder="user123"
                      required
                    />
                  </div>
                </div>

                {/* PWD */}
                <div className="space-y-2">
                  <label className="text-sm font-bold text-emerald-400 uppercase tracking-wider">
                    비밀번호 (PWD)
                  </label>
                  <div className="relative group">
                    <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-emerald-500 transition-all group-focus-within:scale-110 group-focus-within:text-lime-400" />
                    <input
                      type="password"
                      value={formData.password}
                      onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                      className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl pl-12 pr-4 py-4 text-white placeholder-zinc-600 focus:outline-none transition-all"
                      placeholder="8자 이상"
                      required
                    />
                  </div>
                </div>

                {/* Phone */}
                <div className="space-y-2">
                  <label className="text-sm font-bold text-emerald-400 uppercase tracking-wider">
                    휴대폰 번호 (Phone)
                  </label>
                  <div className="relative group">
                    <Phone className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-emerald-500 transition-all group-focus-within:scale-110 group-focus-within:text-lime-400" />
                    <input
                      type="tel"
                      value={formData.phone}
                      onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                      className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl pl-12 pr-4 py-4 text-white placeholder-zinc-600 focus:outline-none transition-all font-mono"
                      placeholder="010-0000-0000"
                      required
                    />
                  </div>
                </div>

                {/* Info Section */}
                <div className="p-4 bg-emerald-500/10 border border-emerald-500/30 rounded-xl">
                  <div className="flex items-center gap-2 mb-2">
                    <CheckCircle className="w-5 h-5 text-emerald-400 animate-pulse" />
                    <span className="text-sm font-bold text-emerald-400">보안 인증 시스템</span>
                  </div>
                  <p className="text-xs text-emerald-500/70">
                    입력하신 모든 정보는 256비트 암호화로 안전하게 보호됩니다.
                  </p>
                </div>

                {/* Submit Button */}
                <button
                  type="submit"
                  className="w-full relative overflow-hidden group rounded-xl"
                >
                  <div className="absolute inset-0 bg-gradient-to-r from-emerald-500 to-lime-500 transition-all duration-300 group-hover:scale-105" />
                  <div className="absolute inset-0 bg-gradient-to-r from-emerald-400 to-lime-400 opacity-0 group-hover:opacity-100 transition-opacity" />
                  <div className="relative flex items-center justify-center gap-3 py-4 text-black font-black text-lg tracking-wider">
                    <Shield className="w-6 h-6" />
                    REGISTER NOW
                  </div>
                </button>

                {/* Login Link */}
                <div className="text-center pt-4">
                  <p className="text-zinc-500 text-sm">
                    이미 계정이 있으신가요?{' '}
                    <a href="#" className="text-emerald-400 hover:text-emerald-300 font-bold transition-colors">
                      로그인하기
                    </a>
                  </p>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;500;600;700;800;900&display=swap');
        
        @keyframes gridMove {
          0% {
            transform: translate(0, 0);
          }
          100% {
            transform: translate(40px, 40px);
          }
        }
      `}</style>
    </div>
  );
}