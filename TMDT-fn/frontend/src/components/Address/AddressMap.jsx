import { MapContainer, TileLayer, Marker, useMap, useMapEvents } from 'react-leaflet';
import { useState, useEffect } from 'react';
import { Loader, MapPin } from 'lucide-react';
import toast from 'react-hot-toast';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
});

function FlyToPosition({ position }) {
    const map = useMap();
    useEffect(() => {
        if (position) {
            map.flyTo(position, 16, { duration: 1.2 });
        }
    }, [position, map]);
    return null;
}

function LocationPicker({ position, setPosition, onAddressFound, API }) {
    useMapEvents({
        async click(e) {
            const { lat, lng } = e.latlng;
            setPosition([lat, lng]);
            try {
                const res = await fetch(`${API}/api/address/reverse-geocode?lat=${lat}&lon=${lng}`);
                const data = await res.json();
                if (data && onAddressFound) {
                    onAddressFound(data);
                }
            } catch (error) {
                console.error('Lỗi gọi API định vị:', error);
            }
        },
    });
    return position ? <Marker position={position} /> : null;
}
const AddressMap = ({ onAddressFound, API, provinces }) => {
    const [position, setPosition] = useState(null);
    const [isLocating, setIsLocating] = useState(false);

    const handleGetCurrentLocation = () => {
        if (!navigator.geolocation) {
            toast.error("Trình duyệt hoặc giao thức kết nối (HTTP) không hỗ trợ định vị.", { id: 'map-toast' });
            return;
        }

        setIsLocating(true);
        navigator.geolocation.getCurrentPosition(
            async (pos) => {
                const [lat, lng] = [pos.coords.latitude, pos.coords.longitude];
                setPosition([lat, lng]);
                try {
                    const response = await fetch(`${API}/api/address/reverse-geocode?lat=${lat}&lon=${lng}`);
                    const data = await response.json();
                    if (data && onAddressFound) {
                        onAddressFound(data);
                        toast.success('Đã định vị và cập nhật địa chỉ!', { id: 'map-toast' });
                    } else {
                        toast.success('Đã định vị thành công!', { id: 'map-toast' });
                    }
                } catch (e) {
                    console.error('error:', e);
                    toast.error('Không thể lấy vị trí hiện tại', { id: 'map-toast' });
                } finally {
                    setIsLocating(false);
                }
            },
            (error) => {
                // error.code mapping: 1 (PERMISSION_DENIED), 2 (POSITION_UNAVAILABLE), 3 (TIMEOUT)
                if (error.code === 1) {
                    toast.error('Vui lòng bật quyền truy cập vị trí trên trình duyệt.', { id: 'map-toast' });
                } else if (error.code === 3) {
                    toast.error('Thời gian chờ tìm vị trí quá lâu, vui lòng thử lại.', { id: 'map-toast' });
                } else {
                    toast.error('Không thể lấy vị trí hiện tại.', { id: 'map-toast' });
                }
                setIsLocating(false);
            },
            { enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }
        );
    };

    return (
        <div className="address-map-wrapper">
            <div className="address-map-header">
                <span className="address-map-label">Vị trí Bản đồ <span className="address-map-optional">(Tùy chọn)</span></span>
                <button
                    type="button"
                    className="btn-locate"
                    onClick={handleGetCurrentLocation}
                    disabled={isLocating}
                >
                    {isLocating ? <Loader className="spin" size={14} /> : <MapPin size={14} />}
                    {isLocating ? 'Đang định vị...' : 'Định vị hiện tại'}
                </button>
            </div>
            <p className="address-map-hint">Nhấp vào bản đồ để ghim vị trí và tự điền địa chỉ.</p>
            <div className="address-map-container">
                <MapContainer
                    center={position || [10.762622, 106.660172]}
                    zoom={13}
                    style={{ height: '100%', width: '100%' }}
                >
                    <TileLayer
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    <FlyToPosition position={position} />
                    <LocationPicker
                        position={position}
                        setPosition={setPosition}
                        onAddressFound={onAddressFound}
                        API={API}
                    />
                </MapContainer>
            </div>

        </div>
    );
}

export default AddressMap;