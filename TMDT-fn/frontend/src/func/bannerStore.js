const STORAGE_KEY = 'tgbd_banners_config';

export const defaultBanners = {
    topBanner: {
        active: true,
        imageUrl: "https://cdnv2.tgdd.vn/mwg-static/tgdd/Banner/8e/9a/8e9a5879b5d607b6132b5a8eea9da3dc.png",
        link: "/"
    },
    sliders: [
        { id: 1, imageUrl: "https://cdn.tgdd.vn/2024/02/banner/720-220-720x220-6.png", link: "/" },
        { id: 2, imageUrl: "https://cdn.tgdd.vn/2024/02/banner/720-220-720x220-5.png", link: "/" }
    ]
};

export const getBanners = () => {
    const data = localStorage.getItem(STORAGE_KEY);
    if (!data) return defaultBanners;
    try {
        return JSON.parse(data);
    } catch {
        return defaultBanners;
    }
};

export const saveBanners = (banners) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(banners));
};
