import React, { useState, useEffect } from 'react';
import {
    DollarSign, ShoppingCart, Users, TrendingUp,
    ArrowUpRight, Clock, Package, CheckCircle,
    BarChart2, RefreshCw, Eye
} from 'lucide-react';
import { Link } from 'react-router-dom';
import './Dashboard.css';

const AdminDashboard = () => {
    const [loading, setLoading] = useState(true);
    const [stats, setStats] = useState(null);
    const [recentOrders, setRecentOrders] = useState([]);
    const API = import.meta.env.VITE_SERVER_API;

    const fetchDashboardData = async () => {
        setLoading(true);
        try {
            // Lấy KPI tổng quát và đơn hàng gần đây
            const [kpiRes, ordersRes] = await Promise.all([
                fetch(`${API}/api/statistic/kpi-dashboard`),
                fetch(`${API}/api/order/all?page=1`)
            ]);

            const kpiData = await kpiRes.json();
            const ordersData = await ordersRes.json();

            // Nếu kpiData là mảng, lấy phần tử đầu tiên, nếu là object thì giữ nguyên
            const actualStats = Array.isArray(kpiData) ? kpiData[0] : kpiData;
            setStats(actualStats);

            if (ordersData.success) {
                setRecentOrders(ordersData.data.slice(0, 5));
            }
        } catch (error) {
            console.error("Dashboard fetch error:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const fmt = (v) => v != null ? Number(v).toLocaleString('vi-VN') : '0';
    const currency = (v) => v != null ? Number(v).toLocaleString('vi-VN') + 'đ' : '0đ';

    if (loading) return (
        <div className="db-loading">
            <RefreshCw className="spin" size={32} />
            <p>Đang tải dữ liệu tổng quan...</p>
        </div>
    );

    return (
        <div className="db-container">
            <div className="db-header">
                <div>
                    <h1 className="db-title">Bảng điều khiển</h1>
                    <p className="db-sub">Chào mừng bạn trở lại, đây là những gì đang diễn ra hôm nay.</p>
                </div>
                <button className="admin-btn admin-btn-outline" onClick={fetchDashboardData}>
                    <RefreshCw size={14} /> Làm mới
                </button>
            </div>

            {/* ── KPI Cards ── */}
            <div className="db-kpi-grid">
                <div className="db-card db-kpi-card">
                    <div className="db-kpi-icon ico-revenue"><DollarSign size={24} /></div>
                    <div className="db-kpi-info">
                        <span className="db-kpi-label">Tổng doanh thu</span>
                        <h3 className="db-kpi-val">{currency(stats?.total_revenue || stats?.revenue)}</h3>
                    </div>
                </div>
                <div className="db-card db-kpi-card">
                    <div className="db-kpi-icon ico-orders"><ShoppingCart size={24} /></div>
                    <div className="db-kpi-info">
                        <span className="db-kpi-label">Tổng đơn hàng</span>
                        <h3 className="db-kpi-val">{fmt(stats?.total_orders || stats?.order_count || stats?.count)}</h3>
                    </div>
                </div>
                <div className="db-card db-kpi-card">
                    <div className="db-kpi-icon ico-users"><Users size={24} /></div>
                    <div className="db-kpi-info">
                        <span className="db-kpi-label">Số khách hàng</span>
                        <h3 className="db-kpi-val">{fmt(stats?.total_customers || stats?.customer_count || stats?.user_count)}</h3>
                    </div>
                </div>
                <div className="db-card db-kpi-card">
                    <div className="db-kpi-icon ico-avg"><TrendingUp size={24} /></div>
                    <div className="db-kpi-info">
                        <span className="db-kpi-label">Giá trị đơn TB</span>
                        <h3 className="db-kpi-val">{currency(Math.round(stats?.average_order_value || stats?.avg_order_value || 0))}</h3>
                    </div>
                </div>
            </div>

            <div className="db-main-grid">
                {/* ── Recent Orders ── */}
                <div className="db-card db-main-card">
                    <div className="db-card-header">
                        <h2>Đơn hàng mới nhất</h2>
                        <Link title="Xem tất cả đơn hàng" to="/admin/orders" className="db-view-all">Tất cả</Link>
                    </div>
                    <div className="db-table-container">
                        <table className="db-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Khách hàng</th>
                                    <th>Ngày</th>
                                    <th>Tổng tiền</th>
                                    <th>Trạng thái</th>
                                </tr>
                            </thead>
                            <tbody>
                                {recentOrders.map(order => (
                                    <tr key={order.order_id}>
                                        <td className="fw-600">#{order.order_id}</td>
                                        <td>{order.customer_name}</td>
                                        <td className="txt-muted">{new Date(order.order_date).toLocaleDateString('vi-VN')}</td>
                                        <td className="fw-500">{Number(order.total_amount).toLocaleString()}đ</td>
                                        <td>
                                            <span className={`db-badge status-${order.status_id}`}>
                                                {order.status_name}
                                            </span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* ── Quick Actions / Summaries ── */}
                <div className="db-side-col">
                    <div className="db-card db-mini-card">
                        <h3>Hoạt động nhanh</h3>
                        <div className="db-actions-list">
                            <Link to="/admin/products" className="db-action-item">
                                <Package size={18} /> Quản lý kho hàng
                            </Link>
                            <Link to="/admin/reports" className="db-action-item">
                                <BarChart2 size={18} /> Xem báo cáo chi tiết
                            </Link>
                            <Link to="/admin/users" className="db-action-item">
                                <Users size={18} /> Danh sách người dùng
                            </Link>
                        </div>
                    </div>

                    <div className="db-card db-mini-card">
                        <h3>Hiệu suất</h3>
                        <div className="db-perf-item">
                            <div className="db-perf-info">
                                <span>Tỷ lệ hoàn thành</span>
                                <span>94%</span>
                            </div>
                            <div className="db-perf-bar"><div className="db-perf-fill" style={{ width: '94%', background: '#10b981' }}></div></div>
                        </div>
                        <div className="db-perf-item">
                            <div className="db-perf-info">
                                <span>Sản phẩm còn hàng</span>
                                <span>82%</span>
                            </div>
                            <div className="db-perf-bar"><div className="db-perf-fill" style={{ width: '82%', background: '#3b82f6' }}></div></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AdminDashboard;
