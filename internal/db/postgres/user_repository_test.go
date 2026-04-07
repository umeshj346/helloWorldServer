package postgres

import (
	"context"
	"errors"
	"net/mail"
	"os"
	"reflect"
	"testing"

	"github.com/umeshj346/helloWorldServer/domain"
	"github.com/umeshj346/helloWorldServer/utils"
)


func Test_Fetch(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}

	testDb := NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()
	
	query := `
		INSERT INTO users (first_name, last_name, email)
		VALUES ($1, $2, $3)
	`
	testFirstName1, testLastName1, testEmail1 := "foo", "bar", "foo@bar.com"
	testFirstName2, testLastName2, testEmail2 := "foo", "baz", "foo@baz.com"
	testDb.Exec(query, testFirstName1, testLastName1, testEmail1)
	testDb.Exec(query, testFirstName2, testLastName2, testEmail2)
	
	testRepo := NewUserPgRepository(testDb)
	
	query = `
		SELECT id, first_name, last_name, email 
		from users
	`
	resultList, err := testRepo.fetch(context.TODO(), query)
	if err != nil {
		t.Fatalf("error fetching users, err: %v", err)
	}
	
	if len(resultList) != 2 {
		t.Fatalf("invalid no. of users present in the database, err: %v", err)
	}
	resultList[0].ID = 0
	resultList[1].ID = 0
	email1, err := mail.ParseAddress(testEmail1)
	email2, err := mail.ParseAddress(testEmail2)

	expectedList := []domain.User {
		{
			FirstName: testFirstName1,
			LastName: testLastName1,
			Email: *email1,
		},
		{
			FirstName: testFirstName2,
			LastName: testLastName2,
			Email: *email2,
		},
	}

	if !reflect.DeepEqual(resultList, expectedList) {
		t.Errorf("added user data is not correct\nexp: %+v\ngot: %+v\n", expectedList, resultList)
	}
}

func Test_Insert(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}

	testDb := NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()
	
	
	user := domain.UserData{
		FirstName: "foo",
		LastName: "bar",
		Email: "foo.bar@eg.com",
	}
	
	testRepo := NewUserPgRepository(testDb)
	err = testRepo.InsertUser(context.TODO(), &user)
	if err != nil {
		t.Fatalf("error fetching users, err: %v", err)
	}

	query := `
		SELECT id, first_name, last_name, email 
		from users
	`
	resultList, err := testRepo.fetch(context.TODO(), query)
	if len(resultList) != 1 {
		t.Fatalf("invalid no. of users present in the database, wanted: 1, got: %v", len(resultList))
	}
	resultList[0].ID = 0
	email, err := mail.ParseAddress(user.Email)

	expectedList := []domain.User {
		{
			FirstName: user.FirstName,
			LastName: user.LastName,
			Email: *email,
		},
	}

	if !reflect.DeepEqual(resultList, expectedList) {
		t.Errorf("added user data is not correct\nexp: %+v\ngot: %+v\n", expectedList, resultList)
	}
}

func Test_GetUserByName(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testRepo := NewUserPgRepository(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	user := domain.UserData{
		FirstName: "foo",
		LastName: "bar",
		Email: "foo.bar@g.com",
	}
	
	err = testRepo.InsertUser(context.TODO(), &user)
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}
	testEmail, err := mail.ParseAddress(user.Email)
	if err != nil {
		t.Fatalf("error parsing email address, err: %v", err)
	}
	testUser := domain.User{
		FirstName: "foo",
		LastName: "bar",
		Email: *testEmail,
	}

	err = testRepo.InsertUser(context.TODO(), 
		&domain.UserData{
			FirstName: "foo",
			LastName: "baz",
			Email: "foo.baz@g.com",
		},
	)

	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	err = testRepo.InsertUser(context.TODO(), 
		&domain.UserData{
			FirstName: "bar",
			LastName: "baz",
			Email: "bar.baz@g.com",
		},
	)
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	err = testRepo.InsertUser(context.TODO(), 
		&domain.UserData{
			FirstName: "baz",
			LastName: "foo",
			Email: "baz.foo@g.com",
		},
	)
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	tests := map[string]struct {
		first 		string
		last 		string
		expected 	*domain.User
		expectedErr error
	} {
		"simple": {
			first: 	"foo",
			last: 	"bar",
			expected: &testUser,
			expectedErr: nil,
		},
		"first name lookup": {
			first: "foo",
			last: "boo",
			expected: nil,
			expectedErr: domain.ErrNoResultFound,
		},
		"last name lookup": {
			first: "boo",
			last: "foo",
			expected: nil,
			expectedErr: domain.ErrNoResultFound,
		},
		"no match lookup": {
			first: "boo",
			last: "boo",
			expected: nil,
			expectedErr: domain.ErrNoResultFound,
		},
		"empty first name": {
			first: "",
			last: "baz",
			expected: nil,
			expectedErr: domain.ErrNoResultFound,
		},
		"empty last name": {
			first: "foo",
			last: "",
			expected: nil,
			expectedErr: domain.ErrNoResultFound,
		},
	}
	for name, test := range tests {
		result, err := testRepo.GetUserByName(context.TODO(), test.first, test.last)
		if result != nil {
			result.ID = 0
		}
		if !reflect.DeepEqual(result, test.expected){
			t.Errorf("%s: invalid result\nexpected: %+v\ngot: %+v\n", 
						name, test.expected, result)
		}

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: invalid error reported\nexpected: %v\ngot: %v", name, test.expectedErr, err)
		}
		
	}
}

func Test_CountOfUsers(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}

	testDb := NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testRepo := NewUserPgRepository(testDb)
	
	mockUsers := []domain.UserData{
		{
			FirstName: "Test",
			LastName: "Man-1",
			Email: "test.man@1.com",
		},
		{
			FirstName: "Test",
			LastName: "Man-2",
			Email: "test.man@3.com",
		},
	}
	
	err = testRepo.InsertUser(context.TODO(), &mockUsers[0])
	if err != nil {
		t.Fatalf("error inserting into database, err: %v", err)
	}

	err = testRepo.InsertUser(context.TODO(), &mockUsers[1])
	if err != nil {
		t.Fatalf("error inserting into database, err: %v", err)
	}

	usersCnt, err := testRepo.GetCountOfUsers(context.TODO())
	if err != nil {
		t.Fatalf("error fetching usersCnt, err: %v", err)
	}
	if usersCnt != 2 {
		t.Errorf("bad users cnt, exp: 2, got: %v", usersCnt)
	}
}
