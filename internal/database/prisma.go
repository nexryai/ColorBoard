package database

import (
	"context"
	"github.com/nexryai/ColorBoard/db"
)

func GetPrismaClient() (*db.PrismaClient, context.Context, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, nil, err
	}

	return client, context.Background(), nil
}
