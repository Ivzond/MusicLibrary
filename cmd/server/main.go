package main

import (
	"MusicLibrary/internal/api"
	"MusicLibrary/internal/migrations"
	"MusicLibrary/pkg"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		pkg.Fatal("Ошибка загрузки файла .env", nil)
	}

	pkg.InitLogger() // Инициализация логгера

	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		pkg.Fatal("В конфигурационном файле не задана DATABASE_URL", nil)
	}
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		pkg.Fatal("Ошибка подключения к базу данных", map[string]interface{}{"error": err})
	}
	defer db.Close()

	err = migrations.ApplyMigrations(db)
	if err != nil {
		pkg.Fatal("Ошибка при применении миграций", map[string]interface{}{"error": err})
	}
	router := api.InitializeRoutes(db) // Используем новую функцию для инициализации роутера

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	pkg.Info("Сервис успешно запущен на порту: "+port, nil)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		pkg.Fatal("Ошибка при запуске HTTP-сервера", map[string]interface{}{"error": err})
	}
}
