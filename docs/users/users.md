# PMII API - Users Module

API Contract untuk modul Users.
Endpoint dalam modul ini dapat diakses oleh **semua authenticated users**.
Memerlukan autentikasi menggunakan Bearer Token JWT.

**Version:** 1.0.0
**Base URL:** `/v1`

---

## Authentication

Modul ini menggunakan **Bearer Authentication (JWT)**. Token diperoleh dari endpoint login.

**Format Header:**
`Authorization: Bearer <token>`

---

## Endpoints

### 1. Get My Profile

Mengambil profil user yang sedang login (berdasarkan token JWT).
Dapat diakses oleh **semua authenticated users** (baik admin maupun user biasa).

- **URL:** `/users/me`
- **Method:** `GET`
- **Security:** `BearerAuth`

#### Responses

| Code | Description |
| :--- | :--- |
| **200** | Berhasil mendapatkan profil |
| **401** | Unauthorized - Token tidak valid, expired, atau tidak disertakan |
| **404** | Not Found - User tidak ditemukan |
| **500** | Internal Server Error - Terjadi kesalahan pada server |

#### Example Response (200 OK)

```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Berhasil mendapatkan profile"
  },
  "data": {
    "id": 1,
    "fullName": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "status": "active",
    "photoUri": "https://example.com/photos/john.jpg"
  }
}
```

#### Example Response (401 Unauthorized)

**Token tidak disertakan:**
```json
{
  "meta": {
    "code": 401,
    "status": "error",
    "message": "Akses ditolak - Token tidak ada"
  },
  "data": null
}
```

**Token tidak valid:**
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

#### Example Response (404 Not Found)
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

---

## Data Models

### UserProfileResponse

| Field | Type | Description | Example |
| :--- | :--- | :--- | :--- |
| `id` | Integer | ID unik user | `1` |
| `fullName` | String | Nama lengkap user | `John Doe` |
| `email` | String | Email user | `john@example.com` |
| `role` | String | Role user in system (`admin`, `user`) | `user` |
| `status` | String | Status akun user (`active`, `inactive`) | `active` |
| `photoUri` | String | URL foto profil user (opsional) | `https://example.com/photos/john.jpg` |

### Meta

| Field | Type | Description | Example |
| :--- | :--- | :--- | :--- |
| `code` | Integer | HTTP status code | `200` |
| `status` | String | Status response (`success`, `error`) | `success` |
| `message` | String | Pesan response | `Data berhasil diambil` |
