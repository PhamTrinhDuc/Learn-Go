import React from 'react';
import ProductCard from './ProductCard';
import '../../pages/Pagelist/ProductList.css';
const ProductGrid = ({ products , column }) => {
    if (products.length === 0) {
        return (
            <div className="empty-state">
                <p>Không tìm thấy sản phẩm phù hợp</p>
            </div>
        );
    }

    return (
        <div className="product-grid" style={{
            "display": "grid",
            "grid-template-columns": `repeat(${column}, 1fr)`,
            "gap": "10px",
            "margin-top": "20px"
        }}>
            {products.map(product => (
                <div key={product.id} className="product-grid__item">
                    <ProductCard product={product} />
                </div>
            ))}
        </div>
    );
};

export default ProductGrid;
