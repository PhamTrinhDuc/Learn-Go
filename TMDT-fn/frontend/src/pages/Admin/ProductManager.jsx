import React, { useState, useEffect, useCallback } from 'react';
import { Plus, Edit, Trash2, Package, CheckCircle, Search, ChevronLeft, ChevronRight, Loader, FolderPlus, Import as ImportIcon, Image, X, Upload, RefreshCw } from 'lucide-react';
import toast from 'react-hot-toast';
import { normalizeProductData } from '../../func/productHelpers';
import SortIcon from '../../components/Admin/SortIcon';
import CategoryManager from '../../components/Admin/CategoryManager';
import ProductModal from '../../components/Admin/ProductModal';
import ImportStock from '../../components/Admin/ImportStock';
import InventoryModal from '../../components/Admin/InventoryModal';
import '../../components/Admin/ProductForm.css';

const ProductManager = () => {
    const [products, setProducts] = useState([]);
    const [rawProducts, setRawProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showProductModal, setShowProductModal] = useState(false);
    const [showCategoryModal, setShowCategoryModal] = useState(false);
    const [showStockModal, setShowStockModal] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [searchTerm, setSearchTerm] = useState('');
    const [categoryFilter, setCategoryFilter] = useState('Tất cả');
    const [showOutOfStockOnly, setShowOutOfStockOnly] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 8;
    const [categories, setCategories] = useState(['Tất cả']);
    const [categoryMap, setCategoryMap] = useState({});
    const [name, setName] = useState('');
    const [basePrice, setBasePrice] = useState('');
    const [weight, setWeight] = useState('');
    const [basePriceNumeric, setBasePriceNumeric] = useState('');
    const [lowStockThreshold, setLowStockThreshold] = useState('');
    const [productData, setProductData] = useState({
        specs: {},
        variants: []
    });
    const [activeCategory, setActiveCategory] = useState('');
    const [isSaving, setIsSaving] = useState(false);
    const [thumbnail, setThumbnail] = useState(null);
    const [sortConfig, setSortConfig] = useState({ key: null, direction: 'asc' });
    const [showInventoryModal, setShowInventoryModal] = useState(false);
    const [selectedProductForInventory, setSelectedProductForInventory] = useState(null);
    const [editingProductId, setEditingProductId] = useState(null);
    const [editInitialData, setEditInitialData] = useState(null);
    const [alertToast, setAlertToast] = useState({ show: false, message: '', type: 'success' });

    const showToast = (message, type = 'success') => {
        setAlertToast({ show: true, message, type });
        setTimeout(() => setAlertToast({ show: false, message: '', type: 'success' }), 3000);
    };
    const resetForm = () => {
        setName('');
        setBasePrice('');
        setWeight('');
        setBasePriceNumeric('');
        setLowStockThreshold('');
        setActiveCategory('');
        setProductData({ specs: {}, variants: [] });
        setThumbnail(null);
        setIsEditing(false);
        setEditingProductId(null);
        setEditInitialData(null);
        setShowInventoryModal(false);
        setSelectedProductForInventory(null);
    };

    const handleSaveProduct = async () => {
        if (!name || !activeCategory) {
            showToast('Vui lòng nhập tên và chọn danh mục!', 'error');
            return;
        }

        setIsSaving(true);
        try {
            const formData = new FormData();
            if (isEditing && editingProductId) {
                formData.append('id', editingProductId);
            }
            formData.append('name', name);
            formData.append('category', Object.keys(categoryMap).find(key => categoryMap[key] === activeCategory) || activeCategory.toLowerCase());
            formData.append('price', basePrice);
            formData.append('weight', weight);
            formData.append('base_price_numeric', basePriceNumeric);
            formData.append('low_stock_threshold', lowStockThreshold);
            if (thumbnail && thumbnail instanceof File) {
                formData.append('thumbnail', thumbnail);
            }
            formData.append('specs', JSON.stringify(productData.specs));
            formData.append('variants', JSON.stringify(productData.variants.map(v => ({
                id: v.id,
                variant_name: v.variant_name,
                price: v.price,
                stock: v.stock,
                sku: v.sku,
                price_base: v.price_base
            }))));

            productData.variants.forEach((v, vIndex) => {
                if (v.image && v.image instanceof File) {
                    formData.append(`variant_image_${vIndex}`, v.image);
                }
                if (v.gallery && Array.isArray(v.gallery)) {
                    v.gallery.forEach((file, gIndex) => {
                        if (file instanceof File) {
                            formData.append(`variant_${vIndex}_gallery_${gIndex}`, file);
                        }
                    });
                }
            });

            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/save`, {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                showToast('Lưu sản phẩm thành công!');
                setShowProductModal(false);
                resetForm();
                loadData();
            } else {
                showToast('Lỗi khi lưu sản phẩm!', 'error');
            }
        } catch (err) {
            console.error('Error saving product:', err);
            showToast('Lỗi kết nối server!', 'error');
        } finally {
            setIsSaving(false);
        }
    };

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const catRes = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/category`);
            const catJson = await catRes.json();
            let catMap = {};
            let allLabels = ['Tất cả'];

            if (catJson.success && Array.isArray(catJson.data)) {
                catJson.data.forEach(c => {
                    if (!c.hasSubmenu) {
                        allLabels.push(c.label);
                    }
                    catMap[c.slug] = c.label;

                    if (c.hasSubmenu && Array.isArray(c.submenu)) {
                        c.submenu.forEach(subGroup => {
                            if (Array.isArray(subGroup.items)) {
                                subGroup.items.forEach(item => {
                                    if (!allLabels.includes(item.label)) {
                                        allLabels.push(item.label);
                                    }
                                    catMap[item.slug] = item.label;
                                });
                            }
                        });
                    }
                });
                setCategories(allLabels);
                setCategoryMap(catMap);
            }

            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/all-products?limit=2000`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

            const result = await response.json();
            if (result.success) {
                const cleanData = normalizeProductData(result.data);
                setRawProducts(cleanData);
                const flattened = [];

                cleanData.forEach(p => {
                    const models = p.modelVariants || [];
                    if (models.length > 0) {
                        models.forEach(m => {
                            const storage = m.model || m.specs?.['Bộ nhớ trong'] || m.specs?.['Dung lượng'] || '';
                            const ram = m.specs?.['Cấu hình & Bộ nhớ']?.['RAM'] || m.specs?.['RAM'] || '';
                            const modelSpecs = [storage, ram ? `RAM ${ram}` : ''].filter(Boolean).join(' - ');

                            if (m.variants && m.variants.length > 0) {
                                m.variants.forEach((v, idx) => {
                                    const rawQty = parseInt(v.quantity ?? v.stock ?? v.inventory ?? 0, 10);
                                    const reservedQty = parseInt(v.reserved ?? 0, 10);
                                    const stock = Math.max(0, rawQty - reservedQty);
                                    const vName = v.variant_name || v.color_name || '';

                                    let cleanName = p.name;
                                    if (modelSpecs) cleanName = cleanName.replace(modelSpecs, '').trim();
                                    if (vName) cleanName = cleanName.replace(vName, '').trim();
                                    cleanName = cleanName.replace(/[\s-]+$/, '').trim();

                                    const threshold = p.low_stock_threshold !== undefined && p.low_stock_threshold !== null ? parseInt(p.low_stock_threshold) : 5;
                                    let status = 'Còn hàng';
                                    if (stock === 0) status = 'Hết hàng';
                                    else if (stock <= threshold) status = 'Sắp hết hàng';

                                    flattened.push({
                                        id: `${p.id}-${m.id || 'm'}-${v.id || idx}-${vName}`,
                                        productId: p.id,
                                        name: cleanName,
                                        model: modelSpecs || 'N/A',
                                        variant: vName || 'Mặc định',
                                        category: catMap[p.category] || p.category_name || 'Khác',
                                        category_id: p.category,
                                        price: v.price || m.calculated_price || m.price || 0,
                                        stock,
                                        status,
                                        image: p.img_thumb || '/placeholder.png',
                                        originalData: p,
                                        modelData: m
                                    });
                                });
                            } else {
                                const rawQty = parseInt(m.quantity ?? m.stock ?? m.inventory ?? 0, 10);
                                const reservedQty = parseInt(m.reserved ?? 0, 10);
                                const stock = Math.max(0, rawQty - reservedQty);
                                let cleanName = p.name;
                                if (modelSpecs) cleanName = cleanName.replace(modelSpecs, '').trim();
                                cleanName = cleanName.replace(/[\s-]+$/, '').trim();

                                const threshold = p.low_stock_threshold !== undefined && p.low_stock_threshold !== null ? parseInt(p.low_stock_threshold) : 5;
                                let status = 'Còn hàng';
                                if (stock === 0) status = 'Hết hàng';
                                else if (stock <= threshold) status = 'Sắp hết hàng';

                                flattened.push({
                                    id: `${p.id}-${m.id || 'm'}`,
                                    productId: p.id,
                                    name: cleanName,
                                    model: modelSpecs || 'N/A',
                                    variant: 'Mặc định',
                                    category: catMap[p.category] || p.category_name || 'Khác',
                                    category_id: p.category,
                                    price: m.calculated_price || m.price || 0,
                                    stock,
                                    status,
                                    image: p.img_thumb || '/placeholder.png',
                                    originalData: p,
                                    modelData: m
                                });
                            }
                        });
                    } else {
                        const rawQty = parseInt(p.quantity ?? p.stock ?? p.inventory ?? 0, 10);
                        const reservedQty = parseInt(p.reserved ?? 0, 10);
                        const stock = Math.max(0, rawQty - reservedQty);

                        const threshold = p.low_stock_threshold !== undefined && p.low_stock_threshold !== null ? parseInt(p.low_stock_threshold) : 5;
                        let status = 'Còn hàng';
                        if (stock === 0) status = 'Hết hàng';
                        else if (stock <= threshold) status = 'Sắp hết hàng';

                        flattened.push({
                            id: p.id,
                            productId: p.id,
                            name: p.name,
                            model: 'N/A',
                            variant: 'Mặc định',
                            category: catMap[p.category] || p.category_name || 'Khác',
                            category_id: p.category,
                            price: p.calculated_price || p.price || 0,
                            stock,
                            status,
                            image: p.img_thumb || '/placeholder.png',
                            originalData: p
                        });
                    }
                });
                setProducts(flattened);
            }
        } catch (err) {
            console.error('Error fetching admin products:', err);
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const handleDelete = (id) => {
        toast((t) => (
            <div className="confirm-toast">
                <span className="confirm-toast-title">Bạn có chắc chắn muốn xóa sản phẩm này?</span>
                <div className="confirm-toast-actions">
                    <button className="btn-confirm-yes" onClick={() => {
                        toast.dismiss(t.id);
                        performDelete(id);
                    }}>Xóa</button>
                    <button className="btn-confirm-no" onClick={() => toast.dismiss(t.id)}>Hủy</button>
                </div>
            </div>
        ), { duration: 6000, position: 'top-center' });
    };

    const performDelete = async (id) => {
        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/${id}`, {
                method: 'DELETE'
            });
            if (response.ok) {
                toast.success('Xóa sản phẩm thành công!', { id: 'prod-delete' });
                setProducts(products.filter(p => p.id !== id));
            } else {
                toast.error('Lỗi khi xóa sản phẩm!', { id: 'prod-delete' });
            }
        } catch (err) {
            toast.error('Lỗi kết nối server!', { id: 'prod-delete' });
        }
    };

    const handleEditProduct = (p) => {
        const original = p.originalData;

        // Find the specific variant that was clicked in the product table
        const clickedVariantName = p.variant; // e.g. "Đen", "Mặc định"
        const allVariants = (original.variants || []);

        // Match by variant_name or color_name (DB column is color_name)
        let matchingVariant = allVariants.find(v =>
            (v.variant_name || v.color_name || 'Mặc định') === clickedVariantName
        );
        // Fallback: use first variant if no exact match
        if (!matchingVariant && allVariants.length > 0) {
            matchingVariant = allVariants[0];
        }

        // Use variant price if available, otherwise fall back to product price
        const displayPrice = (matchingVariant?.price) || original.price || original.calculated_price || '';

        setName(original.name || '');
        setBasePrice(displayPrice);
        setWeight(original.weight || '');
        setBasePriceNumeric(original.base_price_numeric || '');
        setLowStockThreshold(original.low_stock_threshold !== undefined && original.low_stock_threshold !== null ? original.low_stock_threshold : 5);
        const catLabel = categoryMap[original.category] || original.category_name || p.category;
        setActiveCategory(catLabel);

        // Only pass the one clicked variant so the form shows a single-variant editor
        const variantsForEdit = matchingVariant
            ? [{ ...matchingVariant }]
            : allVariants.slice(0, 1).map(v => ({ ...v }));

        const snapshot = {
            specs: JSON.parse(JSON.stringify(original.specs || {})),
            variants: variantsForEdit
        };
        setEditInitialData(snapshot);
        setProductData(snapshot);
        setThumbnail(original.img_thumb || null);
        setEditingProductId(original.id);
        setIsEditing(true);
        setShowProductModal(true);
    };

    const formatPrice = (price) => {
        if (price === undefined || price === null) return '0đ';
        return price.toLocaleString() + 'đ';
    };

    const handleSort = (key) => {
        let direction = 'asc';
        if (sortConfig.key === key && sortConfig.direction === 'asc') {
            direction = 'desc';
        }
        setSortConfig({ key, direction });
    };

    const sortedProducts = [...products].sort((a, b) => {
        if (!sortConfig.key) return 0;

        let aValue = a[sortConfig.key];
        let bValue = b[sortConfig.key];

        if (sortConfig.key === 'price' || sortConfig.key === 'stock') {
            aValue = Number(aValue) || 0;
            bValue = Number(bValue) || 0;
            const diff = aValue - bValue;
            return sortConfig.direction === 'asc' ? diff : -diff;
        } else {
            aValue = String(aValue || '').toLowerCase();
            bValue = String(bValue || '').toLowerCase();
        }

        if (aValue < bValue) {
            return sortConfig.direction === 'asc' ? -1 : 1;
        }
        if (aValue > bValue) {
            return sortConfig.direction === 'asc' ? 1 : -1;
        }
        return 0;
    });

    const filteredProducts = sortedProducts.filter(p => {
        const name = p.name || '';
        const category = p.category || '';
        const lowerSearch = searchTerm.toLowerCase();

        const matchesSearch = name.toLowerCase().includes(lowerSearch) ||
            category.toLowerCase().includes(lowerSearch);
        const matchesCategory = categoryFilter === 'Tất cả' || category === categoryFilter;
        const matchesStock = !showOutOfStockOnly || p.stock === 0;

        return matchesSearch && matchesCategory && matchesStock;
    });

    const refreshData = async () => {
        try {
            const catRes = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/category`);
            const catJson = await catRes.json();
            if (catJson.success && Array.isArray(catJson.data)) {
                let catMap = {};
                let allLabels = ['Tất cả'];
                catJson.data.forEach(c => {
                    if (!c.hasSubmenu) allLabels.push(c.label);
                    catMap[c.slug] = c.label;
                    if (c.hasSubmenu && Array.isArray(c.submenu)) {
                        c.submenu.forEach(group => group.items?.forEach(item => {
                            if (!allLabels.includes(item.label)) allLabels.push(item.label);
                            catMap[item.slug] = item.label;
                        }));
                    }
                });
                setCategories(allLabels);
                setCategoryMap(catMap);
            }
        } catch (err) {
            console.error('Error refreshing categories:', err);
        }
    };

    const totalPages = Math.ceil(filteredProducts.length / itemsPerPage);
    const indexOfLastItem = currentPage * itemsPerPage;
    const indexOfFirstItem = indexOfLastItem - itemsPerPage;
    const currentItems = filteredProducts.slice(indexOfFirstItem, indexOfLastItem);

    const handlePageChange = (pageNumber) => {
        if (typeof pageNumber === 'number') {
            setCurrentPage(pageNumber);
        }
    };

    const getPaginationRange = () => {
        const delta = 1;
        const range = [];
        for (let i = Math.max(2, currentPage - delta); i <= Math.min(totalPages - 1, currentPage + delta); i++) {
            range.push(i);
        }

        if (currentPage - delta > 2) {
            range.unshift('...');
        }
        if (currentPage + delta < totalPages - 1) {
            range.push('...');
        }

        range.unshift(1);
        if (totalPages > 1) {
            range.push(totalPages);
        }

        return range;
    };

    useEffect(() => {
        setCurrentPage(1);
    }, [searchTerm, categoryFilter, showOutOfStockOnly]);

    return (
        <div className="admin-product-manager">
            <div className="admin-card">
                <div className="admin-card-header" style={{ padding: '16px 24px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', width: '100%', gap: '16px' }}>
                        <h2 style={{ whiteSpace: 'nowrap', margin: 0, fontSize: '1.25rem' }}>Quản lý sản phẩm</h2>

                        <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '6px 12px', background: '#f8fafc', border: '1px solid #e2e8f0', borderRadius: '12px', width: '100%', maxWidth: '800px' }}>
                                <div className="admin-search-wrapper" style={{ position: 'relative', flex: 1 }}>
                                    <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--admin-text-muted)' }} />
                                    <input
                                        type="text"
                                        className="admin-form-input"
                                        placeholder="Tìm sản phẩm..."
                                        style={{ paddingLeft: '40px', width: '100%', height: '36px', borderRadius: '8px', border: 'none', background: 'transparent' }}
                                        value={searchTerm}
                                        onChange={(e) => setSearchTerm(e.target.value)}
                                    />
                                </div>
                                <div style={{ width: '1px', height: '24px', background: '#e2e8f0', margin: '0 4px' }}></div>
                                <select
                                    className="admin-form-input"
                                    style={{ width: '130px', height: '32px', fontSize: '0.85rem', padding: '0 8px', border: 'none', background: 'transparent' }}
                                    value={categoryFilter}
                                    onChange={(e) => setCategoryFilter(e.target.value)}
                                >
                                    {categories.map(cat => (
                                        <option key={cat} value={cat}>{cat === 'Tất cả' ? 'Tất cả loại' : cat}</option>
                                    ))}
                                </select>
                                <div style={{ width: '1px', height: '24px', background: '#e2e8f0', margin: '0 4px' }}></div>
                                <button
                                    className={`admin-btn-sm ${showOutOfStockOnly ? 'admin-btn-primary' : 'admin-btn-outline'}`}
                                    style={{ height: '28px', fontSize: '0.75rem', borderRadius: '6px', whiteSpace: 'nowrap' }}
                                    onClick={() => setShowOutOfStockOnly(!showOutOfStockOnly)}
                                >
                                    Hết hàng
                                </button>
                            </div>
                        </div>

                        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '8px' }}>
                            <button className="admin-btn admin-btn-outline" onClick={() => setShowStockModal(true)} style={{ height: '40px', padding: '0 12px', borderRadius: '8px' }} title="Nhập kho">
                                <ImportIcon size={18} />
                                <span className="hidden-mobile">Nhập kho</span>
                            </button>
                            <button className="admin-btn admin-btn-primary" onClick={() => { resetForm(); setShowProductModal(true); }} style={{ height: '40px', padding: '0 16px', borderRadius: '8px', whiteSpace: 'nowrap' }}>
                                <Plus size={18} />
                                <span>Thêm SP</span>
                            </button>
                            <button className="admin-btn admin-btn-outline" onClick={() => setShowCategoryModal(true)} style={{ height: '40px', padding: '0 12px', borderRadius: '8px' }} title="Thêm loại">
                                <FolderPlus size={18} />
                                <span className="hidden-mobile">Loại</span>
                            </button>
                            <button className="admin-btn admin-btn-outline" onClick={() => loadData()} style={{ height: '40px', padding: '0 12px', borderRadius: '8px' }} title="Làm mới">
                                <RefreshCw size={18} className={loading ? 'spin' : ''} />
                                <span className="hidden-mobile">Làm mới</span>
                            </button>
                        </div>
                    </div>
                </div>

                <div className="admin-table-container">
                    {loading ? (
                        <div style={{ padding: '40px', textAlign: 'center' }}>
                            <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: '12px' }} />
                            <p>Đang tải dữ liệu sản phẩm...</p>
                        </div>
                    ) : (
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th>Hình ảnh</th>
                                    <th onClick={() => handleSort('name')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Tên sản phẩm
                                            <SortIcon activeKey={sortConfig.key} columnKey="name" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('category')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Danh mục
                                            <SortIcon activeKey={sortConfig.key} columnKey="category" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('model')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Model
                                            <SortIcon activeKey={sortConfig.key} columnKey="model" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('variant')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Màu sắc
                                            <SortIcon activeKey={sortConfig.key} columnKey="variant" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('price')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Giá bán
                                            <SortIcon activeKey={sortConfig.key} columnKey="price" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('stock')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Tồn kho
                                            <SortIcon activeKey={sortConfig.key} columnKey="stock" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('status')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Trạng thái
                                            <SortIcon activeKey={sortConfig.key} columnKey="status" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th>Thao tác</th>
                                </tr>
                            </thead>
                            <tbody>
                                {currentItems.map((p) => (
                                    <tr key={p.id}>
                                        <td>
                                            <img
                                                src={(p.image || '').startsWith('http') ? p.image : `${import.meta.env.VITE_PHOTO_SERVER_API}${p.image}`}
                                                alt={p.name}
                                                style={{ width: '32px', height: '32px', objectFit: 'cover', borderRadius: '4px' }}
                                            />
                                        </td>
                                        <td style={{ fontWeight: '500', maxWidth: '200px' }}>
                                            <div className="line-clamp-2" title={p.name}>{p.name || 'N/A'}</div>
                                        </td>
                                        <td>{p.category || 'N/A'}</td>
                                        <td style={{ fontSize: '0.85rem' }}>{p.model || 'N/A'}</td>
                                        <td style={{ fontSize: '0.85rem' }}>{p.variant || 'Mặc định'}</td>
                                        <td>
                                            {formatPrice(p.price)}
                                        </td>
                                        <td>{(p.stock ?? 0) >= 1000 ? '999+' : (p.stock ?? 0)}</td>
                                        <td>
                                            <span 
                                                className={`admin-badge ${p.status === 'Còn hàng' ? 'admin-badge-success' : p.status === 'Sắp hết hàng' ? 'admin-badge-warning' : 'admin-badge-danger'}`}
                                                style={p.status === 'Sắp hết hàng' ? { background: 'rgba(245, 158, 11, 0.1)', color: '#f59e0b', border: '1px solid rgba(245, 158, 11, 0.2)' } : {}}
                                            >
                                                {p.status || 'Hết hàng'}
                                            </span>
                                        </td>
                                        <td>
                                            <div style={{ display: 'flex', gap: '8px' }}>
                                                <button
                                                    className="admin-btn admin-btn-outline admin-btn-sm"
                                                    title="Kiểm tra tồn"
                                                    onClick={() => {
                                                        setSelectedProductForInventory(p);
                                                        setShowInventoryModal(true);
                                                    }}
                                                >
                                                    {CheckCircle ? <CheckCircle size={14} /> : 'OK'}
                                                </button>
                                                <button
                                                    className="admin-btn admin-btn-outline admin-btn-sm"
                                                    onClick={() => handleEditProduct(p)}
                                                >
                                                    {Edit ? <Edit size={14} /> : 'EDIT'}
                                                </button>
                                                <button
                                                    className="admin-btn admin-btn-outline admin-btn-sm"
                                                    style={{ color: 'var(--admin-danger)' }}
                                                    onClick={() => handleDelete(p.productId)}
                                                >
                                                    {Trash2 ? <Trash2 size={14} /> : 'DEL'}
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                                {currentItems.length === 0 && (
                                    <tr>
                                        <td colSpan="9" style={{ textAlign: 'center', padding: '24px' }}>Không tìm thấy sản phẩm nào</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    )}
                </div>

                {
                    totalPages > 1 && (
                        <div className="admin-pagination">
                            <button
                                className="admin-page-btn"
                                disabled={currentPage === 1}
                                onClick={() => handlePageChange(currentPage - 1)}
                            >
                                {ChevronLeft ? <ChevronLeft size={16} /> : '<'}
                            </button>

                            {getPaginationRange().map((page, i) => (
                                <button
                                    key={i}
                                    className={`admin-page-btn ${currentPage === page ? 'active' : ''} ${page === '...' ? 'dots' : ''}`}
                                    onClick={() => handlePageChange(page)}
                                    disabled={page === '...'}
                                    style={page === '...' ? { cursor: 'default', border: 'none' } : {}}
                                >
                                    {page}
                                </button>
                            ))}

                            <button
                                className="admin-page-btn"
                                disabled={currentPage === totalPages}
                                onClick={() => handlePageChange(currentPage + 1)}
                            >
                                {ChevronRight ? <ChevronRight size={16} /> : '>'}
                            </button>
                        </div>
                    )
                }
            </div >

            <ProductModal
                isOpen={showProductModal}
                onClose={() => setShowProductModal(false)}
                isEditing={isEditing}
                name={name}
                setName={setName}
                activeCategory={activeCategory}
                setActiveCategory={setActiveCategory}
                categories={categories}
                basePrice={basePrice}
                setBasePrice={setBasePrice}
                basePriceNumeric={basePriceNumeric}
                setBasePriceNumeric={setBasePriceNumeric}
                lowStockThreshold={lowStockThreshold}
                setLowStockThreshold={setLowStockThreshold}
                weight={weight}
                setWeight={setWeight}
                categoryMap={categoryMap}
                setProductData={setProductData}
                thumbnail={thumbnail}
                setThumbnail={setThumbnail}
                handleSaveProduct={handleSaveProduct}
                isSaving={isSaving}
                initialData={editInitialData}
            />

            <CategoryManager
                isOpen={showCategoryModal}
                onClose={() => setShowCategoryModal(false)}
                onSuccess={refreshData}
            />

            <ImportStock
                isOpen={showStockModal}
                onClose={() => setShowStockModal(false)}
                products={rawProducts}
                categories={categories}
                categoryMap={categoryMap}
                onSuccess={() => { setShowStockModal(false); loadData(); }}
            />

            <InventoryModal
                isOpen={showInventoryModal}
                onClose={() => setShowInventoryModal(false)}
                product={selectedProductForInventory}
                formatPrice={formatPrice}
            />

            {alertToast.show && (
                <div className="toast-container" style={{ position: 'fixed', bottom: '20px', left: '50%', transform: 'translateX(-50%)', zIndex: 9999 }}>
                    <div className={`toast ${alertToast.type}`} style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '12px 24px', borderRadius: '12px', background: alertToast.type === 'success' ? '#10b981' : '#ef4444', color: '#fff', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.1)' }}>
                        <div className="toast-icon">
                            {alertToast.type === 'success' ? <CheckCircle size={20} /> : <X size={20} />}
                        </div>
                        <div className="toast-content" style={{ fontWeight: '500' }}>
                            {alertToast.message}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ProductManager;
