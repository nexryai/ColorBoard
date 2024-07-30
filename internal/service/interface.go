package service

import (
	"github.com/nexryai/ColorBoard/db"
	"io"
)

type (
	IUserService    interface{}
	IAuthService    interface{}
	IIdentService   interface{}
	IGalleryService interface {
		CreateGallery(gallery db.GalleryModel) (string, error)
		AddImage(reader io.Reader, gallery db.GalleryModel) (string, error)
	}
	IThumbnailService interface {
		GenerateThumbnail(reader io.Reader) (*[]byte, error)
	}
	IStorageService interface {
		CreateFile(reader io.Reader) (string, error)
		GetFileUrl(param db.ImageEqualsUniqueWhereParam) (string, error)
		DeleteFile(param db.ImageEqualsUniqueWhereParam) error
	}
)
