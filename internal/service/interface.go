package service

import (
	"io"

	"github.com/nexryai/ColorBoard/db"
)

type UserCreateParam struct {
	Name      string
	AuthUID   string
	AvatarUrl string
}

type GalleryCreateParam struct {
	Name   string
	UserId string
}

type (
	IUserService interface {
		CreateUser(user *UserCreateParam) (string, error)
		GetUser(param db.UserEqualsUniqueWhereParam) (*db.UserModel, error)
		UpdateAvatarUrl(param db.UserEqualsUniqueWhereParam, avatarUrl string) error
	}
	IGalleryService interface {
		CreateGallery(gallery *GalleryCreateParam) (string, error)
		GetGallery(param db.GalleryEqualsUniqueWhereParam) (*db.GalleryModel, error)
		/*AddImage(reader io.Reader, gallery *db.GalleryModel) (string, error)*/
	}
	IStorageService interface {
		CreateFile(reader io.Reader) (string, error)
		GetFileUrl(param db.ImageEqualsUniqueWhereParam) (string, error)
		DeleteFile(param db.ImageEqualsUniqueWhereParam) error
	}
)
