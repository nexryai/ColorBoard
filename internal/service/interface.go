package service

import (
	"github.com/nexryai/ColorBoard/db"
	"io"
)

type UserCreateParam struct {
	Name      string
	AuthUID   string
	AvatarUrl string
}

type (
	IUserService interface {
		CreateUser(user *UserCreateParam) (string, error)
		GetUser(param db.UserEqualsUniqueWhereParam) (*db.UserModel, error)
		UpdateAvatarUrl(param db.UserEqualsUniqueWhereParam, avatarUrl string) error
	}
	IGalleryService interface {
		CreateGallery(gallery *db.GalleryModel) (string, error)
		AddImage(reader io.Reader, gallery *db.GalleryModel) (string, error)
	}
	IStorageService interface {
		CreateFile(reader io.Reader) (string, error)
		GetFileUrl(param db.ImageEqualsUniqueWhereParam) (string, error)
		DeleteFile(param db.ImageEqualsUniqueWhereParam) error
	}
)
