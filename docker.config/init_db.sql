-- ============================================================
-- CHUỖI CẮT TÓC — DATABASE SCHEMA
-- PostgreSQL + pgvector
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================================
-- NHÓM 1: CẤU TRÚC CỬA HÀNG
-- ============================================================

CREATE TABLE branch (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(100) NOT NULL,
    address         TEXT NOT NULL,
    phone           VARCHAR(20),
    opening_hours   TEXT,                      -- e.g. "08:00-20:00"
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE service (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             VARCHAR(100) NOT NULL,
    category         VARCHAR(50),               -- e.g. "cut", "color", "treatment"
    description      TEXT,
    duration_minutes INT NOT NULL DEFAULT 30,
    estimated_duration INT NOT NULL,             -- dùng để tính toán lịch trình, có thể khác duration thực tế
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE branch_service_price (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id    UUID NOT NULL REFERENCES branch(id) ON DELETE CASCADE,
    service_id   UUID NOT NULL REFERENCES service(id) ON DELETE CASCADE,
    price        NUMERIC(10,2) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (branch_id, service_id)
);

CREATE TABLE stylist (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id        UUID NOT NULL REFERENCES branch(id) ON DELETE CASCADE,
    name             VARCHAR(100) NOT NULL,
    phone            VARCHAR(20),
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE stylist_schedule (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stylist_id   UUID NOT NULL REFERENCES stylist(id) ON DELETE CASCADE,
    day_of_week  INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time   TIME NOT NULL,
    end_time     TIME NOT NULL,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (stylist_id, day_of_week),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- NHÓM 2: USER (khách hàng + manager + owner)
-- ============================================================

CREATE TABLE users (
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name               VARCHAR(100) NOT NULL,
    phone              VARCHAR(20) UNIQUE NOT NULL,
    email              VARCHAR(100) UNIQUE,
    username           VARCHAR(50) UNIQUE,
    password_hash      TEXT,
    birthday           DATE,
    address            TEXT,
    role               VARCHAR(20) NOT NULL DEFAULT 'customer'
                           CHECK (role IN ('customer', 'manager', 'owner')),
    loyalty_points     INT NOT NULL DEFAULT 0,
    preferred_branch_id UUID REFERENCES branch(id) ON DELETE SET NULL,
    last_visit_at      TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- NHÓM 3: BOOKING
-- ============================================================

CREATE TABLE booking (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    branch_id        UUID NOT NULL REFERENCES branch(id) ON DELETE CASCADE,
    stylist_id       UUID NOT NULL REFERENCES stylist(id) ON DELETE CASCADE,
    service_id       UUID NOT NULL REFERENCES service(id) ON DELETE CASCADE,
    scheduled_at     TIMESTAMPTZ NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 30,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending','confirmed','completed','cancelled','no_show')),
    cancel_reason    TEXT,
    source           VARCHAR(20)                -- 'zalo', 'web', 'agent', 'manual'
                         CHECK (source IN ('zalo','web','agent','manual')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- NHÓM 4: SẢN PHẨM & TỒN KHO
-- ============================================================

CREATE TABLE product (
    id                          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name                        VARCHAR(100) NOT NULL,
    category                    VARCHAR(50),
    price_in                     NUMERIC(10,2) NOT NULL,
    price_out                    NUMERIC(10,2) NOT NULL,
    usage_type                  VARCHAR(20) NOT NULL DEFAULT 'both'
                                    CHECK (usage_type IN ('internal','retail','both')),
    low_stock_threshold_retail   INT NOT NULL DEFAULT 5,
    low_stock_threshold_internal INT NOT NULL DEFAULT 3,
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE inventory (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id   UUID NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    branch_id    UUID NOT NULL REFERENCES branch(id) ON DELETE CASCADE,
    quantity_total    INT NOT NULL DEFAULT 0 CHECK (quantity_total >= 0),
    quantity_retail   INT NOT NULL DEFAULT 0 CHECK (quantity_retail >= 0),
    quantity_internal INT NOT NULL DEFAULT 0 CHECK (quantity_internal >= 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (product_id, branch_id)
);

CREATE TABLE inventory_log (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    action_type  VARCHAR(30) NOT NULL
                     CHECK (action_type IN ('import','sale','internal_use','adjustment')),
    qty_change   INT NOT NULL,                 -- âm = giảm, dương = tăng
    note         TEXT,
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    performer_role VARCHAR(20),               -- 'manager', 'agent'
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- NHÓM 5: ĐƠN HÀNG BÁN LẺ
-- ============================================================

CREATE TABLE orders (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    branch_id      UUID NOT NULL REFERENCES branch(id) ON DELETE CASCADE,
    total_amount   NUMERIC(10,2) NOT NULL DEFAULT 0,
    points_earned  INT NOT NULL DEFAULT 0,
    payment_status BOOLEAN NOT NULL DEFAULT FALSE,
    payment_method VARCHAR(30),               -- 'cash', 'transfer', 'momo'
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id   UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES product(id) ON DELETE CASCADE,
    quantity   INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price NUMERIC(10,2) NOT NULL
);

-- ============================================================
-- NHÓM 6: LOYALTY
-- ============================================================

CREATE TABLE loyalty_transaction (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        VARCHAR(20) NOT NULL
                    CHECK (type IN ('earn','redeem','expire','manual')),
    points      INT NOT NULL,
    ref_type    VARCHAR(20)
                    CHECK (ref_type IN ('booking','orders')),
    ref_id      UUID,                          -- trỏ tới booking.id hoặc order.id
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- NHÓM 7: CONVERSATION & MESSAGES
-- ============================================================

CREATE TABLE conversation (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE messages (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversation(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL
                        CHECK (role IN ('user','assistant','system')),
    content         TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- NHÓM 8: VẬN HÀNH HỆ THỐNG
-- ============================================================

CREATE TABLE knowledge_base (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id  UUID REFERENCES branch(id) ON DELETE CASCADE,  -- NULL = dùng chung
    title      VARCHAR(200) NOT NULL,
    content    TEXT NOT NULL,
    embedding  vector(1024),                   -- OpenAI / Anthropic embedding dim
    metadata   JSONB,
    category   VARCHAR(50),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notify_log (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel    VARCHAR(20) NOT NULL
                   CHECK (channel IN ('zalo','sms','slack')),
    type       VARCHAR(30) NOT NULL,           -- 'reminder', 'birthday', 'loyalty', 'reactivation'
    status     VARCHAR(20) NOT NULL DEFAULT 'sent'
                   CHECK (status IN ('sent','failed')),
    sent_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE agent_action_log (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users(id) ON DELETE SET NULL,
    agent_name  VARCHAR(50) NOT NULL,          -- 'booking', 'faq', 'inventory', 'loyalty', 'analytics'
    action_type VARCHAR(50) NOT NULL,
    payload     JSONB,
    status      VARCHAR(30) NOT NULL DEFAULT 'executed'
                    CHECK (status IN ('executed','pending_approval','approved','rejected')),
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- INDEXES
-- ============================================================

-- User
CREATE UNIQUE INDEX idx_user_phone ON users(phone);

-- Booking
CREATE INDEX idx_booking_stylist_time  ON booking(stylist_id, scheduled_at);
CREATE INDEX idx_booking_branch_time   ON booking(branch_id, scheduled_at);
CREATE INDEX idx_booking_customer      ON booking(user_id);
CREATE INDEX idx_booking_status        ON booking(status);

-- Inventory
CREATE INDEX idx_inventory_branch      ON inventory(branch_id);
CREATE INDEX idx_inventory_product     ON inventory(product_id);

-- Conversation & Messages
CREATE INDEX idx_conversation_user     ON conversation(user_id);
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at);

-- Loyalty
CREATE INDEX idx_loyalty_customer      ON loyalty_transaction(user_id);

-- Notify (tránh spam)
CREATE INDEX idx_notify_user_type      ON notify_log(user_id, type, sent_at);

-- Agent action log
CREATE INDEX idx_agent_log_status      ON agent_action_log(status)
    WHERE status = 'pending_approval';

-- Knowledge base vector search
CREATE INDEX idx_kb_embedding          ON knowledge_base
    USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);