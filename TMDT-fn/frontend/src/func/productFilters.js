export const applyProductFilters = (products, filters) => {
    if (!products || products.length === 0) return [];

    let result = [...products];

    const normalize = (val) => {
        if (!val) return '';
        return String(val)
            .normalize("NFD")
            .replace(/[\u0300-\u036f]/g, "")
            .replace(/đ/g, "d").replace(/Đ/g, "D")
            .toLowerCase()
            .replace(/[^a-z0-9]/g, '');
    };

    Object.entries(filters).forEach(([type, value]) => {
        if (type === 'sort' || type === 'price' || value === 'all' || !value) return;

        const values = Array.isArray(value) ? value : [value];
        if (values.length === 0 || (values.length === 1 && values[0] === 'all')) return;

        const normalizedType = normalize(type);

        result = result.filter(p => {
            if (type === 'brand' || type === 'brands') {
                const pBrand = normalize(p.calculated_brand || p.brand);
                return values.some(v => normalize(v) === pBrand);
            }
            const findSpecValue = (obj, targetNormalized) => {
                if (!obj || typeof obj !== 'object' || Array.isArray(obj)) return null;

                for (const k in obj) {
                    if (normalize(k) === targetNormalized) return obj[k];
                }

                for (const k in obj) {
                    const sub = obj[k];
                    if (sub && typeof sub === 'object' && !Array.isArray(sub)) {
                        const found = findSpecValue(sub, targetNormalized);
                        if (found !== null) return found;
                    }
                }
                return null;
            };

            const specVal = findSpecValue(p.specs || {}, normalizedType) || p[type];

            if (specVal !== undefined && specVal !== null) {
                const getVal = (v) => {
                    if (!v) return '';
                    if (typeof v === 'string') return v;
                    if (Array.isArray(v)) return getVal(v[0]);
                    if (typeof v === 'object') return v.label || v.name || v.displayName || v.value || JSON.stringify(v);
                    return String(v);
                };

                const normalizedSpecVal = normalize(getVal(specVal));
                const match = values.some(v => {
                    const normalizedV = normalize(v);
                    return normalizedSpecVal.includes(normalizedV) || normalizedV.includes(normalizedSpecVal);
                });

                return match;
            }

            return false;
        });
    });
    if (filters.price && filters.price !== 'all') {
        const price = filters.price;
        result = result.filter(p => {
            const pPrice = p.calculated_price;
            if (pPrice === 0) return false;
            switch (price) {
                case '<2': return pPrice < 2000000;
                case '2-4': return pPrice >= 2000000 && pPrice < 4000000;
                case '4-7': return pPrice >= 4000000 && pPrice < 7000000;
                case '7-13': return pPrice >= 7000000 && pPrice < 13000000;
                case '>13': return pPrice >= 13000000;
                default: return true;
            }
        });
    }
    const sortOption = filters.sort;
    if (sortOption === 'price_asc') {
        result.sort((a, b) => a.calculated_price - b.calculated_price);
    } else if (sortOption === 'price_desc') {
        result.sort((a, b) => b.calculated_price - a.calculated_price);
    } else if (sortOption === 'newest') {
        result.sort((a, b) => (b.id || '').localeCompare(a.id || ''));
    }

    return result;
};
