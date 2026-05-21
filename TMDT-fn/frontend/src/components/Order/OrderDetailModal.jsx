import React from 'react';

const OrderDetailModal = ({ selectedOrder, onClose, formatPrice, getImageUrl }) => {
    if (!selectedOrder) return null;

    return (
        <div className="order-modal-overlay" onClick={onClose}>
            <div className="order-modal-content" onClick={e => e.stopPropagation()}>
                <div className="modal-header">
                    <h2>Chi tiết đơn hàng #{selectedOrder.order_id}</h2>
                    <button className="close-btn" onClick={onClose}>&times;</button>
                </div>
                <div className="modal-body">
                    <div className="modal-section">
                        <h3 className="section-title">Thông tin cửa hàng</h3>
                        <div className="delivery-card" style={{ backgroundColor: '#f9f9f9' }}>
                            <div className="delivery-row">
                                <strong>Cửa hàng:</strong> {selectedOrder.shop_name || 'Thế giới di động'}
                            </div>
                            <div className="delivery-row">
                                <strong>Điện thoại:</strong> {selectedOrder.shop_phone || 'Không có'}
                            </div>
                            <div className="delivery-row">
                                <strong>Địa chỉ:</strong> {selectedOrder.shop_address || 'Không có'}
                            </div>
                        </div>
                    </div>

                    <div className="modal-section">
                        <h3 className="section-title">Thông tin giao hàng</h3>
                        <div className="delivery-card">
                            <div className="delivery-row">
                                <strong>Người nhận:</strong> {selectedOrder.receiver_name}
                            </div>
                            <div className="delivery-row">
                                <strong>Điện thoại:</strong> {selectedOrder.receiver_phone}
                            </div>
                            <div className="delivery-row">
                                <strong>Địa chỉ:</strong> {selectedOrder.detail_address}
                            </div>
                            {selectedOrder.order_note && (
                                <div className="delivery-row">
                                    <strong>Ghi chú:</strong> {selectedOrder.order_note}
                                </div>
                            )}
                        </div>
                    </div>

                    <div className="modal-section">
                        <h3 className="section-title">Phương thức thanh toán</h3>
                        <div className="payment-info">
                            <span>{selectedOrder.payment_method === 'COD' ? 'Thanh toán khi nhận hàng (COD)' : selectedOrder.payment_method}</span>
                            <span className={`payment-status ${selectedOrder.payment_status?.toLowerCase()}`}>
                                {selectedOrder.payment_status}
                            </span>
                        </div>
                    </div>

                    <div className="modal-section">
                        <h3 className="section-title">Danh sách sản phẩm</h3>
                        <div className="modal-items">
                            {selectedOrder.items.map((item, idx) => (
                                <div key={idx} className="modal-item">
                                    <div className="item-img-small">
                                        <img
                                            src={getImageUrl(item.img_thumb)}
                                            alt={item.product_name}
                                            onError={(e) => {
                                                const fallback = 'https://placehold.co/200x200?text=No+Image';
                                                if (e.target.src !== fallback) e.target.src = fallback;
                                            }}
                                        />
                                    </div>
                                    <div className="item-info-mini">
                                        <div className="item-name-mini">{item.product_name}</div>
                                        <div className="item-sub-mini">
                                            {item.color_name && <span>{item.color_name}</span>}
                                            <span>x{item.quantity}</span>
                                        </div>
                                    </div>
                                    <div className="item-price-mini" style={{ marginLeft: 'auto', textAlign: 'right' }}>
                                        <div className="item-line-price" style={{ color: '#ee4d2d', fontWeight: '500' }}>
                                            {formatPrice(Number(item.unit_price) * item.quantity)}
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* ── HOÁ ĐƠN CHI TIẾT kiểu Shopee ────────────────────────────────── */}
                    {(() => {
                        const subtotal = selectedOrder.items.reduce(
                            (sum, item) => sum + Number(item.unit_price) * item.quantity, 0
                        );
                        const totalShipping = selectedOrder.items.reduce((s, i) => s + Number(i.shipping_price || 0), 0);
                        const totalShippingDiscount = selectedOrder.items.reduce((s, i) => s + Number(i.shipping_support_price || 0), 0);
                        const totalProductDiscount = selectedOrder.items.reduce((s, i) => s + Number(i.product_support_price || 0), 0);
                        const grandTotal = subtotal + totalShipping - totalShippingDiscount - totalProductDiscount;

                        // Parse note để lấy giao diện label (vd Vận chuyển: Hỏa tốc) cho đẹp
                        const note = selectedOrder.order_note || '';
                        const shippingMatch = note.match(/Vận chuyển:\s*([^|]+)/);
                        const couponMatch = note.match(/Mã giảm giá:\s*([^\s|]+)/);
                        const shippingLabel = shippingMatch ? shippingMatch[1].trim() : null;
                        const couponCode = couponMatch ? couponMatch[1].trim() : null;

                        return (
                            <div className="modal-invoice">
                                <div className="invoice-row">
                                    <span>Tổng tiền hàng</span>
                                    <span>{formatPrice(subtotal)}</span>
                                </div>

                                <div className="invoice-row">
                                    <span>Phí vận chuyển {shippingLabel ? `(${shippingLabel})` : ''}</span>
                                    <span>{formatPrice(totalShipping)}</span>
                                </div>

                                {totalShippingDiscount >= 0 && (
                                    <div className="invoice-row discount">
                                        <span>Ưu đãi phí vận chuyển</span>
                                        <span className="discount-val">-{formatPrice(totalShippingDiscount)}</span>
                                    </div>
                                )}

                                {totalProductDiscount >= 0 && (
                                    <div className="invoice-row discount">
                                        <span>
                                            Voucher từ Shop &nbsp;
                                            {couponCode && <span className="coupon-badge">{couponCode}</span>}
                                        </span>
                                        <span className="discount-val">-{formatPrice(totalProductDiscount)}</span>
                                    </div>
                                )}

                                <div className="invoice-divider" />

                                <div className="invoice-row grand-total">
                                    <span>Thành tiền</span>
                                    <span className="grand-total-price">{formatPrice(grandTotal)}</span>
                                </div>
                            </div>
                        );
                    })()}
                </div>
            </div>
        </div>
    );
};

export default OrderDetailModal;
