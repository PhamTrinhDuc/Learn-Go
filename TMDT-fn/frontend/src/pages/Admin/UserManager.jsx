import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search, UserCircle, Edit, ShieldAlert, CheckCircle, ChevronLeft, ChevronRight, Loader, UserMinus, UserCheck, UserPlus, RefreshCw, Key } from 'lucide-react';
import SortIcon from '../../components/Admin/SortIcon';
import ConfirmStatusModal from '../../components/Admin/ConfirmStatusModal';

const UserManager = () => {
    const navigate = useNavigate();
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [roleFilter, setRoleFilter] = useState('Tất cả');
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 8;
    const [sortConfig, setSortConfig] = useState({ key: null, direction: 'asc' });
    const [selectedUserForRole, setSelectedUserForRole] = useState(null);
    const [newSelectedRole, setNewSelectedRole] = useState('');
    const [roleError, setRoleError] = useState('');
    const [loadingRole, setLoadingRole] = useState(false);
    const [showRoleModal, setShowRoleModal] = useState(false);
    const [userToToggle, setUserToToggle] = useState(null);
    const [isTogglingStatus, setIsTogglingStatus] = useState(false);
    const [showConfirmModal, setShowConfirmModal] = useState(false);
    const [roles, setRoles] = useState([]);
    const [toast, setToast] = useState({ show: false, message: '', type: 'success' });
    const roleOptions = ['Quản trị viên', 'Khách hàng'];

    const [showResetModal, setShowResetModal] = useState(false);
    const [selectedUserForReset, setSelectedUserForReset] = useState(null);
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [resetError, setResetError] = useState('');
    const [loadingReset, setLoadingReset] = useState(false);

    const showToast = (message, type = 'success') => {
        setToast({ show: true, message, type });
        setTimeout(() => setToast({ show: false, message: '', type: 'success' }), 4000);
    };

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            const data = await response.json();
            setUsers(data);
        } catch (err) {
            console.error("Error fetching users:", err);
        } finally {
            setLoading(false);
        }
    };

    const fetchRoles = async () => {
        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/role`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            const data = await response.json();
            if (data.roles) {
                const processedRoles = data.roles.map(item => {
                    if (item.role.toLowerCase() === 'admin') {
                        return 'Quản trị viên';
                    } else if (item.role.toLowerCase() === 'customer') {
                        return 'Khách hàng';
                    }
                    return item.role;
                });
                const uniqueRoles = [...new Set(processedRoles)];
                setRoles(["Tất cả", ...uniqueRoles]);
            }
        } catch (err) {
            console.error("Error fetching roles:", err);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    useEffect(() => {
        fetchRoles();
    }, []);

    const handleReload = () => {
        fetchUsers();
        fetchRoles();
    };

    const toggleStatus = (user) => {
        setUserToToggle(user);
        setShowConfirmModal(true);
    };

    const confirmToggleStatus = async () => {
        if (!userToToggle) return;

        setIsTogglingStatus(true);
        const newStatus = !userToToggle.is_lock;

        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/status/${userToToggle.id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ is_lock: newStatus })
            });

            const result = await response.json();

            if (!response.ok) {
                showToast(result.message || 'Thay đổi trạng thái thất bại!', 'error');
                setIsTogglingStatus(false);
                return;
            }

            setUsers(users.map(user => {
                if (user.id === userToToggle.id) {
                    return { ...user, is_lock: newStatus };
                }
                return user;
            }));

            showToast(`Đã ${newStatus ? 'khóa' : 'mở khóa'} tài khoản thành công!`, 'success');
            setShowConfirmModal(false);
            setUserToToggle(null);
        } catch (error) {
            console.error('Lỗi Toggle Status:', error);
            showToast('Lỗi kết nối máy chủ khi đổi trạng thái!', 'error');
        } finally {
            setIsTogglingStatus(false);
        }
    };

    const handleUpdateRole = async () => {
        if (!selectedUserForRole || !newSelectedRole) return;
        const newRole = newSelectedRole === 'Quản trị viên' ? 'admin' : 'customer';

        if (newRole === selectedUserForRole.role) {
            setShowRoleModal(false);
            setSelectedUserForRole(null);
            return;
        }

        setLoadingRole(true);
        setRoleError('');

        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/role/${selectedUserForRole.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ role: newRole }),
            });

            const result = await response.json();

            if (!response.ok) {
                setRoleError(result.message || 'Cập nhật quyền thất bại!');
                showToast(result.message || 'Cập nhật quyền thất bại!', 'error');
                setLoadingRole(false);
                return;
            }

            setUsers(users.map(user => {
                if (user.id === selectedUserForRole.id) {
                    return { ...user, role: newRole };
                }
                return user;
            }));

            showToast('Cập nhật quyền thành công!', 'success');
            setShowRoleModal(false);
            setSelectedUserForRole(null);
        } catch (error) {
            console.error('Lỗi Update Role:', error);
            setRoleError('Lỗi kết nối máy chủ khi cập nhật quyền!');
            showToast('Lỗi kết nối máy chủ khi cập nhật quyền!', 'error');
        } finally {
            setLoadingRole(false);
        }
    };

    const handleResetPassword = async () => {
        if (!newPassword) {
            setResetError('Vui lòng nhập mật khẩu mới');
            return;
        }
        if (newPassword.length < 6) {
            setResetError('Mật khẩu phải từ 6 ký tự trở lên');
            return;
        }
        if (newPassword !== confirmPassword) {
            setResetError('Mật khẩu xác nhận không khớp');
            return;
        }

        setLoadingReset(true);
        setResetError('');

        try {
            const response = await fetch(`${import.meta.env.VITE_SERVER_API}/api/user/admin-reset-password/${selectedUserForReset.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ password: newPassword }),
            });

            const result = await response.json();

            if (!response.ok) {
                setResetError(result.message || 'Đặt lại mật khẩu thất bại!');
                showToast(result.message || 'Đặt lại mật khẩu thất bại!', 'error');
                setLoadingReset(false);
                return;
            }

            showToast('Đặt lại mật khẩu thành công!', 'success');
            setShowResetModal(false);
            setSelectedUserForReset(null);
            setNewPassword('');
            setConfirmPassword('');
        } catch (error) {
            console.error('Lỗi Reset Password:', error);
            setResetError('Lỗi kết nối máy chủ!');
            showToast('Lỗi kết nối máy chủ!', 'error');
        } finally {
            setLoadingReset(false);
        }
    };

    const handleSort = (key) => {
        let direction = 'asc';
        if (sortConfig.key === key && sortConfig.direction === 'asc') {
            direction = 'desc';
        }
        setSortConfig({ key, direction });
    };

    const sortedUsers = [...users].sort((a, b) => {
        if (!sortConfig.key) return 0;

        let aValue = a[sortConfig.key];
        let bValue = b[sortConfig.key];

        if (sortConfig.key === 'id') {
            aValue = Number(aValue) || 0;
            bValue = Number(bValue) || 0;
        } else if (sortConfig.key === 'joined_date') {
            aValue = new Date(aValue).getTime() || 0;
            bValue = new Date(bValue).getTime() || 0;
        } else {
            aValue = String(aValue || '').toLowerCase();
            bValue = String(bValue || '').toLowerCase();
        }

        if (aValue < bValue) {
            return sortConfig.direction === 'asc' ? -1 : 1;
        }
        if (aValue > bValue) {
            return sortConfig.direction === 'asc' ? 1 : -1;
        }
        return 0;
    });

    const filteredUsers = sortedUsers.filter(user => {
        const matchesSearch = (user.full_name || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
            (user.username || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
            (user.email || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
            (user.num_phone || '').includes(searchTerm);

        const userRoleLabel = user.role === 'admin' ? 'Quản trị viên' : 'Khách hàng';
        const matchesRole = roleFilter === 'Tất cả' || userRoleLabel === roleFilter;

        return matchesSearch && matchesRole;
    });

    const totalPages = Math.ceil(filteredUsers.length / itemsPerPage);
    const indexOfLastItem = currentPage * itemsPerPage;
    const indexOfFirstItem = indexOfLastItem - itemsPerPage;
    const currentUsers = filteredUsers.slice(indexOfFirstItem, indexOfLastItem);

    const handlePageChange = (pageNumber) => {
        setCurrentPage(pageNumber);
    };

    useEffect(() => {
        setCurrentPage(1);
    }, [searchTerm, roleFilter]);

    return (
        <div className="admin-user-manager">
            <div className="admin-card">
                <div className="admin-card-header">
                    <h2>Quản lý người dùng</h2>
                    <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'center' }}>
                        <div className="admin-search-wrapper" style={{ position: 'relative' }}>
                            <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--admin-text-muted)' }} />
                            <input
                                type="text"
                                className="admin-form-input"
                                placeholder="Tìm theo tên, email, SĐT..."
                                style={{ paddingLeft: '40px', width: '250px' }}
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>

                        <select
                            className="admin-form-input"
                            style={{ width: '160px', height: '40px' }}
                            value={roleFilter}
                            onChange={(e) => setRoleFilter(e.target.value)}
                        >
                            {roles.map(role => (
                                <option key={role} value={role}>
                                    {role === 'Tất cả' ? 'Tất cả vai trò' : role}
                                </option>
                            ))}
                        </select>

                        <button
                            className="admin-btn admin-btn-outline"
                            style={{ display: 'flex', alignItems: 'center', gap: '8px', height: '40px' }}
                            onClick={handleReload}
                            disabled={loading}
                        >
                            <RefreshCw size={18} className={loading ? 'spin' : ''} />
                            <span>Làm mới</span>
                        </button>

                        <button
                            className="admin-btn admin-btn-primary"
                            onClick={() => navigate('/admin/users/add')}
                            style={{ height: '40px', padding: '0 16px', borderRadius: '8px', whiteSpace: 'nowrap' }}
                        >
                            <UserPlus size={18} />
                            <span>Thêm người dùng</span>
                        </button>
                    </div>
                </div>

                <div className="admin-table-container">
                    {loading ? (
                        <div style={{ padding: '40px', textAlign: 'center' }}>
                            <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: '12px' }} />
                            <p>Đang tải dữ liệu người dùng...</p>
                        </div>
                    ) : (
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th onClick={() => handleSort('id')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            ID
                                            <SortIcon activeKey={sortConfig.key} columnKey="id" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('full_name')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Người dùng
                                            <SortIcon activeKey={sortConfig.key} columnKey="full_name" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th>Liên hệ</th>
                                    <th onClick={() => handleSort('role')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Vai trò
                                            <SortIcon activeKey={sortConfig.key} columnKey="role" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('joined_date')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Ngày tham gia
                                            <SortIcon activeKey={sortConfig.key} columnKey="joined_date" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('is_lock')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Trạng thái
                                            <SortIcon activeKey={sortConfig.key} columnKey="is_lock" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th>Thao tác</th>
                                </tr>
                            </thead>
                            <tbody>
                                {currentUsers.map((user) => (
                                    <tr key={user.id}>
                                        <td style={{ fontWeight: '600', color: 'var(--admin-primary)' }}>#{user.id}</td>
                                        <td>
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                                <div style={{ width: '32px', height: '32px', borderRadius: '50%', background: 'var(--admin-bg-light)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--admin-primary)' }}>
                                                    <UserCircle size={20} />
                                                </div>
                                                <div style={{ fontWeight: '500' }}>{user.full_name}</div>
                                            </div>
                                        </td>
                                        <td>
                                            <div style={{ fontSize: '0.85rem' }}>{user.email}</div>
                                            <div style={{ fontSize: '0.85rem', color: 'var(--admin-text-muted)' }}>{user.num_phone}</div>
                                        </td>
                                        <td>{user.role === 'admin' ? 'Quản trị viên' : 'Khách hàng'}</td>
                                        <td>{user.joined_date ? new Date(user.joined_date).toLocaleDateString('vi-VN') : '---'}</td>
                                        <td>
                                            <span className={`admin-badge ${!user.is_lock ? 'admin-badge-success' : 'admin-badge-danger'}`}>
                                                {!user.is_lock ? 'Hoạt động' : 'Bị khóa'}
                                            </span>
                                        </td>
                                        <td>
                                            <div style={{ display: 'flex', gap: '8px' }}>
                                                <button
                                                    className="admin-btn admin-btn-outline admin-btn-sm"
                                                    title="Sửa vai trò"
                                                    onClick={() => {
                                                        setSelectedUserForRole(user);
                                                        setNewSelectedRole(user.role === 'admin' ? 'Quản trị viên' : 'Khách hàng');
                                                        setRoleError('');
                                                        setShowRoleModal(true);
                                                    }}
                                                >
                                                    <Edit size={14} />
                                                </button>
                                                <button
                                                    className="admin-btn admin-btn-outline admin-btn-sm"
                                                    title="Đặt lại mật khẩu"
                                                    onClick={() => {
                                                        setSelectedUserForReset(user);
                                                        setNewPassword('');
                                                        setConfirmPassword('');
                                                        setResetError('');
                                                        setShowResetModal(true);
                                                    }}
                                                >
                                                    <Key size={14} />
                                                </button>
                                                <button
                                                    className="admin-btn admin-btn-outline admin-btn-sm"
                                                    style={{ color: !user.is_lock ? 'var(--admin-danger)' : 'var(--admin-success)' }}
                                                    title={!user.is_lock ? 'Khóa tài khoản' : 'Kích hoạt'}
                                                    onClick={() => toggleStatus(user)}
                                                >
                                                    {!user.is_lock ? <UserMinus size={14} /> : <UserCheck size={14} />}
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                                {currentUsers.length === 0 && (
                                    <tr>
                                        <td colSpan="7" style={{ textAlign: 'center', padding: '24px' }}>Không tìm thấy người dùng nào</td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    )}
                </div>

                {totalPages > 1 && (
                    <div className="admin-pagination">
                        <button
                            className="admin-page-btn"
                            disabled={currentPage === 1}
                            onClick={() => handlePageChange(currentPage - 1)}
                        >
                            <ChevronLeft size={16} />
                        </button>

                        {[...Array(totalPages)].map((_, i) => (
                            <button
                                key={i + 1}
                                className={`admin-page-btn ${currentPage === i + 1 ? 'active' : ''}`}
                                onClick={() => handlePageChange(i + 1)}
                            >
                                {i + 1}
                            </button>
                        ))}

                        <button
                            className="admin-page-btn"
                            disabled={currentPage === totalPages}
                            onClick={() => handlePageChange(currentPage + 1)}
                        >
                            <ChevronRight size={16} />
                        </button>
                    </div>
                )}
            </div>

            <ConfirmStatusModal
                isOpen={showConfirmModal}
                user={userToToggle}
                onClose={() => setShowConfirmModal(false)}
                onConfirm={confirmToggleStatus}
                isProcessing={isTogglingStatus}
            />

            {showRoleModal && selectedUserForRole && (
                <div className="admin-modal-overlay">
                    <div className="admin-modal" style={{ maxWidth: '400px' }}>
                        <div className="admin-modal-header">
                            <h2>Thay đổi vai trò</h2>
                            <button className="admin-btn" onClick={() => setShowRoleModal(false)}>×</button>
                        </div>
                        <div className="admin-modal-body">
                            <div style={{ marginBottom: '16px' }}>
                                <p style={{ margin: '0 0 8px 0', fontSize: '0.9rem', color: 'var(--admin-text-muted)' }}>
                                    Thay đổi vai trò cho người dùng: <strong>{selectedUserForRole.full_name}</strong>
                                </p>
                            </div>
                            <div className="admin-form-group">
                                <label className="admin-form-label">Chọn vai trò mới</label>
                                <select
                                    className="admin-form-input"
                                    value={newSelectedRole}
                                    onChange={(e) => setNewSelectedRole(e.target.value)}
                                    disabled={loadingRole}
                                >
                                    {roles.filter(r => r !== 'Tất cả').map(role => (
                                        <option key={role} value={role}>{role}</option>
                                    ))}
                                </select>
                            </div>

                            {roleError && (
                                <div style={{ marginTop: '12px', fontSize: '0.85rem', color: 'var(--admin-danger)', background: 'rgba(239, 68, 68, 0.05)', padding: '8px', borderRadius: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                                    <ShieldAlert size={14} /> {roleError}
                                </div>
                            )}
                        </div>
                        <div className="admin-modal-footer" style={{ gap: '12px' }}>
                            <button
                                className="admin-btn admin-btn-outline"
                                onClick={() => setShowRoleModal(false)}
                                disabled={loadingRole}
                            >
                                Hủy
                            </button>
                            <button
                                className="admin-btn admin-btn-primary"
                                onClick={handleUpdateRole}
                                disabled={loadingRole}
                            >
                                {loadingRole ? <Loader className="spin" size={16} /> : 'Xác nhận đổi'}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {showResetModal && selectedUserForReset && (
                <div className="admin-modal-overlay">
                    <div className="admin-modal" style={{ maxWidth: '400px' }}>
                        <div className="admin-modal-header">
                            <h2>Đặt lại mật khẩu</h2>
                            <button className="admin-btn" onClick={() => setShowResetModal(false)}>×</button>
                        </div>
                        <div className="admin-modal-body">
                            <div style={{ marginBottom: '16px' }}>
                                <p style={{ margin: '0 0 8px 0', fontSize: '0.9rem', color: 'var(--admin-text-muted)' }}>
                                    Đặt lại mật khẩu cho tài khoản: <strong>{selectedUserForReset.username || selectedUserForReset.email}</strong>
                                </p>
                            </div>
                            <div className="admin-form-group" style={{ marginBottom: '12px' }}>
                                <label className="admin-form-label">Mật khẩu mới</label>
                                <input
                                    type="password"
                                    className="admin-form-input"
                                    placeholder="Nhập mật khẩu mới"
                                    value={newPassword}
                                    onChange={(e) => setNewPassword(e.target.value)}
                                    disabled={loadingReset}
                                />
                            </div>
                            <div className="admin-form-group">
                                <label className="admin-form-label">Xác nhận mật khẩu mới</label>
                                <input
                                    type="password"
                                    className="admin-form-input"
                                    placeholder="Xác nhận mật khẩu mới"
                                    value={confirmPassword}
                                    onChange={(e) => setConfirmPassword(e.target.value)}
                                    disabled={loadingReset}
                                />
                            </div>

                            {resetError && (
                                <div style={{ marginTop: '12px', fontSize: '0.85rem', color: 'var(--admin-danger)', background: 'rgba(239, 68, 68, 0.05)', padding: '8px', borderRadius: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                                    <ShieldAlert size={14} /> {resetError}
                                </div>
                            )}
                        </div>
                        <div className="admin-modal-footer" style={{ gap: '12px' }}>
                            <button
                                className="admin-btn admin-btn-outline"
                                onClick={() => setShowResetModal(false)}
                                disabled={loadingReset}
                            >
                                Hủy
                            </button>
                            <button
                                className="admin-btn admin-btn-primary"
                                onClick={handleResetPassword}
                                disabled={loadingReset}
                            >
                                {loadingReset ? <Loader className="spin" size={16} /> : 'Xác nhận'}
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {toast.show && (
                <div style={{ position: 'fixed', bottom: '24px', left: '50%', transform: 'translateX(-50%)', zIndex: 9999 }}>
                    <div style={{
                        display: 'flex', alignItems: 'center', gap: '12px',
                        padding: '12px 24px', borderRadius: '12px',
                        background: toast.type === 'success' ? '#10b981' : '#ef4444',
                        color: '#fff', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.15)',
                        fontWeight: '500', fontSize: '0.95rem'
                    }}>
                        {toast.type === 'success' ? <CheckCircle size={20} /> : <ShieldAlert size={20} />}
                        {toast.message}
                    </div>
                </div>
            )}
        </div>
    );
};

export default UserManager;
