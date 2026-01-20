package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

/* =========================
   RateLimiter tests
   ========================= */

func setupEchoWithRateLimiter(rdb *redis.Client) *echo.Echo {
	e := echo.New()
	e.Use(RateLimiter(rdb))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	return e
}

func TestRateLimiter_AllowsRequestsUnderLimit(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	e := setupEchoWithRateLimiter(rdb)

	for i := 0; i < 60; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestRateLimiter_BlocksAfterLimit(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	e := setupEchoWithRateLimiter(rdb)

	for i := 0; i < 61; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if i == 60 {
			assert.Equal(t, http.StatusTooManyRequests, rec.Code)
		}
	}
}

func TestRateLimiter_TTLIsSet(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	e := setupEchoWithRateLimiter(rdb)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	ttl := mr.TTL("rate_limit:127.0.0.1")
	assert.True(t, ttl > 0 && ttl <= time.Minute)
}

/* =========================
   JWTMiddleware tests
   ========================= */

func setupEchoWithJWT(secret string) *echo.Echo {
	e := echo.New()
	e.Use(JWTMiddleware(secret))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id":    c.Get("user_id"),
			"user_email": c.Get("user_email"),
		})
	})
	return e
}

func generateToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    1,
		"email": "test@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte(secret))
	return s
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	secret := "secret"
	e := setupEchoWithJWT(secret)

	token := generateToken(secret)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "test@example.com")
}

func TestJWTMiddleware_MissingHeader(t *testing.T) {
	e := setupEchoWithJWT("secret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidFormat(t *testing.T) {
	e := setupEchoWithJWT("secret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "BadFormat")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidSignature(t *testing.T) {
	e := setupEchoWithJWT("secret")

	token := generateToken("wrong-secret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	secret := "secret"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    1,
		"email": "test@example.com",
		"exp":   time.Now().Add(-time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte(secret))

	e := setupEchoWithJWT(secret)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+s)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_WrongSigningMethod(t *testing.T) {
	secret := "secret"

	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"id": 1,
	})
	s, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	e := setupEchoWithJWT(secret)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+s)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRateLimiter_RedisError(t *testing.T) {
	// Некорректный адрес Redis → Incr вернёт ошибку
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:0",
	})

	e := echo.New()
	e.Use(RateLimiter(rdb))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Redis error")
}
