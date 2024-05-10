package handlers

import (
	"noto/database"
	"noto/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Comment model.Comment

func CreateComment(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.Comment)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	if Contains(model.EntitiesExt, c.Params("entity")) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "UnknownEntity"})
	}
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)
	newComment := Comment{
		Title_id:   EntityIDFromSlug(c.Params("slug"), c.Params("entity")),
		Title_type: c.Params("entity"),
		Author_id:  uid,
		Author_log: claimer.Issuer,
		Content:    jsonx.Content,
	}
	err := db.Create(&newComment).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorCommentCreate",
		})
	}
	if c.Params("entity") == "news" {
		newsid, _ := strconv.Atoi(c.Params("id"))
		cachenews := News{}
		db.Where("id = ?", newsid).Find(&cachenews)
		cachenews.CommentQty++
		db.Where("id = ?", newsid).Update("commentqty", cachenews.CommentQty)
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "CommentCreated",
	})
}

func GetEntityComments(c *fiber.Ctx) error {
	db := database.DB
	comment := []Comment{}

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

	err := db.Offset(offset).Limit(pageSize).Order("id asc").Where("title_id = ? AND title_type = ?", EntityIDFromSlug(c.Params("slug"), c.Params("entity")), c.Params("entity")).Find(&comment).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "CommentsNotFound",
		})
	}
	comments := make([]model.CommentExp, len(comment))

	for i := 0; i < len(comment); i++ {
		usercache := User{}
		db.Where("id = ?", comment[i].Author_id).Take(&usercache)
		comments[i].ID = comment[i].ID
		comments[i].Title_id = comment[i].Title_id
		comments[i].Title_type = comment[i].Title_type
		comments[i].Author_id = comment[i].Author_id
		comments[i].Content = comment[i].Content
		comments[i].AuthLog = usercache.Login
		comments[i].AuthAvatar = usercache.Avatar
		comments[i].AuthUname = usercache.Username
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    comments,
	})
}

func GetSingleComment(c *fiber.Ctx) error {
	db := database.DB
	comment := Comment{}

	err := db.Where("id = ?", c.Params("id")).Find(&comment).Error
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "CommentNotFound",
		})
	}

	author := User{}
	db.Where("id = ?", comment.Author_id).Take(&author)

	rvexp := model.CommentExp{}
	rvexp.ID = comment.ID
	rvexp.CreatedAt = comment.CreatedAt
	rvexp.Title_id = comment.Title_id
	rvexp.Title_type = comment.Title_type
	rvexp.Author_id = comment.Author_id
	rvexp.AuthUname = author.Username
	rvexp.AuthLog = author.Login
	rvexp.AuthAvatar = author.Avatar
	rvexp.Content = comment.Content

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    rvexp,
	})
}

func UpdateComment(c *fiber.Ctx) error {
	db := database.DB
	jsonx := new(model.Comment)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)

	comment := Comment{}
	err := db.Where("id = ?", c.Params("id")).Find(&comment).Error
	if comment.Author_id != uid || err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UpdateProhibited",
		})
	}
	updatedComment := comment
	if jsonx.Content != "" {
		updatedComment.Content = jsonx.Content
	}
	erx := db.Model(&model.Comment{}).Where("id = ?", c.Params("id")).Updates(Comment{Content: updatedComment.Content}).Error
	if erx != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "CommentUpdateFailed",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "CommentUpdateSuccess",
	})
}

func DeleteComment(c *fiber.Ctx) error {
	db := database.DB
	claimer := RequesterTokenInfo(c)
	uid, _ := strconv.Atoi(claimer.Id)
	if claimer.IssuedAt == 2 {
		err := db.Where("id = ?", c.Params("id")).Delete(&model.Comment{})
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "CommentDeleteFailed",
			})
		}
		if c.Params("entity") == "news" {
			newsid, _ := strconv.Atoi(c.Params("id"))
			cachenews := News{}
			db.Where("id = ?", newsid).Find(&cachenews)
			cachenews.CommentQty--
			db.Where("id = ?", newsid).Update("commentqty", cachenews.CommentQty)
		}
		return c.JSON(fiber.Map{
			"code":    200,
			"message": "CommentDeleteSuccess",
		})
	}
	err := db.Where("author_id = ? AND id = ?", uid, c.Params("id")).Delete(&model.Comment{})
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "CommentDeleteProhibited",
		})
	}
	if c.Params("entity") == "news" {
		newsid, _ := strconv.Atoi(c.Params("id"))
		cachenews := News{}
		db.Where("id = ?", newsid).Find(&cachenews)
		cachenews.CommentQty--
		db.Where("id = ?", newsid).Update("commentqty", cachenews.CommentQty)
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "CommentDeleteSuccess",
	})
}
