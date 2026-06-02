package main

import (
	_ "Q38/docs"

	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

func main() {
	app := iris.New()

	// Настройка Swagger UI
	config := &swagger.Config{
		URL: "http://localhost:8080/swagger/doc.json",
	}

	// Swagger UI маршруты
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
	app.Get("/swagger", swagger.CustomWrapHandler(config, swaggerFiles.Handler))

	// Отдаём статическую HTML-страницу
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("./static/index.html")
	})

	// API эндпоинт
	app.Get("/api/v1/pets/{id:int}", getPet)

	app.Listen(":8080")
}

// Pet - структура для ответа
type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse - структура ошибки
type ErrorResponse struct {
	Message string `json:"message"`
}

// @Summary Get a pet by ID
// @Description Returns a single pet
// @Tags pets
// @Accept json
// @Produce json
// @Param id path int true "Pet ID"
// @Success 200 {object} Pet
// @Failure 400 {object} ErrorResponse
// @Router /pets/{id} [get]
func getPet(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")
	ctx.JSON(Pet{ID: id, Name: "Buddy"})
}
