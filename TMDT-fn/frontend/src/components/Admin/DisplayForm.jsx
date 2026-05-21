import React, { useState } from 'react';
import { Loader } from 'lucide-react';
import toast from 'react-hot-toast';

const API = import.meta.env.VITE_SERVER_API;

const DisplayForm = ({ displayInfo, setIsEditing, onSaveSuccess }) => {
    const initialColumnsHome = displayInfo?.column_in_home || 3;
    const initialItemsHome = displayInfo?.item_in_home || 9;
    const initialColumnsList = displayInfo?.column_in_productlist || 4;
    const initialItemsList = displayInfo?.item_in_productlist || 12;

    const [editForm, setEditForm] = useState({
        column_in_home: initialColumnsHome,
        row_in_home: Math.ceil(initialItemsHome / initialColumnsHome) || 3,
        column_in_productlist: initialColumnsList,
        row_in_productlist: Math.ceil(initialItemsList / initialColumnsList) || 3,
    });

    const [isSaving, setIsSaving] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (isSaving) return;

        setIsSaving(true);
        const submitData = {
            column_in_home: editForm.column_in_home,
            item_in_home: editForm.column_in_home * editForm.row_in_home,
            column_in_productlist: editForm.column_in_productlist,
            item_in_productlist: editForm.column_in_productlist * editForm.row_in_productlist
        };

        try {
            const res = await fetch(`${API}/api/pagination`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(submitData),
            });
            const data = await res.json();

            if (!res.ok) throw new Error(data.message || 'Lưu thất bại');

            toast.success('Cập nhật thông tin trưng bày thành công!');
            if (onSaveSuccess) onSaveSuccess(data?.data?.[0] ?? submitData);
            setIsEditing(false);
        } catch (err) {
            toast.error(err.message || 'Lỗi hệ thống khi lưu cấu hình.');
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="display-form shopee-style-form-wrapper">
            <div className="display-grid">
                <div className="display-section">
                    <div className="section-title">
                        <h4>Trang chủ</h4>
                    </div>
                    <div className="section-content">
                        <div className="form-group">
                            <label>Số cột</label>
                            <input
                                type="number"
                                className="form-input"
                                value={editForm.column_in_home}
                                onChange={e => setEditForm({ ...editForm, column_in_home: parseInt(e.target.value) || 1 })}
                                onFocus={e => e.target.select()}
                                min="1"
                                max="12"
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Số hàng</label>
                            <input
                                type="number"
                                className="form-input"
                                value={editForm.row_in_home}
                                onChange={e => setEditForm({ ...editForm, row_in_home: parseInt(e.target.value) || 1 })}
                                onFocus={e => e.target.select()}
                                min="1"
                                max="20"
                                required
                            />
                        </div>
                        <div style={{ marginTop: '8px', fontSize: '0.85rem', color: 'var(--admin-text-muted)' }}>
                            Tổng số sản phẩm hiển thị: <strong>{editForm.column_in_home * editForm.row_in_home}</strong>
                        </div>
                    </div>
                </div>

                <div className="display-section">
                    <div className="section-title">
                        <h4>Trang danh sách sản phẩm</h4>
                    </div>
                    <div className="section-content">
                        <div className="form-group">
                            <label>Số cột</label>
                            <input
                                type="number"
                                className="form-input"
                                value={editForm.column_in_productlist}
                                onChange={e => setEditForm({ ...editForm, column_in_productlist: parseInt(e.target.value) || 1 })}
                                onFocus={e => e.target.select()}
                                min="1"
                                max="12"
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Số hàng</label>
                            <input
                                type="number"
                                className="form-input"
                                value={editForm.row_in_productlist}
                                onChange={e => setEditForm({ ...editForm, row_in_productlist: parseInt(e.target.value) || 1 })}
                                onFocus={e => e.target.select()}
                                min="1"
                                max="20"
                                required
                            />
                        </div>
                        <div style={{ marginTop: '8px', fontSize: '0.85rem', color: 'var(--admin-text-muted)' }}>
                            Tổng số sản phẩm hiển thị: <strong>{editForm.column_in_productlist * editForm.row_in_productlist}</strong>
                        </div>
                    </div>
                </div>
            </div>

            <div className="form-footer">
                <button type="button" className="btn-cancel" onClick={() => setIsEditing(false)}>
                    Hủy bỏ
                </button>
                <button type="submit" className="btn-save" disabled={isSaving}>
                    {isSaving ? <><Loader className="spin" size={16} /> Đang lưu...</> : 'Lưu Thay Đổi'}
                </button>
            </div>
        </form>
    );
};

export default DisplayForm;
