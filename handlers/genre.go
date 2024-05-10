package handlers

import (
	"noto/database"
	"noto/model"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type Genre model.Genre

func CreateGenre(c *fiber.Ctx) error {
	db := database.DB
	json := new(model.Genre)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	newGenre := Genre{
		Genre_name: json.Genre_name,
	}
	err := db.Create(&newGenre).Error
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "genre can be onlu uniqe",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data": newGenre,
	})
}

func GetGenre(c *fiber.Ctx) error {
	db := database.DB
	Genre := []Genre{}
	p := new(model.Genre)
	if err := c.QueryParser(p); err != nil {
		return err
	}
	stdLimit := 50
	if q, err := strconv.Atoi(c.Query("limit")); err == nil {
		if q < 50 {
			stdLimit = q
		}
	}
	//.Where("studio = ? AND year = ?", p.Studio, p.Year) - для примера потом расширять
	//когда будет готов рейтинг - допилить Order
	db.Model(&model.Genre{}).Order("ID asc").Limit(stdLimit).Where(p).Find(&Genre)

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    Genre,
	})
}