import React, { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { LayoutDashboard, Box, ShoppingCart, Users, BarChart3, LogOut, Settings, Image as ImageIcon, Ticket, Tag, MessageSquare, MessageCircle } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';
import ConfirmModal from '../ReUse/ConfirmModal';

const AdminSidebar = () => {
    const navItems = [
        { name: 'Dashboard', path: '/admin', icon: <LayoutDashboard size={18} /> },
        { name: 'Sản phẩm', path: '/admin/products', icon: <Box size={18} /> },
        //{ name: 'Khuyến mãi', path: '/admin/promotions', icon: <Tag size={18} /> },
        { name: 'Voucher', path: '/admin/vouchers', icon: <Ticket size={18} /> },
        { name: 'Đơn hàng', path: '/admin/orders', icon: <ShoppingCart size={18} /> },
        { name: 'Người dùng', path: '/admin/users', icon: <Users size={18} /> },
        { name: 'Bình luận', path: '/admin/reviews', icon: <MessageSquare size={18} /> },
        { name: 'Tin nhắn', path: '/admin/chat', icon: <MessageCircle size={18} /> },
        { name: 'Banner', path: '/admin/banners', icon: <ImageIcon size={18} /> },
        { name: 'Cấu hình cửa hàng', path: '/admin/store-config', icon: <Settings size={18} /> },
        { name: 'Báo cáo', path: '/admin/reports', icon: <BarChart3 size={18} /> },
    ];

    const { logout } = useAuth();
    const navigate = useNavigate();
    const [showLogoutModal, setShowLogoutModal] = useState(false);

    const handleLogout = () => {
        setShowLogoutModal(true);
    };

    const confirmLogout = () => {
        navigate('/');
    };

    return (
        <aside className="admin-sidebar">
            <div className="admin-sidebar-header">
                Admin Panel
            </div>
            <nav className="admin-sidebar-nav">
                <ul>
                    {navItems.map((item) => (
                        <li key={item.path}>
                            <NavLink
                                to={item.path}
                                className={({ isActive }) => isActive ? 'active' : ''}
                                end={item.path === '/admin'}
                            >
                                {item.icon}
                                <span>{item.name}</span>
                            </NavLink>
                        </li>
                    ))}
                </ul>
            </nav>
            <div className="admin-sidebar-footer">
                <button className="admin-logout-btn" onClick={handleLogout}>
                    <LogOut size={18} />
                    <span>Thoát</span>
                </button>
            </div>

            <ConfirmModal
                isOpen={showLogoutModal}
                onClose={() => setShowLogoutModal(false)}
                onConfirm={confirmLogout}
                title="Rời khỏi trang quản trị"
                message="Bạn có muốn quay về trang chủ không?"
                confirmText="Đồng ý"
                cancelText="Hủy"
            />
        </aside>
    );
};

export default AdminSidebar;
