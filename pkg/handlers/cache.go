package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/unbot2313/distributed-cache/pkg/router"
)

// RegisterCacheHandlers registers all cache-related HTTP routes on the given mux.
func RegisterCacheHandlers(mux *http.ServeMux, cr *router.CacheRouter) {
	mux.HandleFunc("GET /cache/{key}", handleGetKey(cr))
	mux.HandleFunc("PUT /cache/{key}", handleSetKey(cr))
	mux.HandleFunc("DELETE /cache/{key}", handleDeleteKey(cr))
	mux.HandleFunc("GET /health", handleHealth(cr))
	mux.HandleFunc("GET /ring/info", handleRingInfo(cr))
}

func handleGetKey(cr *router.CacheRouter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		value, found, err := cr.Get(r.Context(), key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"key":   key,
			"value": value,
		})
	}
}

func handleSetKey(cr *router.CacheRouter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		value := string(body)
		if value == "" {
			http.Error(w, "body (value) is required", http.StatusBadRequest)
			return
		}

		if err := cr.Set(r.Context(), key, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nodeID, _ := cr.GetNodeForKey(key)
		writeJSON(w, http.StatusCreated, map[string]string{
			"key":    key,
			"status": "ok",
			"node":   nodeID,
		})
	}
}

func handleDeleteKey(cr *router.CacheRouter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		if err := cr.Delete(r.Context(), key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"key":    key,
			"status": "deleted",
		})
	}
}

func handleHealth(cr *router.CacheRouter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := cr.PingAll(r.Context())
		allHealthy := true

		nodeStatuses := make(map[string]string, len(results))
		for id, err := range results {
			if err != nil {
				nodeStatuses[id] = err.Error()
				allHealthy = false
			} else {
				nodeStatuses[id] = "ok"
			}
		}

		status := http.StatusOK
		overallStatus := "healthy"
		if !allHealthy {
			status = http.StatusServiceUnavailable
			overallStatus = "degraded"
		}

		writeJSON(w, status, map[string]interface{}{
			"status": overallStatus,
			"nodes":  nodeStatuses,
		})
	}
}

func handleRingInfo(cr *router.CacheRouter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := cr.GetRingInfo(10000)
		writeJSON(w, http.StatusOK, info)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
