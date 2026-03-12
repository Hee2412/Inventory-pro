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
7. **Session auto-close** (khi hết deadline)
8. **SuperAdmin review** → Approve/Decline reports

### 🔍 Advanced Features
- **Pagination**: Endpoints hỗ trợ phân trang
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