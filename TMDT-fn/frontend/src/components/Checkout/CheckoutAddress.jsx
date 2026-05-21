import React, { useState, useEffect } from 'react';
import { MapPin, Plus, ChevronRight, X, Loader } from 'lucide-react';
import toast from 'react-hot-toast';
import CheckoutAddressForm from '../Address/CheckoutAddressForm.jsx';
import './CheckoutAddress.css';

const CheckoutAddress = ({
    user,
    savedAddresses,
    setSavedAddresses,
    selectedAddressId,
    setSelectedAddressId,
    loadingAddress,
    API
}) => {
    const [showAddressManager, setShowAddressManager] = useState(false);
    const [showAddForm, setShowAddForm] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const [newAddress, setNewAddress] = useState({
        full_name: user?.full_name || '',
        num_phone: user?.num_phone || '',
        province: 'TP. Hồ Chí Minh',
        district: '',
        ward: '',
        detail_address: '',
        is_default: false
    });

    // Reset form khi user thay đổi
    useEffect(() => {
        setNewAddress(prev => ({
            ...prev,
            full_name: user?.full_name || '',
            num_phone: user?.num_phone || ''
        }));
    }, [user]);

    // Handle show Add Form
    const handleOpenAddForm = () => {
        setShowAddForm(true);
        setShowAddressManager(false);
    };

    const handleCloseAddForm = () => {
        setShowAddForm(false);
        if (savedAddresses.length > 0) {
            setShowAddressManager(true);
        }
    };

    const selectedAddress = savedAddresses.find(a => a.id === selectedAddressId);

    if (loadingAddress) {
        return (
            <div className="checkout-address-loading">
                <Loader className="spin" size={24} />
            </div>
        );
    }

    return (
        <div className="checkout-address-section">
            <h3 className="checkout-address-title">
                <MapPin size={18} /> Địa chỉ nhận hàng
            </h3>

            {/* Khối hiển thị địa chỉ đã chọn (hoặc yêu cầu thêm) */}
            <div className="checkout-address-box">
                {savedAddresses.length === 0 ? (
                    <div className="checkout-address-empty">
                        <p>Bạn chưa có địa chỉ giao hàng nào.</p>
                        <button type="button" className="btn-add-address-shopee" onClick={handleOpenAddForm}>
                            <Plus size={16} /> Thêm địa chỉ mới
                        </button>
                    </div>
                ) : selectedAddress ? (
                    <div className="checkout-address-display" onClick={() => setShowAddressManager(true)}>
                        <div className="address-display-content">
                            <div className="address-display-name-phone">
                                <span className="address-name">{selectedAddress.full_name}</span>
                                <span className="address-phone">{selectedAddress.num_phone}</span>
                            </div>
                            <div className="address-display-detail">
                                {selectedAddress.detail_address}, {selectedAddress.ward}, {selectedAddress.province}
                            </div>
                            {selectedAddress.is_default && (
                                <span className="address-default-badge">Mặc định</span>
                            )}
                        </div>
                        <div className="address-display-action">
                            <span className="change-text">Thay đổi</span>
                            <ChevronRight size={18} color="var(--text-muted)" />
                        </div>
                    </div>
                ) : (
                    // Fallback nếu không khớp ID
                    <div className="checkout-address-empty">
                        <p>Vui lòng chọn địa chỉ giao hàng.</p>
                        <button type="button" className="btn-add-address-shopee" onClick={() => setShowAddressManager(true)}>
                            Chọn địa chỉ
                        </button>
                    </div>
                )}
            </div>

            {/* 1. DROPBOX / MODAL: Address Manager (List) */}
            {showAddressManager && (
                <div className="address-modal-overlay" onClick={() => setShowAddressManager(false)}>
                    <div className="address-modal-content" onClick={e => e.stopPropagation()}>
                        <div className="address-modal-header">
                            <h3>Địa Chỉ Của Tôi</h3>
                            <button className="btn-close-modal" onClick={() => setShowAddressManager(false)}>
                                <X size={20} />
                            </button>
                        </div>
                        <div className="address-modal-body">
                            {savedAddresses.map(addr => (
                                <label key={addr.id} className="address-modal-item">
                                    <div className="address-modal-radio">
                                        <input
                                            type="radio"
                                            name="checkout_address"
                                            checked={selectedAddressId === addr.id}
                                            onChange={() => {
                                                setSelectedAddressId(addr.id);
                                                setShowAddressManager(false);
                                            }}
                                        />
                                    </div>
                                    <div className="address-modal-info">
                                        <div className="address-info-header">
                                            <strong>{addr.full_name}</strong>
                                            <span>|</span>
                                            <span>{addr.num_phone}</span>
                                        </div>
                                        <div className="address-info-desc">
                                            {addr.detail_address}
                                        </div>
                                        <div className="address-info-desc">
                                            {addr.ward}, {addr.province}
                                        </div>
                                        {addr.is_default && (
                                            <span className="address-default-badge">Mặc định</span>
                                        )}
                                    </div>
                                </label>
                            ))}
                        </div>
                        <div className="address-modal-footer">
                            <button type="button" className="btn-add-address-modal" onClick={handleOpenAddForm}>
                                <Plus size={16} /> Thêm Địa Chỉ Mới
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* 2. DROPBOX / MODAL: Add New Address Form */}
            {showAddForm && (
                <div className="address-modal-overlay" onClick={handleCloseAddForm}>
                    <div className="address-modal-content address-form-modal" onClick={e => e.stopPropagation()}>
                        <div className="address-modal-header">
                            <h3>Địa chỉ mới</h3>
                            <button className="btn-close-modal" onClick={handleCloseAddForm}>
                                <X size={20} />
                            </button>
                        </div>
                        <div className="address-modal-body">
                            <CheckoutAddressForm
                                newAddress={newAddress}
                                setNewAddress={setNewAddress}
                                // handleSaveAddress={handleSaveAddress}
                                handleCloseAddForm={handleCloseAddForm}
                                isSaving={isSaving}
                                API={API}
                                user={user}
                                setSavedAddresses={setSavedAddresses}
                                setSelectedAddressId={setSelectedAddressId}
                                setShowAddForm={setShowAddForm}
                                setIsSaving={setIsSaving}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default CheckoutAddress;
