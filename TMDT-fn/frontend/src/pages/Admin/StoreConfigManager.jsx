import React, { useState, useEffect, useRef } from 'react';
import { Loader, MapPin, Phone, Store, Mail } from 'lucide-react';
import toast from 'react-hot-toast';
import './StoreConfigManager.css';
import StoreForm from '../../components/Admin/StoreForm';
import DisplayForm from '../../components/Admin/DisplayForm';

const API = import.meta.env.VITE_SERVER_API;

const StoreConfigManager = () => {
    const [storeInfo, setStoreInfo] = useState(null);
    const [displayInfo, setDisplayInfo] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    const [isEditing, setIsEditing] = useState(false);
    const [isEditingDisplay, setIsEditingDisplay] = useState(false);

    const fetchInfo = async () => {
        setIsLoading(true);
        try {
            const [storeRes, displayRes] = await Promise.all([
                fetch(`${API}/api/address/store`),
                fetch(`${API}/api/pagination`)
            ]);
            const storeData = await storeRes.json();
            const displayData = await displayRes.json();
            setDisplayInfo(displayData?.data?.[0] ?? null);
            if (storeData.success && storeData.data.length > 0) {
                setStoreInfo(storeData.data[0]);
            } else {
                setStoreInfo(null);
            }
        } catch (err) {
            console.error('Fetch store error:', err);
            toast.error('Không thể tải thông tin cửa hàng.', { id: 'store-error' });
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchInfo();
    }, []);

    const handleSaveSuccess = (newData) => {
        setStoreInfo(newData);
        setIsEditing(false);
    };

    const handleSaveDisplaySuccess = (newData) => {
        setDisplayInfo(newData);
        setIsEditingDisplay(false);
    };

    const openEditForm = () => {
        setIsEditing(true);
    };

    const openEditDisplayForm = () => {
        setIsEditingDisplay(true);
    };

    return (
        <div className="store-config-manager">
            <div className="page-header">
                <h2>Cấu hình Cửa hàng</h2>
            </div>

            <div className="config-card">
                <div className="card-header">
                    <Store size={20} className="header-icon" />
                    <h3>Thông tin cửa hàng</h3>
                    {!isEditing && (
                        <button className="btn-edit" onClick={openEditForm} disabled={isLoading}>
                            {storeInfo ? 'Sửa thông tin' : '+ Thêm cửa hàng'}
                        </button>
                    )}
                </div>

                {isLoading ? (
                    <div className="loading-state">
                        <Loader className="spin" size={28} />
                        <p>Đang tải thông tin...</p>
                    </div>
                ) : !isEditing ? (
                    storeInfo ? (
                        <div className="info-view-mode">
                            <div className="info-row">
                                <Store className="info-icon" size={18} />
                                <div className="info-content">
                                    <span className="label">Tên gian hàng</span>
                                    <span className="value">{storeInfo.name}</span>
                                </div>
                            </div>
                            <div className="info-row">
                                <Phone className="info-icon" size={18} />
                                <div className="info-content">
                                    <span className="label">Hotline bán hàng</span>
                                    <span className="value hotline">{storeInfo.hotline || '—'}</span>
                                </div>
                            </div>
                            <div className="info-row">
                                <Mail className="info-icon" size={18} />
                                <div className="info-content">
                                    <span className="label">Email liên hệ</span>
                                    <span className="value">{storeInfo.mail || '—'}</span>
                                </div>
                            </div>
                            <div className="info-row align-top">
                                <MapPin className="info-icon" size={18} />
                                <div className="info-content">
                                    <span className="label">Địa chỉ Cửa hàng / Kho</span>
                                    <span className="value">
                                        {storeInfo.road}, {storeInfo.ward}, {storeInfo.province}
                                    </span>
                                </div>
                            </div>
                        </div>
                    ) : (
                        <div className="empty-state">
                            <Store size={48} />
                            <p>Chưa có thông tin cửa hàng nào. Hãy thêm mới!</p>
                        </div>
                    )
                ) : (
                    <StoreForm
                        storeInfo={storeInfo}
                        setIsEditing={setIsEditing}
                        onSaveSuccess={handleSaveSuccess}
                    />
                )}
            </div>
            <div className="config-card display-config">
                <div className="card-header">
                    <Store size={20} className="header-icon" />
                    <h3>Thông tin Trưng bày hàng</h3>
                    {!isEditingDisplay && (
                        <button className="btn-edit" onClick={openEditDisplayForm} disabled={isLoading}>
                            Sửa thông tin
                        </button>
                    )}
                </div>

                <div className="card-body">
                    {isLoading ? (
                        <div className="loading-state">
                            <Loader className="spin" size={28} />
                            <p>Đang tải thông tin...</p>
                        </div>
                    ) : isEditingDisplay ? (
                        <DisplayForm
                            displayInfo={displayInfo}
                            setIsEditing={setIsEditingDisplay}
                            onSaveSuccess={handleSaveDisplaySuccess}
                        />
                    ) : displayInfo && (
                        <div className="display-grid">
                            <div className="display-section">
                                <div className="section-title">
                                    <h4>Trang chủ</h4>
                                </div>
                                <div className="section-content">
                                    <div className="display-item">
                                        <span className="label">Số cột:</span>
                                        <span className="value">{displayInfo.column_in_home}</span>
                                    </div>
                                    <div className="display-item">
                                        <span className="label">Số hàng:</span>
                                        <span className="value">{Math.ceil(displayInfo.item_in_home / displayInfo.column_in_home) || 0}</span>
                                    </div>
                                    <div className="display-item">
                                        <span className="label">Số sản phẩm hiển thị:</span>
                                        <span className="value">{displayInfo.item_in_home}</span>
                                    </div>
                                </div>
                            </div>

                            <div className="display-section">
                                <div className="section-title">
                                    <h4>Trang danh sách sản phẩm</h4>
                                </div>
                                <div className="section-content">
                                    <div className="display-item">
                                        <span className="label">Số cột:</span>
                                        <span className="value">{displayInfo.column_in_productlist}</span>
                                    </div>
                                    <div className="display-item">
                                        <span className="label">Số hàng:</span>
                                        <span className="value">{Math.ceil(displayInfo.item_in_productlist / displayInfo.column_in_productlist) || 0}</span>
                                    </div>
                                    <div className="display-item">
                                        <span className="label">Số sản phẩm hiển thị:</span>
                                        <span className="value">{displayInfo.item_in_productlist}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default StoreConfigManager;
