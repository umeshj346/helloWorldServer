package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
	"testing"

	"github.com/umeshj346/helloWorldServer/internal/users"
)

func Test_HandleWelcome(t *testing.T) {
	w := httptest.NewRecorder()
	
	handleWelcome(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome to my website!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %v, \ngot: %v\n", expectedMessage, w.Body.String())
	}
	
}

func Test_HandleGoodBye(t *testing.T) {
	w := httptest.NewRecorder()
	handleGoodBye(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("GoodBye\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleHelloParameterized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?user=TestMan", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}

}

func Test_HandleHelloParameterized_WrongParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?t=v", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != http.StatusOK {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}

}

func Test_HandleHelloParameterized_EmptyParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)

	desiredCode := http.StatusOK
	if w.Code != http.StatusOK {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}

}

func Test_HandleUserResponseHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/response/TestMan/hello", nil)
	req.SetPathValue("user", "TestMan")

	w := httptest.NewRecorder()

	handleUserResponseHello(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleHelloHeader(t *testing.T) {
	testFirstName, testLastName :=  "Test", "Man"
	testEmail := "foo@bar"

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}
	testManager.AddUser(testFirstName, testLastName, testEmail)

	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)
	req.Header.Set("userFirst", testFirstName)
	req.Header.Set("userLast", testLastName)

	w := httptest.NewRecorder()

	testServer.handleHelloHeader(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello Test Man!\nYour Email is foo@bar")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleHelloHeader_WrongHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)
	req.Header.Set("foo", "bar")

	w := httptest.NewRecorder()

	s := server {
		userManager: users.NewManager(),
	}
	s.handleHelloHeader(w, req)

	desiredCode := http.StatusNotFound
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("no users found\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleHelloHeader_NoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)

	w := httptest.NewRecorder()
	s := server{
		userManager: users.NewManager(),
	}
	s.handleHelloHeader(w, req)

	desiredCode := http.StatusNotFound
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("no users found\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleJson(t *testing.T) {
	testRequest := UserData {
		FirstName: "Test Man",
	}

	marshalledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshalledRequestBody))
	w := httptest.NewRecorder()

	handleJson(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello Test Man!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleJson_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/json", nil)
	w := httptest.NewRecorder()

	handleJson(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("bad request body!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_HandleJson_EmptyName(t *testing.T) {
	var testRequest UserData

	marshalledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshalledRequestBody))
	w := httptest.NewRecorder()

	handleJson(w, req)

	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v but got %v\n. body: %v", 
				desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("invalid username provided!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("wrong message expected: %q, \ngot: %q\n", expectedMessage , w.Body.String())
	}
}

func Test_AddUser(t *testing.T) {
	testUserData := UserData {
		FirstName: "Test",
		LastName: "Man",
		Email: "foo.bar@eg.com",
	}
	marshalledRequestBody, err := json.Marshal(testUserData)
	if err != nil {
		t.Fatalf("error marshalling testUser, err: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/add-user", 
								bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := server {
		userManager: testManager,
	}
	testServer.addUser(w, req)

	desiredCode := http.StatusCreated
	if w.Code != desiredCode {
		t.Errorf("bad response code, exp: %v, but got: %v", desiredCode, w.Code)
	}

	foundUser, err := testServer.userManager.GetUserByName("Test", "Man")
	if err != nil {
		t.Fatalf("error checking whether testUser is added, err: %v", err)
	}
	convUser := convertUserToUserData(foundUser)
	if !reflect.DeepEqual(convUser, &testUserData) {
		t.Errorf("bad retrieve of user\nwanted: %+v\ngot:%+v",  testUserData, *convUser)
	}
}

func Test_AddUser_BadHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/add-user", nil)
	req.Header.Set("Content-Type", "application/html")
	w := httptest.NewRecorder()
	s := server {
		userManager: users.NewManager(),
	}
	s.addUser(w, req)
	desiredCode := http.StatusUnsupportedMediaType
	if w.Code != desiredCode {
		t.Errorf("bad response code, wanted: %v, got: %v", desiredCode, w.Code)
	}

	expectedMessage := []byte("unsupported Content-Type header: \"application/html\"\n")
	if !bytes.Equal(w.Body.Bytes(), expectedMessage) {
		t.Errorf("bad response message\nwanted: %q\ngot:%q", expectedMessage, w.Body.Bytes())
	}
	
}

func Test_GetUser(t *testing.T) {
	testFirstName, testLastName, testEmail:= "Test", "Man", "foo@boo"

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	testQuery := UserData {
		FirstName: testFirstName,
		LastName: testLastName,
	}
	marshalledRequestBody, err := json.Marshal(testQuery)
	if err != nil {
		t.Fatalf("error marshalling testQuery, err: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/get-userr", bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	
	testServer.getUser(w, req)
	
	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, wanted: %v, got: %v", desiredCode, w.Code)
	}

	var resultUserData UserData
	decoder := json.NewDecoder(w.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&resultUserData)
	if err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	expectedUserData := UserData {
		FirstName: testFirstName,
		LastName: testLastName,
		Email: testEmail,
	}

	if !reflect.DeepEqual(resultUserData, expectedUserData) {
		t.Errorf("bad result\nwanted: %+v\ngot: %+v", expectedUserData, resultUserData)
	}
}

func Test_GetUser_BadHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/add-user", nil)
	req.Header.Set("Content-Type", "application/html")
	w := httptest.NewRecorder()
	s := server {
		userManager: users.NewManager(),
	}
	s.getUser(w, req)
	desiredCode := http.StatusUnsupportedMediaType
	if w.Code != desiredCode {
		t.Errorf("bad response code, wanted: %v, got: %v", desiredCode, w.Code)
	}

	expectedMessage := []byte("unsupported Content-Type header: \"application/html\"\n")
	if !bytes.Equal(w.Body.Bytes(), expectedMessage) {
		t.Errorf("bad response message\nwanted: %q\ngot:%q", expectedMessage, w.Body.Bytes())
	}
	
}

func Test_GetUser_NoUser(t *testing.T) {
	testFirstName, testLastName, testEmail:= "Test", "Man", "foo@boo"

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err != nil {
		t.Fatalf("error adding user, err: %v", err)
	}

	testQuery := UserData {
		FirstName: "Fake",
		LastName: "Query",
	}
	marshalledRequestBody, err := json.Marshal(testQuery)
	if err != nil {
		t.Fatalf("error marshalling testQuery, err: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/get-userr", bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	
	testServer.getUser(w, req)
	
	desiredCode := http.StatusNotFound
	if w.Code != desiredCode {
		t.Errorf("bad response code, wanted: %v, got: %v", desiredCode, w.Code)
	}

	expectedMessage := []byte("no users found\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad response message, wanted : %q, got: %q", expectedMessage, w.Body.Bytes())
	}
} 

func Test_ConverUserToUserData(t *testing.T) {
	testFirstName, testLastName := "Test", "Man"
	testEmail, err := mail.ParseAddress("foo.bar@eg.com")
	if err != nil {
		t.Fatalf("error parsing the email address")
	}
	testUser := users.User {
		FirstName: testFirstName,
		LastName: testLastName,
		Email: *testEmail,
	}
	resultUserData := convertUserToUserData(&testUser)
	if resultUserData == nil {
		t.Fatalf("bad conversion from user to userData")
	}

	expectedUserData := UserData {
		FirstName: testFirstName,
		LastName: testLastName,
		Email: testEmail.Address,
	}

	if !reflect.DeepEqual(resultUserData, &expectedUserData) {
		t.Errorf("bad conversion from user to userdata\nexp: %+v\ngot: %+v\n",
					expectedUserData, resultUserData)
	}
}

func Test_ConverUserToUserData_NilUser(t *testing.T) {
	resultUserData := convertUserToUserData(nil)
	if resultUserData != nil {
		t.Errorf("bad conversion from user to userData")
	}
}



