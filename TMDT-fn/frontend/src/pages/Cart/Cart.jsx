import React, { useState } from 'react';
import { useCart } from '../../context/CartContext';
import { useAuth } from '../../context/AuthContext';
import { Link, useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import CartItem from '../../components/Cart/CartItem';
import AuthModal from '../../components/Auth/AuthModal';
import ConfirmModal from '../../components/ReUse/ConfirmModal';
import './Cart.css';
import { Trash2 } from 'lucide-react';

const Cart = () => {
    const { cartItems, updateQuantity, updateCartItemOptions, removeFromCart, removeMultipleFromCart, getCartTotal } = useCart();
    const { isAuthenticated } = useAuth();
    const navigate = useNavigate();
    const [coupon, setCoupon] = useState('');
    const [discount, setDiscount] = useState(0);
    const [showAuthModal, setShowAuthModal] = useState(false);
    const [selectedItems, setSelectedItems] = useState([]);
    const [isChecking, setIsChecking] = useState(false);
    const [confirmConfig, setConfirmConfig] = useState({ isOpen: false, title: '', message: '', onConfirm: null });

    const getItemKey = (item) => {
        const itemId = item.product?.id || item.product?.product_id;
        return `${itemId}-${item.capacity || ''}-${item.color_name || ''}`;
    };

    const handleSelectAll = (e) => {
        if (e.target.checked) {
            setSelectedItems(cartItems.map(item => getItemKey(item)));
        } else {
            setSelectedItems([]);
        }
    };

    const handleSelectItem = (item) => {
        const key = getItemKey(item);
        setSelectedItems(prev => {
            if (prev.includes(key)) {
                return prev.filter(k => k !== key);
            } else {
                return [...prev, key];
            }
        });
    };

    const handleRemoveItem = (index) => {
        setConfirmConfig({
            isOpen: true,
            title: 'Xác nhận xóa',
            message: 'Bạn có chắc chắn muốn xóa sản phẩm này khỏi giỏ hàng?',
            onConfirm: () => {
                removeFromCart(index);
                setConfirmConfig(prev => ({ ...prev, isOpen: false }));
            }
        });
    };

    const handleRemoveSelected = () => {
        if (selectedItems.length === 0) return;
        setConfirmConfig({
            isOpen: true,
            title: 'Xóa nhiều sản phẩm',
            message: `Bạn có chắc muốn xóa ${selectedItems.length} sản phẩm đã chọn?`,
            onConfirm: () => {
                const indicesToRemove = cartItems
                    .map((item, index) => selectedItems.includes(getItemKey(item)) ? index : -1)
                    .filter(index => index !== -1);

                removeMultipleFromCart(indicesToRemove);
                setSelectedItems([]);
                setConfirmConfig(prev => ({ ...prev, isOpen: false }));
            }
        });
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price);
    };

    const handleCheckout = async () => {
        if (selectedItems.length === 0) return;
        if (!isAuthenticated) {
            setShowAuthModal(true);
            return;
        }

        setIsChecking(true);
        try {
            const selectedCartItemsForCheckout = cartItems.filter(item => selectedItems.includes(getItemKey(item)));
            const payloadItems = selectedCartItemsForCheckout.map(item => ({
                variant_id: item.variant_id ?? item.product?.variant_id ?? item.product?.id,
                quantity: item.quantity
            }));

            const res = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/check-inventory`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ items: payloadItems })
            });

            const data = await res.json();

            if (data.success && data.all_in_stock) {
                // ─── ĐỒNG BỘ GIÁ TỪ DATABASE VÀO GIỎ HÀNG TRƯỚC KHI CHECKOUT ───
                const syncedCartItems = selectedCartItemsForCheckout.map(cartItem => {
                    const variantId = cartItem.variant_id ?? cartItem.product?.variant_id ?? cartItem.product?.id;
                    const dbItem = data.data.find(d => d.variant_id === variantId);

                    // Nếu Backend trả về thuộc tính current_price, ta sẽ ghi đè giá mới vào sản phẩm
                    if (dbItem && dbItem.current_price) {
                        return {
                            ...cartItem,
                            price: dbItem.current_price,
                            product: {
                                ...cartItem.product, // Giữ nguyên các thuộc tính khác (bao gồm cả ảnh img_thumb)
                                calculated_price: dbItem.current_price,
                                price: dbItem.current_price
                            }
                        };
                    }
                    return cartItem;
                });
                // ─────────────────────────────────────────────────────────────

                navigate('/checkout', {
                    state: {
                        selectedItems: syncedCartItems, // Gửi mảng đã được cập nhật giá sang trang Checkout
                        discount: discount,
                        coupon: coupon
                    }
                });
            } else {
                let msg = 'Một số sản phẩm không đủ tồn kho. Vui lòng kiểm tra lại!';
                if (data.data && Array.isArray(data.data)) {
                    const outOfStock = data.data.filter(i => !i.is_enough).map(i => `${i.product_name} ${i.color_name ? `(${i.color_name})` : ''}`).join(', ');
                    if (outOfStock) msg = `Sản phẩm không đủ tồn kho:\n ${outOfStock}`;
                }
                toast.dismiss('checkout-err');
                toast.error(msg, { id: 'checkout-err', duration: 4000 });
            }
        } catch (error) {
            console.error('Check inventory error:', error);
            toast.dismiss('checkout-err');
            toast.error('Lỗi khi kiểm tra tồn kho. Vui lòng thử lại sau!', { id: 'checkout-err', duration: 4000 });
        } finally {
            setIsChecking(false);
        }
    };

    const selectedCartItems = cartItems.filter(item => selectedItems.includes(getItemKey(item)));
    const selectedTotal = selectedCartItems.reduce((total, item) => {
        const price = item.price || item.product.calculated_price || item.product.price || 0;
        return total + (price * item.quantity);
    }, 0);

    const finalTotal = selectedTotal - discount;

    if (cartItems.length === 0) {
        return (
            <div className="cart-page empty">
                <div className="container" style={{ textAlign: 'center', padding: '50px' }}>
                    <h2>Giỏ hàng trống</h2>
                    <p>Không có sản phẩm nào trong giỏ hàng của bạn.</p>
                    <Link to="/" className="btn-primary" style={{ display: 'inline-block', width: 'auto', padding: '10px 20px', marginTop: '20px' }}>Tiếp tục mua sắm</Link>
                </div>
            </div>
        );
    }

    return (
        <div className="cart-page">
            <div className="container">
                <h1 className="page-title">Giỏ hàng của bạn</h1>

                <div className="cart-layout">
                    <div className="cart-list">
                        <div className="cart-list-header">
                            <label className="checkbox-container">
                                <input
                                    type="checkbox"
                                    checked={selectedItems.length === cartItems.length && cartItems.length > 0}
                                    onChange={handleSelectAll}
                                />
                                <span className="checkmark"></span>
                                <span className="label-text">Chọn tất cả ({cartItems.length})</span>
                            </label>
                            <button className="delete-btn" onClick={handleRemoveSelected}>
                                <Trash2 size={18} />
                            </button>
                        </div>
                        {cartItems.map((item, index) => (
                            <CartItem
                                key={index}
                                item={item}
                                index={index}
                                isSelected={selectedItems.includes(getItemKey(item))}
                                onSelect={() => handleSelectItem(item)}
                                onUpdateQuantity={updateQuantity}
                                onUpdateOptions={updateCartItemOptions}
                                onRemove={() => handleRemoveItem(index)}
                            />
                        ))}
                    </div>

                    <div className="cart-summary">
                        <div className="summary-row">
                            <span>Tạm tính ({selectedItems.length} sản phẩm):</span>
                            <span>{formatPrice(selectedTotal)}</span>
                        </div>
                        <div className="summary-row">
                            <span>Giảm giá:</span>
                            <span>-{formatPrice(discount)}</span>
                        </div>
                        <div className="summary-row total">
                            <span>Tổng tiền:</span>
                            <span>{formatPrice(finalTotal)}</span>
                        </div>

                        <button
                            className={`btn-checkout ${selectedItems.length === 0 || isChecking ? 'disabled' : ''}`}
                            onClick={handleCheckout}
                            disabled={selectedItems.length === 0 || isChecking}
                        >
                            {isChecking ? 'ĐANG KIỂM TRA...' : `TIẾN HÀNH ĐẶT HÀNG (${selectedItems.length})`}
                        </button>
                    </div>
                </div>
            </div>

            <AuthModal
                isOpen={showAuthModal}
                onClose={() => setShowAuthModal(false)}
                actionName="thanh toán đơn hàng"
                redirectPath="/login"
            />

            <ConfirmModal
                isOpen={confirmConfig.isOpen}
                onClose={() => setConfirmConfig(prev => ({ ...prev, isOpen: false }))}
                onConfirm={confirmConfig.onConfirm}
                title={confirmConfig.title}
                message={confirmConfig.message}
            />
        </div>
    );
};

export default Cart;