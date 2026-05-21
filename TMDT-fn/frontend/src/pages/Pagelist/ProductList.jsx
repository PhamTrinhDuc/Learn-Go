import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import FilterBar from '../../components/Product/FilterBar';
import AdvancedFilterModal from '../../components/Product/AdvancedFilterModal';
import ProductGrid from '../../components/Product/ProductGrid';
import './ProductList.css';
import { normalizeProductData } from "../../func/productHelpers.jsx";
import { applyProductFilters } from '../../func/productFilters';

const ProductList = () => {
    const { slug } = useParams();

    const [initialProducts, setInitialProducts] = useState([]);
    const [filteredProducts, setFilteredProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [categoriesList, setCategoriesList] = useState([]);
    const [brandOptions, setBrandOptions] = useState([{ label: 'Tất cả', value: 'all' }]);
    const [isAdvancedFilterOpen, setIsAdvancedFilterOpen] = useState(false);
    const [column,setColumn] = useState();

    const [filters, setFilters] = useState({
        price: 'all',
        brand: 'all',
        sort: 'featured'
    });

    useEffect(() => {
        setFilters({
            price: 'all',
            brand: 'all',
            sort: 'featured'
        });
        setInitialProducts([]);
        setBrandOptions([{ label: 'Tất cả', value: 'all' }]);
        setPage(1);
    }, [slug]);

    useEffect(() => {
        const fetchFilters = async () => {
            try {
                const apiUrl = `${import.meta.env.VITE_SERVER_API}/api/product/filters/${slug}`;
                const response = await fetch(apiUrl);
                if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);

                const result = await response.json();
                const brandSection = result.filters?.brands || result.filters?.brand;
                const rawBrands = brandSection?.options || [];

                if (rawBrands.length > 0) {
                    const mappedBrands = rawBrands.map(brand => {
                        if (typeof brand === 'string') {
                            return { label: brand, value: brand.toLowerCase() };
                        }
                        return {
                            label: brand.label || brand.displayName || brand.name || brand.toString(),
                            value: (brand.value || brand.slug || brand.id || brand.toString()).toString().toLowerCase()
                        };
                    });

                    mappedBrands.sort((a, b) => {
                        const aOther = a.label.toLowerCase() === 'other' || a.label.toLowerCase() === 'khác';
                        const bOther = b.label.toLowerCase() === 'other' || b.label.toLowerCase() === 'khác';
                        if (aOther && !bOther) return 1;
                        if (!aOther && bOther) return -1;
                        return a.label.localeCompare(b.label);
                    });

                    setBrandOptions([{ label: 'Tất cả', value: 'all' }, ...mappedBrands]);
                }
            } catch (err) {
                console.error('Error fetching brand filters:', err);
            }
        };

        if (slug) fetchFilters();
    }, [slug]);

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/category`);
                const result = await response.json();
                if (result.success) {
                    setCategoriesList(result.data);
                }
            } catch (err) {
                console.error('Error fetching categories:', err);
            }
        };
        fetchCategories();
    }, []);

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                //const limit = 10;
                const isAllProducts = slug === 'all';
                const apiUrl = isAllProducts
                    ? `${import.meta.env.VITE_SERVER_API}/api/product/all-products?page=${page}`//&limit=${limit}
                    : `${import.meta.env.VITE_SERVER_API}/api/product/products/${slug}/filter`;

                let fetchOptions = undefined;

                if (!isAllProducts) {
                    const apiFilters = {};
                    Object.entries(filters).forEach(([k, v]) => {
                        if (k === 'sort' || k === 'price' || v === 'all' || !v) return;
                        const values = Array.isArray(v) ? v : [v];
                        if (values.length === 0 || (values.length === 1 && values[0] === 'all')) return;

                        if (k === 'brand') apiFilters.brands = values;
                        else apiFilters[k] = values;
                    });

                    fetchOptions = {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            page,
                            //limit,
                            sort: filters.sort,
                            filters: apiFilters
                        })
                    };
                }

                const response = await fetch(apiUrl, fetchOptions);
                if (!response.ok) throw new Error(`Lỗi kết nối: ${response.status}`);

                const result = await response.json();
                console.log(result);

                if (result.success) {
                    const cleanData = normalizeProductData(result.data);
                    setInitialProducts(cleanData);
                    if (result.pagination) {
                        setTotalPages(result.pagination.totalPages);
                        setColumn(result.pagination.column);
                        // Reset page to the server-returned page just in case
                        if (result.pagination.currentPage && page !== result.pagination.currentPage) {
                            setPage(result.pagination.currentPage);
                        }
                    }
                } else {
                    setInitialProducts([]);
                    setTotalPages(1);
                }
            } catch (e) {
                console.error("Lỗi fetch data:", e);
                setInitialProducts([]);
                setTotalPages(1);
            } finally {
                setLoading(false);
            }
        }

        if (slug) fetchData();
    }, [slug, page, filters]); // Re-fetch from server when specific filters applied

    useEffect(() => {
        const result = applyProductFilters(initialProducts, filters);
        setFilteredProducts(result);
    }, [filters, initialProducts]);

    const handleFilterChange = (type, value) => {
        setFilters(prev => ({ ...prev, [type]: value }));
        setPage(1);
    };

    const handleApplyFilters = (newFilters) => {
        setFilters(newFilters);
        setPage(1);
        setIsAdvancedFilterOpen(false);
    };

    const handleSortChange = (value) => {
        setFilters(prev => ({ ...prev, sort: value }));
        setPage(1);
    };

    const findCategoryBySlug = (cats, targetSlug) => {
        if (!cats) return null;
        for (const cat of cats) {
            if (cat.slug?.toLowerCase() === targetSlug?.toLowerCase()) return cat;
            if (cat.submenu) {
                for (const group of cat.submenu) {
                    if (group.items) {
                        const found = group.items.find(item => item.slug?.toLowerCase() === targetSlug?.toLowerCase());
                        if (found) return found;
                    }
                }
            }
        }
        return null;
    };
    const activeCategory = findCategoryBySlug(categoriesList, slug);

    const getPaginationItems = (currentPage, totalPages) => {
        if (!totalPages || totalPages <= 1) return [1];
        if (totalPages <= 6) return Array.from({ length: totalPages }, (_, i) => i + 1);

        if (currentPage <= 3) return [1, 2, 3, 4, '...', totalPages];
        if (currentPage >= totalPages - 2) return [1, '...', totalPages - 3, totalPages - 2, totalPages - 1, totalPages];

        return [1, '...', currentPage - 1, currentPage, currentPage + 1, '...', totalPages];
    };

    return (
        <div className="product-list-page">
            <div className="container">
                <div className="product-block">
                    {activeCategory && (
                        <h2 style={{ fontSize: '18px', marginBottom: '15px' }}>
                            {activeCategory.label}
                        </h2>
                    )}

                    <FilterBar
                        brandOptions={brandOptions}
                        activeFilters={filters}
                        onFilterChange={handleFilterChange}
                        onSortChange={handleSortChange}
                        onOpenAdvanced={() => setIsAdvancedFilterOpen(true)}
                    />

                    <AdvancedFilterModal
                        isOpen={isAdvancedFilterOpen}
                        onClose={() => setIsAdvancedFilterOpen(false)}
                        slug={slug}
                        activeFilters={filters}
                        onApplyFilters={handleApplyFilters}
                        initialProducts={initialProducts}
                    />

                    {loading ? (
                        <div style={{ textAlign: 'center', padding: '50px' }}>Đang tải dữ liệu...</div>
                    ) : filteredProducts.length === 0 ? (
                        <div style={{ textAlign: 'center', padding: '50px', color: '#666' }}>
                            Không tìm thấy sản phẩm nào phù hợp.
                        </div>
                    ) : (
                        <>
                            <ProductGrid products={filteredProducts} column={column} />

                            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', marginTop: '40px', gap: '8px', flexWrap: 'wrap' }}>
                                <button
                                    disabled={page === 1}
                                    onClick={() => setPage(page - 1)}
                                    style={{ padding: '8px 16px', borderRadius: '4px', border: '1px solid #ddd', background: page === 1 ? '#f5f5f5' : '#fff', cursor: page === 1 ? 'not-allowed' : 'pointer' }}
                                >
                                    Trước
                                </button>

                                {getPaginationItems(page, totalPages).map((p, index) => (
                                    p === '...' ? (
                                        <span key={`ellipsis-${index}`} style={{ padding: '8px 4px', color: '#666' }}>...</span>
                                    ) : (
                                        <button
                                            key={`page-${p}`}
                                            onClick={() => setPage(p)}
                                            style={{
                                                padding: '8px 16px',
                                                borderRadius: '4px',
                                                border: p === page ? '1px solid var(--primary, #2563eb)' : '1px solid #ddd',
                                                background: p === page ? 'var(--primary, #2563eb)' : '#fff',
                                                color: p === page ? '#fff' : '#333',
                                                cursor: 'pointer',
                                                fontWeight: p === page ? 'bold' : 'normal'
                                            }}
                                        >
                                            {p}
                                        </button>
                                    )
                                ))}

                                <button
                                    disabled={page === totalPages || totalPages === 0}
                                    onClick={() => setPage(page + 1)}
                                    style={{ padding: '8px 16px', borderRadius: '4px', border: '1px solid #ddd', background: page === totalPages || totalPages === 0 ? '#f5f5f5' : '#fff', cursor: page === totalPages || totalPages === 0 ? 'not-allowed' : 'pointer' }}
                                >
                                    Sau
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ProductList;