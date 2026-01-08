import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { User, Mail, Lock, Phone, Shield, CheckCircle, Dumbbell } from 'lucide-react';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({ name: '', id: '', password: '', phone: '' });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    alert('회원가입이 완료되었습니다! 로그인 해주세요.');
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-black text-white relative flex items-center justify-center p-6">
      <div className="absolute inset-0 opacity-10" style={{ backgroundImage: `linear-gradient(#10b981 1px, transparent 1px), linear-gradient(to right, #10b981 1px, transparent 1px)`, backgroundSize: '40px 40px' }} />
      
      <div className="relative z-10 w-full max-w-lg bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-10">
        <div className="text-center mb-8">
          <div className="inline-block p-4 bg-emerald-500 rounded-2xl mb-4">
            <Dumbbell className="w-8 h-8 text-black" />
          </div>
          <h2 className="text-3xl font-black text-white" style={{ fontFamily: 'Orbitron' }}>JOIN REVOLUTION</h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div className="space-y-2">
            <label className="text-sm font-bold text-emerald-500">성명</label>
            <input required type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none" placeholder="홍길동" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-bold text-emerald-500">아이디</label>
            <input required type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none" placeholder="guss_user" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-bold text-emerald-500">비밀번호</label>
            <input required type="password" title="8자 이상" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none" placeholder="••••••••" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-bold text-emerald-500">연락처</label>
            <input required type="tel" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 outline-none" placeholder="010-0000-0000" />
          </div>

          <button type="submit" className="w-full py-4 bg-gradient-to-r from-emerald-500 to-lime-500 text-black font-black rounded-xl mt-6 shadow-lg shadow-emerald-500/20">
            REGISTER NOW
          </button>
        </form>
        
        <p className="text-center mt-6 text-zinc-500 text-sm">
          이미 계정이 있으신가요? <span onClick={() => navigate('/login')} className="text-emerald-400 cursor-pointer font-bold">로그인하기</span>
        </p>
      </div>
    </div>
  );
}