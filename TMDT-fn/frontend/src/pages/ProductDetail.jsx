import React, { useEffect, useState, useCallback } from 'react';
import { useParams, Link, useNavigate, useLocation } from 'react-router-dom';
import { Star, ShoppingCart, ChevronRight, Package, Loader } from 'lucide-react';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import ProductImageSlider from '../components/Product/ProductImageSlider';
import ProductSpecs from '../components/Product/ProductSpecs';
import ProductReviews from '../components/Product/ProductReviews';
import AuthModal from '../components/Auth/AuthModal';
import { findSpecValue } from '../func/productHelpers';
import toast from 'react-hot-toast';
import '../components/Product/ProductDetail.css';

const API = import.meta.env.VITE_SERVER_API;
const PHOTO_API = (import.meta.env.VITE_PHOTO_SERVER_API || 'http://localhost:8081/images').replace(/\/+$/, '');

const normalizePath = (path, category, id) => {
    if (!path) return '';
    if (path.startsWith('http') || path.startsWith('/')) return path;
    return `/${category}/${id}/${path}`;
};

const ProductDetail = () => {
    const { slug, nameId } = useParams();
    const navigate = useNavigate();
    const location = useLocation();
    const { addToCart } = useCart();
    const { isAuthenticated } = useAuth();

    const [productGroup, setProductGroup] = useState(null);
    const [displayProduct, setDisplayProduct] = useState(null);
    const [selectedColor, setSelectedColor] = useState(null);
    const [loading, setLoading] = useState(true);
    const [showAuthModal, setShowAuthModal] = useState(false);
    const [modalAction, setModalAction] = useState('');
    const [isChecking, setIsChecking] = useState(false);
    const [timeLeft, setTimeLeft] = useState({ h: 23, m: 59, s: 59 });
    const [reviews, setReviews] = useState([]);

    useEffect(() => {
        const timer = setInterval(() => {
            setTimeLeft(prev => {
                if (prev.s > 0) return { ...prev, s: prev.s - 1 };
                if (prev.m > 0) return { ...prev, m: prev.m - 1, s: 59 };
                if (prev.h > 0) return { h: prev.h - 1, m: 59, s: 59 };
                return prev;
            });
        }, 1000);
        return () => clearInterval(timer);
    }, []);

    useEffect(() => {
        if (displayProduct?.product_id) {
            const fetchReviews = async () => {
                try {
                    const res = await fetch(`${API}/api/product/reviews/${displayProduct.product_id}`);
                    const result = await res.json();
                    if (result.success) {
                        setReviews(result.data);
                    }
                } catch (error) {
                    console.error("Error fetching product reviews:", error);
                }
            };
            fetchReviews();
        }
    }, [displayProduct?.product_id]);

    const fetchProduct = useCallback(async (identifier) => {
        setLoading(true);
        try {
            const res = await fetch(`${API}/api/product/product/${identifier}`);
            console.log(`${API}/api/product/product/${identifier}`);
            if (!res.ok) throw new Error(`HTTP Error: ${res.status}`);
            const result = await res.json();

            if (!result.success || !result.data) {
                setDisplayProduct(null);
                setProductGroup([]);
                return;
            }

            const data = result.data;
            const cat = data.category_id || data.category || slug;
            const pid = data.product_id;

            // Normalize paths cho sản phẩm chính
            const normalizeProduct = (prod) => {
                const normPaths = (p) => normalizePath(p, cat, prod.product_id);
                return {
                    ...prod,
                    category: cat,
                    img_thumb: normPaths(prod.img_thumb),
                    local_desc_images: (prod.local_desc_images || []).map(normPaths),
                    variants: (prod.variants || []).map(v => ({
                        ...v,
                        local_gallery: (v.local_gallery || []).map(normPaths),
                        stock: v.quantity ?? v.stock ?? 0
                    }))
                };
            };

            const mainProduct = normalizeProduct(data);

            // Normalize versions (các bản ROM/RAM khác cùng base_id)
            const versions = (data.versions || [data]).map(v => {
                const normV = (p) => normalizePath(p, cat, v.product_id);
                return {
                    ...v,
                    category: cat,
                    img_thumb: normV(v.img_thumb),
                };
            });

            setProductGroup(versions);
            setDisplayProduct(mainProduct);
        } catch (err) {
            console.error('Lỗi fetch product detail:', err);
            setDisplayProduct(null);
            setProductGroup([]);
        } finally {
            setLoading(false);
        }
    }, [slug]);

    useEffect(() => {
        fetchProduct(nameId || slug);
    }, [nameId, slug, fetchProduct]);

    const handleVersionSwitch = async (version) => {
        if (version.product_id === displayProduct?.product_id) return;

        navigate(`/${version.category || slug}/${version.name_id}`, { replace: true });

        try {
            const res = await fetch(`${API}/api/product/product/${version.name_id || version.product_id}`);
            const result = await res.json();
            if (result.success && result.data) {
                const cat = result.data.category_id || slug;
                const normV = (p) => normalizePath(p, cat, result.data.product_id);
                const newDisplay = {
                    ...result.data,
                    category: cat,
                    img_thumb: normV(result.data.img_thumb),
                    local_desc_images: (result.data.local_desc_images || []).map(normV),
                    variants: (result.data.variants || []).map(v => ({
                        ...v,
                        local_gallery: (v.local_gallery || []).map(normV),
                        stock: v.quantity ?? v.stock ?? 0
                    }))
                };
                setDisplayProduct(newDisplay);
            }
        } catch (e) {
            console.error('Lỗi khi chuyển phiên bản:', e);
        }
    };
    useEffect(() => {
        if (displayProduct?.variants?.length > 0) {
            setSelectedColor(prevSelected => {
                if (prevSelected) {
                    const matchedVariant = displayProduct.variants.find(v => v.variant_id === prevSelected.variant_id);
                    if (matchedVariant) return matchedVariant;
                }
                return displayProduct.variants[0];
            });
        } else {
            setSelectedColor(null);
        }
    }, [displayProduct]);

    // ─── Helpers ─────────────────────────────────────────────────────────────
    const formatPrice = (price) => {
        if (!price) return 'Liên hệ';
        return new Intl.NumberFormat('vi-VN').format(price) + '₫';
    };

    const getAvailableStock = () => {
        const target = selectedColor || displayProduct;
        if (!target) return 0;
        const rawQty = parseInt(target.quantity ?? target.stock ?? 0, 10);
        const rsvQty = parseInt(target.reserved ?? 0, 10);
        return Math.max(0, rawQty - rsvQty);
    };

    const handleAddToCart = async (noteOrEvent) => {
        const note = typeof noteOrEvent === 'string' ? noteOrEvent : '';
        if (!isAuthenticated) {
            setModalAction(note);
            setShowAuthModal(true);
            return false;
        }

        const targetVariantId = selectedColor?.variant_id || displayProduct?.variant_id || displayProduct?.id;
        setIsChecking(true);
        try {
            const res = await fetch(`${API}/api/product/check-inventory`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ items: [{ variant_id: targetVariantId, quantity: 1 }] })
            });
            const data = await res.json();

            if (data.success && data.all_in_stock) {
                const options = selectedColor ? {
                    variant_id: selectedColor.variant_id || displayProduct.variant_id,
                    color_name: selectedColor.color_name,
                    color_code: selectedColor.color_code,
                    price: selectedColor.price,
                    image: selectedColor.local_gallery?.[0] || displayProduct.img_thumb,
                    capacity: displayProduct.model || findSpecValue(displayProduct.specs, 'Dung lượng lưu trữ')
                } : {
                    variant_id: displayProduct.variant_id,
                    capacity: displayProduct.model || findSpecValue(displayProduct.specs, 'Dung lượng lưu trữ')
                };
                addToCart(displayProduct, options);
                toast.dismiss('pd-cart');
                toast.success('Đã thêm sản phẩm vào giỏ hàng!', { id: 'pd-cart' });
                return true;
            } else {
                fetchProduct(nameId || slug);
                let msg = 'Sản phẩm đã hết hàng trong kho. Vui lòng chọn sản phẩm khác!';
                if (data.data && Array.isArray(data.data)) {
                    const outOfStock = data.data.filter(i => !i.is_enough)
                        .map(i => `${i.product_name} ${i.color_name ? `(${i.color_name})` : ''}`).join(', ');
                    if (outOfStock) msg = `Tồn kho không đủ. Tồn kho còn lại: ${getAvailableStock()}`;
                }
                toast.dismiss('pd-cart');
                toast.error(msg, { id: 'pd-cart' });
                return false;
            }
        } catch (error) {
            console.error('Lỗi khi kiểm tra tồn kho:', error);
            toast.dismiss('pd-cart');
            toast.error('Có lỗi xảy ra khi kiểm tra tồn kho. Vui lòng thử lại!', { id: 'pd-cart' });
            return false;
        } finally {
            setIsChecking(false);
        }
    };

    const handleBuyNow = async () => {
        const success = await handleAddToCart('mua sản phẩm này');
        if (success) navigate('/cart');
    };

    // ─── States & renders ────────────────────────────────────────────────────
    if (loading) {
        return (
            <div style={{ textAlign: 'center', padding: '80px' }}>
                <Loader className="spin" size={40} style={{ color: 'var(--primary-color)' }} />
                <p style={{ marginTop: '12px', color: '#666' }}>Đang tải thông tin sản phẩm...</p>
            </div>
        );
    }

    if (!displayProduct) {
        return (
            <div className="container" style={{ padding: '50px', textAlign: 'center' }}>
                <h2>Sản phẩm không tồn tại hoặc đã hết hàng</h2>
                <Link to="/" style={{ color: '#2f80ed' }}>Quay lại trang chủ</Link>
            </div>
        );
    }

    const brand = findSpecValue(displayProduct.specs, 'Hãng') || findSpecValue(displayProduct.specs, 'Thương hiệu');
    const brandName = Array.isArray(brand) ? brand[0].replace('.', '') : (brand || 'SmartPhone');
    const categoryName = displayProduct.category_name || (displayProduct.category === 'dtdd' ? 'Điện thoại' : displayProduct.category === 'laptop' ? 'Laptop' : 'Phụ kiện');

    const currentImages = selectedColor?.local_gallery?.length > 0
        ? selectedColor.local_gallery
        : (displayProduct.images || [displayProduct.img_thumb || displayProduct.image]);

    const displayPrice = selectedColor?.price || displayProduct.calculated_price || displayProduct.price;

    // Unique capacities from all versions (RAM/ROM selector)
    const uniqueCapacities = productGroup
        ? Array.from(new Map(productGroup.map(v => {
            const cap = v.model || findSpecValue(v.specs, 'Dung lượng lưu trữ') || v.name_id;
            return [cap, v];
        })).entries())
        : [];

    const getProductRating = () => {
        if (!reviews || reviews.length === 0) {
            return displayProduct?.rating || '4.9';
        }
        const sum = reviews.reduce((acc, curr) => acc + curr.rating, 0);
        return (sum / reviews.length).toFixed(1);
    };

    return (
        <div className="product-detail-page">
            <div className="container">
                <nav className="breadcrumbs">
                    <Link to={`/`}>Trang chủ</Link>
                    <ChevronRight size={12} />
                    <Link to={`/${displayProduct.category}`}>{categoryName}</Link>
                </nav>

                <div className="product-detail__header">
                    <h1 className="product-detail__title">{displayProduct.name}</h1>
                    <div className="product-detail__header-row">
                        <div className="product-detail__rating-header">
                            <div className="header-meta-item">
                                <Star size={16} fill="#fb6e2e" stroke="#fb6e2e" />
                                <span style={{ color: '#fb6e2e', fontWeight: 'bold' }}>{getProductRating()}</span>
                            </div>
                            <span className="rating-count">
                                {displayProduct.sold_count ? `Đã bán ${displayProduct.sold_count}` : 'Đã bán 133,7k'}
                            </span>
                        </div>
                    </div>
                </div>

                <div className="product-detail__main">
                    <div className="product-detail__left product-box">
                        <ProductImageSlider
                            images={currentImages}
                            productName={displayProduct.name}
                        />
                        <div id="specs" className="product-detail__bottom">
                            <ProductSpecs specs={displayProduct.specs} />
                            <ProductReviews reviews={reviews} />
                        </div>
                    </div>

                    <div className="product-detail__right product-box--right">
                        {/* ROM / RAM Selector */}
                        {uniqueCapacities.length > 1 && (
                            <div className="product-options">
                                <div className="option-list">
                                    {uniqueCapacities.map(([cap, version]) => {
                                        const isActive = displayProduct.product_id === version.product_id;
                                        return (
                                            <button
                                                key={version.product_id}
                                                className={`option-btn ${isActive ? 'active' : ''}`}
                                                onClick={() => handleVersionSwitch(version)}
                                            >
                                                {Array.isArray(cap) ? cap[0] : cap}
                                            </button>
                                        );
                                    })}
                                </div>
                            </div>
                        )}

                        {/* Color Selector */}
                        {displayProduct.variants?.length > 0 && (
                            <div className="product-options color-options">
                                <div className="option-list">
                                    {displayProduct.variants.map((v, idx) => {
                                        const isActive = selectedColor?.color_name === v.color_name;
                                        return (
                                            <button
                                                key={idx}
                                                className={`option-btn color-btn ${isActive ? 'active' : ''}`}
                                                onClick={() => setSelectedColor(v)}
                                            >
                                                <span className="color-circle" style={{ backgroundColor: v.color_code }}></span>
                                                {v.color_name}
                                            </button>
                                        );
                                    })}
                                </div>
                            </div>
                        )}

                        {/* Stock info */}
                        <div className="product-stock-info" style={{ marginBottom: '15px', display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.95rem', color: '#666' }}>
                            <Package size={18} style={{ color: 'var(--primary-color)' }} />
                            {getAvailableStock() === 0 ? (
                                <span style={{ color: '#ef4444', fontWeight: '600', padding: '4px 10px', background: 'rgba(239, 68, 68, 0.1)', borderRadius: '8px' }}>
                                    Hết hàng
                                </span>
                            ) : getAvailableStock() <= (displayProduct.low_stock_threshold !== undefined && displayProduct.low_stock_threshold !== null ? parseInt(displayProduct.low_stock_threshold) : 5) ? (
                                <span style={{ color: '#f59e0b', fontWeight: '600', padding: '4px 10px', background: 'rgba(245, 158, 11, 0.1)', borderRadius: '8px', border: '1px solid rgba(245, 158, 11, 0.2)' }}>
                                    Sắp hết hàng (Chỉ còn {getAvailableStock()} sản phẩm!)
                                </span>
                            ) : (
                                <span>Tồn kho: <strong>{getAvailableStock()}</strong> sản phẩm sẵn có</span>
                            )}
                        </div>

                        {/* Price */}
                        <div className="warranty-options">
                            <div className="warranty-item active">
                                <span className="warranty-price">
                                    {selectedColor?.formatted_price || displayProduct.formatted_price || formatPrice(displayPrice)}
                                </span>
                            </div>
                        </div>

                        {/* Actions */}
                        <div className="product-actions">
                            <button
                                className="btn-buy-now"
                                onClick={handleBuyNow}
                                disabled={getAvailableStock() < 1 || isChecking}
                            >
                                <strong>{(getAvailableStock() < 1) ? 'HẾT HÀNG' : (isChecking ? 'ĐANG KIỂM TRA...' : 'MUA NGAY')}</strong>
                            </button>
                            <button
                                className="btn-add-cart"
                                onClick={() => handleAddToCart('thêm sản phẩm vào giỏ hàng')}
                                disabled={getAvailableStock() < 1 || isChecking}
                            >
                                <ShoppingCart size={24} />
                                <span>Thêm vào giỏ</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <AuthModal
                isOpen={showAuthModal}
                onClose={() => setShowAuthModal(false)}
                actionName={modalAction}
                redirectPath="/login"
            />
        </div>
    );
};

export default ProductDetail;
