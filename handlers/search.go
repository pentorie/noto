package handlers

import (
	"noto/database"
	"noto/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Search(c *fiber.Ctx) error {
	db := database.DB
	entity := c.Query("entity")
	title := c.Query("title")
	switch entity {
	case "anime":
		data := []model.AnimeShort{}
		err1 := db.Model(&model.Anime{}).Where("title ILIKE ?", "%"+title+"%").
			Or("altt->>'rus' ILIKE ?", "%"+title+"%").
			Or("altt->>'eng' ILIKE ?", "%"+title+"%").
			Or("altt->>'orig' ILIKE ?", "%"+title+"%").
			Or("altt->>'etc' ILIKE ?", "%"+title+"%").Find(&data).Error
		if err1 != gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    200,
				"message": "success",
				"data":    data,
			})
		}
	case "manga":
		data := []model.MangaShort{}
		err1 := db.Model(&model.Manga{}).Where("title ILIKE ?", "%"+title+"%").
			Or("altt->>'rus' ILIKE ?", "%"+title+"%").
			Or("altt->>'eng' ILIKE ?", "%"+title+"%").
			Or("altt->>'orig' ILIKE ?", "%"+title+"%").
			Or("altt->>'etc' ILIKE ?", "%"+title+"%").Find(&data).Error
		if err1 != gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    200,
				"message": "success",
				"data":    data,
			})
		}
	case "ranobe":
		data := []model.RanobeShort{}
		err1 := db.Model(&model.Ranobe{}).Where("title ILIKE ?", "%"+title+"%").
			Or("altt->>'rus' ILIKE ?", "%"+title+"%").
			Or("altt->>'eng' ILIKE ?", "%"+title+"%").
			Or("altt->>'orig' ILIKE ?", "%"+title+"%").
			Or("altt->>'etc' ILIKE ?", "%"+title+"%").Find(&data).Error
		if err1 != gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    200,
				"message": "success",
				"data":    data,
			})
		}

	}

	return c.JSON(fiber.Map{
		"code":    404,
		"message": "IncorrectEntity",
	})
}
