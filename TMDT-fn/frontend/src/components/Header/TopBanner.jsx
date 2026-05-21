import React, { useState, useEffect } from 'react';
import { getBanners } from '../../func/bannerStore';
import './Header.css';

const TopBanner = () => {
    const [config, setConfig] = useState(getBanners().topBanner);

    useEffect(() => {
        const handleStorageChange = () => {
            setConfig(getBanners().topBanner);
        };
        const handleExternalStorageChange = (e) => {
            if (e.key === 'tgbd_banners_config' || !e.key) {
                handleStorageChange();
            }
        };
        window.addEventListener('bannerConfigChanged', handleStorageChange);
        window.addEventListener('storage', handleExternalStorageChange);
        return () => {
            window.removeEventListener('bannerConfigChanged', handleStorageChange);
            window.removeEventListener('storage', handleExternalStorageChange);
        };
    }, []);

    if (!config || !config.active || !config.imageUrl) {
        return null;
    }

    return (
        <a href={config.link || '#'} className="top-banner" style={{ display: 'block' }}>
            <img
                src={config.imageUrl}
                alt="Promotion Banner"
                className="top-banner__img"
            />
        </a>
    );
};

export default TopBanner;