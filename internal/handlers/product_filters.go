package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// parseCategoryFilters normalizes category query parameters, supporting both
// repeated "category" params and comma-separated values.
func parseCategoryFilters(c *gin.Context) []string {
	raw := c.QueryArray("category")
	if len(raw) == 0 {
		if single := strings.TrimSpace(c.Query("category")); single != "" {
			raw = []string{single}
		}
	}

	seen := make(map[string]struct{}, len(raw))
	categories := make([]string, 0, len(raw))

	for _, value := range raw {
		for _, part := range strings.Split(value, ",") {
			name := strings.TrimSpace(part)
			if name == "" {
				continue
			}
			if _, exists := seen[name]; exists {
				continue
			}
			seen[name] = struct{}{}
			categories = append(categories, name)
		}
	}

	return categories
}
