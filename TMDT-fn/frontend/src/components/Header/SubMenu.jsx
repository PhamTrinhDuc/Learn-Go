import React from 'react';
import SubMenuGroup from './SubMenuGroup';
import './Header.css';

const SubMenu = ({ groups, onClose }) => {
    return (
        <div className="submenu-dropdown">
            <div className="submenu-dropdown__content">
                {groups.map((group, idx) => (
                    <SubMenuGroup
                        key={idx}
                        title={group.title}
                        items={group.items}
                        onClose={onClose}
                    />
                ))}
            </div>
        </div>
    );
};

export default SubMenu;
