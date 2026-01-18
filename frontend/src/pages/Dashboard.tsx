import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, MapPin, LogIn, LogOut, Activity, Shield } from 'lucide-react';
import seoulMapImg from '../assets/seoul-map.png'; 
import api from '../api/axios'; 
import CongestionChart from '../components/charts/CongestionChart';

const POSITIONS: { [key: number]: { top: string; left: string } } = {
  1: { top: '78%', left: '68%' }, 2: { top: '32%', left: '28%' }, 3: { top: '35%', left: '72%' }, 
};

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [gyms, setGyms] = useState<any[]>([]);
  const [selectedGym, setSelectedGym] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const fetchGyms = async () => {
    try {
      const response = await api.get('/api/gyms');
      setGyms(response.data);
      if (response.data.length > 0 && !selectedGym) setSelectedGym(response.data[0]);
      setLoading(false);
    } catch (error) { setLoading(false); }
  };

  useEffect(() => {
    fetchGyms();
    const interval = setInterval(fetchGyms, 5000); 
    return () => clearInterval(interval);
  }, [selectedGym]);

  if (loading) return <div className="bg-black h-screen flex items-center justify-center text-emerald-400">LOADING...</div>;

  return (
    <div className="min-h-screen bg-black text-white relative overflow-hidden font-sans">
      <div className="relative z-10 p-6 max-w-7xl mx-auto h-screen flex flex-col">
        {/* 상단 네비바 */}
        <div className="bg-zinc-950/80 border-2 border-emerald-500/30 rounded-3xl p-6 mb-8 flex justify-between items-center shadow-lg">
          <h1 className="text-3xl font-black text-emerald-400" style={{ fontFamily: 'Orbitron' }}>GUSS MAP</h1>
          <button onClick={() => { localStorage.clear(); window.location.reload(); }} className="px-6 py-2 bg-emerald-500 rounded-xl text-black font-bold">LOGOUT</button>
        </div>

        <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-8 overflow-hidden">
          {/* 지도 영역 */}
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl relative">
            <div className="w-full h-full bg-zinc-900 relative">
              <img src={seoulMapImg} className="w-full h-full object-cover opacity-20" />
              {gyms.map((gym) => (
                <button key={gym.guss_number} className="absolute -translate-x-1/2 -translate-y-1/2" 
                  style={{ top: POSITIONS[gym.guss_number]?.top || '50%', left: POSITIONS[gym.guss_number]?.left || '50%' }}
                  onClick={() => setSelectedGym(gym)}>
                  <MapPin className={`w-10 h-10 ${selectedGym?.guss_number === gym.guss_number ? 'text-emerald-400' : 'text-emerald-800'}`} />
                </button>
              ))}
            </div>
          </div>
          
          {/* 우측 패널 */}
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col justify-between">
            {selectedGym ? (
              <div className="flex-1 overflow-y-auto">
                <h2 className="text-2xl font-black text-emerald-400 mb-8 italic"><Shield className="inline mr-2"/> INFO_PANEL</h2>
                <div className="space-y-6">
                  <div className="border-l-4 border-emerald-500 pl-4">
                    <p className="text-zinc-500 text-[10px] uppercase">Center Name</p>
                    <p className="text-2xl font-black">{selectedGym.guss_name}</p>
                  </div>
                  <div className="border-l-4 border-zinc-800 pl-4">
                    <p className="text-zinc-500 text-[10px] uppercase mb-4">Congestion Trend</p>
                    {/* 차트 삽입 구역 */}
                    <div className="h-40 bg-black/50 rounded-xl border border-zinc-800 p-2">
                      <CongestionChart gymId={selectedGym.guss_number} color="#10b981" />
                    </div>
                    <p className="mt-4 text-emerald-400 font-bold">인원: {selectedGym.guss_user_count} / {selectedGym.guss_size}</p>
                  </div>
                </div>
              </div>
            ) : <div className="text-center opacity-30">지점을 선택해 주세요.</div>}
            <button onClick={() => navigate(`/guss?gymId=${selectedGym?.guss_number}`)} className="w-full py-5 bg-zinc-900 border border-emerald-500/30 rounded-2xl font-black text-emerald-400 mt-4">VIEW DETAILS</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;