package gmicro_test

import (
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/varrrro/pay-up/internal/gmicro"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

var db *gorm.DB
var gm *gmicro.GroupsManager
var r *mux.Router

func TestMain(m *testing.M) {
	// Open connection to test DB
	db, _ = gorm.Open("sqlite3", ":memory:")
	defer db.Close()
	db.CreateTable(&group.Group{}, &member.Member{}) // create tables

	// Create manager with test DB connection
	gm = gmicro.NewManager(db)

	r = mux.NewRouter().StrictSlash(true)
	r.Use(gmicro.ContentTypeMiddleware)
	r.HandleFunc("/", gmicro.StatusHandler).Methods("GET")
	r.HandleFunc("/groups", gmicro.GroupsHandler(gm)).Methods("POST")
	r.HandleFunc("/groups/{groupid}", gmicro.GroupHandler(gm)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/groups/{groupid}/members", gmicro.MembersHandler(gm)).Methods("POST")
	r.HandleFunc("/groups/{groupid}/members/{memberid}", gmicro.MemberHandler(gm)).Methods("GET", "PUT", "DELETE")

	// Run tests
	os.Exit(m.Run())
}

func clearDB() {
	db.Delete(&member.Member{})
	db.Delete(&group.Group{})
}
