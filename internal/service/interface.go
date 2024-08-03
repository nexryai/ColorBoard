package service

import (
	"github.com/nexryai/ColorBoard/db"
	"io"
)

type (
	IUserService interface {
		CreateUser(user *db.UserModel) (string, error)
		GetUser(param *db.UserEqualsUniqueWhereParam) (*db.UserModel, error)
	}
	IGalleryService interface {
		CreateGallery(gallery *db.GalleryModel) (string, error)
		AddImage(reader io.Reader, gallery *db.GalleryModel) (string, error)
	}
	IStorageService interface {
		CreateFile(reader io.Reader) (string, error)
		GetFileUrl(param *db.ImageEqualsUniqueWhereParam) (string, error)
		DeleteFile(param *db.ImageEqualsUniqueWhereParam) error
	}
)
