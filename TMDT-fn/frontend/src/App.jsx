import { BrowserRouter, Routes, Route } from 'react-router-dom';
import MainLayout from './components/Layout/MainLayout';
import Home from './pages/Home';
import ProductList from './pages/Pagelist/ProductList';
import Login from './pages/Auth/Login';
import Register from './pages/Auth/Register';
import ForgotPassword from './pages/Auth/ForgotPassword';
import ProductDetail from './pages/ProductDetail';
import Cart from './pages/Cart/Cart';
import Checkout from './pages/Checkout/Checkout';
import OrderSuccess from './pages/Checkout/OrderSuccess';
import OrderFailed from './pages/Checkout/OrderFailed';
import PaymentQR from './pages/Checkout/PaymentQR';
import PayOSReturn from './pages/Checkout/PayOSReturn';
import PayOSCancel from './pages/Checkout/PayOSCancel';
import Profile from './pages/Auth/Profile';
import OrderHistory from './pages/Auth/OrderHistory';
import AddressManager from './pages/Auth/AddressManager';
import SearchPage from './pages/SearchPage';
import { AuthProvider } from './context/AuthContext';
import { CartProvider } from './context/CartContext';
import ProtectedRoute from './components/Auth/ProtectedRoute';
import { Toaster } from 'react-hot-toast';

import AdminLayout from './components/Layout/AdminLayout';
import AdminProductManager from './pages/Admin/ProductManager';
import AdminOrderManager from './pages/Admin/OrderManager';
import AdminReport from './pages/Admin/Report';
import AdminUserManager from './pages/Admin/UserManager';
import AdminAddUser from './pages/Admin/AddUser';
import StoreConfigManager from './pages/Admin/StoreConfigManager';
import BannerManager from './pages/Admin/BannerManager';
import AdminDashboard from './pages/Admin/Dashboard';
import VoucherManager from './pages/Admin/VoucherManager';
import PromotionManager from './pages/Admin/PromotionManager';
import AdminReviewManager from './pages/Admin/ReviewManager';
import AdminChatManager from './pages/Admin/AdminChatManager';

function App() {
  return (
    <AuthProvider>
      <CartProvider>
        <Toaster position="top-center" reverseOrder={false} />
        <BrowserRouter>
          <Routes>
          {/* User Routes */}
          <Route path="/" element={<MainLayout />}>
            <Route index element={<Home />} />
            <Route path="search" element={<SearchPage />} />
            <Route path=":slug" element={<ProductList />} />
            <Route path=":slug/:nameId" element={<ProductDetail />} />
            <Route
              path="history"
              element={
                <ProtectedRoute allowedRoles={['customer', 'admin']}>
                  <OrderHistory />
                </ProtectedRoute>
              }
            />
            <Route
              path="profile"
              element={
                <ProtectedRoute allowedRoles={['customer', 'admin']}>
                  <Profile />
                </ProtectedRoute>
              }
            />
            <Route
              path="addresses"
              element={
                <ProtectedRoute allowedRoles={['customer', 'admin']}>
                  <AddressManager />
                </ProtectedRoute>
              }
            />
            <Route path="cart" element={<Cart />} />
            <Route
              path="checkout"
              element={
                <ProtectedRoute allowedRoles={['customer', 'admin']}>
                  <Checkout />
                </ProtectedRoute>
              }
            />
            <Route path="checkout/success" element={<OrderSuccess />} />
            <Route path="checkout/failed" element={<OrderFailed />} />
            <Route
              path="checkout/payment"
              element={
                <ProtectedRoute allowedRoles={['customer', 'admin']}>
                  <PaymentQR />
                </ProtectedRoute>
              }
            />
            <Route path="checkout/payos-return" element={<PayOSReturn />} />
            <Route path="checkout/payos-cancel" element={<PayOSCancel />} />

            {/* Auth Routes */}
            <Route path="login" element={<Login />} />
            <Route path="register" element={<Register />} />
            <Route path="forgot-password" element={<ForgotPassword />} />
          </Route>

          {/* Admin Routes */}
          <Route
            path="/admin"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <AdminLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<AdminDashboard />} />
            <Route path="products" element={<AdminProductManager />} />
            <Route path="orders" element={<AdminOrderManager />} />
            <Route path="reports" element={<AdminReport />} />
            <Route path="users" element={<AdminUserManager />} />
            <Route path="users/add" element={<AdminAddUser />} />
            <Route path="store-config" element={<StoreConfigManager />} />
            <Route path="banners" element={<BannerManager />} />
            <Route path="vouchers" element={<VoucherManager />} />
            <Route path="promotions" element={<PromotionManager />} />
            <Route path="reviews" element={<AdminReviewManager />} />
            <Route path="chat" element={<AdminChatManager />} />
          </Route>
        </Routes>
      </BrowserRouter>
      </CartProvider>
    </AuthProvider>
  );
}

export default App;
