import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Package, Calendar, DollarSign, Plus, Edit, Trash2, 
  Shield, Activity, TrendingUp
} from 'lucide-react';

export default function AdminPage() {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'equipment' | 'reservation' | 'revenue'>('equipment');
  const [showAddModal, setShowAddModal] = useState(false);
  const [newEquipment, setNewEquipment] = useState({ name: '', category: '', quantity: '' });
  
  // --- [데이터 상태 관리] ---
  const [equipmentList, setEquipmentList] = useState<any[]>([]);
  const [reservationList, setReservationList] = useState<any[]>([]);
  const [revenueList, setRevenueList] = useState<any[]>([]);
  const [stats, setStats] = useState<any>(null); // 백엔드 HandleDashboard 데이터 저장
  const [isLoading, setIsLoading] = useState(false);

  const API_BASE = "http://localhost:9000/api"; // 백엔드 포트 확인 필요 (Go 서버 포트)
  const token = localStorage.getItem('token');

  // 1. 초기 권한 체크 및 데이터 로드
  useEffect(() => {
    const role = localStorage.getItem('userRole');
    const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';

    // 권한 체크 (보안 강화)
    if (!isLoggedIn || role !== 'ADMIN') {
      alert('관리자 전용 구역입니다. 접근 권한이 없습니다.');
      navigate('/');
      return;
    }
    
    // 통계와 리스트 데이터 로드
    fetchDashboardStats();
    fetchTabData();
  }, [navigate, activeTab]);

  // 2. 백엔드 HandleDashboard로부터 통계 가져오기
  const fetchDashboardStats = async () => {
    try {
      const res = await fetch(`${API_BASE}/dashboard`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await res.json();
      setStats(data);
    } catch (err) {
      console.error("통계 로딩 실패:", err);
    }
  };

  // 3. 탭별 데이터 가져오기 (백엔드 핸들러와 매칭)
  const fetchTabData = async () => {
    setIsLoading(true);
    try {
      const gymId = 1; 
      let endpoint = "";
      
      // 백엔드 main.go에 설정된 라우팅 경로와 일치시켜야 함
      if (activeTab === 'equipment') endpoint = `/equipments?gym_id=${gymId}`;
      else if (activeTab === 'reservation') endpoint = `/reservations?gym_id=${gymId}`;
      else endpoint = `/sales?gym_id=${gymId}`;

      const res = await fetch(`${API_BASE}${endpoint}`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await res.json();
      
      if (activeTab === 'equipment') setEquipmentList(data || []);
      else if (activeTab === 'reservation') setReservationList(data || []);
      else setRevenueList(data || []);
    } catch (err) {
      console.error("데이터 로딩 실패:", err);
    } finally {
      setIsLoading(false);
    }
  };

  // 4. 기구 추가 실행
  const handleAddEquipment = async () => {
    if (!newEquipment.name || !newEquipment.category || !newEquipment.quantity) {
      alert('모든 항목을 입력해주세요!');
      return;
    }

    try {
      const res = await fetch(`${API_BASE}/equipments`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          fk_guss_number: 1,
          equip_name: newEquipment.name,
          equip_category: newEquipment.category,
          equip_quantity: parseInt(newEquipment.quantity),
          equip_status: 'active',
          purchase_date: continentalDate() // 날짜 문자열 전송
        })
      });

      if (res.ok) {
        alert('기구가 성공적으로 등록되었습니다.');
        setShowAddModal(false);
        setNewEquipment({ name: '', category: '', quantity: '' });
        fetchTabData(); 
        fetchDashboardStats(); // 통계 갱신
      } else {
        alert('등록 실패: 서버 응답을 확인하세요.');
      }
    } catch (err) {
      alert('서버 통신 오류가 발생했습니다.');
    }
  };

  // 5. 기구 삭제 실행
  const handleDeleteEquipment = async (id: number) => {
    if (!window.confirm('이 기구를 정말 삭제하시겠습니까?')) return;

    try {
      const res = await fetch(`${API_BASE}/equipments/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) {
        alert('삭제되었습니다.');
        fetchTabData();
        fetchDashboardStats();
      }
    } catch (err) {
      alert('삭제 중 오류가 발생했습니다.');
    }
  };

  const handleEditEquipment = (item: any) => {
    const newName = prompt('수정할 기구 이름을 입력하세요:', item.equip_name);
    if (!newName) return;
    alert(`${newName}으로 수정 로직을 호출합니다. (PUT API 필요)`);
  };

  // 통계 계산 (백엔드 데이터 우선, 없으면 로컬 계산)
  const totalRevenue = revenueList.reduce((sum, item) => sum + (item.amount || 0), 0);
  const displayRevenue = stats ? (stats.total_revenue || totalRevenue) : totalRevenue;

  return (
    <div className="min-h-screen bg-black text-white flex font-sans">
      <div className="fixed inset-0 opacity-20 pointer-events-none">
        <div className="absolute inset-0" style={{
          backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`,
          backgroundSize: '40px 40px'
        }} />
      </div>

      <div className="w-80 bg-zinc-950 border-r-2 border-emerald-500/30 p-6 relative z-10 flex flex-col h-screen sticky top-0">
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center shadow-[0_0_15px_rgba(16,185,129,0.5)]">
              <Shield className="w-7 h-7 text-black" strokeWidth={2.5} />
            </div>
            <div>
              <h1 className="text-xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400"
                  style={{ fontFamily: 'Orbitron, sans-serif' }}>
                GUSS ADMIN
              </h1>
              <p className="text-xs text-emerald-500">Management System</p>
            </div>
          </div>
          <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-full" />
        </div>

        <nav className="flex-1 space-y-3">
          <MenuButton active={activeTab === 'equipment'} onClick={() => setActiveTab('equipment')} icon={<Package className="w-6 h-6" />} label="기구 등록" />
          <MenuButton active={activeTab === 'reservation'} onClick={() => setActiveTab('reservation')} icon={<Calendar className="w-6 h-6" />} label="예약 현황 (로그)" />
          <MenuButton active={activeTab === 'revenue'} onClick={() => setActiveTab('revenue')} icon={<DollarSign className="w-6 h-6" />} label="매출 (로그)" />
        </nav>

        <div className="mt-auto bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
          <div className="flex items-center gap-2 mb-3">
            <Activity className="w-4 h-4 text-emerald-400 animate-pulse" />
            <span className="text-sm font-bold text-emerald-400">System Status</span>
          </div>
          <div className="text-xs text-zinc-500 space-y-2 font-mono">
            <StatusRow label="서버 상태" value={stats ? "Online" : "Offline"} />
            <StatusRow label="현재 접속자" value={stats ? `${stats.active_now}명` : "-"} />
            <div className="flex justify-between"><span>마지막 업데이트</span><span>{stats ? "방금 전" : "대기 중"}</span></div>
          </div>
        </div>
      </div>

      <div className="flex-1 relative z-10 overflow-y-auto">
        <div className="p-8">
          <div className="grid grid-cols-3 gap-6 mb-8">
            <StatCard icon={<Package className="w-8 h-8 text-emerald-400" />} label="총 기구" value={`${equipmentList.length}종`} />
            <StatCard icon={<Calendar className="w-8 h-8 text-emerald-400" />} label="시스템 상태" value={stats?.status || "Running"} />
            <StatCard icon={<DollarSign className="w-8 h-8 text-emerald-400" />} label="오늘 매출" value={`${(displayRevenue / 10000).toFixed(0)}만원`} />
          </div>

          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl overflow-hidden shadow-2xl">
            <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500" />
            {isLoading ? (
              <div className="p-20 text-center text-emerald-500 font-bold animate-pulse">데이터를 동기화 중입니다...</div>
            ) : (
              <>
                {activeTab === 'equipment' && (
                  <div className="p-8">
                    <div className="flex items-center justify-between mb-6">
                      <h2 className="text-3xl font-black text-white">기구 관리</h2>
                      <button onClick={() => setShowAddModal(true)} className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold hover:scale-105 transition-all shadow-lg shadow-emerald-500/50">
                        <Plus className="w-5 h-5" /> 기구 추가
                      </button>
                    </div>
                    <div className="space-y-3">
                      {equipmentList.map(item => (
                        <EquipmentItem 
                          key={item.equip_id} 
                          item={item} 
                          onEdit={() => handleEditEquipment(item)}
                          onDelete={() => handleDeleteEquipment(item.equip_id)} 
                        />
                      ))}
                    </div>
                  </div>
                )}
                {activeTab === 'reservation' && <ReservationTable data={reservationList} />}
                {activeTab === 'revenue' && <RevenueTable data={revenueList} total={displayRevenue} />}
              </>
            )}
          </div>
        </div>
      </div>

      {showAddModal && (
        <AddModal 
          newEquipment={newEquipment} 
          setNewEquipment={setNewEquipment} 
          onClose={() => setShowAddModal(false)} 
          onConfirm={handleAddEquipment} 
        />
      )}
      
      <style>{`@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700;900&display=swap');`}</style>
    </div>
  );
}

// --- 아래는 디자인 컴포넌트들 (기존 UI 로직 유지) ---

const MenuButton = ({ active, onClick, icon, label }: any) => (
  <button onClick={onClick} className={`w-full flex items-center gap-3 px-5 py-4 rounded-xl font-bold transition-all ${active ? 'bg-gradient-to-r from-emerald-500 to-lime-500 text-black shadow-lg shadow-emerald-500/50' : 'bg-zinc-900 text-emerald-400 hover:bg-zinc-800 border border-emerald-500/20'}`}>
    {icon} <span>{label}</span>
  </button>
);

const StatCard = ({ icon, label, value }: any) => (
  <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6 hover:border-emerald-500/50 transition-all group">
    <div className="flex items-center justify-between mb-3">
      {icon} <TrendingUp className="w-5 h-5 text-emerald-400 opacity-50 group-hover:opacity-100" />
    </div>
    <p className="text-sm text-emerald-400 mb-1">{label}</p>
    <p className="text-4xl font-black text-white">{value}</p>
  </div>
);

const EquipmentItem = ({ item, onEdit, onDelete }: any) => (
  <div className="bg-zinc-900/50 border border-emerald-500/20 rounded-2xl p-6 hover:border-emerald-500/40 transition-all">
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-4">
        <div className="w-16 h-16 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center text-black font-bold"><Package size={32} /></div>
        <div>
          <h3 className="text-xl font-bold text-white mb-1">{item.equip_name}</h3>
          <div className="flex items-center gap-3">
            <span className="px-3 py-1 bg-emerald-500/20 text-emerald-400 text-sm font-bold rounded-lg border border-emerald-500/30">{item.equip_category}</span>
            <span className="text-white font-semibold">{item.equip_quantity}대</span>
            <span className={`px-3 py-1 text-sm font-bold rounded-lg border ${item.equip_status === 'active' ? 'bg-lime-500/20 text-lime-400 border-lime-500/30' : 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30'}`}>{item.equip_status === 'active' ? '정상' : '점검중'}</span>
          </div>
        </div>
      </div>
      <div className="flex gap-2">
        <button onClick={onEdit} className="p-3 bg-zinc-800 hover:bg-emerald-500/20 border border-emerald-500/30 rounded-xl transition-all"><Edit className="w-5 h-5 text-emerald-400" /></button>
        <button onClick={onDelete} className="p-3 bg-zinc-800 hover:bg-red-500/20 border border-red-500/30 rounded-xl transition-all"><Trash2 className="w-5 h-5 text-red-400" /></button>
      </div>
    </div>
  </div>
);

const ReservationTable = ({ data }: any) => (
  <div className="p-8 overflow-x-auto">
    <h2 className="text-3xl font-black text-white mb-6">예약 현황 로그</h2>
    <table className="w-full">
      <thead>
        <tr className="border-b-2 border-emerald-500/30 text-emerald-400 font-bold">
          <th className="text-left py-4 px-4">회원명</th><th className="text-left py-4 px-4">연락처</th><th className="text-left py-4 px-4">상태</th>
        </tr>
      </thead>
      <tbody>
        {data.map((item: any) => (
          <tr key={item.revs_number} className="border-b border-emerald-500/10 hover:bg-zinc-900/50 transition-colors">
            <td className="py-4 px-4 text-white font-semibold">{item.user_name}</td>
            <td className="py-4 px-4 text-zinc-400 font-mono text-sm">{item.user_phone}</td>
            <td className="py-4 px-4">
              <span className={`px-3 py-1 rounded-lg text-sm font-bold border ${item.revs_status === 'CONFIRMED' ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' : 'bg-red-500/20 text-red-400 border-red-500/30'}`}>{item.revs_status}</span>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

const RevenueTable = ({ data, total }: any) => (
  <div className="p-8 overflow-x-auto">
    <div className="flex justify-between items-center mb-6">
      <h2 className="text-3xl font-black text-white">매출 로그</h2>
      <div className="px-6 py-3 bg-gradient-to-r from-emerald-500/20 to-lime-500/20 border border-emerald-500/30 rounded-xl">
        <p className="text-sm text-emerald-400">오늘 총 매출</p>
        <p className="text-2xl font-black text-white">{total.toLocaleString()}원</p>
      </div>
    </div>
    <table className="w-full">
      <thead>
        <tr className="border-b-2 border-emerald-500/30 text-emerald-400 font-bold">
          <th className="text-left py-4 px-4">종류</th><th className="text-left py-4 px-4">금액</th><th className="text-left py-4 px-4">일시</th>
        </tr>
      </thead>
      <tbody>
        {data.map((item: any, idx: number) => (
          <tr key={idx} className="border-b border-emerald-500/10 hover:bg-zinc-900/50 transition-colors">
            <td className="py-4 px-4 text-white font-semibold">{item.type}</td>
            <td className="py-4 px-4 text-xl font-black text-emerald-400">{item.amount.toLocaleString()}원</td>
            <td className="py-4 px-4 text-zinc-400 text-sm">{item.date}</td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
);

const AddModal = ({ newEquipment, setNewEquipment, onClose, onConfirm }: any) => (
  <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
    <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full shadow-[0_0_50px_rgba(16,185,129,0.2)]">
      <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-t-3xl mb-6" />
      <h3 className="text-2xl font-black text-white mb-6 text-center uppercase tracking-widest" style={{ fontFamily: 'Orbitron' }}>Add New Equipment</h3>
      <div className="space-y-4">
        <Input label="기구 이름" value={newEquipment.name} onChange={(v:any) => setNewEquipment({...newEquipment, name: v})} placeholder="예: 트레드밀" />
        <Input label="카테고리" value={newEquipment.category} onChange={(v:any) => setNewEquipment({...newEquipment, category: v})} placeholder="예: 유산소" />
        <Input label="수량" value={newEquipment.quantity} onChange={(v:any) => setNewEquipment({...newEquipment, quantity: v})} placeholder="예: 10" />
        <div className="flex gap-4 mt-6">
          <button onClick={onClose} className="flex-1 py-3 bg-zinc-900 border border-zinc-700 rounded-xl text-white font-bold">취소</button>
          <button onClick={onConfirm} className="flex-1 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-black shadow-lg shadow-emerald-500/30">등록하기</button>
        </div>
      </div>
    </div>
  </div>
);

const Input = ({ label, value, onChange, placeholder }: any) => (
  <div>
    <label className="text-sm text-emerald-400 font-bold mb-2 block">{label}</label>
    <input type="text" value={value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 text-white focus:outline-none transition-all" />
  </div>
);

const StatusRow = ({ label, value }: any) => (
  <div className="flex justify-between">
    <span>{label}</span>
    <span className="text-emerald-400 flex items-center gap-1">
      <span className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse" /> {value}
    </span>
  </div>
);

// 날짜 포맷 함수 (YYYY-MM-DD)
function continentalDate() {
  const d = new Date(); // new 키워드는 내부 객체 생성에 필요합니다.
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}