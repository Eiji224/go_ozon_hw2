package router

import (
	"encoding/json"
	"net/http"
)

func getCartContent(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"user_id":  r.PathValue("user_id"),
		"is_works": "yes",
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
