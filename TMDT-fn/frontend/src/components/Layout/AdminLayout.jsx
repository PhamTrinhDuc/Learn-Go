import { Outlet } from 'react-router-dom';
import AdminSidebar from '../Admin/AdminSidebar';
import AdminHeader from '../Admin/AdminHeader';
import '../../styles/admin.css';

const AdminLayout = () => {
    return (
        <div className="admin-layout">
            <AdminSidebar />
            <main className="admin-main">
                <AdminHeader />
                <div className="admin-content">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default AdminLayout;
