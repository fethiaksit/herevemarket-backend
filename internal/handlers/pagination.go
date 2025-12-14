package handlers

import (
	"fmt"
	"strconv"
)

const (
	defaultPage  int64 = 1
	defaultLimit int64 = 20
	maxLimit     int64 = 100
)

// parsePaginationParams parses common pagination query parameters with sane defaults and limits.
func parsePaginationParams(pageStr, limitStr string) (int64, int64, error) {
	page := defaultPage
	limit := defaultLimit

	if pageStr != "" {
		parsed, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil || parsed < 1 {
			return 0, 0, fmt.Errorf("invalid page parameter")
		}
		page = parsed
	}

	if limitStr != "" {
		parsed, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || parsed < 1 {
			return 0, 0, fmt.Errorf("invalid limit parameter")
		}
		if parsed > maxLimit {
			parsed = maxLimit
		}
		limit = parsed
	}

	return page, limit, nil
}
