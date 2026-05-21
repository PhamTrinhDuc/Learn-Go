import React from 'react';
import { Link } from 'react-router-dom';
import Modal from '../ReUse/Modal';
import './AuthModal.css';

const AuthModal = ({ isOpen, onClose, actionName, redirectPath }) => {
    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Thông báo">
            <div className="auth-modal-content">
                <p>Vui lòng đăng nhập để {actionName}</p>
                <div className="auth-modal-actions">
                    <button className="btn-back" onClick={onClose}>Trở lại</button>
                    <Link to={redirectPath} className="btn-login-link">
                        <button className="btn-login-submit">Đăng nhập</button>
                    </Link>
                </div>
            </div>
        </Modal>
    );
};

export default AuthModal;
