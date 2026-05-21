import React, { useState, useEffect, useRef } from 'react';
import { Plus, X, Tag, Image, Package, Trash2 } from 'lucide-react';
import './ProductForm.css';

// Search for a spec value at top level OR one level deep inside nested groups
// Checks an array of possible keys (like [field.key, field.label])
const findSpecDetails = (specs, possibleKeys) => {
    if (!specs) return null;
    
    // 1. Exact match first
    for (const targetKey of possibleKeys) {
        if (!targetKey) continue;
        // Check top level
        if (specs[targetKey] !== undefined && specs[targetKey] !== null) {
            return { value: specs[targetKey], originalKey: targetKey, groupName: null };
        }
        // Check nested groups
        for (const groupKey in specs) {
            const group = specs[groupKey];
            if (group && typeof group === 'object' && !Array.isArray(group)) {
                if (group[targetKey] !== undefined && group[targetKey] !== null) {
                    return { value: group[targetKey], originalKey: targetKey, groupName: groupKey };
                }
            }
        }
    }

    // 2. Fuzzy match (substring or synonym)
    const synonyms = {
        'bộ nhớ trong': ['dung lượng lưu trữ', 'rom'],
        'camera sau': ['độ phân giải camera sau', 'camera chính'],
        'camera trước': ['độ phân giải camera trước', 'camera selfie'],
        'màn hình': ['màn hình rộng', 'kích thước màn hình'],
        'pin & sạc': ['dung lượng pin'],
        'chip xử lý (cpu)': ['cpu', 'chip xử lý'],
        'hệ điều hành': ['os']
    };

    for (const targetKey of possibleKeys) {
        if (!targetKey) continue;
        const t = targetKey.toLowerCase();
        
        const searchTerms = [t];
        if (synonyms[t]) searchTerms.push(...synonyms[t]);

        // Check top level
        for (const k of Object.keys(specs)) {
            const kl = k.toLowerCase();
            if (searchTerms.some(term => kl.includes(term) || term.includes(kl))) {
                return { value: specs[k], originalKey: k, groupName: null };
            }
        }
        
        // Check nested groups
        for (const groupKey in specs) {
            const group = specs[groupKey];
            if (group && typeof group === 'object' && !Array.isArray(group)) {
                for (const k of Object.keys(group)) {
                    const kl = k.toLowerCase();
                    if (searchTerms.some(term => kl.includes(term) || term.includes(kl))) {
                        return { value: group[k], originalKey: k, groupName: groupKey };
                    }
                }
            }
        }
    }

    return null;
};

// Flatten objects/arrays to a string to prevent [object Object] in inputs
const formatSpecValue = (val) => {
    if (val === null || val === undefined) return '';
    if (Array.isArray(val)) return val.join(', ');
    if (typeof val === 'object') {
        return Object.values(val).filter(v => typeof v !== 'object').join(', ');
    }
    return String(val);
};

// Resolve an image to a displayable src — handles File objects AND server URL strings
const getImgSrc = (img) => {
    if (!img) return null;
    if (img instanceof File || img instanceof Blob) return URL.createObjectURL(img);
    if (typeof img === 'string') {
        return img.startsWith('http') ? img : `${import.meta.env.VITE_PHOTO_SERVER_API}${img}`;
    }
    return null;
};

const ProductForm = ({ categorySlug, onChange, initialData, isEditing }) => {
    console.log('ProductForm Render - category:', categorySlug, 'initialData:', initialData);
    const [templates, setTemplates] = useState({});
    const [specifications, setSpecifications] = useState(initialData?.specs || {});
    const [customFields, setCustomFields] = useState([]);
    const [variants, setVariants] = useState(() => {
        return (initialData?.variants || []).map((v, i) => ({
            ...v,
            id: v.id || `v-${Date.now()}-${i}`
        }));
    });

    const [tagInputs, setTagInputs] = useState({});
    const [loading, setLoading] = useState(true);
    const prevSlugRef = useRef(categorySlug);
    const isFirstMount = useRef(true);

    const inputRef = useRef({});
    const fieldToGroupRef = useRef({});

    useEffect(() => {
        const fetchTemplates = async () => {
            try {
                const res = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/form-templates`);
                const data = await res.json();
                setTemplates(data);
            } catch (err) {
                console.error('Error fetching form templates:', err);
            } finally {
                setLoading(false);
            }
        };
        fetchTemplates();
    }, [categorySlug]);

    useEffect(() => {
        if (!isEditing || !initialData || Object.keys(templates).length === 0) return;

        const template = templates[categorySlug] || { fields: [] };
        const templateKeys = template.fields.map(f => f.key);

        // Build the group map to remember which category/group each spec belongs to
        const fieldToGroup = {};
        Object.entries(initialData.specs || {}).forEach(([groupName, groupValue]) => {
            if (groupValue && typeof groupValue === 'object' && !Array.isArray(groupValue)) {
                Object.keys(groupValue).forEach(fieldKey => {
                    fieldToGroup[fieldKey] = groupName;
                });
            }
        });
        fieldToGroupRef.current = fieldToGroup;

        const fieldToOriginalKey = {}; // Maps template key -> original key from DB

        // For each template field, look up its value from the (possibly nested) specs
        // We check both the english key AND the vietnamese label
        const matchedSpecs = {};
        template.fields.forEach(f => {
            const match = findSpecDetails(initialData.specs, [f.key, f.label]);
            if (match) {
                matchedSpecs[f.key] = formatSpecValue(match.value);
                fieldToOriginalKey[f.key] = match.originalKey;
                if (match.groupName) {
                    fieldToGroup[f.key] = match.groupName;
                }
            }
        });
        
        // Save the reverse mapping for reconstruction on save
        inputRef.current._keyMapping = fieldToOriginalKey;

        // === DEBUG ===
        console.log('[ProductForm] categorySlug:', categorySlug);
        console.log('[ProductForm] initialData.specs:', JSON.stringify(initialData.specs));
        console.log('[ProductForm] templateKeys:', templateKeys);
        console.log('[ProductForm] matchedSpecs:', matchedSpecs);
        // =============

        // Flatten entire specs tree and collect non-template keys as custom fields
        const flatSpecs = {};
        const recurse = (obj) => {
            Object.entries(obj || {}).forEach(([key, val]) => {
                if (val && typeof val === 'object' && !Array.isArray(val)) {
                    recurse(val);
                } else {
                    flatSpecs[key] = val;
                }
            });
        };
        recurse(initialData.specs || {});

        const extraFields = [];
        Object.entries(flatSpecs).forEach(([key, value]) => {
            if (!templateKeys.includes(key)) {
                // Ensure value is always a string for <input> compatibility
                const strValue = (value === null || value === undefined)
                    ? ''
                    : (typeof value === 'string' ? value : JSON.stringify(value));
                extraFields.push({ id: `cf-${Date.now()}-${key}`, key, value: strValue });
            }
        });

        console.log('[ProductForm] extraFields:', extraFields.map(f => f.key));

        setSpecifications(matchedSpecs);
        setCustomFields(extraFields);

        if (initialData.variants?.length > 0) {
            setVariants(initialData.variants.map((v, i) => ({
                ...v,
                id: v.variant_id || v.id || `v-${Date.now()}-${i}`,
                variant_name: v.variant_name || v.color_name || '',
                price: v.price ?? '',
                stock: v.quantity ?? v.stock ?? v.amount ?? 0,
                sku: v.sku || '',
                price_base: v.price_base ?? 0
            })));
        }
    }, [templates, categorySlug, initialData, isEditing]);

    useEffect(() => {
        if (!isEditing && prevSlugRef.current !== categorySlug && !isFirstMount.current) {
            setSpecifications({});
            setCustomFields([]);
            setTagInputs({});
            setVariants([]);
        }
        prevSlugRef.current = categorySlug;
        isFirstMount.current = false;
    }, [categorySlug, isEditing]);

    useEffect(() => {
        const finalSpecs = {};
        const fieldToGroup = fieldToGroupRef.current || {};
        const keyMapping = inputRef.current._keyMapping || {};

        // Reconstruct specs grouped by category/section
        Object.entries(specifications).forEach(([key, value]) => {
            const groupName = fieldToGroup[key] || 'Thông số kỹ thuật';
            const originalKey = keyMapping[key] || key; // Use DB key if we mapped it (e.g. "Hệ điều hành")
            if (!finalSpecs[groupName]) {
                finalSpecs[groupName] = {};
            }
            finalSpecs[groupName][originalKey] = value;
        });

        customFields.forEach(field => {
            if (field.key && field.value) {
                const groupName = fieldToGroup[field.key] || 'Thông số kỹ thuật';
                if (!finalSpecs[groupName]) {
                    finalSpecs[groupName] = {};
                }
                finalSpecs[groupName][field.key] = field.value;
            }
        });

        if (onChange) {
            onChange({
                specs: finalSpecs,
                variants: variants
            });
        }
    }, [specifications, customFields, variants]);

    const handleTextChange = (key, value) => {
        setSpecifications(prev => ({
            ...prev,
            [key]: value
        }));
    };

    const handleAddTag = (key, e) => {
        if (e.key === 'Enter' || e.key === ',') {
            e.preventDefault();
            const value = tagInputs[key]?.trim();
            if (value) {
                const currentTags = specifications[key] || [];
                if (!currentTags.includes(value)) {
                    setSpecifications(prev => ({
                        ...prev,
                        [key]: [...currentTags, value]
                    }));
                }
                setTagInputs(prev => ({ ...prev, [key]: '' }));
            }
        }
    };

    const removeTag = (key, tagToRemove) => {
        setSpecifications(prev => ({
            ...prev,
            [key]: (prev[key] || []).filter(t => t !== tagToRemove)
        }));
    };

    const addCustomField = () => {
        setCustomFields(prev => [...prev, { id: Date.now(), key: '', value: '' }]);
    };

    const handleCustomFieldChange = (index, type, value) => {
        setCustomFields(prev => prev.map((field, idx) =>
            idx === index ? { ...field, [type]: value } : field
        ));
    };

    const removeCustomField = (index) => {
        setCustomFields(prev => prev.filter((_, i) => i !== index));
    };

    const addVariant = () => {
        setVariants(prev => [...prev, {
            id: `v-${Date.now()}-${prev.length}`,
            variant_name: '',
            price: '',
            stock: '',
            sku: '',
            price_base: 0,
            image: null,
            gallery: []
        }]);
    };

    const handleVariantChange = (index, field, value) => {
        setVariants(prev => prev.map((v, idx) =>
            idx === index ? { ...v, [field]: value } : v
        ));
    };

    const handleVariantImageChange = (index, file) => {
        setVariants(prev => prev.map((v, idx) =>
            idx === index ? { ...v, image: file } : v
        ));
    };

    const handleVariantGalleryAdd = (index, files) => {
        setVariants(prev => prev.map((v, idx) =>
            idx === index ? { ...v, gallery: [...(v.gallery || []), ...Array.from(files)] } : v
        ));
    };

    const removeVariantGalleryImage = (variantIndex, imageIndex) => {
        setVariants(prev => prev.map((v, idx) =>
            idx === variantIndex ? { ...v, gallery: (v.gallery || []).filter((_, i) => i !== imageIndex) } : v
        ));
    };

    const removeVariant = (index) => {
        setVariants(prev => prev.filter((_, i) => i !== index));
    };

    if (loading) return <div className="p-4 text-center">Đang tải biểu mẫu...</div>;

    const currentTemplate = templates[categorySlug] || { fields: [] };

    return (
        <div className="product-form-container">
            <div className="product-form-section">
                <div className="section-header">
                    <h3 className="section-title">Thông số kỹ thuật</h3>
                </div>
                <div className="product-form-grid">
                    {currentTemplate.fields.map((field) => (
                        <div key={field.key} className="form-field-group">
                            <label className="form-field-label">{field.label}</label>

                            {field.type === 'multi' ? (
                                <div className="tag-input-container" onClick={() => inputRef.current[field.key]?.focus()}>
                                    {Array.isArray(specifications[field.key]) && specifications[field.key].map(tag => (
                                        <span key={tag} className="tag-item">
                                            {tag}
                                            <X size={14} className="tag-remove" onClick={(e) => {
                                                e.stopPropagation();
                                                removeTag(field.key, tag);
                                            }} />
                                        </span>
                                    ))}
                                    <input
                                        id={`spec-${field.key}`}
                                        ref={el => inputRef.current[field.key] = el}
                                        type="text"
                                        className="tag-input-field"
                                        placeholder={(!specifications[field.key] || specifications[field.key].length === 0) ? "Nhấn Enter hoặc ',' để thêm..." : ""}
                                        value={tagInputs[field.key] || ''}
                                        onChange={(e) => setTagInputs(prev => ({ ...prev, [field.key]: e.target.value }))}
                                        onKeyDown={(e) => handleAddTag(field.key, e)}
                                    />
                                </div>
                            ) : (
                                <input
                                    id={`spec-${field.key}`}
                                    type="text"
                                    className="form-field-input"
                                    placeholder={`Nhập ${field.label.toLowerCase()}...`}
                                    value={specifications[field.key] || ''}
                                    onChange={(e) => handleTextChange(field.key, e.target.value)}
                                />
                            )}
                        </div>
                    ))}
                </div>

                <div className="custom-params-section">
                    <div className="custom-params-header">
                        <label className="form-field-label" style={{ margin: 0 }}>Thông số bổ sung</label>
                        <button type="button" className="btn-add-param" onClick={addCustomField}>
                            <Plus size={14} />
                            Thêm dòng
                        </button>
                    </div>

                    {customFields.map((field, index) => (
                        <div key={field.id || index} className="custom-field-row">
                            <input
                                type="text"
                                className="form-field-input"
                                style={{ flex: 1 }}
                                placeholder="Tên thông số"
                                value={field.key}
                                onChange={(e) => handleCustomFieldChange(index, 'key', e.target.value)}
                            />
                            <input
                                type="text"
                                className="form-field-input"
                                style={{ flex: 1 }}
                                placeholder="Giá trị"
                                value={field.value}
                                onChange={(e) => handleCustomFieldChange(index, 'value', e.target.value)}
                            />
                            <button type="button" className="btn-remove-param" onClick={() => removeCustomField(index)}>
                                <X size={16} />
                            </button>
                        </div>
                    ))}
                </div>
            </div>

            <div className="product-form-section" style={{ borderTop: '1px solid #f1f5f9', marginTop: '24px', paddingTop: '24px' }}>
                <div className="section-header">
                    <div>
                        <h3 className="section-title">
                            {isEditing ? 'Thông tin biến thể' : 'Phân loại & Biến thể'}
                        </h3>
                        <p className="section-subtitle" style={{ margin: 0 }}>
                            {isEditing
                                ? 'Màu sắc, giá bán và tồn kho của biến thể này'
                                : 'Màu sắc, dung lượng, hình ảnh riêng cho từng loại...'}
                        </p>
                    </div>
                    {!isEditing && (
                        <button type="button" className="btn-add-param" onClick={addVariant}>
                            <Plus size={14} />
                            Thêm biến thể
                        </button>
                    )}
                </div>

                <div className="variants-wrapper">
                    {variants.length === 0 ? (
                        <div className="empty-variants">
                            <Package size={32} style={{ opacity: 0.2, marginBottom: '8px' }} />
                            <p>Sản phẩm này chưa có biến thể nào.</p>
                        </div>
                    ) : (
                        variants.map((variant, index) => (
                            <div key={variant.id ? `variant-${variant.id}` : `idx-${index}`} className="variant-card">
                                <div className="variant-card-main">
                                    <div className="variant-thumb-container">
                                        <div
                                            className="variant-image-uploader-big"
                                            onClick={() => document.getElementById(`variant-image-${variant.id}`).click()}
                                        >
                                            {variant.image ? (
                                                <img src={getImgSrc(variant.image)} alt="V" />
                                            ) : (variant.local_gallery?.[0]) ? (
                                                <img src={getImgSrc(variant.local_gallery[0])} alt="V" />
                                            ) : (
                                                <div className="upload-placeholder">
                                                    <Image size={24} />
                                                    <span>Ảnh đại diện</span>
                                                </div>
                                            )}
                                            <input
                                                type="file"
                                                id={`variant-image-${variant.id}`}
                                                hidden
                                                accept="image/*"
                                                onChange={(e) => {
                                                    if (e.target.files[0]) handleVariantImageChange(index, e.target.files[0]);
                                                }}
                                            />
                                        </div>
                                    </div>

                                    <div className="variant-inputs-grid">
                                        <div className="form-field-group">
                                            <label className="admin-form-label-sm">Tên biến thể</label>
                                            <input
                                                type="text"
                                                className="form-field-input"
                                                placeholder="Đen, 128GB..."
                                                value={variant.variant_name || ''}
                                                onChange={(e) => handleVariantChange(index, 'variant_name', e.target.value)}
                                            />
                                        </div>
                                        <div className="form-field-group">
                                            <label className="admin-form-label-sm">Giá bán (đ)</label>
                                            <input
                                                type="number"
                                                className="form-field-input"
                                                placeholder="Giá riêng"
                                                value={variant.price ?? ''}
                                                onChange={(e) => handleVariantChange(index, 'price', e.target.value)}
                                            />
                                        </div>
                                        <div className="form-field-group">
                                            <label className="admin-form-label-sm">Giá nhập (đ)</label>
                                            <input
                                                type="number"
                                                className="form-field-input"
                                                placeholder="0đ"
                                                value={variant.price_base ?? 0}
                                                disabled={true}
                                                style={{ backgroundColor: '#f1f5f9', cursor: 'not-allowed', color: '#64748b' }}
                                            />
                                        </div>
                                        <div className="form-field-group">
                                            <label className="admin-form-label-sm">Tồn kho</label>
                                            <input
                                                type="number"
                                                className="form-field-input"
                                                placeholder="0"
                                                value={variant.stock ?? ''}
                                                onChange={(e) => handleVariantChange(index, 'stock', e.target.value)}
                                            />
                                        </div>
                                        <div className="form-field-group">
                                            <label className="admin-form-label-sm">SKU</label>
                                            <input
                                                type="text"
                                                className="form-field-input"
                                                placeholder="Mã SKU"
                                                value={variant.sku || ''}
                                                onChange={(e) => handleVariantChange(index, 'sku', e.target.value)}
                                            />
                                        </div>
                                    </div>

                                    {!(isEditing && variants.length === 1) && (
                                        <div className="variant-actions">
                                            <button type="button" className="btn-remove-variant" onClick={() => removeVariant(index)}>
                                                <Trash2 size={18} />
                                            </button>
                                        </div>
                                    )}
                                </div>

                                <div className="variant-gallery-section">
                                    <div className="variant-gallery-header">
                                        <span className="admin-form-label-sm">Bộ sưu tập ảnh cho biến thể này</span>
                                    </div>
                                    <div className="variant-gallery-grid">
                                        {/* Existing server images (read-only preview) */}
                                        {(variant.local_gallery || []).map((imgPath, idx) => (
                                            <div key={`srv-${idx}`} className="gallery-item-mini">
                                                <img src={getImgSrc(imgPath)} alt="gallery" />
                                            </div>
                                        ))}
                                        {/* Newly uploaded files (can be removed) */}
                                        {(variant.gallery || []).map((file, idx) => (
                                            <div key={`new-${idx}`} className="gallery-item-mini">
                                                <img src={getImgSrc(file)} alt="gallery" />
                                                <button
                                                    type="button"
                                                    className="btn-remove-gallery-img"
                                                    onClick={() => removeVariantGalleryImage(index, idx)}
                                                >
                                                    <X size={12} />
                                                </button>
                                            </div>
                                        ))}
                                        <div
                                            className="btn-add-gallery-mini"
                                            onClick={() => document.getElementById(`variant-gallery-${variant.id}`).click()}
                                        >
                                            <Plus size={20} />
                                            <span>Thêm ảnh</span>
                                            <input
                                                type="file"
                                                id={`variant-gallery-${variant.id}`}
                                                hidden
                                                multiple
                                                accept="image/*"
                                                onChange={(e) => handleVariantGalleryAdd(index, e.target.files)}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>
        </div>
    );
};

export default ProductForm;
