import React, { useState, useEffect, useRef } from 'react';
import { X, ShieldCheck, Mail, Loader, CheckCircle2, AlertTriangle } from 'lucide-react';
import toast from 'react-hot-toast';

const API = import.meta.env.VITE_SERVER_API;

const OTPModal = ({ isOpen, onClose, email, onSuccess, actionLabel = 'Xác thực', sendOtpEndpoint = '/api/otp/send-otp', verifyOtpEndpoint = '/api/otp/verify-otp' }) => {
    const [otp, setOtp] = useState(['', '', '', '', '', '']);
    const [isSending, setIsSending] = useState(false);
    const [isVerifying, setIsVerifying] = useState(false);
    const [countdown, setCountdown] = useState(60);
    const [canResend, setCanResend] = useState(false);
    const [error, setError] = useState('');
    const inputRefs = useRef([]);

    useEffect(() => {
        if (isOpen && email) {
            handleSendOTP();
            resetOtp();
        }
    }, [isOpen, email]);

    useEffect(() => {
        let timer;
        if (countdown > 0 && !canResend && isOpen) {
            timer = setInterval(() => {
                setCountdown(prev => {
                    if (prev <= 1) {
                        setCanResend(true);
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        }
        return () => clearInterval(timer);
    }, [countdown, canResend, isOpen]);

    const resetOtp = () => {
        setOtp(['', '', '', '', '', '']);
        setError('');
        if (inputRefs.current[0]) {
            setTimeout(() => inputRefs.current[0].focus(), 100);
        }
    };

    const handleSendOTP = async () => {
        setIsSending(true);
        setError('');
        try {
            const response = await fetch(`${API}${sendOtpEndpoint}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email })
            });
            const data = await response.json();

            if (response.ok) {
                toast.success('Mã OTP đã được gửi đến email của bạn!', { id: 'otp-toast' });
                setCountdown(60);
                setCanResend(false);
            } else {
                setError(data.message || 'Lỗi gửi mã OTP.');
                if (data.message?.includes('sẵn sàng') || data.message?.includes('vui lòng thử lại')) {
                    toast.error(data.message, { id: 'otp-toast' });
                }
            }
        } catch (err) {
            setError('Lỗi kết nối máy chủ khi gửi OTP.');
            console.error(err);
        } finally {
            setIsSending(false);
        }
    };

    const handleVerifyOTP = async (otpString = otp.join('')) => {
        if (otpString.length !== 6) {
            setError('Vui lòng nhập đủ 6 số OTP.');
            return;
        }

        setIsVerifying(true);
        setError('');
        try {
            const response = await fetch(`${API}${verifyOtpEndpoint}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, otp: otpString })
            });
            const data = await response.json();

            if (response.ok && data.success) {
                toast.success('Xác thực OTP thành công!', { id: 'otp-toast' });
                onSuccess(data);
            } else {
                setError(data.message || 'Mã OTP không chính xác hoặc đã hết hạn.');
                resetOtp();
            }
        } catch (err) {
            setError('Lỗi kết nối máy chủ khi xác thực OTP.');
            console.error(err);
        } finally {
            setIsVerifying(false);
        }
    };

    const handleChange = (index, value) => {
        if (!/^\d*$/.test(value)) return;

        const newOtp = [...otp];
        newOtp[index] = value;
        setOtp(newOtp);
        setError('');

        if (value && index < 5) {
            inputRefs.current[index + 1].focus();
        }

        if (index === 5 && value && newOtp.every(digit => digit !== '')) {
            handleVerifyOTP(newOtp.join(''));
        }
    };

    const handleKeyDown = (index, e) => {
        if (e.key === 'Backspace' && !otp[index] && index > 0) {
            inputRefs.current[index - 1].focus();
        } else if (e.key === 'Enter') {
            handleVerifyOTP();
        }
    };

    const handlePaste = (e) => {
        e.preventDefault();
        const pastedData = e.clipboardData.getData('text').slice(0, 6).replace(/\D/g, '');
        if (pastedData) {
            const newOtp = [...otp];
            for (let i = 0; i < pastedData.length; i++) {
                newOtp[i] = pastedData[i];
            }
            setOtp(newOtp);
            if (pastedData.length === 6) {
                inputRefs.current[5].focus();
                handleVerifyOTP(newOtp.join(''));
            } else {
                inputRefs.current[pastedData.length].focus();
            }
        }
    };

    if (!isOpen) return null;

    return (
        <div className="admin-modal-overlay nav-modal-overlay" style={{ zIndex: 99999, display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: 'rgba(0,0,0,0.6)', position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, backdropFilter: 'blur(4px)' }}>
            <div className="admin-modal-content" style={{ background: '#fff', borderRadius: '16px', padding: '32px 24px', width: '90%', maxWidth: '400px', position: 'relative', boxShadow: '0 20px 40px rgba(0,0,0,0.2)' }}>
                <button
                    style={{ position: 'absolute', top: '16px', right: '16px', border: 'none', background: 'transparent', cursor: 'pointer', padding: 4 }}
                    onClick={onClose}
                    disabled={isVerifying}
                >
                    <X size={24} color="#666" />
                </button>

                <div style={{ textAlign: 'center', marginBottom: '24px' }}>
                    <div style={{ background: '#e0f2fe', width: '64px', height: '64px', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
                        <ShieldCheck size={32} color="#0284c7" />
                    </div>
                    <h2 style={{ fontSize: '1.4rem', color: '#111', margin: '0 0 8px', fontWeight: 'bold' }}>Xác thực Email</h2>
                    <p style={{ color: '#555', fontSize: '0.95rem', margin: '0', lineHeight: '1.5' }}>
                        Mã xác thực gồm 6 chữ số đã được gửi tới <br />
                        <strong style={{ color: '#0284c7' }}>{email}</strong>
                    </p>
                </div>

                {error && (
                    <div style={{ background: '#fee2e2', color: '#ef4444', padding: '10px 14px', borderRadius: '8px', fontSize: '0.9rem', marginBottom: '20px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                        <AlertTriangle size={16} />
                        {error}
                    </div>
                )}

                <div style={{ display: 'flex', justifyContent: 'space-between', gap: '8px', marginBottom: '24px' }}>
                    {otp.map((digit, index) => (
                        <input
                            key={index}
                            ref={el => inputRefs.current[index] = el}
                            type="text"
                            maxLength="1"
                            value={digit}
                            onChange={e => handleChange(index, e.target.value)}
                            onKeyDown={e => handleKeyDown(index, e)}
                            onPaste={handlePaste}
                            disabled={isVerifying}
                            style={{
                                width: '48px', height: '56px', fontSize: '1.5rem', textAlign: 'center',
                                border: `2px solid ${digit ? '#0284c7' : '#e5e7eb'}`,
                                borderRadius: '12px', fontWeight: 'bold', color: '#111',
                                transition: 'all 0.2s', outline: 'none', background: '#fff'
                            }}
                            onFocus={e => e.target.select()}
                        />
                    ))}
                </div>

                <button
                    onClick={handleVerifyOTP}
                    disabled={isVerifying || otp.join('').length < 6}
                    style={{
                        width: '100%', padding: '14px 0', borderRadius: '10px', border: 'none',
                        background: (isVerifying || otp.join('').length < 6) ? '#cbd5e1' : '#0ea5e9',
                        color: 'white', fontWeight: 'bold', fontSize: '1.05rem', cursor: (isVerifying || otp.join('').length < 6) ? 'not-allowed' : 'pointer',
                        transition: 'all 0.2s', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px'
                    }}
                >
                    {isVerifying ? <Loader className="spin" size={20} /> : <CheckCircle2 size={20} />}
                    {actionLabel}
                </button>

                <div style={{ textAlign: 'center', marginTop: '20px', fontSize: '0.95rem' }}>
                    <span style={{ color: '#666' }}>Không nhận được mã? </span>
                    {canResend ? (
                        <button
                            onClick={handleSendOTP}
                            disabled={isSending}
                            style={{ background: 'none', border: 'none', color: '#0ea5e9', fontWeight: 'bold', cursor: isSending ? 'not-allowed' : 'pointer', padding: 0 }}
                        >
                            {isSending ? 'Đang gửi...' : 'Gửi lại mã'}
                        </button>
                    ) : (
                        <span style={{ color: '#94a3b8', fontWeight: '500' }}>Gửi lại sau {countdown}s</span>
                    )}
                </div>
            </div>
        </div>
    );
};

export default OTPModal;
