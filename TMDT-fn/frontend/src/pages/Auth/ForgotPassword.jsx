import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { Eye, EyeOff } from 'lucide-react';
import AuthForm from '../../components/Auth/AuthForm';
import OTPModal from '../../components/Auth/OTPModal';

const API = import.meta.env.VITE_SERVER_API;

const ForgotPassword = () => {
    const navigate = useNavigate();
    const [step, setStep] = useState(1);
    const [email, setEmail] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [resetToken, setResetToken] = useState('');
    const [error, setError] = useState('');
    const [showOtpModal, setShowOtpModal] = useState(false);
    const [isUpdating, setIsUpdating] = useState(false);

    // Toggle password visibility
    const [showPassword, setShowPassword] = useState(false);

    const handleStep1 = async (e) => {
        e.preventDefault();
        if (!email || !email.includes('@')) {
            setError('Vui lòng nhập Email hợp lệ');
            return;
        }
        setError('');

        try {
            const checkUrl = `${API}/api/user/check-email`;
            const checkRes = await fetch(checkUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email })
            });
            const checkData = await checkRes.json();

            if (!checkRes.ok) {
                if (checkRes.status === 404 || checkData.exists === false) {
                    setError(checkData.message || 'Email này chưa được đăng ký trong hệ thống!');
                    return;
                }
                if (checkRes.status === 403 || checkData.is_lock) {
                    setError(checkData.message || 'Tài khoản này hiện đang bị khóa. Vui lòng liên hệ quản trị viên!');
                    return;
                }
                setError(checkData.message || 'Lỗi hệ thống khi kiểm tra email!');
                return;
            }

            setShowOtpModal(true);
        } catch (err) {
            setError('Lỗi kết nối khi kiểm tra email.');
            console.error(err);
        }
    };

    const handleOtpSuccess = (data) => {
        setShowOtpModal(false);
        if (data && data.resetToken) {
            setResetToken(data.resetToken);
        }
        setStep(3);
    };

    const handleStep3 = async (e) => {
        e.preventDefault();
        if (!newPassword || !confirmPassword) {
            setError('Vui lòng nhập đầy đủ mật khẩu mới và xác nhận');
            return;
        }
        if (newPassword.length < 6) {
            setError('Mật khẩu mới phải có ít nhất 6 ký tự');
            return;
        }
        if (newPassword !== confirmPassword) {
            setError('Mật khẩu xác nhận không khớp');
            return;
        }
        if (!resetToken) {
            setError('Lỗi kết nối phiên bảo mật. Vui lòng thử lại!');
            return;
        }
        setError('');
        setIsUpdating(true);

        try {
            const response = await fetch(`${API}/api/user/reset-password`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, resetToken, newPassword })
            });

            const data = await response.json();

            if (response.ok) {
                toast.success('Đổi mật khẩu thành công!', { id: 'auth-toast' });
                navigate('/login');
            } else {
                setError(data.message || 'Lỗi khi cập nhật mật khẩu');
            }
        } catch (err) {
            setError('Lỗi kết nối máy chủ');
            console.error(err);
        } finally {
            setIsUpdating(false);
        }
    };

    const renderStepContent = () => {
        switch (step) {
            case 1:
                return (
                    <>
                        <div className="form-group">
                            <label className="form-label">Nhập Email đã đăng ký</label>
                            <input
                                type="email"
                                className="form-input"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="Địa chỉ email"
                            />
                        </div>
                        <button onClick={handleStep1} className="btn-primary">Gửi mã xác thực</button>
                    </>
                );
            case 3:
                return (
                    <>
                        <div className="form-group">
                            <label className="form-label">Mật khẩu mới</label>
                            <div className="password-input-wrapper" style={{ position: 'relative' }}>
                                <input
                                    type={showPassword ? 'text' : 'password'}
                                    className="form-input"
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    placeholder="Nhập mật khẩu ít nhất 6 ký tự"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    style={{ position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', cursor: 'pointer', color: '#666' }}
                                >
                                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                                </button>
                            </div>
                        </div>
                        <div className="form-group">
                            <label className="form-label">Xác nhận mật khẩu mới</label>
                            <div className="password-input-wrapper" style={{ position: 'relative' }}>
                                <input
                                    type={showPassword ? 'text' : 'password'}
                                    className="form-input"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    placeholder="Nhập lại mật khẩu mới"
                                    style={{
                                        borderColor: (confirmPassword && !newPassword.startsWith(confirmPassword)) ? 'red' : ''
                                    }}
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPassword(!showPassword)}
                                    style={{ position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)', background: 'none', border: 'none', cursor: 'pointer', color: '#666' }}
                                >
                                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                                </button>
                            </div>
                            {confirmPassword && !newPassword.startsWith(confirmPassword) && (
                                <span style={{ color: 'red', fontSize: '12px', marginTop: '4px', display: 'block' }}>Mật khẩu xác nhận không khớp (sai ký tự)</span>
                            )}
                        </div>
                        <button onClick={handleStep3} disabled={isUpdating} className="btn-primary">
                            {isUpdating ? 'Đang cập nhật...' : 'Đổi mật khẩu'}
                        </button>
                    </>
                );
            default:
                return null;
        }
    };

    return (
        <AuthForm title="Quên mật khẩu">
            {error && <div className="error-message">{error}</div>}

            {renderStepContent()}

            <div className="auth-links">
                {(!localStorage.getItem('user') && !localStorage.getItem('role')) && (
                    <Link to="/login" className="auth-link">Quay lại đăng nhập</Link>
                )}
            </div>

            <OTPModal
                isOpen={showOtpModal}
                onClose={() => setShowOtpModal(false)}
                email={email}
                onSuccess={handleOtpSuccess}
                actionLabel="Xác thực"
            />
        </AuthForm>
    );
};

export default ForgotPassword;
