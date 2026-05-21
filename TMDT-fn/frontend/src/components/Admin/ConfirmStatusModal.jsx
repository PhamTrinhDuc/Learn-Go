import React from 'react';
import { AlertTriangle, UserCheck, UserMinus, X, Loader } from 'lucide-react';

const ConfirmStatusModal = ({ isOpen, onClose, onConfirm, user, isProcessing }) => {
    if (!isOpen || !user) return null;

    const isLocking = !user.is_lock;

    return (
        <div className="admin-modal-overlay">
            <div className="admin-modal" style={{ maxWidth: '450px' }}>
                <div className="admin-modal-header">
                    <h2 style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                        {isLocking ? (
                            <UserMinus size={22} style={{ color: 'var(--admin-danger)' }} />
                        ) : (
                            <UserCheck size={22} style={{ color: 'var(--admin-success)' }} />
                        )}
                        Xác nhận {isLocking ? 'khóa' : 'mở'} tài khoản
                    </h2>
                    <button className="admin-btn" onClick={onClose} disabled={isProcessing}>
                        <X size={20} />
                    </button>
                </div>
                <div className="admin-modal-body" style={{ padding: '24px', textAlign: 'center' }}>
                    <div style={{
                        width: '64px',
                        height: '64px',
                        borderRadius: '50%',
                        background: isLocking ? 'rgba(239, 68, 68, 0.1)' : 'rgba(34, 197, 94, 0.1)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        margin: '0 auto 16px',
                        color: isLocking ? 'var(--admin-danger)' : 'var(--admin-success)'
                    }}>
                        <AlertTriangle size={32} />
                    </div>
                    <p style={{ margin: '0 0 8px 0', fontSize: '1.1rem', fontWeight: '500' }}>
                        Bạn có chắc chắn muốn {isLocking ? 'khóa' : 'mở khóa'} tài khoản này?
                    </p>
                    <p style={{ margin: 0, color: 'var(--admin-text-muted)', fontSize: '0.95rem' }}>
                        Người dùng: <strong>{user.full_name}</strong> ({user.email})
                    </p>
                    {isLocking && (
                        <p style={{ marginTop: '12px', fontSize: '0.85rem', color: 'var(--admin-danger)', background: 'rgba(239, 68, 68, 0.05)', padding: '8px', borderRadius: '4px' }}>
                            * Khi bị khóa, người dùng này sẽ không thể đăng nhập vào hệ thống.
                        </p>
                    )}
                </div>
                <div className="admin-modal-footer" style={{ justifyContent: 'center', gap: '12px' }}>
                    <button
                        className="admin-btn admin-btn-outline"
                        onClick={onClose}
                        style={{ minWidth: '100px' }}
                        disabled={isProcessing}
                    >
                        Hủy bỏ
                    </button>
                    <button
                        className={`admin-btn ${isLocking ? 'admin-btn-danger' : 'admin-btn-primary'}`}
                        onClick={onConfirm}
                        style={{ minWidth: '100px' }}
                        disabled={isProcessing}
                    >
                        {isProcessing ? <Loader className="spin" size={18} /> : (isLocking ? 'Xác nhận Khóa' : 'Xác nhận Mở')}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ConfirmStatusModal;
