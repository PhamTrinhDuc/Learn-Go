import React, { useState, useEffect, useRef } from 'react';
import { useLocation, useNavigate, Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useCart } from '../../context/CartContext';
import QRCode from 'react-qr-code';
import { io } from 'socket.io-client';
import toast from 'react-hot-toast';
import { Copy, AlertCircle, X, Check } from 'lucide-react';
import CancelQRModal from '../../components/Checkout/CancelQRModal';
import './Checkout.css';

const API = import.meta.env.VITE_SERVER_API;

const PaymentQR = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const { user } = useAuth();
    const { cartItems, clearCart, removeMultipleFromCart } = useCart();

    const qrData = location.state?.qrData;
    const selectedItems = location.state?.selectedItems || [];

    const [countdown, setCountdown] = useState('15:00');
    const [isExpired, setIsExpired] = useState(false);
    const [showCancelModal, setShowCancelModal] = useState(false);
    const isCancelingRef = useRef(false);

    useEffect(() => {
        window.scrollTo(0, 0);
    }, []);

    // Xử lý sao chép văn bản
    const handleCopy = (text, label) => {
        navigator.clipboard.writeText(text).then(() => {
            toast.success(`Đã sao chép ${label}`, { id: 'copy-toast', position: 'top-center' });
        });
    };

    const cleanupCartAndFinish = (orderId, totalAmount) => {
        if (selectedItems.length > 0) {
            const indicesToRemove = cartItems.map((item, idx) => {
                const hit = selectedItems.some(sel => {
                    const selId = sel.product?.id || sel.product?.product_id;
                    const itemId = item.product?.id || item.product?.product_id;
                    return selId === itemId && sel.capacity === item.capacity && sel.color_name === item.color_name;
                });
                return hit ? idx : -1;
            }).filter(i => i !== -1);
            removeMultipleFromCart(indicesToRemove);
        } else {
            clearCart();
        }
        navigate('/checkout/success', { state: { orderSuccess: { orderId, total_amount: totalAmount } } });
    };

    const handleConfirmCancel = async (autoExpireError = null) => {
        if (isCancelingRef.current) return;
        isCancelingRef.current = true;

        try {
            if (user?.id) {
                await fetch(`${API}/api/order/reserve/${user.id}`, { method: 'DELETE' });
            }
        } catch (err) {
            console.error('Lỗi khi hủy giữ hàng chủ động:', err);
        } finally {
            if (autoExpireError) {
                navigate('/checkout/failed', { state: { error: autoExpireError } });
            } else {
                toast('Đã hủy giao dịch thanh toán.', { icon: 'i', id: 'cancel-payment' });
                navigate('/');
            }
        }
    };

    useEffect(() => {
        if (!qrData || !user?.id) return;

        let interval;
        const expireDate = new Date(qrData.expires_at).getTime();

        const updateTimer = () => {
            const now = new Date().getTime();
            const distance = expireDate - now;

            if (distance <= 0) {
                clearInterval(interval);
                setCountdown('00:00');
                setIsExpired(true);
            } else {
                const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
                const seconds = Math.floor((distance % (1000 * 60)) / 1000);
                setCountdown(`${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`);
            }
        };

        updateTimer();
        interval = setInterval(updateTimer, 1000);

        return () => clearInterval(interval);
    }, [qrData, user?.id]);

    useEffect(() => {
        if (isExpired) {
            handleConfirmCancel('Phiên QR Code 15 phút đã hết hạn. Hệ thống tự động hủy giao dịch này.');
        }
    }, [isExpired]);

    useEffect(() => {
        if (!user?.id || !qrData) return;

        const socket = io(API, { transports: ['websocket'] });
        const eventName = `payos_paid_${user.id}`;

        socket.on(eventName, (data) => {
            console.log('[Socket] PayOS paid event received:', data);
            toast.success('Thanh toán thành công! Đơn hàng đã được nhận.', { id: 'payment-success' });
            cleanupCartAndFinish(data.order_id, data.total_amount);
        });

        return () => {
            socket.off(eventName);
            socket.disconnect();
        };
    }, [user?.id, qrData]);

    if (!qrData) {
        return <Navigate to="/cart" replace />;
    }

    const formatPrice = (price) =>
        new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price);

    const bankAccountStr = qrData.account_number || "0968118125";
    const accountNameStr = qrData.account_name || "Trần Việt Hoàng";
    const amountStr = qrData.total_amount;
    const contentStr = qrData.payment_code;

    return (
        <div className="checkout-page" style={{ padding: '40px 20px', background: '#f8f9fa', minHeight: '100vh', display: 'flex', justifyContent: 'center', alignItems: 'flex-start' }}>
            <div className="container" style={{ maxWidth: '800px', width: '100%', background: '#fff', borderRadius: '16px', boxShadow: '0 10px 30px rgba(0,0,0,0.08)', overflow: 'hidden', padding: 0 }}>
                {/* Header */}
                <div style={{ background: '#f1f5f9', padding: '16px 20px', textAlign: 'center', borderBottom: '1px solid #e2e8f0' }}>
                    <h3 style={{ margin: 0, fontSize: '1.2rem', color: '#0f172a', fontWeight: '700' }}>
                        Quét mã hoặc chuyển khoản chính xác số tiền, nội dung bên dưới
                    </h3>
                </div>

                <div style={{ padding: '24px 28px 32px' }}>
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: '32px', alignItems: 'flex-start' }}>

                        {/* LEFT COLUMN: QR Code */}
                        <div style={{ flex: '1 1 300px', display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                            <div style={{ padding: '16px', border: '2px solid #005baa', borderRadius: '16px', display: 'inline-block', marginBottom: '24px', background: 'white' }}>
                                {isExpired ? (
                                    <div style={{ width: 220, height: 220, display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#fff1f0' }}>
                                        <span style={{ color: '#cf1322', fontWeight: 'bold' }}>Mã đã hết hạn</span>
                                    </div>
                                ) : (
                                    <QRCode value={qrData.qr_url} size={220} />
                                )}
                            </div>

                            {/* Bank Logos */}
                            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '16px', fontSize: '0.9rem', color: '#005baa', fontWeight: '700' }}>
                                <span>napas247</span>
                                <div style={{ width: 1, height: 16, background: '#cbd5e1' }} />
                                <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                                    <span style={{ color: '#E3001B', fontSize: '1.2rem' }}>★</span> MB
                                </span>
                            </div>
                            {/* Countdown */}
                            <div style={{ marginTop: '32px' }}>
                                <span style={{ fontSize: '0.9rem', color: '#64748b' }}>Phiên giao dịch kết thúc sau </span>
                                <strong style={{ color: '#ee4d2d', fontSize: '1.05rem' }}>{countdown}</strong>
                            </div>
                        </div>

                        {/* RIGHT COLUMN: Information */}
                        <div style={{ flex: '1 1 340px', background: '#fff', fontSize: '0.95rem', color: '#334155' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
                                <div style={{ width: 40, height: 40, background: '#002888', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#E3001B', fontSize: '1.5rem' }}>★</div>
                                <div>
                                    <div style={{ color: '#64748b', fontSize: '0.85rem' }}>Ngân hàng</div>
                                    <div style={{ fontWeight: '700', color: '#0f172a' }}>Ngân hàng TMCP Quân đội</div>
                                </div>
                            </div>

                            <div style={{ marginBottom: '16px' }}>
                                <div style={{ color: '#64748b', fontSize: '0.85rem', marginBottom: '2px' }}>Chủ tài khoản:</div>
                                <div style={{ fontWeight: '700', color: '#0f172a', textTransform: 'uppercase' }}>{accountNameStr}</div>
                            </div>

                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                                <div>
                                    <div style={{ color: '#64748b', fontSize: '0.85rem', marginBottom: '2px' }}>Số tài khoản:</div>
                                    <div style={{ fontWeight: '700', color: '#0f172a' }}>{bankAccountStr}</div>
                                </div>
                                <button onClick={() => handleCopy(bankAccountStr, 'Số tài khoản')} style={{ border: 'none', background: '#e0f2fe', color: '#0369a1', padding: '6px 12px', borderRadius: '4px', fontSize: '0.8rem', fontWeight: '600', cursor: 'pointer' }}>Sao chép</button>
                            </div>

                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                                <div>
                                    <div style={{ color: '#64748b', fontSize: '0.85rem', marginBottom: '2px' }}>Số tiền:</div>
                                    <div style={{ fontWeight: '700', color: '#0f172a' }}>{formatPrice(amountStr)}</div>
                                </div>
                                <button onClick={() => handleCopy(amountStr.toString(), 'Số tiền')} style={{ border: 'none', background: '#e0f2fe', color: '#0369a1', padding: '6px 12px', borderRadius: '4px', fontSize: '0.8rem', fontWeight: '600', cursor: 'pointer' }}>Sao chép</button>
                            </div>

                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '24px' }}>
                                <div style={{ flex: 1, paddingRight: '12px' }}>
                                    <div style={{ color: '#64748b', fontSize: '0.85rem', marginBottom: '2px' }}>Nội dung:</div>
                                    <div style={{ fontWeight: '700', color: '#0f172a', wordBreak: 'break-word' }}>{contentStr}</div>
                                </div>
                                <button onClick={() => handleCopy(contentStr, 'Nội dung')} style={{ border: 'none', background: '#e0f2fe', color: '#0369a1', padding: '6px 12px', borderRadius: '4px', fontSize: '0.8rem', fontWeight: '600', cursor: 'pointer', marginTop: '16px' }}>Sao chép</button>
                            </div>

                            <div style={{ background: '#f8fafc', padding: '12px', borderRadius: '8px', fontSize: '0.85rem', color: '#475569', border: '1px solid #e2e8f0', marginBottom: '24px', display: 'flex', gap: '8px', alignItems: 'flex-start' }}>
                                <AlertCircle size={16} color="#0284c7" style={{ flexShrink: 0, marginTop: '2px' }} />
                                <div>Lưu ý: Nhập chính xác số tiền <strong style={{ color: '#0f172a' }}>{formatPrice(amountStr).replace(' ₫', '')}</strong>, nội dung <strong style={{ color: '#0f172a' }}>{contentStr}</strong> khi chuyển khoản.</div>
                            </div>
                        </div>
                    </div>

                    <hr style={{ border: 'none', borderTop: '1px solid #e2e8f0', margin: '20px 0 16px' }} />

                    {/* BOTTOM: Countdown & Cancel */}
                    <div style={{ textAlign: 'center', maxWidth: '300px', margin: '0 auto' }}>
                        <button
                            type="button"
                            onClick={() => setShowCancelModal(true)}
                            style={{ width: '100%', padding: '14px', borderRadius: '8px', border: '1px solid #cbd5e1', background: '#fff', color: '#475569', fontWeight: '600', fontSize: '1rem', cursor: 'pointer', transition: 'all 0.2s', boxShadow: '0 2px 4px rgba(0,0,0,0.02)' }}
                            onMouseOver={(e) => { e.target.style.background = '#f1f5f9'; e.target.style.borderColor = '#94a3b8'; }}
                            onMouseOut={(e) => { e.target.style.background = '#fff'; e.target.style.borderColor = '#cbd5e1'; }}
                        >
                            Huỷ
                        </button>
                    </div>

                </div>
            </div>

            <CancelQRModal
                isOpen={showCancelModal}
                onClose={() => setShowCancelModal(false)}
                onConfirm={() => handleConfirmCancel()}
            />

        </div>
    );
};

export default PaymentQR;
