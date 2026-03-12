# Inventory Management System

Hệ thống quản lý kho hàng và đặt hàng cho chuỗi cửa hàng.

## ✨ Features

### 👥 User Management
- Phân quyền: Super Admin, Admin, Store
- CRUD users
- Activate/Deactivate accounts
- Soft delete với hard delete (Super Admin only)

### 📦 Product Management
- CRUD products
- Search & filter products
- Activate/Deactivate products
- Soft delete với hard delete (Super Admin only)

### 🛒 Order Flow
1. **Admin tạo Order Session** (kỳ đặt hàng)
2. **Admin thêm products** vào session
3. **Store đăng nhập** → Xem session
4. **Store tạo/update order** (chọn số lượng)
5. **Store submit order**
6. **Admin review** → Approve hoặc Decline
7. **Admin close session** khi xong

### 📊 Audit Flow (Inventory Checking)
1. **SuperAdmin tạo Audit Session** (kỳ kiểm kê)
2. **SuperAdmin thêm products** cần kiểm
3. **System tự động tạo audit items** cho TẤT CẢ stores
4. **Store điền số lượng thực tế** (actual stock)
5. **System tự động tính variance** (chênh lệch)
6. **Store save** nhiều lần (trước deadline)
7. **Session auto-close** khi hết deadline
8. **SuperAdmin review** → Approve/Decline reports

### 🔍 Advanced Features
- **Pagination**: Tất cả endpoints hỗ trợ phân trang
- **Search & Filter**: Tìm kiếm và lọc dữ liệu
- **Request Logging**: Log tất cả requests
- **JWT Authentication**: Bảo mật với JWT tokens

---

## 🛠 Tech Stack

- **Backend**: Golang 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Architecture**: Clean Architecture (Repository - Service - Handler)

---

## 📁 Project Structure
```
Inventory-pro/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── domain/           # Models
│   ├── dto/
│   │   ├── request/      # Request DTOs
│   │   └── response/     # Response DTOs
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   ├── handler/          # HTTP handlers
│   └── middleware/       # Middleware (Auth, Logging)
├── pkg/
│   ├── database/         # DB connection
│   ├── pagination/       # Pagination helper
│   └── response/         # Response helper
├── .env
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 Setup & Installation

### 1. Clone repository
```bash
git clone https://github.com/Hee2412/inventory-pro.git
cd inventory-pro
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Setup environment
```bash
cp .env.example .env
```

Edit `.env`:
```env
DATABASE_URL=postgres://user:password@localhost:5432/inventory?sslmode=disable
JWT_SECRET=your-super-secret-key-change-this
PORT=8080
```

### 4. Setup Database
```bash
# Create PostgreSQL database
createdb inventory

# Run application (auto-migrate)
go run main.go
```

### 5. Create initial Super Admin
```sql
-- Connect to database
psql inventory

-- Insert super admin
INSERT INTO users (email, password, role, is_active, created_at, updated_at)
VALUES (
  'admin@inventory.com',
  '$2a$10$XYZ...', -- Hash of 'password123'
  'super_admin',
  true,
  NOW(),
  NOW()
);
```

Or use API:
```bash
POST /auth/register
{
  "email": "admin@inventory.com",
  "password": "password123",
  "role": "super_admin"
}
```

---

## 🔐 Authentication

### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "username": "yadmin@inventor.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "admin@inventory.com",
      "role": "super_admin"
    }
  }
}
```

### Use Token
```bash
GET /api/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## 📚 API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication Endpoints
```
POST   /auth/login          # Login
POST   /auth/register       # Register (Admin only)
GET    /api/me             # Get current user
```

### User Management (Admin)
```
GET    /api/admin/users                    # List users (with pagination & filters)
GET    /api/admin/users/:id                # Get user
POST   /api/admin/users/register           # Create user
PUT    /api/admin/users/:id                # Update user
PATCH  /api/admin/users/:id/activate       # Activate user
PATCH  /api/admin/users/:id/deactivate     # Deactivate user
DELETE /api/admin/users/:id                # Soft delete
DELETE /api/superadmin/users/:id/hard      # Hard delete (SuperAdmin)
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)
- `search`: Search by email or store_name
- `role`: Filter by role (store, admin, super_admin)
- `is_active`: Filter by active status (true/false)

### Product Management
```
GET    /api/products                       # List products (public)
GET    /api/products/:id                   # Get product
GET    /api/admin/products                 # List all (with inactive)
POST   /api/admin/products                 # Create product
PUT    /api/admin/products/:id             # Update product
PATCH  /api/admin/products/:id/activate    # Activate
PATCH  /api/admin/products/:id/deactivate  # Deactivate
DELETE /api/admin/products/:id             # Soft delete
DELETE /api/superadmin/products/:id/hard   # Hard delete
```

**Query Parameters:**
- `page`, `limit`: Pagination
- `search`: Search by name or code
- `is_active`: Filter by status
- `min_price`, `max_price`: Price range

### Order Flow
```
# Admin - Session Management
POST   /api/admin/sessions                           # Create session
GET    /api/admin/sessions                           # List sessions
GET    /api/admin/sessions/:sessionId                # Get session
POST   /api/admin/sessions/products                  # Add products to session
DELETE /api/admin/sessions/:sessionId/products/:productId  # Remove product
PATCH  /api/admin/sessions/:sessionId/close          # Close session

# Store - Order Management
GET    /api/store/sessions/:sessionId/order          # Get/Create order
PUT    /api/store/orders/:orderId/items              # Update order items
GET    /api/store/orders/:orderId                    # Get order detail
GET    /api/store/orders                             # List my orders

# Admin - Order Review
GET    /api/admin/sessions/:sessionId/orders         # Orders in session
GET    /api/admin/orders                             # All orders (with filters)
POST   /api/admin/orders/:orderId/approve            # Approve order
POST   /api/admin/orders/:orderId/decline            # Decline order
```

### Audit Flow
```
# SuperAdmin - Audit Session
POST   /api/admin/audit-sessions                     # Create audit session
GET    /api/admin/audit-sessions                     # List sessions
GET    /api/admin/audit-sessions/:sessionId          # Get session
POST   /api/admin/audit-sessions/products            # Add products
DELETE /api/admin/audit-sessions/:sessionId/products/:productId  # Remove
PATCH  /api/admin/audit-sessions/:sessionId/close    # Close session
PUT    /api/admin/audit-sessions/:sessionId          # Update session

# Store - Audit Report
GET    /api/store/audit-sessions/:sessionId/report   # View items
PUT    /api/store/audit-sessions/:sessionId/items    # Update items (batch)
GET    /api/store/audit-reports                      # My audit history

# SuperAdmin - Review
GET    /api/superadmin/audit-sessions/:sessionId/reports           # All reports
GET    /api/superadmin/audit-sessions/:sessionId/stores/:storeId   # Detail
GET    /api/superadmin/audit-sessions/:sessionId/summary           # Summary
POST   /api/superadmin/audit-sessions/:sessionId/stores/:storeId/approve
POST   /api/superadmin/audit-sessions/:sessionId/stores/:storeId/decline
```

---

## 📮 Postman Collection

Import collection: [Download Here](./docs/Inventory-API.postman_collection.json)

---

## 🗄 Database Schema

### Users
```sql
id, email, password, role, store_name, phone, address, is_active, 
created_at, updated_at, deleted_at
```

### Products
```sql
id, product_name, product_code, description, price, unit, is_active,
created_at, updated_at, deleted_at
```

### Order Sessions
```sql
id, title, order_cycle, deadline, delivery_date, status, created_by,
created_at, updated_at, deleted_at
```

### Store Orders
```sql
id, store_id, session_id, status, submitted_at, approved_at, approved_by,
created_at, updated_at
```

### Order Items
```sql
id, order_id, product_id, quantity, product_name, product_code
```

### Audit Sessions
```sql
id, title, audit_type, start_date, end_date, status, created_by,
created_at, updated_at
```

### Audit Reports (Store Audit Reports)
```sql
id, session_id, store_id, product_id, system_stock, actual_stock, variance,
approved_at, approved_by, updated_at
```

---

## 🔄 Workflows

### Order Workflow
```
1. Admin creates Order Session
2. Admin adds Products to Session
3. Stores login and view Session
4. Stores create/update Orders (select quantities)
5. Stores submit Orders
6. Admin reviews and Approves/Declines
7. Admin closes Session
```

### Audit Workflow
```
1. SuperAdmin creates Audit Session (with deadline)
2. SuperAdmin adds Products
3. System auto-creates audit items for ALL stores
4. Stores fill in actual stock (can edit multiple times before deadline)
5. Session auto-closes at deadline
6. SuperAdmin reviews reports and Approves/Declines
```

---

## 👨‍💻 Development

### Run locally
```bash
go run main.go
```

### Build
```bash
go build -o inventory-app
./inventory-app
```

### Test
```bash
go test ./...
```

---

## 📦 Deployment

### Docker
```bash
docker build -t inventory-app .
docker run -p 8080:8080 --env-file .env inventory-app
```

### Docker Compose
```bash
docker-compose up -d
```

---

## 👤 Author

**Huy Truong**
- GitHub: [@Hee2412](https://github.com/Hee2412)
- Email: huytruong2412hee@gmail.com

