import React from 'react';

const AdminHeader = () => {
    return (
        <header className="admin-header">
            <div className="admin-header-title">
                <h1>Bảng điều khiển</h1>
            </div>
            <div className="admin-user-info">
                <span className="admin-user-name">Administrators</span>
                <div className="admin-avatar">A</div>
            </div>
        </header>
    );
};

export default AdminHeader;
