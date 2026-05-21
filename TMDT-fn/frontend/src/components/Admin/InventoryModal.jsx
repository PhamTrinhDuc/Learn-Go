import React from 'react';
import { Package } from 'lucide-react';

const InventoryModal = ({ isOpen, onClose, product, formatPrice = (p) => (p || 0).toLocaleString() + 'đ' }) => {
    if (!isOpen || !product) return null;

    return (
        <div className="admin-modal-overlay">
            <div className="admin-modal" style={{ maxWidth: '600px' }}>
                <div className="admin-modal-header">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                        <Package size={20} style={{ color: 'var(--admin-primary)' }} />
                        <div>
                            <h2 style={{ margin: 0, fontSize: '1.1rem' }}>Chi tiết tồn kho</h2>
                            <p style={{ margin: 0, fontSize: '0.85rem', color: 'var(--admin-text-muted)' }}>{product.name}</p>
                        </div>
                    </div>
                    <button className="admin-btn" onClick={onClose}>×</button>
                </div>
                <div className="admin-modal-body" style={{ padding: '0' }}>
                    <table className="admin-table" style={{ margin: 0, border: 'none' }}>
                        <thead style={{ background: '#f8fafc' }}>
                            <tr>
                                <th style={{ padding: '12px 20px', fontSize: '0.8rem', width: '40%' }}>Tên biến thể</th>
                                <th style={{ padding: '12px 20px', fontSize: '0.8rem', textAlign: 'right', width: '20%' }}>Giá bán</th>
                                <th style={{ padding: '12px 20px', fontSize: '0.8rem', textAlign: 'center', width: '15%' }}>Tồn kho</th>
                                <th style={{ padding: '12px 20px', fontSize: '0.8rem', textAlign: 'right', width: '25%' }}>Trạng thái</th>
                            </tr>
                        </thead>
                        <tbody>
                            {(() => {
                                const allRows = [];
                                const targetVariants = product.modelData
                                    ? [product.modelData]
                                    : (product.modelVariants && product.modelVariants.length > 0
                                        ? product.modelVariants
                                        : [product]);

                                targetVariants.forEach(v => {
                                    const storage = v.model || v.specs?.['Bộ nhớ trong'] || v.specs?.['Dung lượng'] || v.specs?.['Dung lượng lưu trữ'] || '';
                                    const ram = v.specs?.['Cấu hình & Bộ nhớ']?.['RAM'] || v.specs?.['RAM'] || '';
                                    const specs = [storage, ram ? `RAM ${ram}` : ''].filter(Boolean).join(' - ');

                                    if (v.variants && v.variants.length > 0) {
                                        v.variants.forEach(cv => {
                                            const color = cv.color_name || cv.variant_name || '';
                                            const stock = cv.stock ?? cv.quantity ?? cv.inventory ?? v.stock ?? 0;
                                            const price = cv.price || v.calculated_price || v.price || 0;
                                            const displayName = color || specs || 'Mặc định';

                                            allRows.push({
                                                displayName,
                                                sku: cv.sku || v.sku,
                                                price,
                                                stock,
                                                key: `${v.id || 'v'}-${color}-${Math.random()}`
                                            });
                                        });
                                    } else {
                                        const stock = v.stock ?? v.quantity ?? v.inventory ?? 0;
                                        const price = v.calculated_price || v.price || 0;
                                        allRows.push({
                                            displayName: specs || v.variant_name || 'Mặc định',
                                            sku: v.sku,
                                            price,
                                            stock,
                                            key: v.id || Math.random()
                                        });
                                    }
                                });

                                return allRows.map((row) => (
                                    <tr key={row.key}>
                                        <td style={{ padding: '12px 20px' }}>
                                            <div style={{ fontWeight: '600', color: 'var(--admin-text-dark)', fontSize: '0.9rem' }}>{row.displayName}</div>
                                        </td>
                                        <td style={{ padding: '12px 20px', textAlign: 'right', fontWeight: '500', color: '#e11d48' }}>{formatPrice(row.price)}</td>
                                        <td style={{ padding: '12px 20px', textAlign: 'center', fontWeight: 'bold', fontSize: '1rem' }}>{row.stock}</td>
                                        <td style={{ padding: '12px 20px', textAlign: 'right', whiteSpace: 'nowrap' }}>
                                            <span className={`admin-badge ${row.stock > 0 ? 'admin-badge-success' : 'admin-badge-danger'}`} style={{ fontSize: '0.7rem', padding: '2px 8px' }}>
                                                {row.stock > 0 ? 'Còn hàng' : 'Hết hàng'}
                                            </span>
                                        </td>
                                    </tr>
                                ));
                            })()}
                        </tbody>
                    </table>
                </div>
                <div className="admin-modal-footer">
                    <button className="admin-btn admin-btn-primary" style={{ width: '100%' }} onClick={onClose}>Đóng</button>
                </div>
            </div>
        </div>
    );
};

export default InventoryModal;
