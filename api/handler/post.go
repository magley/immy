package handler


import (		
	"immy-api/service"
	_ "immy-api/model"
)

type PostHandler struct {
	PostService 	*service.PostService
}