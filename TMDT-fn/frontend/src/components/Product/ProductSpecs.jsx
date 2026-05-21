import React, { useState } from 'react';
import { ChevronDown, ChevronUp } from 'lucide-react';
import './ProductDetail.css';

const AccordionItem = ({ title, content, initialOpen = false }) => {
    const [isOpen, setIsOpen] = useState(initialOpen);

    const items = (content && typeof content === 'object' && !Array.isArray(content))
        ? Object.entries(content).map(([label, value]) => ({ label, value }))
        : Array.isArray(content)
            ? content.map((v, i) => ({ label: `${title} ${i + 1}`, value: v }))
            : [{ label: title, value: content }];

    if (items.length === 0) return null;

    return (
        <div className={`accordion-item ${isOpen ? 'active' : ''}`}>
            <button className="accordion-header" onClick={() => setIsOpen(!isOpen)}>
                <span>{title}</span>
                {isOpen ? <ChevronUp size={20} className="accordion-icon" /> : <ChevronDown size={20} className="accordion-icon" />}
            </button>
            <div className={`accordion-content ${isOpen ? 'show' : ''}`}>
                <table className="specs-table">
                    <tbody>
                        {items.map((item, index) => (
                            <tr key={index}>
                                <td className="specs-table__label">{item.label}:</td>
                                <td className="specs-table__value">
                                    {Array.isArray(item.value)
                                        ? item.value.map((v, i) => <div key={i}>{v}</div>)
                                        : item.value || 'N/A'}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

const ProductSpecs = ({ specs }) => {
    if (!specs) return null;

    const entries = Object.entries(specs);
    const groupedSections = [];
    const flatSpecs = {};

    entries.forEach(([key, value]) => {
        if (value && typeof value === 'object' && !Array.isArray(value)) {
            groupedSections.push([key, value]);
        } else {
            flatSpecs[key] = value;
        }
    });

    const sections = [...groupedSections];
    if (Object.keys(flatSpecs).length > 0) {
        sections.push(['Thông số kỹ thuật chi tiết', flatSpecs]);
    }

    if (sections.length === 0) return null;

    return (
        <div className="product-specs">
            <div className="specs-accordion">
                {sections.map(([title, content], index) => (
                    <AccordionItem
                        key={index}
                        title={title}
                        content={content}
                        initialOpen={index === 0}
                    />
                ))}
            </div>
        </div>
    );
};

export default ProductSpecs;
