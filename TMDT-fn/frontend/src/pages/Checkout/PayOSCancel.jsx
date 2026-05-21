import React from 'react';
import { useNavigate } from 'react-router-dom';
import { XCircle } from 'lucide-react';

/**
 * Trang được PayOS redirect khi user bấm "Hủy" trên trang thanh toán PayOS
 */
const PayOSCancel = () => {
    const navigate = useNavigate();

    return (
        <div style={{ minHeight: '70vh', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 20, textAlign: 'center', padding: 24 }}>
            <XCircle size={72} color="#ef4444" />
            <h2 style={{ margin: 0, fontSize: '1.8rem', color: '#0f172a' }}>Thanh toán đã bị hủy</h2>
            <p style={{ color: '#64748b', fontSize: '1rem', margin: 0 }}>
                Bạn đã hủy giao dịch thanh toán.<br />
                Hàng đã được giữ cho bạn — hãy quay lại trang thanh toán nếu muốn thử lại.
            </p>
            <div style={{ display: 'flex', gap: 12 }}>
                <button
                    onClick={() => navigate(-1)}
                    style={{ padding: '12px 28px', borderRadius: 10, background: '#6366f1', color: '#fff', fontWeight: 700, fontSize: '1rem', border: 'none', cursor: 'pointer' }}
                >
                    Quay lại thanh toán
                </button>
                <button
                    onClick={() => navigate('/')}
                    style={{ padding: '12px 28px', borderRadius: 10, background: '#e2e8f0', color: '#475569', fontWeight: 700, fontSize: '1rem', border: 'none', cursor: 'pointer' }}
                >
                    Về trang chủ
                </button>
            </div>
        </div>
    );
};

export default PayOSCancel;
