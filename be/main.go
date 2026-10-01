package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"vocabulary/database"
	"vocabulary/route"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := runMigration(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		return
	}

	var err error
	db, err = database.Connect()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	router := gin.Default()
	route.RegisterAPIRoutes(router, db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

func runMigration(args []string) error {
	m, err := migrate.New("file://migrations", database.PostgresMigrationURL())
	if err != nil {
		return err
	}

	if len(args) == 0 || args[0] == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}

		fmt.Println("migration up completed")
		return nil
	}

	if args[0] == "down" {
		if len(args) < 2 {
			return errors.New("VERSION is required. Usage: make migrate-down VERSION=2")
		}

		version, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid VERSION %q: %w", args[1], err)
		}

		if err := m.Migrate(uint(version)); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}

		fmt.Printf("migration changed to version %d\n", version)
		return nil
	}

	return fmt.Errorf("unknown migrate command %q", args[0])
}
