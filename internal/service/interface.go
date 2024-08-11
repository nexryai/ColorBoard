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
		GetGallery(id string) (*db.GalleryModel, error)
		/*AddImage(reader io.Reader, gallery *db.GalleryModel) (string, error)*/
	}
	IStorageService interface {
		CreateFile(reader io.Reader) (string, error)
		GetFileUrl(param db.ImageEqualsUniqueWhereParam) (string, error)
		DeleteFile(param db.ImageEqualsUniqueWhereParam) error
	}
)
