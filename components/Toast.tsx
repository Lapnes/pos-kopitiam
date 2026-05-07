'use client';

import { CheckCircle, XCircle, X } from 'lucide-react';

interface ToastProps {
  message: string;
  isError?: boolean;
  onClose: () => void;
}

export default function Toast({ message, isError, onClose }: ToastProps) {
  return (
    <div className="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 anim-toast-in">
      <div className="flex items-center gap-3 pl-4 pr-3 py-3 rounded-2xl shadow-xl"
        style={{
          background:   isError ? 'var(--danger)' : 'var(--sage-dark)',
          color:        '#fff',
          minWidth:     220,
          maxWidth:     360,
          boxShadow:    isError
            ? '0 8px 30px rgba(192,56,56,0.3)'
            : '0 8px 30px rgba(61,92,61,0.35)',
        }}>
        {isError
          ? <XCircle size={17} style={{ flexShrink: 0, opacity: 0.9 }} />
          : <CheckCircle size={17} style={{ flexShrink: 0, opacity: 0.9 }} />}
        <span style={{ fontSize: 13, fontWeight: 600, flex: 1 }}>{message}</span>
        <button onClick={onClose}
          className="w-6 h-6 flex items-center justify-center rounded-lg transition-all ml-1"
          style={{ background: 'rgba(255,255,255,0.15)', flexShrink: 0 }}
          onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'rgba(255,255,255,0.25)'; }}
          onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'rgba(255,255,255,0.15)'; }}>
          <X size={13} />
        </button>
      </div>
    </div>
  );
}
