package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/kritAsawaniramol/book-store/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var dbConn *postgresDatabase

var lock = &sync.Mutex{}

// Implement Singleton
func NewPostgresDatabase(cfg *config.Config) Database {
	if dbConn == nil {
		lock.Lock()
		defer lock.Unlock()
		if dbConn == nil {
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
				cfg.Db.Host,
				cfg.Db.User,
				cfg.Db.Password,
				cfg.Db.DBName,
				cfg.Db.Port,
				cfg.Db.SSLMode,
				cfg.Db.TimeZone,
			)
			DBLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: time.Second, // Slow SQL threshold
					LogLevel:      logger.Info, // Log level
					Colorful:      true,        // Disable color
				},
			)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger:      DBLogger,
				PrepareStmt: false,
			})
			if err != nil {
				panic("failed to connect database")
			}

			dbConn = &postgresDatabase{Db: db}

		} else {
			log.Println("Database connection already created.")
		}
	} else {
		log.Println("Database connection already created.")
	}
	return dbConn
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return p.Db
}
