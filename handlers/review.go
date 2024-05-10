package handlers

import (
	"noto/database"
	"noto/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Review model.Review

func CreateReview(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.Review)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	if Contains(model.EntityTypes, c.Params("entity")) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "UnknownEntity"})
	}
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)
	newReview := Review{
		Title_id:   EntityIDFromSlug(c.Params("slug"), c.Params("entity")),
		Title_type: c.Params("entity"),
		Author_id:  uid,
		Author_log: claimer.Issuer,
		Content:    jsonx.Content,
		Rated:      jsonx.Rated,
	}
	err := db.Create(&newReview).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorReviewCreate",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "ReviewCreated",
	})
}

func GetTitleReviews(c *fiber.Ctx) error {
	db := database.DB
	review := []Review{}

	//пагинация
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("limit"))
	if (pageSize <= 0) || (pageSize > 5) {
		pageSize = 5
	}
	offset := (page - 1) * pageSize

	err := db.Offset(offset).Limit(pageSize).Order("id asc").Where("title_id = ? AND title_type = ?", EntityIDFromSlug(c.Params("slug"), c.Params("entity")), c.Params("entity")).Find(&review).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "ReviewsNotFound",
		})
	}
	reviews := make([]model.ReviewExp, len(review))

	for i := 0; i < len(review); i++ {
		usercache := User{}
		db.Where("id = ?", review[i].Author_id).Take(&usercache)
		reviews[i].ID = review[i].ID
		reviews[i].Title_id = review[i].Title_id
		reviews[i].Title_type = review[i].Title_type
		reviews[i].Author_id = review[i].Author_id
		reviews[i].Content = review[i].Content
		reviews[i].Rated = review[i].Rated
		reviews[i].AuthLog = usercache.Login
		reviews[i].AuthAvatar = usercache.Avatar
		reviews[i].AuthUname = usercache.Username
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    reviews,
	})
}

func GetAllReviews(c *fiber.Ctx) error {
	db := database.DB
	Review := []Review{}

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

	err := db.Offset(offset).Limit(pageSize).Order("id asc").Find(&Review).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "ReviewsNotFound",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    Review,
	})
}

func GetSingleReview(c *fiber.Ctx) error {
	db := database.DB
	review := Review{}

	err := db.Where("id = ?", c.Params("id")).Find(&review).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "ReviewNotFound",
		})
	}

	author := User{}
	db.Where("id = ?", review.Author_id).Take(&author)

	rvexp := model.ReviewExp{}
	rvexp.ID = review.ID
	rvexp.CreatedAt = review.CreatedAt
	rvexp.Title_id = review.Title_id
	rvexp.Title_type = review.Title_type
	rvexp.Author_id = review.Author_id
	rvexp.AuthUname = author.Username
	rvexp.AuthLog = author.Login
	rvexp.AuthAvatar = author.Avatar
	rvexp.Content = review.Content
	rvexp.Rated = review.Rated

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    rvexp,
	})
}

func UpdateReview(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.Review)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)

	review := Review{}
	err := db.Where("id = ?", c.Params("id")).Find(&review).Error
	if review.Author_id != uid || err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UpdateProhibited",
		})
	}
	updatedReview := review
	if jsonx.Content != "" {
		updatedReview.Content = jsonx.Content
	}
	if jsonx.Rated != 0 {
		updatedReview.Rated = jsonx.Rated
	}
	erx := db.Model(&model.Review{}).Where("id = ?", c.Params("id")).Updates(Review{Content: updatedReview.Content, Rated: updatedReview.Rated}).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ReviewUpdateFailed",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "ReviewUpdateSuccess",
	})
}

func DeleteReview(c *fiber.Ctx) error {
	db := database.DB
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)
	if claimer.IssuedAt == 2 {
		err := db.Where("id = ?", c.Params("id")).Delete(&model.Review{})
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "ReviewDeleteFailed",
			})
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "ReviewDeleteSuccess",
		})
	}
	err := db.Where("author_id = ? AND id = ?", uid, c.Params("id")).Delete(&model.Review{})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ReviewDeleteProhibited",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "ReviewDeleteSuccess",
	})
}
