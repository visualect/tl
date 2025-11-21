package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/visualect/tl/internal/auth"
	"github.com/visualect/tl/internal/handlers"
	"github.com/visualect/tl/internal/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	tasksRepo := repo.NewTasks(db)
	usersRepo := repo.NewUsers(db)

	h := handlers.New(usersRepo, tasksRepo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &handlers.CustomValidator{Validator: validator.New()}

	e.POST("/signup", h.SignUp)
	e.POST("/login", h.Login)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
		SigningKey: auth.Secret,
	}

	authRequired := e.Group("")
	authRequired.Use(echojwt.WithConfig(config))

	authRequired.GET("/tasks", h.GetTasks)
	authRequired.POST("/tasks", h.AddTask)
	authRequired.PATCH("/tasks/:id", h.ToggleCompleteTask)
	authRequired.DELETE("/tasks/:id", h.DeleteTask)

	e.Logger.Fatal(e.Start(":8000"))
}
