import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { CheckCircle2, AlertTriangle, AlertCircle, Eye, EyeOff } from 'lucide-react';
import AuthForm from '../../components/Auth/AuthForm';
import OTPModal from '../../components/Auth/OTPModal';
import '../../styles/Auth.css';
import { handlePhoneChange, validatePhone, validateEmail, handleEmailChange, validateName, handleNameChange } from '../../func/phoneValidation';

const Register = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        full_name: '',
        email: '',
        num_phone: '',
        password: '',
        confirmPassword: ''
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const [toast, setToast] = useState({ show: false, message: '', type: 'success' });
    const [nameError, setNameError] = useState('');
    const [phoneError, setPhoneError] = useState('');
    const [emailError, setEmailError] = useState('');
    const [showOtpModal, setShowOtpModal] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const showToast = (message, type = 'success') => {
        setToast({ show: true, message, type });
        setTimeout(() => setToast({ show: false, message: '', type: 'success' }), 10000);
    };

    const handleChange = (e) => {
        if (e.target.name === 'num_phone') {
            const { cleaned, error } = handlePhoneChange(e.target.value);
            setFormData({ ...formData, num_phone: cleaned });
            setPhoneError(error);
        } else if (e.target.name === 'email') {
            const { cleaned, error } = handleEmailChange(e.target.value);
            setFormData({ ...formData, email: cleaned });
            setEmailError(error);
        } else if (e.target.name === 'full_name') {
            const { cleaned, error } = handleNameChange(e.target.value);
            setFormData({ ...formData, full_name: cleaned });
            setNameError(error);
        } else {
            setFormData({ ...formData, [e.target.name]: e.target.value });
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        const { full_name, email, num_phone, password, confirmPassword } = formData;

        if (!full_name || !email || !num_phone || !password || !confirmPassword) {
            setError('Vui lòng điền đầy đủ các trường');
            return;
        }

        if (password !== confirmPassword) {
            setError('Mật khẩu xác nhận không khớp');
            return;
        }

        const nameErr = validateName(full_name);
        if (nameErr) {
            setNameError(nameErr);
            return;
        }

        const emailErr = validateEmail(email);
        if (emailErr) {
            setEmailError(emailErr);
            return;
        }

        const phoneErr = validatePhone(num_phone);
        if (phoneErr) {
            setPhoneError(phoneErr);
            return;
        }

        setLoading(true);
        try {
            const checkUrl = `${import.meta.env.VITE_SERVER_API}/api/user/check-email`;
            const checkRes = await fetch(checkUrl, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email })
            });

            const checkData = await checkRes.json();
            if (checkData.exists) {
                const errorMsg = checkData.is_lock 
                    ? checkData.message 
                    : 'Email này đã được đăng ký. Vui lòng đăng nhập!';
                setError(errorMsg);
                showToast(errorMsg, 'error');
                setLoading(false);
                return;
            }
        } catch (err) {
            console.log('Bỏ qua check-email do lỗi kết nối hoặc API chưa hỗ trợ', err);
        }
        setLoading(false);

        setShowOtpModal(true);
    };

    const handleRegisterConfirmed = async () => {
        setLoading(true);
        try {
            const { full_name, email, num_phone, password } = formData;
            const payload = {
                fullName: full_name,
                email,
                numPhone: num_phone,
                password
            };

            const apiUrl = `${import.meta.env.VITE_SERVER_API}/api/user/register`;
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload),
            });

            const result = await response.json();

            if (!response.ok) {
                setError(result.message || 'Đăng ký thất bại!');
                showToast(result.message || 'Đăng ký thất bại!', 'error');
                setLoading(false);
                setShowOtpModal(false);
                return;
            }

            setShowOtpModal(false);
            showToast('Đăng ký thành công! Vui lòng đăng nhập.');
            setTimeout(() => {
                navigate('/login');
            }, 2000);
        } catch (err) {
            console.error('Lỗi Register:', err);
            setError('Lỗi kết nối đến máy chủ!');
            showToast('Lỗi kết nối đến máy chủ!', 'error');
            setShowOtpModal(false);
        } finally {
            setLoading(false);
        }
    };

    return (
        <AuthForm title="Đăng ký tài khoản" onSubmit={handleSubmit}>
            {error && <div className="error-message">{error}</div>}

            <div className="form-group">
                <label className="form-label">Họ và Tên</label>
                <input
                    type="text"
                    name="full_name"
                    className="form-input"
                    style={{ borderColor: nameError ? '#ef4444' : '' }}
                    value={formData.full_name}
                    onChange={handleChange}
                    placeholder="Nhập họ tên của bạn"
                    disabled={loading}
                />
                {nameError && (
                    <div style={{ fontSize: '0.78rem', color: '#ef4444', marginTop: '5px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <AlertCircle size={13} /> {nameError}
                    </div>
                )}
            </div>

            <div className="form-group">
                <label className="form-label">Email</label>
                <input
                    type="email"
                    name="email"
                    className="form-input"
                    style={{ borderColor: emailError ? '#ef4444' : '' }}
                    value={formData.email}
                    onChange={handleChange}
                    placeholder="Nhập địa chỉ email"
                    disabled={loading}
                />
                {emailError && (
                    <div style={{ fontSize: '0.78rem', color: '#ef4444', marginTop: '5px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <AlertCircle size={13} /> {emailError}
                    </div>
                )}
            </div>

            <div className="form-group">
                <label className="form-label">Số điện thoại</label>
                <input
                    type="tel"
                    name="num_phone"
                    className="form-input"
                    style={{ borderColor: phoneError ? '#ef4444' : '' }}
                    value={formData.num_phone}
                    onChange={handleChange}
                    placeholder="Nhập số điện thoại"
                    disabled={loading}
                />
                {phoneError && (
                    <div style={{ fontSize: '0.78rem', color: '#ef4444', marginTop: '5px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <AlertCircle size={13} /> {phoneError}
                    </div>
                )}
            </div>

            <div className="form-group">
                <label className="form-label">Mật khẩu</label>
                <div className="password-input-wrapper" style={{ position: 'relative' }}>
                    <input
                        type={showPassword ? 'text' : 'password'}
                        name="password"
                        className="form-input"
                        value={formData.password}
                        onChange={handleChange}
                        placeholder="Nhập mật khẩu"
                        disabled={loading}
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
                <label className="form-label">Xác nhận mật khẩu</label>
                <div className="password-input-wrapper" style={{ position: 'relative' }}>
                    <input
                        type={showPassword ? 'text' : 'password'}
                        name="confirmPassword"
                        className="form-input"
                        value={formData.confirmPassword}
                        onChange={handleChange}
                        placeholder="Nhập lại mật khẩu"
                        disabled={loading}
                        style={{
                            borderColor: (formData.confirmPassword && !formData.password.startsWith(formData.confirmPassword)) ? 'red' : ''
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
                {formData.confirmPassword && !formData.password.startsWith(formData.confirmPassword) && (
                    <span style={{ color: 'red', fontSize: '12px', marginTop: '4px', display: 'block' }}>Mật khẩu xác nhận không khớp (sai ký tự)</span>
                )}
            </div>

            <button type="submit" className="btn-primary" disabled={loading}>
                {loading ? 'Đang xử lý...' : 'Đăng ký'}
            </button>

            <div className="auth-links">
                <span>Bạn đã có tài khoản?</span>
                <Link to="/login" className="auth-link">Đăng nhập</Link>
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

            <OTPModal
                isOpen={showOtpModal}
                onClose={() => setShowOtpModal(false)}
                email={formData.email}
                onSuccess={handleRegisterConfirmed}
                actionLabel="Hoàn tất đăng ký"
            />
        </AuthForm>
    );
};

export default Register;
