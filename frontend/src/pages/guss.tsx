import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Activity, Clock, Users, Calendar, Target, Heart, MapPin, ChevronLeft, Shield, Phone } from 'lucide-react';
import api from '../api/axios';
import StatusModal from './StatusModal';
import { QRCodeSVG } from 'qrcode.react';

export default function GussPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const gymId = searchParams.get('gymId');

  const [gymData, setGymData] = useState<any>(null);
  const [showReservationModal, setShowReservationModal] = useState(false);
  const [selectedTime, setSelectedTime] = useState(''); 
  const [showQRModal, setShowQRModal] = useState(false);
  const [qrValue, setQrValue] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [statusModal, setStatusModal] = useState({ isOpen: false, type: 'SUCCESS' as any, title: '', message: '' });

  const fetchDetail = async () => {
    if (!gymId) return;
    try {
      const response = await api.get(`/api/gyms/${gymId}`);
      setGymData(response.data);
    } catch (error) { console.error("데이터 로딩 실패:", error); }
  };

  useEffect(() => {
    fetchDetail();
    const interval = setInterval(fetchDetail, 5000);
    return () => clearInterval(interval);
  }, [gymId]);

  const handleReservationConfirm = async () => {
    if (!selectedTime || isSubmitting) return;

    setIsSubmitting(true);
    try {
      const today = new Date().toISOString().split('T')[0];
      const response = await api.post('/api/reserve', {
        gym_id: parseInt(gymId || '0'),
        visit_time: `${today} ${selectedTime}:00`
      });
      
      setShowReservationModal(false);
      setQrValue(response.data.qr_data); // 백엔드에서 받은 실제 Check-in URL 설정
      setShowQRModal(true);
    } catch (error: any) {
      const errMsg = error.response?.data?.error || '예약 중 오류가 발생했습니다.';
      setStatusModal({ isOpen: true, type: 'ERROR', title: 'FAILED', message: errMsg });
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!gymData) return <div className="min-h-screen bg-black" />;
  const { gym, congestion } = gymData;

  return (
    <div className="min-h-screen bg-black text-white relative font-sans">
      <div className="relative z-10 p-6 max-w-7xl mx-auto">
        <button onClick={() => navigate('/')} className="text-emerald-400 mb-8 flex items-center gap-2"><ChevronLeft /> BACK TO MAP</button>
        <div className="text-center mb-12">
          <h1 className="text-5xl font-black text-emerald-400 italic" style={{ fontFamily: 'Orbitron' }}>{gym?.guss_name}</h1>
          <p className="text-zinc-500 mt-2">{gym?.guss_address}</p>
        </div>

        {/* 혼잡도 카드 */}
        <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 mb-8">
          <div className="flex justify-between items-center mb-4">
            <span className="flex items-center gap-2 text-emerald-400 font-bold"><Activity /> 실시간 혼잡도</span>
            <span className="text-3xl font-black">{Math.round(congestion * 100)}%</span>
          </div>
          <div className="h-4 bg-zinc-900 rounded-full overflow-hidden">
            <div className="h-full bg-emerald-500 transition-all duration-1000" style={{ width: `${congestion * 100}%` }} />
          </div>
        </div>

        <div className="grid lg:grid-cols-2 gap-8">
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-2xl p-8">
            <p className="text-zinc-500 uppercase text-xs mb-2">Current Capacity</p>
            <p className="text-4xl font-black">{gym?.guss_user_count} / {gym?.guss_size}명</p>
          </div>
          <div className="flex items-end justify-end">
            <button onClick={() => setShowReservationModal(true)} className="px-12 py-6 bg-emerald-500 text-black font-black text-2xl rounded-2xl hover:scale-105 transition-all">지금 예약하기</button>
          </div>
        </div>
      </div>

      {/* 예약 모달 */}
      {showReservationModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <div className="bg-zinc-950 border-2 border-emerald-500/30 rounded-3xl p-8 max-w-md w-full">
            <h3 className="text-2xl font-black text-emerald-400 mb-6 text-center italic">방문 시간 선택</h3>
            <select value={selectedTime} onChange={(e)=>setSelectedTime(e.target.value)} className="w-full bg-black border-2 border-zinc-800 p-4 rounded-xl text-white mb-6 outline-none">
              <option value="">시간대를 선택하세요</option>
              <option value="10:00">10:00 AM</option>
              <option value="19:00">07:00 PM</option>
            </select>
            <div className="flex gap-4">
              <button onClick={()=>setShowReservationModal(false)} className="flex-1 py-4 bg-zinc-900 rounded-xl font-bold">취소</button>
              <button onClick={handleReservationConfirm} disabled={isSubmitting} className="flex-1 py-4 bg-emerald-500 text-black rounded-xl font-black">
                {isSubmitting ? '처리 중...' : '예약 완료'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* QR 모달 */}
      {showQRModal && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <div className="bg-zinc-950 border-2 border-emerald-500/50 rounded-3xl p-8 max-w-sm w-full text-center">
            <Shield size={48} className="mx-auto text-emerald-400 mb-4" />
            <h3 className="text-2xl font-black text-emerald-400 mb-2">ENTRY QR CODE</h3>
            <p className="text-zinc-500 mb-6">센터 입구에서 스캔해 주세요.</p>
            <div className="bg-white p-4 rounded-xl inline-block mb-8">
              <QRCodeSVG value={qrValue} size={180} />
            </div>
            <button onClick={()=>setShowQRModal(false)} className="w-full py-4 bg-zinc-900 border border-emerald-500/30 rounded-xl font-black text-white">닫기</button>
          </div>
        </div>
      )}

      <StatusModal isOpen={statusModal.isOpen} type={statusModal.type} title={statusModal.title} message={statusModal.message} onClose={()=>setStatusModal({...statusModal, isOpen:false})} />
    </div>
  );
}