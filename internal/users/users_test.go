package users

import (
	"errors"
	"net/mail"
	"reflect"
	"sync"
	"testing"
)

func Test_AddUser(t *testing.T) {
	testManager := NewManager()

	testFirstName, testLastName:= "Test", "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing the address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.Address)
	if err != nil  {
		t.Fatalf("error creating user: %v", err)
	}

	if len(testManager.users) != 1 {
		t.Errorf("bad test manager user count, expected: %d, got: %d", 1, len(testManager.users))
		if len(testManager.users) < 1 {
			t.Fatalf("users is empty")
		}
	}
	expUser := User{
		FirstName: testFirstName, 
		LastName: testLastName, 
		Email:  *testEmail,
	}
	foundUser := testManager.users[0]

	if !reflect.DeepEqual(expUser, foundUser) {
		t.Errorf("added user data is not correct\nexp: %+v\ngot: %+v", expUser, foundUser)
	}
}

func Test_Adduser_InvalidEmail(t *testing.T) {
	testManager := NewManager()

	testFirstName, testLastName:= "Test", "Userman"
	testEmail := "foobar"

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil  {
		t.Fatalf("no error for invalid email(%v)", testEmail)
	} else {
		expectedError := "invalid email: foobar"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}

	if len(testManager.users) > 0 {
		t.Errorf("bad test manager user count, exp: %d, got : %d", 0, len(testManager.users))
	}
}

func Test_AddUser_MultipleCalls(t *testing.T) {
	testManager := NewManager()

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
	if len(testManager.users) != 1 {
		t.Errorf("bad test manager user count, exp: %d, got : %d", 0, len(testManager.users))
	}

	
}

func Test_AddUser_EmptyFirstName(t *testing.T) {
	testManager := NewManager()

	testFirstName, testLastName:= "", "Userman"
	testEmail := "foo@bar.com"

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil  {
		t.Fatalf("no error for empty first name")
	} else {
		expectedError := "empty first name"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}

	if len(testManager.users) > 0 {
		t.Errorf("bad test manager user count, exp: %d, got : %d", 0, len(testManager.users))
	}
}

func Test_AddUser_EmptyLastName(t *testing.T) {
	testManager := NewManager()

	testFirstName, testLastName:= "Test", ""
	testEmail := "foo@bar.com"

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil  {
		t.Fatalf("no error for empty last name")
	} else {
		expectedError := "empty last name"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}

	if len(testManager.users) > 0 {
		t.Errorf("bad test manager user count, exp: %d, got : %d", 0, len(testManager.users))
	}
}

func Test_AddUser_DuplicateUser(t *testing.T) {
	testManager := NewManager()
	
	testFirstName, testLastName:= "Test", "UserMan"
	testEmail := "foo@bar.com"

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
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

	if len(testManager.users) > 1 {
		t.Errorf("bad test manager user count, exp: %d, got : %d", 0, len(testManager.users))
	}
}

func Test_GetUserByName(t *testing.T) {
	testManager := NewManager()
	err := testManager.AddUser("foo", "bar", "foo.bar@g.com")
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
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
			expected: &testManager.users[0],
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
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("%s: invalid result\nexpected: %+v\ngot: %+v\n", 
						name, test.expected, result)
		}
		
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: invalid error reported\nexpected: %v\ngot: %v", name, test.expectedErr, err)
		}

		
	}
}


