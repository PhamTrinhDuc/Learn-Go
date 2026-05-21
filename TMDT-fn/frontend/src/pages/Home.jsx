import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ProductCard from '../components/Product/ProductCard';
import { ChevronRight, Loader, Smartphone, Laptop, Headphones, Tablet, Zap } from 'lucide-react';
import { normalizeProductData } from '../func/productHelpers';
import { getBanners } from '../func/bannerStore';
import { Swiper, SwiperSlide } from 'swiper/react';
import { Autoplay, Pagination, Navigation } from 'swiper/modules';
import 'swiper/css';
import 'swiper/css/pagination';
import 'swiper/css/navigation';
import './Home.css';

const IconMap = {
    Smartphone: <Smartphone size={18} />,
    Laptop: <Laptop size={18} />,
    Headphones: <Headphones size={18} />,
    Tablet: <Tablet size={18} />,
    Zap: <Zap size={18} />
};

const Home = () => {
    const [activeTab, setActiveTab] = useState('dtdd');
    const [products, setProducts] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [sliderBanners, setSliderBanners] = useState(getBanners().sliders);
    const navigate = useNavigate();
    const [column, setColumn] = useState('');

    useEffect(() => {
        const handleBannerChange = () => {
            setSliderBanners(getBanners().sliders);
        };
        const handleStorageChange = (e) => {
            if (e.key === 'tgbd_banners_config' || !e.key) {
                handleBannerChange();
            }
        };
        window.addEventListener('bannerConfigChanged', handleBannerChange);
        window.addEventListener('storage', handleStorageChange);
        return () => {
            window.removeEventListener('bannerConfigChanged', handleBannerChange);
            window.removeEventListener('storage', handleStorageChange);
        };
    }, []);

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            try {
                const baseUrl = import.meta.env.VITE_SERVER_API;

                const [prodRes, catRes] = await Promise.all([
                    fetch(`${baseUrl}/api/product/all-products`),//?limit=10
                    fetch(`${baseUrl}/api/product/category`)
                ]);

                const prodResult = await prodRes.json();
                const catResult = await catRes.json();

                console.log(prodResult.column);
                setColumn(prodResult.column);

                if (catResult.success) {
                    setCategories(catResult.data);
                }

                if (prodResult.success) {
                    const groupedByCategory = {};

                    prodResult.data.forEach(categoryBlock => {
                        const rootId = categoryBlock.category_id;

                        const processedCards = categoryBlock.products.map(rawProd => {
                            const versionsArray = rawProd.versions || [rawProd];
                            const normResult = normalizeProductData(versionsArray);
                            return normResult.length > 0 ? normResult[0] : null;
                        }).filter(Boolean);

                        groupedByCategory[rootId] = processedCards;
                    });

                    setProducts(groupedByCategory);

                    if (prodResult.data.length > 0 && !activeTab) {
                        setActiveTab(prodResult.data[0].category_id);
                    }
                }
            } catch (error) {
                console.error('Error fetching data:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    const activeCategory = categories.find(c => c.slug === activeTab);

    const getFilteredProducts = () => {
        if (!activeCategory || !products || !products[activeTab]) return [];
        return products[activeTab]; //.slice(0, 10)
    };

    const filteredProducts = getFilteredProducts();

    return (
        <div className="home-page">
            <div className="container">

                <div style={{ padding: '20px 0' }}>
                    {sliderBanners && sliderBanners.length > 0 ? (
                        <Swiper
                            modules={[Autoplay, Pagination, Navigation]}
                            className="home-banner-slider"
                            spaceBetween={20}
                            slidesPerView={1}
                            navigation
                            pagination={{ clickable: true }}
                            autoplay={{ delay: 3500, disableOnInteraction: false }}
                            loop={true}
                            style={{ borderRadius: '12px', overflow: 'hidden' }}
                        >
                            {sliderBanners.map((slide, idx) => (
                                <SwiperSlide key={slide.id || idx}>
                                    <a href={slide.link || '#'}>
                                        <img src={slide.imageUrl} style={{ width: '100%', display: 'block' }} alt={`Banner ${idx}`} />
                                    </a>
                                </SwiperSlide>
                            ))}
                        </Swiper>
                    ) : (
                        <div style={{ padding: '40px', background: '#f8f9fa', borderRadius: '12px', textAlign: 'center', color: '#666' }}>
                            Không có hình ảnh quảng cáo.
                        </div>
                    )}
                </div>

                <section className="home-promotions">
                    <div className="promotions-header">
                        <div className="home-tabs">
                            {categories.map(tab => (
                                <button
                                    key={tab.slug}
                                    className={`home-tab tab-category ${activeTab === tab.slug ? 'active' : ''}`}
                                    onClick={() => setActiveTab(tab.slug)}
                                >
                                    {tab.icon && IconMap[tab.icon] ? (
                                        <span className="tab-icon-wrapper">{IconMap[tab.icon]}</span>
                                    ) : null}
                                    <span className="tab-label">{tab.label}</span>
                                </button>
                            ))}
                        </div>
                    </div>

                    <div className="promotions-content-box">
                        {loading ? (
                            <div className="loading-container" style={{ padding: '50px', textAlign: 'center' }}>
                                <Loader className="spin" size={32} style={{ color: 'var(--primary-color)' }} />
                                <p style={{ marginTop: '10px' }}>Đang tải dữ liệu...</p>
                            </div>
                        ) : (
                            <>
                                <div className="home-product-grid" style={{
                                    "display": "grid",
                                    "grid-template-columns": `repeat(${column}, 1fr)`,
                                    "gap": "10px"
                                }}>
                                    {filteredProducts.map(product => (
                                        <ProductCard key={product.id} product={product} />
                                    ))}
                                    {filteredProducts.length === 0 && (
                                        <div style={{ gridColumn: '1 / -1', textAlign: 'center', padding: '40px', color: '#666' }}>
                                            Không có sản phẩm nào cho danh mục này.
                                        </div>
                                    )}
                                </div>

                                <div className="home-see-more-container">
                                    <button
                                        className="btn-see-more-white"
                                        onClick={() => navigate(`/${activeTab}`)}
                                    >
                                        Xem thêm {categories.find(t => t.slug === activeTab)?.label}
                                        <ChevronRight size={16} />
                                    </button>
                                </div>
                            </>
                        )}
                    </div>
                </section>
            </div>
        </div>
    );
};

export default Home;
