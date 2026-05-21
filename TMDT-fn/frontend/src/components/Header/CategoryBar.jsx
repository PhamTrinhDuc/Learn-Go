import React, { useEffect, useState } from 'react';
import * as LucideIcons from 'lucide-react';
import { ChevronDown } from 'lucide-react';
import SubMenu from './SubMenu';
import './Header.css';
import { Link } from 'react-router-dom';

const CategoryBar = () => {
    const [categoriesList, setCategoriesList] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [activeSubmenu, setActiveSubmenu] = useState(null);

    const closeSubmenu = () => setActiveSubmenu(null);

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const apiUrl = `${import.meta.env.VITE_SERVER_API}/api/product/category`;
                console.log("Fetching:", apiUrl);
                const response = await fetch(apiUrl);
                if (!response.ok) {
                    throw new Error('Failed to fetch categories');
                }
                const result = await response.json();
                if (result.success && Array.isArray(result.data)) {
                    setCategoriesList(result.data);
                } else {
                    throw new Error('Invalid data format received');
                }
            } catch (err) {
                console.error('Error fetching categories:', err);
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchCategories();
    }, []);

    if (loading) {
        return (
            <div className="category-bar">
                <div className="category-bar__container">
                    <div style={{ padding: '10px', fontSize: '13px' }}>Đang tải danh mục...</div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="category-bar">
                <div className="category-bar__container">
                    <div style={{ padding: '10px', fontSize: '13px', color: 'red' }}>Lỗi: {error}</div>
                </div>
            </div>
        );
    }

    return (
        <div className="category-bar">
            <div className="category-bar__container">
                {categoriesList.map(cat => {
                    const IconComponent = LucideIcons[cat.icon] || LucideIcons.HelpCircle;
                    const hasSubmenu = cat.submenu && cat.submenu.length > 0;
                    const linkTo = `/${cat.slug || '#'}`;

                    return (
                        <div
                            key={cat.category_id || cat.id}
                            className={`category-item ${hasSubmenu ? 'has-submenu' : ''}`}
                            onMouseEnter={() => hasSubmenu && setActiveSubmenu(cat.category_id || cat.id)}
                            onMouseLeave={() => setActiveSubmenu(null)}
                        >
                            <Link to={linkTo} className="category-item__label" style={{ textDecoration: 'none', color: 'inherit' }}>
                                <IconComponent size={18} />
                                <span>{cat.label}</span>
                                {hasSubmenu && <ChevronDown size={14} className="category-arrow" />}
                            </Link>
                            {hasSubmenu && activeSubmenu === cat.id && (
                                <SubMenu groups={cat.submenu} onClose={closeSubmenu} />
                            )}
                        </div>
                    );
                })}
            </div>
        </div>
    );
};

export default CategoryBar;
