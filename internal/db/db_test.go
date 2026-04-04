package db

import (
	"os"
	"testing"

	"github.com/umeshj346/helloWorldServer/utils"
)

func Test_NewPostgresDB_UserDB(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}

	testDBStr := os.Getenv("DATABASE_URL")
	db := NewPostgresDB(testDBStr)
	if db == nil {
		t.Fatalf("error accessing test database")
	}
}

func Test_NewPostgresDB_TestUserDB(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	
	testDBStr := os.Getenv("TEST_DATABASE_URL")
	db := NewPostgresDB(testDBStr)
	if db == nil {
		t.Fatalf("error accessing test database")
	}
}