import React, { useState } from 'react';
import { Filter, ChevronDown, Check, X } from 'lucide-react';
import { SORT_TABS, PRICE_SORT_OPTIONS, PRICE_TAB } from '../../data/sortOptions';
import './FilterBar.css';

const FilterBar = ({
    brandOptions = [],
    activeFilters = {},
    onFilterChange,
    onSortChange,
    onOpenAdvanced
}) => {
    const [isPriceOpen, setIsPriceOpen] = useState(false);

    let priceLabel = PRICE_TAB.label;
    if (activeFilters.sort === 'price_asc') priceLabel = 'Giá thấp - cao';
    if (activeFilters.sort === 'price_desc') priceLabel = 'Giá cao - thấp';

    const isPriceActive = activeFilters.sort === 'price_asc' || activeFilters.sort === 'price_desc';

    const activeCount = Object.entries(activeFilters).reduce((acc, [key, val]) => {
        if (key === 'sort') return acc;
        if (val === 'all' || !val) return acc;
        if (Array.isArray(val)) return acc + val.length;
        return acc + 1;
    }, 0);

    const brandValue = activeFilters.brand || 'all';
    const isBrandSelected = (val) => {
        if (brandValue === 'all') return val === 'all';
        if (Array.isArray(brandValue)) return brandValue.some(v => v.toLowerCase() === val.toLowerCase());
        return brandValue.toLowerCase() === val.toLowerCase();
    };

    const hasBrandFilter = brandValue !== 'all' && (Array.isArray(brandValue) ? brandValue.length > 0 : true);

    const removeFilter = (type, value) => {
        const current = activeFilters[type];
        if (Array.isArray(current)) {
            const next = current.filter(v => v !== value);
            onFilterChange(type, next.length === 0 ? 'all' : next);
        } else {
            onFilterChange(type, 'all');
        }
    };

    return (
        <div className="filter-bar">
            <div className="filter-bar__top">
                <button
                    className={`filter-main-btn ${activeCount > 0 ? 'active' : ''}`}
                    onClick={onOpenAdvanced}
                >
                    <div style={{ position: 'relative', display: 'flex', alignItems: 'center' }}>
                        <Filter size={18} />
                        {activeCount > 0 && <span className="filter-badge-dot">{activeCount}</span>}
                    </div>
                    <span>Lọc</span>
                </button>

                <div className="filter-brands-list">
                    {hasBrandFilter ? (
                        <div className="filter-chips-list">
                            {Object.entries(activeFilters).map(([type, value]) => {
                                if (type === 'sort' || value === 'all' || !value) return null;
                                const values = Array.isArray(value) ? value : [value];
                                return values.map(v => {
                                    const label = type === 'brand' ? (brandOptions.find(b => b.value?.toLowerCase() === v.toLowerCase())?.label || v) : v;
                                    return (
                                        <button key={`${type}-${v}`} className="brand-chip selected-view" onClick={() => removeFilter(type, v)}>
                                            {label} <X size={14} className="brand-close-icon" />
                                        </button>
                                    );
                                });
                            })}
                            <button className="clear-all-text" onClick={() => {
                                Object.keys(activeFilters).forEach(key => {
                                    if (key !== 'sort') onFilterChange(key, 'all');
                                });
                            }}>Xóa tất cả</button>
                        </div>
                    ) : (
                        brandOptions.filter(b => b.value?.toLowerCase() !== 'all').map((brand) => (
                            <button
                                key={brand.value}
                                className={`brand-chip ${isBrandSelected(brand.value) ? 'active' : ''}`}
                                onClick={() => {
                                    const isSelected = isBrandSelected(brand.value);
                                    if (isSelected) {
                                        removeFilter('brand', brand.value);
                                    } else {
                                        const current = activeFilters.brand || 'all';
                                        if (Array.isArray(current)) {
                                            onFilterChange('brand', [...current, brand.value]);
                                        } else if (current === 'all') {
                                            onFilterChange('brand', [brand.value]);
                                        } else {
                                            onFilterChange('brand', [current, brand.value]);
                                        }
                                    }
                                }}
                            >
                                {brand.icon ? <img src={brand.icon} alt="" className="brand-logo" /> : brand.label}
                            </button>
                        ))
                    )}
                </div>
            </div>

            <div className="filter-bar__bottom">
                <span className="sort-label">Sắp xếp theo:</span>

                <div className="sort-tabs">
                    {SORT_TABS.map(tab => (
                        <button
                            key={tab.value}
                            className={`sort-tab ${activeFilters.sort === tab.value ? 'active' : ''}`}
                            onClick={() => onSortChange(tab.value)}
                        >
                            {tab.label}
                        </button>
                    ))}

                    <div className="sort-dropdown-wrapper">
                        <button
                            className={`sort-tab sort-dropdown-btn ${isPriceActive ? 'active' : ''}`}
                            onClick={() => setIsPriceOpen(!isPriceOpen)}
                        >
                            <span>{priceLabel}</span>
                            <ChevronDown size={14} />
                        </button>

                        {isPriceOpen && (
                            <div className="sort-dropdown-menu">
                                {PRICE_SORT_OPTIONS.map(option => (
                                    <div
                                        key={option.value}
                                        className="sort-dropdown-item"
                                        onClick={() => {
                                            onSortChange(option.value);
                                            setIsPriceOpen(false);
                                        }}
                                    >
                                        {option.label}
                                        {activeFilters.sort === option.value && <Check size={14} />}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default FilterBar;
