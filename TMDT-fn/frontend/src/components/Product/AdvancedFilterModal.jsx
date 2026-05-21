import React, { useState, useEffect, useMemo } from 'react';
import { X, Check } from 'lucide-react';
import { applyProductFilters } from '../../func/productFilters';
import './AdvancedFilterModal.css';

const AdvancedFilterModal = ({
    isOpen,
    onClose,
    slug,
    activeFilters = {},
    onApplyFilters,
    initialProducts = []
}) => {
    const [filterData, setFilterData] = useState(null);
    const [loading, setLoading] = useState(false);
    const [tempFilters, setTempFilters] = useState({});

    useEffect(() => {
        setFilterData(null);
    }, [slug]);

    useEffect(() => {
        if (isOpen) {
            setTempFilters({ ...activeFilters });
            document.body.style.overflow = 'hidden';
            if (!filterData) {
                fetchFilterOptions();
            }
        } else {
            document.body.style.overflow = 'unset';
        }
    }, [isOpen, activeFilters, slug, filterData]);

    const fetchFilterOptions = async () => {
        if (!slug) return;
        setLoading(true);
        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/filters/${slug}`);
            const result = await response.json();
            setFilterData(result);
        } catch (err) {
            console.error('Error fetching filters:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleSelect = (type, value) => {
        const normalizedType = type === 'brands' ? 'brand' : type;
        setTempFilters(prev => {
            const current = prev[normalizedType] || [];
            if (Array.isArray(current)) {
                if (current.includes(value)) {
                    const next = current.filter(v => v !== value);
                    return { ...prev, [normalizedType]: next.length === 0 ? 'all' : next };
                } else {
                    return { ...prev, [normalizedType]: [...current, value] };
                }
            } else {
                if (current === value) return { ...prev, [normalizedType]: 'all' };
                return { ...prev, [normalizedType]: [value] };
            }
        });
    };

    const previewCount = useMemo(() => {
        return applyProductFilters(initialProducts, tempFilters).length;
    }, [initialProducts, tempFilters]);

    const apply = () => {
        onApplyFilters(tempFilters);
    };

    const clearAll = () => {
        const resetFilters = { brand: 'all', sort: activeFilters.sort, price: 'all' };
        setTempFilters(resetFilters);
    };

    if (!isOpen) return null;

    return (
        <div className="filter-modal-overlay" onClick={onClose}>
            <div className="filter-modal" onClick={e => e.stopPropagation()}>
                <div className="filter-modal__header">
                    <h3>Tất cả bộ lọc</h3>
                    <button className="close-btn" onClick={onClose}>
                        <X size={20} /> Đóng
                    </button>
                </div>

                <div className="filter-modal__content">
                    {Object.values(tempFilters).flat().some(v => v !== 'all' && v) && (
                        <div className="modal-section selected-section">
                            <div className="section-label">Đã chọn:</div>
                            <div className="selected-chips">
                                {Object.entries(tempFilters).map(([type, value]) => {
                                    if (type === 'sort' || value === 'all' || !value) return null;
                                    const values = Array.isArray(value) ? value : [value];
                                    return values.map(v => (
                                        <button key={`${type}-${v}`} className="modal-chip active" onClick={() => handleSelect(type, v)}>
                                            {v} <X size={14} />
                                        </button>
                                    ));
                                })}
                                <button className="clear-link" onClick={clearAll}>Xóa tất cả</button>
                            </div>
                        </div>
                    )}

                    {loading ? (
                        <div className="loading-spinner">Đang tải bộ lọc...</div>
                    ) : (
                        filterData?.filters && Object.entries(filterData.filters).map(([key, section]) => {
                            return (
                                <div key={key} className="modal-section">
                                    <h4>{section.displayName || key}</h4>
                                    <div className="options-grid">
                                        {section.options?.map(opt => {
                                            const val = typeof opt === 'string' ? opt : opt.label;
                                            const normalizedType = key === 'brands' ? 'brand' : key;
                                            const current = tempFilters[normalizedType] || [];
                                            const isActive = Array.isArray(current) ? current.includes(val) : current === val;

                                            return (
                                                <button
                                                    key={val}
                                                    className={`option-btn ${isActive ? 'active' : ''}`}
                                                    onClick={() => handleSelect(key, val)}
                                                >
                                                    {isActive && <Check size={14} className="check-icon" />}
                                                    {val}
                                                </button>
                                            );
                                        })}
                                    </div>
                                </div>
                            );
                        })
                    )}
                </div>

                <div className="filter-modal__footer">
                    <button className="reset-btn" onClick={clearAll}>Bỏ chọn</button>
                    <button className="submit-btn" onClick={apply}>
                        Xem {previewCount} kết quả
                    </button>
                </div>
            </div>
        </div>
    );
};

export default AdvancedFilterModal;
