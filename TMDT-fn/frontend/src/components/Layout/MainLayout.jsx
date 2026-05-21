import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from '../Header/Header';
import Footer from '../Footer/Footer';
import ChatWidget from '../Chat/ChatWidget';
import './MainLayout.css';

const MainLayout = () => {
    return (
        <div className="main-layout">
            <Header />
            <main className="main-layout__content">
                <Outlet />
            </main>
            <Footer />
            <ChatWidget />
        </div>
    );
};

export default MainLayout;
