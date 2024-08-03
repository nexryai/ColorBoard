package account

import (
	"github.com/nexryai/ColorBoard/db"
	"github.com/nexryai/ColorBoard/internal/database"
)

type UserServices struct{}

func (us *UserServices) CreateUser(user *db.UserModel) (string, error) {
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

func (us *UserServices) GetUser(param *db.UserEqualsUniqueWhereParam) (*db.UserModel, error) {
	prisma, ctx, err := database.GetPrismaClient()
	if err != nil {
		return nil, err
	}

	user, err := prisma.User.FindUnique(*param).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
