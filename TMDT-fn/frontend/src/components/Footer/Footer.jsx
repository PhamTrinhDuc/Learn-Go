import React from 'react';
import './Footer.css';

const Footer = () => {
    return (
        <footer className="footer">
            <div className="container footer__container">
                <div className="footer__column">
                    <h3 className="footer__heading">Thông tin liên hệ</h3>
                    <ul className="footer__list">
                        <li><a href="#">Giới thiệu công ty</a></li>
                        <li><a href="#">Tuyển dụng</a></li>
                        <li><a href="#">Gửi góp ý, khiếu nại</a></li>
                        <li><a href="#">Tìm siêu thị</a></li>
                    </ul>
                </div>
                <div className="footer__column">
                    <h3 className="footer__heading">Hỗ trợ khách hàng</h3>
                    <ul className="footer__list">
                        <li><a href="#">Lịch sử mua hàng</a></li>
                        <li><a href="#">Hướng dẫn mua hàng</a></li>
                        <li><a href="#">Thanh toán</a></li>
                        <li><a href="#">Hóa đơn điện tử</a></li>
                    </ul>
                </div>
                <div className="footer__column">
                    <h3 className="footer__heading">Chính sách</h3>
                    <ul className="footer__list">
                        <li><a href="#">Chính sách bảo hành</a></li>
                        <li><a href="#">Chính sách đổi trả</a></li>
                        <li><a href="#">Giao hàng & lắp đặt</a></li>
                        <li><a href="#">Bảo mật thông tin</a></li>
                    </ul>
                </div>
                <div className="footer__column">
                    <h3 className="footer__heading">Kết nối với chúng tôi</h3>
                    <div className="footer__socials">
                        {/* Placeholders for social icons */}
                        <span>Facebook</span>
                        <span>Youtube</span>
                        <span>Zalo</span>
                    </div>
                </div>
            </div>
            <div className="footer__copyright">
                <div className="container">
                    <p>© 2026 Thegioibatdong. All rights reserved.</p>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
