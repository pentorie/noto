package handlers

import (
	"image"
	"image/jpeg"
	"noto/database"
	"noto/model"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type News model.News

func CreateNews(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.News)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)
	file, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "CoverRequired",
		})
	}

	src, _ := file.Open()

	filename := uuid.New().String() + ".jpg"
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "InvalidCoverExtension",
		})
	}

	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	// Resize the image to 1200px width
	resizedImg := resize.Resize(1200, 0, img, resize.Lanczos3)

	// Save the resized image to a new file
	out, err := os.Create("storage/covers/news/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the resized image to the file
	jpeg.Encode(out, resizedImg, nil)
	newNews := News{
		Author_id: uid,
		AuthLog:   claimer.Issuer,
		Content:   jsonx.Content,
		Title:     jsonx.Title,
		Cover:     "/covers/news/" + filename,
	}
	erx := db.Create(&newNews).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorNewsCreate",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "NewsCreated",
	})
}

func GetNews(c *fiber.Ctx) error {
	db := database.DB
	news := []model.NewsShort{}

	//пагинация
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if (pageSize <= 0) || (pageSize > 15) {
		pageSize = 15
	}
	offset := (page - 1) * pageSize

	err := db.Model(&model.News{}).Offset(offset).Limit(pageSize).Order("id asc").Find(&news).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "NewsNotFound",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    news,
	})
}

func GetSingleNews(c *fiber.Ctx) error {
	db := database.DB
	news := News{}

	err := db.Where("id = ?", c.Params("id")).Find(&news).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "NewsNotFound",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    news,
	})
}

func UpdateNews(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.News)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	news := News{}
	err := db.Where("id = ?", c.Params("id")).Find(&news).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "NewsNotFound",
		})
	}
	updatedNews := news
	if jsonx.Content != "" {
		updatedNews.Content = jsonx.Content
	}
	if jsonx.Title != "" {
		updatedNews.Title = jsonx.Title
	}
	filename := uuid.New().String() + ".jpg"
	if jsonx.Cover != "" {
		file, err := c.FormFile("cover")
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "CoverRequired",
			})
		}

		src, _ := file.Open()

		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "InvalidCoverExtension",
			})
		}

		img, _, err := image.Decode(src)
		if err != nil {
			return err
		}

		// Resize the image to 1200px width
		resizedImg := resize.Resize(1200, 0, img, resize.Lanczos3)

		// Save the resized image to a new file
		out, err := os.Create("storage/covers/news/" + filename)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the resized image to the file
		jpeg.Encode(out, resizedImg, nil)
	}
	erx := db.Model(&model.News{}).Where("id = ?", c.Params("id")).Updates(News{Content: updatedNews.Content, Title: updatedNews.Title, Cover: "/covers/news/" + filename}).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "NewsUpdateFailed",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "NewsUpdateSuccess",
	})
}

func DeleteNews(c *fiber.Ctx) error {
	db := database.DB
	claimer := RequesterTokenInfo(c)
	if claimer.IssuedAt == 2 {
		err := db.Where("id = ?", c.Params("id")).Delete(&model.News{})
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "NewsDeleteFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "NewsDeleteSuccess",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "NewsDeleteProhibited",
	})
}
