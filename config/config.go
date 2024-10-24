package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App    App
		Client Client
		Db     Db
		Kafka  Kafka
		Grpc   Grpc
		Jwt    Jwt
		Admin  Admin
		Stripe Stripe
	}

	App struct {
		Stage string
		Name  string
		Host  string
		Port  int
	}

	Client struct {
		URL string
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
		Url     string
		GroupID string
		ApiKey  string
		Secret  string
	}

	Grpc struct {
		AuthUrl  string
		UserUrl  string
		ShelfUrl string
		BookUrl  string
		OrderUrl string
	}

	Jwt struct {
		AccessSecretKey  string
		AccessDuration   int64
		RefreshSecretKey string
		RefreshDuration  int64
	}

	Admin struct {
		Username string
		Password string
	}

	Stripe struct {
		SecretKey      string
		EndPointSecret string
	}
)

var config *Config

var lock = &sync.Mutex{}

func LoadConfig(path string) *Config {
	if config == nil {
		lock.Lock()
		defer lock.Unlock()
		if config == nil {
			if err := godotenv.Load(path); err != nil {

				log.Fatalf("Error loading .env file: %s", err.Error())
			}
			appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
			if err != nil {
				log.Fatal("Error loading .env file: app's port is invalid")
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

			config = &Config{
				App: App{
					Stage: os.Getenv("APP_STAGE"),
					Name:  os.Getenv("APP_NAME"),
					Host:  os.Getenv("APP_HOST"),
					Port:  appPort,
				},

				Client: Client{
					URL:  os.Getenv("CLIENT_URL"),
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
					Url:     os.Getenv("KAFKA_URL"),
					GroupID: os.Getenv("KAFKA_GROUP_ID"),
					ApiKey:  os.Getenv("KAFKA_API_KEY"),
					Secret:  os.Getenv("KAFKA_API_SECRET"),
				},
				Grpc: Grpc{
					AuthUrl:  os.Getenv("GRPC_AUTH_URL"),
					UserUrl:  os.Getenv("GRPC_USER_URL"),
					ShelfUrl: os.Getenv("GRPC_SHELF_URL"),
					BookUrl:  os.Getenv("GRPC_BOOK_URL"),
					OrderUrl: os.Getenv("GRPC_ORDER_URL"),
				},

				Jwt: Jwt{
					AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
					AccessDuration:   jwtAccessDuration,
					RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
					RefreshDuration:  jwtRefreshDuration,
				},

				Admin: Admin{
					Username: os.Getenv("ADMIN_USERNAME"),
					Password: os.Getenv("ADMIN_PASSWORD"),
				},

				Stripe: Stripe{
					SecretKey:      os.Getenv("STRIPE_SECRET_KEY"),
					EndPointSecret: os.Getenv("STRIPE_ENDPOINT_SECRET"),
				},
			}
		} else {
			log.Println("Config already loaded.")
		}
	} else {
		log.Println("Config already loaded.")
	}
	return config
}
