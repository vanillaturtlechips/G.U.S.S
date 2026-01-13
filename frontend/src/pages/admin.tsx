import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Package, Calendar, DollarSign, Plus, Edit, Trash2, 
  Shield, Activity, TrendingUp, Search, MapPin
} from 'lucide-react';

// 백엔드 domain.Equipment 구조체와 매칭 (JSON 태그 기준)
interface Equipment {
  id: number;
  gym_id: number;
  name: string;
  category: string;
  quantity: number;
  status: string;
  purchaseDate: string;
}

interface Gym {
  guss_number: number;
  guss_name: string;
  guss_address: string;
  guss_status: string;
  guss_user_count?: number;
}

export default function AdminPage() {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'equipment' | 'reservation' | 'revenue'>('equipment');
  const [showAddModal, setShowAddModal] = useState(false);
  const [newEquipment, setNewEquipment] = useState({ name: '', category: '', quantity: '' });
  
  const [gyms, setGyms] = useState<Gym[]>([]); 
  const [selectedGymId, setSelectedGymId] = useState<number | null>(null); 
  const [equipmentList, setEquipmentList] = useState<Equipment[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const API_BASE = "http://localhost:9000/api"; 
  const token = localStorage.getItem('token');

  // 1. 초기 로드: 체육관 목록
  useEffect(() => {
    fetchGyms();
  }, []);

  const fetchGyms = async () => {
    try {
      const res = await fetch(`${API_BASE}/gyms`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await res.json();
      setGyms(data || []);
      if (data && data.length > 0) setSelectedGymId(data[0].guss_number);
    } catch (err) {
      console.error("체육관 로드 실패", err);
    }
  };

  // 2. 데이터 조회 (기구/예약/매출)
  useEffect(() => {
    if (selectedGymId) fetchTabData();
  }, [selectedGymId, activeTab]);

  const fetchTabData = async () => {
    if (!selectedGymId) return;
    setIsLoading(true);
    try {
      if (activeTab === 'equipment') {
        const res = await fetch(`${API_BASE}/equipments?gym_id=${selectedGymId}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();
        setEquipmentList(data || []);
      }
      // 예약/매출은 나중에 다이나모 연결 예정이므로 현재 로직은 유지하거나 Mock 처리
    } catch (err) {
      console.error("데이터 조회 실패", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddEquipment = async () => {
    if (!selectedGymId) return;
    try {
      const res = await fetch(`${API_BASE}/equipments`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          gym_id: selectedGymId,
          name: newEquipment.name,
          category: newEquipment.category,
          quantity: parseInt(newEquipment.quantity),
          status: 'active',
          purchaseDate: new Date().toISOString().split('T')[0]
        })
      });

      if (res.ok) {
        setShowAddModal(false);
        setNewEquipment({ name: '', category: '', quantity: '' });
        fetchTabData(); 
      }
    } catch (err) {
      alert('기구 추가 실패');
    }
  };

  const currentGym = gyms.find(g => g.guss_number === selectedGymId);

  return (
    <div className="min-h-screen bg-black text-white flex font-sans overflow-hidden">
      {/* 화려한 격자 배경 */}
      <div className="fixed inset-0 opacity-20 pointer-events-none z-0">
        <div className="absolute inset-0" style={{
          backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`,
          backgroundSize: '40px 40px'
        }} />
      </div>

      {/* 사이드바 */}
      <div className="w-80 bg-zinc-950 border-r-2 border-emerald-500/30 p-6 relative z-10 flex flex-col h-screen">
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center shadow-[0_0_15px_rgba(16,185,129,0.5)]">
              <Shield className="w-7 h-7 text-black" strokeWidth={2.5} />
            </div>
            <div>
              <h1 className="text-xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400"
                  style={{ fontFamily: 'Orbitron, sans-serif' }}>GUSS ADMIN</h1>
              <p className="text-xs text-emerald-500">Management System</p>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-full" />
        </div>

        <nav className="space-y-3 mb-8">
          <MenuButton active={activeTab === 'equipment'} onClick={() => setActiveTab('equipment')} icon={<Package className="w-6 h-6" />} label="기구 등록" />
          <MenuButton active={activeTab === 'reservation'} onClick={() => setActiveTab('reservation')} icon={<Calendar className="w-6 h-6" />} label="예약 현황 (로그)" />
          <MenuButton active={activeTab === 'revenue'} onClick={() => setActiveTab('revenue')} icon={<DollarSign className="w-6 h-6" />} label="매출 (로그)" />
        </nav>

        {/* 체육관 선택창 */}
        <div className="bg-zinc-900/50 border-2 border-emerald-500/20 rounded-2xl p-4 mb-6">
          <label className="text-[10px] text-emerald-500 font-black uppercase tracking-[0.2em] mb-3 block">Select Gymnasium</label>
          <div className="relative">
            <select 
              value={selectedGymId || ''} 
              onChange={(e) => setSelectedGymId(Number(e.target.value))}
              className="w-full bg-black border border-emerald-500/30 rounded-xl px-4 py-3 text-sm font-bold text-white focus:outline-none focus:border-emerald-500 transition-all appearance-none cursor-pointer"
            >
              {gyms.map(gym => (
                <option key={gym.guss_number} value={gym.guss_number} className="bg-zinc-900">{gym.guss_name}</option>
              ))}
            </select>
            <Search className="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-emerald-500 pointer-events-none" />
          </div>
        </div>

        {/* 시스템 상태 */}
        <div className="mt-auto bg-zinc-900 border border-emerald-500/30 rounded-xl p-4 font-mono text-[10px] text-zinc-500 space-y-2">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-1"><Activity className="w-3 h-3 text-emerald-400" /> SYSTEM</div>
            <div className="text-emerald-400 flex items-center gap-1"><div className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse" /> ONLINE</div>
          </div>
          <div className="flex justify-between"><span>CONN. USERS</span><span className="text-white">{currentGym?.guss_user_count || 12}P</span></div>
        </div>
      </div>

      {/* 메인 콘텐츠 */}
      <div className="flex-1 relative z-10 overflow-y-auto p-8">
        <div className="grid grid-cols-3 gap-6 mb-8">
          <StatCard icon={<Package className="w-8 h-8 text-emerald-400" />} label="총 기구" value={`${equipmentList.length}종`} />
          <StatCard icon={<Activity className="w-8 h-8 text-emerald-400" />} label="현재 이용객" value={`${currentGym?.guss_user_count || 0}명`} />
          <StatCard icon={<TrendingUp className="w-8 h-8 text-emerald-400" />} label="지점 상태" value={currentGym?.guss_status || "Running"} />
        </div>

        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl overflow-hidden shadow-2xl">
          <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500" />
          
          <div className="p-8">
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-3xl font-black text-white uppercase tracking-tighter" style={{ fontFamily: 'Orbitron' }}>
                  {activeTab === 'equipment' ? '기구 관리' : activeTab === 'reservation' ? '예약 내역' : '매출 로그'}
                </h2>
                <p className="text-emerald-500 font-bold mt-1">{currentGym?.guss_name}</p>
              </div>
              
              {/* 기구 추가 버튼: onClick 이벤트 직접 연결 */}
              {activeTab === 'equipment' && (
                <button 
                  onClick={() => setShowAddModal(true)} 
                  className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-black hover:scale-105 transition-all shadow-[0_0_20px_rgba(16,185,129,0.4)]"
                >
                  <Plus className="w-5 h-5" /> 기구 추가
                </button>
              )}
            </div>

            {isLoading ? (
              <div className="py-20 text-center text-emerald-500 animate-pulse font-bold tracking-widest">DATA SYNCING...</div>
            ) : activeTab === 'equipment' ? (
              <div className="space-y-4">
                {equipmentList.length === 0 ? (
                  <div className="py-20 text-center text-zinc-700 border-2 border-dashed border-zinc-900 rounded-3xl font-bold">등록된 기구가 없습니다.</div>
                ) : (
                  equipmentList.map(item => <EquipmentItem key={item.id} item={item} />)
                )}
              </div>
            ) : activeTab === 'reservation' ? (
              <ReservationTable /> 
            ) : (
              <RevenueTable />
            )}
          </div>
        </div>
      </div>

      {/* 모달 */}
      {showAddModal && (
        <AddModal 
          onClose={() => setShowAddModal(false)} 
          onConfirm={handleAddEquipment} 
          data={newEquipment} 
          setData={setNewEquipment} 
        />
      )}

      <style>{`@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700;900&display=swap');`}</style>
    </div>
  );
}

/* --- 스타일 컴포넌트 --- */

const MenuButton = ({ active, onClick, icon, label }: any) => (
  <button onClick={onClick} className={`w-full flex items-center gap-3 px-5 py-4 rounded-xl font-bold transition-all ${active ? 'bg-gradient-to-r from-emerald-500 to-lime-500 text-black shadow-lg shadow-emerald-500/50' : 'bg-zinc-900 text-emerald-400 hover:bg-zinc-800 border border-emerald-500/20'}`}>
    {icon} <span>{label}</span>
  </button>
);

const StatCard = ({ icon, label, value }: any) => (
  <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6 hover:border-emerald-500/50 transition-all relative overflow-hidden group">
    <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-30 transition-opacity">{icon}</div>
    <div className="relative z-10">
      <div className="flex items-center justify-between mb-3 text-emerald-400">{icon} <TrendingUp className="w-5 h-5 opacity-50" /></div>
      <p className="text-xs text-emerald-500/70 font-black uppercase tracking-widest mb-1">{label}</p>
      <p className="text-4xl font-black text-white">{value}</p>
    </div>
  </div>
);

const EquipmentItem = ({ item }: { item: Equipment }) => (
  <div className="bg-zinc-900/50 border border-emerald-500/20 rounded-2xl p-6 flex items-center justify-between hover:border-emerald-500/40 transition-all group">
    <div className="flex items-center gap-5">
      <div className="w-14 h-14 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center text-black shadow-lg shadow-emerald-500/20">
        <Package size={28} />
      </div>
      <div>
        <h3 className="text-lg font-bold text-white group-hover:text-emerald-400 transition-colors">{item.name}</h3>
        <div className="flex items-center gap-3 mt-1.5">
          <span className="px-2 py-0.5 bg-emerald-500/10 text-emerald-500 text-[10px] font-black rounded border border-emerald-500/20 uppercase tracking-tighter">{item.category}</span>
          <span className="text-zinc-500 text-xs font-bold">{item.quantity}대 보유</span>
          <span className={`text-xs font-black flex items-center gap-1 ${item.status === 'active' ? 'text-lime-500' : 'text-amber-500'}`}>
            <div className={`w-1.5 h-1.5 rounded-full ${item.status === 'active' ? 'bg-lime-500' : 'bg-amber-500 animate-pulse'}`} /> 
            {item.status === 'active' ? '정상' : '점검중'}
          </span>
        </div>
      </div>
    </div>
    <div className="flex gap-2">
      <button className="p-3 bg-zinc-800 hover:bg-emerald-500/20 border border-emerald-500/20 rounded-xl transition-all text-emerald-400"><Edit size={18} /></button>
      <button className="p-3 bg-zinc-800 hover:bg-red-500/20 border border-red-500/20 rounded-xl transition-all text-red-500"><Trash2 size={18} /></button>
    </div>
  </div>
);

// Mock 예약 테이블
const ReservationTable = () => (
  <div className="overflow-x-auto">
    <table className="w-full text-left">
      <thead>
        <tr className="text-emerald-500 font-black border-b border-emerald-500/30 text-xs uppercase tracking-widest">
          <th className="p-4">User</th><th className="p-4">Time</th><th className="p-4">Status</th>
        </tr>
      </thead>
      <tbody className="text-zinc-400 text-sm">
        {[1,2,3].map(i => (
          <tr key={i} className="border-b border-zinc-900 hover:bg-zinc-900/30 transition-colors">
            <td className="p-4 font-bold text-white">Guest_{i}</td>
            <td className="p-4">2026-01-13 18:30</td>
            <td className="p-4"><span className="px-2 py-1 bg-emerald-500/10 text-emerald-400 rounded text-[10px] font-black border border-emerald-500/20">CONFIRMED</span></td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

// Mock 매출 테이블
const RevenueTable = () => (
  <div className="space-y-4">
    <div className="p-6 bg-emerald-500/5 border border-emerald-500/20 rounded-2xl flex justify-between items-center">
      <span className="text-emerald-500 font-black uppercase tracking-widest">Today Total</span>
      <span className="text-3xl font-black text-white">500,000 KRW</span>
    </div>
    <div className="grid grid-cols-2 gap-4">
      {[1,2,3,4].map(i => (
        <div key={i} className="p-4 bg-zinc-900/50 border border-zinc-800 rounded-xl flex justify-between">
          <span className="text-zinc-500 text-xs font-bold">Transaction_{i}</span>
          <span className="text-emerald-400 font-black">125,000원</span>
        </div>
      ))}
    </div>
  </div>
);

const AddModal = ({ onClose, onConfirm, data, setData }: any) => (
  <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-xl">
    <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full shadow-[0_0_50px_rgba(16,185,129,0.2)]">
      <h3 className="text-2xl font-black text-white mb-6 text-center tracking-widest uppercase" style={{ fontFamily: 'Orbitron' }}>New Equipment</h3>
      <div className="space-y-4">
        <input placeholder="이름 (예: 트레드밀)" className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-white focus:border-emerald-500 outline-none" 
               onChange={e => setData({...data, name: e.target.value})} />
        <input placeholder="카테고리 (예: 유산소)" className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-white focus:border-emerald-500 outline-none" 
               onChange={e => setData({...data, category: e.target.value})} />
        <input type="number" placeholder="수량" className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-white focus:border-emerald-500 outline-none" 
               onChange={e => setData({...data, quantity: e.target.value})} />
        <div className="flex gap-3 mt-8">
          <button onClick={onClose} className="flex-1 py-4 bg-zinc-800 text-zinc-400 font-bold rounded-xl">CANCEL</button>
          <button onClick={onConfirm} className="flex-1 py-4 bg-gradient-to-r from-emerald-500 to-lime-500 text-black font-black rounded-xl shadow-lg shadow-emerald-500/30">REGISTER</button>
        </div>
      </div>
    </div>
  </div>
);