import { useState } from 'react';
import { Package, Calendar, DollarSign, Plus, Edit, Trash2, Shield, Activity, TrendingUp } from 'lucide-react';

export default function AdminPage() {
  const [activeTab, setActiveTab] = useState<'equipment' | 'reservation' | 'revenue'>('equipment');

  return (
    <div className="min-h-screen bg-black text-white flex">
      {/* 사이드바 */}
      <div className="w-72 bg-zinc-950 border-r-2 border-emerald-500/30 p-8 sticky top-0 h-screen">
        <div className="flex items-center gap-3 mb-12">
          <Shield className="w-10 h-10 text-emerald-500" />
          <h1 className="text-xl font-black text-emerald-400" style={{ fontFamily: 'Orbitron' }}>GUSS ADMIN</h1>
        </div>

        <nav className="space-y-4">
          <button 
            onClick={() => setActiveTab('equipment')}
            className={`w-full flex items-center gap-4 p-4 rounded-xl font-bold transition-all ${activeTab === 'equipment' ? 'bg-emerald-500 text-black shadow-lg shadow-emerald-500/30' : 'bg-zinc-900 text-emerald-400'}`}
          >
            <Package /> 기구 등록
          </button>
          <button 
            onClick={() => setActiveTab('reservation')}
            className={`w-full flex items-center gap-4 p-4 rounded-xl font-bold transition-all ${activeTab === 'reservation' ? 'bg-emerald-500 text-black shadow-lg shadow-emerald-500/30' : 'bg-zinc-900 text-emerald-400'}`}
          >
            <Calendar /> 예약 현황
          </button>
          <button 
            onClick={() => setActiveTab('revenue')}
            className={`w-full flex items-center gap-4 p-4 rounded-xl font-bold transition-all ${activeTab === 'revenue' ? 'bg-emerald-500 text-black shadow-lg shadow-emerald-500/30' : 'bg-zinc-900 text-emerald-400'}`}
          >
            <DollarSign /> 매출 로그
          </button>
        </nav>
      </div>

      {/* 메인 콘텐츠 */}
      <div className="flex-1 p-12">
        <div className="flex justify-between items-center mb-12">
          <h2 className="text-4xl font-black text-white">{activeTab.toUpperCase()} MANAGEMENT</h2>
          <div className="flex items-center gap-4 bg-zinc-900 p-4 rounded-2xl border border-emerald-500/20">
            <Activity className="text-emerald-400 animate-pulse" />
            <span className="font-bold text-emerald-400">SYSTEM ONLINE</span>
          </div>
        </div>

        {/* 대시보드 요약 */}
        <div className="grid grid-cols-3 gap-8 mb-12">
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8">
            <p className="text-emerald-400 text-sm font-bold mb-2">Total Equipment</p>
            <p className="text-4xl font-black">24 Units</p>
          </div>
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8">
            <p className="text-emerald-400 text-sm font-bold mb-2">Daily Revenue</p>
            <p className="text-4xl font-black">₩ 1,250,000</p>
          </div>
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8">
            <p className="text-emerald-400 text-sm font-bold mb-2">Active Users</p>
            <p className="text-4xl font-black">47 Members</p>
          </div>
        </div>

        {/* 탭 내용 (목업 테이블) */}
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl overflow-hidden">
          <table className="w-full text-left">
            <thead className="bg-emerald-500/10 text-emerald-400 font-bold border-b-2 border-emerald-500/30">
              <tr>
                <th className="p-6">CATEGORY</th>
                <th className="p-6">NAME</th>
                <th className="p-6">STATUS</th>
                <th className="p-6">ACTION</th>
              </tr>
            </thead>
            <tbody className="text-zinc-400">
              <tr className="border-b border-zinc-900">
                <td className="p-6">Cardio</td>
                <td className="p-6 text-white font-bold">Treadmill XL-500</td>
                <td className="p-6"><span className="px-3 py-1 bg-emerald-500/20 text-emerald-400 rounded-lg text-xs">Active</span></td>
                <td className="p-6 flex gap-4"><Edit className="w-5 h-5 cursor-pointer"/> <Trash2 className="w-5 h-5 text-red-500 cursor-pointer" /></td>
              </tr>
              {/* 추가 데이터... */}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}