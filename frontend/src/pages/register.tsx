// src/pages/Register.tsx (예시)
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';

const Register: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    user_name: '',
    user_phone: '',
    user_id: '',
    user_pw: ''
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // 백엔드의 HandleRegister 호출
      await api.post('/api/register', formData);
      alert('회원가입 성공! 로그인 페이지로 이동합니다.');
      navigate('/login');
    } catch (error: any) {
      alert(error.response?.data || '회원가입 실패');
    }
  };

  return (
    <div className="min-h-screen bg-black text-white flex items-center justify-center p-6">
      <form onSubmit={handleSubmit} className="bg-zinc-950 border-2 border-emerald-500/30 p-8 rounded-3xl w-full max-w-md">
        <h2 className="text-3xl font-black text-emerald-400 mb-8 text-center" style={{ fontFamily: 'Orbitron' }}>JOIN GUSS</h2>
        <div className="space-y-4">
          <input 
            type="text" placeholder="이름" 
            className="w-full bg-black border-2 border-zinc-800 p-3 rounded-xl focus:border-emerald-500 outline-none"
            onChange={(e) => setFormData({...formData, user_name: e.target.value})}
            required
          />
          <input 
            type="text" placeholder="전화번호" 
            className="w-full bg-black border-2 border-zinc-800 p-3 rounded-xl focus:border-emerald-500 outline-none"
            onChange={(e) => setFormData({...formData, user_phone: e.target.value})}
            required
          />
          <input 
            type="text" placeholder="아이디" 
            className="w-full bg-black border-2 border-zinc-800 p-3 rounded-xl focus:border-emerald-500 outline-none"
            onChange={(e) => setFormData({...formData, user_id: e.target.value})}
            required
          />
          <input 
            type="password" placeholder="비밀번호" 
            className="w-full bg-black border-2 border-zinc-800 p-3 rounded-xl focus:border-emerald-500 outline-none"
            onChange={(e) => setFormData({...formData, user_pw: e.target.value})}
            required
          />
        </div>
        <button type="submit" className="w-full mt-8 py-4 bg-gradient-to-r from-emerald-500 to-lime-500 text-black font-black rounded-xl hover:scale-105 transition-all">
          회원가입 하기
        </button>
      </form>
    </div>
  );
};

export default Register;