import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { 
  Activity, TrendingUp, Clock, Users, 
  Calendar, Target, Heart, MapPin, ChevronLeft 
} from 'lucide-react';
import api from '../api/axios'; // 설정한 axios 인스턴스

export default function GussPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const gymId = searchParams.get('gymId'); // URL 파라미터에서 gymId 추출

  // 상태 관리 (초기값 null)
  const [gymData, setGymData] = useState<any>(null);
  const [showReservationModal, setShowReservationModal] = useState(false);
  const [selectedTime, setSelectedTime] = useState('');

  // 로그인 상태 확인
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';

  /**
   * [실시간 데이터 패칭]
   * 백엔드의 { "gym": {...}, "congestion": 0 } 구조를 처리합니다.
   */
  const fetchDetail = async () => {
    if (!gymId) return;
    try {
      const response = await api.get(`/api/gyms/${gymId}`);
      setGymData(response.data);
    } catch (error) {
      console.error("데이터 로딩 실패:", error);
    }
  };

  useEffect(() => {
    fetchDetail();
    const interval = setInterval(fetchDetail, 5000); // 5초마다 실시간 갱신
    return () => clearInterval(interval);
  }, [gymId]);

  // 예약 신청 함수
  const handleReservationConfirm = async () => {
    if (!selectedTime) {
      alert('시간대를 선택해주세요!');
      return;
    }

    try {
      // 백엔드 예약 API 호출
      await api.post('/api/reserve', {
        fk_guss_number: parseInt(gymId || '0')
      });
      
      alert(`🎉 예약이 완료되었습니다!\n방문 예정 시간: ${selectedTime}`);
      setShowReservationModal(false);
      fetchDetail(); // 예약 후 즉시 인원수 업데이트
    } catch (error: any) {
      if (error.response?.status === 401) {
        alert('로그인이 필요한 서비스입니다. 로그인 페이지로 이동합니다.');
        navigate('/login');
      } else {
        alert('예약 처리 중 오류가 발생했습니다.');
      }
    }
  };

  // 로딩 상태 디자인
  if (!gymData) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center text-emerald-400 font-black tracking-widest">
        LOADING GYM DATA...
      </div>
    );
  }

  // 백엔드에서 계산해준 혼잡도 수치
  const utilization = Math.round(gymData.congestion * 100) || 0;
  // 실제 지점 데이터
  const gym = gymData.gym;

  return (
    <div className="min-h-screen bg-black text-white relative overflow-x-hidden font-sans">
      {/* 배경 그리드 (기존 디자인 유지) */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute inset-0" style={{
          backgroundImage: `linear-gradient(to right, #10b981 1px, transparent 1px), linear-gradient(to bottom, #10b981 1px, transparent 1px)`,
          backgroundSize: '40px 40px'
        }} />
      </div>

      <div className="relative z-10 p-6 max-w-7xl mx-auto">
        {/* 상단 네비게이션 */}
        <button 
          onClick={() => navigate('/')}
          className="flex items-center gap-2 text-emerald-400 hover:text-white transition-colors mb-8 font-bold"
        >
          <ChevronLeft className="w-5 h-5" /> BACK TO MAP
        </button>

        <div className="mb-8 text-center">
          <h1 className="text-5xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400 mb-2"
              style={{ fontFamily: 'Orbitron, sans-serif' }}>
            {gym?.guss_name?.toUpperCase() || "GYM STATUS"}
          </h1>
          <p className="text-emerald-400">실시간 혼잡도 및 예약 시스템</p>
        </div>

        {/* 실시간 혼잡도 그래프 (데이터 연동) */}
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 mb-8 shadow-[0_0_50px_rgba(16,185,129,0.05)]">
          <div className="flex items-center gap-3 mb-6">
            <Activity className="w-6 h-6 text-emerald-400" />
            <h3 className="text-xl font-bold">현재 실시간 혼잡도</h3>
          </div>
          <div className="relative h-12 bg-zinc-900 rounded-2xl overflow-hidden border border-emerald-500/20">
            <div 
              className="absolute inset-y-0 left-0 bg-gradient-to-r from-emerald-500 via-lime-400 to-yellow-500 transition-all duration-1000"
              style={{ width: `${utilization}%` }}
            >
              <div className="absolute inset-0 bg-white/20 animate-pulse" />
            </div>
          </div>
          <div className="flex justify-between mt-4 text-emerald-400 font-bold">
            <span>쾌적</span>
            <span className="text-3xl font-black">{utilization}%</span>
            <span className="text-red-500">혼잡</span>
          </div>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* 통계 패널 */}
          <div className="space-y-6">
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
              <div className="flex items-center gap-2 mb-4 text-emerald-400">
                <Users className="w-5 h-5"/> <span className="font-bold">현재 인원</span>
              </div>
              <p className="text-4xl font-black">
                {gym?.guss_user_count} / {gym?.guss_size}명
              </p>
            </div>
            <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
              <div className="flex items-center gap-2 mb-4 text-emerald-400">
                <TrendingUp className="w-5 h-5"/> <span className="font-bold">피크 시간대</span>
              </div>
              <p className="text-xl font-bold">18:00 - 21:00</p>
            </div>
          </div>

          {/* 시설 정보 및 예약 섹션 */}
          <div className="lg:col-span-2 bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 flex flex-col justify-between">
            <div>
              <h2 className="text-2xl font-black mb-6 flex items-center gap-2">
                <MapPin className="text-emerald-400" /> 시설 이용 안내
              </h2>
              <p className="text-zinc-400 mb-6 font-bold">{gym?.guss_address}</p>
              <ul className="space-y-4 text-zinc-400">
                <li className="flex items-center gap-3"><Heart className="w-4 h-4 text-emerald-500"/> 유산소 존: 트레드밀 15대 상시 가동</li>
                <li className="flex items-center gap-3"><Target className="w-4 h-4 text-emerald-500"/> 프리웨이트: 덤벨 최대 50kg 구비</li>
                <li className="flex items-center gap-3"><Clock className="w-4 h-4 text-emerald-500"/> 예약 취소는 1시간 전까지만 가능</li>
              </ul>
            </div>
            
            <div className="mt-12 flex justify-end">
              <button 
                onClick={() => setShowReservationModal(true)}
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
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full shadow-2xl">
            <h3 className="text-2xl font-black text-emerald-400 mb-6 text-center">방문 예정 시간</h3>
            <div className="space-y-6">
              <select 
                value={selectedTime}
                onChange={(e) => setSelectedTime(e.target.value)}
                className="w-full bg-black border-2 border-zinc-800 rounded-xl p-4 text-white focus:border-emerald-500 outline-none appearance-none"
              >
                <option value="">시간대를 선택하세요</option>
                <option value="10:00">10:00 AM</option>
                <option value="14:00">02:00 PM</option>
                <option value="19:00">07:00 PM</option>
              </select>
              <div className="flex gap-4">
                <button 
                  onClick={() => setShowReservationModal(false)} 
                  className="flex-1 py-4 bg-zinc-900 rounded-xl font-bold hover:bg-zinc-800 transition-all"
                >
                  취소
                </button>
                <button 
                  onClick={handleReservationConfirm} 
                  className="flex-1 py-4 bg-emerald-500 text-black rounded-xl font-black hover:bg-emerald-400 transition-all"
                >
                  예약 완료
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}