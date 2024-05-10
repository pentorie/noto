package handlers

import (
	"strings"
	"time"

	"image"
	"image/jpeg"
	"noto/database"
	"noto/model"
	"os"
	"strconv"

	"github.com/nfnt/resize"

	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Anime model.Anime
type Genres model.Genres
type AnimeTimings model.AnimeTimings

func UpdateAnnounces(c *fiber.Ctx) error {
	db := database.DB
	animes := []AnimeTimings{}
	db.Model(&model.Anime{}).Where("status = ?", "announce").Find(&animes)
	for i := 0; i < len(animes); i++ {
		ctimeu := time.Unix(time.Now().Unix(), 0)
		cdate := ctimeu.Format("2 January 2006")
		ttime := animes[i].AiredOn.Format("2 January 2006")
		//если тайтл меняет статус с анонса на онгоинг
		if cdate == ttime {
			anime := Anime{
				ID:          animes[i].ID,
				CurrEpisode: 1,
				NextEpisode: animes[i].AiredOn.AddDate(0, 0, animes[i].Period),
				Status:      "ongoing",
			}
			db.Save(&anime)
		}
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "AnnouncedTitlesUpdated",
	})
}

func UpdateOngoings(c *fiber.Ctx) error {
	db := database.DB
	animes := []AnimeTimings{}
	db.Model(&model.Anime{}).Where("status = ?", "ongoing").Find(&animes)
	for i := 0; i < len(animes); i++ {
		//если текущая дата равна дате старта + периодичность * кол-во вышедших эпизодов
		if time.Now().Format("2 June 2000") == (animes[i].AiredOn.
			AddDate(0, 0, animes[i].Period*animes[i].CurrEpisode)).Format("2 June 2000") && (time.Now().Format("2 June 2000") != (animes[i].AiredOn.AddDate(0, 0, animes[i].Period*animes[i].Episodes)).Format("2 June 2000")) {
			anime := Anime{
				ID:          animes[i].ID,
				CurrEpisode: animes[i].CurrEpisode + 1,
				NextEpisode: animes[i].AiredOn.AddDate(0, 0, animes[i].Period*animes[i].CurrEpisode),
				Status:      "ongoing",
			}
			db.Save(&anime)
		}
		if time.Now().Format("2 June 2000") == (animes[i].AiredOn.AddDate(0, 0, animes[i].Period*animes[i].Episodes)).Format("2 June 2000") {
			anime := Anime{
				ID:          animes[i].ID,
				Status:      "aired",
				CurrEpisode: animes[i].Episodes,
			}
			db.Save(&anime)
		}
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "OngoingTitlesUpdated",
	})
}

func CreateAnime(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(Anime)
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
	out, err := os.Create("storage/covers/" + filename)
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
	if Contains(model.TypeAnime, jsonx.Type) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "Invalid type Anime"})
	}

	newAnime := Anime{
		Title:       jsonx.Title,
		Slug:        strings.ToLower(jsonx.Slug),
		Altt:        jsonx.Altt,
		Studio:      jsonx.Studio,
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
		Period:      jsonx.Period,
		NextEpisode: jsonx.NextEpisode,
	}
	err = db.Create(&newAnime).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorAnimeCreate",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "AnimeCreated",
	})
}

func GetAnimes(c *fiber.Ctx) error {
	db := database.DB
	Anime := []AnimeShort{}

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
	types := model.TypeAnime
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
	qb := db.Model(&model.Anime{}).Where("status IN ? AND type IN ?", statuses, types)
	for _, g := range genres {
		qb = qb.Where("genres LIKE ?", "%"+g+"%")
	}
	for _, q := range themes {
		qb = qb.Where("themes LIKE ?", "%"+q+"%")
	}

	qb.Offset(offset).Limit(pageSize).Order(orderreq[0] + " " + orderreq[1]).Find(&Anime)

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    Anime,
	})
}

func GetAnimeBySlug(c *fiber.Ctx) error {
	db := database.DB
	slug := c.Params("slug")
	anime := Anime{}
	if slugint, err := strconv.Atoi(slug); err == nil {
		query := Anime{ID: slugint}
		err := db.Take(&anime, &query).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "AnimeNotFound",
			})
		}
		return c.Status(fiber.StatusOK).JSON(anime)
	}
	query := Anime{Slug: slug}
	err := db.Take(&anime, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "AnimeNotFound",
		})
	}
	return c.Status(fiber.StatusOK).JSON(anime)
}
