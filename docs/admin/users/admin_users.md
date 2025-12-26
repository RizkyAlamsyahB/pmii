# PMII API - Admin Users Module

> **Version:** 1.0.0  
> **Contact:** PMII Development Team

## Deskripsi

API Contract untuk modul **Admin Users**. Semua endpoint dalam modul ini **hanya dapat diakses oleh Admin**. Memerlukan autentikasi menggunakan Bearer Token JWT.

---

## 🔐 Autentikasi

Semua endpoint memerlukan autentikasi menggunakan **Bearer Token JWT**.

**Format Header:**
```
Authorization: Bearer <token>
```

Token didapatkan dari endpoint login.

---

## 📍 Base URL

```
/v1
```

---

## 📑 Daftar Endpoint

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/admin/users` | Mengambil daftar semua user dengan pagination |
| `POST` | `/admin/users` | Membuat user baru |
| `GET` | `/admin/users/{id}` | Mengambil detail user berdasarkan ID |
| `PUT` | `/admin/users/{id}` | Mengupdate data user berdasarkan ID |
| `DELETE` | `/admin/users/{id}` | Menghapus user berdasarkan ID (soft delete) |

---

## 📖 Detail Endpoint

### 1. Get All Users

Mengambil daftar semua user dengan pagination.

**Endpoint:**
```
GET /v1/admin/users
```

**Akses:** Hanya Admin

#### Query Parameters

| Parameter | Tipe | Required | Default | Deskripsi |
|-----------|------|----------|---------|-----------|
| `page` | integer | ❌ | 1 | Nomor halaman (minimum: 1) |
| `limit` | integer | ❌ | 20 | Jumlah item per halaman (minimum: 1, maximum: 100) |

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Data user berhasil diambil",
        "pagination": {
            "page": 1,
            "limit": 20,
            "total": 100,
            "lastPage": 5
        }
    },
    "data": {
        "users": [
            {
                "id": 1,
                "fullName": "John Doe",
                "email": "john@example.com",
                "role": "admin",
                "status": "active",
                "photoUri": "https://example.com/photos/john.jpg"
            },
            {
                "id": 2,
                "fullName": "Jane Smith",
                "email": "jane@example.com",
                "role": "author",
                "status": "active",
                "photoUri": ""
            }
        ],
        "total": 2
    }
}
```

#### Response Error

| Code | Deskripsi |
|------|-----------|
| 401 | Unauthorized - Token tidak valid atau tidak ada |
| 403 | Forbidden - Tidak memiliki akses ke resource ini |
| 500 | Internal Server Error |

---

### 2. Create User

Membuat user baru di sistem. Mendukung upload foto profile menggunakan `multipart/form-data`.

**Endpoint:**
```
POST /v1/admin/users
```

**Akses:** Hanya Admin

#### Request Body

**Content-Type:** `multipart/form-data` atau `application/json`

| Field | Tipe | Required | Deskripsi |
|-------|------|----------|-----------|
| `full_name` | string | ✅ | Nama lengkap user (2-100 karakter) |
| `email` | string | ✅ | Email user (harus unik, format email valid) |
| `password` | string | ✅ | Password user (minimal 8 karakter) |
| `photo` | file | ❌ | Foto profil user (maksimal 5MB) - hanya untuk `multipart/form-data` |

#### Contoh Request (JSON)

```json
{
    "full_name": "New User",
    "email": "newuser@example.com",
    "password": "password123"
}
```

#### Response Success (201)

```json
{
    "meta": {
        "code": 201,
        "status": "success",
        "message": "User berhasil dibuat"
    },
    "data": {
        "id": 3,
        "fullName": "New User",
        "email": "newuser@example.com",
        "role": "author",
        "status": "active",
        "photoUri": "https://example.com/photos/newuser.jpg"
    }
}
```

#### Response Error

| Code | Deskripsi | Contoh Message |
|------|-----------|----------------|
| 400 | Validation Error | "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag" |
| 400 | Email Already Exists | "Email sudah digunakan" |
| 400 | Invalid Password | "Password tidak valid" |
| 401 | Unauthorized | "Akses ditolak - Token tidak valid" |
| 403 | Forbidden | "Akses ditolak - Hanya admin yang dapat mengakses" |
| 500 | Internal Server Error | "Terjadi kesalahan server" |

---

### 3. Get User by ID

Mengambil detail user berdasarkan ID.

**Endpoint:**
```
GET /v1/admin/users/{id}
```

**Akses:** Hanya Admin

#### Path Parameters

| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id` | integer | ✅ | ID unik user (minimum: 1) |

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Profil berhasil diambil"
    },
    "data": {
        "id": 1,
        "fullName": "John Doe",
        "email": "john@example.com",
        "role": "admin",
        "status": "active",
        "photoUri": "https://example.com/photos/john.jpg"
    }
}
```

#### Response Error

| Code | Deskripsi | Contoh Message |
|------|-----------|----------------|
| 400 | Invalid ID | "ID tidak valid" |
| 401 | Unauthorized | "Akses ditolak - Token tidak valid" |
| 403 | Forbidden | "Akses ditolak - Hanya admin yang dapat mengakses" |
| 404 | Not Found | "User tidak ditemukan" |
| 500 | Internal Server Error | "Terjadi kesalahan server" |

---

### 4. Update User by ID

Mengupdate data user berdasarkan ID. Semua field bersifat opsional - hanya field yang disertakan yang akan di-update. Mendukung upload foto profile menggunakan `multipart/form-data`.

**Endpoint:**
```
PUT /v1/admin/users/{id}
```

**Akses:** Hanya Admin

#### Path Parameters

| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id` | integer | ✅ | ID unik user (minimum: 1) |

#### Request Body

**Content-Type:** `multipart/form-data` atau `application/json`

| Field | Tipe | Required | Deskripsi |
|-------|------|----------|-----------|
| `full_name` | string | ❌ | Nama lengkap user (2-100 karakter) |
| `email` | string | ❌ | Email user (harus unik, format email valid) |
| `password` | string | ❌ | Password baru user (minimal 8 karakter) |
| `role` | integer | ❌ | Role user: `1` = Admin, `2` = Author |
| `is_active` | boolean | ❌ | Status aktif user (true/false) |
| `photo` | file | ❌ | Foto profil user (maksimal 5MB) - hanya untuk `multipart/form-data` |

#### Contoh Request (JSON)

```json
{
    "full_name": "John Doe Updated",
    "email": "john.updated@example.com",
    "role": 2,
    "is_active": true
}
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "User berhasil diupdate"
    },
    "data": {
        "id": 1,
        "fullName": "John Doe Updated",
        "email": "john.updated@example.com",
        "role": "admin",
        "status": "active",
        "photoUri": "https://example.com/photos/john-new.jpg"
    }
}
```

#### Response Error

| Code | Deskripsi | Contoh Message |
|------|-----------|----------------|
| 400 | Invalid ID | "ID tidak valid" |
| 400 | Email Already Used | "Email sudah digunakan user lain" |
| 400 | Invalid Password | "Password tidak valid" |
| 401 | Unauthorized | "Akses ditolak - Token tidak valid" |
| 403 | Forbidden | "Akses ditolak - Hanya admin yang dapat mengakses" |
| 404 | Not Found | "User tidak ditemukan" |
| 500 | Internal Server Error | "Terjadi kesalahan server" |

---

### 5. Delete User by ID

Menghapus user berdasarkan ID (soft delete).

**Endpoint:**
```
DELETE /v1/admin/users/{id}
```

**Akses:** Hanya Admin

#### Path Parameters

| Parameter | Tipe | Required | Deskripsi |
|-----------|------|----------|-----------|
| `id` | integer | ✅ | ID unik user (minimum: 1) |

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "User berhasil dihapus"
    },
    "data": null
}
```

#### Response Error

| Code | Deskripsi | Contoh Message |
|------|-----------|----------------|
| 400 | Invalid ID | "ID tidak valid" |
| 401 | Unauthorized | "Akses ditolak - Token tidak valid" |
| 403 | Forbidden | "Akses ditolak - Hanya admin yang dapat mengakses" |
| 404 | Not Found | "User tidak ditemukan" |
| 500 | Internal Server Error | "Terjadi kesalahan server" |

---

## 📦 Data Models

### User Object

Object user yang dikembalikan dari API:

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `id` | integer | ID user |
| `fullName` | string | Nama lengkap user |
| `email` | string | Email user |
| `role` | string | Role user (`admin` / `author`) |
| `status` | string | Status user (`active` / `inactive`) |
| `photoUri` | string | URL foto profil user |

### Pagination Object

Object pagination yang dikembalikan untuk response yang mendukung pagination:

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `page` | integer | Halaman saat ini |
| `limit` | integer | Jumlah item per halaman |
| `total` | integer | Total semua item |
| `lastPage` | integer | Halaman terakhir |

### Meta Object

Object meta yang selalu ada di setiap response:

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `code` | integer | HTTP status code |
| `status` | string | Status response (`success` / `error`) |
| `message` | string | Pesan response |
| `pagination` | object | Object pagination (opsional) |

---

## ⚠️ Error Responses

### 401 Unauthorized

```json
{
    "meta": {
        "code": 401,
        "status": "error",
        "message": "Akses ditolak - Token tidak valid"
    },
    "data": null
}
```

### 403 Forbidden

```json
{
    "meta": {
        "code": 403,
        "status": "error",
        "message": "Akses ditolak - Hanya admin yang dapat mengakses"
    },
    "data": null
}
```

### 404 Not Found

```json
{
    "meta": {
        "code": 404,
        "status": "error",
        "message": "User tidak ditemukan"
    },
    "data": null
}
```

### 500 Internal Server Error

```json
{
    "meta": {
        "code": 500,
        "status": "error",
        "message": "Terjadi kesalahan server"
    },
    "data": null
}
```

---

## 📝 Catatan

1. Semua endpoint memerlukan autentikasi Bearer Token JWT.
2. Hanya user dengan role **Admin** yang dapat mengakses endpoint-endpoint ini.
3. Upload foto profil menggunakan `multipart/form-data` dengan maksimal ukuran **5MB**.
4. Penghapusan user menggunakan metode **soft delete**.
5. Role yang tersedia: `Admin (1)` dan `Author (2)`.
