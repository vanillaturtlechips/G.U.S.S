import React from 'react';
import { Search, MapPin, LogIn, Activity } from 'lucide-react';

const dashboard: React.FC = () => {
  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden font-sans">
      {/* 1️⃣ 배경 애니메이션 그리드 */}
      <div className="absolute inset-0 opacity-20" style={{ 
        backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`, 
        backgroundSize: '40px 40px' 
      }} />

      {/* 2️⃣ 메인 콘텐츠 */}
      <div className="relative z-10 p-6 max-w-7xl mx-auto">
        {/* 상단 바 (Top Bar) */}
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center">
          <div className="flex items-center gap-4 flex-1">
            <h1 className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400" 
                style={{ fontFamily: 'Orbitron, sans-serif' }}>GUSS MAP</h1>
            <div className="relative flex-1 max-w-md ml-4">
              <Search className="absolute left-3 top-3 w-5 h-5 text-zinc-500" />
              <input type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-10 py-2 text-white outline-none transition-all" placeholder="헬스장 검색..." />
            </div>
          </div>
          <button className="px-6 py-2 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all shadow-lg shadow-emerald-500/50 flex items-center gap-2">
            <LogIn className="w-5 h-5" /> LOGIN
          </button>
        </div>

        {/* 지도 및 정보창 구역 */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-4 h-[500px] relative">
            <div className="w-full h-full bg-zinc-900 rounded-xl flex items-center justify-center overflow-hidden">
               {/* 사진 속 지도 구현 부분 */}
               <img src="seoul_map.png" alt="서울 지도" className="opacity-40 w-full h-full object-cover" />
               <MapPin className="absolute top-1/2 left-1/2 w-10 h-10 text-emerald-500 animate-bounce" />
            </div>
          </div>
          
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8">
            <h2 className="text-xl font-bold text-emerald-400 mb-6 flex items-center gap-2" style={{ fontFamily: 'Orbitron' }}>
              <Activity className="w-5 h-5" />상 세 정 보
            </h2>
            <div className="space-y-4">
              <p className="text-zinc-400">주소: 서울특별시 ...</p>
              <p className="text-zinc-400">전화번호: 02-123-4567</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default dashboard;