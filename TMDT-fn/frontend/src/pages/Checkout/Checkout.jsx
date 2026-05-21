import React, { useState, useEffect, useRef } from 'react';
import { useCart } from '../../context/CartContext';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import toast from 'react-hot-toast';
import { Loader, MapPin, CheckCircle, ScanLine } from 'lucide-react';
import CheckoutAddress from '../../components/Checkout/CheckoutAddress';
import './Checkout.css';

const API = import.meta.env.VITE_SERVER_API;

const Checkout = () => {
    const { cartItems, clearCart, removeMultipleFromCart } = useCart();
    const { user } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const selectedItems = location.state?.selectedItems || [];
    const displayItems = selectedItems.length > 0 ? selectedItems : cartItems;

    const [couponCode, setCouponCode] = useState(location.state?.coupon || '');
    const [discountValue, setDiscountValue] = useState(location.state?.discount || 0);
    const [inputCoupon, setInputCoupon] = useState(location.state?.coupon || '');

    const [shippingVoucherCode, setShippingVoucherCode] = useState('');
    const [inputShipping, setInputShipping] = useState('');
    const [shippingDiscountValue, setShippingDiscountValue] = useState(0);

    const [savedAddresses, setSavedAddresses] = useState([]);
    const [selectedAddressId, setSelectedAddressId] = useState(null);
    const [loadingAddress, setLoadingAddress] = useState(true);

    const [shippingOptions] = useState([
        { id: 'standard', name: 'Giao hàng tiêu chuẩn', fee: 30000, time: '3-4 ngày' },
        { id: 'express', name: 'Giao hàng hỏa tốc', fee: 55000, time: '1-2 ngày' },
    ]);
    const [selectedShipping, setSelectedShipping] = useState('standard');
    const [paymentMethod, setPaymentMethod] = useState('QR');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [userNote, setUserNote] = useState('');

    useEffect(() => {
        if (!user?.id) { setLoadingAddress(false); return; }
        fetch(`${API}/api/user/address/${user.id}`)
            .then(r => r.json())
            .then(data => {
                if (data.success && data.data.length > 0) {
                    const sorted = [...data.data].sort((a, b) => b.is_default - a.is_default);
                    setSavedAddresses(sorted);
                    setSelectedAddressId(sorted[0].id);
                } else {
                    setSavedAddresses([]);
                }
            })
            .catch(() => setSavedAddresses([]))
            .finally(() => setLoadingAddress(false));
    }, [user]);

    const formatPrice = (price) =>
        new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price);

        const PHOTO_API = (import.meta.env.VITE_PHOTO_SERVER_API || 'http://localhost:8081/images').replace(/\/+$/, '');

        const getImageUrl = (img) => {
            if (!img || img.endsWith('/')) return '/placeholder-image.png';
            if (img.startsWith('http')) return img;

            const cleanImg = img.startsWith('/') ? img : `/${img}`;
            console.log(`${PHOTO_API}${cleanImg}`);
            return `${PHOTO_API}${cleanImg}`;
    };

    const getShippingFee = () => {
        const opt = shippingOptions.find(o => o.id === selectedShipping);
        return opt ? opt.fee : 0;
    };

    const checkoutTotal = displayItems.reduce((sum, item) => {
        const price = item.price || item.product.calculated_price || item.product.price || 0;
        return sum + price * item.quantity;
    }, 0);

    const totalShippingToPay = Math.max(0, getShippingFee() - shippingDiscountValue);
    const totalToPay = Math.max(0, checkoutTotal - discountValue) + totalShippingToPay;

    // Re-apply shipping voucher if active and shipping method changes
    useEffect(() => {
        if (shippingVoucherCode) {
            fetch(`${API}/api/voucher/apply`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    code: shippingVoucherCode,
                    orderTotal: getShippingFee()
                })
            })
            .then(r => r.json())
            .then(data => {
                if (data.success && data.discountTarget === 'shipping') {
                    setShippingDiscountValue(data.discountAmount);
                } else {
                    setShippingDiscountValue(0);
                }
            })
            .catch(err => {
                console.error(err);
                setShippingDiscountValue(0);
            });
        } else {
            setShippingDiscountValue(0);
        }
    }, [selectedShipping, shippingVoucherCode]);

    const handleApplyProductVoucher = async () => {
        if (!inputCoupon.trim()) {
            setCouponCode('');
            setDiscountValue(0);
            return;
        }
        try {
            const res = await fetch(`${API}/api/voucher/apply`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    code: inputCoupon.trim(),
                    orderTotal: checkoutTotal
                })
            });
            const data = await res.json();
            if (!res.ok || !data.success) {
                toast.error(data.message || 'Mã giảm giá không hợp lệ!', { id: 'promo-toast' });
                return;
            }
            if (data.discountTarget !== 'product') {
                toast.error('Mã này là mã miễn phí vận chuyển, vui lòng nhập ở ô bên dưới!', { id: 'promo-toast' });
                return;
            }
            setCouponCode(data.voucherCode);
            setDiscountValue(data.discountAmount);
            toast.success(`Áp dụng mã ${data.voucherCode} thành công!`, { id: 'promo-toast' });
        } catch (err) {
            console.error(err);
            toast.error('Có lỗi xảy ra khi áp dụng mã!', { id: 'promo-toast' });
        }
    };

    const handleApplyShippingVoucher = async () => {
        if (!inputShipping.trim()) {
            setShippingVoucherCode('');
            setShippingDiscountValue(0);
            return;
        }
        try {
            const res = await fetch(`${API}/api/voucher/apply`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    code: inputShipping.trim(),
                    orderTotal: getShippingFee()
                })
            });
            const data = await res.json();
            if (!res.ok || !data.success) {
                toast.error(data.message || 'Mã freeship không hợp lệ!', { id: 'ship-toast' });
                return;
            }
            if (data.discountTarget !== 'shipping') {
                toast.error('Mã này là mã giảm giá sản phẩm, vui lòng nhập ở ô bên trên!', { id: 'ship-toast' });
                return;
            }
            setShippingVoucherCode(data.voucherCode);
            setShippingDiscountValue(data.discountAmount);
            toast.success(`Áp dụng mã ${data.voucherCode} thành công!`, { id: 'ship-toast' });
        } catch (err) {
            console.error(err);
            toast.error('Có lỗi xảy ra khi áp dụng mã!', { id: 'ship-toast' });
        }
    };

    // Gán phí vận chuyển và giảm giá theo tỉ lệ (như Shopee) để khi lưu order_details mọi item đều có phí tương ứng
    const buildPayloadItems = () => {
        const totalFee = getShippingFee() || 0;
        const totalProductDiscount = Math.min(discountValue || 0, checkoutTotal);
        const totalShippingDiscount = Math.min(shippingDiscountValue || 0, totalFee);

        let sumPrice = displayItems.reduce((sum, item) => {
            const p = item.price || item.product?.calculated_price || item.product?.price || 0;
            return sum + (p * item.quantity);
        }, 0);

        if (sumPrice === 0) sumPrice = 1;

        let accumulatedFee = 0;
        let accumulatedProductDiscount = 0;
        let accumulatedShippingDiscount = 0;

        const payload = displayItems.map((item, index) => {
            let vId = item.variant_id || item.selectedColor?.variant_id || item.product?.variant_id;
            const matchColor = item.color_name || item.color;
            if (!vId && item.product?.variants) {
                const match = item.product.variants.find(v => v.color_name === matchColor);
                vId = match ? (match.variant_id || match.id) : (item.product.variants[0]?.variant_id || item.product.variants[0]?.id);
            }
            if (!vId) throw new Error(`Sản phẩm ${item.product?.name} thiếu thông tin phân loại.`);

            const p = item.price || item.product?.calculated_price || item.product?.price || 0;
            const lineValue = p * item.quantity;
            const isLast = index === displayItems.length - 1;

            const myShipping = isLast ? totalFee - accumulatedFee : Math.round((lineValue / sumPrice) * totalFee);
            const myProductDiscount = isLast ? totalProductDiscount - accumulatedProductDiscount : Math.round((lineValue / sumPrice) * totalProductDiscount);
            const myShippingDiscount = isLast ? totalShippingDiscount - accumulatedShippingDiscount : Math.round((lineValue / sumPrice) * totalShippingDiscount);

            accumulatedFee += myShipping;
            accumulatedProductDiscount += myProductDiscount;
            accumulatedShippingDiscount += myShippingDiscount;

            return {
                variant_id: Number(vId),
                quantity: Number(item.quantity),
                shipping_price: myShipping,
                shipping_support_price: myShippingDiscount,
                product_support_price: myProductDiscount
            };
        });

        console.log("PAYLOAD SẼ GỬI LÊN BACKEND (Kiểm tra shipping_price):", payload);
        return payload;
    };

    // Gộp ghi chú: mã giảm giá + ghi chú người dùng + hình thức vận chuyển
    const buildNote = () => {
        const shippingOpt = shippingOptions.find(o => o.id === selectedShipping);
        const parts = [
            couponCode ? `Mã SP: ${couponCode}` : '',
            shippingVoucherCode ? `Mã Freeship: ${shippingVoucherCode}` : '',
            userNote,
            shippingOpt ? `Giao hàng: ${shippingOpt.name}` : ''
        ];
        return parts.filter(Boolean).join(' | ');
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

    // ─── LUỒNG 1: COD TRỰC TIẾP (KHÔNG GIỮ HÀNG) ──────────────────────────────────
    const processDirectCOD = async (addressId) => {
        setIsSubmitting(true);
        try {
            const appliedVouchers = [couponCode, shippingVoucherCode].filter(Boolean).join(', ');
            // checkoutDirect API: { user_id, items[{variant_id, quantity}], address_id, payment_method, note, voucherCode }
            const res = await fetch(`${API}/api/order/checkout-direct`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    user_id: Number(user.id || user.user_id),
                    address_id: Number(addressId),
                    payment_method: paymentMethod.toUpperCase(),
                    note: buildNote(),
                    items: buildPayloadItems(),
                    voucherCode: appliedVouchers || null,
                }),
            });
            const orderData = await res.json();
            if (!res.ok) throw new Error(orderData.error || 'Đặt hàng thất bại');
            cleanupCartAndFinish(orderData.order_id, orderData.total_amount);
        } catch (err) {
            toast.error(err.message || 'Lỗi đặt hàng COD', { id: 'checkout-error', duration: 5000 });
            navigate('/checkout/failed', { state: { error: err.message || 'Lỗi đặt hàng COD' } });
        } finally {
            setIsSubmitting(false);
        }
    };

    // ─── LUỒNG 2: THANH TOÁN QR (PHASE 1: GIỮ HÀNG VÀ CHUYỂN TRANG) ──────────
    const openQRModalAndReserve = async (addressId) => {
        setIsSubmitting(true);
        try {
            const payloadItems = buildPayloadItems();
            const reserveRes = await fetch(`${API}/api/order/reserve`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    user_id: Number(user.id || user.user_id),
                    address_id: Number(addressId),
                    payment_method: paymentMethod.toUpperCase(),
                    note: buildNote(),
                    items: payloadItems,
                    voucherCode: couponCode || shippingVoucherCode || null,
                }),
            });
            const reserveData = await reserveRes.json();

            if (!reserveRes.ok) {
                throw new Error(reserveData.error || 'Kho hàng thay đổi, không thể tạo phiên thanh toán.');
            }

            // Xử lý đơn hàng 0đ: Gọi API thanh toán trực tiếp không qua PayOS
            if (reserveData.isFree) {
                const freeRes = await fetch(`${API}/api/order/checkout-free-order`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        user_id: Number(user.id || user.user_id),
                        address_id: Number(addressId),
                        payment_method: paymentMethod.toUpperCase(),
                        note: buildNote(),
                        items: payloadItems,
                    }),
                });
                const freeData = await freeRes.json();
                if (!freeRes.ok) throw new Error(freeData.error || 'Lỗi xác nhận đơn hàng 0đ');

                toast.success('Xác nhận đơn hàng thành công!', { id: 'checkout-success' });
                cleanupCartAndFinish(freeData.order_id, freeData.total_amount);
                return;
            }

            navigate('/checkout/payment', {
                state: { qrData: reserveData, selectedItems: displayItems }
            });
        } catch (err) {
            toast.error(err.message, { id: 'checkout-error', duration: 5000 });
            setIsSubmitting(false);
        }
    };

    const handleSubmit = async (e) => {
        if (e && e.preventDefault) e.preventDefault();
        if (isSubmitting) return;
        if (displayItems.length === 0) {
            toast.error('Không có sản phẩm nào để đặt hàng!', { id: 'no-items-error' });
            return;
        }

        if (!selectedAddressId) {
            toast.error('Vui lòng chọn địa chỉ giao hàng!', { id: 'no-address-error' });
            return;
        }

        if (paymentMethod.toUpperCase() === 'COD') {
            await processDirectCOD(selectedAddressId);
        } else {
            await openQRModalAndReserve(selectedAddressId);
        }
    };

    if (cartItems.length === 0 && displayItems.length === 0) {
        return (
            <div className="checkout-page">
                <div className="container" style={{ padding: '60px', textAlign: 'center' }}>
                    <h2>Giỏ hàng trống</h2>
                    <button className="btn-primary" style={{ marginTop: 16 }} onClick={() => navigate('/')}>Mua sắm ngay</button>
                </div>
            </div>
        );
    }


    const selectedAddress = savedAddresses.find(a => a.id === selectedAddressId);

    return (
        <div className="checkout-page">
            <div className="container">
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h1 className="page-title">Thanh toán</h1>
                </div>
                <form onSubmit={handleSubmit} className="checkout-layout">

                    {/* ===== LEFT ===== */}
                    <div className="checkout-form">
                        <CheckoutAddress
                            user={user}
                            savedAddresses={savedAddresses}
                            setSavedAddresses={setSavedAddresses}
                            selectedAddressId={selectedAddressId}
                            setSelectedAddressId={setSelectedAddressId}
                            loadingAddress={loadingAddress}
                            API={API}
                        />

                        {/* 2. PRODUCT LIST */}
                        <section className="checkout-section">
                            <h3>2. Danh sách sản phẩm ({displayItems.length})</h3>
                            <div className="checkout-items-list">
                                {displayItems.map((item, idx) => {
                                    console.log(item)
                                    const price = item.price || item.product.calculated_price || item.product.price || 0;
                                    return (
                                        <div key={idx} className="checkout-item-row">
                                            <div className="item-img">
                                                <img
                                                    src={getImageUrl(item.image || item.product.img_thumb || item.product.image)}
                                                    alt={item.product.name}
                                                />
                                            </div>
                                            <div className="item-details">
                                                <div className="item-name">{item.product.name}</div>
                                                <div className="item-meta">
                                                    {item.color_name && <span>Màu: {item.color_name}</span>}
                                                    {item.capacity && <span> | {item.capacity}</span>}
                                                </div>
                                                <div className="item-price">
                                                    {item.quantity} x {formatPrice(price)}
                                                </div>
                                            </div>
                                            <div className="item-total">{formatPrice(price * item.quantity)}</div>
                                        </div>
                                    );
                                })}
                            </div>
                        </section>

                        {/* 3. GHI CHÚ */}
                        <section className="checkout-section">
                            <h3>3. Ghi chú đơn hàng</h3>
                            <textarea
                                className="checkout-note-input"
                                placeholder="Ghi chú thêm cho người bán (ví dụ: giao giờ hành chính, gọi trước khi giao...)"
                                value={userNote}
                                onChange={(e) => setUserNote(e.target.value)}
                                rows={3}
                            ></textarea>
                        </section>

                        {/* 4. SHIPPING */}
                        <section className="checkout-section">
                            <h3>4. Hình thức giao hàng</h3>
                            <div className="shipping-options">
                                {shippingOptions.map(opt => (
                                    <label
                                        key={opt.id}
                                        className={`shipping-option ${selectedShipping === opt.id ? 'selected' : ''}`}
                                    >
                                        <input
                                            type="radio"
                                            name="shipping"
                                            value={opt.id}
                                            checked={selectedShipping === opt.id}
                                            onChange={() => setSelectedShipping(opt.id)}
                                        />
                                        <div className="shipping-info">
                                            <strong>{opt.name}</strong>
                                            <span>{opt.time}</span>
                                        </div>
                                        <div className="shipping-fee">{formatPrice(opt.fee)}</div>
                                    </label>
                                ))}
                            </div>
                        </section>

                        {/* 5. PAYMENT */}
                        <section className="checkout-section">
                            <h3>5. Phương thức thanh toán</h3>
                            <div className="payment-methods">
                                {[
                                    { id: 'QR', label: 'Thanh toán quét mã (Ngân hàng / Ví điện tử)' },
                                    { id: 'COD', label: 'Thanh toán khi nhận hàng (COD)' },
                                ].map(m => (
                                    <label
                                        key={m.id}
                                        className={`payment-method ${paymentMethod === m.id ? 'selected' : ''}`}
                                    >
                                        <input
                                            type="radio"
                                            name="payment"
                                            value={m.id}
                                            checked={paymentMethod === m.id}
                                            onChange={() => setPaymentMethod(m.id)}
                                        />
                                        <span>{m.label}</span>
                                    </label>
                                ))}
                            </div>
                        </section>
                    </div>

                    {/* ===== RIGHT SIDEBAR ===== */}
                    <div className="checkout-sidebar">
                        <div className="checkout-summary">
                            <h3>Tóm tắt đơn hàng</h3>

                            {selectedAddress && (
                                <div className="summary-address-preview">
                                    <div className="summary-section-label">Giao đến:</div>
                                    <div className="summary-address-recipient">
                                        <strong>{selectedAddress.full_name}</strong> · {selectedAddress.num_phone}
                                    </div>
                                    <div className="summary-address-text">
                                        {selectedAddress.detail_address}, {selectedAddress.district}, {selectedAddress.province}
                                    </div>
                                </div>
                            )}

                            <div className="order-items-scroll">
                                {displayItems.map((item, idx) => {
                                    const price = item.price || item.product.calculated_price || item.product.price || 0;
                                    return (
                                        <div key={idx} className="order-item-mini">
                                            <div className="item-name">
                                                <span>{item.quantity}x {item.product.name}</span>
                                                {item.color_name && <small>Màu: {item.color_name}</small>}
                                                {item.capacity && <small>{item.capacity}</small>}
                                            </div>
                                            <span>{formatPrice(price * item.quantity)}</span>
                                        </div>
                                    );
                                })}
                            </div>
                            <hr />
                            <div className="summary-row">
                                <span>Tổng tiền hàng</span>
                                <span>{formatPrice(checkoutTotal)}</span>
                            </div>
                            <div className="summary-row">
                                <span>Phí vận chuyển</span>
                                <span>{formatPrice(getShippingFee())}</span>
                            </div>

                            {/* VOUCHER INPUTS */}
                            <div className="voucher-section" style={{ margin: '16px 0', borderTop: '1px dashed #e2e8f0', borderBottom: '1px dashed #e2e8f0', padding: '16px 0' }}>
                                <div style={{ marginBottom: '12px' }}>
                                    <label style={{ display: 'block', fontSize: '0.85rem', color: '#64748b', fontWeight: 600, marginBottom: '6px' }}>Mã giảm giá sản phẩm</label>
                                    <div style={{ display: 'flex', gap: '8px' }}>
                                        <input
                                            type="text"
                                            placeholder="VD: GIAM20K, PROMO10"
                                            value={inputCoupon}
                                            onChange={(e) => setInputCoupon(e.target.value.toUpperCase())}
                                            style={{ flex: 1, padding: '8px 12px', border: '1px solid #d9d9d9', borderRadius: '6px', fontSize: '0.9rem', outline: 'none' }}
                                        />
                                        <button
                                            type="button"
                                            onClick={handleApplyProductVoucher}
                                            style={{ padding: '0 12px', background: '#0f172a', color: 'white', borderRadius: '6px', border: 'none', cursor: 'pointer', fontSize: '0.85rem', fontWeight: '600' }}
                                        >Áp dụng</button>
                                    </div>
                                </div>
                                <div>
                                    <label style={{ display: 'block', fontSize: '0.85rem', color: '#64748b', fontWeight: 600, marginBottom: '6px' }}>Mã miễn phí vận chuyển</label>
                                    <div style={{ display: 'flex', gap: '8px' }}>
                                        <input
                                            type="text"
                                            placeholder="VD: FREESHIP, SHIP50"
                                            value={inputShipping}
                                            onChange={(e) => setInputShipping(e.target.value.toUpperCase())}
                                            style={{ flex: 1, padding: '8px 12px', border: '1px solid #d9d9d9', borderRadius: '6px', fontSize: '0.9rem', outline: 'none' }}
                                        />
                                        <button
                                            type="button"
                                            onClick={handleApplyShippingVoucher}
                                            style={{ padding: '0 12px', background: '#0f172a', color: 'white', borderRadius: '6px', border: 'none', cursor: 'pointer', fontSize: '0.85rem', fontWeight: '600' }}
                                        >Áp dụng</button>
                                    </div>
                                </div>
                            </div>

                            {discountValue > 0 && (
                                <div className="summary-row">
                                    <span>Giảm giá SP {couponCode && `(${couponCode})`}</span>
                                    <span style={{ color: '#059669', fontWeight: '600' }}>-{formatPrice(discountValue)}</span>
                                </div>
                            )}
                            {shippingDiscountValue > 0 && (
                                <div className="summary-row">
                                    <span>Miễn phí vận chuyển {shippingVoucherCode && `(${shippingVoucherCode})`}</span>
                                    <span style={{ color: '#059669', fontWeight: '600' }}>-{formatPrice(shippingDiscountValue)}</span>
                                </div>
                            )}
                            <div className="summary-row total">
                                <span>Thành tiền</span>
                                <span>{formatPrice(totalToPay)}</span>
                            </div>

                            <button
                                type="submit"
                                className="btn-confirm-order"
                                disabled={isSubmitting}
                            >
                                {isSubmitting
                                    ? <><Loader className="spin" size={16} /> Đang xử lý...</>
                                    : 'HOÀN TẤT ĐƠN HÀNG'}
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Checkout;
