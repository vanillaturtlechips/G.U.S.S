import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, MapPin, LogIn, LogOut, Activity, Shield, Phone } from 'lucide-react';
import seoulMapImg from '../assets/seoul-map.png'; 
import api from '../api/axios';

const POSITIONS: { [key: number]: { top: string; left: string } } = {
  1: { top: '78%', left: '68%' }, 
  2: { top: '32%', left: '28%' }, 
  3: { top: '35%', left: '72%' }, 
  4: { top: '62%', left: '44%' }, 
  5: { top: '72%', left: '85%' }, 
};

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [gyms, setGyms] = useState<any[]>([]);
  const [selectedGym, setSelectedGym] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [searchKeyword, setSearchKeyword] = useState('');

  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';
  const userRole = localStorage.getItem('userRole');

  const fetchGyms = async (keyword = '') => {
    try {
      const response = await api.get(`/api/gyms?search=${keyword}`);
      const data = response.data;
      setGyms(Array.isArray(data) ? data : []);
      if (Array.isArray(data) && data.length > 0 && !selectedGym) {
        setSelectedGym(data[0]);
      }
      setLoading(false);
    } catch (error) {
      console.error("데이터 로딩 실패:", error);
      setGyms([]);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGyms(searchKeyword);
    const interval = setInterval(() => fetchGyms(searchKeyword), 5000); 
    return () => clearInterval(interval);
  }, [selectedGym, searchKeyword]);

  const handleAuthAction = () => {
    if (isLoggedIn) {
      localStorage.clear();
      window.location.reload(); 
    } else {
      navigate('/login');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') fetchGyms(searchKeyword);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-black flex flex-col items-center justify-center text-emerald-400 font-black">
        <Activity className="w-12 h-12 animate-spin mb-4" />
        <p className="tracking-[0.5em]">GUSS SYSTEM LOADING...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden font-sans">
      <div className="absolute inset-0 opacity-20" style={{ 
        backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`, 
        backgroundSize: '40px 40px' 
      }} />

      <div className="relative z-10 p-6 max-w-7xl mx-auto h-screen flex flex-col">
        <div className="bg-zinc-950/80 backdrop-blur-sm border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center shadow-[0_0_30px_rgba(16,185,129,0.1)]">
          <div className="flex items-center gap-4 flex-1">
            <h1 className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400 cursor-pointer" 
                style={{ fontFamily: 'Orbitron, sans-serif' }}
                onClick={() => {setSearchKeyword(''); fetchGyms('');}}>
              GUSS MAP
            </h1>
            <div className="relative flex-1 max-w-md ml-4">
              <Search className="absolute left-3 top-3 w-5 h-5 text-zinc-500" />
              <input 
                type="text" 
                value={searchKeyword}
                onChange={(e) => setSearchKeyword(e.target.value)}
                onKeyDown={handleKeyDown}
                className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-10 py-2 text-white outline-none transition-all" 
                placeholder="지점명을 검색하고 엔터를 누르세요..." 
              />
            </div>
          </div>

          <div className="flex items-center gap-3">
            {isLoggedIn && userRole === 'ADMIN' && (
              <button onClick={() => navigate('/admin')} className="px-6 py-2 bg-purple-600/20 border-2 border-purple-500/50 rounded-xl text-purple-400 font-bold hover:bg-purple-500 hover:text-white transition-all">
                ADMIN PAGE
              </button>
            )}
            <button onClick={handleAuthAction} className="px-6 py-2 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all flex items-center gap-2">
              {isLoggedIn ? <LogOut className="w-5 h-5" /> : <LogIn className="w-5 h-5" />}
              {isLoggedIn ? 'LOGOUT' : 'LOGIN'}
            </button>
          </div>
        </div>

        <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-8 overflow-hidden mb-4">
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-4 relative overflow-hidden">
            <div className="w-full h-full bg-zinc-900 rounded-xl flex items-center justify-center overflow-hidden border border-zinc-800 relative shadow-inner">
              <img src={seoulMapImg} alt="서울 지도" className="w-full h-full object-cover opacity-30 grayscale brightness-75" />
              {gyms.map((gym) => (
                <button
                  key={gym.guss_number}
                  className="absolute -translate-x-1/2 -translate-y-1/2 group z-20"
                  style={{ top: POSITIONS[gym.guss_number]?.top || '50%', left: POSITIONS[gym.guss_number]?.left || '50%' }}
                  onClick={() => setSelectedGym(gym)}
                >
                  <MapPin className={`w-10 h-10 transition-all duration-300 ${selectedGym?.guss_number === gym.guss_number ? 'text-emerald-400 scale-125 drop-shadow-[0_0_20px_rgba(16,185,129,1)]' : 'text-emerald-600 opacity-80 group-hover:text-emerald-300'}`} />
                </button>
              ))}
            </div>
          </div>
          
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col justify-between shadow-[0_0_50px_rgba(0,0,0,0.5)]">
            {selectedGym ? (
              <div className="animate-in fade-in duration-500 h-full flex flex-col">
                <h2 className="text-2xl font-bold text-emerald-400 mb-8 flex items-center gap-3 uppercase italic" style={{ fontFamily: 'Orbitron' }}>
                  <Shield className="w-6 h-6" /> INFO_PANEL
                </h2>
                <div className="space-y-6 flex-1 overflow-y-auto pr-2 custom-scrollbar">
                  <div className="border-l-4 border-emerald-500 pl-4">
                    <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Center Name</p>
                    <p className="text-2xl font-black">{selectedGym.guss_name}</p>
                    <div className="flex items-center gap-2 mt-1 text-zinc-400">
                      <Phone className="w-3 h-3 text-emerald-500" />
                      <p className="text-sm font-medium tracking-tighter">{selectedGym.guss_phone || '전화번호 미등록'}</p>
                    </div>
                  </div>
                  <div className="border-l-4 border-zinc-800 pl-4">
                    <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Live Status</p>
                    <p className="text-emerald-400 font-bold flex items-center gap-2 uppercase italic">
                      <Activity className="w-4 h-4 animate-pulse" /> {selectedGym.guss_status}
                    </p>
                    <p className="text-zinc-400 mt-2 text-sm">
                      현재 이용 인원: <span className="text-white font-black text-lg">{selectedGym.guss_user_count}</span> 명 / {selectedGym.guss_size}
                    </p>
                  </div>
                </div>
              </div>
            ) : <div className="flex flex-col items-center justify-center h-full opacity-30"><Shield className="w-12 h-12 mb-4" /><p className="font-bold tracking-widest">지점을 선택해주세요.</p></div>}
            <button onClick={() => navigate(`/guss?gymId=${selectedGym?.guss_number}`)} className="w-full py-5 mt-6 bg-zinc-900 border border-emerald-500/30 rounded-2xl text-white font-bold hover:bg-emerald-500/10 hover:border-emerald-500 transition-all uppercase tracking-[0.2em]" style={{ fontFamily: 'Orbitron' }}>
              Check Details
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
