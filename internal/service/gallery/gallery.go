package gallery

import (
	"github.com/nexryai/ColorBoard/db"
	"github.com/nexryai/ColorBoard/internal/database"
	"github.com/nexryai/ColorBoard/internal/service"
)

type GalleryService struct {}

func (gs *GalleryService) CreateGallery(gallery *service.GalleryCreateParam) (string, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return "", err
	} else {
		defer prisma.Prisma.Disconnect()
	}

	created, err := prisma.Gallery.CreateOne(
		db.Gallery.Name.Set(gallery.Name),
		db.Gallery.User.Link(
			db.User.ID.Equals(gallery.UserId),
		),
		db.Gallery.IsPublic.Set(gallery.IsPublic),
	).Exec(ctx)

	if err != nil {
		return "", err
	}

	return created.ID, nil
}

func (gs *GalleryService) GetGallery(param db.GalleryEqualsUniqueWhereParam) (*db.GalleryModel, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return nil, err
	} else {
		defer prisma.Prisma.Disconnect()
	}

	user, err := prisma.Gallery.FindUnique(param).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewGalleryService() *GalleryService {
	return &GalleryService{}
}