package handlers

import (
	"noto/database"
	"noto/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User model.User
type Marks model.Marks
type Altt model.Altt

func Contains(a []string, x string) bool {
	for _, n := range a {
		if n == x {
			return true
		}
	}
	return false
}

func RequesterTokenInfo(c *fiber.Ctx) *jwt.StandardClaims {
	cookie := c.Cookies("jwt")

	token, _ := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})

	return token.Claims.(*jwt.StandardClaims)
}

func EntityConstruct(ent string) interface{} {
	switch ent {
	case "anime":
		return &model.Anime{}
	case "manga":
		return &model.Manga{}
	case "ranobe":
		return &model.Ranobe{}
	}
	return nil
}

func EntityConstructWID(ent string, idx int) interface{} {
	switch ent {
	case "anime":
		return model.Anime{ID: idx}
	case "manga":
		return model.Manga{ID: idx}
	case "ranobe":
		return model.Ranobe{ID: idx}
	}
	return nil
}

func ListedQueryConstruct(ent string, qty int) interface{} {
	switch ent {
	case "anime":
		return make([]model.Anime, qty)
	case "manga":
		return make([]model.Manga, qty)
	case "ranobe":
		return make([]model.Ranobe, qty)
	}
	return nil
}

func EntityIDFromSlug(cslug string, ent string) int {
	if slugint, err := strconv.Atoi(cslug); err == nil {
		return slugint
	} else {
		db := database.DB
		switch ent {
		case "anime":
			anime := Anime{}
			query := Anime{Slug: cslug}
			err := db.Take(&anime, &query).Error
			if err == gorm.ErrRecordNotFound {
				return 0
			} else {
				return anime.ID
			}
		case "manga":
			manga := Manga{}
			query := Manga{Slug: cslug}
			err := db.Take(&manga, &query).Error
			if err == gorm.ErrRecordNotFound {
				return 0
			} else {
				return manga.ID
			}
		case "ranobe":
			ranobe := Ranobe{}
			query := Ranobe{Slug: cslug}
			err := db.Take(&ranobe, &query).Error
			if err == gorm.ErrRecordNotFound {
				return 0
			} else {
				return ranobe.ID
			}
		}
		return 0
	}
}
