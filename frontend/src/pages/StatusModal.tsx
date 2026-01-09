import React, { ReactElement } from 'react';
import { Shield, CheckCircle, AlertTriangle, Activity, X } from 'lucide-react';

interface StatusModalProps {
  isOpen: boolean;
  onClose: () => void;
  type: 'SUCCESS' | 'ERROR' | 'WAITING' | 'AUTH';
  title: string;
  message: string;
  onConfirm?: () => void;
}

const StatusModal = ({ isOpen, onClose, type, title, message, onConfirm }: StatusModalProps) => {
  if (!isOpen) return null;

  const themes = {
    SUCCESS: { color: 'text-emerald-400', border: 'border-emerald-500/50', icon: <CheckCircle /> },
    ERROR: { color: 'text-red-500', border: 'border-red-500/50', icon: <AlertTriangle /> },
    WAITING: { color: 'text-yellow-400', border: 'border-yellow-500/50', icon: <Activity /> },
    AUTH: { color: 'text-emerald-400', border: 'border-emerald-500/50', icon: <Shield /> }
  };

  const current = themes[type];

  return (
    <div className="fixed inset-0 z-[999] flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
      <div className={`w-full max-w-sm bg-zinc-950 border-2 ${current.border} rounded-3xl p-8 relative shadow-[0_0_50px_rgba(16,185,129,0.1)]`}>
        <button onClick={onClose} className="absolute top-4 right-4 text-zinc-500 hover:text-white">
          <X size={24} />
        </button>
        <div className="text-center">
          <div className={`${current.color} w-16 h-16 mx-auto mb-6 flex items-center justify-center`}>
            {React.cloneElement(current.icon as ReactElement<any>, { size: 48 })}
          </div>
          <h3 className={`text-2xl font-black mb-2 tracking-tighter uppercase ${current.color}`} style={{ fontFamily: 'Orbitron' }}>
            {title}
          </h3>
          <p className="text-zinc-400 font-medium leading-relaxed mb-8 whitespace-pre-wrap">{message}</p>
          <button 
            onClick={onConfirm || onClose}
            className="w-full py-4 bg-zinc-900 border border-emerald-500/30 rounded-xl text-white font-black hover:bg-emerald-500/10 transition-all"
          >
            CONFIRM
          </button>
        </div>
      </div>
    </div>
  );
};

export default StatusModal;