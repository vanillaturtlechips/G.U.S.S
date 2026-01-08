import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, MapPin, LogIn, Activity, Shield } from 'lucide-react';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden font-sans">
      {/* 배경 애니메이션 그리드 */}
      <div className="absolute inset-0 opacity-20" style={{ 
        backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`, 
        backgroundSize: '40px 40px' 
      }} />

      <div className="relative z-10 p-6 max-w-7xl mx-auto">
        {/* 상단 바 */}
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center">
          <div className="flex items-center gap-4 flex-1">
            <h1 className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400" 
                style={{ fontFamily: 'Orbitron, sans-serif' }}>GUSS MAP</h1>
            <div className="relative flex-1 max-w-md ml-4">
              <Search className="absolute left-3 top-3 w-5 h-5 text-zinc-500" />
              <input type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-10 py-2 text-white outline-none transition-all" placeholder="헬스장 검색..." />
            </div>
          </div>
          <button onClick={() => navigate('/login')} className="px-6 py-2 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all shadow-lg shadow-emerald-500/50 flex items-center gap-2">
            <LogIn className="w-5 h-5" /> LOGIN
          </button>
        </div>

        {/* 지도 및 정보창 */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-4 h-[550px] relative">
            <div className="w-full h-full bg-zinc-900 rounded-xl flex items-center justify-center overflow-hidden border border-zinc-800">
               <img src="/seoul_map.png" alt="서울 지도" className="opacity-40 w-full h-full object-cover" />
               <MapPin className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-12 h-12 text-emerald-500 animate-bounce cursor-pointer" 
                       onClick={() => navigate('/guss')} />
            </div>
          </div>
          
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col">
            <h2 className="text-2xl font-bold text-emerald-400 mb-8 flex items-center gap-3" style={{ fontFamily: 'Orbitron' }}>
              <Shield className="w-6 h-6" /> INFO_PANEL
            </h2>
            <div className="space-y-8 flex-1">
              <div className="border-l-4 border-emerald-500 pl-4">
                <p className="text-zinc-500 text-xs">Center Name</p>
                <p className="text-xl font-bold">Trinity Fitness</p>
              </div>
              <div className="border-l-4 border-zinc-800 pl-4">
                <p className="text-zinc-500 text-xs">Location</p>
                <p className="text-white">서울시 강남구</p>
              </div>
              <div className="border-l-4 border-zinc-800 pl-4">
                <p className="text-zinc-500 text-xs">Status</p>
                <p className="text-emerald-400 font-bold flex items-center gap-2">
                  <Activity className="w-4 h-4 animate-pulse" /> OPEN NOW
                </p>
              </div>
            </div>
            <button 
              onClick={() => navigate('/guss')}
              className="w-full mt-10 py-4 bg-zinc-900 border border-emerald-500/30 rounded-xl text-white font-bold hover:bg-emerald-500/10 transition-all"
            >
              DETAIL DATA
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;