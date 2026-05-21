import React, { useState } from 'react';
import { Plus, Trash2, Save, Info, CheckCircle2 } from 'lucide-react';
import './CategoryManager.css';

const CategoryManager = ({ isOpen, onClose, onSuccess }) => {
    const [categoryName, setCategoryName] = useState('');
    const [categorySlug, setCategorySlug] = useState('');
    const [fields, setFields] = useState([
        { label: '', key: '', type: 'text' }
    ]);
    const [loading, setLoading] = useState(false);
    const [toast, setToast] = useState({ show: false, message: '', type: 'success' });

    if (!isOpen) return null;

    const showToast = (message, type = 'success') => {
        setToast({ show: true, message, type });
        setTimeout(() => setToast({ show: false, message: '', type: 'success' }), 3000);
    };

    const handleAddField = () => {
        setFields([...fields, { label: '', key: '', type: 'text' }]);
    };

    const handleRemoveField = (index) => {
        if (fields.length > 1) {
            const newFields = fields.filter((_, i) => i !== index);
            setFields(newFields);
        } else {
            showToast('Phải có ít nhất một thông số', 'error');
        }
    };

    const handleFieldChange = (index, key, value) => {
        const newFields = [...fields];
        newFields[index][key] = value;
        setFields(newFields);
    };

    const resetForm = () => {
        setCategoryName('');
        setCategorySlug('');
        setFields([{ label: '', key: '', type: 'text' }]);
    };

    const handleSave = async () => {
        if (!categoryName || !categorySlug) {
            showToast('Vui lòng nhập tên và mã danh mục', 'error');
            return;
        }

        const isFieldsValid = fields.every(f => f.label.trim() && f.key.trim());
        if (!isFieldsValid) {
            showToast('Vui lòng điền đầy đủ thông tin các thông số', 'error');
            return;
        }

        setLoading(true);
        try {
            const payload = {
                slug: categorySlug,
                label: categoryName,
                fields: fields
            };

            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/form-templates`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload),
            });

            const result = await response.json();

            if (result.success) {
                showToast('Lưu cấu trúc danh mục thành công!');
                resetForm();
                if (onSuccess) onSuccess();
                setTimeout(() => onClose(), 1500);
            } else {
                showToast(result.message || 'Có lỗi xảy ra khi lưu', 'error');
            }
        } catch (error) {
            console.error('Error saving category:', error);
            showToast('Không thể kết nối với máy chủ', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="category-manager-overlay" onClick={onClose}>
            <div className="category-manager-modal" onClick={(e) => e.stopPropagation()}>
                <div className="category-manager-header-row">
                    <h1>Tạo loại sản phẩm mới</h1>
                    <button className="close-btn" onClick={onClose}>&times;</button>
                </div>
                <p className="subtitle">Định nghĩa cấu trúc thông số kỹ thuật cho các dòng sản phẩm của bạn.</p>

                <div className="category-card shadow-none">
                    <div className="form-section">
                        <h2 className="form-section-title">
                            <Info size={20} style={{ color: '#3b82f6' }} />
                            Thông tin cơ bản
                        </h2>
                        <div className="input-grid">
                            <div className="form-group">
                                <label className="form-label">Tên danh mục</label>
                                <input
                                    type="text"
                                    className="form-input"
                                    placeholder="VD: Đồng hồ thông minh"
                                    value={categoryName}
                                    onChange={(e) => setCategoryName(e.target.value)}
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Mã danh mục (Slug)</label>
                                <input
                                    type="text"
                                    className="form-input"
                                    placeholder="VD: smart-watch"
                                    value={categorySlug}
                                    onChange={(e) => setCategorySlug(e.target.value)}
                                />
                            </div>
                        </div>
                    </div>

                    <div className="form-section">
                        <h2 className="form-section-title">
                            <Plus size={20} style={{ color: '#10b981' }} />
                            Cấu trúc thông số (Fields)
                        </h2>
                        <div className="field-list">
                            {fields.map((field, index) => (
                                <div key={index} className="field-row">
                                    <div className="form-group">
                                        <label className="form-label">Tên thông số</label>
                                        <input
                                            type="text"
                                            className="form-input"
                                            placeholder="VD: Màn hình"
                                            value={field.label}
                                            onChange={(e) => handleFieldChange(index, 'label', e.target.value)}
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label className="form-label">Mã thông số (Key)</label>
                                        <input
                                            type="text"
                                            className="form-input"
                                            placeholder="VD: screen"
                                            value={field.key}
                                            onChange={(e) => handleFieldChange(index, 'key', e.target.value)}
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label className="form-label">Loại dữ liệu</label>
                                        <select
                                            className="form-input"
                                            value={field.type}
                                            onChange={(e) => handleFieldChange(index, 'type', e.target.value)}
                                        >
                                            <option value="text">Văn bản (Single)</option>
                                            <option value="multi">Nhiều giá trị (Tags)</option>
                                        </select>
                                    </div>
                                    <button
                                        className="btn btn-danger"
                                        onClick={() => handleRemoveField(index)}
                                        title="Xóa dòng"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </div>
                            ))}
                        </div>
                        <button className="btn add-field-btn" onClick={handleAddField}>
                            <Plus size={18} />
                            Thêm dòng thông số
                        </button>
                    </div>

                    <div className="form-footer">
                        <button className="btn btn-outline" style={{ marginRight: '12px' }} onClick={onClose}>Hủy</button>
                        <button
                            className="btn btn-primary save-btn"
                            onClick={handleSave}
                            disabled={loading}
                        >
                            {loading ? 'Đang lưu...' : (
                                <>
                                    <Save size={20} />
                                    Lưu cấu trúc
                                </>
                            )}
                        </button>
                    </div>
                </div>

                {toast.show && (
                    <div className={`toast ${toast.type === 'error' ? 'bg-red-500' : ''}`} style={{ backgroundColor: toast.type === 'error' ? '#ef4444' : '#10b981', bottom: '10px' }}>
                        {toast.type === 'success' ? <CheckCircle2 size={18} /> : <Info size={18} />}
                        <span>{toast.message}</span>
                    </div>
                )}
            </div>
        </div>
    );
};

export default CategoryManager;
