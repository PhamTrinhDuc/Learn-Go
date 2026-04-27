# User Flow - Complete Implementation

## Overview
Implemented a complete CRUD user management system with JWT authentication and role-based access control for the salon management backend.

## Project Structure

### 1. Domain Layer (`domain/user.go`)
**Entities & Interfaces:**

```go
// User Entity
type User struct {
    ID                uuid.UUID   // Unique identifier
    Email             string      // User email (unique)
    Password          string      // Hashed password (not exposed in JSON)
    Name              string      // User's full name
    Phone             string      // User's phone number
    Birthday          *time.Time  // Optional birthday
    Address           *string     // Optional address
    Role              UserRole    // User role (admin, manager, stylist, customer)
    LoyaltyPoints     int         // Accumulated loyalty points
    PreferredBranchID *uuid.UUID  // Optional preferred branch
    LastVisitAt       *time.Time  // Last visit timestamp
    IsActive          bool        // Soft delete flag
    CreatedAt         time.Time   // Creation timestamp
    UpdatedAt         time.Time   // Last update timestamp
}

// User Roles
const (
    RoleAdmin    UserRole = "admin"      // Full system access
    RoleManager  UserRole = "manager"    // Branch management
    RoleStylist  UserRole = "stylist"    // Service provider
    RoleCustomer UserRole = "customer"   // Customer
)
```

**Interfaces:**

```go
// UserRepository - Data access layer
type UserRepository interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, id uuid.UUID, user *User) (*User, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, page, limit int) ([]*User, int64, error)
    ListByRole(ctx context.Context, role UserRole, page, limit int) ([]*User, int64, error)
}

// UserUsecase - Business logic layer
type UserUsecase interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    Update(ctx context.Context, id uuid.UUID, user *User) (*User, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[User], error)
    ListByRole(ctx context.Context, role UserRole, page, limit int) (*utils.PaginatedResponse[User], error)
    Authenticate(ctx context.Context, email, password string) (*User, string, error)
}
```

### 2. Repository Layer (`repository/user.go`)
**Features:**
- ✅ Create user with UUID generation
- ✅ Get user by ID
- ✅ Get user by email (for authentication)
- ✅ Update user information
- ✅ Soft delete (mark is_active = false)
- ✅ List all active users with pagination
- ✅ List users filtered by role with pagination
- ✅ PostgreSQL with pgx/v5 driver
- ✅ Automatic timestamp management

**Key Methods:**
```go
Create()      // Insert new user with hashed password
GetByEmail()  // For login flow - only returns active users
Update()      // Updates allowed fields, rehashes password if provided
Delete()      // Soft delete by setting is_active = false
List()        // Paginated user listing (ORDER BY created_at DESC)
ListByRole()  // Filter by role with pagination
```

### 3. Usecase Layer (`usecase/user.go`)
**Business Logic:**

```go
// Create - Register new user
- Validates email, password, name, phone (required fields)
- Hashes password using bcrypt
- Sets default role as 'customer'
- Sets is_active = true

// GetByID - Fetch user details
- Validates UUID format

// Update - Modify user information
- Fetches existing user first
- Only updates specific fields
- Re-hashes password if provided

// Delete - Soft delete user
- Marks user as inactive

// List - Get paginated user list
- Uses utils.NormalizePagination() for page/limit validation
- Returns PaginatedResponse with total count

// ListByRole - Get users filtered by role
- Similar to List but filters by UserRole

// Authenticate - Login user
- Validates email and password
- Retrieves user from database by email
- Verifies password using bcrypt.CompareHashAndPassword()
- Generates JWT token valid for 24 hours
- Updates last_visit_at timestamp
- Returns user + token
```

**JWT Token Structure:**
```json
{
  "user_id": "uuid-string",
  "email": "user@example.com",
  "role": "customer",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### 4. Controller Layer (`controller/user.go`)
**API Endpoints:**

**Authentication (Public)**
```
POST /auth/register
- Body: {email, password, name, phone}
- Returns: User object (201 Created)

POST /auth/login  
- Body: {email, password}
- Returns: {user: User, token: string} (200 OK)
```

**Authentication (Public)**
```
POST /auth/register
- Body: {email, password, name, phone}
- Returns: User object (201 Created)
- Who: Everyone

POST /auth/login  
- Body: {email, password}
- Returns: {user: User, token: string} (200 OK)
- Who: Everyone
```

**User Profile (Self-Service)**
```
GET /users/me
- Returns: Current user's full profile
- Who: Any authenticated user
- Header: Authorization: Bearer <token>

PUT /users/:id
- Body: User object (partial update)
- Returns: Updated User
- Who: User can update self, Admin can update anyone
- Restriction: Non-admin users can only update their own profile
```

**User Management - Admin/Manager Only**
```
GET /users
- Query: ?page=1&limit=10
- Returns: PaginatedResponse with users list
- Who: Admin + Manager only
- Header: Authorization: Bearer <admin-token>

GET /users/:id
- Returns: User object
- Who: Any authenticated user can view any user

POST /users
- Body: User object
- Returns: Created User (201)
- Who: Admin only

DELETE /users/:id
- Returns: 204 No Content
- Who: Admin only
```

**User Filtering - Admin/Manager Only**
```
GET /users/role/:role
- Query: ?page=1&limit=10
- Params: role = "admin" | "manager" | "stylist" | "customer"
- Returns: PaginatedResponse filtered by role
- Who: Admin + Manager only
```

**Request/Response Examples:**

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@salon.com",
    "password": "SecurePassword123",
    "name": "John Doe",
    "phone": "+1234567890"
  }'

# Response
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "customer@salon.com",
    "name": "John Doe",
    "phone": "+1234567890",
    "role": "customer",
    "is_active": true,
    "loyalty_points": 0,
    "created_at": "2026-04-27T10:30:00Z",
    "updated_at": "2026-04-27T10:30:00Z"
  }
}

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@salon.com",
    "password": "SecurePassword123"
  }'

# Response
{
  "data": {
    "user": {...user object...},
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ..."
  }
}

# Get User (requires Authorization header)
curl -X GET http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ..."
```

### 5. Middleware (`api/middleware/auth.go`)
**Authentication Features:**
- ✅ JWT token validation (HMAC-SHA256)
- ✅ Bearer token parsing
- ✅ Automatic claims extraction
- ✅ Required authentication middleware
- ✅ Optional authentication middleware
- ✅ Role-based access control (RBAC)
- ✅ Context injection for route handlers

**Middleware Usage:**
```go
// Required auth
protected := router.Group("/api")
protected.Use(authMiddleware.Handler())
{
    protected.GET("/users", controller.ListUsers)
}

// Admin only
adminOnly := router.Group("/api")
adminOnly.Use(authMiddleware.Handler())
adminOnly.Use(authMiddleware.RequireRole("admin"))
{
    adminOnly.POST("/users", controller.CreateUser)
}

// Optional auth
public := router.Group("/api")
public.Use(authMiddleware.OptionalHandler())
{
    public.GET("/public-data", controller.GetPublicData)
}
```

### 6. Routes (`api/route/app.go`)
**Complete Route Registration with Role-Based Access Control:**

```go
// Public authentication routes (no auth required)
POST   /api/auth/login
POST   /api/auth/register

// Protected user management routes (requires authentication)
GET    /api/users/me                // Get current user's info (any authenticated user)
GET    /api/users                   // List all users (admin/manager only)
GET    /api/users/:id               // Get specific user (any authenticated user)
POST   /api/users                   // Create user (admin only)
PUT    /api/users/:id               // Update user (user can update self, admin can update any)
DELETE /api/users/:id               // Delete user (admin only)
GET    /api/users/role/:role        // List by role (admin/manager only)
```

**Role-Based Access Control:**

```go
// Public endpoints (NO authentication required)
POST /auth/login       → Everyone
POST /auth/register    → Everyone

// Protected endpoints (authentication required, but no role restriction)
GET /users/me          → Any authenticated user
GET /users/:id         → Any authenticated user (can view other users)

// Admin/Manager endpoints
GET  /users            → Admin + Manager only
GET  /users/role/:role → Admin + Manager only

// Admin-only endpoints
POST   /users          → Admin only (create user)
DELETE /users/:id      → Admin only (delete user)

// Flexible update (special logic)
PUT /users/:id         → Users can update themselves, Admin can update anyone
                       → Returns 400 if non-admin tries to update someone else's profile
```

## Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    birthday TIMESTAMP,
    address TEXT,
    role VARCHAR(50) NOT NULL DEFAULT 'customer',
    loyalty_points INTEGER DEFAULT 0,
    preferred_branch_id UUID,
    last_visit_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
```

## Security Features

✅ **Password Security:**
- Passwords hashed using bcrypt with cost factor 10
- Password never exposed in JSON responses
- Password comparison timing attack resistant

✅ **JWT Security:**
- HMAC-SHA256 signing
- 24-hour expiration
- Claims include user_id, email, and role
- Secret stored in environment variable

✅ **Authorization:**
- Role-based access control
- Middleware enforces authentication before protected routes
- Optional authentication for public endpoints
- Role validation middleware for admin endpoints

✅ **Data Protection:**
- Soft deletes (users not permanently deleted)
- Timestamps for audit trail
- Context-based user identification

## User Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      USER MANAGEMENT FLOW                        │
└─────────────────────────────────────────────────────────────────┘

1. REGISTRATION FLOW
   ┌──────────────┐
   │   Register   │
   └──────┬───────┘
          │ Email, Password, Name, Phone
          ▼
   ┌──────────────────────┐
   │ Validate Input       │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Hash Password        │
   │ (bcrypt)             │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Create User (DB)     │
   │ role=customer        │
   │ is_active=true       │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Return User Object   │
   │ (201 Created)        │
   └──────────────────────┘

2. LOGIN FLOW
   ┌──────────────┐
   │   Login      │
   └──────┬───────┘
          │ Email, Password
          ▼
   ┌──────────────────────┐
   │ Validate Input       │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Get User by Email    │
   │ (is_active=true)     │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Compare Password     │
   │ (bcrypt)             │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Generate JWT Token   │
   │ (24h expiry)         │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Update last_visit_at │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Return User + Token  │
   │ (200 OK)             │
   └──────────────────────┘

3. PROTECTED REQUEST FLOW
   ┌──────────────────────┐
   │ API Request with JWT │
   │ in Authorization     │
   │ header               │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Auth Middleware      │
   │ Validate Token       │
   └──────┬───────────────┘
          │
          ├─ Invalid? ──→ Return 401 Unauthorized
          │
          ▼
   ┌──────────────────────┐
   │ Extract Claims to    │
   │ Context              │
   │ (user_id, role)      │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Route Handler        │
   │ Process Request      │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Return Response      │
   │ (200/201/204/etc)    │
   └──────────────────────┘

4. ROLE-BASED ACCESS CONTROL
   ┌──────────────────────┐
   │ Request to Admin     │
   │ Resource             │
   └──────┬───────────────┘
          │
          ▼
   ┌──────────────────────┐
   │ Check User Role      │
   │ from Context         │
   └──────┬───────────────┘
          │
          ├─ Not Admin? ──→ Return 403 Forbidden
          │
          ▼
   ┌──────────────────────┐
   │ Allow Request        │
   └──────────────────────┘
```

## Usage Examples

### 1. Register New User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "stylist@salon.com",
    "password": "SecurePass123",
    "name": "Jane Smith",
    "phone": "+1987654321"
  }'
# Response: User object with role=customer (201 Created)
```

### 2. Login and Get Token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "stylist@salon.com",
    "password": "SecurePass123"
  }'

# Response includes JWT token
# {
#   "data": {
#     "user": {...},
#     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
#   }
# }
```

### 3. User Accesses Own Profile (Self-Service)
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get current user info
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer $TOKEN"

# Response: Current user's complete profile
```

### 4. User Updates Own Profile
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X PUT http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith Updated",
    "phone": "+1111111111"
  }'
# Response: Updated User object (200 OK)
```

### 5. User Tries to Update Another User (FAILS - Authorization)
```bash
CUSTOMER_TOKEN="..."  # Regular customer token

# Try to update admin user profile - SHOULD FAIL
curl -X PUT http://localhost:8080/api/users/admin-uuid \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Hacked"}'

# Response: 400 Bad Request
# {
#   "error": "you can only update your own profile"
# }
```

### 6. Admin Lists All Users
```bash
ADMIN_TOKEN="..."  # Admin user token

curl -X GET "http://localhost:8080/api/users?page=1&limit=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Response: PaginatedResponse with all active users
```

### 7. Manager Views Users by Role
```bash
MANAGER_TOKEN="..."  # Manager user token

curl -X GET "http://localhost:8080/api/users/role/stylist?page=1&limit=5" \
  -H "Authorization: Bearer $MANAGER_TOKEN"

# Response: PaginatedResponse with only stylists
```

### 8. Admin Creates New User
```bash
ADMIN_TOKEN="..."

curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newmanager@salon.com",
    "password": "SecurePass456",
    "name": "New Manager",
    "phone": "+1555555555",
    "role": "manager"
  }'

# Response: Created User object (201)
```

### 9. Admin Deletes User
```bash
ADMIN_TOKEN="..."

curl -X DELETE http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Response: 204 No Content (user is soft-deleted)
```

### 10. Regular User Tries to List All Users (FAILS - Authorization)
```bash
CUSTOMER_TOKEN="..."  # Customer token (no admin/manager role)

curl -X GET "http://localhost:8080/api/users?page=1&limit=10" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN"

# Response: 403 Forbidden
# {
#   "error": "Insufficient permissions"
# }
```

### 11. Customer Views Specific User Info (Allowed)
```bash
CUSTOMER_TOKEN="..."

curl -X GET http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer $CUSTOMER_TOKEN"

# Response: User object (200 OK) - any user can view any other user's profile
```

### 12. Invalid Token (FAILS - Authentication)
```bash
# Missing or invalid token
curl -X GET http://localhost:8080/api/users/me

# Response: 401 Unauthorized
# {
#   "error": "Authorization header required"
# }

# Expired/invalid token
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer invalid.token.here"

# Response: 401 Unauthorized
# {
#   "error": "Invalid token: ..."
# }

## Configuration

**Environment Variables:**
```bash
# JWT secret for token signing
JWT_SECRET=your-super-secret-key-min-32-chars

# Database connection
DATABASE_URL=postgres://user:password@localhost/salon_db
```

## Testing

### Unit Tests (Usecase Layer)
```go
// Test password hashing
// Test email validation
// Test token generation
// Test pagination
```

### Integration Tests (Controller + Repository)
```go
// Test registration flow
// Test login flow
// Test JWT validation
// Test CRUD operations
```

### Load Tests
```bash
# Test concurrent user creation
# Test login performance
# Test pagination under load
```

## Future Enhancements

- [ ] Email verification during registration
- [ ] Refresh token support
- [ ] Two-factor authentication
- [ ] OAuth2 integration
- [ ] Password reset flow
- [ ] User profile management
- [ ] Avatar upload
- [ ] Social login (Google, Facebook)
- [ ] User activity logging
- [ ] Session management

## Completion Status

✅ **Complete User CRUD Flow:**
- ✅ User registration
- ✅ User login with JWT
- ✅ User profile view
- ✅ User profile update
- ✅ User deletion (soft delete)
- ✅ List all users (paginated)
- ✅ Filter users by role (paginated)
- ✅ JWT authentication middleware
- ✅ Role-based access control
- ✅ Password security (bcrypt)
- ✅ Database schema ready

**All compilation errors resolved** - Ready for testing and integration!
