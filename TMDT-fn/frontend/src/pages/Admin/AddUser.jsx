import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    ChevronLeft, Save, X, UserCircle, Mail, Phone,
    Shield, Calendar, Users, Lock, KeyRound, Loader, AlertCircle, CheckCircle2
} from 'lucide-react';
import { handlePhoneChange, validatePhone, validateEmail, handleEmailChange, validateName, handleNameChange } from '../../func/phoneValidation';

const AddUser = () => {
    const navigate = useNavigate();
    const [isSaving, setIsSaving] = useState(false);
    const [phoneError, setPhoneError] = useState('');
    const [emailError, setEmailError] = useState('');
    const [nameError, setNameError] = useState('');
    const [toast, setToast] = useState({ show: false, message: '', type: 'success' });

    const [formData, setFormData] = useState({
        password: '',
        full_name: '',
        email: '',
        num_phone: '',
        role: 'customer',
        is_lock: false,
        dob: '',
        gender: ''
    });

    const showToast = (message, type = 'success') => {
        setToast({ show: true, message, type });
        setTimeout(() => setToast({ show: false, message: '', type: 'success' }), 4000);
    };

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        if (name === 'num_phone') {
            const { cleaned, error } = handlePhoneChange(value);
            setFormData(prev => ({ ...prev, num_phone: cleaned }));
            setPhoneError(error);
        } else if (name === 'email') {
            const { cleaned, error } = handleEmailChange(value);
            setFormData(prev => ({ ...prev, email: cleaned }));
            setEmailError(error);
        } else if (name === 'full_name') {
            const { cleaned, error } = handleNameChange(value);
            setFormData(prev => ({ ...prev, full_name: cleaned }));
            setNameError(error);
        } else {
            setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!formData.password || !formData.full_name || !formData.email) {
            showToast('Vui lòng điền đầy đủ các trường bắt buộc!', 'error');
            return;
        }

        const nameErr = validateName(formData.full_name);
        if (nameErr) { setNameError(nameErr); return; }

        if (formData.num_phone) {
            const phoneErr = validatePhone(formData.num_phone);
            if (phoneErr) { setPhoneError(phoneErr); return; }
        }

        const emailErr = validateEmail(formData.email);
        if (emailErr) {
            setEmailError(emailErr);
            return;
        }

        setIsSaving(true);
        try {
            const payload = {
                password: formData.password,
                fullName: formData.full_name,
                email: formData.email,
                numPhone: formData.num_phone || null,
                role: formData.role,
                isLock: formData.is_lock,
                dob: formData.dob || null,
                gender: formData.gender || null
            };

            const res = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            const result = await res.json();

            if (res.ok) {
                showToast('Thêm người dùng thành công!');
                setTimeout(() => navigate('/admin/users'), 1500);
            } else {
                showToast(result.message || 'Lỗi khi thêm người dùng!', 'error');
            }
        } catch (err) {
            console.error('Add user error:', err);
            showToast('Lỗi kết nối server!', 'error');
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="admin-add-user">
            <div style={{ marginBottom: '20px', display: 'flex', alignItems: 'center', gap: '12px' }}>
                <button
                    className="admin-btn admin-btn-outline"
                    onClick={() => navigate('/admin/users')}
                    style={{ padding: '8px' }}
                >
                    <ChevronLeft size={20} />
                </button>
                <h2 style={{ margin: 0 }}>Thêm người dùng mới</h2>
            </div>

            <div className="admin-card" style={{ maxWidth: '860px' }}>
                <form onSubmit={handleSubmit} style={{ padding: '28px' }}>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>

                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <KeyRound size={15} /> Mật khẩu <span style={{ color: '#ef4444' }}>*</span>
                            </label>
                            <input
                                type="password"
                                name="password"
                                className="admin-form-input"
                                placeholder="Mật khẩu..."
                                value={formData.password}
                                onChange={handleChange}
                                required
                                autoComplete="new-password"
                            />
                        </div>
                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <UserCircle size={15} /> Họ và tên <span style={{ color: '#ef4444' }}>*</span>
                            </label>
                            <input
                                type="text"
                                name="full_name"
                                className="admin-form-input"
                                style={{ borderColor: nameError ? '#ef4444' : '' }}
                                placeholder="Nguyễn Văn A"
                                value={formData.full_name}
                                onChange={handleChange}
                                required
                            />
                            {nameError && (
                                <div style={{ fontSize: '0.75rem', color: '#ef4444', marginTop: '4px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                                    <AlertCircle size={12} /> {nameError}
                                </div>
                            )}
                        </div>

                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Mail size={15} /> Email <span style={{ color: '#ef4444' }}>*</span>
                            </label>
                            <input
                                type="email"
                                name="email"
                                className="admin-form-input"
                                style={{ borderColor: emailError ? '#ef4444' : '' }}
                                placeholder="example@email.com"
                                value={formData.email}
                                onChange={handleChange}
                                required
                            />
                            {emailError && (
                                <div style={{ fontSize: '0.75rem', color: '#ef4444', marginTop: '4px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                                    <AlertCircle size={12} /> {emailError}
                                </div>
                            )}
                        </div>

                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Phone size={15} /> Số điện thoại
                            </label>
                            <input
                                type="tel"
                                name="num_phone"
                                className="admin-form-input"
                                style={{ borderColor: phoneError ? '#ef4444' : '' }}
                                placeholder="0xxxxxxxxx"
                                value={formData.num_phone}
                                onChange={handleChange}
                            />
                            {phoneError && (
                                <div style={{ fontSize: '0.75rem', color: '#ef4444', marginTop: '4px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                                    <AlertCircle size={12} /> {phoneError}
                                </div>
                            )}
                        </div>

                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Shield size={15} /> Vai trò
                            </label>
                            <select
                                name="role"
                                className="admin-form-input"
                                value={formData.role}
                                onChange={handleChange}
                            >
                                <option value="customer">Khách hàng</option>
                                <option value="admin">Quản trị viên</option>
                            </select>
                        </div>

                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Calendar size={15} /> Ngày sinh
                            </label>
                            <input
                                type="date"
                                name="dob"
                                className="admin-form-input"
                                value={formData.dob}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="admin-form-group">
                            <label className="admin-form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Users size={15} /> Giới tính
                            </label>
                            <select
                                name="gender"
                                className="admin-form-input"
                                value={formData.gender}
                                onChange={handleChange}
                            >
                                <option value="">-- Chọn giới tính --</option>
                                <option value="Nam">Nam</option>
                                <option value="Nữ">Nữ</option>
                                <option value="Không rõ">Không rõ</option>
                            </select>
                        </div>

                        <div className="admin-form-group" style={{ gridColumn: '1 / -1' }}>
                            <label style={{ display: 'flex', alignItems: 'center', gap: '10px', cursor: 'pointer', userSelect: 'none' }}>
                                <input
                                    type="checkbox"
                                    name="is_lock"
                                    checked={formData.is_lock}
                                    onChange={handleChange}
                                    style={{ width: '16px', height: '16px', accentColor: '#ef4444', cursor: 'pointer' }}
                                />
                                <Lock size={15} style={{ color: formData.is_lock ? '#ef4444' : '#94a3b8' }} />
                                <span style={{ color: formData.is_lock ? '#ef4444' : 'inherit', fontWeight: '500' }}>
                                    Khóa tài khoản ngay khi tạo
                                </span>
                            </label>
                        </div>

                    </div>

                    <div style={{ marginTop: '32px', display: 'flex', gap: '12px', justifyContent: 'flex-end' }}>
                        <button
                            type="button"
                            className="admin-btn admin-btn-outline"
                            onClick={() => navigate('/admin/users')}
                            disabled={isSaving}
                        >
                            <X size={18} /> Hủy
                        </button>
                        <button
                            type="submit"
                            className="admin-btn admin-btn-primary"
                            disabled={isSaving}
                        >
                            {isSaving ? <Loader className="spin" size={18} /> : <Save size={18} />}
                            {isSaving ? 'Đang lưu...' : 'Lưu người dùng'}
                        </button>
                    </div>
                </form>
            </div>

            {toast.show && (
                <div style={{ position: 'fixed', bottom: '24px', left: '50%', transform: 'translateX(-50%)', zIndex: 9999 }}>
                    <div style={{
                        display: 'flex', alignItems: 'center', gap: '12px',
                        padding: '12px 24px', borderRadius: '12px',
                        background: toast.type === 'success' ? '#10b981' : '#ef4444',
                        color: '#fff', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.15)',
                        fontWeight: '500', fontSize: '0.95rem'
                    }}>
                        {toast.type === 'success' ? <CheckCircle2 size={20} /> : <AlertCircle size={20} />}
                        {toast.message}
                    </div>
                </div>
            )}
        </div>
    );
};

export default AddUser;
