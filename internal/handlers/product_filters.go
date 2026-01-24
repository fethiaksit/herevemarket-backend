package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func parseCategoryIDQuery(c *gin.Context) ([]primitive.ObjectID, error) {
	rawValues := c.QueryArray("categoryId")
	if len(rawValues) == 0 {
		if raw := strings.TrimSpace(c.Query("categoryId")); raw != "" {
			rawValues = []string{raw}
		}
	}

	if len(rawValues) == 0 {
		return nil, nil
	}

	values := make([]string, 0, len(rawValues))
	for _, raw := range rawValues {
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			value := strings.TrimSpace(part)
			if value == "" {
				continue
			}
			values = append(values, value)
		}
	}

	if len(values) == 0 {
		return nil, nil
	}

	return parseCategoryObjectIDs(values)
}
