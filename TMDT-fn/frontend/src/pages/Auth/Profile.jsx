import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { User, Mail, Phone, Calendar, Users, Loader, AlertCircle, UserCircle, Shield, CheckCircle2, AlertTriangle, Lock, Eye, EyeOff } from 'lucide-react';
import OTPModal from '../../components/Auth/OTPModal';
import '../../styles/Auth.css';
import { handlePhoneChange, validatePhone, validateEmail, handleEmailChange, validateName, handleNameChange } from '../../func/phoneValidation';

const Profile = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isSaving, setIsSaving] = useState(false);
    const [error, setError] = useState('');
    const [toast, setToast] = useState({ show: false, message: '', type: 'success' });
    const [nameError, setNameError] = useState('');
    const [phoneError, setPhoneError] = useState('');
    const [emailError, setEmailError] = useState('');

    const [showChangePasswordModal, setShowChangePasswordModal] = useState(false);
    const [showOtpModal, setShowOtpModal] = useState(false);
    const [newPassword, setNewPassword] = useState('');
    const [confirmNewPassword, setConfirmNewPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [newPasswordError, setNewPasswordError] = useState('');
    const [confirmPasswordError, setConfirmPasswordError] = useState('');

    const showToast = (message, type = 'success') => {
        setToast({ show: true, message, type });
        setTimeout(() => setToast({ show: false, message: '', type: 'success' }), 5000);
    };

    const handlePasswordContinue = async () => {
        let hasError = false;
        setNewPasswordError('');
        setConfirmPasswordError('');

        if (!newPassword) {
            setNewPasswordError('Vui lòng nhập mật khẩu mới');
            hasError = true;
        } else if (newPassword.length < 6) {
            setNewPasswordError('Mật khẩu phải có ít nhất 6 ký tự');
            hasError = true;
        }

        if (!confirmNewPassword) {
            setConfirmPasswordError('Vui lòng nhập lại mật khẩu');
            hasError = true;
        } else if (newPassword !== confirmNewPassword) {
            setConfirmPasswordError('Mật khẩu xác nhận không khớp');
            hasError = true;
        }

        if (hasError) return;

        setIsSaving(true);
        try {
            const displayUser = profileData || user;
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/check-password-validity`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: displayUser.email, newPassword })
            });

            const data = await response.json();

            if (!response.ok) {
                setNewPasswordError(data.message || 'Mật khẩu không hợp lệ');
                return;
            }

            if (!data.isValid) {
                setNewPasswordError(data.message || 'Mật khẩu mới không được trùng với mật khẩu hiện tại!');
                return;
            }

            setShowChangePasswordModal(false);
            setShowOtpModal(true);
        } catch (err) {
            showToast('Lỗi khi kiểm tra mật khẩu', 'error');
        } finally {
            setIsSaving(false);
        }
    };

    const handlePasswordChangeConfirmed = async (verifyData) => {
        setShowOtpModal(false);
        try {
            const displayUser = profileData || user;
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/reset-password`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email: displayUser.email, resetToken: verifyData?.resetToken, newPassword })
            });
            const data = await response.json();
            if (response.ok) {
                showToast('Đổi mật khẩu thành công!');
                setNewPassword('');
                setConfirmNewPassword('');
            } else {
                showToast(data.message || 'Lỗi khi đổi mật khẩu', 'error');
            }
        } catch (err) {
            showToast('Lỗi kết nối máy chủ', 'error');
        }
    };

    useEffect(() => {
        const fetchProfile = async () => {
            if (!user?.id) {
                setLoading(false);
                return;
            }

            try {
                const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/information/${user.id}`);
                const result = await response.json();

                if (!response.ok) {
                    throw new Error(result.message || 'Không thể tải thông tin hồ sơ');
                }

                setProfileData(result);
            } catch (err) {
                console.error('Error fetching profile:', err);
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchProfile();
    }, [user?.id]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        if (name === 'num_phone') {
            const { cleaned, error } = handlePhoneChange(value);
            setProfileData(prev => ({ ...(prev || user), num_phone: cleaned }));
            setPhoneError(error);
        } else if (name === 'email') {
            const { cleaned, error } = handleEmailChange(value);
            setProfileData(prev => ({ ...(prev || user), email: cleaned }));
            setEmailError(error);
        } else if (name === 'full_name') {
            const { cleaned, error } = handleNameChange(value);
            setProfileData(prev => ({ ...(prev || user), full_name: cleaned }));
            setNameError(error);
        } else {
            setProfileData(prev => ({ ...(prev || user), [name]: value }));
        }
    };

    const handleSave = async () => {
        const nameErr = validateName((profileData || user)?.full_name || '');
        if (nameErr) { setNameError(nameErr); return; }

        const phoneErr = validatePhone((profileData || user)?.num_phone || '');
        if (phoneErr) { setPhoneError(phoneErr); return; }



        setIsSaving(true);
        try {
            const displayUser = profileData || user;
            const payload = {
                fullName: displayUser.full_name,
                numPhone: displayUser.num_phone,
                dob: displayUser.dob,
                gender: displayUser.gender,
                email: displayUser.email
            };

            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/information/${user.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload),
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.message || 'Cập nhật thất bại');
            }

            showToast('Cập nhật thông tin thành công!');
        } catch (err) {
            console.error('Error updating profile:', err);
            showToast(err.message || 'Lỗi khi cập nhật thông tin', 'error');
        } finally {
            setIsSaving(false);
        }
    };

    const formatDateForInput = (dateString) => {
        if (!dateString) return '';
        try {
            const date = new Date(dateString);
            return date.toISOString().split('T')[0];
        } catch (e) {
            return '';
        }
    };

    if (!user) return null;

    if (loading) {
        return (
            <div className="container" style={{ padding: '120px 15px', textAlign: 'center' }}>
                <Loader className="spin" size={40} style={{ color: 'var(--color-primary)', marginBottom: '15px' }} />
                <p>Đang tải thông tin cá nhân...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="container" style={{ padding: '120px 15px', textAlign: 'center' }}>
                <AlertCircle size={40} style={{ color: '#ef4444', marginBottom: '15px' }} />
                <p style={{ color: '#ef4444' }}>{error}</p>
                <button
                    className="btn-primary"
                    style={{ marginTop: '20px', width: 'auto', padding: '10px 25px' }}
                    onClick={() => window.location.reload()}
                >
                    Thử lại
                </button>
            </div>
        );
    }

    const displayUser = profileData || user;

    return (
        <div className="container" style={{ padding: '100px 15px' }}>
            <div className="auth-container" style={{ background: 'none', padding: 0, minHeight: 'auto' }}>
                <div className="auth-form" style={{ maxWidth: '700px', textAlign: 'left' }}>
                    <h2 className="auth-title">Thông tin tài khoản</h2>

                    <div className="profile-section" style={{ marginBottom: '30px', display: 'flex', alignItems: 'center', gap: '25px' }}>
                        <div className="user-avatar" style={{ width: '90px', height: '90px', fontSize: '36px', background: 'var(--color-primary)', color: '#000' }}>
                            {(displayUser.full_name || displayUser.username || 'U').charAt(0).toUpperCase()}
                        </div>
                        <div>
                            <h3 style={{ fontSize: '26px', margin: 0, fontWeight: '700' }}>{displayUser.full_name || displayUser.username}</h3>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: '#666', marginTop: '8px' }}>
                                <Calendar size={14} />
                                <span>Thành viên từ {displayUser.joined_date ? new Date(displayUser.joined_date).toLocaleDateString('vi-VN') : '---'}</span>
                            </div>
                        </div>
                    </div>

                    <div className="profile-info-grid" style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
                        <div className="info-item">
                            <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <User size={16} /> Họ và Tên
                            </label>
                            <input
                                type="text"
                                name="full_name"
                                className={`form-input ${nameError ? 'input-error' : ''}`}
                                value={displayUser.full_name || ''}
                                onChange={handleChange}
                                disabled={isSaving}
                            />
                            {nameError && (
                                <div className="field-error-msg">
                                    <AlertCircle size={13} /> {nameError}
                                </div>
                            )}
                        </div>

                        {/* <div className="info-item">
                            <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <UserCircle size={16} /> Tên tài khoản
                            </label>
                            <input type="text" className="form-input" value={displayUser.username || ''} readOnly style={{ background: '#f5f5f5', cursor: 'not-allowed' }} />
                        </div> */}

                        <div className="info-item">
                            <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Mail size={16} /> Email
                            </label>
                            <input
                                type="text"
                                name="email"
                                className="form-input"
                                value={displayUser.email || ''}
                                readOnly
                                style={{ background: '#f5f5f5', cursor: 'not-allowed', color: '#666' }}
                            />
                        </div>

                        <div className="info-item">
                            <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Phone size={16} /> Số điện thoại
                            </label>
                            <input
                                type="tel"
                                name="num_phone"
                                className={`form-input ${phoneError ? 'input-error' : ''}`}
                                value={displayUser.num_phone || ''}
                                onChange={handleChange}
                                placeholder="Nhập số điện thoại"
                                disabled={isSaving}
                            />
                            {phoneError && (
                                <div className="field-error-msg">
                                    <AlertCircle size={13} /> {phoneError}
                                </div>
                            )}
                        </div>

                        <div className="info-item">
                            <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Calendar size={16} /> Ngày sinh
                            </label>
                            <input
                                type="date"
                                name="dob"
                                className="form-input"
                                value={formatDateForInput(displayUser.dob)}
                                onChange={handleChange}
                                disabled={isSaving}
                            />
                        </div>

                        <div className="info-item">
                            <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Users size={16} /> Giới tính
                            </label>
                            <select
                                name="gender"
                                className="form-input"
                                value={displayUser.gender || ''}
                                onChange={handleChange}
                                disabled={isSaving}
                            >
                                <option value="">Chọn giới tính</option>
                                <option value="Nam">Nam</option>
                                <option value="Nữ">Nữ</option>
                                <option value="Không rõ">Không rõ</option>
                            </select>
                        </div>
                    </div>

                    <div style={{ display: 'flex', gap: '15px' }}>
                        <button
                            className="btn-primary"
                            style={{ marginTop: '30px' }}
                            onClick={handleSave}
                            disabled={isSaving}
                        >
                            {isSaving ? 'Đang lưu...' : 'Lưu thay đổi'}
                        </button>
                        <button
                            className="btn-outline"
                            style={{ marginTop: '30px', border: '1px solid #ddd', background: 'none' }}
                            onClick={() => setShowChangePasswordModal(true)}
                            disabled={isSaving}
                        >
                            Đổi mật khẩu
                        </button>
                    </div>
                </div>
            </div>

            {toast.show && (
                <div className="toast-container">
                    <div className={`toast ${toast.type}`}>
                        <div className="toast-icon">
                            {toast.type === 'success' ? <CheckCircle2 size={24} /> : <AlertTriangle size={24} />}
                        </div>
                        <div className="toast-content">
                            {toast.message}
                        </div>
                    </div>
                </div>
            )}

            {/* Change Password Modal */}
            {showChangePasswordModal && (
                <div className="admin-modal-overlay nav-modal-overlay" style={{ zIndex: 9999, display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: 'rgba(0,0,0,0.6)', position: 'fixed', top: 0, left: 0, right: 0, bottom: 0 }}>
                    <div className="admin-modal-content" style={{ background: '#fff', borderRadius: '12px', padding: '24px', width: '90%', maxWidth: '400px' }}>
                        <h3 style={{ margin: '0 0 16px', fontSize: '1.25rem' }}>Đổi mật khẩu</h3>
                        <div className="form-group">
                            <label className="form-label">Mật khẩu mới</label>
                            <div className="password-input-wrapper" style={{ position: 'relative' }}>
                                <input
                                    type={showPassword ? 'text' : 'password'}
                                    className={`form-input ${newPasswordError ? 'input-error' : ''}`}
                                    value={newPassword}
                                    onChange={(e) => {
                                        setNewPassword(e.target.value);
                                        if (newPasswordError) setNewPasswordError('');
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            handlePasswordContinue();
                                        }
                                    }}
                                    placeholder="Nhập ít nhất 6 ký tự"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    style={{ position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', cursor: 'pointer', color: '#666' }}
                                >
                                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                                </button>
                            </div>
                            {newPasswordError && (
                                <div className="field-error-msg">
                                    <AlertCircle size={13} /> {newPasswordError}
                                </div>
                            )}
                        </div>
                        <div className="form-group">
                            <label className="form-label">Xác nhận mật khẩu mới</label>
                            <div className="password-input-wrapper" style={{ position: 'relative' }}>
                                <input
                                    type={showConfirmPassword ? 'text' : 'password'}
                                    className={`form-input ${(confirmPasswordError || (confirmNewPassword && !newPassword.startsWith(confirmNewPassword))) ? 'input-error' : ''}`}
                                    value={confirmNewPassword}
                                    onChange={(e) => {
                                        setConfirmNewPassword(e.target.value);
                                        if (confirmPasswordError) setConfirmPasswordError('');
                                    }}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            handlePasswordContinue();
                                        }
                                    }}
                                    placeholder="Nhập lại mật khẩu mới"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                    style={{ position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', cursor: 'pointer', color: '#666' }}
                                >
                                    {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                                </button>
                            </div>
                            {confirmPasswordError && (
                                <div className="field-error-msg">
                                    <AlertCircle size={13} /> {confirmPasswordError}
                                </div>
                            )}
                            {confirmNewPassword && !confirmPasswordError && !newPassword.startsWith(confirmNewPassword) && (
                                <div className="field-error-msg" style={{ fontSize: '11px' }}>
                                    Mật khẩu xác nhận không khớp (sai ký tự)
                                </div>
                            )}
                        </div>
                        <div style={{ display: 'flex', gap: '12px', marginTop: '20px', alignItems: 'center' }}>
                            <button
                                className="btn-outline"
                                style={{ flex: 1, padding: '10px 0', height: '42px', display: 'flex', alignItems: 'center', justifyContent: 'center', marginTop: 0 }}
                                onClick={() => {
                                    setShowChangePasswordModal(false);
                                    setNewPassword('');
                                    setConfirmNewPassword('');
                                    setNewPasswordError('');
                                    setConfirmPasswordError('');
                                }}
                            >
                                Hủy
                            </button>
                            <button
                                className="btn-primary"
                                style={{ flex: 1, padding: '10px 0', height: '42px', display: 'flex', alignItems: 'center', justifyContent: 'center', marginTop: 0 }}
                                onClick={handlePasswordContinue}
                            >
                                Tiếp tục
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* OTP Modal */}
            <OTPModal
                isOpen={showOtpModal}
                onClose={() => setShowOtpModal(false)}
                email={displayUser.email}
                onSuccess={handlePasswordChangeConfirmed}
                actionLabel="Xác nhận đổi mật khẩu"
            />
        </div>
    );
};

export default Profile;
