package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/Vovarama1992/go-bagdoor-bot/docs"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/db"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/flight"
	httpauth "github.com/Vovarama1992/go-bagdoor-bot/internal/http/auth"
	httpflight "github.com/Vovarama1992/go-bagdoor-bot/internal/http/flight"
	httporder "github.com/Vovarama1992/go-bagdoor-bot/internal/http/order"
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

	userRepo := user.NewPostgresRepository(pool)
	userService := user.NewService(userRepo)

	orderRepo := order.NewPostgresRepository(pool)
	orderService := order.NewService(orderRepo, userRepo)

	flightRepo := flight.NewPostgresRepository(pool)
	flightService := flight.NewService(flightRepo)

	s3 := storage.NewS3Uploader()

	mux := http.NewServeMux()

	httpauth.RegisterRoutes(mux, httpauth.AuthDeps{
		UserService: userService,
	})

	httporder.RegisterRoutes(mux, httporder.OrderDeps{
		UserService:  userService,
		OrderService: orderService,
		Uploader:     s3,
	})

	httpflight.RegisterRoutes(mux, httpflight.FlightDeps{
		UserService:   userService,
		FlightService: flightService,
	})

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("HTTP-сервер запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
