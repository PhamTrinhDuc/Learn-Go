import React, { useState } from 'react';
import { Star, PlusCircle } from 'lucide-react';
import { Link } from 'react-router-dom';
import './ProductCard.css';

const ProductCard = ({ product }) => {

    const [selectedVariant, setSelectedVariant] = useState(product);


    const formatPrice = (price) => {
        if (!price) return 'Liên hệ';
        return new Intl.NumberFormat('vi-VN').format(price) + '₫';
    };

    const handleVariantClick = (e, variant) => {
        e.preventDefault();
        e.stopPropagation();
        setSelectedVariant(variant);
    };

    const displayProduct = selectedVariant;
    const hasVariants = product.modelVariants && product.modelVariants.length > 1;

    const findSpec = (targetKey) => {
        const specs = displayProduct.specs;
        if (!specs) return null;
        if (specs[targetKey]) return specs[targetKey];
        for (const key in specs) {
            if (specs[key]?.[targetKey]) return specs[key][targetKey];
        }
        return null;
    };

    return (
        <Link
            to={`/${displayProduct.category || 'dtdd'}/${displayProduct.name_id || displayProduct.id}`}
            state={{ product }}
            className="product-card"
        >
            <div className="product-card__top">
                <span className="label-installment">Trả chậm 0% trả trước 0đ</span>
            </div>

            <div className="product-card__image-container">
                <img
                    src={(displayProduct.img_thumb || '').startsWith('http') ? displayProduct.img_thumb : `${import.meta.env.VITE_PHOTO_SERVER_API}${displayProduct.img_thumb}`}
                    alt={displayProduct.name}
                    className="product-card__image"
                />
            </div>

            <div className="product-card__info">
                <h3 className="product-card__name">{displayProduct.name}</h3>

                <div className="product-card__specs">
                    {(() => {
                        const val = findSpec('Độ phân giải màn hình');
                        if (!val) return null;
                        const text = Array.isArray(val) ? val[0] : val;
                        return (
                            <span className="spec-tag">
                                {text.replace(/\(\d+\s*x\s*\d+\s*[pP]ixels\)/g, '').trim()}
                            </span>
                        );
                    })()}

                    {(() => {
                        const val = findSpec('Màn hình rộng');
                        if (!val) return null;
                        const text = Array.isArray(val) ? val[0] : val;
                        const match = text.match(/(\d+\.?\d*\s*")/);
                        return (
                            <span className="spec-tag">
                                {match ? match[1].trim() : text.split('-')[0].trim()}
                            </span>
                        );
                    })()}
                </div>

                {hasVariants && (
                    <div className="product-card__variants">
                        {product.modelVariants.map((variant) => {
                            const variantLabel = variant.model || (Array.isArray(variant.specs?.['Dung lượng lưu trữ']) ? variant.specs?.['Dung lượng lưu trữ'][0] : variant.specs?.['Dung lượng lưu trữ']) || 'N/A';

                            return (
                                <button
                                    key={variant.id}
                                    className={`variant-btn ${selectedVariant.id === variant.id ? 'active' : ''}`}
                                    onClick={(e) => handleVariantClick(e, variant)}
                                >
                                    {variantLabel}
                                </button>
                            );
                        })}
                    </div>
                )}

                <div className="product-card__pricing">
                    <span className="product-card__price">
                        {displayProduct.formatted_price || formatPrice(displayProduct.calculated_price)}
                    </span>
                    {(displayProduct.old_price || displayProduct.discount_percent) && (
                        <div className="product-card__discount">
                            {displayProduct.old_price && <span className="product-card__old-price">{formatPrice(displayProduct.old_price)}</span>}
                            {displayProduct.discount_percent && <span className="discount-percent">-{displayProduct.discount_percent}%</span>}
                        </div>
                    )}
                </div>

                {displayProduct.promo_text && <div className="product-card__gift">{displayProduct.promo_text}</div>}

                <div className="product-card__footer">
                    <div className="rating-box">
                        <Star size={12} fill="#fb6e2e" stroke="none" />
                        <span>{displayProduct.rating || '4.9'}</span>
                    </div>
                    <div className="compare-box">
                        <span>Đã bán {displayProduct.sold_count || '0'}</span>
                    </div>
                </div>
            </div>
        </Link>
    );
};

export default ProductCard;
