package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/varrrro/pay-up/internal/gmicro"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

func main() {
	// Open database connection
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Create tables
	db.CreateTable(&group.Group{}, &member.Member{})

	// Create data manager using database connection
	gm := gmicro.NewManager(db)

	// Build handler functions using data manager
	h := gmicro.NewHandlers(gm)

	// Build router with handlers
	r := gmicro.NewRouter(h)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
