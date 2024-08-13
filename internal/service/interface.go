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
	Name     string
	IsPublic bool
	UserId   string
}

type (
	IUserService interface {
		CreateUser(user *UserCreateParam) (string, error)
		// ToDo: paramやめて識別子をuserIdで統一する（現状AuthUIDで見つけてる部分がある）
		GetUser(param db.UserEqualsUniqueWhereParam) (*db.UserModel, error)
		UpdateAvatarUrl(param db.UserEqualsUniqueWhereParam, avatarUrl string) error
	}
	IGalleryService interface {
		CreateGallery(gallery *GalleryCreateParam) (string, error)
		GetGallery(userId string, id string) (*db.GalleryModel, error)
		GetGalleriesByUserId(userId string) (*[]db.GalleryModel, error)
		AddImage(reader io.Reader, thumbReader io.Reader, userId string, galleryId string, blurash string, w int, h int) (string, error)
	}
	IStorage interface {
		CreateFile(reader io.Reader, userId string) (string, error)
		GetFileUrl(id string, userId string) (string, error)
		DeleteFile(id string, userId string) error
	}
)
