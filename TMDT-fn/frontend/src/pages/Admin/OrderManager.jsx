import React, { useState, useEffect } from 'react';
import { Eye, CheckCircle, XCircle, Printer, ChevronLeft, ChevronRight, Loader, Search, RefreshCw, Truck, Package, AlertTriangle } from 'lucide-react';
import SortIcon from '../../components/Admin/SortIcon';
import OrderDetailsModal from '../../components/Admin/OrderDetailsModal';
import toast from 'react-hot-toast';

const API = import.meta.env.VITE_SERVER_API || 'http://localhost:8080';

const OrderManager = () => {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [currentPage, setCurrentPage] = useState(1);
    const [pagination, setPagination] = useState({ total_pages: 1, total_items: 0 });
    const [searchTerm, setSearchTerm] = useState('');
    const [statusFilter, setStatusFilter] = useState('Tất cả');
    const [fromDate, setFromDate] = useState('');
    const [toDate, setToDate] = useState('');
    const [sortConfig, setSortConfig] = useState({ key: 'order_date', direction: 'desc' });
    const [selectedOrderDetails, setSelectedOrderDetails] = useState(null);
    const [showDetailsModal, setShowDetailsModal] = useState(false);
    const [loadingDetails, setLoadingDetails] = useState(false);
    const [isUpdatingStatus, setIsUpdatingStatus] = useState(null);
    const [cancelModal, setCancelModal] = useState({ open: false, orderId: null });

    const orderStatuses = [
        'Tất cả',
        'Chờ xác nhận',
        'Chờ lấy hàng',
        'Chờ giao hàng',
        'Đã giao',
        'Trả hàng',
        'Đã hủy'
    ];


    const fetchOrders = async (page = 1, isSilent = false) => {
        if (!isSilent) setLoading(true);
        try {
            const res = await fetch(`${API}/api/order/all?page=${page}`);
            const data = await res.json();
            if (data.success) {
                setOrders(data.data);
                setPagination(data.pagination);
            } else {
                toast.error(data.message || 'Lỗi khi tải đơn hàng', { id: 'fetch-orders-error' });
            }
        } catch (error) {
            console.error('Fetch Error:', error);
            if (!isSilent) toast.error('Không thể kết nối đến máy chủ', { id: 'fetch-orders-error' });
        } finally {
            if (!isSilent) setLoading(false);
        }
    };

    useEffect(() => {
        fetchOrders(currentPage);
        const interval = setInterval(() => fetchOrders(currentPage, true), 30000);
        return () => clearInterval(interval);
    }, [currentPage]);

    const updateStatus = async (id, newStatusId) => {
        if (isUpdatingStatus === id) return;
        setIsUpdatingStatus(id);
        try {
            const adminId = JSON.parse(localStorage.getItem('user')).id;
            const res = await fetch(`${API}/api/order/status/${id}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ order_id: id, status_id: newStatusId, admin_id: adminId })
            });
            const data = await res.json();
            if (data.success) {
                toast.success('Cập nhật trạng thái thành công', { id: `status-${id}` });
                fetchOrders(currentPage, true);
            } else {
                toast.error(data.message || 'Lỗi cập nhật', { id: `status-${id}` });
            }
        } catch (error) {
            toast.error('Lỗi kết nối', { id: `status-${id}` });
        } finally {
            setIsUpdatingStatus(null);
        }
    };

    const handleViewDetails = (order) => {
        setSelectedOrderDetails(order);
        setShowDetailsModal(true);
    };

    const handleSort = (key) => {
        let direction = 'asc';
        if (sortConfig.key === key && sortConfig.direction === 'asc') {
            direction = 'desc';
        }
        setSortConfig({ key, direction });
    };

    const sortedOrders = [...orders].sort((a, b) => {
        if (!sortConfig.key) return 0;

        let aValue = a[sortConfig.key];
        let bValue = b[sortConfig.key];

        if (sortConfig.key === 'total_amount') {
            aValue = Number(aValue) || 0;
            bValue = Number(bValue) || 0;
        } else if (sortConfig.key === 'order_date') {
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

    const getStatusBadge = (status) => {
        switch (status) {
            case 'Chờ xác nhận': return 'admin-badge-warning';
            case 'Chờ lấy hàng': return 'admin-badge-blue';
            case 'Chờ giao hàng': return 'admin-badge-blue';
            case 'Đã giao': return 'admin-badge-success';
            case 'Trả hàng': return 'admin-badge-warning';
            case 'Đã hủy': return 'admin-badge-danger';
            default: return '';
        }
    };

    const filteredOrders = sortedOrders.filter(order => {
        const matchesSearch = String(order.order_id).toLowerCase().includes(searchTerm.toLowerCase()) ||
            (order.customer_name || '').toLowerCase().includes(searchTerm.toLowerCase());
        const matchesStatus = statusFilter === 'Tất cả' || order.status_name === statusFilter;

        let matchesDate = true;
        if (order.order_date) {
            const orderDateTime = new Date(order.order_date).getTime();

            if (fromDate) {
                const start = new Date(fromDate).getTime();
                if (orderDateTime < start) matchesDate = false;
            }
            if (toDate) {
                const end = new Date(toDate).setHours(23, 59, 59, 999);
                if (orderDateTime > end) matchesDate = false;
            }
        }

        return matchesSearch && matchesStatus && matchesDate;
    });

    const totalPages = pagination.total_pages;
    const currentOrders = filteredOrders;

    const handlePageChange = (pageNumber) => {
        setCurrentPage(pageNumber);
    };

    useEffect(() => {
        setCurrentPage(1);
    }, [searchTerm, statusFilter, fromDate, toDate]);

    const handlePrint = () => {
        window.print();
    };

    return (
        <div className="admin-order-manager">
            <div className="admin-card">
                <div className="admin-card-header" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '16px 24px', gap: '24px' }}>
                    <h2 style={{ margin: 0, whiteSpace: 'nowrap' }}>Quản lý đơn hàng</h2>

                    <div style={{ flex: 1, maxWidth: '800px', backgroundColor: 'var(--admin-bg-soft)', padding: '12px', borderRadius: '12px', border: '1px solid var(--admin-border)', display: 'flex', flexDirection: 'column', gap: '12px' }}>
                        <div className="admin-search-wrapper" style={{ position: 'relative', width: '100%' }}>
                            <Search size={18} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--admin-text-muted)' }} />
                            <input
                                type="text"
                                className="admin-form-input"
                                placeholder="Tìm mã đơn hoặc khách hàng..."
                                style={{ paddingLeft: '40px', width: '100%', height: '36px', borderRadius: '8px', border: '1px solid #e5e7eb' }}
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                        </div>

                        <div style={{ display: 'flex', gap: '16px', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'center' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <span style={{ fontSize: '0.85rem', color: 'var(--admin-text-muted)', whiteSpace: 'nowrap' }}>Trạng thái:</span>
                                <select
                                    className="admin-form-input"
                                    style={{ width: '140px', height: '32px', padding: '0 8px', fontSize: '0.9rem', border: '1px solid #e5e7eb' }}
                                    value={statusFilter}
                                    onChange={(e) => setStatusFilter(e.target.value)}
                                >
                                    {orderStatuses.map(status => (
                                        <option key={status} value={status}>
                                            {status === 'Tất cả' ? 'Tất cả' : status}
                                        </option>
                                    ))}
                                </select>
                            </div>

                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <span style={{ fontSize: '0.85rem', color: 'var(--admin-text-muted)', whiteSpace: 'nowrap' }}>Từ:</span>
                                <input
                                    type="date"
                                    className="admin-form-input"
                                    style={{ height: '32px', padding: '0 8px', width: '140px', fontSize: '0.9rem', border: '1px solid #e5e7eb' }}
                                    value={fromDate}
                                    onChange={(e) => setFromDate(e.target.value)}
                                />
                            </div>

                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <span style={{ fontSize: '0.85rem', color: 'var(--admin-text-muted)', whiteSpace: 'nowrap' }}>Đến:</span>
                                <input
                                    type="date"
                                    className="admin-form-input"
                                    style={{ height: '32px', padding: '0 8px', width: '140px', fontSize: '0.9rem', border: '1px solid #e5e7eb' }}
                                    value={toDate}
                                    onChange={(e) => setToDate(e.target.value)}
                                />
                            </div>

                            {(searchTerm || statusFilter !== 'Tất cả' || fromDate || toDate) && (
                                <button
                                    style={{ background: 'none', border: 'none', color: '#ef4444', cursor: 'pointer', fontSize: '0.85rem', padding: '0 8px', fontWeight: '500' }}
                                    onClick={() => {
                                        setSearchTerm('');
                                        setStatusFilter('Tất cả');
                                        setFromDate('');
                                        setToDate('');
                                    }}
                                >
                                    Xóa lọc
                                </button>
                            )}
                        </div>
                    </div>

                    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                        <button className="admin-btn admin-btn-outline" style={{ whiteSpace: 'nowrap', width: '100%' }}>Xuất Excel</button>
                        <button
                            className="admin-btn admin-btn-outline"
                            style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px', width: '100%' }}
                            onClick={() => {
                                const btnIcon = document.getElementById('refresh-icon');
                                if (btnIcon) btnIcon.classList.add('spin');
                                fetchOrders(currentPage, true).finally(() => {
                                    if (btnIcon) btnIcon.classList.remove('spin');
                                });
                            }}
                        >
                            <RefreshCw id="refresh-icon" size={16} />
                            Làm mới
                        </button>
                    </div>
                </div>

                <div className="admin-table-container">
                    {loading ? (
                        <div style={{ padding: '40px', textAlign: 'center' }}>
                            <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: '12px' }} />
                            <p>Đang tải dữ liệu đơn hàng...</p>
                        </div>
                    ) : (
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th onClick={() => handleSort('order_id')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Mã đơn
                                            <SortIcon activeKey={sortConfig.key} columnKey="order_id" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('customer_name')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Khách hàng
                                            <SortIcon activeKey={sortConfig.key} columnKey="customer_name" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('order_date')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Ngày đặt
                                            <SortIcon activeKey={sortConfig.key} columnKey="order_date" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('total_amount')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Tổng tiền
                                            <SortIcon activeKey={sortConfig.key} columnKey="total_amount" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('payment_method')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            PT Thanh toán
                                            <SortIcon activeKey={sortConfig.key} columnKey="payment_method" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th onClick={() => handleSort('status_name')} style={{ cursor: 'pointer' }}>
                                        <div style={{ display: 'flex', alignItems: 'center' }}>
                                            Trạng thái
                                            <SortIcon activeKey={sortConfig.key} columnKey="status_name" direction={sortConfig.direction} />
                                        </div>
                                    </th>
                                    <th>Thao tác</th>
                                    <th style={{ textAlign: 'center' }}>In</th>
                                </tr>
                            </thead>
                            <tbody>
                                {currentOrders.map((order) => {
                                    const orderSubtotal = Number(order.total_amount);
                                    const orderShipping = order.items ? order.items.reduce((s, i) => s + Number(i.shipping_price || 0), 0) : 0;
                                    const orderShippingDiscount = order.items ? order.items.reduce((s, i) => s + Number(i.shipping_support_price || 0), 0) : 0;
                                    const orderProductDiscount = order.items ? order.items.reduce((s, i) => s + Number(i.product_support_price || 0), 0) : 0;
                                    const finalGrandTotal = orderSubtotal + orderShipping - orderShippingDiscount - orderProductDiscount;

                                    return (
                                        <tr key={order.order_id}>
                                            <td style={{ fontWeight: '600', color: 'var(--admin-primary)' }}>#{order.order_id}</td>
                                            <td>{order.customer_name}</td>
                                            <td>{new Date(order.order_date).toLocaleString('vi-VN')}</td>
                                            <td style={{ fontWeight: '500' }}>{finalGrandTotal.toLocaleString()}đ</td>
                                            <td>{order.payment_method}</td>
                                            <td>
                                                <span className={`admin-badge ${getStatusBadge(order.status_name)}`}>
                                                    {order.status_name}
                                                </span>
                                            </td>
                                            <td>
                                                <div style={{ display: 'flex', gap: '8px' }}>
                                                    <button
                                                        className="admin-btn admin-btn-outline admin-btn-sm"
                                                        title="Xem chi tiết đơn hàng"
                                                        onClick={() => handleViewDetails(order)}
                                                    >
                                                        <Eye size={14} />
                                                    </button>
                                                    {/* 1. Nút Duyệt Đơn (Từ Chờ xác nhận -> Chờ lấy hàng) */}
                                                    {order.status_id === 1 && (
                                                        <button
                                                            className="admin-btn admin-btn-outline admin-btn-sm"
                                                            style={{ color: 'var(--admin-success)' }}
                                                            title="Bấm để Xác nhận đơn & báo kho chuẩn bị Hàng"
                                                            onClick={() => updateStatus(order.order_id, 2)}
                                                            disabled={isUpdatingStatus === order.order_id}
                                                        >
                                                            {isUpdatingStatus === order.order_id ? <Loader className="spin" size={14} /> : <CheckCircle size={14} />}
                                                        </button>
                                                    )}

                                                    {/* 2. Nút Bàn giao Shipper (Từ Chờ lấy hàng -> Chờ/Đang giao hàng) */}
                                                    {order.status_id === 2 && (
                                                        <button
                                                            className="admin-btn admin-btn-outline admin-btn-sm"
                                                            style={{ color: '#ce930a' }}
                                                            title="Bàn giao Shipper / Đang giao hàng"
                                                            onClick={() => updateStatus(order.order_id, 3)}
                                                            disabled={isUpdatingStatus === order.order_id}
                                                        >
                                                            {isUpdatingStatus === order.order_id ? <Loader className="spin" size={14} /> : <Package size={14} />}
                                                        </button>
                                                    )}

                                                    {/* 3. Nút Xác nhận Thành Công (Từ Đang giao hàng -> Đã giao) */}
                                                    {order.status_id === 3 && (
                                                        <button
                                                            className="admin-btn admin-btn-outline admin-btn-sm"
                                                            style={{ color: 'var(--admin-primary)' }}
                                                            title="Đánh dấu đã Giao Hàng thành công"
                                                            onClick={() => updateStatus(order.order_id, 4)}
                                                            disabled={isUpdatingStatus === order.order_id}
                                                        >
                                                            {isUpdatingStatus === order.order_id ? <Loader className="spin" size={14} /> : <Truck size={14} />}
                                                        </button>
                                                    )}

                                                    {/* Chỉ cho phép Hủy khi đơn còn nằm ở nhà (Chờ xác nhận & Chờ lấy hàng). Đã đi giao hoặc Hủy/Thành công rồi sẽ KHÔNG được thấy nút Hủy nữa! */}
                                                    {order.status_id === 1 && (
                                                        <button
                                                            className="admin-btn admin-btn-outline admin-btn-sm"
                                                            style={{ color: 'var(--admin-danger)' }}
                                                            title="Hủy đơn giao dịch"
                                                            onClick={() => setCancelModal({ open: true, orderId: order.order_id })}
                                                        >
                                                            <XCircle size={14} />
                                                        </button>
                                                    )}
                                                </div>
                                            </td>
                                            <td style={{ textAlign: 'center' }}>
                                                <button className="admin-btn admin-btn-outline admin-btn-sm" title="In hóa đơn" onClick={handlePrint}>
                                                    <Printer size={14} />
                                                </button>
                                            </td>
                                        </tr>
                                    );
                                })}
                                {currentOrders.length === 0 && (
                                    <tr>
                                        <td colSpan="8" style={{ textAlign: 'center', padding: '24px' }}>Không tìm thấy đơn hàng nào khớp với bộ lọc</td>
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

            <OrderDetailsModal
                selectedOrderDetails={selectedOrderDetails}
                showDetailsModal={showDetailsModal}
                setShowDetailsModal={setShowDetailsModal}
                getStatusBadge={getStatusBadge}
                loadingDetails={loadingDetails}
            />

            {/* ── Confirm Cancel Modal ── */}
            {cancelModal.open && (
                <div className="admin-modal-overlay" onClick={() => setCancelModal({ open: false, orderId: null })}>
                    <div className="admin-modal" style={{ maxWidth: '420px' }} onClick={e => e.stopPropagation()}>
                        <div className="admin-modal-header" style={{ borderBottom: '1px solid var(--admin-border)' }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                <span style={{ color: '#ef4444', display: 'flex' }}><AlertTriangle size={20} /></span>
                                <h2 style={{ margin: 0, fontSize: '1rem' }}>Xác nhận hủy đơn hàng</h2>
                            </div>
                            <button
                                className="admin-btn"
                                style={{ padding: '4px 8px', fontSize: '1.2rem', lineHeight: 1 }}
                                onClick={() => setCancelModal({ open: false, orderId: null })}
                            >×</button>
                        </div>
                        <div className="admin-modal-body">
                            <p style={{ margin: 0, color: 'var(--admin-text-main)', lineHeight: 1.6 }}>
                                Bạn có chắc chắn muốn <strong style={{ color: '#ef4444' }}>hủy đơn hàng #{cancelModal.orderId}</strong> không?
                                <br />
                                <span style={{ fontSize: '0.85rem', color: 'var(--admin-text-muted)' }}>Thao tác này không thể hoàn tác.</span>
                            </p>
                        </div>
                        <div className="admin-modal-footer">
                            <button
                                className="admin-btn admin-btn-outline"
                                onClick={() => setCancelModal({ open: false, orderId: null })}
                                disabled={isUpdatingStatus === cancelModal.orderId}
                            >
                                Không, giữ lại
                            </button>
                            <button
                                className="admin-btn"
                                style={{ background: '#ef4444', color: 'white', border: 'none' }}
                                disabled={isUpdatingStatus === cancelModal.orderId}
                                onClick={() => {
                                    updateStatus(cancelModal.orderId, 6);
                                    setCancelModal({ open: false, orderId: null });
                                }}
                            >
                                {isUpdatingStatus === cancelModal.orderId
                                    ? <Loader className="spin" size={15} />
                                    : <><XCircle size={15} /> Xác nhận hủy</>}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OrderManager;
