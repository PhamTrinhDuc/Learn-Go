import React from 'react';
import { Star } from 'lucide-react';
import './ProductDetail.css';

const ProductReviews = ({ reviews = [] }) => {
    const formatDate = (dateStr) => {
        if (!dateStr) return '';
        return new Date(dateStr).toLocaleDateString('vi-VN', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    return (
        <div className="product-reviews">
            <h3 className="section-title">Đánh giá & Bình luận ({reviews.length})</h3>

            <div className="review-list">
                {reviews.length === 0 ? (
                    <p className="no-reviews" style={{ color: '#777', fontStyle: 'italic', padding: '15px 0' }}>
                        Chưa có đánh giá nào cho sản phẩm này.
                    </p>
                ) : (
                    reviews.map((rev) => (
                        <div key={rev.review_id} className="review-item">
                            <div className="review-item__header">
                                <strong>{rev.full_name || 'Khách hàng'}</strong>
                                <div className="review-item__stars" style={{ display: 'flex', gap: '2px', marginLeft: '10px' }}>
                                    {[1, 2, 3, 4, 5].map(i => (
                                        <Star 
                                            key={i} 
                                            size={14} 
                                            fill={i <= rev.rating ? "var(--color-star)" : "none"} 
                                            stroke={i <= rev.rating ? "var(--color-star)" : "#ccc"} 
                                        />
                                    ))}
                                </div>
                                <span className="review-item__date" style={{ marginLeft: 'auto', fontSize: '12px', color: '#999' }}>
                                    {formatDate(rev.created_at)}
                                </span>
                            </div>
                            <p className="review-item__content">{rev.comment}</p>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

export default ProductReviews;
