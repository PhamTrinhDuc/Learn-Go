import { Link } from 'react-router-dom';
import './Header.css'; // Re-use header CSS for now since classes are defined there

const SubMenuGroup = ({ title, items, onClose }) => {
    return (
        <div className="submenu-group">
            <h4 className="submenu-group__title">{title}</h4>
            <div className="submenu-group__list">
                {items.map((item, index) => (
                    <Link
                        key={index}
                        to={`/${item.slug || '#'}`}
                        className="submenu-item"
                        onClick={onClose}
                    >
                        {item.label}
                    </Link>
                ))}
            </div>
        </div>
    );
};

export default SubMenuGroup;
