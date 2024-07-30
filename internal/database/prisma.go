package database

import "github.com/nexryai/ColorBoard/db"

func GetPrismaClient() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}
