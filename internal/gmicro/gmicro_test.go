package gmicro_test

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/varrrro/pay-up/internal/gmicro"
	"github.com/varrrro/pay-up/internal/gmicro/group"
	"github.com/varrrro/pay-up/internal/gmicro/member"
)

var gm *gmicro.GroupsManager

func TestMain(m *testing.M) {
	// Open connection to test DB
	db, _ := gorm.Open("sqlite3", ":memory:")
	defer db.Close()

	// Create tables
	db.CreateTable(&group.Group{}, &member.Member{})

	// Create manager with test DB connection
	gm = gmicro.NewManager(db)

	os.Exit(m.Run())
}
