package shared

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultListLimit  uint64 = 20
	DefaultListOffset uint64 = 0
)

type ListParams struct {
	Limit  uint64
	Offset uint64
	Search string
}

func NewListParamsFromRequest(r *http.Request) ListParams {
	return ListParams{
		Limit:  parseUintQuery(r, "limit", DefaultListLimit),
		Offset: parseUintQuery(r, "offset", DefaultListOffset),
		Search: strings.ToLower(strings.TrimSpace(r.URL.Query().Get("search"))),
	}
}

func parseUintQuery(r *http.Request, key string, fallback uint64) uint64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}
