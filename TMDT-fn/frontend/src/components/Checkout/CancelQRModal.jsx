import React, { useState } from 'react';
import { X } from 'lucide-react';

const CancelQRModal = ({ isOpen, onClose, onConfirm }) => {
    const [isLoading, setIsLoading] = useState(false);

    if (!isOpen) return null;

    const handleConfirm = async () => {
        if (isLoading) return;
        setIsLoading(true);
        try {
            await onConfirm();
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.6)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 9999, backdropFilter: 'blur(4px)' }}>
            <div style={{ background: '#fff', padding: '24px', borderRadius: '16px', width: '90%', maxWidth: '360px', textAlign: 'center', animation: 'scaleUp 0.2s ease-out' }}>
                <div style={{ width: 64, height: 64, borderRadius: '50%', background: '#fee2e2', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
                    <X size={32} color="#ef4444" />
                </div>
                <h3 style={{ margin: '0 0 12px', fontSize: '1.25rem', color: '#0f172a' }}>Xác nhận hủy đơn hàng?</h3>
                <p style={{ color: '#64748b', fontSize: '0.95rem', margin: '0 0 24px', lineHeight: 1.5 }}>
                    Bạn có chắc chắn muốn hủy giao dịch này không? Đơn hàng sẽ bị xóa khỏi hệ thống.
                </p>
                <div style={{ display: 'flex', gap: '12px' }}>
                    <button
                        onClick={onClose}
                        disabled={isLoading}
                        style={{ flex: 1, padding: '12px', borderRadius: '8px', background: '#f1f5f9', color: '#475569', border: 'none', fontWeight: '600', cursor: isLoading ? 'not-allowed' : 'pointer', opacity: isLoading ? 0.7 : 1 }}
                    >
                        Đóng
                    </button>
                    <button
                        onClick={handleConfirm}
                        disabled={isLoading}
                        style={{ flex: 1, padding: '12px', borderRadius: '8px', background: isLoading ? '#f87171' : '#ef4444', color: '#fff', border: 'none', fontWeight: '600', cursor: isLoading ? 'not-allowed' : 'pointer', opacity: isLoading ? 0.7 : 1 }}
                    >
                        {isLoading ? 'Đang hủy...' : 'Hủy đơn hàng'}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default CancelQRModal;
