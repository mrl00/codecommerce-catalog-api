package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"codecommerceapi/internal/service"
)

func parsePaginationParams(r *http.Request) service.PaginationParams {
	q := r.URL.Query()
	page := 1
	perPage := 10

	if v := q.Get("page"); v != "" {
		fmt.Sscanf(v, "%d", &page)
	}
	if v := q.Get("per_page"); v != "" {
		fmt.Sscanf(v, "%d", &perPage)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return service.PaginationParams{Page: page, PerPage: perPage}
}

func errorResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
