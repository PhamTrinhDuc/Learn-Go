import React from 'react';
import { X, Upload, Loader } from 'lucide-react';
import toast from 'react-hot-toast';
import ProductForm from './ProductForm';

const ProductModal = ({
    isOpen,
    onClose,
    isEditing,
    name,
    setName,
    activeCategory,
    setActiveCategory,
    categories,
    basePrice,
    setBasePrice,
    basePriceNumeric,
    setBasePriceNumeric,
    lowStockThreshold,
    setLowStockThreshold,
    weight,
    setWeight,
    categoryMap,
    setProductData,
    thumbnail,
    setThumbnail,
    handleSaveProduct,
    isSaving,
    initialData
}) => {
    const [formKey, setFormKey] = React.useState(0);

    if (!isOpen) return null;

    const handlePasteJSON = (e) => {
        const pastedText = e.clipboardData.getData('text');
        try {
            const data = JSON.parse(pastedText);
            const pName = data.name || data.product_name || data.title;
            const pPrice = data.price || data.calculated_price || data.base_price || data.gia;
            const pWeight = data.weight || data.khoi_luong || data.mass;
            const pSpecs = data.specs || data.specifications || data.details || data.thong_so || {};
            const pVariants = data.variants || data.modelVariants || data.options || data.bien_the || [];
            const catVal = data.category_id || data.category || data.danh_muc;
            const pBasePriceNumeric = data.base_price_numeric || data.cost || data.original_price || data.gia_goc;
            const pThreshold = data.low_stock_threshold || data.nguong_canh_bao;

            if (pName) setName(pName);
            if (pPrice) setBasePrice(pPrice);
            if (pWeight) setWeight(pWeight);
            if (pBasePriceNumeric) setBasePriceNumeric(pBasePriceNumeric);
            if (pThreshold) setLowStockThreshold(pThreshold);

            if (catVal) {
                if (categoryMap[catVal]) {
                    setActiveCategory(categoryMap[catVal]);
                } else {
                    const foundLabel = Object.values(categoryMap).find(label =>
                        label.toLowerCase() === String(catVal).toLowerCase()
                    );
                    if (foundLabel) setActiveCategory(foundLabel);
                }
            }
            const newProductData = {
                specs: pSpecs,
                variants: pVariants.map((v, i) => ({
                    ...v,
                    id: v.variant_id || v.id || `v-${Date.now()}-${i}`,
                    stock: v.quantity ?? v.stock ?? v.amount ?? 0
                }))
            };

            setProductData(newProductData);
            setFormKey(prev => prev + 1);
            toast.success('Đã tự động điền thông tin từ JSON!', { id: 'prod-json' });
        } catch (err) {
            console.error('Paste JSON Error:', err);
            toast.error('Dữ liệu JSON không hợp lệ', { id: 'prod-json' });
        }
    };

    const getThumbnailSrc = (thumb) => {
        if (!thumb) return '';
        if (typeof thumb === 'string') {
            return thumb.startsWith('http') ? thumb : `${import.meta.env.VITE_PHOTO_SERVER_API}${thumb}`;
        }
        return URL.createObjectURL(thumb);
    };

    return (
        <div className="admin-modal-overlay">
            <div className="admin-modal" style={{ maxWidth: '800px' }}>
                <div className="admin-modal-header">
                    <h2>{isEditing ? 'Sửa sản phẩm' : 'Thêm sản phẩm mới'}</h2>
                    <button className="admin-btn" onClick={onClose}>×</button>
                </div>
                <div className="admin-modal-body">
                    <div className="admin-form-grid">
                        <div className="admin-form-group">
                            <label className="admin-form-label">Tên sản phẩm</label>
                            <input
                                id="product-name"
                                type="text"
                                className="admin-form-input"
                                placeholder="Nhập tên sản phẩm hoặc dán JSON..."
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                onPaste={handlePasteJSON}
                            />
                        </div>
                        <div className="admin-form-group">
                            <label className="admin-form-label">Danh mục</label>
                            <select
                                className="admin-form-input"
                                value={activeCategory}
                                onChange={(e) => setActiveCategory(e.target.value)}
                            >
                                <option value="">Chọn danh mục...</option>
                                {categories.filter(cat => cat !== 'Tất cả').map(cat => (
                                    <option key={cat} value={cat}>{cat}</option>
                                ))}
                            </select>
                        </div>
                        <div className="admin-form-group">
                            <label className="admin-form-label">Giá bán</label>
                            <input
                                id="product-price"
                                type="number"
                                className="admin-form-input"
                                placeholder="0đ"
                                value={basePrice}
                                onChange={(e) => setBasePrice(e.target.value)}
                            />
                        </div>
                        <div className="admin-form-group">
                            <label className="admin-form-label">Khối lượng (g)</label>
                            <input
                                id="product-weight"
                                type="number"
                                step="any"
                                min="0"
                                className="admin-form-input"
                                placeholder="Nhập khối lượng (VD: 3)"
                                value={weight}
                                onChange={(e) => setWeight(e.target.value)}
                            />
                        </div>

                        <div className="admin-form-group">
                            <label className="admin-form-label">Ngưỡng cảnh báo sắp hết hàng</label>
                            <input
                                id="product-low-stock-threshold"
                                type="number"
                                min="0"
                                className="admin-form-input"
                                placeholder="Mặc định: 5"
                                value={lowStockThreshold || ''}
                                onChange={(e) => setLowStockThreshold(e.target.value)}
                            />
                        </div>
                    </div>

                    {activeCategory && (
                        <div style={{ marginTop: '20px' }}>
                            <h3 className="admin-form-label" style={{ marginBottom: '12px', color: 'var(--admin-primary)' }}>Thông số kỹ thuật</h3>
                            <ProductForm
                                key={formKey}
                                isEditing={isEditing}
                                categorySlug={Object.keys(categoryMap).find(key => categoryMap[key] === activeCategory) || activeCategory.toLowerCase()}
                                onChange={setProductData}
                                initialData={initialData}
                            />
                        </div>
                    )}

                    <div className="admin-form-group" style={{ marginTop: '20px' }}>
                        <label className="admin-form-label">Ảnh đại diện (Thumbnail)</label>
                        <div
                            style={{
                                border: '2px dashed var(--admin-border)',
                                padding: '16px',
                                textAlign: 'center',
                                borderRadius: '12px',
                                background: thumbnail ? '#f8fafc' : 'transparent',
                                cursor: 'pointer',
                                position: 'relative',
                                minHeight: '120px',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center'
                            }}
                            onClick={() => document.getElementById('thumb-upload').click()}
                        >
                            {thumbnail ? (
                                <div style={{ position: 'relative', width: '100px', height: '100px' }}>
                                    <img
                                        src={getThumbnailSrc(thumbnail)}
                                        alt="Thumb preview"
                                        style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: '8px' }}
                                    />
                                    <button
                                        className="admin-btn-sm"
                                        style={{ position: 'absolute', top: '-8px', right: '-8px', background: 'var(--admin-danger)', color: '#fff', borderRadius: '50%', padding: '2px' }}
                                        onClick={(e) => { e.stopPropagation(); setThumbnail(null); }}
                                    >
                                        <X size={14} />
                                    </button>
                                </div>
                            ) : (
                                <div className="admin-text-muted">
                                    <Upload size={24} style={{ marginBottom: '8px', opacity: 0.5 }} />
                                    <p style={{ margin: 0, fontSize: '0.85rem' }}>Click để tải ảnh thu nhỏ</p>
                                </div>
                            )}
                            <input
                                type="file"
                                id="thumb-upload"
                                hidden
                                accept="image/*"
                                onChange={(e) => e.target.files[0] && setThumbnail(e.target.files[0])}
                            />
                        </div>
                    </div>
                </div>
                <div className="admin-modal-footer">
                    <button className="admin-btn admin-btn-outline" onClick={onClose}>Hủy</button>
                    <button className="admin-btn admin-btn-primary" onClick={handleSaveProduct} disabled={isSaving}>
                        {isSaving ? <Loader className="spin" size={16} /> : 'Lưu sản phẩm'}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ProductModal;
