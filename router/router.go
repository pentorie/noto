package router

import (
	"noto/handlers"
	"noto/middleware"

	"github.com/gofiber/fiber/v2"
)

func Initalize(router *fiber.App) {

	router.Use(middleware.Security)

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Noto.moe API Services")
	})

	router.Use(middleware.Json)

	//anime handlers
	anime := router.Group("/anime")
	anime.Post("/", middleware.AuthenticatedAdmin, handlers.CreateAnime)
	anime.Get("/", handlers.GetAnimes)
	anime.Get("/:slug", handlers.GetAnimeBySlug)

	manga := router.Group("/manga")
	manga.Post("/", middleware.AuthenticatedAdmin, handlers.CreateManga)
	manga.Get("/", handlers.GetMangas)
	manga.Get("/:slug", handlers.GetMangaBySlug)

	ranobe := router.Group("/ranobe")
	ranobe.Post("/", middleware.AuthenticatedAdmin, handlers.CreateRanobe)
	ranobe.Get("/", handlers.GetRanobes)
	ranobe.Get("/:slug", handlers.GetRanobeBySlug)

	//entity updates
	router.Get("/:entity/updates", middleware.Authenticated, handlers.RequestAllUpdates)
	router.Get("/:entity/:slug/upd", middleware.Authenticated, handlers.RequestTitleUpdates)
	router.Get("/:entity/:slug/upd/:updateid", middleware.Authenticated, handlers.RequestUpdateComparision)
	router.Patch("/:entity/:slug/upd/:updateid/accept", middleware.AuthenticatedAdmin, handlers.AcceptUpdate)
	router.Patch("/:entity/:slug/upd/:updateid/decline", middleware.AuthenticatedAdmin, handlers.DeclineRequest)

	anime.Post("/:slug", middleware.Authenticated, handlers.UpdateAnimeEntity)
	manga.Post("/:slug", middleware.Authenticated, handlers.UpdateMangaEntity)
	ranobe.Post("/:slug", middleware.Authenticated, handlers.UpdateRanobeEntity)

	//userlist and search
	router.Get("/:login/list/:entity", handlers.RetrieveList)
	router.Post("/:entity/:slug/onlist", middleware.Authenticated, handlers.AddEntityOnList)
	router.Post("/:entity/:slug/fromlist", middleware.Authenticated, handlers.RemoveEntityFromList)
	router.Get("/search", handlers.Search)

	//authentication
	authent := router.Group("/account")
	authent.Post("/register", handlers.Register)
	authent.Patch("/setdefaults", handlers.UserDefaults)
	authent.Post("/login", handlers.Login)
	authent.Post("/logout", handlers.Logout)
	authent.Post("/addAdmin", handlers.AddAdmin)
	authent.Get("/validate", middleware.Validate)

	router.Patch("/updongoing", handlers.UpdateOngoings)
	router.Patch("/updannounce", handlers.UpdateAnnounces)

	//profile handlers
	router.Get("/u/:login", handlers.GetProfile)                                   //профиль пользователя
	router.Patch("/u/:login/edit", middleware.Authenticated, handlers.EditProfile) //редактирование профиля
	router.Get("/settings", middleware.Authenticated, handlers.GetSettings)
	router.Patch("/settings", middleware.Authenticated, handlers.EditSettings) //редактирование настроек

	//news handlers
	news := router.Group("/news")
	news.Post("/", middleware.AuthenticatedAdmin, handlers.CreateNews)
	news.Get("/", handlers.GetNews)
	news.Get("/:id", handlers.GetSingleNews)                                //получить одну новость
	news.Patch("/:id", handlers.UpdateNews, middleware.AuthenticatedAdmin)  //обновить новость
	news.Delete("/:id", handlers.DeleteNews, middleware.AuthenticatedAdmin) //удалить новость

	//comment handlers
	comment := router.Group("/comment")
	router.Post("/:entity/:slug/comment", middleware.Authenticated, handlers.CreateComment)
	router.Get("/:entity/:slug/comments", handlers.GetEntityComments) //получить несколько комментариев
	comment.Get("/:id", handlers.GetSingleComment)
	comment.Patch("/:id", middleware.Authenticated, handlers.UpdateComment)
	comment.Delete("/:id", middleware.Authenticated, handlers.DeleteComment)

	//review handlers
	review := router.Group("/review")
	router.Post("/:entity/:slug/review", middleware.Authenticated, handlers.CreateReview)
	router.Get("/:entity/:slug/reviews", handlers.GetTitleReviews)
	review.Patch("/:id", middleware.Authenticated, handlers.UpdateReview)
	review.Get(":id", handlers.GetSingleReview)
	review.Delete("/:id", middleware.Authenticated, handlers.DeleteReview)

	genre := router.Group("/genre")
	genre.Post("/", middleware.AuthenticatedAdmin, handlers.CreateGenre)
	genre.Get("/", handlers.GetGenre)

	router.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})

}
