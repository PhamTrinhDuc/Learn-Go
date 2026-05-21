import React from 'react';
import { useLocation, useNavigate, Navigate } from 'react-router-dom';
import { XCircle } from 'lucide-react';
import './Checkout.css';

const OrderFailed = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const errorMsg = location.state?.error || 'Đã có lỗi xảy ra trong quá trình xử lý đơn hàng.';

    return (
        <div className="checkout-page">
            <div className="container">
                <div className="order-success-screen" style={{ textAlign: 'center', padding: '60px 20px', background: '#fff', borderRadius: '16px', boxShadow: '0 10px 30px rgba(0,0,0,0.05)' }}>
                    <XCircle size={72} color="#ee4d2d" style={{ margin: '0 auto 20px', display: 'block' }} />
                    <h2 style={{ color: '#ee4d2d', marginBottom: '16px', fontSize: '1.8rem' }}>Đặt hàng không thành công</h2>
                    <p style={{ color: '#64748b', marginBottom: '32px', fontSize: '1.1rem' }}>{errorMsg}</p>
                    <div style={{ display: 'flex', gap: '20px', justifyContent: 'center', flexWrap: 'wrap', marginTop: '40px' }}>
                        <button 
                            style={{ padding: '14px 32px', borderRadius: '8px', fontWeight: '600', border: 'none', background: '#0284c7', color: 'white', cursor: 'pointer', minWidth: '220px', fontSize: '1rem', boxShadow: '0 4px 12px rgba(2, 132, 199, 0.2)' }}
                            onClick={() => navigate('/cart')}
                        >
                            Kiểm tra giỏ hàng
                        </button>
                        <button 
                            style={{ padding: '14px 32px', borderRadius: '8px', fontWeight: '600', border: '1px solid #cbd5e1', background: 'white', color: '#475569', cursor: 'pointer', minWidth: '220px', fontSize: '1rem' }}
                            onClick={() => navigate('/')}
                        >
                            Về trang chủ
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default OrderFailed;
