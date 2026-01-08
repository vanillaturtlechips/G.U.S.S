import { useState } from 'react';
import {
  Package, Calendar, DollarSign, Plus, Edit, Trash2, 
  Search, Download, Eye, Shield, Activity, TrendingUp
} from 'lucide-react';

export default function AdminPage() {
  const [activeTab, setActiveTab] = useState<'equipment' | 'reservation' | 'revenue'>('equipment');
  const [showAddModal, setShowAddModal] = useState(false);
  const [newEquipment, setNewEquipment] = useState({ name: '', category: '', quantity: '' });

  // 1. 오빠가 준 목업 데이터 그대로 유지 (기구, 예약, 매출)
  const equipment = [
    { id: 1, name: '트레드밀', category: '유산소', quantity: 10, status: 'active', purchaseDate: '2023-01-15' },
    { id: 2, name: '벤치프레스', category: '웨이트', quantity: 6, status: 'active', purchaseDate: '2023-02-20' },
    { id: 3, name: '스쿼트랙', category: '웨이트', quantity: 4, status: 'maintenance', purchaseDate: '2023-03-10' },
    { id: 4, name: '레그프레스', category: '머신', quantity: 3, status: 'active', purchaseDate: '2023-04-05' },
    { id: 5, name: '일립티컬', category: '유산소', quantity: 6, status: 'active', purchaseDate: '2023-01-15' },
  ];

  const reservations = [
    { id: 1, member: '김민수', phone: '010-1234-5678', time: '14:00-15:00', date: '2024-01-08', status: 'confirmed' },
    { id: 2, member: '이지은', phone: '010-2345-6789', time: '15:00-17:00', date: '2024-01-08', status: 'confirmed' },
    { id: 3, member: '박준혁', phone: '010-3456-7890', time: '16:00-17:00', date: '2024-01-08', status: 'pending' },
    { id: 4, member: '최서연', phone: '010-4567-8901', time: '17:00-18:00', date: '2024-01-08', status: 'cancelled' },
    { id: 5, member: '정우성', phone: '010-5678-9012', time: '18:00-20:00', date: '2024-01-08', status: 'confirmed' },
  ];

  const revenue = [
    { id: 1, member: '홍길동', type: '월회원권', amount: 150000, method: '카드', date: '2024-01-08 14:23', status: 'completed' },
    { id: 2, member: '강수진', type: '3개월권', amount: 400000, method: '계좌이체', date: '2024-01-08 13:15', status: 'completed' },
    { id: 3, member: '윤서준', type: 'PT 10회', amount: 500000, method: '카드', date: '2024-01-08 11:47', status: 'completed' },
    { id: 4, member: '한지민', type: '일일권', amount: 15000, method: '현금', date: '2024-01-08 10:20', status: 'completed' },
    { id: 5, member: '김태희', type: '월회원권', amount: 150000, method: '카드', date: '2024-01-08 09:35', status: 'completed' },
  ];

  const totalRevenue = revenue.reduce((sum, item) => sum + item.amount, 0);
  const totalEquipment = equipment.reduce((sum, item) => sum + item.quantity, 0);

  const handleAddEquipment = () => {
    if (!newEquipment.name || !newEquipment.category || !newEquipment.quantity) {
      alert('모든 항목을 입력해주세요!');
      return;
    }
    alert(`기구 등록 완료!\n이름: ${newEquipment.name}\n카테고리: ${newEquipment.category}\n수량: ${newEquipment.quantity}`);
    setShowAddModal(false);
    setNewEquipment({ name: '', category: '', quantity: '' });
  };

  return (
    <div className="min-h-screen bg-black text-white flex">
      {/* 2. 애니메이션 그리드 배경 */}
      <div className="fixed inset-0 opacity-20 pointer-events-none">
        <div className="absolute inset-0" style={{
          backgroundImage: `
            linear-gradient(to right, #10b981 1px, transparent 1px),
            linear-gradient(to bottom, #10b981 1px, transparent 1px)
          `,
          backgroundSize: '40px 40px'
        }} />
      </div>

      {/* 3. 왼쪽 사이드바 (수직 메뉴) */}
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
          <MenuButton 
            active={activeTab === 'equipment'} 
            onClick={() => setActiveTab('equipment')} 
            icon={<Package className="w-6 h-6" />} 
            label="기구 등록" 
          />
          <MenuButton 
            active={activeTab === 'reservation'} 
            onClick={() => setActiveTab('reservation')} 
            icon={<Calendar className="w-6 h-6" />} 
            label="예약 현황 (로그)" 
          />
          <MenuButton 
            active={activeTab === 'revenue'} 
            onClick={() => setActiveTab('revenue')} 
            icon={<DollarSign className="w-6 h-6" />} 
            label="매출 (로그)" 
          />
        </nav>

        {/* 시스템 상태 표시기 */}
        <div className="mt-auto bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
          <div className="flex items-center gap-2 mb-3">
            <Activity className="w-4 h-4 text-emerald-400 animate-pulse" />
            <span className="text-sm font-bold text-emerald-400">System Status</span>
          </div>
          <div className="text-xs text-zinc-500 space-y-2 font-mono">
            <StatusRow label="서버 상태" value="Online" />
            <StatusRow label="데이터베이스" value="Connected" />
            <div className="flex justify-between"><span>마지막 업데이트</span><span>방금 전</span></div>
          </div>
        </div>
      </div>

      {/* 4. 메인 콘텐츠 영역 */}
      <div className="flex-1 relative z-10 overflow-y-auto">
        <div className="p-8">
          {/* 상단 통계 카드 */}
          <div className="grid grid-cols-3 gap-6 mb-8">
            <StatCard icon={<Package className="w-8 h-8 text-emerald-400" />} label="총 기구" value={`${totalEquipment}대`} />
            <StatCard icon={<Calendar className="w-8 h-8 text-emerald-400" />} label="오늘 예약" value={`${reservations.length}건`} />
            <StatCard icon={<DollarSign className="w-8 h-8 text-emerald-400" />} label="오늘 매출" value={`${(totalRevenue / 10000).toFixed(0)}만원`} />
          </div>

          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl overflow-hidden shadow-2xl">
            <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500" />

            {activeTab === 'equipment' && (
              <div className="p-8">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-3xl font-black text-white">기구 관리</h2>
                  <button onClick={() => setShowAddModal(true)} className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-xl text-black font-bold transition-all shadow-lg shadow-emerald-500/50 hover:scale-105">
                    <Plus className="w-5 h-5" /> 기구 추가
                  </button>
                </div>
                <div className="space-y-3">
                  {equipment.map(item => <EquipmentItem key={item.id} item={item} />)}
                </div>
              </div>
            )}

            {activeTab === 'reservation' && <ReservationTable data={reservations} />}
            {activeTab === 'revenue' && <RevenueTable data={revenue} total={totalRevenue} />}
          </div>
        </div>
      </div>

      {/* 5. 모달 (디자인 오빠 취향대로!) */}
      {showAddModal && <AddModal 
        newEquipment={newEquipment} 
        setNewEquipment={setNewEquipment} 
        onClose={() => setShowAddModal(false)} 
        onConfirm={handleAddEquipment} 
      />}

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700;900&display=swap');
      `}</style>
    </div>
  );
}

// --- 서브 컴포넌트들 (오빠 코드 깔끔하게 정리!) ---

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

const EquipmentItem = ({ item }: any) => (
  <div className="bg-zinc-900/50 border border-emerald-500/20 rounded-2xl p-6 hover:border-emerald-500/40 transition-all">
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-4">
        <div className="w-16 h-16 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center text-black font-bold">
          <Package size={32} />
        </div>
        <div>
          <h3 className="text-xl font-bold text-white mb-1">{item.name}</h3>
          <div className="flex items-center gap-3">
            <span className="px-3 py-1 bg-emerald-500/20 text-emerald-400 text-sm font-bold rounded-lg border border-emerald-500/30">{item.category}</span>
            <span className="text-white font-semibold">{item.quantity}대</span>
            <span className={`px-3 py-1 text-sm font-bold rounded-lg border ${item.status === 'active' ? 'bg-lime-500/20 text-lime-400 border-lime-500/30' : 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30'}`}>
              {item.status === 'active' ? '정상' : '점검중'}
            </span>
          </div>
        </div>
      </div>
      <div className="flex gap-2">
        <button className="p-3 bg-zinc-800 hover:bg-emerald-500/20 border border-emerald-500/30 rounded-xl transition-all"><Edit className="w-5 h-5 text-emerald-400" /></button>
        <button className="p-3 bg-zinc-800 hover:bg-red-500/20 border border-red-500/30 rounded-xl transition-all"><Trash2 className="w-5 h-5 text-red-400" /></button>
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
          <th className="text-left py-4 px-4">회원명</th><th className="text-left py-4 px-4">연락처</th>
          <th className="text-left py-4 px-4">예약 시간</th><th className="text-left py-4 px-4">상태</th>
        </tr>
      </thead>
      <tbody>
        {data.map((item: any) => (
          <tr key={item.id} className="border-b border-emerald-500/10 hover:bg-zinc-900/50 transition-colors">
            <td className="py-4 px-4 text-white font-semibold">{item.member}</td>
            <td className="py-4 px-4 text-zinc-400 font-mono text-sm">{item.phone}</td>
            <td className="py-4 px-4 text-white font-semibold">{item.time}</td>
            <td className="py-4 px-4">
              <span className={`px-3 py-1 rounded-lg text-sm font-bold border ${item.status === 'confirmed' ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' : 'bg-red-500/20 text-red-400 border-red-500/30'}`}>
                {item.status === 'confirmed' ? '확정' : '취소'}
              </span>
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
          <th className="text-left py-4 px-4">회원명</th><th className="text-left py-4 px-4">금액</th><th className="text-left py-4 px-4">일시</th>
        </tr>
      </thead>
      <tbody>
        {data.map((item: any) => (
          <tr key={item.id} className="border-b border-emerald-500/10 hover:bg-zinc-900/50 transition-colors">
            <td className="py-4 px-4 text-white font-semibold">{item.member}</td>
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
      <h3 className="text-2xl font-black text-white mb-6 text-center uppercase tracking-widest">Add New Equipment</h3>
      <div className="space-y-4">
        <Input label="기구 이름" value={newEquipment.name} onChange={(v:any) => setNewEquipment({...newEquipment, name: v})} placeholder="예: 트레드밀" />
        <div className="flex gap-4">
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