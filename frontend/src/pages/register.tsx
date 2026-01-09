import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';
import StatusModal from './StatusModal';

const Register: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    user_name: '',
    user_phone: '',
    user_id: '',
    user_pw: ''
  });
  const [statusModal, setStatusModal] = useState({ isOpen: false, type: 'SUCCESS' as any, title: '', message: '' });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/api/register', formData);
      setStatusModal({
        isOpen: true,
        type: 'SUCCESS',
        title: 'JOIN SUCCESS',
        message: '성공적으로 등록되었습니다.\n로그인 페이지로 이동합니다.'
      });
    } catch (error: any) {
      setStatusModal({
        isOpen: true,
        type: 'ERROR',
        title: 'FAILED',
        message: error.response?.data || '가입 정보 중복 또는 서버 오류입니다.'
      });
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

      <StatusModal 
        isOpen={statusModal.isOpen}
        type={statusModal.type}
        title={statusModal.title}
        message={statusModal.message}
        onClose={() => setStatusModal({ ...statusModal, isOpen: false })}
        onConfirm={() => {
          if (statusModal.type === 'SUCCESS') navigate('/login');
          else setStatusModal({ ...statusModal, isOpen: false });
        }}
      />
    </div>
  );
};

export default Register;