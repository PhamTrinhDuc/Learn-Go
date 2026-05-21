import React, { useState, useEffect, useRef } from 'react';
import { useAuth } from '../../context/AuthContext';
import { io } from 'socket.io-client';
import { Send, MessageSquare, User } from 'lucide-react';
import './AdminChatManager.css';

const API = import.meta.env.VITE_SERVER_API;

const AdminChatManager = () => {
    const { user } = useAuth();
    const [rooms, setRooms] = useState([]);
    const [activeRoom, setActiveRoom] = useState(null);
    const [messages, setMessages] = useState([]);
    const [inputValue, setInputValue] = useState('');

    const globalSocketRef = useRef(null);
    const messagesEndRef = useRef(null);

    // 1. Fetch active rooms list
    const fetchRooms = async () => {
        try {
            const res = await fetch(`${API}/api/chat/rooms`);
            const result = await res.json();
            if (result.success) {
                setRooms(result.data);
            }
        } catch (error) {
            console.error('Error fetching rooms list:', error);
        }
    };

    // 2. Setup global socket to listen to rooms updates
    useEffect(() => {
        fetchRooms();

        const socket = io(API, { transports: ['websocket'] });
        globalSocketRef.current = socket;

        socket.on('chat_rooms_updated', () => {
            fetchRooms();
        });

        return () => {
            if (globalSocketRef.current) {
                globalSocketRef.current.disconnect();
                globalSocketRef.current = null;
            }
        };
    }, []);

    // 3. Setup socket room-joining and message sync when activeRoom changes
    useEffect(() => {
        if (!activeRoom || !globalSocketRef.current) {
            setMessages([]);
            return;
        }

        const setupRoom = async () => {
            try {
                // Fetch initial messages for active room
                const res = await fetch(`${API}/api/chat/messages/${activeRoom.room_id}`);
                const result = await res.json();
                if (result.success) {
                    setMessages(result.data);
                }

                // Join room on Socket
                globalSocketRef.current.emit('join_chat_room', activeRoom.room_id);

                // Mark messages in the room as read
                globalSocketRef.current.emit('mark_messages_read', { room_id: activeRoom.room_id, role: 'admin' });
            } catch (error) {
                console.error('Error opening chat room:', error);
            }
        };

        setupRoom();

        // Bind message listener for the active room
        const socket = globalSocketRef.current;
        const handleNewMessage = (newMsg) => {
            if (newMsg.room_id === activeRoom.room_id) {
                setMessages((prev) => {
                    if (prev.some(m => m.message_id === newMsg.message_id)) return prev;
                    return [...prev, newMsg];
                });
                
                // Immediately mark as read since admin has this room open
                if (newMsg.sender_role === 'customer') {
                    socket.emit('mark_messages_read', { room_id: activeRoom.room_id, role: 'admin' });
                }
            }
        };

        socket.on('new_chat_message', handleNewMessage);

        return () => {
            socket.off('new_chat_message', handleNewMessage);
        };
    }, [activeRoom]);

    // Scroll to bottom helper
    useEffect(() => {
        if (messages.length > 0) {
            scrollToBottom();
        }
    }, [messages]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputValue.trim() || !activeRoom || !globalSocketRef.current) return;

        const messageData = {
            room_id: activeRoom.room_id,
            sender_id: user.id,
            sender_role: 'admin',
            message_text: inputValue.trim()
        };

        globalSocketRef.current.emit('send_chat_message', messageData);
        setInputValue('');
        scrollToBottom();
    };

    const getInitials = (name) => {
        if (!name) return 'KH';
        const parts = name.trim().split(' ');
        if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
        return (parts[parts.length - 2][0] + parts[parts.length - 1][0]).toUpperCase();
    };

    const formatTime = (timestamp) => {
        if (!timestamp) return '';
        try {
            const date = new Date(timestamp);
            return date.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' });
        } catch (e) {
            return '';
        }
    };

    return (
        <div className="admin-chat-container">
            {/* Sidebar list */}
            <div className="admin-chat-sidebar">
                <div className="admin-chat-sidebar-header">
                    <h3>Khách hàng hỗ trợ</h3>
                </div>
                <div className="admin-chat-rooms-list">
                    {rooms.length === 0 ? (
                        <div style={{ textAlign: 'center', color: '#8c8c8c', marginTop: '30px', fontSize: '13px' }}>
                            Không có yêu cầu chat nào.
                        </div>
                    ) : (
                        rooms.map((r) => (
                            <div
                                key={r.room_id}
                                className={`admin-chat-room-item ${activeRoom?.room_id === r.room_id ? 'active' : ''}`}
                                onClick={() => setActiveRoom(r)}
                            >
                                <div className="admin-chat-room-avatar">
                                    {getInitials(r.customer_name)}
                                </div>
                                <div className="admin-chat-room-details">
                                    <div className="admin-chat-room-info">
                                        <span className="admin-chat-room-name">{r.customer_name}</span>
                                        <span className="admin-chat-room-time">{formatTime(r.updated_at)}</span>
                                    </div>
                                    <div className="admin-chat-room-preview-row">
                                        <span className="admin-chat-room-last-msg">
                                            {r.last_message || 'Chưa có tin nhắn'}
                                        </span>
                                        {r.unread_count > 0 && (
                                            <span className="admin-chat-room-badge">{r.unread_count}</span>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>

            {/* Chat Workspace */}
            <div className="admin-chat-area">
                {!activeRoom ? (
                    <div className="admin-chat-empty-state">
                        <MessageSquare size={64} strokeWidth={1.5} />
                        <p>Chọn một cuộc trò chuyện để bắt đầu trả lời khách hàng.</p>
                    </div>
                ) : (
                    <>
                        {/* Header */}
                        <div className="admin-chat-area-header">
                            <div className="admin-chat-header-user">
                                <div className="admin-chat-header-avatar">
                                    {getInitials(activeRoom.customer_name)}
                                </div>
                                <div>
                                    <h4 className="admin-chat-header-name">{activeRoom.customer_name}</h4>
                                    <p className="admin-chat-header-email">{activeRoom.customer_email}</p>
                                </div>
                            </div>
                        </div>

                        {/* Messages panel */}
                        <div className="admin-chat-messages">
                            {messages.length === 0 ? (
                                <div style={{ textAlign: 'center', color: '#8c8c8c', marginTop: '20px' }}>
                                    Chưa có tin nhắn nào. Bắt đầu cuộc trò chuyện.
                                </div>
                            ) : (
                                messages.map((msg, index) => (
                                    <div key={msg.message_id || index} className={`admin-chat-msg-row ${msg.sender_role}`}>
                                        <div className="admin-chat-msg-bubble">
                                            {msg.message_text}
                                            <span className="admin-chat-msg-time">{formatTime(msg.created_at)}</span>
                                        </div>
                                    </div>
                                ))
                            )}
                            <div ref={messagesEndRef} />
                        </div>

                        {/* Input bar */}
                        <form className="admin-chat-input-form" onSubmit={handleSendMessage}>
                            <input
                                type="text"
                                placeholder="Nhập câu trả lời..."
                                value={inputValue}
                                onChange={(e) => setInputValue(e.target.value)}
                                maxLength={1000}
                            />
                            <button className="admin-chat-input-send" type="submit" disabled={!inputValue.trim()}>
                                <span>Gửi</span>
                                <Send size={16} />
                            </button>
                        </form>
                    </>
                )}
            </div>
        </div>
    );
};

export default AdminChatManager;
