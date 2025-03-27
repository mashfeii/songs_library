package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mashfeii/songs_library/config"
	client "github.com/mashfeii/songs_library/internal/api"
	"github.com/mashfeii/songs_library/internal/application"
	"github.com/mashfeii/songs_library/internal/infrastructure/database"
	"github.com/mashfeii/songs_library/internal/infrastructure/handlers"
	"github.com/sirupsen/logrus"

	_ "github.com/mashfeii/songs_library/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func migrateDB(connectionString string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	m, err := migrate.New(
		"file://sql/migrations",
		connectionString,
	)
	if err != nil {
		logrus.Fatal("could not create migration: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatal("could not apply migration: ", err)
	}

	logrus.Info("Database migrated successfully")
}

func initRouting(r *gin.Engine, service *application.SongsService) {
	r.GET("/songs", handlers.GetSongs(service))
	r.GET("/songs/:id/verses", handlers.GetSongVerses(service))
	r.POST("/songs", handlers.AddSong(service))
	r.PUT("/songs/:id", handlers.UpdateSong(service))
	r.DELETE("/songs/:id", handlers.DeleteSong(service))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func main() {
	config, err := config.NewConfigFromFile("app")
	if err != nil {
		logrus.Fatal("Loading .env file: ", err)
	}

	migrateDB(config.ToDSN())

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.ToDSN())
	if err != nil {
		logrus.Fatal("Creating connection pool: ", err)
	}
	defer pool.Close()

	logrus.Info("Database connection pool created")

	externalClient, err := client.NewClientWithResponses(config.APIEndpoint)
	if err != nil {
		logrus.Fatal("Creating external client: ", err)
	}

	repository := database.NewPgxRepository(pool)
	service := application.NewSongsService(
		repository,
		externalClient,
	)

	r := gin.Default()
	initRouting(r, service)

	logrus.Info("Starting server on port ", config.ServingPort)

	logrus.Error(r.Run(fmt.Sprintf(":%d", config.ServingPort)))
}
