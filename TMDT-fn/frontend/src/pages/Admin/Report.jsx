import React, { useState, useEffect, useCallback } from 'react';
import {
    TrendingUp, ShoppingCart, Users, CheckCircle,
    Clock, DollarSign, BarChart3, Star, MapPin,
    ArrowUpRight, ArrowDownRight, RefreshCw, AlertCircle,
    Eye, X, Loader
} from 'lucide-react';
import './Report.css';

// ── helper ───────────────────────────────────────────────────────────────────
const KpiCard = ({ label, value, icon: Icon, color, sub }) => (
    <div className="admin-card rpt-kpi">
        <div className="rpt-kpi-icon" style={{ backgroundColor: `${color}15`, color }}>
            <Icon size={24} />
        </div>
        <div className="rpt-kpi-content">
            <span className="rpt-kpi-label">{label}</span>
            <h3 className="rpt-kpi-value">{value}</h3>
            {sub && <span className="rpt-kpi-sub">{sub}</span>}
        </div>
    </div>
);

const SectionTitle = ({ icon: Icon, title, color }) => (
    <div className="rpt-section-title">
        <Icon size={20} style={{ color }} />
        <h2 style={{ borderLeft: `4px solid ${color}` }}>{title}</h2>
    </div>
);

// ── main component ────────────────────────────────────────────────────────────
const AdminReport = () => {
    const [range, setRange] = useState('month');
    const [topPeriod, setTopPeriod] = useState('all');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const [salesOverview, setSalesOverview] = useState([]);
    const [topSellers, setTopSellers] = useState([]);
    const [paymentStats, setPaymentStats] = useState([]);
    const [discountStats, setDiscountStats] = useState([]);
    const [procurement, setProcurement] = useState([]);
    const [loyalCustomers, setLoyalCustomers] = useState([]);
    const [demographics, setDemographics] = useState({ geoDistribution: [], genderDist: [] });
    const [kpi, setKpi] = useState(null);

    const [selectedProvince, setSelectedProvince] = useState('');
    const [provinceCustomers, setProvinceCustomers] = useState([]);
    const [loadingCustomers, setLoadingCustomers] = useState(false);
    const [showProvinceModal, setShowProvinceModal] = useState(false);

    const fetchAll = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const BASE = import.meta.env.VITE_SERVER_API;

            // Helper to fetch safely
            const safeFetch = async (url) => {
                try {
                    const res = await fetch(url);
                    if (!res.ok) {
                        console.error(`Fetch error ${res.status} for: ${url}`);
                        return null;
                    }
                    return await res.json();
                } catch (e) {
                    console.error(`Network error for: ${url}`, e);
                    return null;
                }
            };

            const [sales, top, pay, disc, proc, loyal, demo, kpiData] = await Promise.all([
                safeFetch(`${BASE}/api/statistic/sales-overview?range=${range}`),
                safeFetch(`${BASE}/api/statistic/top-sellers?limit=10&period=${topPeriod}`),
                safeFetch(`${BASE}/api/statistic/payment-revenue`),
                safeFetch(`${BASE}/api/statistic/discount-analysis`),
                safeFetch(`${BASE}/api/statistic/procurement-report`),
                safeFetch(`${BASE}/api/statistic/loyal-customers`),
                safeFetch(`${BASE}/api/statistic/customer-demographics`),
                safeFetch(`${BASE}/api/statistic/kpi-dashboard`),
            ]);

            if (sales) setSalesOverview(Array.isArray(sales) ? sales : []);
            if (top) setTopSellers(Array.isArray(top) ? top : []);
            if (pay) setPaymentStats(Array.isArray(pay) ? pay : []);
            if (disc) setDiscountStats(Array.isArray(disc) ? disc : []);
            if (proc) setProcurement(Array.isArray(proc) ? proc : []);
            if (loyal) setLoyalCustomers(Array.isArray(loyal) ? loyal : []);
            if (demo) setDemographics(demo && typeof demo === 'object' ? demo : { geoDistribution: [], genderDist: [] });
            if (kpiData) setKpi(kpiData || null);

        } catch (err) {
            console.error("General Fetch error:", err);
            setError("Có lỗi hệ thống khi tải báo cáo. Vui lòng kiểm tra Server.");
        } finally {
            setLoading(false);
        }
    }, [range, topPeriod]);

    useEffect(() => { fetchAll(); }, [fetchAll]);

    const handleViewCustomersByProvince = async (province) => {
        setSelectedProvince(province);
        setShowProvinceModal(true);
        setLoadingCustomers(true);
        try {
            const BASE = import.meta.env.VITE_SERVER_API;
            const res = await fetch(`${BASE}/api/statistic/customers-by-province?province=${encodeURIComponent(province)}`);
            if (res.ok) {
                const data = await res.json();
                setProvinceCustomers(data);
            } else {
                setProvinceCustomers([]);
            }
        } catch (e) {
            console.error("Lỗi fetch khách hàng theo tỉnh thành:", e);
            setProvinceCustomers([]);
        } finally {
            setLoadingCustomers(false);
        }
    };

    const fmt = (v) => v != null ? Number(v).toLocaleString('vi-VN') + ' đ' : '0 đ';
    const fmtNum = (v) => v != null ? Number(v).toLocaleString('vi-VN') : '0';

    if (loading) return (
        <div className="rpt-loading">
            <RefreshCw className="spin" size={32} />
            <p>Đang tổng hợp dữ liệu...</p>
        </div>
    );

    if (error) return (
        <div className="rpt-error">
            <AlertCircle color="#dc2626" size={48} />
            <h2>Rất tiếc!</h2>
            <p>{error}</p>
            <button className="admin-btn admin-btn-primary" onClick={fetchAll}>Thử lại</button>
        </div>
    );

    // Tính toán tỷ lệ chốt đơn từ kpiData
    const completedOrders = salesOverview.length > 0 ? salesOverview.reduce((sum, s) => sum + Number(s.completed_orders || 0), 0) : 0;
    const totalOrders = salesOverview.length > 0 ? salesOverview.reduce((sum, s) => sum + Number(s.total_orders || 0), 0) : 0;
    const completionRate = totalOrders > 0 ? ((completedOrders / totalOrders) * 100).toFixed(1) : '0';

    return (
        <div className="rpt-container">
            <div className="rpt-header">
                <div>
                    <h1>Báo cáo & Thống kê</h1>
                    <p>Dữ liệu thời gian thực từ hệ thống kinh doanh</p>
                </div>
                {/* <div className="rpt-range-tabs">
                    {[
                        { k: 'day', v: 'Hôm nay' },
                        { k: 'month', v: 'Tháng này' },
                        { k: 'year', v: 'Năm nay' }
                    ].map(item => (
                        <button
                            key={item.k}
                            className={`rpt-range-btn ${range === item.k ? 'active' : ''}`}
                            onClick={() => setRange(item.k)}
                        >
                            {item.v}
                        </button>
                    ))}
                </div> */}
            </div>

            {/* ── KPI Grid ── */}
            <SectionTitle icon={BarChart3} title="Tổng quan" color="#2563eb" />
            <div className="rpt-grid-4">
                <KpiCard
                    label="Tổng doanh thu"
                    value={fmt(kpi?.total_revenue)}
                    icon={DollarSign}
                    color="#2563eb"
                    sub="Revenue"
                />
                <KpiCard
                    label="Giá trị đơn TB"
                    value={fmt(kpi?.average_order_value)}
                    icon={ShoppingCart}
                    color="#7c3aed"
                    sub="Average Order Value"
                />
                <KpiCard
                    label="Tỷ lệ đơn đã giao"
                    value={`${completionRate}%`}
                    icon={CheckCircle}
                    color="#0891b2"
                    sub={`${fmtNum(completedOrders)} / ${fmtNum(totalOrders)} đơn`}
                />
            </div>

            <div className="rpt-main-grid">
                {/* ── Left Column ── */}
                <div className="rpt-col">
                    {/* Top Sellers */}
                    <SectionTitle icon={Star} title="Top sản phẩm bán chạy" color="#d97706" />
                    <div className="admin-card rpt-mb">
                        <div className="admin-card-header">
                            <div>
                                <h2>Top 10 sản phẩm</h2>
                                <div className="rpt-range-tabs" style={{ marginTop: 8 }}>
                                    {[
                                        { k: 'all', v: 'Tất cả' },
                                        { k: 'week', v: 'Tuần' },
                                        { k: 'month', v: 'Tháng' },
                                        { k: 'year', v: 'Năm' }
                                    ].map(item => (
                                        <button
                                            key={item.k}
                                            className={`rpt-range-btn ${topPeriod === item.k ? 'active' : ''}`}
                                            onClick={() => setTopPeriod(item.k)}
                                            style={{ fontSize: '0.75rem', padding: '2px 8px' }}
                                        >
                                            {item.v}
                                        </button>
                                    ))}
                                </div>
                            </div>
                            <span className="admin-badge admin-badge-warning">{topSellers.length} sản phẩm</span>
                        </div>
                        <div className="admin-table-container">
                            {topSellers.length === 0 ? (
                                <div className="rpt-empty">Chưa có dữ liệu sản phẩm</div>
                            ) : (
                                <table className="admin-table">
                                    <thead>
                                        <tr>
                                            <th>#</th>
                                            <th>Sản phẩm</th>
                                            <th style={{ textAlign: 'center' }}>Đã bán</th>
                                            <th style={{ textAlign: 'right' }}>Doanh thu</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {topSellers.map((p, i) => (
                                            <tr key={p.product_id}>
                                                <td>{i + 1}</td>
                                                <td className="fw-500">{p.name}</td>
                                                <td style={{ textAlign: 'center' }}>{fmtNum(p.total_sold)}</td>
                                                <td style={{ textAlign: 'right' }} className="rpt-revenue">{fmt(p.total_revenue)}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>
                    </div>

                    {/* Payment & Support Analysis */}
                    {/* <SectionTitle icon={CheckCircle} title="Thanh toán & Chiết khấu" color="#059669" />
                    <div className="admin-card">
                        <div className="admin-card-header">
                            <h2>Phân tích hiệu quả</h2>
                        </div>
                        <div className="admin-table-container">
                            {paymentStats.length === 0 ? (
                                <div className="rpt-empty">Chưa có dữ liệu thanh toán</div>
                            ) : (
                                <table className="admin-table">
                                    <thead>
                                        <tr>
                                            <th>P.Thức</th>
                                            <th>Số đơn</th>
                                            <th>Doanh thu</th>
                                            <th>Ship hỗ trợ</th>
                                            <th>Tỷ lệ CK</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {paymentStats.map((pm, i) => {
                                            const disc = discountStats.find(d => d.payment_method === pm.payment_method) || {};
                                            return (
                                                <tr key={i}>
                                                    <td>
                                                        <span className="admin-badge admin-badge-blue">
                                                            {pm.payment_method || '—'}
                                                        </span>
                                                    </td>
                                                    <td>{fmtNum(pm.total_orders)}</td>
                                                    <td className="rpt-revenue">{fmt(pm.total_revenue)}</td>
                                                    <td>{fmt(disc.total_shipping_support || 0)}</td>
                                                    <td>{disc.discount_rate_percent ? `${disc.discount_rate_percent}%` : '0%'}</td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            )}
                        </div>
                    </div> */}
                </div>

                {/* ── Right Column ── */}
                <div className="rpt-col">
                    {/* Procurement */}
                    {/* <SectionTitle icon={RefreshCw} title="Chi phí nhập hàng" color="#6366f1" />
                    <div className="admin-card rpt-mb">
                        <div className="admin-table-container">
                            <table className="admin-table">
                                <thead>
                                    <tr>
                                        <th>Nhà cung cấp</th>
                                        <th>Số hóa đơn</th>
                                        <th>Tổng chi</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {procurement.map((s, i) => (
                                        <tr key={i}>
                                            <td className="fw-500">{s.supplier_name}</td>
                                            <td>{s.total_invoices}</td>
                                            <td style={{ color: '#dc2626' }}>{fmt(s.total_spend)}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div> */}

                    {/* Loyal Customers */}
                    <SectionTitle icon={Users} title="Khách hàng tiêu biểu" color="#ec4899" />
                    <div className="admin-card rpt-mb">
                        <div className="admin-table-container">
                            <table className="admin-table">
                                <thead>
                                    <tr>
                                        <th>Khách hàng</th>
                                        <th>Đơn</th>
                                        <th>Tổng chi</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {loyalCustomers.map((c) => (
                                        <tr key={c.id}>
                                            <td>
                                                <div className="rpt-user-cell">
                                                    <span className="fw-600">{c.full_name}</span>
                                                    <span className="txt-muted">{c.email}</span>
                                                </div>
                                            </td>
                                            <td>{c.total_orders}</td>
                                            <td className="rpt-revenue">{fmt(c.total_spent)}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>

                    {/* Demographics */}
                    <SectionTitle icon={MapPin} title="Phân bố khách hàng" color="#14b8a6" />
                    <div className="admin-card">
                        <div className="admin-table-container">
                            {demographics.geoDistribution.length === 0 ? (
                                <div className="rpt-empty">Chưa có dữ liệu vị trí</div>
                            ) : (
                                <table className="admin-table">
                                    <thead>
                                        <tr>
                                            <th>Tỉnh / Thành phố</th>
                                            <th style={{ textAlign: 'center' }}>Khách hàng</th>
                                            <th style={{ textAlign: 'center' }}>Tổng đơn hàng</th>
                                            <th style={{ textAlign: 'center' }}>Thành công</th>
                                            <th style={{ textAlign: 'center' }}>Tỷ lệ giao hàng</th>
                                            <th style={{ textAlign: 'center' }}>Thao tác</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {demographics.geoDistribution.map((loc, i) => {
                                            const total = Number(loc.total_orders || 0);
                                            const completed = Number(loc.completed_orders || 0);
                                            const rate = total > 0 ? ((completed / total) * 100).toFixed(1) : '0.0';
                                            return (
                                                <tr key={i}>
                                                    <td className="fw-600" style={{ color: 'var(--admin-primary)' }}>
                                                        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                                                            <MapPin size={16} style={{ color: '#14b8a6' }} />
                                                            {loc.province}
                                                        </div>
                                                    </td>
                                                    <td style={{ textAlign: 'center' }} className="fw-500">{fmtNum(loc.count)}</td>
                                                    <td style={{ textAlign: 'center' }}>{fmtNum(total)}</td>
                                                    <td style={{ textAlign: 'center', color: 'var(--admin-success)' }}>{fmtNum(completed)}</td>
                                                    <td style={{ textAlign: 'center' }}>
                                                        <span className={`admin-badge ${Number(rate) >= 80 ? 'admin-badge-success' : Number(rate) >= 50 ? 'admin-badge-blue' : 'admin-badge-danger'}`}
                                                              style={Number(rate) >= 80 ? {} : Number(rate) >= 50 ? {} : { background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444' }}>
                                                            {rate}%
                                                        </span>
                                                    </td>
                                                    <td style={{ textAlign: 'center' }}>
                                                        <button
                                                            className="admin-btn admin-btn-outline admin-btn-sm"
                                                            title="Xem danh sách khách hàng"
                                                            onClick={() => handleViewCustomersByProvince(loc.province)}
                                                        >
                                                            <Eye size={14} />
                                                            <span style={{ marginLeft: '4px', fontSize: '0.8rem' }}>Khách hàng</span>
                                                        </button>
                                                    </td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            )}
                        </div>
                    </div>
                </div>
            </div>
            {showProvinceModal && (
                <div className="admin-modal-overlay">
                    <div className="admin-modal" style={{ maxWidth: '750px', width: '95%' }}>
                        <div className="admin-modal-header">
                            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <MapPin size={22} style={{ color: '#14b8a6' }} />
                                <h2>Danh sách khách hàng tại: {selectedProvince}</h2>
                            </div>
                            <button className="admin-btn" onClick={() => setShowProvinceModal(false)}>
                                <X size={20} />
                            </button>
                        </div>
                        <div className="admin-modal-body" style={{ maxHeight: '480px', overflowY: 'auto' }}>
                            {loadingCustomers ? (
                                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', padding: '40px' }}>
                                    <Loader className="spin" size={32} style={{ color: 'var(--admin-primary)', marginBottom: '12px' }} />
                                    <p style={{ marginTop: '8px' }}>Đang tải danh sách khách hàng...</p>
                                </div>
                            ) : provinceCustomers.length === 0 ? (
                                <div style={{ textAlign: 'center', padding: '30px', color: 'var(--admin-text-muted)' }}>
                                    Không tìm thấy khách hàng nào có địa chỉ tại đây.
                                </div>
                            ) : (
                                <div className="admin-table-container">
                                    <table className="admin-table">
                                        <thead>
                                            <tr>
                                                <th>Mã KH</th>
                                                <th>Tên khách hàng</th>
                                                <th>Email / SĐT</th>
                                                <th style={{ textAlign: 'center' }}>Tổng đơn</th>
                                                <th style={{ textAlign: 'center' }}>Thành công</th>
                                                <th style={{ textAlign: 'center' }}>Tỷ lệ mua</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {provinceCustomers.map((cust) => {
                                                const total = Number(cust.total_orders || 0);
                                                const completed = Number(cust.completed_orders || 0);
                                                const rate = total > 0 ? ((completed / total) * 100).toFixed(1) : '0.0';
                                                return (
                                                    <tr key={cust.id}>
                                                        <td style={{ fontWeight: '600', color: 'var(--admin-primary)' }}>#{cust.id}</td>
                                                        <td>
                                                            <div className="fw-600">{cust.full_name}</div>
                                                            <div style={{ fontSize: '0.75rem', color: 'var(--admin-text-muted)' }}>
                                                                Tham gia: {cust.joined_date ? new Date(cust.joined_date).toLocaleDateString('vi-VN') : '—'}
                                                            </div>
                                                        </td>
                                                        <td>
                                                            <div>{cust.email}</div>
                                                            <div style={{ fontSize: '0.8rem', color: 'var(--admin-text-muted)' }}>{cust.num_phone || '—'}</div>
                                                        </td>
                                                        <td style={{ textAlign: 'center' }} className="fw-500">{fmtNum(total)}</td>
                                                        <td style={{ textAlign: 'center', color: 'var(--admin-success)' }} className="fw-600">{fmtNum(completed)}</td>
                                                        <td style={{ textAlign: 'center' }}>
                                                            <span className={`admin-badge ${Number(rate) >= 80 ? 'admin-badge-success' : Number(rate) >= 50 ? 'admin-badge-blue' : 'admin-badge-danger'}`}>
                                                                {rate}%
                                                            </span>
                                                        </td>
                                                    </tr>
                                                );
                                            })}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </div>
                        <div className="admin-modal-footer">
                            <button className="admin-btn admin-btn-outline" onClick={() => setShowProvinceModal(false)}>
                                Đóng
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminReport;
