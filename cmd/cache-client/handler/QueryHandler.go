package handler

import (
	"encoding/json"
	"net/http"

	"github.com/unbot2313/distributed-cache/pkg/services"
)

type QueryHandler interface {
	HandleUpload(w http.ResponseWriter, r *http.Request)
}

type queryHandler struct {
	uploadService services.UploadService
}

func NewQueryHandler(uploadService services.UploadService) QueryHandler {
	return &queryHandler{
		uploadService: uploadService,
	}
}

func (q *queryHandler) HandleUpload(w http.ResponseWriter, r *http.Request){
	var body services.UploadBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response, err := json.Marshal(
			services.UploadResponse{
			Success: false,
			Message: "Invalid request body",
		})
		if err != nil {
			http.Error(w, "Error interno del servidor al parsear el response", 500)
			return
		}
		w.Write(response)
	}

	if err := q.uploadService.Upload(body); err != nil {

		response, err := json.Marshal( services.UploadResponse{
			Success: false,
			Message: "Failed to upload",
		})
		if err != nil {
			http.Error(w, "Error interno del servidor al parsear el response", 500)
			return
		}
		w.Write(response)
	}

	response, err := json.Marshal(
		services.UploadResponse{
		Success: true,
		Message: "Upload successful",
	})

	if err != nil {
			http.Error(w, "Error interno del servidor al parsear el response", 500)
			return
		}
		w.Write(response)

}
