import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Package, Calendar, DollarSign, Plus, Edit, Trash2, 
  Shield, Activity, TrendingUp, Search, MapPin
} from 'lucide-react';

/* --- 인터페이스 정의 --- */
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
  const [editingEquipment, setEditingEquipment] = useState<Equipment | null>(null); // 수정용 상태
  const [newEquipment, setNewEquipment] = useState({ name: '', category: '', quantity: '' });
  
  const [gyms, setGyms] = useState<Gym[]>([]); 
  const [selectedGymId, setSelectedGymId] = useState<number | null>(null); 
  const [equipmentList, setEquipmentList] = useState<Equipment[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const API_BASE = "http://localhost:9000/api"; 
  const token = localStorage.getItem('token');

  useEffect(() => { fetchGyms(); }, []);

  const fetchGyms = async () => {
    try {
      const res = await fetch(`${API_BASE}/gyms`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await res.json();
      setGyms(data || []);
      if (data && data.length > 0) setSelectedGymId(data[0].guss_number);
    } catch (err) { console.error("체육관 로드 실패", err); }
  };

  useEffect(() => { if (selectedGymId) fetchTabData(); }, [selectedGymId, activeTab]);

  const fetchTabData = async () => {
    if (!selectedGymId) return;
    setIsLoading(true);
    try {
      if (activeTab === 'equipment') {
        const res = await fetch(`${API_BASE}/equipments?gymId=${selectedGymId}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();
        setEquipmentList(data || []);
      }
    } catch (err) { console.error("데이터 조회 실패", err); }
    finally { setIsLoading(false); }
  };

  // [기능 추가] 기구 삭제
  const handleDeleteEquipment = async (id: number) => {
    if (!window.confirm("이 기구를 삭제하시겠습니까?")) return;
    try {
      const res = await fetch(`${API_BASE}/equipments/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) fetchTabData();
      else alert("삭제 실패");
    } catch (err) { alert("삭제 통신 오류"); }
  };

  // [기능 추가] 기구 등록 및 수정 통합 핸들러
  const handleSaveEquipment = async () => {
    if (!selectedGymId) return;
    
    // 수정이면 PUT, 등록이면 POST
    const isEdit = !!editingEquipment;
    const url = isEdit ? `${API_BASE}/equipments/${editingEquipment.id}` : `${API_BASE}/equipments`;
    const method = isEdit ? 'PUT' : 'POST';

    try {
      const res = await fetch(url, {
        method: method,
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
          purchaseDate: isEdit ? editingEquipment.purchaseDate : new Date().toISOString().split('T')[0]
        })
      });

      if (res.ok) {
        setShowAddModal(false);
        setEditingEquipment(null);
        setNewEquipment({ name: '', category: '', quantity: '' });
        fetchTabData(); 
      }
    } catch (err) { alert('처리 실패'); }
  };

  const openEditModal = (item: Equipment) => {
    setEditingEquipment(item);
    setNewEquipment({ 
      name: item.name, 
      category: item.category, 
      quantity: item.quantity.toString() 
    });
    setShowAddModal(true);
  };

  const currentGym = gyms.find(g => g.guss_number === selectedGymId);

  return (
    <div className="min-h-screen bg-black text-white flex font-sans overflow-hidden">
      {/* 배경/사이드바 로직 생략 (기존과 동일) */}
      <div className="fixed inset-0 opacity-20 pointer-events-none z-0">
        <div className="absolute inset-0" style={{
          backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`,
          backgroundSize: '40px 40px'
        }} />
      </div>

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
      </div>

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
              
              {activeTab === 'equipment' && (
                <button 
                  onClick={() => {
                    setEditingEquipment(null);
                    setNewEquipment({ name: '', category: '', quantity: '' });
                    setShowAddModal(true);
                  }} 
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
                  equipmentList.map(item => (
                    <EquipmentItem 
                      key={item.id} 
                      item={item} 
                      onEdit={() => openEditModal(item)} 
                      onDelete={() => handleDeleteEquipment(item.id)} 
                    />
                  ))
                )}
              </div>
            ) : null}
          </div>
        </div>
      </div>

      {showAddModal && (
        <AddModal 
          onClose={() => { setShowAddModal(false); setEditingEquipment(null); }} 
          onConfirm={handleSaveEquipment} 
          data={newEquipment} 
          setData={setNewEquipment}
          isEdit={!!editingEquipment}
        />
      )}
    </div>
  );
}

/* --- 하단 컴포넌트에 이벤트 연결 --- */

const EquipmentItem = ({ item, onEdit, onDelete }: { item: Equipment, onEdit: any, onDelete: any }) => (
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
      {/* 이제 onClick 이벤트가 연결되었습니다 */}
      <button 
        onClick={onEdit}
        className="p-3 bg-zinc-800 hover:bg-emerald-500/20 border border-emerald-500/20 rounded-xl transition-all text-emerald-400"
      >
        <Edit size={18} />
      </button>
      <button 
        onClick={onDelete}
        className="p-3 bg-zinc-800 hover:bg-red-500/20 border border-red-500/20 rounded-xl transition-all text-red-500"
      >
        <Trash2 size={18} />
      </button>
    </div>
  </div>
);

// AddModal에도 타이틀 분기 처리 추가
const AddModal = ({ onClose, onConfirm, data, setData, isEdit }: any) => (
  <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-xl">
    <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full shadow-[0_0_50px_rgba(16,185,129,0.2)]">
      <h3 className="text-2xl font-black text-white mb-6 text-center tracking-widest uppercase" style={{ fontFamily: 'Orbitron' }}>
        {isEdit ? 'Update Equipment' : 'New Equipment'}
      </h3>
      <div className="space-y-4">
        <input value={data.name} placeholder="이름" className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-white outline-none" 
               onChange={e => setData({...data, name: e.target.value})} />
        <input value={data.category} placeholder="카테고리" className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-white outline-none" 
               onChange={e => setData({...data, category: e.target.value})} />
        <input value={data.quantity} type="number" placeholder="수량" className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-4 text-white outline-none" 
               onChange={e => setData({...data, quantity: e.target.value})} />
        <div className="flex gap-3 mt-8">
          <button onClick={onClose} className="flex-1 py-4 bg-zinc-800 text-zinc-400 font-bold rounded-xl">CANCEL</button>
          <button onClick={onConfirm} className="flex-1 py-4 bg-gradient-to-r from-emerald-500 to-lime-500 text-black font-black rounded-xl">
            {isEdit ? 'UPDATE' : 'REGISTER'}
          </button>
        </div>
      </div>
    </div>
  </div>
);

// (StatCard, MenuButton 등은 기존과 동일)
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