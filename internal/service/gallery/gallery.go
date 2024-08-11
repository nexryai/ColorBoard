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

func (gs *GalleryService) GetGallery(userId string, id string) (*db.GalleryModel, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return nil, err
	} else {
		defer prisma.Prisma.Disconnect()
	}

	// FindUniqueが使えない
	found, err := prisma.Gallery.FindFirst(
		db.Gallery.And(
			db.Gallery.ID.Equals(id),
			db.Gallery.UserID.Equals(userId),
		),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return found, nil
}

func (gs *GalleryService) GetGalleriesByUserId(userId string) (*[]db.GalleryModel, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return nil, err
	} else {
		defer prisma.Prisma.Disconnect()
	}

	found, err := prisma.Gallery.FindMany(
		db.Gallery.UserID.Contains(userId),
	).With(
		// Used for a thumbnail
		db.Gallery.Images.Fetch().Take(1),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func NewGalleryService() *GalleryService {
	return &GalleryService{}
}