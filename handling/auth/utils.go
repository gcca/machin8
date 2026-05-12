package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const keyPrefix = "m8k_"
const sessionTTL = 7 * 24 * time.Hour

func hashKey(raw, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}

func Authenticate(db *sql.DB, username, password string) (*int64, error) {
	var id int64
	err := db.QueryRow(
		`SELECT id FROM auth.user
		 WHERE username = $1 AND password = crypt($2, password)`,
		username, password,
	).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	if _, err = db.Exec(
		`UPDATE auth.user SET last_logged_in = NOW() WHERE id = $1`, id,
	); err != nil {
		return nil, fmt.Errorf("update last_logged_in: %w", err)
	}
	return &id, nil
}

func CreateAPIKey(db *sql.DB, userID int64, name, secret string, expiry *time.Time) (rawKey string, err error) {
	rawKey = keyPrefix + uuid.NewString()
	hashed := hashKey(rawKey, secret)

	var q string
	var args []any
	if expiry != nil {
		q = `INSERT INTO auth.apikey (hashed_key, user_id, name, expiry) VALUES ($1, $2, $3, $4)`
		args = []any{hashed, userID, name, *expiry}
	} else {
		q = `INSERT INTO auth.apikey (hashed_key, user_id, name) VALUES ($1, $2, $3)`
		args = []any{hashed, userID, name}
	}

	if _, err = db.Exec(q, args...); err != nil {
		return "", fmt.Errorf("insert apikey: %w", err)
	}
	return rawKey, nil
}

func LogIn(c *gin.Context, db *sql.DB, userID int64) error {
	key := "m8s_" + uuid.NewString()
	expiry := time.Now().Add(sessionTTL)

	if _, err := db.Exec(
		`INSERT INTO auth.session (user_id, key, expires) VALUES ($1, $2, $3)`,
		userID, key, expiry,
	); err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	c.SetCookie("session", key, int(sessionTTL.Seconds()), "/", "", false, true)

	return nil
}
