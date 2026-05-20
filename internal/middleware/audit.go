package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log mutations (POST, PUT, DELETE, PATCH)
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
			c.Next()
			return
		}

		// Read and buffer the request body
		var reqBody []byte
		if c.Request.Body != nil {
			var err error
			reqBody, err = io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}
		}

		// Process the request
		c.Next()

		// Skip logging if the request was unsuccessful
		status := c.Writer.Status()
		if status < 200 || status >= 300 {
			return
		}

		// Get authenticated user ID
		userIDVal, exists := c.Get("userID")
		if !exists {
			return
		}

		userIDStr, ok := userIDVal.(string)
		if !ok || userIDStr == "" {
			return
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			return
		}

		// Determine action
		action := string(models.AuditUpdate)
		if method == "POST" {
			action = string(models.AuditCreate)
		} else if method == "DELETE" {
			action = string(models.AuditDelete)
		}

		// Determine entity and entity ID from path
		path := c.Request.URL.Path
		entity := "unknown"
		var entityUUID uuid.UUID

		// Paths look like: /api/v1/menus or /api/v1/menus/123-456
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 3 {
			entity = parts[2]
			if len(parts) >= 4 {
				if id, err := uuid.Parse(parts[3]); err == nil {
					entityUUID = id
				}
			}
		}

		// Build and save the audit log
		log := models.AuditLog{
			UserID:    userUUID,
			Action:    action,
			Entity:    entity,
			OldValue:  "null",
			NewValue:  "null",
			Timestamp: time.Now(),
		}
		if len(reqBody) > 0 {
			log.NewValue = string(reqBody)
		}
		log.ID = uuid.New()
		log.CreatedAt = time.Now()
		log.UpdatedAt = time.Now()

		if entityUUID != uuid.Nil {
			log.EntityID = entityUUID
		}

		// Log background insert to avoid blocking client response
		go func(l models.AuditLog) {
			db.Create(&l)
		}(log)
	}
}
