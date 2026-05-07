'use client';

import { useEffect } from 'react';
import { X } from 'lucide-react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

export default function Modal({ isOpen, onClose, children }: ModalProps) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
    if (isOpen) document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center anim-fade-in"
      style={{ background: 'rgba(42,26,8,0.45)', backdropFilter: 'blur(6px)' }}
      onClick={onClose}>
      <div className="anim-pop-in relative rounded-2xl shadow-2xl w-full max-w-md mx-4 p-6"
        style={{ background: 'var(--surface)', border: '1px solid var(--border)' }}
        onClick={(e) => e.stopPropagation()}>
        <button onClick={onClose}
          className="absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-lg transition-all"
          style={{ background: 'var(--surface3)', color: 'var(--text-3)' }}
          onMouseEnter={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--border-med)'; (e.currentTarget as HTMLElement).style.color = 'var(--text)'; }}
          onMouseLeave={(e) => { (e.currentTarget as HTMLElement).style.background = 'var(--surface3)'; (e.currentTarget as HTMLElement).style.color = 'var(--text-3)'; }}>
          <X size={16} />
        </button>
        {children}
      </div>
    </div>
  );
}
