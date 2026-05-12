package ngin

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"machin8/core/conf"
)

func OwnerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(int64)

		notebookID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid notebook id"})
			return
		}

		db, err := sql.Open("postgres", conf.Settings.DBUrl)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
			return
		}
		defer db.Close()

		var count int
		err = db.QueryRow(`
			SELECT COUNT(*)
			FROM ngin.notebook n
			JOIN org.analist a ON a.id = n.analist_id
			WHERE n.id = $1 AND a.user_id = $2`,
			notebookID, userID,
		).Scan(&count)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
