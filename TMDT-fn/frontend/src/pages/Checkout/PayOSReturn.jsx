import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { CheckCircle } from 'lucide-react';

/**
 * Trang được PayOS redirect về sau khi thanh toán thành công
 * Socket từ webhook sẽ đã tạo đơn hàng → trang này chỉ hiển thị thông báo.
 */
const PayOSReturn = () => {
    const navigate = useNavigate();
    const [params] = useSearchParams();
    const orderCode = params.get('orderCode');
    const status = params.get('status');

    useEffect(() => {
        // Tự redirect về lịch sử đơn hàng sau 3 giây
        const t = setTimeout(() => navigate('/history'), 3000);
        return () => clearTimeout(t);
    }, [navigate]);

    return (
        <div style={{ minHeight: '70vh', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: 20, textAlign: 'center', padding: 24 }}>
            <CheckCircle size={72} color="#10b981" />
            <h2 style={{ margin: 0, fontSize: '1.8rem', color: '#0f172a' }}>Thanh toán thành công!</h2>
            <p style={{ color: '#64748b', fontSize: '1rem', margin: 0 }}>
                {status === 'PAID' ? `Giao dịch #${orderCode} đã hoàn tất.` : 'Đơn hàng của bạn đang được xử lý.'}
                <br />Bạn sẽ được chuyển về lịch sử đơn hàng trong giây lát...
            </p>
            <button
                onClick={() => navigate('/history')}
                style={{ padding: '12px 32px', borderRadius: 10, background: '#10b981', color: '#fff', fontWeight: 700, fontSize: '1rem', border: 'none', cursor: 'pointer' }}
            >
                Xem đơn hàng ngay
            </button>
        </div>
    );
};

export default PayOSReturn;
