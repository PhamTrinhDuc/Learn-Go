import React from 'react';
import Modal from './Modal';
import './ConfirmModal.css';

const ConfirmModal = ({ isOpen, onClose, onConfirm, title, message, confirmText = 'Xác nhận', cancelText = 'Hủy' }) => {
    return (
        <Modal isOpen={isOpen} onClose={onClose} title={title}>
            <div className="confirm-modal-content">
                <p className="confirm-message">{message}</p>
                <div className="confirm-actions">
                    <button className="btn-cancel" onClick={onClose}>
                        {cancelText}
                    </button>
                    <button className="btn-confirm" onClick={onConfirm}>
                        {confirmText}
                    </button>
                </div>
            </div>
        </Modal>
    );
};

export default ConfirmModal;
