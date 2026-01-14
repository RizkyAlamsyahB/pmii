package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRequireRole_ValidRolePasses menguji admin user akses admin route
func TestRequireRole_ValidRolePasses(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set user_role di context (dari AuthMiddleware)
	c.Set("user_role", "1") // Admin role

	// Handler yang akan dipanggil jika middleware pass
	handlerCalled := false
	testHandler := func(c *gin.Context) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "success"})
	}

	// Jalankan middleware
	middleware := RequireRole("1")
	middleware(c)

	// Jika middleware pass, lanjut ke handler
	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if c.IsAborted() {
		t.Error("Expected middleware to pass, but it aborted")
	}
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestRequireRole_InvalidRoleBlocked menguji user biasa akses admin route (harus 403)
func TestRequireRole_InvalidRoleBlocked(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set user_role di context (User role, bukan Admin)
	c.Set("user_role", "2") // User role

	// Handler yang seharusnya TIDAK dipanggil
	handlerCalled := false
	testHandler := func(_ *gin.Context) {
		handlerCalled = true
	}

	// Jalankan middleware
	middleware := RequireRole("1") // Require Admin
	middleware(c)

	// Middleware harus abort
	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if !c.IsAborted() {
		t.Error("Expected middleware to abort, but it passed")
	}
	if handlerCalled {
		t.Error("Expected handler NOT to be called")
	}
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

// TestRequireRole_NoAuthContext menguji request tanpa user_role di context (harus 401)
func TestRequireRole_NoAuthContext(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// TIDAK set user_role di context (simulasi belum auth)

	// Handler yang seharusnya TIDAK dipanggil
	handlerCalled := false
	testHandler := func(_ *gin.Context) {
		handlerCalled = true
	}

	// Jalankan middleware
	middleware := RequireRole("1")
	middleware(c)

	// Middleware harus abort
	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if !c.IsAborted() {
		t.Error("Expected middleware to abort, but it passed")
	}
	if handlerCalled {
		t.Error("Expected handler NOT to be called")
	}
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// TestRequireRole_WrongRoleType menguji context value bukan string (harus 500)
func TestRequireRole_WrongRoleType(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set user_role dengan tipe yang salah (int bukan string)
	c.Set("user_role", 1) // Tipe salah: int bukan string

	// Handler yang seharusnya TIDAK dipanggil
	handlerCalled := false
	testHandler := func(_ *gin.Context) {
		handlerCalled = true
	}

	// Jalankan middleware
	middleware := RequireRole("1")
	middleware(c)

	// Middleware harus abort
	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if !c.IsAborted() {
		t.Error("Expected middleware to abort, but it passed")
	}
	if handlerCalled {
		t.Error("Expected handler NOT to be called")
	}
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// TestRequireAnyRole_AdminAccess menguji RequireAnyRole dengan admin role
func TestRequireAnyRole_AdminAccess(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set admin role
	c.Set("user_role", "1")

	// Handler
	handlerCalled := false
	testHandler := func(c *gin.Context) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "success"})
	}

	// Jalankan middleware - izinkan admin dan user
	middleware := RequireAnyRole("1", "2")
	middleware(c)

	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if c.IsAborted() {
		t.Error("Expected middleware to pass for admin role")
	}
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestRequireAnyRole_UserAccess menguji RequireAnyRole dengan user role
func TestRequireAnyRole_UserAccess(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set user role
	c.Set("user_role", "2")

	// Handler
	handlerCalled := false
	testHandler := func(c *gin.Context) {
		handlerCalled = true
		c.JSON(200, gin.H{"message": "success"})
	}

	// Jalankan middleware - izinkan admin dan user
	middleware := RequireAnyRole("1", "2")
	middleware(c)

	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if c.IsAborted() {
		t.Error("Expected middleware to pass for user role")
	}
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestRequireAnyRole_UnauthorizedRole menguji RequireAnyRole dengan role yang tidak diizinkan
func TestRequireAnyRole_UnauthorizedRole(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set role yang tidak ada dalam allowed list
	c.Set("user_role", "3") // Role tidak dikenal

	// Handler
	handlerCalled := false
	testHandler := func(_ *gin.Context) {
		handlerCalled = true
	}

	// Jalankan middleware - hanya izinkan admin
	middleware := RequireAnyRole("1")
	middleware(c)

	if !c.IsAborted() {
		testHandler(c)
	}

	// Validasi
	if !c.IsAborted() {
		t.Error("Expected middleware to abort for unauthorized role")
	}
	if handlerCalled {
		t.Error("Expected handler NOT to be called")
	}
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}
