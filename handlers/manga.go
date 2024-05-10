package handlers

import (
	"strings"

	"noto/database"
	"noto/model"
	"strconv"
	"image"
	"image/jpeg"
	"os"
	"github.com/nfnt/resize"

	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateManga(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(Manga)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "InvalidJSON",
		})
	}
	jsonx.Slug = strings.Trim(jsonx.Slug, " /,.%():!;[]{}")

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

		// Resize the image to 400px width
		resizedImg := resize.Resize(500, 0, img, resize.Lanczos3)

		// Save the resized image to a new file
		out, err := os.Create("storage/covers/"+filename)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the resized image to the file
		jpeg.Encode(out, resizedImg, nil)
	if Contains(model.StatusTypes, jsonx.Status) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "Invalid status"})
	}
	if Contains(model.Rating, jsonx.Rating) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "Invalid rating"})
	}
	if Contains(model.TypeManga, jsonx.Type) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "Invalid type Manga"})
	}

	newManga := Manga{
		Title:       jsonx.Title,
		Slug:        strings.ToLower(jsonx.Slug),
		Altt:    jsonx.Altt,
		Author:      jsonx.Author,
		Type:        jsonx.Type,
		Status:      jsonx.Status,
		Rating:      jsonx.Rating,
		Cover:       "covers/" + filename,
		Genres:      jsonx.Genres,
		Themes:      jsonx.Themes,
		Duration:    jsonx.Duration,
		Sources:     jsonx.Sources,
		Description: jsonx.Description,
		ExtLinks:    jsonx.ExtLinks,
		AiredOn:     jsonx.AiredOn,
		AiredEnd:    jsonx.AiredEnd,
		Episodes:    jsonx.Episodes,
		CurrEpisode: jsonx.CurrEpisode,
	}
	err = db.Create(&newManga).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorMangaCreate",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "MangaCreated",
	})
}

func GetMangas(c *fiber.Ctx) error {
	db := database.DB
	Manga := []model.MangaShort{}
	p := new(model.Manga)
	if err := c.QueryParser(p); err != nil {
		return err
	}

	//пагинация
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if (pageSize <= 0) || (pageSize > 60) {
		pageSize = 60
	}
	offset := (page - 1) * pageSize

	orderreq := strings.Split(c.Query("order"), ",")
	if len(orderreq) < 2 || c.Query("order") == "" {
		orderreq = []string{"mark_mean", "desc"}
	}

	//изобретено укропами но работает нормально: проверка статуса и типа из модели attr.go
	statusesstr := c.Query("status")
	statuses := model.StatusTypes
	if statusesstr != "" {
		statuses = strings.Split(statusesstr, ",")
	}

	typesstr := c.Query("type")
	types := model.TypeManga
	if typesstr != "" {
		types = strings.Split(typesstr, ",")
	}

	genresstr := c.Query("genre")

	var genres []string
	if genresstr != "" {
		genres = strings.Split(genresstr, ",")
	}

	themesstr := c.Query("theme")
	var themes []string
	if themesstr != "" {
		themes = strings.Split(themesstr, ",")
	}

	// query builder object and chain the conditions
	qb := db.Model(&model.Manga{}).Where("status IN ? AND type IN ?", statuses, types)
	for _, g := range genres {
		qb = qb.Where("genres LIKE ?", "%"+g+"%")
	}
	for _, q := range themes {
		qb = qb.Where("themes LIKE ?", "%"+q+"%")
	}

	qb.Offset(offset).Limit(pageSize).Order(orderreq[0] + " " + orderreq[1]).Find(&Manga)

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    Manga,
	})
}

func GetMangaBySlug(c *fiber.Ctx) error {
	db := database.DB
	slug := c.Params("slug")
	manga := Manga{}
	if slugint, err := strconv.Atoi(slug); err == nil {
		query := Manga{ID: slugint}
		err := db.Take(&manga, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "MangaNotFound",
			})
		}
		return c.Status(fiber.StatusOK).JSON(manga)
	}
	query := Manga{Slug: slug}
	err := db.Take(&manga, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "MangaNotFound",
		})
	}
	return c.Status(fiber.StatusOK).JSON(manga)
}
