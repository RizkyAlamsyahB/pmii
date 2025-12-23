# PMII API - Public Home Module

> **Version:** 1.0.0  
> **Contact:** PMII Development Team

## Deskripsi

API Contract untuk modul **Public Home**. Modul ini menyediakan data konten untuk halaman utama (landing page) aplikasi. Semua endpoint dalam modul ini dapat diakses secara publik tanpa memerlukan autentikasi.

---

## 📍 Base URL

```
/v1
```

---

## 📑 Daftar Endpoint

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/home/hero` | Mengambil data hero section (headline) |
| `GET` | `/home/latest-news` | Mengambil daftar berita terbaru |
| `GET` | `/home/about-us` | Mengambil data ringkasan 'Tentang Kami' |
| `GET` | `/home/why` | Mengambil data 'Mengapa PMII?' |
| `GET` | `/home/testimonial` | Mengambil daftar testimoni |
| `GET` | `/home/faq` | Mengambil daftar FAQ |
| `GET` | `/home/cta` | Mengambil data Call to Action |

---

## 📖 Detail Endpoint

### 1. Get Hero Section

Mengambil daftar post populer yang dijadikan headline di halaman utama.

**Endpoint:**
```
GET /v1/home/hero
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data hero section"
    },
    "data": [
        {
            "id": 1,
            "title": "Judul Headline",
            "slug": "judul-headline",
            "excerpt": "Ringkasan konten...",
            "imageUrl": "https://res.cloudinary.com/...",
            "publishedAt": "2023-12-23",
            "category": {
                "id": 1,
                "name": "Opini"
            },
            "authorId": 1,
            "tags": ["Kaderisasi", "Nasional"]
        }
    ]
}
```

---

### 2. Get Latest News Section

Mengambil 5 berita terbaru yang dipublikasikan.

**Endpoint:**
```
GET /v1/home/latest-news
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data latest news section"
    },
    "data": [
        {
            "id": 1,
            "title": "Berita Terbaru",
            "featured_image": "https://res.cloudinary.com/...",
            "created_at": "2023-12-23T10:00:00Z",
            "updated_at": "2023-12-23T10:00:00Z",
            "total_views": 150
        }
    ]
}
```

---

### 3. Get About Us Section

Mengambil informasi singkat mengenai profil organisasi untuk ditampilkan di landing page.

**Endpoint:**
```
GET /v1/home/about-us
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data about us section"
    },
    "data": {
        "title": "Sekilas tentang sejarah...",
        "subtitle": "Tentang PMII",
        "description": "Organisasi mahasiswa berbasis nilai...",
        "image_uri": "https://res.cloudinary.com/..."
    }
}
```

---

### 4. Get Why Section

Mengambil data poin-poin alasan atau nilai unggul organisasi.

**Endpoint:**
```
GET /v1/home/why
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data why section"
    },
    "data": {
        "title": "Nilai dan keunggulan...",
        "subtitle": "Mengapa PMII?",
        "description": null,
        "data": [
            {
                "title": "Perkaderan Berjenjang",
                "description": "Pembentukan karakter...",
                "icon_uri": "https://res.cloudinary.com/..."
            }
        ]
    }
}
```

---

### 5. Get Testimonial Section

Mengambil daftar testimoni (limit 7).

**Endpoint:**
```
GET /v1/home/testimonial
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data testimonial section"
    },
    "data": [
        {
            "testimoni": "PMII memberikan saya ruang tumbuh...",
            "name": "Ahmad",
            "status": "Alumni",
            "career": "CEO at Tech Corp",
            "image_uri": "https://res.cloudinary.com/..."
        }
    ]
}
```

---

### 6. Get FAQ Section

Mengambil daftar pertanyaan yang sering diajukan.

**Endpoint:**
```
GET /v1/home/faq
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data faq section"
    },
    "data": {
        "title": "Jawaban atas pertanyaan umum...",
        "subtitle": "FAQ",
        "description": null,
        "data": [
            {
                "question": "Apa itu PMII?",
                "answer": "PMII adalah organisasi..."
            }
        ]
    }
}
```

---

### 7. Get CTA Section

Mengambil konten ajakan bergabung.

**Endpoint:**
```
GET /v1/home/cta
```

#### Response Success (200)

```json
{
    "meta": {
        "code": 200,
        "status": "success",
        "message": "Berhasil mengambil data cta section"
    },
    "data": {
        "title": "Ambil langkah pertama...",
        "subtitle": "Siap Bergabung dengan PMII?"
    }
}
```

---

## 📦 Data Models

### Meta Object

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `code` | integer | HTTP status code |
| `status` | string | Status response (`success` / `error`) |
| `message` | string | Pesan response |

---

## ⚠️ Error Responses

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

1. Semua endpoint bersifat **Public**.
2. Tidak ada pagination pada endpoint-endpoint ini karena jumlah data telah dibatasi di sisi server untuk kebutuhan layout landing page.
3. Media (gambar/ikon) dikembalikan dalam bentuk URL lengkap dari Cloudinary.
