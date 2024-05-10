package database

import (
	"log"
	"os"

	"noto/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectDB() {
	var err error // define error here to prevent overshadowing the global DB

	env := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(env), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(
		&model.User{}, &model.Session{}, &model.Product{}, &model.Anime{},
		&model.Ranobe{}, &model.Manga{}, &model.News{}, &model.Review{},
		&model.Genre{}, &model.Comment{}, &model.Settings{}, &model.UpdateAnime{},
		&model.UpdateManga{}, &model.UpdateRanobe{})
	if err != nil {
		log.Fatal(err)
	}

}
