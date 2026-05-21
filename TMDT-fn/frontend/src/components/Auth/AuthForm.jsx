import React from 'react';
import '../../styles/Auth.css';

const AuthForm = ({ title, children, onSubmit }) => {
    return (
        <div className="auth-container">
            <div className="auth-form">
                <h2 className="auth-title">{title}</h2>
                <form onSubmit={onSubmit}>
                    {children}
                </form>
            </div>
        </div>
    );
};

export default AuthForm;
