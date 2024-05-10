package handlers

import (
	"strconv"
	"time"

	"noto/database"
	"noto/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type user model.User

const SecretKey = "secret"

func Profile(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user model.User

	database.DB.Where("id = ?", claims.Id).First(&user)

	return c.JSON(user)
}

func Register(c *fiber.Ctx) error {
	db := database.DB
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) //GenerateFromPassword returns the bcrypt hash of the password at the given cost i.e. (14 in our case).

	newUser := user{
		Email:    data["email"],
		Login:    data["login"],
		Password: password,
	}
	err := db.Create(&newUser).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "NotUniqueLoginOrMail",
		})
	}

	return c.JSON(newUser)
}

func AddAdmin(c *fiber.Ctx) error {
	db := database.DB
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) //GenerateFromPassword returns the bcrypt hash of the password at the given cost i.e. (14 in our case).

	newUser := user{
		Email:    data["email"],
		Login:    data["login"],
		Role:     2,
		Password: password,
	}
	err := db.Create(&newUser).Error
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(newUser)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user model.User

	database.DB.Where("login = ?", data["login"]).First(&user) //Check the email is present in the DB

	if user.ID == 0 { //If the ID return is '0' then there is no such email present in the DB
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	} // If the email is present in the DB then compare the Passwords and if incorrect password then return error.

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        strconv.Itoa(int(user.ID)),            //id contains the ID of the user.
		IssuedAt:  int64(user.Role),                      //issuedAt contains the role of the user.
		Issuer:    user.Login,                            //issuer contains the login of the user.
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //Adds time to the token i.e. 24 hours.
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	expiration := time.Now().Add(24 * time.Hour)
	sesstype, erx := strconv.ParseBool(data["savesession"])
	if erx != nil {
		sesstype = false
	}
	if sesstype {
		expiration = time.Now().Add(30 * 24 * time.Hour)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expiration,
		Path:     "/",
		Domain:   ".noto.moe",
		SameSite: "none",
		HTTPOnly: true,
		Secure:   true,
	} //Creates the cookie to be passed.

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		Domain:   ".noto.moe",
		SameSite: "none",
		Expires:  time.Now().Add(-time.Hour * 1), //Sets the expiry time an hour ago in the past.
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})

}
