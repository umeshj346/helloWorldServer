package users

import (
	"context"
	"errors"
	"net/mail"
	"reflect"
	"testing"

	"github.com/umeshj346/helloWorldServer/domain"
	"github.com/umeshj346/helloWorldServer/users/mocks"
)

var ctx context.Context =context.TODO()

func Test_AddUser(t *testing.T) {
	repo := mocks.MockRepo{}
	testService := Service{repo: &repo}
	testUserData := domain.UserData {
		FirstName: "foo",
		LastName: "bar",
		Email: "foo.bar@eg.com",
	}
	err := testService.AddUser(ctx, &testUserData)
	if err != nil {
		t.Fatalf("error adding User, err: %v", err)
	}
	if cnt, _ := repo.GetCountOfUsers(ctx); cnt != 1 {
		t.Fatalf("bad user count, exp: %v, got: %v", 1, cnt)
	}
	email, err := mail.ParseAddress(testUserData.Email)
	if err != nil {
		t.Fatalf("error parsing email, err: %v", err)
	}
	testUser := domain.User {
		FirstName: testUserData.FirstName,
		LastName: testUserData.LastName,
		Email: *email,
	}
	foundUser, err := repo.GetUserByName(ctx, testUser.FirstName, testUser.LastName)
	if err != nil {
		t.Fatalf("user is not found!!")
	}
	if !reflect.DeepEqual(foundUser, testUser) {
		t.Errorf("bad user added in the database, err: %v", err)
	}
}

func Test_Adduser_InvalidEmail(t *testing.T) {
	repo := mocks.MockRepo{}
	testService := Service{repo: &repo}
	testUserData := domain.UserData {
		FirstName: "foo",
		LastName: "bar",
		Email: "foo.b",
	}
	err := testService.AddUser(ctx, &testUserData)
	if err == nil  {
		t.Fatalf("no error for invalid email")
	} else {
		if !errors.Is(err, domain.ErrInvalidEmail) {
			t.Errorf("bad error test, wanted: %v, got: %v", domain.ErrInvalidEmail, err)
		}
	}
}

func Test_AddUser_EmptyFirstName(t *testing.T) {
	repo := mocks.MockRepo{}
	testService := Service{repo: &repo}
	testUserData := domain.UserData {
		FirstName: "",
		LastName: "bar",
		Email: "foo.bar@eg.com",
	}
	err := testService.AddUser(ctx, &testUserData)
	if err == nil  {
		t.Fatalf("no error for empty first name")
	} else {
		expectedError := "empty first name"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}
}
func Test_AddUser_EmptyLastName(t *testing.T) {
	repo := mocks.MockRepo{}
	testService := Service{repo: &repo}
	testUserData := domain.UserData {
		FirstName: "foo",
		LastName: "",
		Email: "foo.bar@eg.com",
	}
	err := testService.AddUser(ctx, &testUserData)
	if err == nil  {
		t.Fatalf("no error for empty last name")
	} else {
		expectedError := "empty last name"
		if err.Error() != expectedError {
			t.Errorf("bad error test, wanted: %s, got: %s", expectedError, err.Error())
		}
	}
}

func Test_AddUser_DuplicateUser(t *testing.T) {
	repo := mocks.MockRepo{}
	testService := Service{repo: &repo}
	testUserData := domain.UserData {
		FirstName: "foo",
		LastName: "bar",
		Email: "foo.bar@eg.com",
	}
	err := testService.AddUser(ctx, &testUserData)
	if err != nil {
		t.Fatalf("error adding User, err: %v", err)
	}

	err = testService.AddUser(ctx, &testUserData)
	if err == nil  {
		t.Fatalf("no error for adding duplicate user")
	} else {
		if !errors.Is(err, domain.ErrUserAlreadyExists) {
			t.Errorf("bad error test, wanted: %s, got: %s", domain.ErrUserAlreadyExists, err)
		}
	}

	if cnt, _ := repo.GetCountOfUsers(ctx); cnt != 1 {
		t.Fatalf("bad user count, exp: %v, got: %v", 1, cnt)
	}
}