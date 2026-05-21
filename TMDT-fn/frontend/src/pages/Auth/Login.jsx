import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import AuthForm from '../../components/Auth/AuthForm';
import { useAuth } from '../../context/AuthContext';
import { GoogleLogin } from '@react-oauth/google';

const Login = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { login } = useAuth();
    const from = location.state?.from?.pathname || "/";

    const [formData, setFormData] = useState({
        username: '',
        password: ''
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        if (!formData.username || !formData.password) {
            setError('Vui lòng nhập đầy đủ thông tin');
            return;
        }
        setLoading(true);
        try {
            const apiUrl = `${import.meta.env.VITE_SERVER_API}/api/user/login`;
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });
            const result = await response.json();
            if (!response.ok) {
                setError(result.message || 'Đăng nhập thất bại!');
                setLoading(false);
                return;
            }
            login(result.user);
            navigate(from, { replace: true });
        } catch (err) {
            console.error('Lỗi Login:', err);
            setError('Lỗi kết nối đến máy chủ!');
        } finally {
            setLoading(false);
        }
    };

    const handleGoogleLogin = async (credential) => {
        setError('');
        setLoading(true);
        try {
            const apiUrl = `${import.meta.env.VITE_SERVER_API}/api/user/google/login`;
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ idToken: credential }),
            });
            const result = await response.json();
            if (!response.ok) {
                setError(result.message || 'Đăng nhập bằng Google thất bại!');
                setLoading(false);
                return;
            }
            login(result.user);
            navigate(from, { replace: true });
        } catch (err) {
            console.error('Lỗi Login Google:', err);
            setError('Lỗi kết nối đến máy chủ!');
        } finally {
            setLoading(false);
        }
    };

    return (
        <AuthForm title="Đăng nhập" onSubmit={handleSubmit}>
            {error && <div className="error-message">{error}</div>}

            <div className="form-group">
                <label className="form-label">Email hoặc Số điện thoại</label>
                <input
                    type="text"
                    name="username"
                    className="form-input"
                    value={formData.username}
                    onChange={handleChange}
                    placeholder="Nhập email hoặc SĐT"
                    disabled={loading}
                />
            </div>

            <div className="form-group">
                <label className="form-label">Mật khẩu</label>
                <input
                    type="password"
                    name="password"
                    className="form-input"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="Nhập mật khẩu"
                    disabled={loading}
                />
            </div>

            <button type="submit" className="btn-primary" disabled={loading} style={{ marginBottom: '25px' }}>
                {loading ? 'Đang đăng nhập...' : 'Đăng nhập'}
            </button>

            <div style={{ display: 'flex', alignItems: 'center', marginBottom: '25px' }}>
                <div style={{ flex: 1, height: '1px', backgroundColor: '#e2e8f0' }}></div>
                <span style={{ padding: '0 15px', color: '#64748b', fontSize: '13px', fontWeight: '600' }}>HOẶC</span>
                <div style={{ flex: 1, height: '1px', backgroundColor: '#e2e8f0' }}></div>
            </div>

            <div style={{ display: 'flex', justifyContent: 'center', width: '100%' }}>
                <GoogleLogin
                    onSuccess={credentialResponse => {
                        // credentialResponse.credential chính là ID Token
                        console.log("Gửi Token này về Backend:", credentialResponse.credential);
                        if (typeof handleGoogleLogin === 'function') {
                            handleGoogleLogin(credentialResponse.credential);
                        }
                    }}
                    onError={() => {
                        console.log('Login Failed');
                    }}
                />
            </div>

            <div className="auth-links">
                <Link to="/forgot-password" title="Quên mật khẩu?" className="auth-link">Quên mật khẩu?</Link>
                <span>|</span>
                <Link to="/register" title="Đăng ký tài khoản" className="auth-link">Đăng ký mới</Link>
            </div>
        </AuthForm>
    );
};

export default Login;
