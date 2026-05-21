import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { X, User, FileText, MapPin, Clock, LogOut, ChevronRight, Settings } from 'lucide-react';
import * as LucideIcons from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import './Header.css';

const MobileMenu = ({ isOpen, onClose }) => {
    const { user, isAuthenticated, logout } = useAuth();
    const [categories, setCategories] = useState([]);

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/product/category`);
                const result = await response.json();
                if (result.success) setCategories(result.data);
            } catch (err) {
                console.error(err);
            }
        };
        if (isOpen) fetchCategories();
    }, [isOpen]);

    if (!isOpen) return null;

    return (
        <div className={`mobile-menu-overlay ${isOpen ? 'show' : ''}`} onClick={onClose}>
            <div className={`mobile-menu-drawer ${isOpen ? 'slide-in' : ''}`} onClick={e => e.stopPropagation()}>
                <div className="mobile-menu-header">
                    <div className="mobile-menu-logo">Thegioibatdong</div>
                    <button className="mobile-menu-close" onClick={onClose}>
                        <X size={24} />
                    </button>
                </div>

                <div className="mobile-menu-body">
                    {isAuthenticated ? (
                        <div className="mobile-user-section">
                            <div className="mobile-user-info">
                                <div className="mobile-user-avatar">
                                    {(user?.full_name || 'U').charAt(0).toUpperCase()}
                                </div>
                                <div className="mobile-user-details">
                                    <p className="mobile-user-name">{user?.full_name}</p>
                                    <p className="mobile-user-email">{user?.email}</p>
                                </div>
                            </div>
                            <div className="mobile-user-actions">
                                <Link to="/profile" onClick={onClose} className="mobile-action-item">
                                    <User size={18} />
                                    <span>Tài khoản của tôi</span>
                                </Link>
                                <Link to="/addresses" onClick={onClose} className="mobile-action-item">
                                    <MapPin size={18} />
                                    <span>Sổ địa chỉ</span>
                                </Link>
                                <Link to="/history" onClick={onClose} className="mobile-action-item">
                                    <Clock size={18} />
                                    <span>Lịch sử đơn hàng</span>
                                </Link>
                                {user?.role === 'admin' && (
                                    <Link to="/admin" onClick={onClose} className="mobile-action-item admin-link">
                                        <Settings size={18} />
                                        <span>Trang quản trị</span>
                                    </Link>
                                )}
                            </div>
                        </div>
                    ) : (
                        <div className="mobile-login-section">
                            <p>Đăng nhập để nhận nhiều ưu đãi!</p>
                            <Link to="/login" onClick={onClose} className="mobile-btn-login">Đăng nhập</Link>
                            <Link to="/register" onClick={onClose} className="mobile-btn-register">Đăng ký thành viên</Link>
                        </div>
                    )}

                    <div className="mobile-section-divider"></div>

                    <div className="mobile-categories">
                        <h4 className="section-title">Danh mục sản phẩm</h4>
                        <div className="category-list">
                            {categories.map(cat => {
                                const IconComp = LucideIcons[cat.icon] || LucideIcons.Package;
                                return (
                                    <Link key={cat.id} to={`/${cat.slug}`} onClick={onClose} className="mobile-cat-item">
                                        <div className="cat-icon-wrapper">
                                            <IconComp size={20} />
                                        </div>
                                        <span>{cat.label}</span>
                                        <ChevronRight size={16} />
                                    </Link>
                                );
                            })}
                        </div>
                    </div>
                </div>

                {isAuthenticated && (
                    <div className="mobile-menu-footer">
                        <button className="mobile-logout-btn" onClick={() => { logout(); onClose(); }}>
                            <LogOut size={18} />
                            <span>Đăng xuất</span>
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
};

export default MobileMenu;
