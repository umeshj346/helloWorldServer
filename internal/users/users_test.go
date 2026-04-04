package users

import (
	"errors"
	"net/mail"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/umeshj346/helloWorldServer/internal/db"
	"github.com/umeshj346/helloWorldServer/utils"
)

func Test_AddUser(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testFirstName, testLastName:= "Test", "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing the address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.Address)
	if err != nil  {
		t.Fatalf("error creating user: %v", err)
	}

	sizeQuery := `
		SELECT COUNT(*) FROM users
	`
	row := testManager.db.QueryRow(sizeQuery)
	var tbLen int
	err = row.Scan(&tbLen)
	if err != nil {
		t.Fatalf("error reading no. of entries in users table: %v", err)
	}
	if tbLen != 1 {
		t.Fatalf("bad test Manager users count, wanted: 1, got: %v", tbLen)
	}

	query := `
		SELECT first_name, last_name, email
		FROM users
		WHERE first_name = $1 and last_name = $2 and email = $3
	`
	var foundUser User
	var email string

	row = testManager.db.QueryRow(query, testFirstName, testLastName, testEmail.Address)
	err = row.Scan(&foundUser.FirstName, &foundUser.LastName, &email)
	if err != nil {
		t.Fatalf("error extracting data from row: %v", err)
	}

	parsedResultEmail, err := mail.ParseAddress(email)
	if err != nil {
		t.Fatalf("error parsing the email: %v", err)
	}
	foundUser.Email = *parsedResultEmail

	expUser := User{
		FirstName: testFirstName, 
		LastName: testLastName, 
		Email:  *testEmail,
	}

	if !reflect.DeepEqual(expUser, foundUser) {
		t.Errorf("added user data is not correct\nexp: %+v\ngot: %+v", expUser, foundUser)
	}
}

func Test_Adduser_InvalidEmail(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testFirstName, testLastName:= "Test", "Userman"
	testEmail := "foobar"

	err = testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil  {
		t.Fatalf("no error for invalid email(%v)", testEmail)
	} else {
		expectedError := "invalid email: foobar"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}

	sizeQuery := `
		SELECT COUNT(*) FROM users
	`
	row := testManager.db.QueryRow(sizeQuery)
	var tbLen int
	err = row.Scan(&tbLen)
	if err != nil {
		t.Fatalf("error reading no. of entries in users table: %v", err)
	}
	if tbLen != 0 {
		t.Fatalf("bad test Manager users count, wanted: 1, got: %v", tbLen)
	}
}

func Test_AddUser_MultipleCalls(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testFirstName, testLastName, testEmail := "foo", "bar", "foo@bar"

	numRequests := 10000
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			_ = testManager.AddUser(testFirstName, testLastName, testEmail)
		}()
	}

	wg.Wait()
	sizeQuery := `
		SELECT COUNT(*) FROM users
	`
	row := testManager.db.QueryRow(sizeQuery)
	var tbLen int
	err = row.Scan(&tbLen)
	if err != nil {
		t.Fatalf("error reading no. of entries in users table: %v", err)
	}
	if tbLen != 1 {
		t.Fatalf("bad test Manager users count, wanted: 1, got: %v", tbLen)
	}	
}

func Test_AddUser_EmptyFirstName(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testFirstName, testLastName:= "", "Userman"
	testEmail := "foo@bar.com"

	err = testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil  {
		t.Fatalf("no error for empty first name")
	} else {
		expectedError := "empty first name"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}

	sizeQuery := `
		SELECT COUNT(*) FROM users
	`
	row := testManager.db.QueryRow(sizeQuery)
	var tbLen int
	err = row.Scan(&tbLen)
	if err != nil {
		t.Fatalf("error reading no. of entries in users table: %v", err)
	}
	if tbLen != 0 {
		t.Fatalf("bad test Manager users count, wanted: 1, got: %v", tbLen)
	}
}

func Test_AddUser_EmptyLastName(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testFirstName, testLastName:= "Test", ""
	testEmail := "foo@bar.com"

	err = testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil  {
		t.Fatalf("no error for empty last name")
	} else {
		expectedError := "empty last name"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}

	sizeQuery := `
		SELECT COUNT(*) FROM users
	`
	row := testManager.db.QueryRow(sizeQuery)
	var tbLen int
	err = row.Scan(&tbLen)
	if err != nil {
		t.Fatalf("error reading no. of entries in users table: %v", err)
	}
	if tbLen != 0 {
		t.Fatalf("bad test Manager users count, wanted: 1, got: %v", tbLen)
	}
}

func Test_AddUser_DuplicateUser(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()

	testFirstName, testLastName:= "Test", "UserMan"
	testEmail := "foo@bar.com"

	err = testManager.AddUser(testFirstName, testLastName, testEmail)
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail)
	if errors.Is(err, nil)  {
		t.Fatalf("no error for duplicate user")
	} else {
		expectedError := "user with this name already exists"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err)
		}
	}

	sizeQuery := `
		SELECT COUNT(*) FROM users
	`
	row := testManager.db.QueryRow(sizeQuery)
	var tbLen int
	err = row.Scan(&tbLen)
	if err != nil {
		t.Fatalf("error reading no. of entries in users table: %v", err)
	}
	if tbLen != 1 {
		t.Fatalf("bad test Manager users count, wanted: 1, got: %v", tbLen)
	}
}

func Test_GetUserByName(t *testing.T) {
	err := utils.LoadEnv()
	if err != nil {
		t.Fatalf("error opening .env file, err: %v", err)
	}
	testDb := db.NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	testManager := NewManager(testDb)
	defer func() {
		testDb.Exec(`DELETE FROM users`)
	}()
	
	err = testManager.AddUser("foo", "bar", "foo.bar@g.com")
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}
	testEmail, err := mail.ParseAddress("foo.bar@g.com")
	if err != nil {
		t.Fatalf("error parsing email address, err: %v", err)
	}
	testUser := User{
		FirstName: "foo",
		LastName: "bar",
		Email: *testEmail,
	}

	err = testManager.AddUser("foo", "baz", "foo.baz@g.com")
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	err = testManager.AddUser("bar", "baz", "bar.baz@g.com")
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	err = testManager.AddUser("baz", "foo", "baz.foo@g.com")
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	tests := map[string]struct {
		first 		string
		last 		string
		expected 	*User
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
			expectedErr: ErrNoResultFound,
		},
		"last name lookup": {
			first: "boo",
			last: "foo",
			expected: nil,
			expectedErr: ErrNoResultFound,
		},
		"no match lookup": {
			first: "boo",
			last: "boo",
			expected: nil,
			expectedErr: ErrNoResultFound,
		},
		"empty first name": {
			first: "",
			last: "baz",
			expected: nil,
			expectedErr: ErrNoResultFound,
		},
		"empty last name": {
			first: "foo",
			last: "",
			expected: nil,
			expectedErr: ErrNoResultFound,
		},
	}
	for name, test := range tests {
		result, err := testManager.GetUserByName(test.first, test.last)
		if !reflect.DeepEqual(result, test.expected){
			t.Errorf("%s: invalid result\nexpected: %+v\ngot: %+v\n", 
						name, test.expected, result)
		}

		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: invalid error reported\nexpected: %v\ngot: %v", name, test.expectedErr, err)
		}
		
	}
}


