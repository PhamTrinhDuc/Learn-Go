import React, { useState, useEffect } from 'react';
import { Search, Plus, Edit, Trash2, CheckCircle, ShieldAlert, Loader, Ticket, Calendar, DollarSign, Percent, AlertCircle } from 'lucide-react';
import toast from 'react-hot-toast';

const VoucherManager = () => {
    const [vouchers, setVouchers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState('Tất cả'); // 'Tất cả', 'Chưa diễn ra', 'Đang chạy', 'Đã hết hạn'
    
    // Modal & Form state
    const [showModal, setShowModal] = useState(false);
    const [modalMode, setModalMode] = useState('create'); // 'create' or 'edit'
    const [editingVoucherId, setEditingVoucherId] = useState(null);
    const [formData, setFormData] = useState({
        code: '',
        name: '',
        start_date: '',
        end_date: '',
        discount_type: 'amount', // 'amount' or 'percent'
        discount_value: '',
        discount_target: 'product', // 'product' or 'shipping'
        min_order_value: '',
        max_discount_amount: ''
    });

    const [formError, setFormError] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // Delete state
    const [showDeleteModal, setShowDeleteModal] = useState(false);
    const [voucherToDelete, setVoucherToDelete] = useState(null);
    const [deleting, setDeleting] = useState(false);

    const API = import.meta.env.VITE_SERVER_API;

    const fetchVouchers = async () => {
        setLoading(true);
        try {
            const res = await fetch(`${API}/api/voucher`);
            const data = await res.json();
            if (data.success) {
                setVouchers(data.data);
            } else {
                toast.error(data.message || 'Lỗi khi lấy danh sách mã giảm giá');
            }
        } catch (err) {
            console.error('Error fetching vouchers:', err);
            toast.error('Lỗi kết nối máy chủ');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchVouchers();
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

    const getVoucherStatus = (voucher) => {
        const now = new Date();
        const start = new Date(voucher.start_date);
        const end = new Date(voucher.end_date);

        if (now < start) {
            return { label: 'Chưa diễn ra', color: 'admin-badge-warning', type: 'future' };
        } else if (now >= start && now <= end) {
            return { label: 'Đang hoạt động', color: 'admin-badge-success', type: 'active' };
        } else {
            return { label: 'Đã hết hạn', color: 'admin-badge-danger', type: 'expired' };
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

        if (name === 'discount_value' && value) {
            const num = Number(value);
            if (formData.discount_type === 'percent' && num >= 100) {
                toast.error('Tỷ lệ giảm giá (%) phải nhỏ hơn 100%!');
                setFormError('Tỷ lệ giảm giá (%) phải nhỏ hơn 100%!');
                setFormData(prev => ({ ...prev, discount_value: '99' }));
                return;
            } else if (num <= 0) {
                toast.error('Mức giảm giá phải lớn hơn 0!');
                setFormError('Mức giảm giá phải lớn hơn 0!');
                setFormData(prev => ({ ...prev, discount_value: '1' }));
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

    const openCreateModal = () => {
        setModalMode('create');
        setEditingVoucherId(null);
        setFormData({
            code: '',
            name: '',
            start_date: '',
            end_date: '',
            discount_type: 'amount',
            discount_value: '',
            discount_target: 'product',
            min_order_value: '0',
            max_discount_amount: ''
        });
        setFormError('');
        setShowModal(true);
    };

    const openEditModal = (voucher) => {
        setModalMode('edit');
        setEditingVoucherId(voucher.voucher_id);
        
        // Format ISO dates to datetime-local input format (YYYY-MM-DDTHH:MM)
        const startIso = new Date(voucher.start_date).toISOString().slice(0, 16);
        const endIso = new Date(voucher.end_date).toISOString().slice(0, 16);

        setFormData({
            code: voucher.code,
            name: voucher.name,
            start_date: startIso,
            end_date: endIso,
            discount_type: voucher.discount_type,
            discount_value: voucher.discount_value.toString(),
            discount_target: voucher.discount_target,
            min_order_value: voucher.min_order_value.toString(),
            max_discount_amount: voucher.max_discount_amount ? voucher.max_discount_amount.toString() : ''
        });
        setFormError('');
        setShowModal(true);
    };

    const validateForm = () => {
        const { code, name, start_date, end_date, discount_value, min_order_value } = formData;
        if (!code.trim() || !name.trim() || !start_date || !end_date || !discount_value) {
            return 'Vui lòng nhập đầy đủ các thông tin bắt buộc!';
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

        if (Number(discount_value) <= 0) {
            return 'Mức giảm giá phải lớn hơn 0!';
        }

        if (formData.discount_type === 'percent' && Number(discount_value) >= 100) {
            return 'Tỷ lệ giảm giá (%) phải nhỏ hơn 100%!';
        }

        if (Number(min_order_value) < 0) {
            return 'Giá trị tối thiểu đơn hàng không thể âm!';
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
            discount_value: Number(formData.discount_value),
            min_order_value: Number(formData.min_order_value || 0),
            max_discount_amount: formData.max_discount_amount ? Number(formData.max_discount_amount) : null
        };

        const url = modalMode === 'create' 
            ? `${API}/api/voucher/create`
            : `${API}/api/voucher/edit/${editingVoucherId}`;
        
        const method = modalMode === 'create' ? 'POST' : 'PUT';

        try {
            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            const data = await res.json();

            if (data.success) {
                toast.success(data.message || 'Lưu thành công!');
                setShowModal(false);
                fetchVouchers();
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

    const handleDeleteClick = (voucher) => {
        setVoucherToDelete(voucher);
        setShowDeleteModal(true);
    };

    const confirmDelete = async () => {
        if (!voucherToDelete) return;
        setDeleting(true);
        try {
            const res = await fetch(`${API}/api/voucher/delete/${voucherToDelete.voucher_id}`, {
                method: 'DELETE'
            });
            const data = await res.json();

            if (data.success) {
                toast.success(data.message || 'Xóa voucher thành công!');
                setShowDeleteModal(false);
                setVoucherToDelete(null);
                fetchVouchers();
            } else {
                toast.error(data.message || 'Lỗi khi xóa voucher!');
            }
        } catch (err) {
            console.error(err);
            toast.error('Lỗi kết nối máy chủ!');
        } finally {
            setDeleting(false);
        }
    };

    // Filters and Search logic
    const filteredVouchers = vouchers.filter(v => {
        const matchesSearch = v.code.toLowerCase().includes(searchTerm.toLowerCase()) || 
                              v.name.toLowerCase().includes(searchTerm.toLowerCase());
        
        const status = getVoucherStatus(v);
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

    return (
        <div className="admin-voucher-manager">
            <div className="admin-card">
                <div className="admin-card-header">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        <Ticket size={24} style={{ color: 'var(--admin-primary)' }} />
                        <h2>Quản lý Voucher giảm giá</h2>
                    </div>
                    <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'center' }}>
                        <div className="admin-search-wrapper" style={{ position: 'relative' }}>
                            <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--admin-text-muted)' }} />
                            <input
                                type="text"
                                className="admin-form-input"
                                placeholder="Tìm mã code hoặc tên..."
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
                            <option value="Đang chạy">Đang hoạt động</option>
                            <option value="Đã hết hạn">Đã hết hạn</option>
                        </select>

                        <button
                            className="admin-btn admin-btn-primary"
                            onClick={openCreateModal}
                            style={{ height: '40px', padding: '0 16px', borderRadius: '8px', display: 'flex', alignItems: 'center', gap: '8px' }}
                        >
                            <Plus size={18} />
                            <span>Tạo Voucher mới</span>
                        </button>
                    </div>
                </div>

                <div className="admin-table-container">
                    {loading ? (
                        <div style={{ padding: '40px', textAlign: 'center' }}>
                            <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: '12px' }} />
                            <p>Đang tải dữ liệu voucher...</p>
                        </div>
                    ) : (
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th>Mã Code</th>
                                    <th>Tên Voucher</th>
                                    <th>Thời gian áp dụng</th>
                                    <th>Giảm giá</th>
                                    <th>Đối tượng</th>
                                    <th>Hạn mức tối thiểu</th>
                                    <th>Trạng thái</th>
                                    <th>Thao tác</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredVouchers.map((voucher) => {
                                    const status = getVoucherStatus(voucher);
                                    return (
                                        <tr key={voucher.voucher_id}>
                                            <td style={{ fontWeight: '700', color: 'var(--admin-primary)' }}>
                                                <span style={{ display: 'inline-block', padding: '4px 8px', borderRadius: '6px', background: 'rgba(29, 78, 216, 0.08)' }}>
                                                    {voucher.code}
                                                </span>
                                            </td>
                                            <td>
                                                <div style={{ fontWeight: '500' }}>{voucher.name}</div>
                                            </td>
                                            <td>
                                                <div style={{ fontSize: '0.8rem', color: 'var(--admin-text-muted)', display: 'flex', flexDirection: 'column', gap: '2px' }}>
                                                    <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                                                        <Calendar size={12} /> Bắt đầu: {formatDateLocal(voucher.start_date)}
                                                    </span>
                                                    <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                                                        <Calendar size={12} /> Kết thúc: {formatDateLocal(voucher.end_date)}
                                                    </span>
                                                </div>
                                            </td>
                                            <td style={{ fontWeight: '600' }}>
                                                {voucher.discount_type === 'percent' ? (
                                                    <span style={{ display: 'flex', alignItems: 'center', gap: '2px', color: '#eab308' }}>
                                                        <Percent size={14} /> {voucher.discount_value}%
                                                        {voucher.max_discount_amount && (
                                                            <small style={{ color: 'var(--admin-text-muted)', display: 'block', fontSize: '0.75rem', fontWeight: 'normal' }}>
                                                                (Tối đa {formatPrice(voucher.max_discount_amount)})
                                                            </small>
                                                        )}
                                                    </span>
                                                ) : (
                                                    <span style={{ display: 'flex', alignItems: 'center', gap: '2px', color: '#10b981' }}>
                                                        <DollarSign size={14} /> -{formatPrice(voucher.discount_value)}
                                                    </span>
                                                )}
                                            </td>
                                            <td>
                                                <span className={`admin-badge ${voucher.discount_target === 'product' ? 'admin-badge-success' : 'admin-badge-warning'}`}>
                                                    {voucher.discount_target === 'product' ? 'Sản phẩm' : 'Phí vận chuyển'}
                                                </span>
                                            </td>
                                            <td style={{ color: 'var(--admin-text-muted)' }}>
                                                {voucher.min_order_value > 0 ? formatPrice(voucher.min_order_value) : '0đ'}
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
                                                        onClick={() => openEditModal(voucher)}
                                                    >
                                                        <Edit size={14} />
                                                    </button>
                                                    <button
                                                        className="admin-btn admin-btn-outline admin-btn-sm"
                                                        style={{ color: 'var(--admin-danger)' }}
                                                        title="Xóa"
                                                        onClick={() => handleDeleteClick(voucher)}
                                                    >
                                                        <Trash2 size={14} />
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    );
                                })}
                                {filteredVouchers.length === 0 && (
                                    <tr>
                                        <td colSpan="8" style={{ textAlign: 'center', padding: '32px' }}>
                                            <AlertCircle size={24} style={{ color: 'var(--admin-text-muted)', margin: '0 auto 8px' }} />
                                            <p style={{ color: 'var(--admin-text-muted)' }}>Không tìm thấy mã giảm giá nào phù hợp!</p>
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
                    <div className="admin-modal" style={{ maxWidth: '600px', width: '100%' }}>
                        <div className="admin-modal-header">
                            <h2>{modalMode === 'create' ? 'Tạo mã Voucher mới' : 'Cập nhật Voucher'}</h2>
                            <button className="admin-btn" onClick={() => setShowModal(false)}>×</button>
                        </div>
                        <form onSubmit={handleSubmit} noValidate>
                            <div className="admin-modal-body" style={{ display: 'flex', flexDirection: 'column', gap: '16px', maxHeight: '75vh', overflowY: 'auto' }}>
                                
                                {formError && (
                                    <div style={{ fontSize: '0.85rem', color: 'var(--admin-danger)', background: 'rgba(239, 68, 68, 0.05)', padding: '10px', borderRadius: '6px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                                        <ShieldAlert size={16} /> {formError}
                                    </div>
                                )}

                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Mã Code*</label>
                                        <input
                                            type="text"
                                            className="admin-form-input"
                                            name="code"
                                            value={formData.code}
                                            onChange={(e) => setFormData(p => ({ ...p, code: e.target.value.replace(/\s/g, '').toUpperCase() }))}
                                            placeholder="VD: GIAM20K, SUMMERSALE"
                                            required
                                            disabled={modalMode === 'edit'} // Lock code on edit
                                        />
                                    </div>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Tên Voucher/Chiến dịch*</label>
                                        <input
                                            type="text"
                                            className="admin-form-input"
                                            name="name"
                                            value={formData.name}
                                            onChange={handleInputChange}
                                            placeholder="VD: Khuyến mãi mùa hè 20K"
                                            required
                                        />
                                    </div>
                                </div>

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

                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Mục tiêu áp dụng*</label>
                                        <select
                                            className="admin-form-input"
                                            name="discount_target"
                                            value={formData.discount_target}
                                            onChange={handleInputChange}
                                        >
                                            <option value="product">Giảm giá tiền sản phẩm</option>
                                            <option value="shipping">Giảm phí vận chuyển</option>
                                        </select>
                                    </div>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Loại giảm giá*</label>
                                        <select
                                            className="admin-form-input"
                                            name="discount_type"
                                            value={formData.discount_type}
                                            onChange={handleInputChange}
                                        >
                                            <option value="amount">Theo số tiền cố định (đ)</option>
                                            <option value="percent">Theo phần trăm (%)</option>
                                        </select>
                                    </div>
                                </div>

                                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Mức giảm giá*</label>
                                        <input
                                            type="number"
                                            className="admin-form-input"
                                            name="discount_value"
                                            value={formData.discount_value}
                                            onChange={handleInputChange}
                                            onFocus={e => e.target.select()}
                                            placeholder={formData.discount_type === 'percent' ? 'VD: 15 (%)' : 'VD: 50000 (đ)'}
                                            min="1"
                                            max={formData.discount_type === 'percent' ? '99' : undefined}
                                            required
                                        />
                                    </div>
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Đơn hàng tối thiểu áp dụng (đ)</label>
                                        <input
                                            type="number"
                                            className="admin-form-input"
                                            name="min_order_value"
                                            value={formData.min_order_value}
                                            onChange={handleInputChange}
                                            onFocus={e => e.target.select()}
                                            placeholder="VD: 150000"
                                            min="0"
                                        />
                                    </div>
                                </div>

                                {formData.discount_type === 'percent' && (
                                    <div className="admin-form-group">
                                        <label className="admin-form-label">Số tiền giảm giá tối đa (đ) - Để trống nếu không giới hạn</label>
                                        <input
                                            type="number"
                                            className="admin-form-input"
                                            name="max_discount_amount"
                                            value={formData.max_discount_amount}
                                            onChange={handleInputChange}
                                            onFocus={e => e.target.select()}
                                            placeholder="VD: 50000"
                                            min="1"
                                        />
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
            {showDeleteModal && voucherToDelete && (
                <div className="admin-modal-overlay">
                    <div className="admin-modal" style={{ maxWidth: '400px' }}>
                        <div className="admin-modal-header">
                            <h2>Xác nhận xóa voucher</h2>
                            <button className="admin-btn" onClick={() => setShowDeleteModal(false)}>×</button>
                        </div>
                        <div className="admin-modal-body">
                            <p>Bạn có chắc chắn muốn xóa mã giảm giá <strong>{voucherToDelete.code}</strong> khỏi hệ thống?</p>
                            <p style={{ fontSize: '0.85rem', color: 'var(--admin-danger)', marginTop: '8px' }}>
                                Lưu ý: Thao tác này không thể hoàn tác!
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
                                {deleting ? <Loader className="spin" size={16} /> : 'Đồng ý xóa'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default VoucherManager;
