import React, { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { Package, Truck, CheckCircle, Clock, ChevronRight, Search, ShoppingBag, Star } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import OrderDetailModal from '../../components/Order/OrderDetailModal';
import toast from 'react-hot-toast';
import './OrderHistory.css';

const API = import.meta.env.VITE_SERVER_API;

const OrderHistory = () => {
    const { user } = useAuth();
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState(0);
    const [selectedOrder, setSelectedOrder] = useState(null);
    const navigate = useNavigate();

    // States for review modal
    const [showReviewModal, setShowReviewModal] = useState(false);
    const [reviewOrder, setReviewOrder] = useState(null);
    const [reviewRating, setReviewRating] = useState(5);
    const [reviewComment, setReviewComment] = useState('');
    const [isSubmittingReview, setIsSubmittingReview] = useState(false);

    const tabs = [
        { id: 0, label: 'Tất cả' },
        { id: 1, label: 'Chờ xác nhận' },
        { id: 2, label: 'Chờ lấy hàng' },
        { id: 3, label: 'Chờ giao hàng' },
        { id: 4, label: 'Đã giao' },
        { id: 5, label: 'Trả hàng' },
        { id: 6, label: 'Đã hủy' }
    ];

    useEffect(() => {
        if (user?.id) {
            fetchOrders();
            const interval = setInterval(() => {
                if (document.visibilityState === 'visible') {
                    fetchOrders(true);
                }
            }, 20000);
            return () => clearInterval(interval);
        }
    }, [user]);

    const fetchOrders = async (isSilent = false) => {
        if (!isSilent) setLoading(true);
        try {
            const res = await fetch(`${API}/api/order/purchase-history/${user.id}`);
            const result = await res.json();
            if (result.success) {
                setOrders(result.data);
            }
        } catch (error) {
            console.error('Error fetching orders:', error);
        } finally {
            if (!isSilent) setLoading(false);
        }
    };

    const handleOpenReview = (order) => {
        setReviewOrder(order);
        setReviewRating(5);
        setReviewComment('');
        setShowReviewModal(true);
    };

    const handleCloseReview = () => {
        setShowReviewModal(false);
        setReviewOrder(null);
    };

    const handleRepurchase = (order) => {
        if (order.items && order.items.length > 0) {
            const firstItem = order.items[0];
            navigate(`/${firstItem.category_id || 'product'}/${firstItem.product_id}`);
        }
    };

    const submitReview = async () => {
        if (!reviewComment.trim()) {
            toast.error("Vui lòng nhập bình luận đánh giá");
            return;
        }
        if (reviewComment.length > 255) {
            toast.error("Bình luận không được vượt quá 255 ký tự");
            return;
        }
        if (isSubmittingReview) return;

        setIsSubmittingReview(true);
        try {
            const res = await fetch(`${API}/api/product/reviews`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    user_id: user.id,
                    order_id: reviewOrder.order_id,
                    rating: reviewRating,
                    comment: reviewComment
                })
            });

            const result = await res.json();
            if (result.success) {
                toast.success("Đánh giá đơn hàng thành công!");
                handleCloseReview();
                fetchOrders(true);
            } else {
                toast.error(result.message || "Gửi đánh giá thất bại");
            }
        } catch (error) {
            console.error("Error submitting review:", error);
            toast.error("Lỗi hệ thống khi gửi đánh giá");
        } finally {
            setIsSubmittingReview(false);
        }
    };

    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price);
    };

    const formatDate = (dateStr) => {
        return new Date(dateStr).toLocaleDateString('vi-VN', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const getImageUrl = (img) => {
        const fallback = 'https://placehold.co/200x200?text=No+Image';
        if (!img || img === '/placeholder.jpg' || img === '/placeholder.png') return fallback;
        if (img.startsWith('http')) return img;
        const photoApi = import.meta.env.VITE_PHOTO_SERVER_API || 'http://localhost:8081/images';
        const apiBase = photoApi.replace(/\/+$/, '');
        const imagePath = img.startsWith('/') ? img : `/${img}`;
        return `${apiBase}${imagePath}`;
    };

    const filteredOrders = activeTab === 0
        ? orders
        : orders.filter(order => order.status_id === activeTab);

    if (loading) {
        return (
            <div className="order-history-loading">
                <div className="spinner"></div>
                <p>Đang tải lịch sử đơn hàng...</p>
            </div>
        );
    }

    return (
        <div className="order-history-page">
            <div className="container">
                <div className="order-history-header">
                    <h1>Lịch sử đơn hàng</h1>
                </div>

                <div className="order-tabs">
                    {tabs.map(tab => (
                        <button
                            key={tab.id}
                            className={`tab-item ${activeTab === tab.id ? 'active' : ''}`}
                            onClick={() => setActiveTab(tab.id)}
                        >
                            {tab.label}
                        </button>
                    ))}
                </div>

                <div className="orders-list">
                    {filteredOrders.length === 0 ? (
                        <div className="empty-orders">
                            <ShoppingBag size={64} opacity={0.2} />
                            <p>Chưa có đơn hàng nào</p>
                            <button className="btn-shop-now" onClick={() => navigate('/')}>Mua sắm ngay</button>
                        </div>
                    ) : (
                        filteredOrders.map(order => {
                            const orderSubtotal = Number(order.total_amount);
                            const orderShipping = order.items.reduce((s, i) => s + Number(i.shipping_price || 0), 0);
                            const orderShippingDiscount = order.items.reduce((s, i) => s + Number(i.shipping_support_price || 0), 0);
                            const orderProductDiscount = order.items.reduce((s, i) => s + Number(i.product_support_price || 0), 0);
                            const finalGrandTotal = orderSubtotal + orderShipping - orderShippingDiscount - orderProductDiscount;

                            return (
                                <div key={order.order_id} className="order-card">
                                    <div className="order-card-header">
                                        <div className="order-id">
                                            <Package size={16} />
                                            <span>Đơn hàng: #{order.order_id}</span>
                                        </div>
                                        <div className={`order-status-badge status-${order.status_id}`}>
                                            {order.status_name}
                                        </div>
                                    </div>

                                    <div className="order-items">
                                        {order.items.map((item, idx) => (
                                            <div key={idx} className="order-item">
                                                <div className="item-img">
                                                    <img
                                                        src={getImageUrl(item.img_thumb)}
                                                        alt={item.product_name}
                                                        onError={(e) => {
                                                            const fallback = 'https://placehold.co/200x200?text=No+Image';
                                                            if (e.target.src !== fallback) {
                                                                e.target.src = fallback;
                                                            }
                                                        }}
                                                    />
                                                </div>
                                                <div className="item-details">
                                                    <h3 className="item-name">{item.product_name}</h3>
                                                    <div className="item-specs">
                                                        <span className="item-variant">
                                                            Phân loại:
                                                            {item.color_code && (
                                                                <span className="color-dot" style={{ backgroundColor: item.color_code }}></span>
                                                            )}
                                                            {item.color_name}
                                                        </span>
                                                        <span className="item-qty">x{item.quantity}</span>
                                                    </div>
                                                </div>
                                                <div className="item-price">
                                                    {formatPrice(item.unit_price)}
                                                </div>
                                            </div>
                                        ))}
                                    </div>

                                    <div className="order-card-footer">
                                        <div className="order-date">
                                            <Clock size={14} />
                                            <span>Ngày đặt: {formatDate(order.order_date)}</span>
                                        </div>
                                        <div className="order-summary">
                                            <div className="total-label">Thành tiền:</div>
                                            <div className="total-amount">{formatPrice(finalGrandTotal)}</div>
                                        </div>
                                    </div>
                                    <div className="order-card-actions">
                                        {order.status_id === 4 ? (
                                            order.is_reviewed ? (
                                                <>
                                                    <button className="btn-detail" onClick={() => setSelectedOrder(order)}> Xem chi tiết </button>
                                                    <button className="btn-repurchase" onClick={() => handleRepurchase(order)}>Mua lại</button>
                                                </>
                                            ) : (
                                                <button className="btn-confirm-review" onClick={() => handleOpenReview(order)}>Xác nhận nhận hàng</button>
                                            )
                                        ) : (
                                            <button className="btn-detail" onClick={() => setSelectedOrder(order)}> Xem chi tiết </button>
                                        )}
                                    </div>
                                </div>
                            );
                        })
                    )}
                </div>
            </div>

            {selectedOrder && (
                <OrderDetailModal
                    selectedOrder={selectedOrder}
                    onClose={() => setSelectedOrder(null)}
                    formatPrice={formatPrice}
                    getImageUrl={getImageUrl}
                />
            )}

            {showReviewModal && reviewOrder && (
                <div className="order-modal-overlay">
                    <div className="order-modal-content review-modal">
                        <div className="modal-header">
                            <h2>Đánh giá đơn hàng #{reviewOrder.order_id}</h2>
                            <button className="close-btn" onClick={handleCloseReview}>&times;</button>
                        </div>
                        <div className="modal-body">
                            <div className="review-items-list" style={{ marginBottom: '20px', borderBottom: '1px solid #f1f5f9', paddingBottom: '15px' }}>
                                {reviewOrder.items.map((item, idx) => (
                                    <div key={idx} className="review-item-row" style={{ display: 'flex', alignItems: 'center', gap: '15px', marginBottom: '10px' }}>
                                        <img src={getImageUrl(item.img_thumb)} alt={item.product_name} style={{ width: '50px', height: '50px', objectFit: 'contain', border: '1px solid #eee', borderRadius: '4px' }} />
                                        <div className="review-item-info">
                                            <div className="review-item-name" style={{ fontSize: '14px', fontWeight: '500' }}>{item.product_name}</div>
                                            {item.color_name && <div className="review-item-variant" style={{ fontSize: '12px', color: '#888' }}>Phân loại: {item.color_name}</div>}
                                        </div>
                                    </div>
                                ))}
                            </div>

                            <div className="review-stars-section" style={{ textAlign: 'center', marginBottom: '20px' }}>
                                <p className="rating-label" style={{ fontWeight: '500', color: '#555', marginBottom: '10px' }}>Vui lòng đánh giá dịch vụ và sản phẩm:</p>
                                <div className="rating-stars" style={{ display: 'flex', justifyContent: 'center', gap: '8px', marginBottom: '8px' }}>
                                    {[1, 2, 3, 4, 5].map((star) => (
                                        <Star
                                            key={star}
                                            size={32}
                                            fill={star <= reviewRating ? "#fb6e2e" : "none"}
                                            stroke={star <= reviewRating ? "#fb6e2e" : "#ccc"}
                                            onClick={() => setReviewRating(star)}
                                            style={{ cursor: 'pointer', transition: 'transform 0.1s' }}
                                        />
                                    ))}
                                </div>
                                <span className="rating-text" style={{ fontSize: '14px', fontWeight: '600', color: '#fb6e2e' }}>
                                    {reviewRating === 5 && 'Tuyệt vời'}
                                    {reviewRating === 4 && 'Rất tốt'}
                                    {reviewRating === 3 && 'Bình thường'}
                                    {reviewRating === 2 && 'Không hài lòng'}
                                    {reviewRating === 1 && 'Tệ'}
                                </span>
                            </div>

                            <div className="review-comment-section" style={{ marginBottom: '20px', position: 'relative' }}>
                                <textarea
                                    value={reviewComment}
                                    onChange={(e) => setReviewComment(e.target.value.slice(0, 255))}
                                    placeholder="Chia sẻ cảm nhận của bạn về sản phẩm và dịch vụ (tối đa 255 ký tự)..."
                                    rows={4}
                                    className="review-textarea"
                                    style={{
                                        width: '100%',
                                        padding: '12px',
                                        borderRadius: '6px',
                                        border: '1px solid #d1d5db',
                                        fontSize: '14px',
                                        lineHeight: '1.5',
                                        resize: 'none',
                                        boxSizing: 'border-box'
                                    }}
                                />
                                <div className="char-counter" style={{ textAlign: 'right', fontSize: '12px', color: '#888', marginTop: '4px' }}>
                                    {reviewComment.length}/255
                                </div>
                            </div>

                            <div className="review-modal-actions" style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
                                <button className="btn-cancel-review" onClick={handleCloseReview} style={{ padding: '8px 16px', border: '1px solid #d1d5db', background: 'white', borderRadius: '4px', cursor: 'pointer' }}>Hủy</button>
                                <button
                                    className="btn-submit-review"
                                    onClick={submitReview}
                                    disabled={!reviewComment.trim()}
                                    style={{
                                        padding: '8px 16px',
                                        border: 'none',
                                        background: reviewComment.trim() ? '#ee4d2d' : '#ccc',
                                        color: 'white',
                                        borderRadius: '4px',
                                        cursor: reviewComment.trim() ? 'pointer' : 'not-allowed',
                                        fontWeight: '500'
                                    }}
                                >
                                    Gửi đánh giá
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OrderHistory;
