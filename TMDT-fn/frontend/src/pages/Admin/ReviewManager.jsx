import React, { useState, useEffect } from 'react';
import { Check, Trash2, MessageSquare, Star } from 'lucide-react';
import toast from 'react-hot-toast';
import './ReviewManager.css';

const API = import.meta.env.VITE_SERVER_API;

const ReviewManager = () => {
    const [reviews, setReviews] = useState([]);
    const [loading, setLoading] = useState(true);
    const [filterStatus, setFilterStatus] = useState('all'); // 'all', 'pending', 'approved'

    useEffect(() => {
        fetchReviews();
    }, []);

    const fetchReviews = async () => {
        setLoading(true);
        try {
            const res = await fetch(`${API}/api/product/reviews/admin/all`);
            const result = await res.json();
            if (result.success) {
                setReviews(result.data);
            } else {
                toast.error(result.message || 'Lỗi tải danh sách bình luận');
            }
        } catch (error) {
            console.error('Error fetching reviews:', error);
            toast.error('Lỗi kết nối máy chủ');
        } finally {
            setLoading(false);
        }
    };

    const handleApprove = async (reviewId) => {
        try {
            const res = await fetch(`${API}/api/product/reviews/${reviewId}/approve`, {
                method: 'PATCH'
            });
            const result = await res.json();
            if (result.success) {
                toast.success('Duyệt bình luận thành công!');
                setReviews(prev => prev.map(r => r.review_id === reviewId ? { ...r, status: 'approved' } : r));
            } else {
                toast.error(result.message || 'Duyệt thất bại');
            }
        } catch (error) {
            console.error('Error approving review:', error);
            toast.error('Lỗi kết nối máy chủ');
        }
    };

    const handleDelete = async (reviewId) => {
        if (!window.confirm('Bạn có chắc chắn muốn xóa bình luận này?')) return;

        try {
            const res = await fetch(`${API}/api/product/reviews/${reviewId}`, {
                method: 'DELETE'
            });
            const result = await res.json();
            if (result.success) {
                toast.success('Xóa bình luận thành công!');
                setReviews(prev => prev.filter(r => r.review_id !== reviewId));
            } else {
                toast.error(result.message || 'Xóa thất bại');
            }
        } catch (error) {
            console.error('Error deleting review:', error);
            toast.error('Lỗi kết nối máy chủ');
        }
    };

    const filteredReviews = reviews.filter(r => {
        if (filterStatus === 'all') return true;
        return r.status === filterStatus;
    });

    const formatDate = (dateStr) => {
        return new Date(dateStr).toLocaleDateString('vi-VN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <div className="review-manager-page">
            <div className="manager-header">
                <div className="title-section">
                    <h1>Quản lý bình luận</h1>
                    <p>Phê duyệt hoặc xóa các đánh giá sản phẩm của khách hàng.</p>
                </div>
            </div>

            <div className="filter-tabs">
                <button 
                    className={`tab-btn ${filterStatus === 'all' ? 'active' : ''}`}
                    onClick={() => setFilterStatus('all')}
                >
                    Tất cả ({reviews.length})
                </button>
                <button 
                    className={`tab-btn ${filterStatus === 'pending' ? 'active' : ''}`}
                    onClick={() => setFilterStatus('pending')}
                >
                    Chờ duyệt ({reviews.filter(r => r.status === 'pending').length})
                </button>
                <button 
                    className={`tab-btn ${filterStatus === 'approved' ? 'active' : ''}`}
                    onClick={() => setFilterStatus('approved')}
                >
                    Đã duyệt ({reviews.filter(r => r.status === 'approved').length})
                </button>
            </div>

            <div className="manager-card">
                {loading ? (
                    <div className="loading-state" style={{ padding: '40px', textAlignment: 'center', color: '#666' }}>Đang tải danh sách bình luận...</div>
                ) : filteredReviews.length === 0 ? (
                    <div className="empty-state" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', padding: '40px', color: '#888' }}>
                        <MessageSquare size={48} opacity={0.2} style={{ marginBottom: '10px' }} />
                        <p>Không tìm thấy bình luận nào</p>
                    </div>
                ) : (
                    <div className="table-responsive">
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th>Đơn hàng</th>
                                    <th>Khách hàng</th>
                                    <th>Sản phẩm</th>
                                    <th>Đánh giá</th>
                                    <th style={{ width: '35%' }}>Bình luận</th>
                                    <th>Thời gian</th>
                                    <th>Trạng thái</th>
                                    <th>Hành động</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filteredReviews.map((rev) => (
                                    <tr key={rev.review_id}>
                                        <td>#{rev.order_id}</td>
                                        <td className="font-medium">{rev.customer_name}</td>
                                        <td>{rev.product_name}</td>
                                        <td>
                                            <div className="star-display">
                                                {[1, 2, 3, 4, 5].map(i => (
                                                    <Star 
                                                        key={i} 
                                                        size={14} 
                                                        fill={i <= rev.rating ? "#fb6e2e" : "none"} 
                                                        stroke={i <= rev.rating ? "#fb6e2e" : "#ccc"} 
                                                    />
                                                ))}
                                            </div>
                                        </td>
                                        <td className="comment-cell">{rev.comment || <em style={{ color: '#aaa' }}>Không có bình luận</em>}</td>
                                        <td>{formatDate(rev.created_at)}</td>
                                        <td>
                                            <span className={`status-tag ${rev.status}`}>
                                                {rev.status === 'approved' ? 'Đã duyệt' : 'Chờ duyệt'}
                                            </span>
                                        </td>
                                        <td>
                                            <div className="action-buttons">
                                                {rev.status === 'pending' && (
                                                    <button 
                                                        className="btn-action approve" 
                                                        onClick={() => handleApprove(rev.review_id)}
                                                        title="Phê duyệt bình luận"
                                                    >
                                                        <Check size={14} /> Duyệt
                                                    </button>
                                                )}
                                                <button 
                                                    className="btn-action delete" 
                                                    onClick={() => handleDelete(rev.review_id)}
                                                    title="Xóa bình luận"
                                                >
                                                    <Trash2 size={14} /> Xóa
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ReviewManager;
