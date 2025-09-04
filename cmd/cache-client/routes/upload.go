package routes

import (
	"net/http"

	"github.com/unbot2313/distributed-cache/cmd/cache-client/handler"
)

func InitializeUploadRoute(q handler.QueryHandler){
	http.HandleFunc("/put/", q.HandleUpload)
}