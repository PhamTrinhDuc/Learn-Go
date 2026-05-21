const parsePriceString = (str) => {
    if (!str) return 0;
    const clean = String(str).replace(/[^\d]/g, '');
    return parseInt(clean) || 0;
};

const extractPrice = (product) => {
    let p = 0;
    if (typeof product.price === 'number' && product.price > 0) {
        p = product.price;
    } else if (product.price) {
        p = parsePriceString(product.price);
    }

    if (p === 0 && product.price_str) {
        p = parsePriceString(product.price_str);
    }

    if (p === 0 && Array.isArray(product.variants) && product.variants.length > 0) {
        const variantPrices = product.variants
            .map(v => {
                if (typeof v.price === 'number' && v.price > 0) return v.price;
                if (v.price) return parsePriceString(v.price);
                if (v.price_str) return parsePriceString(v.price_str);
                return 0;
            })
            .filter(price => price > 0);

        if (variantPrices.length > 0) {
            p = Math.min(...variantPrices);
        }
    }

    if (p === 0 && product.description) {
        const match = product.description.match(/([\d\.]+)(?=\s*?[đdĐ])/);
        if (match) {
            const val = parseInt(match[1].replace(/\./g, ''));
            if (val > 100000) p = val;
        }
    }

    return p;
};

export const findSpecValue = (specs, targetKey) => {
    if (!specs) return null;

    if (specs[targetKey]) return specs[targetKey];

    for (const key in specs) {
        const sub = specs[key];
        if (sub && typeof sub === 'object' && !Array.isArray(sub)) {
            if (sub[targetKey]) return sub[targetKey];
        }
    }
    return null;
};

const extractBrand = (product) => {
    try {
        const getStringValue = (val) => {
            if (!val) return '';
            if (typeof val === 'string') return val;
            if (Array.isArray(val)) return getStringValue(val[0]);
            if (typeof val === 'object') {
                return getStringValue(val.label || val.name || val.displayName || val.value || '');
            }
            return String(val);
        };

        let brand = getStringValue(product.brand);
        if (!brand) {
            const specs = product.specs || {};
            brand = getStringValue(findSpecValue(specs, 'Hãng') || findSpecValue(specs, 'Thương hiệu'));
        }

        if (!brand && product.name) {
            brand = product.name.trim().split(' ')[0];
        }

        if (brand) {
            let cleanBrand = brand
                .replace(/\s*\(.*?\)/g, '')
                .replace(/[.,]/g, '')
                .trim();

            if (cleanBrand.toLowerCase() === 'macbook') return 'Apple';
            return cleanBrand;
        }

        return 'other';
    } catch (e) {
        return 'other';
    }
};

export const normalizeProductData = (rawData) => {
    if (!rawData) return [];
    let items = Array.isArray(rawData) ? rawData : (rawData.data || []);
    if (!Array.isArray(items)) return [];
    items = items.flatMap(item => {
        if (item.products && Array.isArray(item.products)) {
            return item.products.flatMap(p => Array.isArray(p.versions) ? p.versions : [p]);
        }
        if (Array.isArray(item.versions)) {
            return item.versions;
        }
        return [item];
    });

    const groupedMap = new Map();

    const formatPriceDisplay = (price, priceStr) => {
        if (priceStr && priceStr.length > 3) return priceStr;
        if (price > 0) return price.toLocaleString('vi-VN') + '₫';
        return 'Liên hệ';
    };

    items.forEach(product => {
        const id = product.product_id || product.id;
        const category = product.category_id || product.category;
        const price = extractPrice(product);
        const brand = extractBrand(product).toLowerCase();
        const normalizePath = (path) => {
            if (!path) return '';
            if (path.startsWith('http') || path.startsWith('/')) return path;
            return `/${category}/${id}/${path}`;
        };

        let thumb = normalizePath(product.img_thumb);

        const local_desc_images = (product.local_desc_images || []).map(normalizePath);

        const processedVariants = (product.variants || []).map(variant => {
            const vPrice = variant.price || 0;
            const numericVPrice = typeof vPrice === 'number' ? vPrice : parsePriceString(vPrice);
            const rawQty = parseInt(variant.quantity ?? variant.stock ?? variant.amount ?? 0, 10);
            const rsvQty = parseInt(variant.reserved ?? 0, 10);
            return {
                ...variant,
                local_gallery: (variant.local_gallery || []).map(normalizePath),
                stock: Math.max(0, rawQty - rsvQty),
                calculated_price: numericVPrice,
                formatted_price: formatPriceDisplay(numericVPrice, variant.price_str)
            };
        });

        const normalizedItem = {
            ...product,
            id: id,
            category: category,
            calculated_price: price,
            calculated_brand: brand,
            img_thumb: thumb,
            local_desc_images,
            variants: processedVariants,
            formatted_price: formatPriceDisplay(price, product.price_str)
        };

        const uniqueKey = product.base_id || product.name.split(/\d+GB|\d+TB/)[0].trim();

        if (!groupedMap.has(uniqueKey)) {
            groupedMap.set(uniqueKey, {
                main: normalizedItem,
                variants: [normalizedItem]
            });
        } else {
            const entry = groupedMap.get(uniqueKey);
            entry.variants.push(normalizedItem);

            if (price > 0 && (entry.main.calculated_price === 0 || price < entry.main.calculated_price)) {
                entry.main = normalizedItem;
            }
        }
    });

    return Array.from(groupedMap.values()).map(entry => ({
        ...entry.main,
        modelVariants: entry.variants.sort((a, b) => a.calculated_price - b.calculated_price)
    }));
};