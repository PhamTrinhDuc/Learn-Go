import React from 'react';
import { XCircle, Loader } from 'lucide-react';

const OrderDetailsModal = ({ selectedOrderDetails, showDetailsModal, setShowDetailsModal, getStatusBadge, loadingDetails }) => {
    if (!showDetailsModal) return null;

    return (
        <div className="admin-modal-overlay" style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', zIndex: 9999, display: 'flex', justifyContent: 'center', alignItems: 'center' }} onClick={() => setShowDetailsModal(false)}>
            <div className="admin-modal-content" style={{ background: '#fff', width: '90%', maxWidth: '850px', maxHeight: '90vh', overflowY: 'auto', borderRadius: '12px', padding: '24px', position: 'relative', boxShadow: '0 10px 25px rgba(0,0,0,0.2)' }} onClick={e => e.stopPropagation()}>
                <button style={{ position: 'absolute', top: 16, right: 16, background: 'transparent', border: 'none', cursor: 'pointer', padding: 4 }} onClick={() => setShowDetailsModal(false)}>
                    <XCircle size={24} color="#666" />
                </button>

                <h2 style={{ fontSize: '1.4rem', borderBottom: '1px solid #ebebeb', paddingBottom: 16, marginBottom: 24, color: 'var(--admin-text)' }}>
                    Chi tiết hóa đơn {selectedOrderDetails ? <span style={{ color: 'var(--admin-primary)' }}>#{selectedOrderDetails.order_id}</span> : '...'}
                </h2>

                {loadingDetails ? (
                    <div style={{ padding: '60px 0', textAlign: 'center' }}>
                        <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: 12 }} />
                        <p style={{ color: '#666' }}>Đang tải dữ liệu chi tiết...</p>
                    </div>
                ) : selectedOrderDetails ? (
                    <div className="order-details-body">
                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '20px', marginBottom: '24px' }}>

                            {/* 1. Khung Thông tin Cửa Hàng (Shop) */}
                            <div style={{ background: '#f8f9fa', padding: '16px', borderRadius: '8px', border: '1px solid #eaeaea' }}>
                                <h4 style={{ marginBottom: '12px', fontSize: '1.05rem', color: '#333' }}>Thông tin shop</h4>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem' }}><strong style={{ color: '#555', width: '85px', display: 'inline-block' }}>Cửa hàng:</strong> {selectedOrderDetails.shop_name || 'Hệ thống'}</p>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem' }}><strong style={{ color: '#555', width: '85px', display: 'inline-block' }}>Điện thoại:</strong> {selectedOrderDetails.shop_phone || 'Không có'}</p>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem', display: 'flex' }}><strong style={{ color: '#555', width: '85px', flexShrink: 0 }}>Địa chỉ:</strong> <span>{selectedOrderDetails.shop_address || 'Không có thông tin'}</span></p>
                            </div>

                            {/* 2. Khung Thông tin Giao Hàng (Khách hàng) */}
                            <div style={{ background: '#f8f9fa', padding: '16px', borderRadius: '8px', border: '1px solid #eaeaea' }}>
                                <h4 style={{ marginBottom: '12px', fontSize: '1.05rem', color: '#333' }}>Thông tin khách hàng</h4>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem' }}><strong style={{ color: '#555', width: '85px', display: 'inline-block' }}>Khách hàng:</strong> {selectedOrderDetails.customer_name || selectedOrderDetails.full_name}</p>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem' }}><strong style={{ color: '#555', width: '85px', display: 'inline-block' }}>Điện thoại:</strong> {selectedOrderDetails.customer_phone}</p>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem', display: 'flex' }}><strong style={{ color: '#555', width: '85px', flexShrink: 0 }}>Địa chỉ:</strong> <span>{selectedOrderDetails.customer_address}, {selectedOrderDetails.ward}, {selectedOrderDetails.province}</span></p>
                            </div>

                            {/* 3. Khung Thông tin Giao Dịch */}
                            <div style={{ background: '#f8f9fa', padding: '16px', borderRadius: '8px', border: '1px solid #eaeaea' }}>
                                <h4 style={{ marginBottom: '12px', fontSize: '1.05rem', color: '#333' }}>Thông tin giao dịch</h4>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem' }}><strong style={{ color: '#555', width: '85px', display: 'inline-block' }}>Ngày tạo:</strong> {new Date(selectedOrderDetails.order_date).toLocaleString('vi-VN')}</p>
                                <p style={{ margin: '4px 0', fontSize: '0.95rem' }}><strong style={{ color: '#555', width: '85px', display: 'inline-block' }}>Thanh toán:</strong> {selectedOrderDetails.payment_method}</p>
                                <div style={{ margin: '4px 0', fontSize: '0.95rem', display: 'flex', alignItems: 'center' }}>
                                    <strong style={{ color: '#555', width: '85px' }}>Trạng thái:</strong>
                                    <span className={`admin-badge ${getStatusBadge(selectedOrderDetails.status_name)}`}>
                                        {selectedOrderDetails.status_name}
                                    </span>
                                </div>
                            </div>
                        </div>

                        <h4 style={{ marginBottom: '12px', fontSize: '1.05rem', color: '#333' }}>Danh sách sản phẩm ({selectedOrderDetails.items?.length || 0})</h4>
                        <div style={{ overflowX: 'auto', border: '1px solid #eaeaea', borderRadius: '8px' }}>
                            <table className="admin-table" style={{ margin: 0, width: '100%' }}>
                                <thead style={{ background: '#f1f3f5' }}>
                                    <tr>
                                        <th style={{ padding: '12px', textAlign: 'left' }}>Tên sản phẩm</th>
                                        <th style={{ padding: '12px', textAlign: 'center' }}>Đơn giá</th>
                                        <th style={{ padding: '12px', textAlign: 'center' }}>S.Lượng</th>
                                        <th style={{ padding: '12px', textAlign: 'right' }}>Thành tiền</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {selectedOrderDetails.items && selectedOrderDetails.items.length > 0 ? (
                                        selectedOrderDetails.items.map((item, idx) => (
                                            <tr key={idx} style={{ borderBottom: '1px solid #eaeaea' }}>
                                                <td style={{ padding: '12px' }}>
                                                    <div style={{ fontWeight: 500, color: '#222' }}>{item.product_name}</div>
                                                    <div style={{ marginTop: '4px', display: 'flex', gap: '12px', fontSize: '0.85rem' }}>
                                                        {item.color_name && <span style={{ color: '#666', background: '#e9ecef', padding: '2px 6px', borderRadius: '4px' }}>Màu: <strong>{item.color_name}</strong></span>}
                                                        {item.capacity && <span style={{ color: '#666', background: '#e9ecef', padding: '2px 6px', borderRadius: '4px' }}>Dung lượng: <strong>{item.capacity}</strong></span>}
                                                    </div>
                                                </td>
                                                <td style={{ padding: '12px', textAlign: 'center', color: '#555' }}>{Number(item.unit_price).toLocaleString()}đ</td>
                                                <td style={{ padding: '12px', textAlign: 'center', fontWeight: '500' }}>{item.quantity}</td>
                                                <td style={{ padding: '12px', textAlign: 'right', fontWeight: 600, color: '#e11b1e' }}>{(Number(item.unit_price) * item.quantity).toLocaleString()}đ</td>
                                            </tr>
                                        ))
                                    ) : (
                                        <tr><td colSpan="4" style={{ textAlign: 'center', padding: '20px' }}>Không có chi tiết sản phẩm.</td></tr>
                                    )}
                                </tbody>
                            </table>
                        </div>

                        {(() => {
                            const subtotal = selectedOrderDetails.items ? selectedOrderDetails.items.reduce((s, i) => s + Number(i.unit_price) * i.quantity, 0) : Number(selectedOrderDetails.total_amount);
                            const totalShipping = selectedOrderDetails.items ? selectedOrderDetails.items.reduce((s, i) => s + Number(i.shipping_price || 0), 0) : 0;
                            const totalShippingDiscount = selectedOrderDetails.items ? selectedOrderDetails.items.reduce((s, i) => s + Number(i.shipping_support_price || 0), 0) : 0;
                            const totalProductDiscount = selectedOrderDetails.items ? selectedOrderDetails.items.reduce((s, i) => s + Number(i.product_support_price || 0), 0) : 0;
                            const grandTotal = subtotal + totalShipping - totalShippingDiscount - totalProductDiscount;

                            return (
                                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '24px', paddingTop: '16px', borderTop: '2px dashed #ddd' }}>
                                    <div style={{ width: '400px' }}>
                                        {selectedOrderDetails.note && (
                                            <div style={{ marginBottom: '16px', background: '#fff3cd', color: '#856404', padding: '12px', borderRadius: '4px', fontSize: '0.9rem' }}>
                                                <strong>Ghi chú:</strong> {selectedOrderDetails.note}
                                            </div>
                                        )}

                                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8, fontSize: '0.95rem', color: '#555' }}>
                                            <span>Tổng tiền hàng:</span>
                                            <span>{subtotal.toLocaleString()} đ</span>
                                        </div>

                                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8, fontSize: '0.95rem', color: '#555' }}>
                                            <span>Phí vận chuyển:</span>
                                            <span>{totalShipping.toLocaleString()} đ</span>
                                        </div>

                                        {(totalShippingDiscount > 0 || totalProductDiscount > 0) && (
                                            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8, fontSize: '0.95rem', color: 'var(--admin-success)' }}>
                                                <span>Tổng chiết khấu:</span>
                                                <span>-{(totalShippingDiscount + totalProductDiscount).toLocaleString()} đ</span>
                                            </div>
                                        )}

                                        <div style={{ borderTop: '1px solid #eaeaea', margin: '8px 0' }}></div>

                                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 12, fontSize: '1.2rem', fontWeight: 'bold' }}>
                                            <span style={{ color: '#444' }}>Thành tiền:</span>
                                            <span style={{ color: 'var(--admin-primary)', fontSize: '1.5rem' }}>{grandTotal.toLocaleString()} đ</span>
                                        </div>
                                    </div>
                                </div>
                            );
                        })()}
                    </div>
                ) : (
                    <div style={{ padding: '40px', textAlign: 'center', color: '#888' }}>Không tìm thấy thông tin sản phẩm.</div>
                )}

                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '24px' }}>
                    <button className="admin-btn admin-btn-outline" onClick={() => setShowDetailsModal(false)} style={{ padding: '8px 24px', fontWeight: '500' }}>Đóng lại</button>
                </div>
            </div>
        </div>
    );
};

export default OrderDetailsModal;
