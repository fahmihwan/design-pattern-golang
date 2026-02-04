package main

import (
	"best-pattern/internal/config"
	"best-pattern/internal/handler"
	"best-pattern/internal/repository"
	"best-pattern/internal/service"
	"best-pattern/internal/util"
	"best-pattern/pkg/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func main() {

	cfg, err := config.LoadConfig()

	// // // Create a new DI container
	container := dig.New()

	err = container.Provide(func() *config.Config { return cfg })
	if err != nil {
		panic(fmt.Sprintf("Failed to provide configuration: %v", err))
	}

	// Provide PostgreSQL database
	err = container.Provide(database.NewPostgresDB)
	if err != nil {
		panic(fmt.Sprintf("Failed to provide PostgreSQL database: %v", err))
	}
	_ = container.Invoke(func(db *gorm.DB) {
		// middleware.InitializeValidator(db)
	})

	//Provide JWT Manager
	err = container.Provide(func(cfg *config.Config) *util.JWTManager {
		exp := time.Duration(cfg.JWT.ExpiryMin) * time.Minute
		return util.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Issuer, exp)
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to provide JWTManager: %v", err))
	}
	// ===================================================
	// Provide router (chi.Mux implements http.Handler)
	provideRepositories(container)

	provideServices(container)

	provideHandler(container)

	//  provide router
	if err := container.Provide(handler.NewRouter); err != nil {
		panic(fmt.Sprintf("Failed to provide router: %v", err))
	}

	// server: minta *chi.Mux
	if err := container.Provide(func(r *chi.Mux) *http.Server {
		return &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
	}); err != nil {
		panic(fmt.Sprintf("Failed to provide server: %v", err))
	}

	if err := container.Invoke(func(srv *http.Server) {
		log.Println("listening on http://localhost:8080")
		log.Fatal(srv.ListenAndServe())
	}); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
func provideRepositories(container *dig.Container) {
	var err error

	err = container.Provide(func(db *gorm.DB) repository.BookRepository {
		return *repository.NewBookRepository(db)
	})

	err = container.Provide(func(db *gorm.DB) repository.UserRepository {
		return *repository.NewUserRepository(db)
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to provide BookRepository: %v", err))
	}

	err = container.Provide(func(db *gorm.DB) repository.Repository {
		return repository.Repository{
			Book: repository.NewBookRepository(db),
			User: repository.NewUserRepository(db),
		}
	})
}

func provideServices(container *dig.Container) {
	var err error

	err = container.Provide(func(db *gorm.DB, repo repository.Repository) *service.BookService {
		return service.NewBookService(repo)
	})

	err = container.Provide(func(db *gorm.DB, repo repository.Repository) *service.UserService {
		return service.NewUserService(repo)
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to provide BookService: %v", err))
	}

	err = container.Provide(func(db *gorm.DB, repo repository.Repository) service.Service {
		return service.Service{
			Book: service.NewBookService(repo),
			User: service.NewUserService(repo),
		}
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to provide Sector: %v", err))
	}
}

func provideHandler(container *dig.Container) {
	// 1) provide semua sub-handler (nanti tinggal tambah baris di sini)
	if err := container.Provide(handler.NewUserHandler); err != nil {
		panic(fmt.Sprintf("Failed to provide UserHandler: %v", err))
	}
	if err := container.Provide(handler.NewBookHandler); err != nil {
		panic(fmt.Sprintf("Failed to provide UserHandler: %v", err))
	}

	// 2) provide agregator HandlerInteface (isi field-fieldnya dari DI)
	if err := container.Provide(func(
		BookService *service.BookService,
		UserService *service.UserService,
		jwtManager *util.JWTManager,
	) *handler.HandlerInteface {
		return &handler.HandlerInteface{
			UserHandler: handler.NewUserHandler(UserService, jwtManager),
			BookHandler: handler.NewBookHandler(BookService),
		}
	}); err != nil {
		panic(fmt.Sprintf("Failed to provide HandlerInteface: %v", err))
	}
}
