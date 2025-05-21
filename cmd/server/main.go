package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/Vovarama1992/go-bagdoor-bot/docs"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/flight"
	httpPkg "github.com/Vovarama1992/go-bagdoor-bot/internal/http"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/user"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}

	pool, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Репозитории и сервисы
	userRepo := user.NewPostgresRepository(pool)
	userService := user.NewService(userRepo)

	orderRepo := order.NewPostgresRepository(pool)
	orderService := order.NewService(orderRepo, userRepo)

	flightRepo := flight.NewPostgresRepository(pool)
	flightService := flight.NewService(flightRepo)

	s3 := storage.NewS3Uploader()

	// Роутинг
	mux := http.NewServeMux()

	orderDeps := httpPkg.OrderDeps{
		UserService:  userService,
		OrderService: orderService,
		Uploader:     s3,
	}

	flightDeps := httpPkg.FlightDeps{
		UserService:   userService,
		FlightService: flightService,
	}

	httpPkg.RegisterRoutes(mux, orderDeps)
	httpPkg.RegisterFlightRoutes(mux, flightDeps)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("HTTP-сервер запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
