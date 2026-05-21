import { handleEmailChange, handlePhoneChange, validateEmail, validatePhone } from "../../func/phoneValidation.js";
import AddressMap from "../Address/AddressMap.jsx";
import { Loader } from "lucide-react";
import React, { useEffect, useState, useRef } from "react";
import toast from "react-hot-toast";

const API = import.meta.env.VITE_SERVER_API;
const LOCATION_API = `${API}/api/address`;

const EMPTY_FORM = {
    id: null,
    name: '',
    hotline: '',
    province: '',
    ward: '',
    road: '',
    mail: '',
};


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

const StoreForm = ({ storeInfo, setIsEditing, onSaveSuccess }) => {
    const [editForm, setEditForm] = useState(storeInfo ? {
        id: storeInfo.id,
        name: storeInfo.name,
        hotline: storeInfo.hotline || '',
        province: storeInfo.province,
        ward: storeInfo.ward,
        road: storeInfo.road,
        mail: storeInfo.mail || ''
    } : EMPTY_FORM);

    const [isSaving, setIsSaving] = useState(false);
    const [phoneError, setPhoneError] = useState('');
    const [emailError, setEmailError] = useState('');

    const [provinces, setProvinces] = useState([]);
    const [wards, setWards] = useState([]);

    const provinceOptions = provinces.map(p => ({ value: p.code, label: p.name }));
    const wardOptions = wards.map(w => ({ value: w.code, label: w.name }));

    useEffect(() => {
        const fetchProvince = async () => {
            try {
                const res = await fetch(`${LOCATION_API}/provinces`);
                const data = await res.json();
                setProvinces(data);
            } catch (err) {
                console.error('Fetch province error:', err);
            }
        };
        fetchProvince();
    }, []);

    useEffect(() => {
        if (provinces.length > 0 && editForm.province) {
            const p = provinces.find(prov => prov.name === editForm.province);
            if (p) {
                fetch(`${LOCATION_API}/wards/${p.code}`)
                    .then(res => res.json())
                    .then(data => setWards(data))
                    .catch(() => { });
            }
        }
    }, [provinces, editForm.province]);

    const handleChangeProvince = async (selectedProvince) => {
        if (!selectedProvince) return;
        const { value: provinceCode, label: provinceName } = selectedProvince;
        setEditForm(prev => ({ ...prev, province: provinceName, ward: '' }));
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
        setEditForm(prev => ({ ...prev, ward: selectedWard.label }));
    };

    const handleMapAddressFound = ({ detail, province, ward }) => {
        const fuzzyMatch = (list, name) =>
            list.find(item =>
                item.name.toLowerCase().includes(name.toLowerCase()) ||
                name.toLowerCase().includes(item.name.toLowerCase())
            );

        const matchedProvince = province ? fuzzyMatch(provinces, province) : null;
        const matchedWard = ward ? fuzzyMatch(wards, ward) : null;

        setEditForm(prev => ({
            ...prev,
            road: detail || prev.road,
            province: matchedProvince?.name || prev.province,
            ward: matchedWard?.name || ward || prev.ward,
        }));
    };

    const handleSaveStoreConfig = async (e) => {
        e.preventDefault();
        if (isSaving) return;

        if (!editForm.name || !editForm.province || !editForm.ward || !editForm.road) {
            toast.error('Vui lòng điền đủ: Tên cửa hàng, Tỉnh/Thành, Phường/Xã và Địa chỉ đường!', { id: 'store-error' });
            return;
        }

        const pErr = validatePhone(editForm.hotline);
        setPhoneError(pErr);
        if (pErr) return;

        const eErr = validateEmail(editForm.mail);
        setEmailError(eErr);
        if (eErr) return;

        setIsSaving(true);
        try {
            const payload = {
                id: editForm.id || undefined,
                name: editForm.name,
                hotline: editForm.hotline || null,
                province: editForm.province,
                ward: editForm.ward,
                road: editForm.road,
                mail: editForm.mail || null,
            };

            const res = await fetch(`${API}/api/address/store`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            const data = await res.json();

            if (!res.ok) throw new Error(data.message || 'Lưu thất bại');

            toast.success(data.message || 'Cập nhật thông tin cửa hàng thành công!', { id: 'store-success' });
            if (onSaveSuccess) onSaveSuccess(data.data);
            setIsEditing(false);
        } catch (err) {
            toast.error(err.message || 'Lỗi hệ thống khi lưu cửa hàng.', { id: 'store-error' });
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="info-edit-mode">
            <form id="edit-store-form" onSubmit={handleSaveStoreConfig} className="shopee-style-form-wrapper">

                <div className="form-row">
                    <div className="form-group flex-1">
                        <label>Tên cửa hàng <span className="required">*</span></label>
                        <input
                            type="text"
                            className="form-input"
                            placeholder="Nhập tên cửa hàng / thương hiệu"
                            value={editForm.name}
                            onChange={e => setEditForm(prev => ({ ...prev, name: e.target.value }))}
                            required
                        />
                    </div>
                    <div className="form-group flex-1">
                        <label>Số điện thoại Hotline</label>
                        <input
                            type="tel"
                            className={`form-input ${phoneError ? 'input-error' : ''}`}
                            placeholder="VD: 0912345678"
                            value={editForm.hotline}
                            onChange={e => {
                                const { cleaned, error } = handlePhoneChange(e.target.value);
                                setEditForm(prev => ({ ...prev, hotline: cleaned }));
                                setPhoneError(error);
                            }}
                        />
                        {phoneError && <div className="field-error-msg">{phoneError}</div>}
                    </div>
                </div>
                <div className="form-group">
                    <label>Email liên hệ</label>
                    <input
                        type="email"
                        className={`form-input ${emailError ? 'input-error' : ''}`}
                        placeholder="VD: lienhe@cuahang.com"
                        value={editForm.mail}
                        onChange={e => {
                            const { cleaned, error } = handleEmailChange(e.target.value);
                            setEditForm(prev => ({ ...prev, mail: cleaned }));
                            setEmailError(error);
                        }}
                    />
                    {emailError && <div className="field-error-msg">{emailError}</div>}
                </div>

                <label className="section-label">Địa chỉ cửa hàng</label>

                <div className="map-wrapper">
                    <AddressMap onAddressFound={handleMapAddressFound} API={API} provinces={provinces} />
                </div>

                <div className="form-group">
                    <CustomSelect
                        options={provinceOptions}
                        value={editForm.province}
                        onChange={handleChangeProvince}
                        placeholder="Chọn Tỉnh/Thành phố *"
                    />
                </div>
                <div className="form-group">
                    <CustomSelect
                        options={wardOptions}
                        value={editForm.ward}
                        onChange={handleChangeWard}
                        placeholder="Chọn Phường/Xã *"
                        disabled={!editForm.province}
                    />
                </div>
                <div className="form-group">
                    <textarea
                        className="form-input"
                        placeholder="Số nhà, tên đường... *"
                        value={editForm.road}
                        onChange={e => setEditForm(prev => ({ ...prev, road: e.target.value }))}
                        required
                        rows="3"
                    />
                </div>

                <div className="form-footer">
                    <button type="button" className="btn-cancel" onClick={() => {
                        setIsEditing(false);
                    }}>
                        Hủy bỏ
                    </button>
                    <button type="submit" className="btn-save" disabled={isSaving}>
                        {isSaving ? <><Loader className="spin" size={16} /> Đang lưu...</> : 'Lưu Thay Đổi'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default StoreForm;