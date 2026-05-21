import React, { useState, useEffect } from 'react';
import { Trash2, Plus, Minus, ChevronDown, Check } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useCart } from '../../context/CartContext';
import '../../pages/Cart/Cart.css';

const CartItem = ({ item, index, isSelected, onSelect, onUpdateQuantity, onUpdateOptions, onRemove }) => {
    const { product, quantity, color } = item;
    const { cartItems, updateQuantity, setQuantity } = useCart();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [tempQuantity, setTempQuantity] = useState(quantity);
    const [liveProduct, setLiveProduct] = useState(product);

    useEffect(() => {
        const fetchLiveStock = async () => {
            try {
                const identifier = product.name_id || product.id || product.product_id;
                if (!identifier) return;
                const res = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/product/${identifier}`);
                const result = await res.json();
                if (result.success && result.data) {
                    setLiveProduct(result.data);
                }
            } catch (e) {
                console.error('Lỗi lấy tồn kho CartItem:', e);
            }
        };
        fetchLiveStock();
    }, [product.name_id, product.id, product.product_id]);

    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(price);
    };

    const currentPrice = item.price || liveProduct.calculated_price || liveProduct.price;
    const currentImage = item.image || liveProduct.img_thumb || liveProduct.image;
    const currentColor = item.color_name || color;

    const variants = liveProduct.variants || [];
    const currentVariant = variants.find(v => v.color_name === currentColor);
    const currentStockRaw = currentVariant ? (currentVariant.quantity ?? currentVariant.stock ?? 0) : (liveProduct.quantity ?? liveProduct.stock ?? 0);
    const currentReservedRaw = currentVariant ? (currentVariant.reserved ?? 0) : (liveProduct.reserved ?? 0);
    const currentStock = Math.max(0, parseInt(currentStockRaw, 10) - parseInt(currentReservedRaw, 10));

    const takenColors = cartItems
        .filter((cartItem, i) =>
            i !== index &&
            cartItem.product.id === product.id &&
            cartItem.capacity === item.capacity
        )
        .map(cartItem => cartItem.color_name);

    const availableVariants = variants.filter(v => !takenColors.includes(v.color_name));

    const handleColorSelect = (v) => {
        onUpdateOptions(index, {
            color_name: v.color_name,
            color_code: v.color_code,
            price: v.price,
            image: v.local_gallery?.[0] || product.img_thumb || product.image
        });
        setIsDropdownOpen(false);
    };

    const handleQuantityBlur = () => {
        let val = parseInt(tempQuantity);
        if (isNaN(val) || val < 1) val = 1;
        if (val > currentStock) val = currentStock;

        setTempQuantity(val);
        setQuantity(index, val);
    };

    const handleQuantityKeyDown = (e) => {
        if (e.key === 'Enter') {
            handleQuantityBlur();
            e.target.blur();
        }
    };

    const handleUpdateQuantity = (delta) => {
        const newQty = quantity + delta;
        if (newQty >= 1 && newQty <= currentStock) {
            updateQuantity(index, delta);
            setTempQuantity(newQty);
        }
    };

    const getImageUrl = (img) => {
        if (!img) return '';
        return img.startsWith('http') ? img : `${import.meta.env.VITE_PHOTO_SERVER_API}${img}`;
    };

    const productPath = `/${product.category}/${product.name_id || product.id}`;

    return (
        <div className={`cart-item ${isSelected ? 'selected' : ''}`}>
            <div className="cart-item__checkbox">
                <label className="checkbox-container">
                    <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() => onSelect()}
                    />
                    <span className="checkmark"></span>
                </label>
            </div>
            <div className="cart-item__image">
                <Link to={productPath} state={{ product }}>
                    <img src={getImageUrl(currentImage)} alt={product.name} />
                </Link>
            </div>
            <div className="cart-item__info">
                <h3 className="cart-item__name">
                    <Link to={productPath} state={{ product }}>{product.name}</Link>
                </h3>
                <div className="cart-item__options">
                    {variants.length > 0 && (
                        <div className="color-selector-wrapper">
                            <button
                                className="color-selector-btn"
                                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                            >
                                {currentColor ? `Màu ${currentColor}` : 'Chọn màu'} <ChevronDown size={14} />
                            </button>

                            {isDropdownOpen && (
                                <div className="color-dropdown">
                                    {availableVariants.map((v, idx) => (
                                        <div
                                            key={idx}
                                            className={`color-dropdown-item ${currentColor === v.color_name ? 'active' : ''}`}
                                            onClick={() => handleColorSelect(v)}
                                        >
                                            <div className="color-dropdown-item__img">
                                                <img src={getImageUrl(v.local_gallery?.[0] || product.img_thumb)} alt={v.color_name} />
                                            </div>
                                            <div className="color-dropdown-item__info">
                                                <span className="color-name">{v.color_name}</span>
                                            </div>
                                            {currentColor === v.color_name && <Check size={14} className="check-icon" />}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}
                    {item.capacity && <span> • {item.capacity}</span>}
                </div>
                <div className="cart-item__price">
                    {formatPrice(currentPrice)}
                </div>
            </div>
            <div className="cart-item__controls">
                <div className="cart-item__stock-status">
                    Tồn kho: {currentStock}
                </div>
                <div className="quantity-control">
                    <button onClick={() => handleUpdateQuantity(-1)} disabled={quantity <= 1}><Minus size={16} /></button>
                    <input
                        type="number"
                        className="quantity-input"
                        value={tempQuantity}
                        onChange={(e) => setTempQuantity(e.target.value)}
                        onBlur={handleQuantityBlur}
                        onKeyDown={handleQuantityKeyDown}
                        min="1"
                        max={currentStock}
                    />
                    <button
                        onClick={() => handleUpdateQuantity(1)}
                        disabled={quantity >= currentStock}
                        title={quantity >= currentStock ? "Hết hàng" : ""}
                    >
                        <Plus size={16} />
                    </button>
                </div>
                <button className="delete-btn" onClick={() => onRemove(index)}>
                    <Trash2 size={18} />
                    <span>Xóa</span>
                </button>
            </div>
            {isDropdownOpen && <div className="dropdown-overlay" onClick={() => setIsDropdownOpen(false)}></div>}
        </div>
    );
};

export default CartItem;
