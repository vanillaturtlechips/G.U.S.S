import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, MapPin, LogIn, Activity, Shield } from 'lucide-react';
import seoulMapImg from '../assets/seoul-map.png'; 

// 1. 마커를 한강(물 위)에서 도심 육지 안쪽으로 깊숙이 재배치했습니다.
const GYM_DATA = [
  { id: 'gangnam', name: 'Trinity Fitness 강남', location: '서울시 강남구 테헤란로', status: 'OPEN NOW', top: '78%', left: '68%', members: 42 },
  { id: 'hongdae', name: 'GUSS 홍대점', location: '서울시 마포구 양화로', status: 'OPEN NOW', top: '32%', left: '28%', members: 15 },
  { id: 'seongsu', name: 'GUSS 성수 스튜디오', location: '서울시 성동구 성수이로', status: 'OPEN NOW', top: '35%', left: '72%', members: 28 },
  { id: 'yeouido', name: 'GUSS 여의도 본점', location: '서울시 영등포구 여의나루로', status: 'OPEN NOW', top: '62%', left: '44%', members: 54 },
  { id: 'jamsil', name: 'GUSS 잠실 센터', location: '서울시 송파구 올림픽로', status: 'OPEN NOW', top: '72%', left: '85%', members: 31 },
  { id: 'jongno', name: 'GUSS 종로점', location: '서울시 종로구 종로', status: 'OPEN NOW', top: '28%', left: '50%', members: 19 },
];

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [selectedGym, setSelectedGym] = useState(GYM_DATA[0]);

  const handleDetailView = () => {
    navigate(`/guss?gymId=${selectedGym.id}`);
  };

  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden font-sans">
      {/* 배경 그리드 디자인 유지 */}
      <div className="absolute inset-0 opacity-20" style={{ 
        backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`, 
        backgroundSize: '40px 40px' 
      }} />

      <div className="relative z-10 p-6 max-w-7xl mx-auto h-screen flex flex-col">
        {/* 상단 바 디자인 */}
        <div className="bg-zinc-950/80 backdrop-blur-sm border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center shadow-[0_0_30px_rgba(16,185,129,0.1)]">
          <div className="flex items-center gap-4 flex-1">
            <h1 className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400" 
                style={{ fontFamily: 'Orbitron, sans-serif' }}>GUSS MAP</h1>
            <div className="relative flex-1 max-w-md ml-4">
              <Search className="absolute left-3 top-3 w-5 h-5 text-zinc-500" />
              <input type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-10 py-2 text-white outline-none transition-all" placeholder="주변 헬스장 검색..." />
            </div>
          </div>
          <button onClick={() => navigate('/login')} className="px-6 py-2 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all flex items-center gap-2">
            <LogIn className="w-5 h-5" /> LOGIN
          </button>
        </div>

        <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-8 overflow-hidden">
          {/* 지도 섹션 */}
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-4 relative overflow-hidden">
            <div className="w-full h-full bg-zinc-900 rounded-xl flex items-center justify-center overflow-hidden border border-zinc-800 relative">
              <img src={seoulMapImg} alt="서울 지도" className="w-full h-full object-cover opacity-30 grayscale brightness-75" />
              
              {/* 마커 가시성 강화 */}
              {GYM_DATA.map((gym) => (
                <button
                  key={gym.id}
                  className="absolute -translate-x-1/2 -translate-y-1/2 group z-20"
                  style={{ top: gym.top, left: gym.left }}
                  onClick={() => setSelectedGym(gym)}
                >
                  <div className="relative">
                    {/* 2. 마커 색상을 불투명한 emerald-500으로 변경하고 강한 그림자를 추가했습니다. */}
                    <MapPin 
                      className={`w-10 h-10 transition-all duration-300 drop-shadow-[0_4px_8px_rgba(0,0,0,0.8)] ${
                        selectedGym.id === gym.id 
                        ? 'text-emerald-400 scale-125 drop-shadow-[0_0_20px_rgba(16,185,129,1)]' 
                        : 'text-emerald-500 opacity-100 group-hover:text-emerald-300 group-hover:scale-110' 
                      }`} 
                      fill={selectedGym.id === gym.id ? "#10b98166" : "none"}
                    />
                  </div>
                </button>
              ))}
            </div>
          </div>
          
          {/* 정보 패널 (INFO_PANEL) */}
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col justify-between">
            <div className="animate-in fade-in duration-500">
              <h2 className="text-2xl font-bold text-emerald-400 mb-8 flex items-center gap-3" style={{ fontFamily: 'Orbitron' }}>
                <Shield className="w-6 h-6" /> INFO_PANEL
              </h2>
              
              <div className="space-y-8">
                <div className="border-l-4 border-emerald-500 pl-4">
                  <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Center Name</p>
                  <p className="text-2xl font-black">{selectedGym.name}</p>
                </div>
                
                <div className="border-l-4 border-zinc-800 pl-4">
                  <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Location</p>
                  <p className="text-white font-medium">{selectedGym.location}</p>
                </div>

                <div className="border-l-4 border-zinc-800 pl-4">
                  <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Live Status</p>
                  <p className="text-emerald-400 font-bold flex items-center gap-2">
                    <Activity className="w-4 h-4 animate-pulse" /> {selectedGym.status}
                  </p>
                  <p className="text-zinc-500 text-xs italic mt-1">현재 이용 중: {selectedGym.members}명</p>
                </div>
              </div>
            </div>

            <button 
              onClick={handleDetailView}
              className="w-full py-4 bg-zinc-900 border border-emerald-500/30 rounded-xl text-white font-bold hover:bg-emerald-500/10 hover:border-emerald-500 transition-all shadow-[0_0_15px_rgba(16,185,129,0.1)]"
            >
              상세 데이터 확인하기
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;