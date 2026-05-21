import React, { useState } from 'react';
import { Star, BarChart2 } from 'lucide-react';
import './ProductDetail.css';

const getImageUrl = (img) => {
    if (!img) return '';
    return img.startsWith('http') ? img : `${import.meta.env.VITE_PHOTO_SERVER_API}${img}`;
};

const ProductImageSlider = ({ images, productName }) => {
    const [mainImage, setMainImage] = useState(images[0]);
    const [activeTab, setActiveTab] = useState('featured');

    React.useEffect(() => {
        setMainImage(images[0]);
        setActiveTab('featured');
    }, [images]);

    return (
        <div className="product-slider">
            <div className="product-slider__main">
                <img src={getImageUrl(mainImage)} alt={productName} />
            </div>

            <div className="product-slider__thumbs-container">
                <div className="product-slider__thumbs">
                    <div
                        className={`product-slider__thumb-item featured ${activeTab === 'featured' ? 'active' : ''}`}
                        onClick={() => {
                            setActiveTab('featured');
                            setMainImage(images[0]);
                        }}
                    >
                        <div className="thumb-featured-icon">
                            <Star size={16} />
                        </div>
                        <span>Nổi bật</span>
                    </div>

                    {images.slice(0, 10).map((img, index) => (
                        <div
                            key={index}
                            className={`product-slider__thumb-item ${activeTab === `img-${index}` ? 'active' : ''}`}
                            onClick={() => {
                                setActiveTab(`img-${index}`);
                                setMainImage(img);
                            }}
                        >
                            <img src={getImageUrl(img)} alt={`${productName} thumbnail ${index}`} />
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default ProductImageSlider;
