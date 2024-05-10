package handlers

import (
	"noto/database"
	"noto/model"
	"path/filepath"
	"strconv"
	"image"
	"image/jpeg"
	"os"
	"github.com/nfnt/resize"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Settings model.Settings
type UserPublic model.UserPublic

func UserIDFromSlug(cslug string) int {
	if slugint, err := strconv.Atoi(cslug); err == nil {
		return slugint
	} else {
		db := database.DB
		user := User{}
		query := User{Login: cslug}
		err := db.Take(&user, &query).Error
		if err == gorm.ErrRecordNotFound {
			return 0
		} else {
			return user.ID
		}
	}
}

func GetProfile(c *fiber.Ctx) error {
	db := database.DB
	login := c.Params("login")
	user := UserPublic{}
	query := UserPublic{Login: login}
	err := db.Model(&model.User{}).Take(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func EditProfile(c *fiber.Ctx) error {
	//парсинг тела запроса на редактирование профиля
	jsonx := new(User)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	db := database.DB
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims) //claims.ID для получения ID пользователя
	uid, _ := strconv.Atoi(claims.Id)

	cache := User{}
	query := User{ID: uid}

	err = db.First(&cache, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "UserNotFound",
		})
	}

	//чеки для запросов через API чтобы пустые запросы на поля не изменяли данные
	file, err := c.FormFile("avatar")
	if err == nil {
		src, _ := file.Open()

	filename := uuid.New().String() + ".jpg"
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "InvalidAvatarExtension",
		})
	}
	
	img, _, err := image.Decode(src)
		if err != nil {
			return err
		}

		// Resize the image to 400px width
		resizedImg := resize.Resize(400, 0, img, resize.Lanczos3)

		// Save the resized image to a new file
		out, err := os.Create("storage/avatars/"+filename)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the resized image to the file
		jpeg.Encode(out, resizedImg, nil)
		cache.Avatar = "avatars/" + filename
	}
	
	if jsonx.Username != "" {
		cache.Username = jsonx.Username
	}
	if jsonx.Email != "" {
		cache.Email = jsonx.Email
	}

	cache.Description = jsonx.Description
	cache.About.Name = jsonx.About.Name
	cache.About.Age = jsonx.About.Age
	cache.About.Link = jsonx.About.Link
	cache.About.Gender = jsonx.About.Gender
	cache.About.City = jsonx.About.City

	db.Where("id = ?", uid).Save(&cache)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})

}

// get current user settings
func GetSettings(c *fiber.Ctx) error {
	db := database.DB
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims) //claims.ID для получения ID пользователя
	uid, _ := strconv.Atoi(claims.Id)

	settings := Settings{}
	query := Settings{UserID: uid}

	erx := db.Take(&settings, &query).Error
	if erx == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "UserNotFound",
		})
	}

	return c.Status(fiber.StatusOK).JSON(settings)
}

func EditSettings(c *fiber.Ctx) error {
	jsonx := new(Settings)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	db := database.DB
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims) //claims.ID для получения ID пользователя
	uid, _ := strconv.Atoi(claims.Id)

	cache := Settings{}
	query := Settings{UserID: uid}

	err = db.First(&cache, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "UserNotFound",
		})
	}

	if jsonx.ColorScheme != "" {
		cache.ColorScheme = jsonx.ColorScheme
	}
	if jsonx.Language.InterfaceLanguage != "" {
		cache.Language.InterfaceLanguage = jsonx.Language.InterfaceLanguage
	}
	if jsonx.Language.GenresLanguage != "" {
		cache.Language.GenresLanguage = jsonx.Language.GenresLanguage
	}
	if jsonx.Language.PrefMainTitle != "" {
		cache.Language.PrefMainTitle = jsonx.Language.PrefMainTitle
	}
	if jsonx.Language.PrefSubTitle != "" {
		cache.Language.PrefSubTitle = jsonx.Language.PrefSubTitle
	}
	if jsonx.Language.PrefDescription != "" {
		cache.Language.PrefDescription = jsonx.Language.PrefDescription
	}

	db.Where("user_id = ?", uid).Save(&cache)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}

// стандартные настройки сразу после регистрации
func UserDefaults(c *fiber.Ctx) error {
	db := database.DB

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Unauthorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims) //claims.ID для получения ID пользователя
	uid, _ := strconv.Atoi(claims.Id)

	defaults := Settings{}
	defaults.UserID = uid
	defaults.DataCache = false
	defaults.Language.GenresLanguage = "ru"
	defaults.Language.InterfaceLanguage = "ru"
	defaults.Language.PrefMainTitle = "ru"
	defaults.Language.PrefSubTitle = "jp-ro"
	defaults.Language.PrefDescription = "ru"

	// Add a unique constraint to the table on the user_id column

	// Specify the unique constraint in the ON CONFLICT clause
	cx := db.Model(&Settings{}).Where("user_id = ?", uid).Save(&defaults).Error
	if cx != nil {
		db.Create(&defaults)
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}
