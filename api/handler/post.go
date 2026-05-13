package handler


import (		
	_ "immy-api/service"
	"immy-api/repo"
	_ "immy-api/model"
)

type PostHandler struct {
	PostRepo 	*repo.PostRepo
}