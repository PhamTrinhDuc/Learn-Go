import React, { useEffect } from 'react';
import { useLocation, useNavigate, Navigate } from 'react-router-dom';
import { CheckCircle } from 'lucide-react';
import './Checkout.css';

const OrderSuccess = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const orderSuccess = location.state?.orderSuccess;

    if (!orderSuccess) {
        return <Navigate to="/" replace />;
    }

    useEffect(() => {
        window.scrollTo(0, 0);
    }, []);

    return (
        <div style={{ background: '#f1f5f9', minHeight: '70vh', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '40px 20px' }}>
            <div style={{ background: '#fff', borderRadius: '24px', padding: '48px 32px', maxWidth: '540px', width: '100%', textAlign: 'center', boxShadow: '0 20px 40px -10px rgba(0,0,0,0.08)', animation: 'fadeIn 0.5s ease-out' }}>

                {/* Visual Icon Header */}
                <div style={{ width: '100px', height: '100px', background: '#ecfdf5', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 28px', boxShadow: '0 0 0 16px #f0fdf4' }}>
                    <CheckCircle size={54} color="#10b981" strokeWidth={2.5} />
                </div>

                <h2 style={{ color: '#0f172a', margin: '0 0 16px', fontSize: '2.2rem', fontWeight: '800', letterSpacing: '-0.5px' }}>
                    Đặt hàng thành công!
                </h2>

                <p style={{ color: '#475569', margin: '0 0 32px', lineHeight: '1.6', fontSize: '1.05rem' }}>
                    Cảm ơn bạn đã tiếp tục tin tưởng và mua sắm tại hệ thống. Đơn hàng của bạn sẽ được xử lý và giao trong thời gian sớm nhất.
                </p>

                {/* Order ID Box */}
                <div style={{ background: '#f8fafc', borderRadius: '16px', padding: '20px', marginBottom: '32px', border: '1px dashed #cbd5e1' }}>
                    <p style={{ margin: 0, fontSize: '1rem', color: '#64748b', fontWeight: '500' }}>Mã tra cứu đơn hàng</p>
                    <p style={{ margin: '8px 0 0', fontSize: '1.75rem', color: '#0284c7', fontWeight: '800', letterSpacing: '1px', fontFamily: 'monospace' }}>
                        #{orderSuccess.orderId || orderSuccess.order_id}
                    </p>
                </div>

                {/* Unified Action Buttons */}
                <div style={{ display: 'flex', gap: '16px', justifyContent: 'center', flexWrap: 'wrap' }}>
                    <button
                        style={{ flex: '1 1', minWidth: '200px', padding: '16px 24px', borderRadius: '12px', fontWeight: '700', border: 'none', background: '#0284c7', color: 'white', cursor: 'pointer', fontSize: '1.05rem', boxShadow: '0 4px 12px rgba(2, 132, 199, 0.2)', transition: 'background 0.2s' }}
                        onClick={() => navigate('/')}
                        onMouseOver={(e) => e.target.style.background = '#0369a1'}
                        onMouseOut={(e) => e.target.style.background = '#0284c7'}
                    >
                        Tiếp tục mua sắm
                    </button>

                    <button
                        style={{ flex: '1 1', minWidth: '200px', padding: '16px 24px', borderRadius: '12px', fontWeight: '700', border: '2px solid #e2e8f0', background: 'white', color: '#0f172a', cursor: 'pointer', fontSize: '1.05rem', transition: 'all 0.2s' }}
                        onClick={() => navigate('/history')}
                        onMouseOver={(e) => { e.target.style.background = '#f8fafc'; e.target.style.borderColor = '#cbd5e1'; }}
                        onMouseOut={(e) => { e.target.style.background = 'white'; e.target.style.borderColor = '#e2e8f0'; }}
                    >
                        Xem đơn hàng
                    </button>
                </div>

            </div>
        </div>
    );
};

export default OrderSuccess;
