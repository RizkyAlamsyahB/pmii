# PMII API - Admin Activity Logs Module

> **Version:** 1.0.0  
> **Contact:** PMII Development Team

## Deskripsi

API Contract untuk modul **Admin Activity Logs**. Endpoint ini digunakan untuk audit aktivitas admin dan author di sistem. Semua endpoint dalam modul ini **hanya dapat diakses oleh Admin**. Memerlukan autentikasi menggunakan Bearer Token JWT.

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
| `GET` | `/admin/activity-logs` | Mengambil daftar activity log dengan pagination dan filter |

---

## 📖 Detail Endpoint

### Get All Activity Logs

Mengambil daftar semua activity log dengan pagination dan berbagai opsi filter.

**Endpoint:**
```
GET /v1/admin/activity-logs
```

**Akses:** Hanya Admin

#### Query Parameters

| Parameter | Tipe | Required | Default | Deskripsi |
|-----------|------|----------|---------|-----------|
| `page` | integer | ❌ | 1 | Nomor halaman (minimum: 1) |
| `limit` | integer | ❌ | 30 | Jumlah item per halaman (minimum: 1, maximum: 100) |
| `user_id` | integer | ❌ | - | Filter berdasarkan ID user |
| `module` | string | ❌ | - | Filter berdasarkan module |
| `action_type` | string | ❌ | - | Filter berdasarkan tipe aksi |
| `start_date` | string | ❌ | - | Filter dari tanggal (format: YYYY-MM-DD) |
| `end_date` | string | ❌ | - | Filter sampai tanggal (format: YYYY-MM-DD) |
| `search` | string | ❌ | - | Pencarian di field description |

#### Nilai yang Valid untuk Filter

**Module:**
| Nilai | Deskripsi |
|-------|-----------|
| `user` | Modul user management |
| `post` | Modul artikel/berita |
| `category` | Modul kategori |
| `tags` | Modul tag |
| `testimoni` | Modul testimonial |
| `members` | Modul member |
| `teams` | Modul tim |
| `dokumen` | Modul dokumen |
| `settings` | Modul pengaturan site |
| `auth` | Modul autentikasi |

**Action Type:**
| Nilai | Deskripsi |
|-------|-----------|
| `create` | Membuat data baru |
| `read` | Membaca data |
| `update` | Mengupdate data |
| `delete` | Menghapus data |
| `login` | Login ke sistem |
| `logout` | Logout dari sistem |

#### Contoh Request

```bash
# Default (page 1, limit 30)
GET /v1/admin/activity-logs

# Dengan pagination
GET /v1/admin/activity-logs?page=2&limit=10

# Filter by module
GET /v1/admin/activity-logs?module=post

# Filter by action type
GET /v1/admin/activity-logs?action_type=delete

# Filter by date range
GET /v1/admin/activity-logs?start_date=2025-12-01&end_date=2025-12-31

# Filter by user
GET /v1/admin/activity-logs?user_id=1

# Search di deskripsi
GET /v1/admin/activity-logs?search=password

# Kombinasi filter
GET /v1/admin/activity-logs?module=user&action_type=update&search=password
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Activity logs berhasil dimuat",
        "pagination": {
            "page": 1,
            "limit": 30,
            "total": 150,
            "lastPage": 5
        }
    },
    "data": [
        {
            "id": 1,
            "user_id": 1,
            "user": {
                "id": 1,
                "full_name": "Admin PMII",
                "email": "admin@pmii.or.id"
            },
            "action_type": "create",
            "module": "post",
            "description": "Created new post: Sejarah PMII",
            "target_id": 123,
            "old_value": null,
            "new_value": {
                "title": "Sejarah PMII"
            },
            "ip_address": "192.168.1.1",
            "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
            "created_at": "2025-12-30T09:00:00Z"
        },
        {
            "id": 2,
            "user_id": 1,
            "user": {
                "id": 1,
                "full_name": "Admin PMII",
                "email": "admin@pmii.or.id"
            },
            "action_type": "update",
            "module": "user",
            "description": "Changed password for user: John Doe",
            "target_id": 5,
            "old_value": null,
            "new_value": null,
            "ip_address": "192.168.1.1",
            "user_agent": "Mozilla/5.0",
            "created_at": "2025-12-30T08:30:00Z"
        }
    ]
}
```

#### Response Error

| Code | Deskripsi |
|------|-----------|
| 401 | Unauthorized - Token tidak valid atau tidak ada |
| 403 | Forbidden - Tidak memiliki akses ke resource ini |
| 500 | Internal Server Error |

---

## 📦 Data Models

### Activity Log Object

Object activity log yang dikembalikan dari API:

| Field | Tipe | Nullable | Deskripsi |
|-------|------|----------|-----------|
| `id` | integer | ❌ | ID activity log |
| `user_id` | integer | ❌ | ID user yang melakukan aksi |
| `user` | object | ❌ | Object informasi user |
| `action_type` | string | ❌ | Tipe aksi (create/read/update/delete/login/logout) |
| `module` | string | ❌ | Module tempat aksi dilakukan |
| `description` | string | ✅ | Deskripsi aktivitas |
| `target_id` | integer | ✅ | ID entitas yang menjadi target aksi |
| `old_value` | object | ✅ | Nilai sebelum perubahan (untuk update) |
| `new_value` | object | ✅ | Nilai setelah perubahan (untuk create/update) |
| `ip_address` | string | ✅ | IP address user |
| `user_agent` | string | ✅ | User agent browser |
| `created_at` | string | ❌ | Waktu aktivitas (format: ISO 8601) |

### Activity Log User Info Object

Object informasi user yang ada di dalam activity log:

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `id` | integer | ID user |
| `full_name` | string | Nama lengkap user |
| `email` | string | Email user |

### Pagination Object

Object pagination yang dikembalikan untuk response:

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

### 500 Internal Server Error

```json
{
    "meta": {
        "code": 500,
        "status": "error",
        "message": "Gagal mengambil data activity logs"
    },
    "data": null
}
```

---

## 📝 Catatan

1. Semua endpoint memerlukan autentikasi Bearer Token JWT.
2. Hanya user dengan role **Admin** yang dapat mengakses endpoint ini.
3. Activity log diurutkan berdasarkan waktu terbaru (descending).
4. Field `old_value` dan `new_value` berisi JSON object untuk menyimpan data sebelum dan sesudah perubahan.
5. Gunakan filter `module` dan `action_type` untuk mempersempit hasil pencarian.
6. Format tanggal untuk filter `start_date` dan `end_date` adalah `YYYY-MM-DD`.
