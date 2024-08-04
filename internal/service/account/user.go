package account

import (
	"github.com/nexryai/ColorBoard/db"
	"github.com/nexryai/ColorBoard/internal/database"
	"github.com/nexryai/ColorBoard/internal/service"
)

type UserServices struct{}

func (us *UserServices) CreateUser(user *service.UserCreateParam) (string, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return "", err
	}

	created, err := prisma.User.CreateOne(
		db.User.Name.Set(user.Name),
		db.User.AuthUID.Set(user.AuthUID),
		db.User.Galleries.Link(nil),
		db.User.AvatarURL.Set("hoge"),
	).Exec(ctx)

	if err != nil {
		return "", err
	}

	return created.ID, nil
}

func (us *UserServices) GetUser(param db.UserEqualsUniqueWhereParam) (*db.UserModel, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return nil, err
	}

	user, err := prisma.User.FindUnique(param).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServices) UpdateAvatarUrl(param db.UserEqualsUniqueWhereParam, avatarUrl string) error {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return err
	}

	_, err = prisma.User.FindUnique(param).Update(db.User.AvatarURL.Set(avatarUrl)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewUserServices() *UserServices {
	return &UserServices{}
}
