import { React, useState, useEffect, useRef } from 'react';
import { Loader } from 'lucide-react';
import toast from 'react-hot-toast';
import AddressMap from './AddressMap';
import { validatePhone, validateName, handlePhoneChange, handleNameChange } from '../../func/phoneValidation';

const CheckoutAddressForm = ({
    newAddress,
    setNewAddress,
    handleCloseAddForm,
    isSaving,
    API,
    user,
    setSavedAddresses,
    setSelectedAddressId,
    setShowAddForm,
    setIsSaving,
    isEditing = false,
    editingId = null,
    isDefault = false,
    fetchAddresses = null,
    onDelete = null
}) => {

    const LOCATION_API = `${API}/api/address`;
    const [provinces, setProvinces] = useState([]);
    const [wards, setWards] = useState([]);
    const [nameError, setNameError] = useState('');
    const [phoneError, setPhoneError] = useState('');

    const provinceOptions = provinces.map(p => ({ value: p.code, label: p.name }));
    const wardOptions = wards.map(w => ({ value: w.code, label: w.name }));
    const CustomSelect = ({ options, value, onChange, placeholder, disabled }) => {
        const [isOpen, setIsOpen] = useState(false);
        const [searchTerm, setSearchTerm] = useState('');
        const containerRef = useRef(null);

        const filteredOptions = options.filter(opt =>
            opt.label?.toLowerCase().includes(searchTerm.toLowerCase())
        );

        useEffect(() => {
            const handleClickOutside = (event) => {
                if (containerRef.current && !containerRef.current.contains(event.target)) {
                    setIsOpen(false);
                    setSearchTerm('');
                }
            };
            document.addEventListener('mousedown', handleClickOutside);
            return () => document.removeEventListener('mousedown', handleClickOutside);
        }, []);

        return (
            <div className={`custom-select-container ${disabled ? 'disabled' : ''}`} ref={containerRef}>
                <div className="custom-select-input-wrapper">
                    <input
                        type="text"
                        className="form-input"
                        placeholder={placeholder}
                        value={isOpen ? searchTerm : (value || '')}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        onFocus={() => !disabled && setIsOpen(true)}
                        readOnly={!isOpen}
                        autoComplete="off"
                    />
                    <div className={`arrow ${isOpen ? 'up' : 'down'}`}></div>
                </div>
                {isOpen && !disabled && (
                    <div className="custom-select-dropdown">
                        {filteredOptions.length > 0 ? (
                            filteredOptions.map(opt => (
                                <div
                                    key={opt.value}
                                    className="custom-select-option"
                                    onClick={() => {
                                        onChange(opt);
                                        setIsOpen(false);
                                        setSearchTerm('');
                                    }}
                                >
                                    {opt.label}
                                </div>
                            ))
                        ) : (
                            <div className="custom-select-no-options">Không tìm thấy</div>
                        )}
                    </div>
                )}
            </div>
        );
    };


    useEffect(() => {
        const fetchProvince = async () => {
            try {
                const res = await fetch(`${LOCATION_API}/provinces`);
                const data = await res.json();
                setProvinces(data);
            } catch (err) {
                console.error('Fetch province error:', err);
            }
        }
        fetchProvince();
    }, []);

    const handleChangeProvince = async (selectedProvince) => {
        if (!selectedProvince) return;
        const { value: provinceCode, label: provinceName } = selectedProvince;
        setNewAddress({ ...newAddress, province: provinceName, ward: '' });
        try {
            const res = await fetch(`${LOCATION_API}/wards/${provinceCode}`);
            const data = await res.json();
            setWards(data);
        } catch (err) {
            console.error('Fetch wards error:', err);
        }
    };

    const handleChangeWard = (selectedWard) => {
        if (!selectedWard) return;
        setNewAddress({ ...newAddress, ward: selectedWard.label });
    };

    useEffect(() => {
        if (provinces.length > 0 && newAddress.province && wards.length === 0) {
            const p = provinces.find(prov => prov.name === newAddress.province);
            if (p) {
                fetch(`${LOCATION_API}/wards/${p.code}`)
                    .then(res => res.json())
                    .then(data => setWards(data))
                    .catch(() => { });
            }
        }
    }, [provinces, newAddress.province]);

    const handleMapAddressFound = ({ detail, province, ward }) => {
        const fuzzyMatch = (list, name) =>
            list.find(item =>
                item.name.toLowerCase().includes(name.toLowerCase()) ||
                name.toLowerCase().includes(item.name.toLowerCase())
            );

        const matchedProvince = province ? fuzzyMatch(provinces, province) : null;
        const matchedWard = ward ? fuzzyMatch(wards, ward) : null;

        setNewAddress(prev => ({
            ...prev,
            detail_address: detail || prev.detail_address,
            province: matchedProvince?.name || prev.province,
            ward: matchedWard?.name || ward || prev.ward,
        }));
    };

    const handleSaveAddress = async (e) => {
        e.preventDefault();
        if (isSaving) return;

        if (!newAddress.full_name || !newAddress.num_phone || !newAddress.detail_address) {
            toast.error('Vui lòng điền đầy đủ thông tin địa chỉ!', { id: 'addr-error' });
            return;
        }

        const nErr = validateName(newAddress.full_name);
        setNameError(nErr);
        if (nErr) return;

        const pErr = validatePhone(newAddress.num_phone);
        setPhoneError(pErr);
        if (pErr) return;

        setIsSaving(true);
        try {
            const url = isEditing
                ? `${API}/api/user/address/${editingId}`
                : `${API}/api/user/address`;
            const method = isEditing ? 'PUT' : 'POST';

            const addrRes = await fetch(url, {
                method: method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ ...newAddress, id_user: user.id }),
            });
            const addrData = await addrRes.json();

            if (!addrRes.ok) throw new Error(addrData.message || 'Lưu địa chỉ thất bại');

            toast.success(isEditing ? 'Cập nhật địa chỉ thành công' : 'Thêm địa chỉ thành công', { id: 'addr-success' });

            if (fetchAddresses) {
                await fetchAddresses();
            } else {
                const fetchRes = await fetch(`${API}/api/user/address/${user.id}`);
                const fetchData = await fetchRes.json();
                if (fetchData.success && fetchData.data.length > 0) {
                    const sorted = [...fetchData.data].sort((a, b) => b.is_default - a.is_default);
                    setSavedAddresses(sorted);
                    if (!isEditing) {
                        setSelectedAddressId(addrData.id || sorted[0].id);
                    }
                }
            }

            setShowAddForm(false);
        } catch (err) {
            toast.error(err.message, { id: 'addr-error' });
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="shopee-style-form-wrapper">
            <div id="add-address-form">
                <div className="form-row">
                    <div className="form-group flex-1">
                        <input
                            type="text"
                            className={`form-input ${nameError ? 'input-error' : ''}`}
                            placeholder="Họ và tên"
                            value={newAddress.full_name}
                            onChange={e => {
                                const { cleaned, error } = handleNameChange(e.target.value);
                                setNewAddress({ ...newAddress, full_name: cleaned });
                                setNameError(error);
                            }}
                            required
                        />
                        {nameError && <div className="field-error-msg">{nameError}</div>}
                    </div>
                    <div className="form-group flex-1">
                        <input
                            type="tel"
                            className={`form-input ${phoneError ? 'input-error' : ''}`}
                            placeholder="Số điện thoại"
                            value={newAddress.num_phone}
                            onChange={e => {
                                const { cleaned, error } = handlePhoneChange(e.target.value);
                                setNewAddress({ ...newAddress, num_phone: cleaned });
                                setPhoneError(error);
                            }}
                            required
                        />
                        {phoneError && <div className="field-error-msg">{phoneError}</div>}
                    </div>
                </div>

                <AddressMap onAddressFound={handleMapAddressFound} API={API} provinces={provinces} /><br />

                <div className="form-group">
                    <CustomSelect
                        options={provinceOptions}
                        value={newAddress.province}
                        onChange={handleChangeProvince}
                        placeholder="Chọn Tỉnh/Thành phố"
                    />
                </div>
                <div className="form-group">
                    <CustomSelect
                        options={wardOptions}
                        value={newAddress.ward}
                        onChange={handleChangeWard}
                        placeholder="Chọn Phường/Xã"
                        disabled={!newAddress.province}
                    />
                </div>
                <div className="form-group">
                    <textarea
                        className="form-input"
                        placeholder="Địa chỉ cụ thể (Số nhà, tên đường...)"
                        value={newAddress.detail_address}
                        onChange={e => setNewAddress({ ...newAddress, detail_address: e.target.value })}
                        required
                        rows="2"
                    />
                </div>

                <div
                    onClick={() => {
                        if (!(isEditing && isDefault)) {
                            setNewAddress({ ...newAddress, is_default: !newAddress.is_default });
                        }
                    }}
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        padding: '16px',
                        backgroundColor: 'white',
                        borderRadius: '8px',
                        boxShadow: '0 1px 4px rgba(0,0,0,0.08)',
                        border: '1px solid #f0f0f0',
                        marginTop: '16px',
                        cursor: (isEditing && isDefault) ? 'not-allowed' : 'pointer',
                        opacity: (isEditing && isDefault) ? 0.6 : 1,
                        userSelect: 'none'
                    }}
                >
                    <span style={{ fontSize: '0.95rem', fontWeight: 500, color: '#333' }}>
                        Đặt làm địa chỉ mặc định
                    </span>

                    <div style={{
                        position: 'relative',
                        width: '50px',
                        height: '30px',
                        backgroundColor: newAddress.is_default ? '#34C759' : '#e9e9ea',
                        borderRadius: '999px',
                        transition: 'background-color 0.3s',
                        flexShrink: 0
                    }}>
                        <div style={{
                            position: 'absolute',
                            top: '2px',
                            left: '2px',
                            width: '26px',
                            height: '26px',
                            backgroundColor: 'white',
                            borderRadius: '50%',
                            boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
                            transition: 'transform 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
                            transform: newAddress.is_default ? 'translateX(20px)' : 'translateX(0)',
                        }} />
                    </div>
                </div>
            </div>
            <div className="address-modal-footer form-footer">
                {isEditing && onDelete && (
                    <button type="button" className="btn-delete-in-form" onClick={onDelete}>Xóa địa chỉ</button>
                )}
                <button type="button" className="btn-cancel-modal" onClick={handleCloseAddForm}>Trở lại</button>
                <button type="button" className="btn-save-modal" onClick={handleSaveAddress} disabled={isSaving}>
                    {isSaving ? <Loader className="spin" size={16} /> : 'Hoàn thành'}
                </button>
            </div>
        </div>
    );
};

export default CheckoutAddressForm;
