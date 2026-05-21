import React from 'react';

const SortIcon = ({ activeKey, columnKey, direction }) => {
    const isActive = activeKey === columnKey;

    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            fontSize: '10px',
            lineHeight: '1',
            marginLeft: '8px',
            userSelect: 'none'
        }}>
            <span style={{
                color: (isActive && direction === 'asc') ? 'var(--admin-primary)' : '#cbd5e1'
            }}>▲</span>
            <span style={{
                color: (isActive && direction === 'desc') ? 'var(--admin-primary)' : '#cbd5e1',
                marginTop: '-4px'
            }}>▼</span>
        </div>
    );
};

export default SortIcon;
