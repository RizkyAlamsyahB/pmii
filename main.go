package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Router
	r := gin.Default()

	// 2. Setup CORS Middleware (PENTING!)
	// Agar Frontend (staging.pmii.id) tidak diblokir saat request ke API ini.
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Di production nanti ganti domain spesifik
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 3. Route Health Check (Root)
	r.GET("/", func(c *gin.Context) {
		// Mengambil Environment Variable dari Docker (Contoh pembuktian config server masuk)
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "Localhost (Dev Mode)"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":          "success",
			"message":         "🚀 API Staging PMII Live!",
			"service":         "Backend Go Gin",
			"db_connected_to": dbHost,
		})
	})

	// 4. Route Ping (Buat tes koneksi ringan)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 5. Jalankan Server di Port 8080 (Wajib sesuai Dockerfile)
	r.Run(":8080")
}
