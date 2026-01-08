import { useState } from 'react';
import {
  Package, Calendar, DollarSign, Plus, Edit, Trash2, 
  Search, Download, Eye, Shield, Activity, TrendingUp
} from 'lucide-react';

export default function AdminPage() {
  const [activeTab, setActiveTab] = useState<'equipment' | 'reservation' | 'revenue'>('equipment');
  const [showAddModal, setShowAddModal] = useState(false);
  const [newEquipment, setNewEquipment] = useState({ name: '', category: '', quantity: '' });

  // Mock data
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
      {/* Animated Background */}
      <div className="fixed inset-0 opacity-20 pointer-events-none">
        <div className="absolute inset-0" style={{
          backgroundImage: `
            linear-gradient(to right, #10b981 1px, transparent 1px),
            linear-gradient(to bottom, #10b981 1px, transparent 1px)
          `,
          backgroundSize: '40px 40px'
        }} />
      </div>

      {/* Left Sidebar - Vertical Menu */}
      <div className="w-80 bg-zinc-950 border-r-2 border-emerald-500/30 p-6 relative z-10 flex flex-col h-screen sticky top-0">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="w-12 h-12 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center">
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

        {/* Menu Items */}
        <nav className="flex-1 space-y-3">
          <button
            onClick={() => setActiveTab('equipment')}
            className={`w-full flex items-center gap-3 px-5 py-4 rounded-xl font-bold transition-all ${
              activeTab === 'equipment'
                ? 'bg-gradient-to-r from-emerald-500 to-lime-500 text-black shadow-lg shadow-emerald-500/50'
                : 'bg-zinc-900 text-emerald-400 hover:bg-zinc-800 border border-emerald-500/20'
            }`}
          >
            <Package className="w-6 h-6" />
            <span>기구 등록</span>
          </button>

          <button
            onClick={() => setActiveTab('reservation')}
            className={`w-full flex items-center gap-3 px-5 py-4 rounded-xl font-bold transition-all ${
              activeTab === 'reservation'
                ? 'bg-gradient-to-r from-emerald-500 to-lime-500 text-black shadow-lg shadow-emerald-500/50'
                : 'bg-zinc-900 text-emerald-400 hover:bg-zinc-800 border border-emerald-500/20'
            }`}
          >
            <Calendar className="w-6 h-6" />
            <span>예약 현황 (로그)</span>
          </button>

          <button
            onClick={() => setActiveTab('revenue')}
            className={`w-full flex items-center gap-3 px-5 py-4 rounded-xl font-bold transition-all ${
              activeTab === 'revenue'
                ? 'bg-gradient-to-r from-emerald-500 to-lime-500 text-black shadow-lg shadow-emerald-500/50'
                : 'bg-zinc-900 text-emerald-400 hover:bg-zinc-800 border border-emerald-500/20'
            }`}
          >
            <DollarSign className="w-6 h-6" />
            <span>매출 (로그)</span>
          </button>
        </nav>

        {/* System Status */}
        <div className="mt-auto">
          <div className="bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
            <div className="flex items-center gap-2 mb-3">
              <Activity className="w-4 h-4 text-emerald-400 animate-pulse" />
              <span className="text-sm font-bold text-emerald-400">System Status</span>
            </div>
            <div className="text-xs text-zinc-500 space-y-2">
              <div className="flex justify-between">
                <span>서버 상태</span>
                <span className="text-emerald-400 flex items-center gap-1">
                  <span className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></span>
                  Online
                </span>
              </div>
              <div className="flex justify-between">
                <span>데이터베이스</span>
                <span className="text-emerald-400 flex items-center gap-1">
                  <span className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></span>
                  Connected
                </span>
              </div>
              <div className="flex justify-between">
                <span>마지막 업데이트</span>
                <span className="text-zinc-400">방금 전</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 relative z-10 overflow-y-auto">
        <div className="p-8">
          {/* Stats Bar */}
          <div className="grid grid-cols-3 gap-6 mb-8">
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6 hover:border-emerald-500/50 transition-all">
              <div className="flex items-center justify-between mb-3">
                <Package className="w-8 h-8 text-emerald-400" />
                <TrendingUp className="w-5 h-5 text-emerald-400" />
              </div>
              <p className="text-sm text-emerald-400 mb-1">총 기구</p>
              <p className="text-4xl font-black text-white">{totalEquipment}대</p>
            </div>

            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6 hover:border-emerald-500/50 transition-all">
              <div className="flex items-center justify-between mb-3">
                <Calendar className="w-8 h-8 text-emerald-400" />
                <TrendingUp className="w-5 h-5 text-emerald-400" />
              </div>
              <p className="text-sm text-emerald-400 mb-1">오늘 예약</p>
              <p className="text-4xl font-black text-white">{reservations.length}건</p>
            </div>

            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6 hover:border-emerald-500/50 transition-all">
              <div className="flex items-center justify-between mb-3">
                <DollarSign className="w-8 h-8 text-emerald-400" />
                <TrendingUp className="w-5 h-5 text-emerald-400" />
              </div>
              <p className="text-sm text-emerald-400 mb-1">오늘 매출</p>
              <p className="text-4xl font-black text-white">{(totalRevenue / 10000).toFixed(0)}만원</p>
            </div>
          </div>

          {/* Content based on selected tab */}
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl overflow-hidden">
            <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500" />

            {/* Equipment Tab */}
            {activeTab === 'equipment' && (
              <div className="p-8">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-3xl font-black text-white">기구 관리</h2>
                  <button
                    onClick={() => setShowAddModal(true)}
                    className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 hover:from-emerald-600 hover:to-lime-600 rounded-xl text-black font-bold transition-all shadow-lg shadow-emerald-500/50"
                  >
                    <Plus className="w-5 h-5" />
                    기구 추가
                  </button>
                </div>

                <div className="space-y-3">
                  {equipment.map((item) => (
                    <div
                      key={item.id}
                      className="bg-zinc-900/50 border border-emerald-500/20 rounded-2xl p-6 hover:border-emerald-500/40 transition-all"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <div className="w-16 h-16 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-xl flex items-center justify-center">
                            <Package className="w-8 h-8 text-black" />
                          </div>
                          <div>
                            <h3 className="text-xl font-bold text-white mb-1">{item.name}</h3>
                            <div className="flex items-center gap-3">
                              <span className="px-3 py-1 bg-emerald-500/20 text-emerald-400 text-sm font-bold rounded-lg border border-emerald-500/30">
                                {item.category}
                              </span>
                              <span className="text-white font-semibold">{item.quantity}대</span>
                              <span className={`px-3 py-1 text-sm font-bold rounded-lg border ${
                                item.status === 'active'
                                  ? 'bg-lime-500/20 text-lime-400 border-lime-500/30'
                                  : 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30'
                              }`}>
                                {item.status === 'active' ? '정상' : '점검중'}
                              </span>
                              <span className="text-zinc-500 text-sm">구매일: {item.purchaseDate}</span>
                            </div>
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <button className="p-3 bg-zinc-800 hover:bg-emerald-500/20 border border-emerald-500/30 rounded-xl transition-all group">
                            <Edit className="w-5 h-5 text-emerald-400 group-hover:scale-110 transition-transform" />
                          </button>
                          <button className="p-3 bg-zinc-800 hover:bg-red-500/20 border border-red-500/30 rounded-xl transition-all group">
                            <Trash2 className="w-5 h-5 text-red-400 group-hover:scale-110 transition-transform" />
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Reservation Tab */}
            {activeTab === 'reservation' && (
              <div className="p-8">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-3xl font-black text-white">예약 현황 로그</h2>
                  <div className="flex gap-3">
                    <button className="flex items-center gap-2 px-4 py-2 bg-zinc-900 hover:bg-zinc-800 border border-emerald-500/30 rounded-xl text-emerald-400 font-bold transition-all">
                      <Search className="w-4 h-4" />
                      검색
                    </button>
                    <button className="flex items-center gap-2 px-4 py-2 bg-zinc-900 hover:bg-zinc-800 border border-emerald-500/30 rounded-xl text-emerald-400 font-bold transition-all">
                      <Download className="w-4 h-4" />
                      내보내기
                    </button>
                  </div>
                </div>

                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b-2 border-emerald-500/30">
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">회원명</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">연락처</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">예약 시간</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">날짜</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">상태</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">관리</th>
                      </tr>
                    </thead>
                    <tbody>
                      {reservations.map((item) => (
                        <tr key={item.id} className="border-b border-emerald-500/10 hover:bg-zinc-900/50 transition-colors">
                          <td className="py-4 px-4 text-white font-semibold">{item.member}</td>
                          <td className="py-4 px-4 text-zinc-400 font-mono text-sm">{item.phone}</td>
                          <td className="py-4 px-4 text-white font-semibold">{item.time}</td>
                          <td className="py-4 px-4 text-zinc-400">{item.date}</td>
                          <td className="py-4 px-4">
                            <span className={`px-3 py-1 rounded-lg text-sm font-bold border ${
                              item.status === 'confirmed' ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' :
                              item.status === 'pending' ? 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30' :
                              'bg-red-500/20 text-red-400 border-red-500/30'
                            }`}>
                              {item.status === 'confirmed' ? '확정' :
                               item.status === 'pending' ? '대기' : '취소'}
                            </span>
                          </td>
                          <td className="py-4 px-4">
                            <button className="p-2 bg-zinc-800 hover:bg-emerald-500/20 border border-emerald-500/30 rounded-lg transition-all">
                              <Eye className="w-4 h-4 text-emerald-400" />
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* Summary */}
                <div className="mt-6 grid grid-cols-3 gap-4">
                  <div className="bg-emerald-500/10 border border-emerald-500/30 rounded-xl p-4">
                    <p className="text-sm text-emerald-400 mb-1">확정된 예약</p>
                    <p className="text-2xl font-black text-white">
                      {reservations.filter(r => r.status === 'confirmed').length}건
                    </p>
                  </div>
                  <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-4">
                    <p className="text-sm text-yellow-400 mb-1">대기 중</p>
                    <p className="text-2xl font-black text-white">
                      {reservations.filter(r => r.status === 'pending').length}건
                    </p>
                  </div>
                  <div className="bg-red-500/10 border border-red-500/30 rounded-xl p-4">
                    <p className="text-sm text-red-400 mb-1">취소됨</p>
                    <p className="text-2xl font-black text-white">
                      {reservations.filter(r => r.status === 'cancelled').length}건
                    </p>
                  </div>
                </div>
              </div>
            )}

            {/* Revenue Tab */}
            {activeTab === 'revenue' && (
              <div className="p-8">
                <div className="flex items-center justify-between mb-6">
                  <h2 className="text-3xl font-black text-white">매출 로그</h2>
                  <div className="flex items-center gap-4">
                    <div className="px-6 py-3 bg-gradient-to-r from-emerald-500/20 to-lime-500/20 border border-emerald-500/30 rounded-xl">
                      <p className="text-sm text-emerald-400 mb-1">오늘 총 매출</p>
                      <p className="text-3xl font-black text-white">{totalRevenue.toLocaleString()}원</p>
                    </div>
                    <button className="flex items-center gap-2 px-4 py-2 bg-zinc-900 hover:bg-zinc-800 border border-emerald-500/30 rounded-xl text-emerald-400 font-bold transition-all">
                      <Download className="w-4 h-4" />
                      리포트
                    </button>
                  </div>
                </div>

                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b-2 border-emerald-500/30">
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">회원명</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">상품/서비스</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">금액</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">결제방법</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">일시</th>
                        <th className="text-left py-4 px-4 text-emerald-400 font-bold">상태</th>
                      </tr>
                    </thead>
                    <tbody>
                      {revenue.map((item) => (
                        <tr key={item.id} className="border-b border-emerald-500/10 hover:bg-zinc-900/50 transition-colors">
                          <td className="py-4 px-4 text-white font-semibold">{item.member}</td>
                          <td className="py-4 px-4 text-zinc-400">{item.type}</td>
                          <td className="py-4 px-4 text-2xl font-black text-emerald-400">
                            {item.amount.toLocaleString()}원
                          </td>
                          <td className="py-4 px-4">
                            <span className="px-3 py-1 bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 rounded-lg text-sm font-bold">
                              {item.method}
                            </span>
                          </td>
                          <td className="py-4 px-4 text-zinc-400 text-sm">{item.date}</td>
                          <td className="py-4 px-4">
                            <span className="px-3 py-1 bg-lime-500/20 text-lime-400 border border-lime-500/30 rounded-lg text-sm font-bold">
                              완료
                            </span>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {/* Revenue Chart/Stats */}
                <div className="mt-6 grid grid-cols-4 gap-4">
                  <div className="bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
                    <p className="text-sm text-zinc-400 mb-1">평균 거래액</p>
                    <p className="text-2xl font-black text-white">
                      {revenue.length > 0 
                        ? Math.round(totalRevenue / revenue.length).toLocaleString() 
                        : 0}원
                    </p>
                  </div>
                  <div className="bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
                    <p className="text-sm text-zinc-400 mb-1">총 거래 건수</p>
                    <p className="text-2xl font-black text-white">{revenue.length}건</p>
                  </div>
                  <div className="bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
                    <p className="text-sm text-zinc-400 mb-1">카드 결제</p>
                    <p className="text-2xl font-black text-white">
                      {revenue.filter(r => r.method === '카드').length}건
                    </p>
                  </div>
                  <div className="bg-zinc-900 border border-emerald-500/30 rounded-xl p-4">
                    <p className="text-sm text-zinc-400 mb-1">이번 달 예상</p>
                    <p className="text-2xl font-black text-emerald-400">
                      {Math.round(totalRevenue * 30 / 10000)}만원
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Add Equipment Modal */}
      {showAddModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full shadow-2xl">
            <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-t-3xl mb-6" />
            
            <h3 className="text-2xl font-black text-white mb-6 text-center">새 기구 등록</h3>
            
            <div className="space-y-4">
              <div>
                <label className="text-sm text-emerald-400 font-bold mb-2 block">기구 이름</label>
                <input
                  type="text"
                  value={newEquipment.name}
                  onChange={(e) => setNewEquipment({ ...newEquipment, name: e.target.value })}
                  placeholder="예: 트레드밀"
                  className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 text-white placeholder-zinc-600 focus:outline-none transition-all"
                />
              </div>

              <div>
                <label className="text-sm text-emerald-400 font-bold mb-2 block">카테고리</label>
                <select 
                  value={newEquipment.category}
                  onChange={(e) => setNewEquipment({ ...newEquipment, category: e.target.value })}
                  className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 text-white focus:outline-none transition-all"
                >
                  <option value="">선택하세요</option>
                  <option value="유산소">유산소</option>
                  <option value="웨이트">웨이트</option>
                  <option value="머신">머신</option>
                </select>
              </div>

              <div>
                <label className="text-sm text-emerald-400 font-bold mb-2 block">수량</label>
                <input
                  type="number"
                  value={newEquipment.quantity}
                  onChange={(e) => setNewEquipment({ ...newEquipment, quantity: e.target.value })}
                  placeholder="예: 5"
                  className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl px-4 py-3 text-white placeholder-zinc-600 focus:outline-none transition-all"
                />
              </div>

              <div className="flex gap-3 mt-6">
                <button
                  onClick={() => {
                    setShowAddModal(false);
                    setNewEquipment({ name: '', category: '', quantity: '' });
                  }}
                  className="flex-1 py-3 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 rounded-xl text-white font-bold transition-all"
                >
                  취소
                </button>
                <button
                  onClick={handleAddEquipment}
                  className="flex-1 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 hover:from-emerald-600 hover:to-lime-600 rounded-xl text-black font-bold transition-all shadow-lg shadow-emerald-500/50"
                >
                  등록하기
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;500;600;700;800;900&display=swap');
      `}</style>
    </div>
  );
}