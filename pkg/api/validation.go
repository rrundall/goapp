package api

import (
	"fmt"
	"goapp/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ValidateContentType will check for required content-type in header and return false if not exist
func ValidateContentType(c *gin.Context) bool {
	if ct := c.Request.Header.Get("Content-Type"); ct != "application/json" {
		log.Error().Msgf("%s: %s", config.UnsupportedContentType, ct)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": config.UnsupportedContentType})
		return false
	}
	return true
}

// WarnFieldsCannotBeEmpty will display warning for fields that did not get updates
func WarnFieldsCannotBeEmpty(f []string) string {
	if len(f) > 0 {
		log.Warn().Msgf("%s %s", config.FieldsBeEmptyWarningMsg,
			strings.ReplaceAll(strings.Join(f, ", "), ",omitempty", ""))
		return fmt.Sprintf("%s %s", config.FieldsBeEmptyWarningMsg,
			strings.ReplaceAll(strings.Join(f, ", "), ",omitempty", ""))
	}
	return ""
}

// ValidateRowsAffected will define what message need to show up either success data update or no data update
// base on row affected number.
func ValidateRowsAffected(c *gin.Context, rowsAffected int64, msg string) {
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"message": config.NoDataUpdateWarningMsg, "rows_affected": rowsAffected})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": msg, "rows_affected": rowsAffected})
	}
}

// WarnEmptyData will display warning for no data to pass in query
func WarnEmptyData(c *gin.Context, f []string) bool {
	if len(f) == 0 {
		log.Warn().Msgf("%s %s", config.NoQueryDataPassedWarningMsg)
		c.JSON(http.StatusOK, gin.H{"message": config.NoQueryDataPassedWarningMsg})
		return false
	}
	return true
}
