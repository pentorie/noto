package handlers

import (
	"image"
	"image/jpeg"
	"noto/database"
	"noto/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateAnime model.UpdateAnime
type UpdateManga model.UpdateManga
type UpdateRanobe model.UpdateRanobe

func UpdateAnimeEntity(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(Anime)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Unauthorized",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	ulogin := claims.Issuer

	found := model.AnimeWO{}
	query := Anime{ID: EntityIDFromSlug(c.Params("slug"), "anime")}
	err = db.Model(&model.Anime{}).First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "AnimeNotFound",
		})
	}
	if jsonx.Title != "" {
		found.Title = jsonx.Title
	}
	if jsonx.Studio != "" {
		found.Studio = jsonx.Studio
	}
	if jsonx.Slug != "" {
		jsonx.Slug = strings.Trim(jsonx.Slug, " /,.%():!;[]{}")
		found.Slug = strings.ToLower(jsonx.Slug)
	}
	if jsonx.Altt.Eng != "" {
		found.Altt.Eng = jsonx.Altt.Eng
	}
	if jsonx.Altt.Orig != "" {
		found.Altt.Orig = jsonx.Altt.Orig
	}
	if jsonx.Altt.Rus != "" {
		found.Altt.Rus = jsonx.Altt.Rus
	}
	if jsonx.Altt.Etc != "" {
		found.Altt.Etc = jsonx.Altt.Etc
	}
	if jsonx.Description.Eng != "" {
		found.Description.Eng = jsonx.Description.Eng
	}
	if jsonx.Description.Orig != "" {
		found.Description.Orig = jsonx.Description.Orig
	}
	if jsonx.Description.Rus != "" {
		found.Description.Rus = jsonx.Description.Rus
	}
	if jsonx.Type != "" {
		found.Type = jsonx.Type
	}
	file, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "CoverRequired",
		})
	}
	if jsonx.Cover != "" || err == nil {
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
		found.Cover = "/covers" + filename
	}
	if jsonx.Rating != "" {
		found.Rating = jsonx.Rating
	}
	if jsonx.Genres != nil {
		found.Genres = jsonx.Genres
	}
	if jsonx.Themes != nil {
		found.Themes = jsonx.Themes
	}
	if jsonx.AiredOn.IsZero() == false {
		found.AiredOn = jsonx.AiredOn
	}
	if jsonx.Episodes != 0 {
		found.Episodes = jsonx.Episodes
	}
	if jsonx.CurrEpisode != 0 {
		found.CurrEpisode = jsonx.CurrEpisode
	}
	if jsonx.Period != 0 {
		found.Period = jsonx.Period
	}
	if jsonx.NextEpisode.IsZero() == false {
		found.NextEpisode = jsonx.NextEpisode
	}
	if jsonx.Sources.Kodik != "" {
		found.Sources.Kodik = jsonx.Sources.Kodik
	}
	if jsonx.ExtLinks.Official != "" {
		found.ExtLinks.Official = jsonx.ExtLinks.Official
	}
	if jsonx.ExtLinks.Youtube != "" {
		found.ExtLinks.Youtube = jsonx.ExtLinks.Youtube
	}
	if jsonx.ExtLinks.Twitter != "" {
		found.ExtLinks.Twitter = jsonx.ExtLinks.Twitter
	}
	if jsonx.ExtLinks.Shikimori != "" {
		found.ExtLinks.Shikimori = jsonx.ExtLinks.Shikimori
	}
	if jsonx.ExtLinks.MAL != "" {
		found.ExtLinks.MAL = jsonx.ExtLinks.MAL
	}
	if jsonx.ExtLinks.Kitsu != "" {
		found.ExtLinks.Kitsu = jsonx.ExtLinks.Kitsu
	}
	if jsonx.ExtLinks.AniList != "" {
		found.ExtLinks.AniList = jsonx.ExtLinks.AniList
	}
	if claims.IssuedAt == 2 { //если роль запрашивающего - админ, то обновлять без премодерации
		erx := db.Model(&model.Anime{}).Where("id = ?", EntityIDFromSlug(c.Params("slug"), "anime")).Updates(&found).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "TitleUpdated",
		})
	}
	updquery := UpdateAnime{
		AuthorLog: ulogin,
		TitleID:   EntityIDFromSlug(c.Params("slug"), "anime"),
		Changes:   found,
	}
	erx := db.Create(&updquery).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UpdRequestFailed",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "UpdateRequestCreated",
	})
}

func UpdateMangaEntity(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.MangaWO)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Unauthorized",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	ulogin := claims.Issuer

	found := model.MangaWO{}
	query := Manga{ID: EntityIDFromSlug(c.Params("slug"), "manga")}
	err = db.Model(&model.Manga{}).First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "MangaNotFound",
		})
	}
	if jsonx.Title != "" {
		found.Title = jsonx.Title
	}
	if jsonx.Author != "" {
		found.Author = jsonx.Author
	}
	if jsonx.Slug != "" {
		jsonx.Slug = strings.Trim(jsonx.Slug, " /,.%():!;[]{}")
		found.Slug = strings.ToLower(jsonx.Slug)
	}
	if jsonx.Altt.Eng != "" {
		found.Altt.Eng = jsonx.Altt.Eng
	}
	if jsonx.Altt.Orig != "" {
		found.Altt.Orig = jsonx.Altt.Orig
	}
	if jsonx.Altt.Rus != "" {
		found.Altt.Rus = jsonx.Altt.Rus
	}
	if jsonx.Altt.Etc != "" {
		found.Altt.Etc = jsonx.Altt.Etc
	}
	if jsonx.Description.Eng != "" {
		found.Description.Eng = jsonx.Description.Eng
	}
	if jsonx.Description.Orig != "" {
		found.Description.Orig = jsonx.Description.Orig
	}
	if jsonx.Description.Rus != "" {
		found.Description.Rus = jsonx.Description.Rus
	}
	if jsonx.Type != "" {
		found.Type = jsonx.Type
	}
	file, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "CoverRequired",
		})
	}
	if jsonx.Cover != "" || err == nil {
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
		found.Cover = "/covers" + filename
	}
	if jsonx.Rating != "" {
		found.Rating = jsonx.Rating
	}
	if jsonx.Genres != nil {
		found.Genres = jsonx.Genres
	}
	if jsonx.Themes != nil {
		found.Themes = jsonx.Themes
	}
	if jsonx.AiredOn.IsZero() == false {
		found.AiredOn = jsonx.AiredOn
	}
	if jsonx.Episodes != 0 {
		found.Episodes = jsonx.Episodes
	}
	if jsonx.CurrEpisode != 0 {
		found.CurrEpisode = jsonx.CurrEpisode
	}
	if jsonx.Sources.Mangadex != "" {
		found.Sources.Mangadex = jsonx.Sources.Mangadex
	}
	if jsonx.Sources.Readmanga != "" {
		found.Sources.Readmanga = jsonx.Sources.Readmanga
	}
	if jsonx.ExtLinks.Official != "" {
		found.ExtLinks.Official = jsonx.ExtLinks.Official
	}
	if jsonx.ExtLinks.Youtube != "" {
		found.ExtLinks.Youtube = jsonx.ExtLinks.Youtube
	}
	if jsonx.ExtLinks.Twitter != "" {
		found.ExtLinks.Twitter = jsonx.ExtLinks.Twitter
	}
	if jsonx.ExtLinks.Shikimori != "" {
		found.ExtLinks.Shikimori = jsonx.ExtLinks.Shikimori
	}
	if jsonx.ExtLinks.MAL != "" {
		found.ExtLinks.MAL = jsonx.ExtLinks.MAL
	}
	if jsonx.ExtLinks.Kitsu != "" {
		found.ExtLinks.Kitsu = jsonx.ExtLinks.Kitsu
	}
	if jsonx.ExtLinks.AniList != "" {
		found.ExtLinks.AniList = jsonx.ExtLinks.AniList
	}
	if claims.IssuedAt == 2 { //если роль запрашивающего - админ, то обновлять без премодерации
		erx := db.Model(&model.Manga{}).Where("id = ?", EntityIDFromSlug(c.Params("slug"), "manga")).Updates(&found).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "TitleUpdated",
		})
	}
	updquery := UpdateManga{
		AuthorLog: ulogin,
		TitleID:   EntityIDFromSlug(c.Params("slug"), "manga"),
		Changes:   found,
	}
	erx := db.Create(&updquery).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UpdRequestFailed",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "UpdateRequestCreated",
	})
}

func UpdateRanobeEntity(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(Ranobe)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Unauthorized",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	ulogin := claims.Issuer

	found := model.RanobeWO{}
	query := Ranobe{ID: EntityIDFromSlug(c.Params("slug"), "ranobe")}
	err = db.Model(&model.Ranobe{}).First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "AnimeNotFound",
		})
	}
	if jsonx.Title != "" {
		found.Title = jsonx.Title
	}
	if jsonx.Author != "" {
		found.Author = jsonx.Author
	}
	if jsonx.Slug != "" {
		jsonx.Slug = strings.Trim(jsonx.Slug, " /,.%():!;[]{}")
		found.Slug = strings.ToLower(jsonx.Slug)
	}
	if jsonx.Altt.Eng != "" {
		found.Altt.Eng = jsonx.Altt.Eng
	}
	if jsonx.Altt.Orig != "" {
		found.Altt.Orig = jsonx.Altt.Orig
	}
	if jsonx.Altt.Rus != "" {
		found.Altt.Rus = jsonx.Altt.Rus
	}
	if jsonx.Altt.Etc != "" {
		found.Altt.Etc = jsonx.Altt.Etc
	}
	if jsonx.Description.Eng != "" {
		found.Description.Eng = jsonx.Description.Eng
	}
	if jsonx.Description.Orig != "" {
		found.Description.Orig = jsonx.Description.Orig
	}
	if jsonx.Description.Rus != "" {
		found.Description.Rus = jsonx.Description.Rus
	}
	if jsonx.Type != "" {
		found.Type = jsonx.Type
	}
	file, err := c.FormFile("cover")
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "CoverRequired",
		})
	}
	if jsonx.Cover != "" || err == nil {
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
		found.Cover = "/covers" + filename
	}
	if jsonx.Rating != "" {
		found.Rating = jsonx.Rating
	}
	if jsonx.Genres != nil {
		found.Genres = jsonx.Genres
	}
	if jsonx.Themes != nil {
		found.Themes = jsonx.Themes
	}
	if jsonx.AiredOn.IsZero() == false {
		found.AiredOn = jsonx.AiredOn
	}
	if jsonx.Episodes != 0 {
		found.Episodes = jsonx.Episodes
	}
	if jsonx.CurrEpisode != 0 {
		found.CurrEpisode = jsonx.CurrEpisode
	}
	if jsonx.Sources.Ranobes != "" {
		found.Sources.Ranobes = jsonx.Sources.Ranobes
	}
	if jsonx.ExtLinks.Official != "" {
		found.ExtLinks.Official = jsonx.ExtLinks.Official
	}
	if jsonx.ExtLinks.Youtube != "" {
		found.ExtLinks.Youtube = jsonx.ExtLinks.Youtube
	}
	if jsonx.ExtLinks.Twitter != "" {
		found.ExtLinks.Twitter = jsonx.ExtLinks.Twitter
	}
	if jsonx.ExtLinks.Shikimori != "" {
		found.ExtLinks.Shikimori = jsonx.ExtLinks.Shikimori
	}
	if jsonx.ExtLinks.MAL != "" {
		found.ExtLinks.MAL = jsonx.ExtLinks.MAL
	}
	if jsonx.ExtLinks.Kitsu != "" {
		found.ExtLinks.Kitsu = jsonx.ExtLinks.Kitsu
	}
	if jsonx.ExtLinks.AniList != "" {
		found.ExtLinks.AniList = jsonx.ExtLinks.AniList
	}
	if claims.IssuedAt == 2 { //если роль запрашивающего - админ, то обновлять без премодерации
		erx := db.Model(&model.Ranobe{}).Where("id = ?", EntityIDFromSlug(c.Params("slug"), "ranobe")).Updates(&found).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "TitleUpdated",
		})
	}
	updquery := UpdateRanobe{
		AuthorLog: ulogin,
		TitleID:   EntityIDFromSlug(c.Params("slug"), "ranobe"),
		Changes:   found,
	}
	erx := db.Create(&updquery).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UpdRequestFailed",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "UpdateRequestCreated",
	})
}

// обновления отдельного тайтла
func RequestTitleUpdates(c *fiber.Ctx) error {
	db := database.DB
	entity := c.Params("entity")
	switch entity {
	case "anime":
		titleid := EntityIDFromSlug(c.Params("slug"), "anime")
		upd := []model.UpdateAnimeShort{}
		err := db.Model(&model.UpdateAnime{}).Order("req_status asc").Where("title_id = ?", titleid).Find(&upd).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "UpdatesNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    upd,
		})
	case "manga":
		titleid := EntityIDFromSlug(c.Params("slug"), "manga")
		upd := []model.UpdateMangaShort{}
		err := db.Model(&model.UpdateManga{}).Order("req_status asc").Where("title_id = ?", titleid).Find(&upd).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "UpdatesNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    upd,
		})
	case "ranobe":
		titleid := EntityIDFromSlug(c.Params("slug"), "ranobe")
		upd := []model.UpdateRanobeShort{}
		err := db.Model(&model.UpdateRanobe{}).Order("req_status asc").Where("title_id = ?", titleid).Find(&upd).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "UpdatesNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    upd,
		})
	}
	return c.JSON(fiber.Map{
		"code":    400,
		"message": "UnknownEntity",
	})

}

// все заявки на обновление
func RequestAllUpdates(c *fiber.Ctx) error {
	db := database.DB
	entity := c.Params("entity")
	sort := c.Query("sort")
	var sortarr []string
	if sort == "" {
		sort = "createdat,desc"
	}
	sortarr = strings.Split(sort, ",")
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if (pageSize <= 0) || (pageSize > 60) {
		pageSize = 60
	}
	offset := (page - 1) * pageSize
	switch entity {
	case "anime":
		upd := model.UpdateAnimeShort{}
		err := db.Model(&model.UpdateAnime{}).Offset(offset).Limit(pageSize).Order(sortarr[0] + " " + sortarr[1]).Find(&upd).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "UpdatesNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    upd,
		})
	case "manga":
		upd := model.UpdateMangaShort{}
		err := db.Model(&model.UpdateManga{}).Offset(offset).Limit(pageSize).Order(sortarr[0] + " " + sortarr[1]).Find(&upd).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "UpdatesNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    upd,
		})
	case "ranobe":
		upd := model.UpdateRanobeShort{}
		err := db.Model(&model.UpdateRanobe{}).Offset(offset).Limit(pageSize).Order(sortarr[0] + " " + sortarr[1]).Find(&upd).Error
		if err == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{
				"code":    404,
				"message": "UpdatesNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "success",
			"data":    upd,
		})
	}
	return c.JSON(fiber.Map{
		"code":    400,
		"message": "UnknownEntity",
	})
}

// выводит обновление только по изменённым полям
func RequestUpdateComparision(c *fiber.Ctx) error {
	db := database.DB
	entity := c.Params("entity")
	updid, _ := strconv.Atoi(c.Params("updateid"))
	switch entity {
	case "anime":
		titleid := EntityIDFromSlug(c.Params("slug"), "anime")
		upd := UpdateAnime{}
		query1 := UpdateAnime{TitleID: titleid}
		anime := Anime{}
		query2 := Anime{ID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestNotFound",
			})
		}
		erx := db.Find(&anime, &query2).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "AnimeNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":      200,
			"message":   "success",
			"createdat": upd.CreatedAt,
			"titleid":   upd.TitleID,
			"author":    upd.AuthorLog,
			"status":    upd.ReqStatus,
			"current":   anime,
			"updated":   upd.Changes,
		})
	case "manga":
		titleid := EntityIDFromSlug(c.Params("slug"), "manga")
		upd := UpdateManga{}
		query1 := UpdateManga{TitleID: titleid}
		manga := Manga{}
		query2 := Manga{ID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestNotFound",
			})
		}
		erx := db.Find(&manga, &query2).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "AnimeNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":      200,
			"message":   "success",
			"createdat": upd.CreatedAt,
			"titleid":   upd.TitleID,
			"author":    upd.AuthorLog,
			"status":    upd.ReqStatus,
			"current":   manga,
			"updated":   upd.Changes,
		})
	case "ranobe":
		titleid := EntityIDFromSlug(c.Params("slug"), "ranobe")
		upd := UpdateRanobe{}
		query1 := UpdateRanobe{TitleID: titleid}
		ranobe := Ranobe{}
		query2 := Ranobe{ID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestNotFound",
			})
		}
		erx := db.Find(&ranobe, &query2).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "AnimeNotFound",
			})
		}
		return c.JSON(fiber.Map{
			"code":      200,
			"message":   "success",
			"createdat": upd.CreatedAt,
			"titleid":   upd.TitleID,
			"author":    upd.AuthorLog,
			"status":    upd.ReqStatus,
			"current":   ranobe,
			"updated":   upd.Changes,
		})
	}
	return c.JSON(fiber.Map{
		"code":    400,
		"message": "UnknownEntity",
	})
}

// одобрить реквест
func AcceptUpdate(c *fiber.Ctx) error {
	db := database.DB
	entity := c.Params("entity")
	updid, _ := strconv.Atoi(c.Params("updateid"))
	switch entity {
	case "anime":
		titleid := EntityIDFromSlug(c.Params("slug"), "anime")
		upd := UpdateAnime{}
		query1 := UpdateAnime{TitleID: titleid}
		anime := Anime{}
		query2 := Anime{ID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestNotFound",
			})
		}
		erx := db.Find(&anime, &query2).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "AnimeNotFound",
			})
		}
		ery := db.Save(&UpdateAnime{ID: updid, ReqStatus: 1}).Error
		if ery != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RequestUpdateFailed",
			})
		}
		erz := db.Model(&model.Anime{}).Where("id = ?", titleid).Updates(&upd.Changes).Error
		if erz != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "UpdatedSuccessfully",
		})
	case "manga":
		titleid := EntityIDFromSlug(c.Params("slug"), "manga")
		upd := UpdateManga{}
		query1 := UpdateManga{TitleID: titleid}
		manga := Manga{}
		query2 := Manga{ID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestNotFound",
			})
		}
		erx := db.Find(&manga, &query2).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "MangaNotFound",
			})
		}
		ery := db.Save(&UpdateManga{ID: updid, ReqStatus: 1}).Error
		if ery != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RequestUpdateFailed",
			})
		}
		erz := db.Model(&model.Manga{}).Where("id = ?", titleid).Updates(&upd.Changes).Error
		if erz != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "UpdatedSuccessfully",
		})
	case "ranobe":
		titleid := EntityIDFromSlug(c.Params("slug"), "ranobe")
		upd := UpdateRanobe{}
		query1 := UpdateRanobe{TitleID: titleid}
		ranobe := Ranobe{}
		query2 := Ranobe{ID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestNotFound",
			})
		}
		erx := db.Find(&ranobe, &query2).Error
		if erx != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RanobeNotFound",
			})
		}
		ery := db.Save(&UpdateRanobe{ID: updid, ReqStatus: 1}).Error
		if ery != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RequestUpdateFailed",
			})
		}
		erz := db.Model(&model.Ranobe{}).Where("id = ?", titleid).Updates(&upd.Changes).Error
		if erz != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "UpdateRequestFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "UpdatedSuccessfully",
		})
	}
	return c.JSON(fiber.Map{
		"code":    400,
		"message": "UnknownEntity",
	})
}

// отклонить реквест
func DeclineRequest(c *fiber.Ctx) error {
	db := database.DB
	entity := c.Params("entity")
	updid, _ := strconv.Atoi(c.Params("updateid"))
	switch entity {
	case "anime":
		titleid := EntityIDFromSlug(c.Params("slug"), "anime")
		upd := UpdateAnime{}
		query1 := UpdateAnime{TitleID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "DeclineRequestNotFound",
			})
		}
		ery := db.Save(&UpdateAnime{ID: updid, ReqStatus: 2}).Error
		if ery != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RequestDeclineFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "DeclinedSuccessfully",
		})
	case "manga":
		titleid := EntityIDFromSlug(c.Params("slug"), "manga")
		upd := UpdateManga{}
		query1 := UpdateManga{TitleID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "DeclineRequestNotFound",
			})
		}
		ery := db.Save(&UpdateManga{ID: updid, ReqStatus: 2}).Error
		if ery != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RequestDeclineFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "DeclinedSuccessfully",
		})
	case "ranobe":
		titleid := EntityIDFromSlug(c.Params("slug"), "ranobe")
		upd := UpdateRanobe{}
		query1 := UpdateRanobe{TitleID: titleid}

		err := db.Where("id = ?", updid).Find(&upd, &query1).Error
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "DeclineRequestNotFound",
			})
		}
		ery := db.Save(&UpdateRanobe{ID: updid, ReqStatus: 2}).Error
		if ery != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "RequestDeclineFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "DeclinedSuccessfully",
		})
	}
	return c.JSON(fiber.Map{
		"code":    400,
		"message": "UnknownEntity",
	})
}
