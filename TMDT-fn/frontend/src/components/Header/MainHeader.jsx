import React, { useState, useRef, useEffect } from 'react';
import { Search, ShoppingCart, Menu, User, FileText, ChevronDown, LogOut, Clock, MapPin } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '../../context/CartContext';
import { useAuth } from '../../context/AuthContext';
import './Header.css';

const MainHeader = ({ isMobileMenuOpen, setIsMobileMenuOpen }) => {
    const { getCartCount } = useCart();
    const { user, logout, isAuthenticated } = useAuth();
    const navigate = useNavigate();
    const [isUserDropdownOpen, setIsUserDropdownOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const dropdownRef = useRef(null);
    const cartCount = getCartCount();

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsUserDropdownOpen(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const handleLogout = () => {
        logout();
        setIsUserDropdownOpen(false);
        navigate('/');
    };

    const handleSearch = (e) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
        }
    };

    return (
        <div className="header-main-fixed">
            <div className="container header__container">
                <button
                    className="header__mobile-toggle"
                    onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                >
                    <Menu size={24} />
                </button>

                <div className="header__logo">
                    <Link to="/">
                        <span className="header__logo-text">Thegioibatdong</span>
                    </Link>
                </div>

                <form className="header__search" onSubmit={handleSearch}>
                    <input
                        type="text"
                        className="header__search-input"
                        placeholder="Bạn tìm gì..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                    <button type="submit" className="header__search-btn">
                        <Search size={20} />
                    </button>
                </form>

                <div className="header__actions">
                    {/* <Link to="/history" className="header__action-item hidden-mobile">
                        <FileText size={20} />
                        <span>Lịch sử đơn hàng</span>
                    </Link> */}
                    <Link to="/cart" className="header__action-item">
                        <div className="header__cart-wrapper">
                            <ShoppingCart size={24} />
                            {cartCount > 0 && <span className="header__cart-badge">{cartCount}</span>}
                        </div>
                        <span className="hidden-mobile">Giỏ hàng</span>
                    </Link>

                    {isAuthenticated ? (
                        <div className="header__user-dropdown-container" ref={dropdownRef}>
                            <div
                                className="header__action-item user-active"
                                onClick={() => setIsUserDropdownOpen(!isUserDropdownOpen)}
                            >
                                <div className="user-avatar">
                                    {(user?.full_name || 'U').charAt(0).toUpperCase()}
                                </div>
                                <span className="hidden-mobile">{user?.full_name || 'Người dùng'}</span>
                                <ChevronDown size={14} className={isUserDropdownOpen ? 'rotate' : ''} />
                            </div>

                            {isUserDropdownOpen && (
                                <div className="header__user-dropdown-menu">
                                    <div className="dropdown-user-info">
                                        <p className="user-name">{user?.full_name || 'Người dùng'}</p>
                                        <p className="user-email">{user?.email || ''}</p>
                                    </div>
                                    <div className="dropdown-divider"></div>
                                    <Link to="/profile" className="dropdown-item" onClick={() => setIsUserDropdownOpen(false)}>
                                        <User size={18} />
                                        <span>Thông tin người dùng</span>
                                    </Link>
                                    <Link to="/addresses" className="dropdown-item" onClick={() => setIsUserDropdownOpen(false)}>
                                        <MapPin size={18} />
                                        <span>Địa chỉ của tôi</span>
                                    </Link>
                                    <Link to="/history" className="dropdown-item" onClick={() => setIsUserDropdownOpen(false)}>
                                        <Clock size={18} />
                                        <span>Lịch sử đơn hàng</span>
                                    </Link>
                                    {user?.role === 'admin' && (
                                        <Link to="/admin" className="dropdown-item" onClick={() => setIsUserDropdownOpen(false)}>
                                            <FileText size={18} />
                                            <span>Quản trị (Admin)</span>
                                        </Link>
                                    )}
                                    <div className="dropdown-divider"></div>
                                    <button className="dropdown-item logout-btn" onClick={handleLogout}>
                                        <LogOut size={18} />
                                        <span>Đăng xuất</span>
                                    </button>
                                </div>
                            )}
                        </div>
                    ) : (
                        <Link to="/login" className="header__action-item hidden-mobile">
                            <User size={20} />
                            <span>Đăng nhập</span>
                        </Link>
                    )}
                </div>
            </div>
        </div>
    );
};

export default MainHeader;
