import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, Activity, TrendingUp, Clock, Users, Award, Calendar, Target, Heart, MapPin } from 'lucide-react';

export default function GussPage() {
  const navigate = useNavigate();
  const [showReservationModal, setShowReservationModal] = useState(false);
  const [selectedTime, setSelectedTime] = useState('');
  const [selectedDuration, setSelectedDuration] = useState(1);

  // 임시 로그인 상태 (나중에 전역 상태/localStorage로 관리)
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';

  const userUtilization = 58; // 현재 혼잡도 예시

  const handleReservationClick = () => {
    if (!isLoggedIn) {
      alert('로그인이 필요한 서비스입니다. 로그인 페이지로 이동합니다.');
      navigate('/login');
    } else {
      setShowReservationModal(true);
    }
  };

  const handleReservationConfirm = () => {
    if (!selectedTime) {
      alert('시간대를 선택해주세요!');
      return;
    }
    alert(`🎉 예약이 완료되었습니다!\n시간: ${selectedTime}\n이용시간: ${selectedDuration}시간`);
    setShowReservationModal(false);
  };

  return (
    <div className="min-h-screen bg-black text-white relative overflow-x-hidden">
      {/* 배경 그리드 */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute inset-0" style={{
          backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`,
          backgroundSize: '40px 40px'
        }} />
      </div>

      <div className="relative z-10 p-6 max-w-7xl mx-auto">
        <div className="mb-8 text-center">
          <h1 className="text-5xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400 mb-2"
              style={{ fontFamily: 'Orbitron, sans-serif' }}>GYM STATUS</h1>
          <p className="text-emerald-400">실시간 혼잡도 및 예약 시스템</p>
        </div>

        {/* 혼잡도 막대 그래프 */}
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 mb-8">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-6 h-6 text-emerald-400" />
            <h3 className="text-xl font-bold">현재 실시간 혼잡도</h3>
          </div>
          <div className="relative h-12 bg-zinc-900 rounded-2xl overflow-hidden border border-emerald-500/20">
            <div 
              className="absolute inset-y-0 left-0 bg-gradient-to-r from-emerald-500 via-lime-400 to-yellow-500 transition-all duration-1000"
              style={{ width: `${userUtilization}%` }}
            >
              <div className="absolute inset-0 bg-white/20 animate-pulse" />
            </div>
          </div>
          <div className="flex justify-between mt-4 text-emerald-400 font-bold">
            <span>쾌적</span>
            <span className="text-3xl font-black">{userUtilization}%</span>
            <span className="text-red-500">혼잡</span>
          </div>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* 통계 패널 */}
          <div className="space-y-6">
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
              <div className="flex items-center gap-2 mb-4 text-emerald-400"><Users className="w-5 h-5"/> <span className="font-bold">현재 인원</span></div>
              <p className="text-4xl font-black">47 / 80명</p>
            </div>
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
              <div className="flex items-center gap-2 mb-4 text-emerald-400"><TrendingUp className="w-5 h-5"/> <span className="font-bold">피크 시간대</span></div>
              <p className="text-xl font-bold">18:00 - 21:00</p>
            </div>
          </div>

          {/* 메인 정보 및 예약 버튼 */}
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col justify-between">
            <div>
              <h2 className="text-2xl font-black mb-6 flex items-center gap-2"><MapPin className="text-emerald-400" /> 시설 이용 안내</h2>
              <ul className="space-y-4 text-zinc-400">
                <li className="flex items-center gap-3"><Heart className="w-4 h-4 text-emerald-500"/> 유산소 존: 트레드밀 15대 상시 가동</li>
                <li className="flex items-center gap-3"><Target className="w-4 h-4 text-emerald-500"/> 프리웨이트: 덤벨 최대 50kg 구비</li>
                <li className="flex items-center gap-3"><Clock className="w-4 h-4 text-emerald-500"/> 예약 취소는 1시간 전까지만 가능</li>
              </ul>
            </div>
            
            <div className="mt-12 flex justify-end">
              <button 
                onClick={handleReservationClick}
                className="px-10 py-5 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-2xl text-black font-black text-xl hover:scale-105 transition-all shadow-xl shadow-emerald-500/40 flex items-center gap-3"
              >
                <Calendar className="w-6 h-6" /> 지금 예약하기
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* 예약 모달 */}
      {showReservationModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full">
            <h3 className="text-2xl font-black text-emerald-400 mb-6 text-center">시간대 선택</h3>
            <div className="space-y-6">
              <select 
                value={selectedTime}
                onChange={(e) => setSelectedTime(e.target.value)}
                className="w-full bg-black border-2 border-zinc-800 rounded-xl p-4 text-white focus:border-emerald-500 outline-none"
              >
                <option value="">예약 시간을 선택하세요</option>
                <option value="10:00">10:00 AM</option>
                <option value="14:00">02:00 PM</option>
                <option value="19:00">07:00 PM</option>
              </select>
              <div className="flex gap-4">
                <button onClick={() => setShowReservationModal(false)} className="flex-1 py-4 bg-zinc-900 rounded-xl font-bold">취소</button>
                <button onClick={handleReservationConfirm} className="flex-1 py-4 bg-emerald-500 text-black rounded-xl font-black">예약 완료</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}