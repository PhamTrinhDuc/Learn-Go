import React, { useState, useEffect } from 'react';
import { Search, Plus, Edit, Trash2, CheckCircle, ShieldAlert, Loader, Tag, Calendar, Percent, AlertCircle, ShoppingBag } from 'lucide-react';
import toast from 'react-hot-toast';

const PromotionManager = () => {
    const [promotions, setPromotions] = useState([]);
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState('Tất cả'); // 'Tất cả', 'Chưa diễn ra', 'Đang chạy', 'Đã hết hạn'

    // Modal & Form state
    const [showModal, setShowModal] = useState(false);
    const [modalMode, setModalMode] = useState('create'); // 'create' or 'edit'
    const [editingPromotionId, setEditingPromotionId] = useState(null);
    const [formData, setFormData] = useState({
        product_id: '',
        discount_percent: '',
        start_date: '',
        end_date: ''
    });

    // Product search inside Form
    const [formProductSearch, setFormProductSearch] = useState('');
    const [selectedProduct, setSelectedProduct] = useState(null);
    const [showProductDropdown, setShowProductDropdown] = useState(false);

    const [formError, setFormError] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // Delete state
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [promotionToDelete, setPromotionToDelete] = useState(null);
    const [deleting, setDeleting] = useState(false);

    const API = import.meta.env.VITE_SERVER_API;

    const fetchData = async () => {
        setLoading(true);
        try {
            // 1. Fetch promotions
            const promoRes = await fetch(`${API}/api/promotion`);
            const promoData = await promoRes.json();
            
            // 2. Fetch products for selection
            const prodRes = await fetch(`${API}/api/product/all-products?limit=2000`);
            const prodData = await prodRes.json();

            if (promoData.success) {
                setPromotions(promoData.data);
            } else {
                toast.error(promoData.message || 'Lỗi lấy danh sách chương trình khuyến mãi');
            }

            if (prodData.success) {
                // Ensure unique products
                const uniqueProds = [];
                const seenIds = new Set();
                prodData.data.forEach(p => {
                    if (!seenIds.has(p.product_id)) {
                        seenIds.add(p.product_id);
                        uniqueProds.push(p);
                    }
                });
                setProducts(uniqueProds);
            }
        } catch (err) {
            console.error('Error loading promotion data:', err);
            toast.error('Lỗi kết nối máy chủ');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, []);

    const formatPrice = (price) =>
        new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price);

    const getMinStartDate = () => {
        const now = new Date();
        const year = now.getFullYear();
        const month = String(now.getMonth() + 1).padStart(2, '0');
        const day = String(now.getDate()).padStart(2, '0');
        const hours = String(now.getHours()).padStart(2, '0');
        const minutes = String(now.getMinutes()).padStart(2, '0');
        return `${year}-${month}-${day}T${hours}:${minutes}`;
    };

    const formatDateLocal = (dateString) => {
        if (!dateString) return '---';
        const d = new Date(dateString);
        return d.toLocaleString('vi-VN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const getPromotionStatus = (promo) => {
        const now = new Date();
        const start = new Date(promo.start_date);
        const end = new Date(promo.end_date);

        if (now < start) {
            return { label: 'Chưa diễn ra', color: 'admin-badge-warning', type: 'future' };
        } else if (now >= start && now <= end) {
            return { label: 'Đang chạy', color: 'admin-badge-success', type: 'active' };
        } else {
            return { label: 'Đã kết thúc', color: 'admin-badge-danger', type: 'expired' };
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;

        if (name === 'start_date' && value) {
            const selected = new Date(value);
            const now = new Date();
            // Allow a small buffer of 1 minute
            if (selected < new Date(now.getTime() - 60000)) {
                toast.error('Thời gian bắt đầu không thể ở trong quá khứ!');
                setFormError('Thời gian bắt đầu không thể ở trong quá khứ!');
                const minVal = getMinStartDate();
                setFormData(prev => {
                    const updated = {
                        ...prev,
                        start_date: minVal
                    };
                    if (updated.end_date && new Date(updated.end_date) <= new Date(minVal)) {
                        updated.end_date = minVal;
                    }
                    return updated;
                });
                return;
            } else {
                setFormError('');
                setFormData(prev => {
                    const updated = {
                        ...prev,
                        [name]: value
                    };
                    if (updated.end_date && new Date(updated.end_date) <= new Date(value)) {
                        updated.end_date = value;
                    }
                    return updated;
                });
                return;
            }
        }

        if (name === 'end_date' && value) {
            const end = new Date(value);
            const startLimitStr = formData.start_date || getMinStartDate();
            const startLimit = new Date(startLimitStr);
            if (end <= startLimit) {
                toast.error('Thời gian kết thúc phải sau thời gian bắt đầu!');
                setFormError('Thời gian kết thúc phải sau thời gian bắt đầu!');
                setFormData(prev => ({
                    ...prev,
                    end_date: startLimitStr
                }));
                return;
            } else {
                setFormError('');
            }
        }

        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleProductSelect = (product) => {
        setSelectedProduct(product);
        setFormData(prev => ({
            ...prev,
            product_id: product.product_id
        }));
        setFormProductSearch(product.name);
        setShowProductDropdown(false);
    };

    const openCreateModal = () => {
        setModalMode('create');
        setEditingPromotionId(null);
        setSelectedProduct(null);
        setFormProductSearch('');
        setFormData({
            product_id: '',
            discount_percent: '',
            start_date: '',
            end_date: ''
        });
        setFormError('');
        setShowModal(true);
    };

    const openEditModal = (promo) => {
        setModalMode('edit');
        setEditingPromotionId(promo.promotion_id);

        const matchedProd = products.find(p => p.product_id === promo.product_id);
        if (matchedProd) {
            setSelectedProduct(matchedProd);
            setFormProductSearch(matchedProd.name);
        } else {
            setSelectedProduct({ product_id: promo.product_id, name: promo.product_name, price: promo.product_price });
            setFormProductSearch(promo.product_name);
        }

        // Format ISO dates to datetime-local input format (YYYY-MM-DDTHH:MM)
        const startIso = new Date(promo.start_date).toISOString().slice(0, 16);
        const endIso = new Date(promo.end_date).toISOString().slice(0, 16);

        setFormData({
            product_id: promo.product_id,
            discount_percent: promo.discount_percent.toString(),
            start_date: startIso,
            end_date: endIso
        });
        setFormError('');
        setShowModal(true);
    };

    const validateForm = () => {
        const { product_id, discount_percent, start_date, end_date } = formData;
        if (!product_id || !discount_percent || !start_date || !end_date) {
            return 'Vui lòng điền đầy đủ các thông tin bắt buộc!';
        }

        const discount = parseInt(discount_percent, 10);
        if (isNaN(discount) || discount < 1 || discount > 99) {
            return 'Phần trăm giảm giá phải từ 1 đến 99%!';
        }

        const start = new Date(start_date);
        const end = new Date(end_date);
        const now = new Date();

        if (start >= end) {
            return 'Thời gian bắt đầu phải trước thời gian kết thúc!';
        }

        if (modalMode === 'create' && start < now) {
            return 'Thời gian bắt đầu không thể ở trong quá khứ!';
        }

        return '';
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const error = validateForm();
        if (error) {
            setFormError(error);
            return;
        }

        setSubmitting(true);
        setFormError('');

        const payload = {
            ...formData,
            discount_percent: parseInt(formData.discount_percent, 10)
        };

        const url = modalMode === 'create'
            ? `${API}/api/promotion/create`
            : `${API}/api/promotion/edit/${editingPromotionId}`;
        
        const method = modalMode === 'create' ? 'POST' : 'PUT';

        try {
            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            const data = await res.json();

            if (data.success) {
                toast.success(data.message || 'Lưu chương trình thành công!');
                setShowModal(false);
                fetchData();
            } else {
                setFormError(data.message || 'Lỗi khi thực hiện thao tác!');
            }
        } catch (err) {
            console.error(err);
            setFormError('Lỗi kết nối máy chủ!');
        } finally {
            setSubmitting(false);
        }
    };

    const handleDeleteClick = (promo) => {
        setPromotionToDelete(promo);
        setShowDeleteModal(true);
    };

    const confirmDelete = async () => {
        if (!promotionToDelete) return;
        setDeleting(true);
        try {
            const res = await fetch(`${API}/api/promotion/delete/${promotionToDelete.promotion_id}`, {
                method: 'DELETE'
            });
            const data = await res.json();

            if (data.success) {
                toast.success(data.message || 'Xóa khuyến mãi thành công!');
                setShowDeleteModal(false);
                setPromotionToDelete(null);
                fetchData();
            } else {
                toast.error(data.message || 'Lỗi khi xóa khuyến mãi!');
            }
        } catch (err) {
            console.error(err);
            toast.error('Lỗi kết nối máy chủ!');
        } finally {
            setDeleting(false);
        }
    };

    // Filter and Search logic
    const filteredPromotions = promotions.filter(p => {
        const matchesSearch = (p.product_name || '').toLowerCase().includes(searchTerm.toLowerCase()) || 
                              (p.product_id || '').toLowerCase().includes(searchTerm.toLowerCase());
        
        const status = getPromotionStatus(p);
        let matchesStatus = true;
        if (statusFilter === 'Chưa diễn ra') {
            matchesStatus = status.type === 'future';
        } else if (statusFilter === 'Đang chạy') {
            matchesStatus = status.type === 'active';
        } else if (statusFilter === 'Đã hết hạn') {
            matchesStatus = status.type === 'expired';
        }

        return matchesSearch && matchesStatus;
    });

    // Form product search list
    const filteredFormProducts = products.filter(p => 
        p.name.toLowerCase().includes(formProductSearch.toLowerCase())
    );

    return (
        <div className="admin-promotion-manager">
            <div className="admin-card">
                <div className="admin-card-header">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        <Tag size={24} style={{ color: 'var(--admin-primary)' }} />
                        <h2>Chiến dịch Khuyến mãi Sản phẩm</h2>
                    </div>
                    <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'center' }}>
                        <div className="admin-search-wrapper" style={{ position: 'relative' }}>
                            <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--admin-text-muted)' }} />
                            <input
                                type="text"
                                className="admin-form-input"
                                placeholder="Tìm sản phẩm..."
                                style={{ paddingLeft: '40px', width: '220px' }}
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>

                        <select
                            className="admin-form-input"
                            style={{ width: '160px', height: '40px' }}
                            value={statusFilter}
                            onChange={(e) => setStatusFilter(e.target.value)}
                        >
                            <option value="Tất cả">Tất cả trạng thái</option>
                            <option value="Chưa diễn ra">Chưa diễn ra</option>
                            <option value="Đang chạy">Đang diễn ra</option>
                            <option value="Đã hết hạn">Đã kết thúc</option>
                        </select>

                        <button
                            className="admin-btn admin-btn-primary"
                            onClick={openCreateModal}
                            style={{ height: '40px', padding: '0 16px', borderRadius: '8px', display: 'flex', alignItems: 'center', gap: '8px' }}
                        >
                            <Plus size={18} />
                            <span>Tạo khuyến mãi mới</span>
                        </button>
                    </div>
                </div>

                <div className="admin-table-container">
                    {loading ? (
                        <div style={{ padding: '40px', textAlign: 'center' }}>
                            <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: '12px' }} />
                            <p>Đang tải dữ liệu chương trình khuyến mãi...</p>
                        </div>
                    ) : (
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th>Sản phẩm áp dụng</th>
                                    <th>Giá gốc</th>
                                    <th>% Giảm giá</th>
                                    <th>Giá khuyến mãi</th>
                                    <th>Thời gian áp dụng</th>
                                    <th>Trạng thái</th>
                                    <th>Thao tác</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredPromotions.map((promo) => {
                                    const status = getPromotionStatus(promo);
                                    const price = Number(promo.product_price || 0);
                                    const discountedPrice = Math.round(price * (1 - promo.discount_percent / 100));

                                    return (
                                        <tr key={promo.promotion_id}>
                                            <td>
                                                <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                                    <div style={{ width: '36px', height: '36px', borderRadius: '6px', background: 'var(--admin-bg-light)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--admin-primary)' }}>
                                                        <ShoppingBag size={18} />
                                                    </div>
                                                    <div>
                                                        <div style={{ fontWeight: '600' }}>{promo.product_name || 'Sản phẩm đã bị xóa'}</div>
                                                        <small style={{ color: 'var(--admin-text-muted)', fontSize: '0.75rem' }}>ID: {promo.product_id || 'N/A'}</small>
                                                    </div>
                                                </div>
                                            </td>
                                            <td style={{ color: 'var(--admin-text-muted)', textDecoration: 'line-through' }}>
                                                {formatPrice(price)}
                                            </td>
                                            <td>
                                                <span style={{ fontWeight: '700', color: 'var(--admin-danger)', display: 'flex', alignItems: 'center', gap: '2px' }}>
                                                    <Percent size={14} /> -{promo.discount_percent}%
                                                </span>
                                            </td>
                                            <td style={{ fontWeight: '700', color: '#10b981' }}>
                                                {formatPrice(discountedPrice)}
                                            </td>
                                            <td>
                                                <div style={{ fontSize: '0.8rem', color: 'var(--admin-text-muted)', display: 'flex', flexDirection: 'column', gap: '2px' }}>
                                                    <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                                                        <Calendar size={12} /> Bắt đầu: {formatDateLocal(promo.start_date)}
                                                    </span>
                                                    <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                                                        <Calendar size={12} /> Kết thúc: {formatDateLocal(promo.end_date)}
                                                    </span>
                                                </div>
                                            </td>
                                            <td>
                                                <span className={`admin-badge ${status.color}`}>
                                                    {status.label}
                                                </span>
                                            </td>
                                            <td>
                                                <div style={{ display: 'flex', gap: '8px' }}>
                                                    <button
                                                        className="admin-btn admin-btn-outline admin-btn-sm"
                                                        title="Sửa"
                                                        onClick={() => openEditModal(promo)}
                                                    >
                                                        <Edit size={14} />
                                                    </button>
                                                    <button
                                                        className="admin-btn admin-btn-outline admin-btn-sm"
                                                        style={{ color: 'var(--admin-danger)' }}
                                                        title="Xóa"
                                                        onClick={() => handleDeleteClick(promo)}
                                                    >
                                                        <Trash2 size={14} />
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    );
                                })}
                                {filteredPromotions.length === 0 && (
                                    <tr>
                                        <td colSpan="7" style={{ textAlign: 'center', padding: '32px' }}>
                                            <AlertCircle size={24} style={{ color: 'var(--admin-text-muted)', margin: '0 auto 8px' }} />
                                            <p style={{ color: 'var(--admin-text-muted)' }}>Không tìm thấy chương trình khuyến mãi nào!</p>
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    )}
                </div>
            </div>

            {/* CREATE / EDIT MODAL */}
            {showModal && (
                <div className="admin-modal-overlay">
                    <div className="admin-modal" style={{ maxWidth: '550px', width: '100%' }}>
                        <div className="admin-modal-header">
                            <h2>{modalMode === 'create' ? 'Tạo chương trình khuyến mãi' : 'Cập nhật Khuyến mãi'}</h2>
                            <button className="admin-btn" onClick={() => setShowModal(false)}>×</button>
                        </div>
                        <form onSubmit={handleSubmit} noValidate>
                            <div className="admin-modal-body" style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                                
                                {formError && (
                                    <div style={{ fontSize: '0.85rem', color: 'var(--admin-danger)', background: 'rgba(239, 68, 68, 0.05)', padding: '10px', borderRadius: '6px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                                        <ShieldAlert size={16} /> {formError}
                                    </div>
                                )}

                                <div className="admin-form-group" style={{ position: 'relative' }}>
                                    <label className="admin-form-label">Chọn sản phẩm áp dụng*</label>
                                    <input
                                        type="text"
                                        className="admin-form-input"
                                        placeholder="Tìm tên sản phẩm cần giảm giá..."
                                        value={formProductSearch}
                                        onChange={(e) => {
                                            setFormProductSearch(e.target.value);
                                            setShowProductDropdown(true);
                                        }}
                                        onFocus={() => setShowProductDropdown(true)}
                                        required
                                        disabled={modalMode === 'edit'} // Lock product on edit
                                    />
                                    {showProductDropdown && modalMode === 'create' && (
                                        <div style={{
                                            position: 'absolute', top: '100%', left: 0, right: 0,
                                            maxHeight: '200px', overflowY: 'auto', background: '#fff',
                                            border: '1px solid #d9d9d9', borderRadius: '6px', zIndex: 10,
                                            boxShadow: '0 4px 6px -1px rgba(0,0,0,0.1)'
                                        }}>
                                            {filteredFormProducts.slice(0, 15).map(p => (
                                                <div
                                                    key={p.product_id}
                                                    onClick={() => handleProductSelect(p)}
                                                    style={{ padding: '8px 12px', cursor: 'pointer', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', fontSize: '0.9rem' }}
                                                    onMouseDown={(e) => e.preventDefault()} // Keep focus
                                                >
                                                    <span>{p.name}</span>
                                                    <span style={{ color: 'var(--admin-primary)', fontWeight: '600' }}>{formatPrice(p.price)}</span>
                                                </div>
                                            ))}
                                            {filteredFormProducts.length === 0 && (
                                                <div style={{ padding: '8px 12px', color: 'var(--admin-text-muted)', textAlign: 'center' }}>Không thấy sản phẩm nào</div>
                                            )}
                                        </div>
                                    )}
                                </div>

                                {selectedProduct && (
                                    <div style={{ background: 'var(--admin-bg-light)', padding: '12px', borderRadius: '6px', borderLeft: '4px solid var(--admin-primary)' }}>
                                        <div style={{ fontWeight: '600', fontSize: '0.9rem' }}>Sản phẩm đã chọn: {selectedProduct.name}</div>
                                        <div style={{ fontSize: '0.85rem', color: 'var(--admin-text-muted)', marginTop: '4px' }}>
                                            Giá niêm yết gốc: <strong>{formatPrice(selectedProduct.price)}</strong>
                                        </div>
                                    </div>
                                )}

                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Thời gian bắt đầu*</label>
                                        <input
                                            type="datetime-local"
                                            className="admin-form-input"
                                            name="start_date"
                                            value={formData.start_date}
                                            onChange={handleInputChange}
                                            min={getMinStartDate()}
                                            required
                                        />
                                    </div>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Thời gian kết thúc*</label>
                                        <input
                                            type="datetime-local"
                                            className="admin-form-input"
                                            name="end_date"
                                            value={formData.end_date}
                                            onChange={handleInputChange}
                                            min={formData.start_date || getMinStartDate()}
                                            required
                                        />
                                    </div>
                                </div>

                                <div className="admin-form-group">
                                    <label className="admin-form-label">Tỷ lệ giảm giá (%)*</label>
                                    <input
                                        type="number"
                                        className="admin-form-input"
                                        name="discount_percent"
                                        value={formData.discount_percent}
                                        onChange={handleInputChange}
                                        onFocus={e => e.target.select()}
                                        placeholder="Nhập số từ 1 đến 99"
                                        min="1"
                                        max="99"
                                        required
                                    />
                                </div>

                                {selectedProduct && formData.discount_percent && (
                                    <div style={{ display: 'flex', justifyContent: 'space-between', background: 'rgba(16, 185, 129, 0.08)', padding: '12px', borderRadius: '6px', color: '#065f46', fontSize: '0.9rem', fontWeight: '500' }}>
                                        <span>Giá sau giảm:</span>
                                        <span>
                                            {formatPrice(Math.round(selectedProduct.price * (1 - parseInt(formData.discount_percent) / 100)))}
                                        </span>
                                    </div>
                                )}
                            </div>
                            <div className="admin-modal-footer" style={{ gap: '12px' }}>
                                <button
                                    type="button"
                                    className="admin-btn admin-btn-outline"
                                    onClick={() => setShowModal(false)}
                                    disabled={submitting}
                                >
                                    Hủy
                                </button>
                                <button
                                    type="submit"
                                    className="admin-btn admin-btn-primary"
                                    disabled={submitting}
                                >
                                    {submitting ? <Loader className="spin" size={16} /> : 'Lưu lại'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {/* CONFIRM DELETE MODAL */}
            {showDeleteModal && promotionToDelete && (
                <div className="admin-modal-overlay">
                    <div className="admin-modal" style={{ maxWidth: '400px' }}>
                        <div className="admin-modal-header">
                            <h2>Xác nhận xóa khuyến mãi</h2>
                            <button className="admin-btn" onClick={() => setShowDeleteModal(false)}>×</button>
                        </div>
                        <div className="admin-modal-body">
                            <p>Bạn có chắc muốn xóa khuyến mãi của sản phẩm <strong>{promotionToDelete.product_name}</strong>?</p>
                            <p style={{ fontSize: '0.85rem', color: 'var(--admin-danger)', marginTop: '8px' }}>
                                Thao tác này sẽ phục hồi sản phẩm về giá gốc ngay khi chiến dịch kết thúc hoặc bị xóa!
                            </p>
                        </div>
                        <div className="admin-modal-footer" style={{ gap: '12px' }}>
                            <button
                                className="admin-btn admin-btn-outline"
                                onClick={() => setShowDeleteModal(false)}
                                disabled={deleting}
                            >
                                Hủy
                            </button>
                            <button
                                className="admin-btn admin-btn-primary"
                                style={{ background: 'var(--admin-danger)' }}
                                onClick={confirmDelete}
                                disabled={deleting}
                            >
                                {deleting ? <Loader className="spin" size={16} /> : 'Xác nhận xóa'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default PromotionManager;
