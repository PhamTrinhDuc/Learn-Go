import React, { useState, useEffect } from 'react';
import { MapPin, Plus, Loader, ChevronLeft, X } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import toast from 'react-hot-toast';
import './AddressManager.css';
import CheckoutAddressForm from '../../components/Address/CheckoutAddressForm.jsx';

const API = import.meta.env.VITE_SERVER_API;

const AddressManager = () => {
    const { user } = useAuth();
    const [addresses, setAddresses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showAddressForm, setShowAddressForm] = useState(false);
    const [editingAddress, setEditingAddress] = useState(null);
    const [isSaving, setIsSaving] = useState(false);

    const [newAddress, setNewAddress] = useState({
        full_name: '',
        num_phone: '',
        province: '',
        district: '',
        ward: '',
        detail_address: '',
        is_default: false
    });

    useEffect(() => {
        fetchAddresses();
    }, [user?.id]);

    const fetchAddresses = async () => {
        if (!user?.id) return;
        setLoading(true);
        try {
            const res = await fetch(`${API}/api/user/address/${user.id}`);
            const data = await res.json();
            if (data.success) {
                setAddresses(data.data);
            }
        } catch (err) {
            console.error('Fetch addresses error:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleOpenForm = (addr = null) => {
        if (addr) {
            setEditingAddress(addr);
            setNewAddress({
                full_name: addr.full_name,
                num_phone: addr.num_phone,
                province: addr.province,
                district: addr.district,
                ward: addr.ward,
                detail_address: addr.detail_address,
                is_default: addr.is_default
            });
        } else {
            setEditingAddress(null);
            setNewAddress({
                full_name: user?.full_name || '',
                num_phone: user?.num_phone || '',
                province: '',
                district: '',
                ward: '',
                detail_address: '',
                is_default: addresses.length === 0
            });
        }
        setShowAddressForm(true);
    };

    const handleCloseForm = () => {
        setShowAddressForm(false);
        setEditingAddress(null);
    };

    const handleDelete = (addr, e = null) => {
        if (e) e.stopPropagation();

        if (addr.is_default) {
            toast.error('Không thể xóa địa chỉ mặc định!', { id: 'addr-error' });
            return;
        }

        toast.custom((t) => (
            <div className={`confirm-modal-overlay ${t.visible ? 'active' : ''}`} onClick={() => toast.dismiss(t.id)}>
                <div className="confirm-toast" id="confirm-delete" onClick={e => e.stopPropagation()}>
                    <span className="confirm-toast-title">Bạn có chắc chắn muốn xóa địa chỉ này?</span>
                    <div className="confirm-toast-actions">
                        <button className="btn-confirm-yes" onClick={async () => {
                            toast.dismiss(t.id);
                            await performDelete(addr.id);
                            if (showAddressForm) handleCloseForm();
                        }}>Xóa</button>
                        <button className="btn-confirm-no" onClick={() => toast.dismiss(t.id)}>Hủy</button>
                    </div>
                </div>
            </div>
        ), { id: 'confirm-delete', duration: 6000 });
    };

    const performDelete = async (id) => {
        try {
            const res = await fetch(`${API}/api/user/address/${id}`, { method: 'DELETE' });
            if (res.ok) {
                toast.success('Đã xóa địa chỉ', { id: 'addr-success' });
                setAddresses(prev => prev.filter(a => a.id !== id));
            }
        } catch (err) {
            toast.error('Lỗi khi xóa địa chỉ', { id: 'addr-error' });
        }
    };

    const setAsDefault = async (id, e) => {
        e.stopPropagation();
        try {
            const res = await fetch(`${API}/api/user/address/${id}/default`, {
                method: 'PATCH',
                body: JSON.stringify({ userId: user.id }),
                headers: { 'Content-Type': 'application/json' }
            });
            if (res.ok) {
                toast.success('Đã chọn làm địa chỉ mặc định', { id: 'addr-success' });
                fetchAddresses();
            }
        } catch (err) {
            toast.error('Lỗi khi cập nhật địa chỉ mặc định', { id: 'addr-error' });
        }
    };

    if (loading && !showAddressForm) {
        return (
            <div className="address-manager-loading">
                <Loader className="spin" size={32} />
                <p>Đang tải địa chỉ...</p>
            </div>
        );
    }

    return (
        <div className="address-manager">
            <div className="address-list-view">
                <div className="address-header">
                    <h1>Địa chỉ của Tôi</h1>
                </div>

                <div className="address-items">
                    {addresses.length === 0 ? (
                        <div className="empty-addresses">
                            <MapPin size={48} opacity={0.2} />
                            <p>Bạn chưa có địa chỉ nào</p>
                        </div>
                    ) : (
                        addresses.map(addr => (
                            <div key={addr.id} className={`address-card ${addr.is_default ? 'is-default' : ''}`}>
                                <div className="address-card-main">
                                    <div className="addr-user-info">
                                        <span className="addr-name">{addr.full_name}</span>
                                        <span className="addr-divider">|</span>
                                        <span className="addr-phone">{addr.num_phone}</span>
                                    </div>
                                    <div className="addr-text">
                                        {addr.detail_address}
                                    </div>
                                    <div className="addr-text">
                                        {addr.ward}, {addr.province}
                                    </div>
                                    {addr.is_default && (
                                        <span className="addr-default-badge">Mặc định</span>
                                    )}
                                </div>
                                <div className="address-card-actions">
                                    <button className="btn-edit" onClick={() => handleOpenForm(addr)}>Sửa</button>
                                </div>
                            </div>
                        ))
                    )}
                </div>

                <button className="btn-add-new-shopee" onClick={() => handleOpenForm()}>
                    <Plus size={18} /> Thêm Địa Chỉ Mới
                </button>
            </div>

            {showAddressForm && (
                <div className="address-modal-overlay">
                    <div className="address-modal-content address-form-modal" onClick={e => e.stopPropagation()}>
                        <div className="address-modal-header">
                            <h3>{editingAddress ? 'Cập nhật địa chỉ' : 'Địa chỉ mới'}</h3>
                            <button className="btn-close-modal" onClick={handleCloseForm}>
                                <X size={20} />
                            </button>
                        </div>
                        <div className="address-modal-body">
                            <CheckoutAddressForm
                                newAddress={newAddress}
                                setNewAddress={setNewAddress}
                                handleCloseAddForm={handleCloseForm}
                                isSaving={isSaving}
                                API={API}
                                user={user}
                                setSavedAddresses={setAddresses}
                                setSelectedAddressId={() => { }}
                                setShowAddForm={setShowAddressForm}
                                setIsSaving={setIsSaving}
                                isEditing={!!editingAddress}
                                editingId={editingAddress?.id}
                                isDefault={editingAddress?.is_default}
                                fetchAddresses={fetchAddresses}
                                onDelete={editingAddress ? () => handleDelete(editingAddress) : null}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AddressManager;