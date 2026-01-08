import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, MapPin, LogIn, LogOut, Activity, Shield } from 'lucide-react';
import seoulMapImg from '../assets/seoul-map.png'; // 서울 지도 이미지 경로
import api from '../api/axios'; // 우리가 만든 axios 설정 파일

/** * [데이터 매핑] DB의 guss_number를 지도상의 좌표와 연결합니다. 
 * 마커가 바다에 빠지지 않도록 육지 안쪽 좌표로 고정했습니다.
 */
const POSITIONS: { [key: number]: { top: string; left: string } } = {
  1: { top: '78%', left: '68%' }, // 강남
  2: { top: '32%', left: '28%' }, // 홍대
  3: { top: '35%', left: '72%' }, // 성수
  4: { top: '62%', left: '44%' }, // 여의도
  5: { top: '72%', left: '85%' }, // 잠실
};

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  
  // 상태 관리
  const [gyms, setGyms] = useState<any[]>([]);
  const [selectedGym, setSelectedGym] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  
  // 로그인 상태 확인 (localStorage 활용)
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';

  /** * [실시간 데이터 패칭] 
   * 백엔드 API에서 체육관 목록을 가져와 상태를 업데이트합니다.
   */
  const fetchGyms = async () => {
    try {
      const response = await api.get('/api/gyms');
      const data = response.data;
      setGyms(data);
      
      // 처음 로딩 시 첫 번째 체육관을 기본 선택
      if (data.length > 0 && !selectedGym) {
        setSelectedGym(data[0]);
      }
      setLoading(false);
    } catch (error) {
      console.error("데이터 로딩 실패:", error);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGyms();
    // 5초마다 데이터를 갱신하여 인원수 변화를 실시간으로 반영합니다.
    const interval = setInterval(fetchGyms, 5000); 
    return () => clearInterval(interval);
  }, [selectedGym]);

  // 로그인/로그아웃 처리
  const handleAuthAction = () => {
    if (isLoggedIn) {
      localStorage.removeItem('token');
      localStorage.removeItem('isLoggedIn');
      window.location.reload(); // 상태 반영을 위한 새로고침
    } else {
      navigate('/login');
    }
  };

  const handleDetailView = () => {
    if (selectedGym) {
      navigate(`/guss?gymId=${selectedGym.guss_number}`);
    }
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
      {/* 배경 애니메이션 그리드 [디자인 유지] */}
      <div className="absolute inset-0 opacity-20" style={{ 
        backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`, 
        backgroundSize: '40px 40px' 
      }} />

      <div className="relative z-10 p-6 max-w-7xl mx-auto h-screen flex flex-col">
        {/* 상단 바 [디자인 유지] */}
        <div className="bg-zinc-950/80 backdrop-blur-sm border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center shadow-[0_0_30px_rgba(16,185,129,0.1)]">
          <div className="flex items-center gap-4 flex-1">
            <h1 className="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400" 
                style={{ fontFamily: 'Orbitron, sans-serif' }}>GUSS MAP</h1>
            <div className="relative flex-1 max-w-md ml-4">
              <Search className="absolute left-3 top-3 w-5 h-5 text-zinc-500" />
              <input type="text" className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-10 py-2 text-white outline-none transition-all" placeholder="주변 헬스장 검색..." />
            </div>
          </div>
          <button 
            onClick={handleAuthAction} 
            className="px-6 py-2 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all flex items-center gap-2"
          >
            {isLoggedIn ? <LogOut className="w-5 h-5" /> : <LogIn className="w-5 h-5" />}
            {isLoggedIn ? 'LOGOUT' : 'LOGIN'}
          </button>
        </div>

        <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-8 overflow-hidden mb-4">
          {/* 지도 섹션 [디자인 유지] */}
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-4 relative overflow-hidden">
            <div className="w-full h-full bg-zinc-900 rounded-xl flex items-center justify-center overflow-hidden border border-zinc-800 relative">
              <img src={seoulMapImg} alt="서울 지도" className="w-full h-full object-cover opacity-30 grayscale brightness-75" />
              
              {/* 실시간 마커 렌더링 */}
              {gyms.map((gym) => (
                <button
                  key={gym.guss_number}
                  className="absolute -translate-x-1/2 -translate-y-1/2 group z-20"
                  style={{ 
                    top: POSITIONS[gym.guss_number]?.top || '50%', 
                    left: POSITIONS[gym.guss_number]?.left || '50%' 
                  }}
                  onClick={() => setSelectedGym(gym)}
                >
                  <MapPin 
                    className={`w-10 h-10 transition-all duration-300 drop-shadow-[0_4px_8px_rgba(0,0,0,0.8)] ${
                      selectedGym?.guss_number === gym.guss_number 
                      ? 'text-emerald-400 scale-125 drop-shadow-[0_0_20px_rgba(16,185,129,1)]' 
                      : 'text-emerald-600 opacity-80 group-hover:text-emerald-300 group-hover:scale-110' 
                    }`} 
                    fill={selectedGym?.guss_number === gym.guss_number ? "#10b98166" : "none"}
                  />
                </button>
              ))}
            </div>
          </div>
          
          {/* 정보 패널 (INFO_PANEL) [데이터 연동] */}
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col justify-between">
            {selectedGym ? (
              <div className="animate-in fade-in duration-500">
                <h2 className="text-2xl font-bold text-emerald-400 mb-8 flex items-center gap-3" style={{ fontFamily: 'Orbitron' }}>
                  <Shield className="w-6 h-6" /> INFO_PANEL
                </h2>
                
                <div className="space-y-8">
                  <div className="border-l-4 border-emerald-500 pl-4">
                    <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Center Name</p>
                    <p className="text-2xl font-black">{selectedGym.guss_name}</p>
                  </div>
                  
                  <div className="border-l-4 border-zinc-800 pl-4">
                    <p className="text-zinc-500 text-[10px] uppercase tracking-widest mb-1">Live Status</p>
                    <p className="text-emerald-400 font-bold flex items-center gap-2">
                      <Activity className="w-4 h-4 animate-pulse" /> {selectedGym.guss_status?.toUpperCase()}
                    </p>
                    <p className="text-zinc-400 mt-2">
                      현재 이용 인원: <span className="text-white font-black">{selectedGym.guss_user_count}</span> 명 / {selectedGym.guss_size}
                    </p>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-full opacity-30">
                <Shield className="w-12 h-12 mb-4" />
                <p>지점을 선택해주세요.</p>
              </div>
            )}

            <button 
              onClick={handleDetailView}
              className="w-full py-4 bg-zinc-900 border border-emerald-500/30 rounded-xl text-white font-bold hover:bg-emerald-500/10 hover:border-emerald-500 transition-all"
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