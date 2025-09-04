package services

// mover a otro
type UploadBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}


type uploadService struct {}

type UploadService interface {
	Upload(body UploadBody) error	
}

func newUploadService() UploadService {
	return &uploadService{}
}

func (*uploadService) Upload (body UploadBody) error {
	
	// upload
	
	return nil
}