package app

import (
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"auth-service/internal/utils"
	"os"
	"strconv"
	"time"
)

type Container struct {
	Logger     *logger.Log
	Handler    *handlers.Handler
	usersRepo  *postgres.UserRepository
	tokensRepo *postgres.RefreshTokenRepository
}

func NewContainer(logger *logger.Log) *Container {
	userDBURL := os.Getenv("USER_DB_URL")
	if userDBURL == "" {
		userDBURL = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
	}
	tokenDBURL := os.Getenv("TOKEN_DB_URL")
	if tokenDBURL == "" {
		tokenDBURL = userDBURL
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}
	accessTTL := time.Hour
	if ttl := os.Getenv("JWT_ACCESS_TTL"); ttl != "" {
		if parsed, err := strconv.Atoi(ttl); err == nil && parsed > 0 {
			accessTTL = time.Duration(parsed) * time.Second
		}
	}
	bcryptCost := 12
	if cost := os.Getenv("BCRYPT_COST"); cost != "" {
		if parsed, err := strconv.Atoi(cost); err == nil && parsed >= 4 {
			bcryptCost = parsed
		}
	}

	usersRepo, err := postgres.NewUserRepository(userDBURL)
	if err != nil {
		logger.Logger.Error("failed to initialize user repository", "error", err)
		os.Exit(1)
	}
	tokensRepo, err := postgres.NewRefreshTokenRepository(tokenDBURL)
	if err != nil {
		logger.Logger.Error("failed to initialize refresh token repository", "error", err)
		os.Exit(1)
	}

	hasher := utils.NewBcryptHasher(bcryptCost)
	jwtService := utils.NewJWTService(jwtSecret, accessTTL)
	svc := service.NewService(hasher, jwtService, tokensRepo, usersRepo)
	handler := handlers.NewHandler(logger, svc)

	return &Container{
		Logger:     logger,
		Handler:    handler,
		usersRepo:  usersRepo,
		tokensRepo: tokensRepo,
	}
}
