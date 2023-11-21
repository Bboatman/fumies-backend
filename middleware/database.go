package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	"database/sql"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func DatabaseConnection() gin.HandlerFunc {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSL")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, db, sslMode)

	Migrate(connStr)

	pool, _ := Connect(&gin.Context{}, connStr)

	return func(c *gin.Context) {
		c.Set("db_pool", pool)
		c.Next()
	}
}

const MIGRATION_DIR = "/migrations"

func Connect(ctx *gin.Context, connStr string) (*pgxpool.Pool, string) {
	var err error
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
	return pool, connStr
}

func Migrate(connectionStr string) error {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Default().Printf("Unable to connect to db: %+v", err)
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Default().Printf("Unable to establish driver: %+v", err)
		return err
	}

	dbName := os.Getenv("DB_NAME")
	root, _ := os.Getwd()

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", root),
		dbName, driver)

	if err != nil {
		log.Default().Printf("Unable to migrate: %+v", err)
		return err
	}

	err = m.Up()
	if err != nil {
		log.Default().Printf("Migrations failed to apply, no change: %+v", err)
		return err
	}
	return nil
}
