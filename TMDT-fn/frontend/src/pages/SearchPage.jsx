import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import ProductGrid from '../components/Product/ProductGrid';
import { normalizeProductData } from "../func/productHelpers.jsx";
import './Pagelist/ProductList.css';

const SearchPage = () => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [column, setColumn] = useState();
    const location = useLocation();

    const query = new URLSearchParams(location.search).get('q');

    useEffect(() => {
        const fetchSearchResults = async () => {
            if (!query) return;
            setLoading(true);
            try {
                const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/search`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ keyword: query })
                });
                
                const result = await response.json();
                
                if (result.success && result.data) {
                    // API mới trả về cấu trúc grouped (versions) chuẩn
                    setProducts(normalizeProductData(result.data));
                    if (result.column) {
                        setColumn(result.column);
                    }
                } else {
                    setProducts([]);
                }
            } catch (error) {
                console.error("Search error:", error);
                setProducts([]);
            } finally {
                setLoading(false);
            }
        };

        fetchSearchResults();
    }, [query]);

    return (
        <div className="product-list-page" style={{ padding: '40px 0' }}>
            <div className="container">
                <div style={{ marginBottom: '30px' }}>
                    <h1 style={{ fontSize: '24px', fontWeight: 'bold' }}>
                        Kết quả tìm kiếm cho: "{query}"
                    </h1>
                    <p style={{ color: '#666', marginTop: '5px' }}>
                        Tìm thấy {products.length} sản phẩm phù hợp
                    </p>
                </div>

                {loading ? (
                    <div style={{ textAlign: 'center', padding: '100px 0' }}>Đang tìm kiếm...</div>
                ) : products.length > 0 ? (
                    <ProductGrid products={products} column={column} />
                ) : (
                    <div style={{ textAlign: 'center', padding: '100px 0', color: '#666' }}>
                        Không tìm thấy sản phẩm nào phù hợp với từ khóa của bạn.
                    </div>
                )}
            </div>
        </div>
    );
};

export default SearchPage;
