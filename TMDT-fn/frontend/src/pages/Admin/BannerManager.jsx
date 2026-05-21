import React, { useState, useEffect, useRef } from 'react';
import toast from 'react-hot-toast';
import { Save, Plus, Trash2, Image as ImageIcon, Upload, X } from 'lucide-react';
import { getBanners, saveBanners } from '../../func/bannerStore';
import './StoreConfigManager.css';

const BannerManager = () => {
    const [config, setConfig] = useState({ topBanner: {}, sliders: [] });
    const topBannerInputRef = useRef(null);
    const sliderInputRefs = useRef({});

    useEffect(() => {
        const stored = getBanners();
        setConfig(stored);
    }, []);

    const handleSave = () => {
        try {
            saveBanners(config);
            toast.success("Đã lưu thiết lập các banner thành công!");
            window.dispatchEvent(new Event('bannerConfigChanged'));
        } catch (error) {
            console.error("Lỗi khi lưu banner vào localStorage:", error);
            toast.error("Không thể lưu banner! Dung lượng hình ảnh quá lớn, vui lòng nén hoặc giảm kích thước ảnh xuống dưới 1MB và thử lại.");
        }
    };

    const handleTopBannerChange = (field, value) => {
        setConfig((prev) => ({
            ...prev,
            topBanner: { ...prev.topBanner, [field]: value },
        }));
    };

    const handleFileChange = (e, targetType, sliderId = null) => {
        const file = e.target.files[0];
        if (!file) return;

        if (!file.type.startsWith('image/')) {
            toast.error("Vui lòng chọn tệp hình ảnh!");
            return;
        }

        // Limit size to 2MB to keep localStorage performance healthy
        if (file.size > 2 * 1024 * 1024) {
            toast.error("Kích thước ảnh quá lớn! Vui lòng chọn ảnh dưới 2MB.");
            return;
        }

        const reader = new FileReader();
        reader.onload = (event) => {
            const base64Url = event.target.result;
            if (targetType === 'topBanner') {
                setConfig(prev => ({
                    ...prev,
                    topBanner: { ...prev.topBanner, imageUrl: base64Url }
                }));
                toast.success("Đã chọn ảnh Top Banner mới!");
            } else if (targetType === 'slider' && sliderId) {
                setConfig(prev => ({
                    ...prev,
                    sliders: prev.sliders.map(s =>
                        s.id === sliderId ? { ...s, imageUrl: base64Url } : s
                    )
                }));
                toast.success("Đã chọn ảnh slide mới!");
            }
        };
        reader.readAsDataURL(file);
    };

    const handleRemoveTopBannerImage = () => {
        setConfig(prev => ({
            ...prev,
            topBanner: { ...prev.topBanner, imageUrl: '' }
        }));
        toast.success("Đã xóa ảnh Top Banner!");
    };

    const handleAddSlider = () => {
        const newSlide = {
            id: Date.now(),
            imageUrl: '',
            link: '/'
        };
        setConfig((prev) => ({
            ...prev,
            sliders: [...prev.sliders, newSlide]
        }));
    };

    const handleRemoveSlider = (id) => {
        setConfig((prev) => ({
            ...prev,
            sliders: prev.sliders.filter(s => s.id !== id)
        }));
        toast.success("Đã xóa slide!");
    };

    const handleUpdateSlider = (id, field, value) => {
        setConfig((prev) => ({
            ...prev,
            sliders: prev.sliders.map(s =>
                s.id === id ? { ...s, [field]: value } : s
            )
        }));
    };

    const handleRemoveSliderImage = (id) => {
        setConfig(prev => ({
            ...prev,
            sliders: prev.sliders.map(s =>
                s.id === id ? { ...s, imageUrl: '' } : s
            )
        }));
        toast.success("Đã xóa ảnh slide!");
    };

    return (
        <div className="store-config-manager">
            <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h2>Quản lý Banner (Top & Slider)</h2>
                <button onClick={handleSave} className="btn-save" style={{ display: 'flex', gap: '8px', alignItems: 'center', boxShadow: '0 2px 4px rgba(59, 130, 246, 0.2)' }}>
                    <Save size={18} /> Lưu Thay Đổi
                </button>
            </div>

            <div className="config-card" style={{ marginBottom: '24px' }}>
                <div className="card-header">
                    <ImageIcon size={20} className="header-icon" />
                    <h3>Top Banner (Thanh Ngang Trên Cùng)</h3>
                </div>

                <div style={{ padding: '24px' }} className="shopee-style-form-wrapper">
                    <div style={{ marginBottom: '16px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                        <input
                            type="checkbox"
                            id="activeTopBanner"
                            checked={config.topBanner.active || false}
                            onChange={(e) => handleTopBannerChange('active', e.target.checked)}
                            style={{ width: '18px', height: '18px', cursor: 'pointer' }}
                        />
                        <label htmlFor="activeTopBanner" style={{ fontWeight: '600', cursor: 'pointer', margin: 0 }}>Hiển thị Top Banner</label>
                    </div>

                    <div className="form-group" style={{ marginTop: '20px' }}>
                        <label style={{ fontWeight: '600', marginBottom: '10px' }}>Hình ảnh Top Banner</label>
                        
                        <input 
                            type="file" 
                            ref={topBannerInputRef} 
                            style={{ display: 'none' }} 
                            accept="image/*"
                            onChange={(e) => handleFileChange(e, 'topBanner')}
                        />

                        {config.topBanner.imageUrl ? (
                            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                                <div style={{ border: '1px solid #e2e8f0', padding: '6px', borderRadius: '8px', background: '#f8fafc', boxShadow: 'inset 0 2px 4px rgba(0,0,0,0.02)' }}>
                                    <img 
                                        src={config.topBanner.imageUrl} 
                                        alt="Top Banner Preview" 
                                        style={{ width: '100%', height: 'auto', maxHeight: '120px', objectFit: 'cover', borderRadius: '4px' }} 
                                    />
                                </div>
                                <div style={{ display: 'flex', gap: '10px' }}>
                                    <button 
                                        type="button"
                                        onClick={() => topBannerInputRef.current?.click()}
                                        style={{ background: '#fff', color: '#374151', border: '1px solid #d1d5db', padding: '8px 14px', borderRadius: '6px', fontSize: '13px', display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer', fontWeight: '500', transition: 'all 0.15s' }}
                                        onMouseEnter={(e) => e.currentTarget.style.background = '#f9fafb'}
                                        onMouseLeave={(e) => e.currentTarget.style.background = '#fff'}
                                    >
                                        <Upload size={16} /> Thay đổi ảnh
                                    </button>
                                    <button 
                                        type="button"
                                        onClick={handleRemoveTopBannerImage}
                                        style={{ background: '#fef2f2', color: '#ef4444', border: '1px solid #fee2e2', padding: '8px 14px', borderRadius: '6px', fontSize: '13px', display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer', fontWeight: '500' }}
                                    >
                                        <Trash2 size={16} /> Xóa ảnh
                                    </button>
                                </div>
                            </div>
                        ) : (
                            <div 
                                onClick={() => topBannerInputRef.current?.click()}
                                style={{ 
                                    border: '2px dashed #cbd5e1', 
                                    borderRadius: '12px', 
                                    padding: '30px 20px', 
                                    textAlign: 'center', 
                                    cursor: 'pointer', 
                                    background: '#f8fafc',
                                    transition: 'all 0.2s',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    alignItems: 'center',
                                    gap: '8px'
                                }}
                                onMouseEnter={(e) => {
                                    e.currentTarget.style.borderColor = '#3b82f6';
                                    e.currentTarget.style.background = '#eff6ff';
                                }}
                                onMouseLeave={(e) => {
                                    e.currentTarget.style.borderColor = '#cbd5e1';
                                    e.currentTarget.style.background = '#f8fafc';
                                }}
                            >
                                <Upload size={32} style={{ color: '#64748b' }} />
                                <span style={{ fontSize: '14px', fontWeight: '500', color: '#475569' }}>Nhấp vào đây để chọn ảnh từ thiết bị</span>
                                <span style={{ fontSize: '12px', color: '#64748b' }}>Hỗ trợ tệp PNG, JPG, JPEG (Kích thước khuyên dùng tỷ lệ ngang dài, dưới 2MB)</span>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            <div className="config-card">
                <div className="card-header" style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                        <ImageIcon size={20} className="header-icon" style={{ margin: 0 }} />
                        <h3 style={{ margin: 0 }}>Banner Slider (Trang Chủ)</h3>
                    </div>
                    <button onClick={handleAddSlider} style={{ background: '#0284c7', color: 'white', border: 'none', padding: '6px 12px', borderRadius: '6px', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '4px', fontWeight: '500', fontSize: '13px' }}>
                        <Plus size={16} /> Thêm ảnh
                    </button>
                </div>

                <div style={{ padding: '24px' }} className="shopee-style-form-wrapper">
                    {config.sliders.length === 0 ? (
                        <p style={{ color: '#64748b', textAlign: 'center', padding: '20px 0' }}>Chưa có slide nào. Vui lòng thêm slide mới!</p>
                    ) : (
                        <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '16px' }}>
                            {config.sliders.map((slider, idx) => (
                                <div key={slider.id} style={{ border: '1px solid #e2e8f0', borderRadius: '8px', padding: '16px', display: 'flex', gap: '16px', background: '#f8fafc', position: 'relative' }}>
                                    
                                    <input 
                                        type="file" 
                                        ref={el => sliderInputRefs.current[slider.id] = el}
                                        style={{ display: 'none' }} 
                                        accept="image/*"
                                        onChange={(e) => handleFileChange(e, 'slider', slider.id)}
                                    />

                                    {/* Cột trái: Ảnh preview & Nút thay đổi */}
                                    <div style={{ flex: '0 0 200px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                                        <div style={{ width: '100%', height: '110px', background: '#e2e8f0', borderRadius: '6px', overflow: 'hidden', display: 'flex', alignItems: 'center', justifyContent: 'center', border: '1px solid #cbd5e1', boxShadow: 'inset 0 2px 4px rgba(0,0,0,0.02)' }}>
                                            {slider.imageUrl ? (
                                                <img src={slider.imageUrl} alt={`Slide ${idx + 1}`} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                                            ) : (
                                                <div style={{ textAlign: 'center', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '4px', padding: '10px' }}>
                                                    <ImageIcon size={24} style={{ color: '#94a3b8' }} />
                                                    <span style={{ fontSize: '0.75rem', color: '#64748b' }}>Chưa có hình ảnh</span>
                                                </div>
                                            )}
                                        </div>
                                        
                                        <div style={{ display: 'flex', gap: '6px' }}>
                                            <button 
                                                type="button"
                                                onClick={() => sliderInputRefs.current[slider.id]?.click()}
                                                style={{ flex: 1, background: '#fff', color: '#475569', border: '1px solid #cbd5e1', padding: '6px 8px', borderRadius: '4px', fontSize: '12px', fontWeight: '500', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '4px', cursor: 'pointer', transition: 'all 0.15s' }}
                                                onMouseEnter={(e) => e.currentTarget.style.background = '#f1f5f9'}
                                                onMouseLeave={(e) => e.currentTarget.style.background = '#fff'}
                                            >
                                                <Upload size={13} /> {slider.imageUrl ? 'Thay ảnh' : 'Tải ảnh lên'}
                                            </button>
                                            {slider.imageUrl && (
                                                <button 
                                                    type="button"
                                                    onClick={() => handleRemoveSliderImage(slider.id)}
                                                    style={{ background: '#fef2f2', color: '#ef4444', border: '1px solid #fee2e2', padding: '6px 8px', borderRadius: '4px', fontSize: '12px', cursor: 'pointer' }}
                                                    title="Xóa ảnh"
                                                >
                                                    <X size={14} />
                                                </button>
                                            )}
                                        </div>
                                    </div>

                                    {/* Cột giữa: Cấu hình liên kết */}
                                    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                                        <div className="form-group" style={{ marginBottom: '0' }}>
                                            <label style={{ fontSize: '13px', fontWeight: '600', color: '#475569', marginBottom: '6px' }}>Đường dẫn khi click (URL hoặc Slug)</label>
                                            <input
                                                type="text"
                                                className="form-input"
                                                value={slider.link || '/'}
                                                onChange={(e) => handleUpdateSlider(slider.id, 'link', e.target.value)}
                                                placeholder="Ví dụ: / hoặc /op-lung-iphone"
                                                style={{ padding: '8px 12px', fontSize: '13px' }}
                                            />
                                        </div>
                                    </div>

                                    {/* Cột phải: Nút xóa slide toàn bộ */}
                                    <button
                                        onClick={() => handleRemoveSlider(slider.id)}
                                        style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer', alignSelf: 'center', padding: '8px', transition: 'all 0.2s', borderRadius: '50%', width: '36px', height: '36px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                                        onMouseEnter={(e) => e.currentTarget.style.background = '#fef2f2'}
                                        onMouseLeave={(e) => e.currentTarget.style.background = 'none'}
                                        title="Xóa slide"
                                    >
                                        <Trash2 size={18} />
                                    </button>

                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default BannerManager;
