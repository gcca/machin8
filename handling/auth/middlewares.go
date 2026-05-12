package auth

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"machin8/core/conf"
)

func TokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing api key"})
			return
		}
		rawKey := strings.TrimPrefix(header, "Bearer ")

		db, err := sql.Open("postgres", conf.Settings.DBUrl)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
			return
		}
		defer db.Close()

		hashed := hashKey(rawKey, conf.Settings.Secret)
		var userID int64
		err = db.QueryRow(
			`SELECT user_id FROM auth.apikey
			 WHERE hashed_key = $1 AND active = TRUE AND revoked = FALSE
			   AND (expiry IS NULL OR expiry > NOW())`,
			hashed,
		).Scan(&userID)
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authentication error"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func LogInRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawSession, err := c.Cookie("session")
		if err != nil || rawSession == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
			return
		}

		db, err := sql.Open("postgres", conf.Settings.DBUrl)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
			return
		}
		defer db.Close()

		var userID int64
		err = db.QueryRow(
			`SELECT user_id FROM auth.session
			 WHERE key = $1 AND expires > NOW()`,
			rawSession,
		).Scan(&userID)
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authentication error"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
