package main

import (
	"log"
	"os"

	"github.com/kritAsawaniramol/book-store/config"
	"github.com/kritAsawaniramol/book-store/module/auth"
	"github.com/kritAsawaniramol/book-store/module/book"
	"github.com/kritAsawaniramol/book-store/module/order"
	"github.com/kritAsawaniramol/book-store/module/shelf"
	"github.com/kritAsawaniramol/book-store/module/user"
	"github.com/kritAsawaniramol/book-store/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	db := database.NewPostgresDatabase(cfg)

	migration(db, cfg)
}

func migration(db database.Database, cfg *config.Config) {
	switch cfg.App.Name {
	case "auth":
		authMigration(db)
	case "user":
		userMigration(db, cfg)
	case "shelf":
		shelfMigration(db)
	case "book":
		bookMigration(db)
	case "order":
		orderMigration(db)
	}
}

func authMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(
		&auth.Credential{},
	)

	if err != nil {
		panic(err)
	}

	log.Println("Auth database migration completed!")
}

func userMigration(db database.Database, cfg *config.Config) {
	err := db.GetDb().AutoMigrate(
		&user.User{},
		&user.Role{},
		&user.UserTransactions{},
		&user.TopUpOrder{},
	)
	if err != nil {
		panic(err)
	}

	// roles := []user.Role{{RoleTitle: "admin"}, {RoleTitle: "customer"}}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	// users := []user.User{{Username: cfg.Admin.Username, Password: string(hashedPassword), RoleID: 1}}
	
	
	if err := db.GetDb().FirstOrCreate(&user.Role{}, &user.Role{RoleTitle: "admin"}).Error; err != nil {
		panic(err)
	}
	if err := db.GetDb().FirstOrCreate(&user.Role{}, &user.Role{RoleTitle: "customer"}).Error; err != nil {
		panic(err)
	}


	// if err := db.GetDb().CreateInBatches(roles, 2).Error; err != nil {
	// 	panic(err)
	// }

	if err := db.GetDb().FirstOrCreate(&user.User{}, &user.User{Username: cfg.Admin.Username, Password: string(hashedPassword), RoleID: 1}).Error; err != nil {
		panic(err)
	}

	// if err := db.GetDb().Create(users).Error; err != nil {
	// 	panic(err)
	// }

	log.Println("User database migration completed!")
}

func shelfMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(
		&shelf.Shelves{},
	)

	if err != nil {
		panic(err)
	}

	log.Println("Shelf database migration completed!")
}

func bookMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(
		&book.Books{},
		&book.Tags{},
	)

	if err != nil {
		panic(err)
	}

	log.Println("Book database migration  completed!")
}

func orderMigration(db database.Database) {
	err := db.GetDb().AutoMigrate(
		&order.Orders{},
		&order.OrdersBooks{},
	)


	if err != nil {
		panic(err)
	}

	log.Println("Order database migration completed!")
}
