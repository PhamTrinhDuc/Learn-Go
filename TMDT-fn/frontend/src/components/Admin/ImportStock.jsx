import React, { useState, useEffect, useRef } from 'react';
import { Search, Trash2, X, Plus, Package, ChevronDown, ChevronRight, AlertCircle, CheckCircle2, ChevronUp } from 'lucide-react';
import './ImportStock.css';
import { handlePhoneChange, validatePhone } from '../../func/phoneValidation';

const ImportStock = ({ isOpen, onClose, products = [], categories = [], categoryMap = {}, onSuccess }) => {
    const [supplier, setSupplier] = useState({ name: '', address: '', phone: '' });
    const [phoneError, setPhoneError] = useState('');
    const [searchTerm, setSearchTerm] = useState('');
    const [categoryFilter, setCategoryFilter] = useState('Tất cả');
    const [note, setNote] = useState('');
    const [importItems, setImportItems] = useState([]);
    const [expandedProductId, setExpandedProductId] = useState(null);
    const [toasts, setToasts] = useState([]);
    const [showPriceConfirm, setShowPriceConfirm] = useState(false);
    const [pendingPriceUpdates, setPendingPriceUpdates] = useState([]);
    const [pricingStep, setPricingStep] = useState('ask');

    const leftScrollRef = useRef(null);
    const rightScrollRef = useRef(null);
    const mainScrollRef = useRef(null);

    useEffect(() => {
    }, [isOpen, importItems, searchTerm, categoryFilter]);

    const resetForm = () => {
        setSupplier({ name: '', address: '', phone: '' });
        setPhoneError('');
        setNote('');
        setImportItems([]);
        setSearchTerm('');
        setCategoryFilter('Tất cả');
        setExpandedProductId(null);
    };

    const scrollToTop = (ref) => {
        if (ref.current) {
            ref.current.scrollTo({ top: 0, behavior: 'smooth' });
        }
    };

    const scrollToBottom = (ref) => {
        if (ref.current) {
            ref.current.scrollTo({ top: ref.current.scrollHeight, behavior: 'smooth' });
        }
    };

    if (!isOpen) return null;

    const showToast = (message, type = 'success') => {
        const id = Date.now();
        setToasts(prev => [...prev, { id, message, type }]);
        setTimeout(() => {
            setToasts(prev => prev.filter(t => t.id !== id));
        }, 3000);
    };

    const filteredProducts = products.filter(p => {
        const matchesSearch = (p.name || '').toLowerCase().includes(searchTerm.toLowerCase());
        const pCatLabel = categoryMap[p.category] || p.category_name || p.category || '';
        const matchesCategory = categoryFilter === 'Tất cả' || pCatLabel === categoryFilter;
        return matchesSearch && matchesCategory;
    });

    const handleAddVariant = (product, variant) => {
        const vId = variant ? (variant.id ?? variant.variant_id ?? null) : null;
        const itemId = vId !== null ? `${product.id}-${vId}` : product.id;

        const existingItem = importItems.find(item => item.uniqueId === itemId);
        if (existingItem) {
            setImportItems(items => items.map(item =>
                item.uniqueId === itemId ? { ...item, importQty: item.importQty + 1 } : item
            ));
            showToast(`Đã tăng số lượng cho ${existingItem.name}`, 'warning');
            return;
        }

        const getVariantName = (v) => {
            if (!v) return '';
            const possibleNames = [v.color_name, v.variant_name, v.color, v.name, v.label];
            for (const n of possibleNames) {
                if (n && String(n).trim()) return String(n).trim();
            }
            return 'Mặc định';
        };

        const getPriceBaseOld = () => {
            if (variant) return variant.price_base ?? 0;
            const firstVar = product.originalData?.variants?.[0];
            if (firstVar) return firstVar.price_base ?? 0;
            return product.originalData?.base_price_numeric ?? product.base_price_numeric ?? 0;
        };

        const newItem = {
            uniqueId: itemId,
            productId: product.productId || product.id,
            vId: vId,
            name: product.originalData?.name || product.name,
            model: product.modelData?.model || product.originalData?.model || product.model || '',
            variantName: variant ? getVariantName(variant) : '',
            currentStock: variant ? (variant.stock ?? variant.quantity ?? 0) : product.stock,
            importQty: 1,
            priceImport: 0,
            currentSalePrice: variant ? (variant.price || variant.calculated_price || 0) : (product.price || product.calculated_price || 0),
            priceBaseOld: getPriceBaseOld()
        };

        setImportItems([...importItems, newItem]);
        showToast(`Đã thêm vào danh sách`, 'success');
    };

    const calculateNewPriceBase = (item) => {
        const oldPriceBase = parseInt(item.priceBaseOld) || 0;
        const oldQty = parseInt(item.currentStock) || 0;
        const importPrice = parseInt(item.priceImport) || 0;
        const importQty = parseInt(item.importQty) || 0;
        const totalQty = oldQty + importQty;
        if (totalQty <= 0) return importPrice;
        return Math.round(((oldPriceBase * oldQty) + (importPrice * importQty)) / totalQty);
    };

    const handleUpdatePendingSalePrice = (uniqueId, value) => {
        const cleanVal = parseInt(String(value).replace(/[^\d]/g, '')) || 0;
        setPendingPriceUpdates(prev => prev.map(item =>
            item.uniqueId === uniqueId ? { ...item, newSalePrice: Math.max(0, cleanVal) } : item
        ));
    };

    const handleUpdateQty = (uniqueId, qty) => {
        setImportItems(items => items.map(item =>
            item.uniqueId === uniqueId ? { ...item, importQty: Math.max(0, parseInt(qty) || 0) } : item
        ));
    };

    const handleUpdatePrice = (uniqueId, value) => {
        const cleanVal = parseInt(String(value).replace(/[^\d]/g, '')) || 0;
        setImportItems(items => items.map(item =>
            item.uniqueId === uniqueId ? { ...item, priceImport: Math.max(0, cleanVal) } : item
        ));
    };

    const handleRemoveItem = (uniqueId) => {
        setImportItems(items => items.filter(item => item.uniqueId !== uniqueId));
    };

    const handleConfirmPriceUpdate = async () => {
        try {
            const updates = pendingPriceUpdates.map(item => ({
                productId: item.productId,
                variantId: item.vId ?? null,
                newPrice: item.newSalePrice
            }));
            await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/update-price`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ updates })
            });
            showToast('Đã cập nhật giá bán thành công!', 'success');
        } catch (e) {
            showToast('Lỗi khi cập nhật giá bán!', 'error');
        } finally {
            setShowPriceConfirm(false);
            setPendingPriceUpdates([]);
            setTimeout(() => {
                resetForm();
                if (onSuccess) onSuccess(); else onClose();
            }, 1000);
        }
    };

    const handleSkipPriceUpdate = () => {
        setShowPriceConfirm(false);
        setPendingPriceUpdates([]);
        resetForm();
        if (onSuccess) onSuccess(); else onClose();
    };

    const handleImport = async () => {
        if (importItems.length === 0) {
            showToast('Vui lòng chọn sản phẩm!', 'error');
            return;
        }
        if (!supplier.name) {
            showToast('Chưa nhập nhà cung cấp!', 'error');
            return;
        }

        const payload = {
            supplier: {
                name: supplier.name,
                address: supplier.address,
                phone: supplier.phone
            },
            note: note,
            items: importItems.map(item => ({
                productId: item.productId,
                variantId: item.vId ?? null,
                quantity: item.importQty,
                currentStock: item.currentStock ?? 0,
                priceImport: item.priceImport || 0,
                name: item.name,
                model: item.model || '',
                variantName: item.variantName || ''
            }))
        };

        try {
            const res = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/import`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                showToast('Nhập kho thành công!', 'success');
                // Prepare pricing updates for the post-import confirmation dialog
                const updates = importItems.map(item => ({
                    ...item,
                    newSalePrice: item.currentSalePrice // Default to current sale price
                }));
                setPendingPriceUpdates(updates);
                setPricingStep('ask');
                setShowPriceConfirm(true);
            } else {
                const err = await res.json().catch(() => ({}));
                showToast(err.message || 'Lỗi khi gửi yêu cầu nhập kho!', 'error');
            }
        } catch (e) {
            console.error('Import stock error:', e);
            showToast('Lỗi kết nối server!', 'error');
        }
    };

    return (
        <div className="admin-modal-overlay">
            <div className="admin-toast-container">
                {toasts.map(toast => (
                    <div key={toast.id} className={`admin-toast admin-toast-${toast.type}`}>
                        <div className="admin-toast-icon">
                            {toast.type === 'success' && <CheckCircle2 size={20} color="#10b981" />}
                            {toast.type === 'warning' && <AlertCircle size={20} color="#f59e0b" />}
                            {toast.type === 'error' && <X size={20} color="#ef4444" />}
                        </div>
                        <div className="admin-toast-content">{toast.message}</div>
                    </div>
                ))}
            </div>

            <div className="admin-modal" style={{ maxWidth: '1200px', width: '95%', maxHeight: '95vh', display: 'flex', flexDirection: 'column' }}>
                <div className="admin-modal-header">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                        <Package size={24} style={{ color: 'var(--admin-primary)' }} />
                        <h2 style={{ margin: 0 }}>Nhập kho sản phẩm</h2>
                    </div>
                    <button className="admin-btn" onClick={onClose}><X size={20} /></button>
                </div>

                <div className="admin-modal-body" style={{ overflow: 'hidden', display: 'flex', flexDirection: 'column', position: 'relative', flex: 1, minHeight: 0 }}>
                    <div className="import-stock-container">
                        <div className="import-stock-header-grid">
                            <div className="form-field-group">
                                <label className="form-field-label">Nhà cung cấp</label>
                                <input
                                    type="text"
                                    className="form-field-input"
                                    placeholder="Tên..."
                                    value={supplier.name}
                                    onChange={(e) => setSupplier({ ...supplier, name: e.target.value })}
                                />
                            </div>
                            <div className="form-field-group">
                                <label className="form-field-label">Địa chỉ</label>
                                <input
                                    type="text"
                                    className="form-field-input"
                                    placeholder="Địa chỉ..."
                                    value={supplier.address}
                                    onChange={(e) => setSupplier({ ...supplier, address: e.target.value })}
                                />
                            </div>
                            <div className="form-field-group">
                                <label className="form-field-label">Điện thoại</label>
                                <input
                                    type="tel"
                                    className="form-field-input"
                                    placeholder="SĐT..."
                                    value={supplier.phone}
                                    style={{ borderColor: phoneError ? '#ef4444' : '' }}
                                    onChange={(e) => {
                                        const { cleaned, error } = handlePhoneChange(e.target.value);
                                        setSupplier({ ...supplier, phone: cleaned });
                                        setPhoneError(error);
                                    }}
                                />
                                {phoneError && (
                                    <div style={{ fontSize: '0.75rem', color: '#ef4444', marginTop: '4px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                                        <AlertCircle size={12} /> {phoneError}
                                    </div>
                                )}
                            </div>
                            <div className="form-field-group">
                                <label className="form-field-label">Ghi chú</label>
                                <input
                                    type="text"
                                    className="form-field-input"
                                    placeholder="Nhập ghi chú cho hóa đơn..."
                                    value={note}
                                    onChange={(e) => setNote(e.target.value)}
                                />
                            </div>
                        </div>

                        <div className="import-stock-main-layout" ref={mainScrollRef}>
                            <div className="import-stock-right-col">
                                <div className="import-list-table-container" ref={rightScrollRef}>
                                    <table className="import-list-table">
                                        <thead>
                                            <tr>
                                                <th>Sản phẩm</th>
                                                <th>Model</th>
                                                <th>Biến thể</th>
                                                <th style={{ textAlign: 'center' }}>Tồn</th>
                                                <th style={{ width: '130px', textAlign: 'center' }}>Giá nhập</th>
                                                <th style={{ width: '100px', textAlign: 'center' }}>Số lượng</th>
                                                <th></th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {importItems.map(item => (
                                                <tr key={item.uniqueId}>
                                                    <td style={{ fontWeight: 500 }}>{item.name}</td>
                                                    <td>{item.model || '-'}</td>
                                                    <td>
                                                        {(item.vId !== null && item.vId !== undefined) ? (
                                                            <span className="admin-badge admin-badge-outline">
                                                                {item.variantName || 'Mặc định'}
                                                            </span>
                                                        ) : (
                                                            <span style={{ color: '#94a3b8' }}>-</span>
                                                        )}
                                                    </td>
                                                    <td style={{ textAlign: 'center' }}>{item.currentStock}</td>
                                                    <td>
                                                        <div style={{ display: 'flex', justifyContent: 'center' }}>
                                                            <input
                                                                type="text"
                                                                className="import-qty-input"
                                                                placeholder="VNĐ"
                                                                style={{ width: '100px', textAlign: 'right' }}
                                                                value={item.priceImport > 0 ? new Intl.NumberFormat('vi-VN').format(item.priceImport) : ''}
                                                                onChange={(e) => handleUpdatePrice(item.uniqueId, e.target.value)}
                                                                onFocus={(e) => e.target.select()}
                                                                onClick={(e) => e.target.select()}
                                                            />
                                                        </div>
                                                    </td>
                                                    <td>
                                                        <div style={{ display: 'flex', justifyContent: 'center' }}>
                                                            <input
                                                                type="number"
                                                                className="import-qty-input"
                                                                min="1"
                                                                value={item.importQty}
                                                                onChange={(e) => handleUpdateQty(item.uniqueId, e.target.value)}
                                                                onFocus={(e) => e.target.select()}
                                                                onClick={(e) => e.target.select()}
                                                            />
                                                        </div>
                                                    </td>
                                                    <td>
                                                        <button
                                                            className="btn-remove-import-item"
                                                            onClick={() => handleRemoveItem(item.uniqueId)}
                                                        >
                                                            <Trash2 size={16} />
                                                        </button>
                                                    </td>
                                                </tr>
                                            ))}
                                            {importItems.length === 0 && (
                                                <tr>
                                                    <td colSpan="7" style={{ textAlign: 'center', padding: '40px', color: '#94a3b8' }}>
                                                        <Package size={48} style={{ opacity: 0.2, marginBottom: '12px' }} />
                                                        <div>Chưa có sản phẩm nào</div>
                                                    </td>
                                                </tr>
                                            )}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                            <div className="import-stock-left-col">
                                <div className="import-stock-search-filter">
                                    <div className="admin-search-wrapper" style={{ position: 'relative' }}>
                                        <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: '#94a3b8' }} />
                                        <input
                                            type="text"
                                            className="admin-form-input"
                                            placeholder="Tìm sản phẩm..."
                                            style={{ paddingLeft: '40px', width: '100%' }}
                                            value={searchTerm}
                                            onChange={(e) => setSearchTerm(e.target.value)}
                                        />
                                    </div>
                                    <select
                                        className="admin-form-input"
                                        value={categoryFilter}
                                        onChange={(e) => setCategoryFilter(e.target.value)}
                                    >
                                        <option value="Tất cả">Tất cả danh mục</option>
                                        {categories.filter(c => c !== 'Tất cả').map(cat => (
                                            <option key={cat} value={cat}>{cat}</option>
                                        ))}
                                    </select>
                                </div>

                                <div className="product-selection-list" ref={leftScrollRef}>
                                    {filteredProducts.map(p => {
                                        const variants = p.modelData?.variants || p.originalData?.variants || p.variants || [];
                                        const hasVariants = variants.length > 0;
                                        const isExpanded = expandedProductId === p.id;

                                        return (
                                            <div
                                                key={p.id}
                                                className={`product-item-card ${isExpanded ? 'expanded' : ''}`}
                                                onClick={() => {
                                                    if (hasVariants) setExpandedProductId(isExpanded ? null : p.id);
                                                    else handleAddVariant(p, null);
                                                }}
                                            >
                                                <div className="product-item-info">
                                                    <img
                                                        src={(() => {
                                                            const src = p.img_thumb || p.image || '';
                                                            return src.startsWith('http') ? src : `${import.meta.env.VITE_PHOTO_SERVER_API}${src}`;
                                                        })()}
                                                        className="product-item-thumb"
                                                        alt={p.name}
                                                    />
                                                    <div className="product-item-details">
                                                        <div className="product-item-name">{p.name}</div>
                                                        <div className="product-item-meta">
                                                            {(() => {
                                                                const allVars = (p.modelVariants || [p]).flatMap(m => m.variants?.length > 0 ? m.variants : [m]);
                                                                const totalStock = allVars.reduce((s, v) => s + (v.stock ?? v.quantity ?? 0), 0);
                                                                const catLabel = categoryMap[p.category] || p.category_name || p.category || '';
                                                                return `Tồn: ${totalStock} | ${catLabel}`;
                                                            })()}
                                                        </div>
                                                    </div>
                                                    {hasVariants && (
                                                        isExpanded ? <ChevronDown size={18} color="#64748b" /> : <ChevronRight size={18} color="#64748b" />
                                                    )}
                                                </div>

                                                {hasVariants && isExpanded && (
                                                    <div className="variant-selection-list" onClick={(e) => e.stopPropagation()}>
                                                        {variants.map((v, idx) => {
                                                            const vLabel = [v.color_name, v.variant_name, v.color, v.name, v.label]
                                                                .find(n => n && String(n).trim()) || `Biến thể ${idx + 1}`;
                                                            return (
                                                                <div key={v.id || v.variant_id || idx} className="variant-selection-item">
                                                                    <span>{vLabel} (Tồn: {v.stock ?? v.quantity ?? 0})</span>
                                                                    <button
                                                                        className="btn-add-variant"
                                                                        onClick={() => handleAddVariant(p, v)}
                                                                    >
                                                                        Chọn
                                                                    </button>
                                                                </div>
                                                            );
                                                        })}
                                                    </div>
                                                )}
                                            </div>
                                        );
                                    })}
                                    {filteredProducts.length === 0 && (
                                        <div style={{ textAlign: 'center', padding: '20px', color: '#94a3b8' }}>
                                            Không có sản phẩm phù hợp
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="admin-modal-footer">
                    <button className="admin-btn admin-btn-outline" onClick={onClose}>Hủy</button>
                    <button
                        className="admin-btn admin-btn-primary"
                        onClick={handleImport}
                    >
                        Nhập kho
                    </button>
                </div>
            </div>

            {/* Price Update Confirmation Modal */}
            {showPriceConfirm && (
                <div className="admin-modal-overlay" style={{ zIndex: 9999, position: 'fixed' }}>
                    <div className="admin-modal" style={{ maxWidth: pricingStep === 'edit' ? '800px' : '650px', width: '90%' }}>
                        <div className="admin-modal-header">
                            <h2>{pricingStep === 'edit' ? 'Cập nhật giá bán mới' : 'Cập nhật giá bán?'}</h2>
                        </div>
                        <div className="admin-modal-body">
                            {pricingStep === 'ask' ? (
                                <>
                                    <p style={{ marginBottom: '16px', fontSize: '0.95rem', lineHeight: '1.5', color: '#334155' }}>
                                        Nhập kho thành công! Sản phẩm có giá gốc và giá bán hiện tại như bên dưới. Bạn có muốn cập nhật giá bán mới không?
                                    </p>
                                    <div style={{ maxHeight: '300px', overflowY: 'auto', borderRadius: '8px', border: '1px solid #e2e8f0', marginBottom: '16px' }}>
                                        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                                            <thead>
                                                <tr style={{ background: '#f8fafc', borderBottom: '1px solid #cbd5e1' }}>
                                                    <th style={{ padding: '10px 12px', textAlign: 'left', fontWeight: '600' }}>Sản phẩm / Biến thể</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'right', fontWeight: '600', width: '120px' }}>Giá nhập cũ</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'right', fontWeight: '600', width: '120px' }}>Giá nhập mới</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'right', fontWeight: '600', width: '140px' }}>Giá bán hiện tại</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {pendingPriceUpdates.map(item => (
                                                    <tr key={item.uniqueId} style={{ borderBottom: '1px solid #f1f5f9' }}>
                                                        <td style={{ padding: '10px 12px', fontWeight: '500' }}>
                                                            {item.name}
                                                            {item.variantName && (
                                                                <span className="admin-badge admin-badge-outline" style={{ marginLeft: '8px', fontSize: '0.75rem' }}>
                                                                    {item.variantName}
                                                                </span>
                                                            )}
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'right', color: '#64748b' }}>
                                                            {new Intl.NumberFormat('vi-VN').format(item.priceBaseOld)}đ
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'right', color: '#10b981', fontWeight: '600' }}>
                                                            {new Intl.NumberFormat('vi-VN').format(calculateNewPriceBase(item))}đ
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'right', color: '#475569' }}>
                                                            {new Intl.NumberFormat('vi-VN').format(item.currentSalePrice)}đ
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                </>
                            ) : (
                                <>
                                    <p style={{ marginBottom: '16px', fontSize: '0.95rem', lineHeight: '1.5', color: '#334155' }}>
                                        Thay đổi giá bán mới cho các sản phẩm/biến thể dưới đây (Giá nhập cũ, Giá nhập mới và Giá bán hiện tại không được phép thay đổi):
                                    </p>
                                    <div style={{ maxHeight: '350px', overflowY: 'auto', borderRadius: '8px', border: '1px solid #e2e8f0', marginBottom: '16px' }}>
                                        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                                            <thead>
                                                <tr style={{ background: '#f8fafc', borderBottom: '1px solid #cbd5e1' }}>
                                                    <th style={{ padding: '10px 12px', textAlign: 'left', fontWeight: '600' }}>Sản phẩm / Biến thể</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'center', fontWeight: '600', width: '115px' }}>Giá nhập cũ</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'center', fontWeight: '600', width: '115px' }}>Giá nhập mới</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'center', fontWeight: '600', width: '115px' }}>Giá bán hiện tại</th>
                                                    <th style={{ padding: '10px 12px', textAlign: 'center', fontWeight: '600', width: '135px' }}>Giá bán mới</th>
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {pendingPriceUpdates.map(item => (
                                                    <tr key={item.uniqueId} style={{ borderBottom: '1px solid #f1f5f9' }}>
                                                        <td style={{ padding: '10px 12px', fontWeight: '500' }}>
                                                            {item.name}
                                                            {item.variantName && (
                                                                <span className="admin-badge admin-badge-outline" style={{ marginLeft: '8px', fontSize: '0.75rem' }}>
                                                                    {item.variantName}
                                                                </span>
                                                            )}
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'center' }}>
                                                            <input
                                                                type="text"
                                                                className="import-qty-input"
                                                                style={{ width: '95px', textAlign: 'right', backgroundColor: '#f1f5f9', cursor: 'not-allowed', color: '#64748b', border: '1px solid #cbd5e1' }}
                                                                value={`${new Intl.NumberFormat('vi-VN').format(item.priceBaseOld)}đ`}
                                                                disabled={true}
                                                            />
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'center' }}>
                                                            <input
                                                                type="text"
                                                                className="import-qty-input"
                                                                style={{ width: '95px', textAlign: 'right', backgroundColor: '#f1f5f9', cursor: 'not-allowed', color: '#10b981', border: '1px solid #cbd5e1', fontWeight: '600' }}
                                                                value={`${new Intl.NumberFormat('vi-VN').format(calculateNewPriceBase(item))}đ`}
                                                                disabled={true}
                                                            />
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'center' }}>
                                                            <input
                                                                type="text"
                                                                className="import-qty-input"
                                                                style={{ width: '95px', textAlign: 'right', backgroundColor: '#f1f5f9', cursor: 'not-allowed', color: '#64748b', border: '1px solid #cbd5e1' }}
                                                                value={`${new Intl.NumberFormat('vi-VN').format(item.currentSalePrice)}đ`}
                                                                disabled={true}
                                                            />
                                                        </td>
                                                        <td style={{ padding: '10px 12px', textAlign: 'center' }}>
                                                            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '4px' }}>
                                                                <input
                                                                    type="text"
                                                                    className="import-qty-input"
                                                                    style={{ width: '100px', textAlign: 'right', border: '1px solid #3b82f6', fontWeight: '600' }}
                                                                    placeholder="VNĐ"
                                                                    value={item.newSalePrice > 0 ? new Intl.NumberFormat('vi-VN').format(item.newSalePrice) : ''}
                                                                    onChange={(e) => handleUpdatePendingSalePrice(item.uniqueId, e.target.value)}
                                                                    onFocus={(e) => e.target.select()}
                                                                    onClick={(e) => e.target.select()}
                                                                />
                                                                <span style={{ fontSize: '0.85rem', color: '#64748b' }}>đ</span>
                                                            </div>
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                </>
                            )}
                        </div>
                        <div className="admin-modal-footer" style={{ gap: '12px', justifyContent: 'flex-end' }}>
                            {pricingStep === 'ask' ? (
                                <>
                                    <button className="admin-btn admin-btn-outline" onClick={handleSkipPriceUpdate}>
                                        Bỏ qua, giữ giá cũ
                                    </button>
                                    <button className="admin-btn admin-btn-primary" onClick={() => setPricingStep('edit')}>
                                        Có, cập nhật giá bán
                                    </button>
                                </>
                            ) : (
                                <>
                                    <button className="admin-btn admin-btn-outline" onClick={handleSkipPriceUpdate}>
                                        Hủy / Giữ giá cũ
                                    </button>
                                    <button className="admin-btn admin-btn-primary" onClick={handleConfirmPriceUpdate}>
                                        Cập nhật giá bán mới
                                    </button>
                                </>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ImportStock;
