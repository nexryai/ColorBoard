package gallery

import (
	"io"

	"github.com/nexryai/ColorBoard/db"
	"github.com/nexryai/ColorBoard/internal/database"
	"github.com/nexryai/ColorBoard/internal/service"
)

type GalleryService struct {
	storage service.IStorage
}

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

func (gs *GalleryService) GetGallery(userId string, id string, page int) (*db.GalleryModel, error) {
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
	).With(
		db.Gallery.Images.Fetch().OrderBy(db.Image.ID.Order(db.DESC)).Skip((page*20)-20,).Take(20),
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
	).OrderBy(
		db.Gallery.ID.Order(db.DESC),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (gs *GalleryService) AddImage(reader io.Reader, sha256Hash string, thumbReader io.Reader, userId string, galleryId string, blurhash string, w int, h int) (string, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return "", err
	} else {
		defer prisma.Prisma.Disconnect()
	}

	fileId, err := gs.storage.CreateFile(reader, userId)
	if err != nil {
		return "", err
	}

	thumbnailId, err := gs.storage.CreateFile(thumbReader, userId)
	if err != nil {
		return "", err
	}

	created, err := prisma.Image.CreateOne(
		db.Image.StorageKey.Set(fileId),
		db.Image.Sha256Hash.Set(sha256Hash),
		db.Image.ThumbnailKey.Set(thumbnailId),
		db.Image.Blurhash.Set(blurhash),
		db.Image.User.Link(
			db.User.ID.Equals(userId),
		),
		db.Image.Gallery.Link(
			db.Gallery.ID.Equals(galleryId),
		),
		db.Image.Width.Set(w),
		db.Image.Height.Set(h),
	).Exec(ctx)

	if err != nil {
		return "", err
	}

	return created.ID, nil
}

func NewGalleryService(storage service.IStorage) *GalleryService {
	return &GalleryService{storage: storage}
}
