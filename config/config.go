package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      App
		Client   Client
		Db       Db
		Kafka    Kafka
		Grpc     Grpc
		Jwt      Jwt
		Sessions Sessions
		Google   Google
		Facebook Facebook
	}

	App struct {
		Stage string
		Name  string
		Host  string
		Port  int
	}

	Client struct {
		Host string
		Port int
	}

	// Database
	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}

	Kafka struct {
		Url    string
		ApiKey string
		Secret string
	}

	Grpc struct {
		AuthUrl    string
		UserUrl    string
		ShelfUrl   string
		BookUrl    string
		PaymentUrl string
	}

	Jwt struct {
		AccessSecretKey  string
		AccessDuration   int64
		RefreshSecretKey string
		RefreshDuration  int64
	}

	Sessions struct {
		Secret string
		MaxAge int
	}

	Google struct {
		ClientID     string
		ClientSecret string
	}

	Facebook struct {
		ClientID     string
		ClientSecret string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file")
	}
	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal("Error loading .env file: app's port is invalid")
	}
	clientPort, err := strconv.Atoi(os.Getenv("CLIENT_PORT"))
	if err != nil {
		log.Fatal("Error loading .env file: client's port is invalid")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Error loading .env file: db's port is invalid")
	}

	jwtAccessDuration, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_DURATION"), 10, 64)
	if err != nil {
		log.Fatal(`Error loading .env file: db's "jwt access duration" is invalid`)
	}

	jwtRefreshDuration, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_DURATION"), 10, 64)
	if err != nil {
		log.Fatal(`Error loading .env file: db's "jwt refresh duration" is invalid`)
	}

	sessionsMaxAge, err := strconv.Atoi(os.Getenv("SESSIONS_MAX_AGE"))
	if err != nil {
		log.Fatal(`Error loading .env file: db's "sessions max age" is invalid`)
	}

	return Config{
		App: App{
			Stage: os.Getenv("APP_STAGE"),
			Name:  os.Getenv("APP_NAME"),
			Host:  os.Getenv("APP_HOST"),
			Port:  appPort,
		},

		Client: Client{
			Host: os.Getenv("CLIENT_HOST"),
			Port: clientPort,
		},

		Db: Db{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			TimeZone: os.Getenv("DB_TIMEZONE"),
		},

		Kafka: Kafka{
			Url:    os.Getenv("KAFKA_URL"),
			ApiKey: os.Getenv("KAFKA_API_KEY"),
			Secret: os.Getenv("KAFKA_SECRET"),
		},
		Grpc: Grpc{
			AuthUrl:    os.Getenv("GRPC_AUTH_URL"),
			UserUrl:    os.Getenv("GRPC_USER_URL"),
			ShelfUrl:   os.Getenv("GRPC_SHELF_URL"),
			BookUrl:    os.Getenv("GRPC_BOOK_URL"),
			PaymentUrl: os.Getenv("GRPC_PAYMENT_URL"),
		},

		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			AccessDuration:   jwtAccessDuration,
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			RefreshDuration:  jwtRefreshDuration,
		},

		Sessions: Sessions{
			Secret: os.Getenv("SESSIONS_SECRET"),
			MaxAge: sessionsMaxAge,
		},

		Google: Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},

		Facebook: Facebook{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		},
	}
}
