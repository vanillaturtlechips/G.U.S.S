import { useState } from 'react';
import { Search, Activity, TrendingUp, Clock, Users, Award, Calendar, Target, Heart, MapPin } from 'lucide-react';

export default function GussPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [showReservationModal, setShowReservationModal] = useState(false);
  const [selectedTime, setSelectedTime] = useState('');
  const [selectedDuration, setSelectedDuration] = useState(1);

  // 현재 인원
  const currentUsers = 47;
  const maxUsers = 80;
  const userUtilization = Math.round((currentUsers / maxUsers) * 100);

  // 혼잡도 상태
  const getCongestionStatus = () => {
    if (userUtilization < 40) return { label: '쾌적', color: 'emerald' };
    if (userUtilization < 70) return { label: '보통', color: 'yellow' };
    return { label: '혼잡', color: 'red' };
  };

  const congestion = getCongestionStatus();

  // 시간대 생성 (00:00 - 23:00)
  const timeSlots = Array.from({ length: 24 }, (_, i) => {
    const hour = i.toString().padStart(2, '0');
    return `${hour}:00`;
  });

  const handleReservation = () => {
    if (!selectedTime) {
      alert('시간대를 선택해주세요!');
      return;
    }
    alert(`예약 완료!\n시간: ${selectedTime}\n이용시간: ${selectedDuration}시간`);
    setShowReservationModal(false);
    setSelectedTime('');
    setSelectedDuration(1);
  };

  return (
    <div className="min-h-screen bg-black text-white relative overflow-x-hidden">
      {/* Animated Background */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute inset-0" style={{
          backgroundImage: `
            linear-gradient(to right, #10b981 1px, transparent 1px),
            linear-gradient(to bottom, #10b981 1px, transparent 1px)
          `,
          backgroundSize: '40px 40px'
        }} />
      </div>

      {/* Gradient Effects */}
      <div className="absolute top-0 left-1/3 w-96 h-96 bg-emerald-500/20 rounded-full blur-3xl animate-pulse" />
      <div className="absolute bottom-0 right-1/3 w-96 h-96 bg-lime-500/20 rounded-full blur-3xl animate-pulse" 
           style={{ animationDelay: '1s' }} />

      <div className="relative z-10 min-h-screen p-6">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-8 text-center">
            <h1 className="text-5xl font-black text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-lime-400 mb-2"
                style={{ fontFamily: 'Orbitron, sans-serif' }}>
              GYM STATUS
            </h1>
            <p className="text-emerald-400">실시간 헬스장 현황 및 예약</p>
          </div>

          {/* Top Section - Search & Congestion Bar */}
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-6 mb-6">
            <div className="grid md:grid-cols-2 gap-6">
              {/* 검색 */}
              <div>
                <div className="flex items-center gap-2 mb-3">
                  <Search className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-bold text-white">정보 검색</h3>
                </div>
                <div className="relative">
                  <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-emerald-500" />
                  <input
                    type="text"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="운동 프로그램, 시설 정보 검색..."
                    className="w-full bg-black border-2 border-zinc-800 focus:border-emerald-500 rounded-xl pl-12 pr-4 py-3 text-white placeholder-zinc-600 focus:outline-none transition-all font-mono"
                  />
                </div>
              </div>

              {/* 혼잡도 막대 */}
              <div>
                <div className="flex items-center gap-2 mb-3">
                  <Activity className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-bold text-white">현재 혼잡도</h3>
                </div>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-emerald-400 font-bold">쾌적</span>
                    <span className={`text-2xl font-black ${
                      congestion.color === 'emerald' ? 'text-emerald-400' :
                      congestion.color === 'yellow' ? 'text-yellow-400' :
                      'text-red-400'
                    }`}>
                      {congestion.label}
                    </span>
                    <span className="text-red-400 font-bold">혼잡</span>
                  </div>
                  <div className="relative h-6 bg-zinc-900 rounded-full overflow-hidden border border-emerald-500/30">
                    <div 
                      className={`absolute inset-y-0 left-0 rounded-full transition-all duration-1000 ${
                        congestion.color === 'emerald' ? 'bg-gradient-to-r from-emerald-500 to-lime-500' :
                        congestion.color === 'yellow' ? 'bg-gradient-to-r from-yellow-500 to-orange-500' :
                        'bg-gradient-to-r from-red-500 to-orange-500'
                      }`}
                      style={{ width: `${userUtilization}%` }}
                    >
                      <div className="absolute inset-0 bg-white/20 animate-pulse" />
                    </div>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-zinc-500">0%</span>
                    <span className="text-white font-bold">{userUtilization}%</span>
                    <span className="text-zinc-500">100%</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Main Content */}
          <div className="grid lg:grid-cols-3 gap-6">
            {/* Left Sidebar - Stats */}
            <div className="space-y-4">
              {/* 현재 이용 인원 */}
              <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
                <div className="flex items-center gap-2 mb-4">
                  <Users className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-bold text-white">현재 인원</h3>
                </div>
                <div className="text-center">
                  <p className="text-5xl font-black text-white mb-2">{currentUsers}</p>
                  <p className="text-zinc-500 mb-4">/ {maxUsers}명</p>
                  <div className="h-2 bg-zinc-900 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-gradient-to-r from-emerald-500 to-lime-500 transition-all duration-1000"
                      style={{ width: `${userUtilization}%` }}
                    />
                  </div>
                </div>
              </div>

              {/* 오늘 통계 */}
              <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
                <div className="flex items-center gap-2 mb-4">
                  <TrendingUp className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-bold text-white">오늘 통계</h3>
                </div>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-zinc-400 text-sm">총 이용자</span>
                    <span className="text-white font-bold">156명</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-zinc-400 text-sm">평균 운동시간</span>
                    <span className="text-white font-bold">87분</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-zinc-400 text-sm">피크 시간</span>
                    <span className="text-white font-bold">18:00-20:00</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-zinc-400 text-sm">현재 대기</span>
                    <span className="text-emerald-400 font-bold">없음</span>
                  </div>
                </div>
              </div>

              {/* 추천 시간대 */}
              <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
                <div className="flex items-center gap-2 mb-4">
                  <Award className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-bold text-white">여유 시간대</h3>
                </div>
                <div className="space-y-2">
                  {[
                    { time: '06:00 - 09:00', status: '매우 쾌적' },
                    { time: '14:00 - 16:00', status: '쾌적' },
                    { time: '21:00 - 23:00', status: '보통' }
                  ].map((slot, idx) => (
                    <div key={idx} className="p-3 bg-emerald-500/10 border border-emerald-500/30 rounded-xl">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <Clock className="w-4 h-4 text-emerald-400" />
                          <span className="text-emerald-400 font-semibold text-sm">{slot.time}</span>
                        </div>
                        <span className="text-xs text-zinc-500">{slot.status}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* 운영 시간 */}
              <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-6">
                <div className="flex items-center gap-2 mb-4">
                  <Target className="w-5 h-5 text-emerald-400" />
                  <h3 className="text-lg font-bold text-white">운영 시간</h3>
                </div>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-zinc-400">평일</span>
                    <span className="text-white font-bold">06:00 - 23:00</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-zinc-400">주말</span>
                    <span className="text-white font-bold">08:00 - 22:00</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Center & Right - Main Info Area */}
            <div className="lg:col-span-2">
              <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8 min-h-[600px] flex flex-col">
                <h2 className="text-3xl font-black text-white mb-6 text-center">헬스장 이용 안내</h2>

                {/* Facility Info */}
                <div className="flex-1 space-y-6">
                  {/* 시설 정보 */}
                  <div className="bg-zinc-900/50 border border-emerald-500/20 rounded-2xl p-6">
                    <h3 className="text-xl font-bold text-emerald-400 mb-4 flex items-center gap-2">
                      <MapPin className="w-5 h-5" />
                      시설 정보
                    </h3>
                    <div className="grid md:grid-cols-2 gap-4 text-sm">
                      <div className="space-y-2">
                        <div className="flex items-center gap-2">
                          <span className="text-emerald-400">•</span>
                          <span className="text-white">유산소 기구 24대</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-emerald-400">•</span>
                          <span className="text-white">웨이트 존 30평</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-emerald-400">•</span>
                          <span className="text-white">머신 기구 20대</span>
                        </div>
                      </div>
                      <div className="space-y-2">
                        <div className="flex items-center gap-2">
                          <span className="text-emerald-400">•</span>
                          <span className="text-white">락커룸 & 샤워실</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-emerald-400">•</span>
                          <span className="text-white">무료 Wi-Fi</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-emerald-400">•</span>
                          <span className="text-white">주차 가능</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* 프로그램 */}
                  <div className="bg-zinc-900/50 border border-emerald-500/20 rounded-2xl p-6">
                    <h3 className="text-xl font-bold text-emerald-400 mb-4 flex items-center gap-2">
                      <Heart className="w-5 h-5" />
                      운동 프로그램
                    </h3>
                    <div className="grid md:grid-cols-2 gap-3">
                      {[
                        { name: 'PT (개인 트레이닝)', time: '사전 예약' },
                        { name: '그룹 필라테스', time: '월/수/금 19:00' },
                        { name: '요가 클래스', time: '화/목 18:00' },
                        { name: '크로스핏', time: '토 10:00' }
                      ].map((program, idx) => (
                        <div key={idx} className="p-3 bg-emerald-500/10 rounded-xl border border-emerald-500/30">
                          <p className="text-white font-semibold text-sm">{program.name}</p>
                          <p className="text-emerald-400 text-xs mt-1">{program.time}</p>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* 이용 안내 */}
                  <div className="bg-gradient-to-r from-emerald-500/10 to-lime-500/10 border border-emerald-500/30 rounded-2xl p-6">
                    <h3 className="text-lg font-bold text-white mb-3">예약 및 이용 안내</h3>
                    <div className="space-y-2 text-sm text-emerald-400/90">
                      <div>• 예약은 00시부터 23시까지 가능합니다</div>
                      <div>• 최소 1시간, 최대 4시간까지 이용 가능합니다</div>
                      <div>• 예약 시간 10분 전까지 입장해주세요</div>
                      <div>• 예약 취소는 이용 시간 1시간 전까지 가능합니다</div>
                    </div>
                  </div>
                </div>

                {/* 예약 버튼 - 우측 하단 */}
                <div className="flex justify-end mt-6">
                  <button
                    onClick={() => setShowReservationModal(true)}
                    className="flex items-center gap-3 px-8 py-4 bg-gradient-to-r from-emerald-500 to-lime-500 hover:from-emerald-600 hover:to-lime-600 rounded-xl text-black font-black text-lg transition-all transform hover:scale-105 shadow-lg shadow-emerald-500/50"
                  >
                    <Calendar className="w-6 h-6" />
                    시간대 예약하기
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Reservation Modal */}
      {showReservationModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full shadow-2xl">
            <div className="h-1 bg-gradient-to-r from-emerald-500 to-lime-500 rounded-t-3xl mb-6" />
            
            <div className="text-center mb-6">
              <div className="inline-block p-4 bg-gradient-to-br from-emerald-500 to-lime-500 rounded-2xl mb-4">
                <Calendar className="w-10 h-10 text-black" />
              </div>
              <h3 className="text-2xl font-black text-white mb-2">시간대 예약</h3>
              <p className="text-emerald-400">원하시는 시간과 이용 시간을 선택하세요</p>
            </div>

            <div className="space-y-4 mb-6">
              {/* 예약 시간 선택 */}
              <div className="bg-zinc-900 rounded-xl p-4">
                <label className="text-zinc-400 text-sm block mb-2">예약 시작 시간</label>
                <select 
                  value={selectedTime}
                  onChange={(e) => setSelectedTime(e.target.value)}
                  className="w-full bg-black border border-emerald-500/30 rounded-lg px-4 py-3 text-white font-bold focus:outline-none focus:border-emerald-500"
                >
                  <option value="">시간 선택</option>
                  {timeSlots.map((time) => (
                    <option key={time} value={time}>{time}</option>
                  ))}
                </select>
              </div>

              {/* 이용 시간 선택 */}
              <div className="bg-zinc-900 rounded-xl p-4">
                <label className="text-zinc-400 text-sm block mb-2">이용 시간</label>
                <div className="grid grid-cols-4 gap-2">
                  {[1, 2, 3, 4].map((hour) => (
                    <button
                      key={hour}
                      onClick={() => setSelectedDuration(hour)}
                      className={`py-3 rounded-lg font-bold transition-all ${
                        selectedDuration === hour
                          ? 'bg-gradient-to-r from-emerald-500 to-lime-500 text-black'
                          : 'bg-zinc-800 text-white hover:bg-zinc-700'
                      }`}
                    >
                      {hour}시간
                    </button>
                  ))}
                </div>
              </div>

              {/* 예약 정보 요약 */}
              <div className="bg-emerald-500/10 border border-emerald-500/30 rounded-xl p-4">
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-zinc-400">시작 시간</span>
                    <span className="text-white font-bold">{selectedTime || '-'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-zinc-400">이용 시간</span>
                    <span className="text-white font-bold">{selectedDuration}시간</span>
                  </div>
                  {selectedTime && (() => {
                    const startHour = parseInt(selectedTime.split(':')[0]);
                    const endHour = (startHour + selectedDuration) % 24;
                    return (
                      <div className="flex justify-between">
                        <span className="text-zinc-400">종료 시간</span>
                        <span className="text-emerald-400 font-bold">
                          {String(endHour).padStart(2, '0')}:00
                          {startHour + selectedDuration >= 24 && ' (익일)'}
                        </span>
                      </div>
                    );
                  })()}
                </div>
              </div>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => {
                  setShowReservationModal(false);
                  setSelectedTime('');
                  setSelectedDuration(1);
                }}
                className="flex-1 py-3 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 rounded-xl text-white font-bold transition-all"
              >
                취소
              </button>
              <button
                onClick={handleReservation}
                className="flex-1 py-3 bg-gradient-to-r from-emerald-500 to-lime-500 hover:from-emerald-600 hover:to-lime-600 rounded-xl text-black font-bold transition-all shadow-lg shadow-emerald-500/50"
              >
                예약 확인
              </button>
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