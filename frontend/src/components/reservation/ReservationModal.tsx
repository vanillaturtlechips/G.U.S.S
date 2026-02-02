// src/components/reservation/ReservationModal.tsx
import React, { useState } from 'react';
import { X, Clock, CheckCircle } from 'lucide-react';

interface ReservationModalProps {
  gymName: string;
  onClose: () => void;
  onConfirm: (time: string) => void;
}

const ReservationModal: React.FC<ReservationModalProps> = ({ gymName, onClose, onConfirm }) => {
  const [selectedTime, setSelectedTime] = useState<string | null>(null);

  // 00:00 ~ 23:30까지 30분 단위 시간 생성
  const timeSlots = Array.from({ length: 48 }, (_, i) => {
    const hour = Math.floor(i / 2).toString().padStart(2, '0');
    const min = (i % 2 === 0 ? '00' : '30');
    return `${hour}:${min}`;
  });

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
      {/* 배경 블러 처리 */}
      <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={onClose} />
      
      {/* 모달 창 [디자인 유지] */}
      <div className="relative w-full max-w-lg bg-zinc-950 border-2 border-emerald-500/50 rounded-3xl shadow-[0_0_50px_rgba(16,185,129,0.2)] overflow-hidden animate-in zoom-in duration-300">
        <div className="p-6 border-b border-emerald-500/20 flex justify-between items-center">
          <h2 className="text-xl font-black text-emerald-400 italic uppercase" style={{ fontFamily: 'Orbitron' }}>
            Reserve_Session
          </h2>
          <button onClick={onClose} className="text-zinc-500 hover:text-white transition-colors">
            <X className="w-6 h-6" />
          </button>
        </div>

        <div className="p-6">
          <p className="text-zinc-400 mb-2 text-xs uppercase tracking-widest">Selected Center</p>
          <p className="text-2xl font-black mb-6">{gymName}</p>

          <p className="text-zinc-400 mb-4 text-xs uppercase tracking-widest flex items-center gap-2">
            <Clock className="w-4 h-4" /> Select Time Slot (30min unit)
          </p>

          {/* 시간 선택 그리드 */}
          <div className="grid grid-cols-4 gap-2 max-h-60 overflow-y-auto pr-2 custom-scrollbar">
            {timeSlots.map((time) => (
              <button
                key={time}
                onClick={() => setSelectedTime(time)}
                className={`py-2 rounded-lg font-bold transition-all border ${
                  selectedTime === time
                    ? 'bg-emerald-500 text-black border-emerald-500'
                    : 'bg-zinc-900 text-zinc-400 border-zinc-800 hover:border-emerald-500/50'
                }`}
              >
                {time}
              </button>
            ))}
          </div>
        </div>

        <div className="p-6 bg-zinc-900/50 flex gap-3">
          <button 
            onClick={onClose}
            className="flex-1 py-4 bg-zinc-800 rounded-xl font-bold hover:bg-zinc-700 transition-all"
          >
            CANCEL
          </button>
          <button 
            disabled={!selectedTime}
            onClick={() => selectedTime && onConfirm(selectedTime)}
            className={`flex-1 py-4 rounded-xl font-bold flex items-center justify-center gap-2 transition-all ${
              selectedTime 
                ? 'bg-emerald-500 text-black hover:scale-105' 
                : 'bg-zinc-700 text-zinc-500 cursor-not-allowed'
            }`}
          >
            <CheckCircle className="w-5 h-5" /> CONFIRM
          </button>
        </div>
      </div>
    </div>
  );
};

export default ReservationModal;