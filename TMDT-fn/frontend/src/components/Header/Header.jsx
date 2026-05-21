import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import TopBanner from './TopBanner';
import CategoryBar from './CategoryBar';
import MainHeader from './MainHeader';
import MobileMenu from './MobileMenu';
import './Header.css';

const Header = () => {
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    return (
        <div className="header-wrapper">
            <TopBanner />
            <MainHeader
                isMobileMenuOpen={isMobileMenuOpen}
                setIsMobileMenuOpen={setIsMobileMenuOpen}
            />
            <CategoryBar />

            <MobileMenu
                isOpen={isMobileMenuOpen}
                onClose={() => setIsMobileMenuOpen(false)}
            />
        </div>
    );
};

export default Header;
