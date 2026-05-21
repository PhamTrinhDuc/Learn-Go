import React, { createContext, useState, useEffect, useContext } from 'react';
import toast from 'react-hot-toast';
import { useAuth } from './AuthContext';

const CartContext = createContext();
const API = import.meta.env.VITE_SERVER_API;

export const useCart = () => useContext(CartContext);

export const CartProvider = ({ children }) => {
    const { user, isAuthenticated } = useAuth();
    const [cartItems, setCartItems] = useState([]);
    const [loading, setLoading] = useState(false);

    const userId = user ? Number(user.id || user.user_id) : null;

    // Helper to fetch cart from database
    const fetchCartFromDB = async (uid) => {
        try {
            setLoading(true);
            const res = await fetch(`${API}/api/cart/${uid}`);
            const data = await res.json();
            if (data.success) {
                setCartItems(data.cart);
            }
        } catch (err) {
            console.error("Failed to fetch cart from DB:", err);
        } finally {
            setLoading(false);
        }
    };

    // 1. Initial cart load
    useEffect(() => {
        if (isAuthenticated && userId) {
            // Logged in -> Load from DB
            fetchCartFromDB(userId);
        } else {
            // Not logged in -> Load from localStorage
            const savedCart = localStorage.getItem('cart');
            setCartItems(savedCart ? JSON.parse(savedCart) : []);
        }
    }, [isAuthenticated, userId]);

    // 2. Sync localStorage cart to DB when user signs in
    useEffect(() => {
        const syncCartOnLogin = async () => {
            if (isAuthenticated && userId) {
                const localCart = localStorage.getItem('cart');
                if (localCart) {
                    try {
                        const parsedCart = JSON.parse(localCart);
                        if (parsedCart.length > 0) {
                            const res = await fetch(`${API}/api/cart/sync`, {
                                method: 'POST',
                                headers: { 'Content-Type': 'application/json' },
                                body: JSON.stringify({
                                    user_id: userId,
                                    items: parsedCart
                                })
                            });
                            const data = await res.json();
                            if (data.success) {
                                localStorage.removeItem('cart');
                            }
                        }
                    } catch (err) {
                        console.error("Failed to sync cart on login:", err);
                    }
                }
                // Always fetch fresh cart from DB after login/sync
                fetchCartFromDB(userId);
            }
        };

        syncCartOnLogin();
    }, [isAuthenticated, userId]);

    // 3. Write to localStorage only when user is not authenticated
    useEffect(() => {
        if (!isAuthenticated) {
            localStorage.setItem('cart', JSON.stringify(cartItems));
        }
    }, [cartItems, isAuthenticated]);

    // 4. Add to cart
    const addToCart = async (product, options = {}) => {
        const variantId = Number(options.variant_id || product.variant_id || product.id || product.product_id);
        if (isNaN(variantId)) {
            toast.error("Không tìm thấy biến thể sản phẩm hợp lệ!");
            return;
        }

        if (isAuthenticated && userId) {
            try {
                const res = await fetch(`${API}/api/cart/add`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        user_id: userId,
                        variant_id: variantId,
                        quantity: 1
                    })
                });
                const data = await res.json();
                if (data.success) {
                    await fetchCartFromDB(userId);
                } else {
                    toast.error(data.error || "Lỗi thêm giỏ hàng");
                }
            } catch (err) {
                console.error(err);
                toast.error("Lỗi kết nối tới server.");
            }
        } else {
            // Local fallback logic
            setCartItems(prev => {
                const existingItem = prev.find(item => {
                    const itemId = item.variant_id || item.product?.id || item.product?.product_id;
                    return Number(itemId) === variantId;
                });

                if (existingItem) {
                    return prev.map(item => {
                        const itemId = item.variant_id || item.product?.id || item.product?.product_id;
                        return Number(itemId) === variantId
                            ? { ...item, quantity: item.quantity + 1 }
                            : item;
                    });
                } else {
                    return [...prev, { product, quantity: 1, variant_id: variantId, ...options }];
                }
            });
        }
    };

    // 5. Remove from cart
    const removeFromCart = async (index) => {
        const targetItem = cartItems[index];
        if (!targetItem) return;

        const variantId = Number(targetItem.variant_id);

        if (isAuthenticated && userId && !isNaN(variantId)) {
            try {
                const res = await fetch(`${API}/api/cart/remove`, {
                    method: 'DELETE',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        user_id: userId,
                        variant_id: variantId
                    })
                });
                const data = await res.json();
                if (data.success) {
                    await fetchCartFromDB(userId);
                } else {
                    toast.error(data.error || "Lỗi xóa sản phẩm");
                }
            } catch (err) {
                console.error(err);
                toast.error("Lỗi kết nối server.");
            }
        } else {
            setCartItems(prev => prev.filter((_, i) => i !== index));
        }
    };

    // 6. Remove multiple items from cart
    const removeMultipleFromCart = async (indices) => {
        if (isAuthenticated && userId) {
            try {
                const promises = indices.map(idx => {
                    const item = cartItems[idx];
                    if (item && item.variant_id) {
                        return fetch(`${API}/api/cart/remove`, {
                            method: 'DELETE',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({
                                user_id: userId,
                                variant_id: Number(item.variant_id)
                            })
                        });
                    }
                    return Promise.resolve();
                });
                await Promise.all(promises);
                await fetchCartFromDB(userId);
            } catch (err) {
                console.error("Error batch removing:", err);
            }
        } else {
            setCartItems(prev => prev.filter((_, i) => !indices.includes(i)));
        }
    };

    // 7. Update quantity using Delta (+1 / -1)
    const updateQuantity = async (index, delta) => {
        const targetItem = cartItems[index];
        if (!targetItem) return;

        const variantId = Number(targetItem.variant_id);
        const newQty = targetItem.quantity + delta;
        if (newQty < 1) return;

        if (isAuthenticated && userId && !isNaN(variantId)) {
            try {
                const res = await fetch(`${API}/api/cart/update`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        user_id: userId,
                        variant_id: variantId,
                        quantity: newQty
                    })
                });
                const data = await res.json();
                if (data.success) {
                    await fetchCartFromDB(userId);
                } else {
                    toast.error(data.error || "Lỗi cập nhật số lượng");
                }
            } catch (err) {
                console.error(err);
                toast.error("Lỗi kết nối server.");
            }
        } else {
            setCartItems(prev => prev.map((item, i) => i === index ? { ...item, quantity: newQty } : item));
        }
    };

    // 8. Set exact quantity
    const setQuantity = async (index, value) => {
        const targetItem = cartItems[index];
        if (!targetItem) return;

        const variantId = Number(targetItem.variant_id);
        let val = parseInt(value);
        if (isNaN(val) || val < 1) val = 1;

        if (isAuthenticated && userId && !isNaN(variantId)) {
            try {
                const res = await fetch(`${API}/api/cart/update`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        user_id: userId,
                        variant_id: variantId,
                        quantity: val
                    })
                });
                const data = await res.json();
                if (data.success) {
                    await fetchCartFromDB(userId);
                } else {
                    toast.error(data.error || "Lỗi cập nhật số lượng");
                }
            } catch (err) {
                console.error(err);
                toast.error("Lỗi kết nối server.");
            }
        } else {
            setCartItems(prev => prev.map((item, i) => i === index ? { ...item, quantity: val } : item));
        }
    };

    // 9. Clear entire cart
    const clearCart = async () => {
        if (isAuthenticated && userId) {
            try {
                const res = await fetch(`${API}/api/cart/clear/${userId}`, {
                    method: 'DELETE'
                });
                const data = await res.json();
                if (data.success) {
                    setCartItems([]);
                }
            } catch (err) {
                console.error("Error clearing cart from DB:", err);
            }
        } else {
            setCartItems([]);
        }
    };

    // 10. Update variant options (change color/capacity -> swap variant_id)
    const updateCartItemOptions = async (index, newOptions) => {
        const targetItem = cartItems[index];
        if (!targetItem) return;

        if (isAuthenticated && userId) {
            const oldVariantId = Number(targetItem.variant_id);
            const newVariantId = Number(newOptions.variant_id);

            if (oldVariantId === newVariantId || isNaN(newVariantId)) return;

            try {
                // Remove old variant, add new one with current quantity
                await fetch(`${API}/api/cart/remove`, {
                    method: 'DELETE',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ user_id: userId, variant_id: oldVariantId })
                });

                await fetch(`${API}/api/cart/add`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ user_id: userId, variant_id: newVariantId, quantity: targetItem.quantity })
                });

                await fetchCartFromDB(userId);
            } catch (err) {
                console.error("Error updating cart option in DB:", err);
            }
        } else {
            setCartItems(prev => prev.map((item, i) => i === index ? { ...item, ...newOptions } : item));
        }
    };

    const getCartTotal = () => {
        return cartItems.reduce((total, item) => {
            const price = item.price || item.product?.calculated_price || item.product?.price || 0;
            return total + (price * item.quantity);
        }, 0);
    };

    const getCartCount = () => {
        return cartItems.reduce((total, item) => total + item.quantity, 0);
    };

    return (
        <CartContext.Provider value={{
            cartItems,
            addToCart,
            removeFromCart,
            removeMultipleFromCart,
            updateQuantity,
            setQuantity,
            updateCartItemOptions,
            clearCart,
            getCartTotal,
            getCartCount,
            loading
        }}>
            {children}
        </CartContext.Provider>
    );
};
