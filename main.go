package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/umeshj346/helloWorldServer/internal/db"
	"github.com/umeshj346/helloWorldServer/internal/users"
)

type UserData struct {
    FirstName   string
    LastName    string
    Email       string
}

type server struct {
    userManager *users.Manager
}

func main() {
    err := godotenv.Load()
    if err != nil {
        slog.Error("error loading env file")
        os.Exit(1)
    }
    db := db.NewPostgresDB(os.Getenv("DATABASE_URL"))
    manager := users.NewManager(db)
    defer manager.Shutdown()

    s := server {
        userManager: manager,
    }

    mux := http.NewServeMux()

    httpServer := &http.Server{
        Addr:    ":80",
        Handler: mux,
    }

    mux.HandleFunc("/{$}", handleWelcome)
    mux.HandleFunc("/goodbye/", handleGoodBye)
    mux.HandleFunc("/hello/", handleHelloParameterized)
    mux.HandleFunc("/response/{user}/hello/", handleUserResponseHello)
    mux.HandleFunc("/user/hello", s.handleHelloHeader)
    mux.HandleFunc("POST /json", handleJson)
    mux.HandleFunc("POST /add-user", s.addUser)
    mux.HandleFunc("POST /get-user", s.getUser)

    go func() {
        slog.Info("starting server")
        err := httpServer.ListenAndServe()
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            slog.Error("HTTP server error", "err", err)
            os.Exit(2)
        }
    }()

    var wg sync.WaitGroup
    wg.Add(1)
    go func ()  {
        defer wg.Done()
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

        <-sigChan
        slog.Info("shutting down server")
         
        shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownCancel()
        err := httpServer.Shutdown(shutdownCtx)
        if err != nil {
            slog.Error("error shutting down HTTP server", "err", err)
        }
    }()
    wg.Wait()
    slog.Info("server shutdown is complete")
}

func writeHelloUser(w http.ResponseWriter, userName string) {
    var output bytes.Buffer
    output.WriteString("Hello ")
    output.WriteString(userName)
    output.WriteString("!\n")

    _, err := w.Write(output.Bytes())
    if err != nil {
        slog.Error("error writing response body", "err", err)
        return
    }
}

func writeHelloUserData(w http.ResponseWriter, userData *UserData) {
    output := fmt.Sprintf("Hello %v %v!\nYour Email is %v",
                            userData.FirstName, userData.LastName, userData.Email)  

    _, err := w.Write([]byte(output))
    if err != nil {
        slog.Error("error writing response body", "err", err)
        return
    }
}

func convertUserToUserData(u *users.User) *UserData {
    if u == nil {
        return nil
    }
    return &UserData{
        FirstName: u.FirstName,
        LastName: u.LastName,
        Email: u.Email.Address,
    }
}

func handleWelcome(w http.ResponseWriter, _ *http.Request) {
    n, err := w.Write([]byte("Welcome to my website!\n"))

    if err != nil {
        slog.Error("error writing response", "err", err)
        return
    }
    fmt.Printf("%d bytes written\n", n)
}

func handleGoodBye(w http.ResponseWriter, _ *http.Request) {
    n, err := w.Write([]byte("GoodBye\n"))

    if err != nil {
        slog.Error("error writing response", "err", err)
        return
    }
    fmt.Printf("%d bytes written\n", n)
}

func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    userList := params["user"]

    userName := "User"
    if len(userList) > 0 {
        userName = userList[0]
    }

    // Why do it this way?
    writeHelloUser(w, userName)
}

func handleUserResponseHello(w http.ResponseWriter, r *http.Request) {
    userName := r.PathValue("user")

    writeHelloUser(w, userName)
}

func (s *server) handleHelloHeader(w http.ResponseWriter, r *http.Request) {
    firstName, lastName := r.Header.Get("userFirst"), r.Header.Get("userLast")
    user, err := s.userManager.GetUserByName(firstName, lastName)

    if err != nil {
        if err == users.ErrNoResultFound {
            http.Error(w, "no users found", http.StatusNotFound)
        } else {
            slog.Error("error retrieving user", "err", err)
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }
    converted := convertUserToUserData(user)
    writeHelloUserData(w, converted)
}

func handleJson(w http.ResponseWriter, r *http.Request) {
    byteData, err := io.ReadAll(r.Body)

    if err != nil {
        slog.Error("error reading response body", "err", err)
        http.Error(w, "bad request body", http.StatusBadRequest)
        return
    }

    if len(byteData) == 0 {
        http.Error(w, "bad request body!", http.StatusBadRequest)
        return
    }

    var reqData UserData
    err = json.Unmarshal(byteData, &reqData) 
    if err != nil {
        slog.Error("error unmarshalling request body", "err", err)
        http.Error(w, "error parsing request json ", http.StatusBadRequest)
        return
    }
    if reqData.FirstName == "" {
        http.Error(w, "invalid username provided!", http.StatusBadRequest)
        return
    }

    writeHelloUser(w, reqData.FirstName)

}

func (s *server) addUser(w http.ResponseWriter, r *http.Request) {
    contentType := r.Header.Get("Content-Type")
    if contentType != "application/json" {
        http.Error(w, fmt.Sprintf("unsupported Content-Type header: %q", contentType), http.StatusUnsupportedMediaType)
        return
    }
    requestBody := http.MaxBytesReader(w, r.Body, 1048576)

    decoder := json.NewDecoder(requestBody)
    decoder.DisallowUnknownFields()

    var u UserData
    err := decoder.Decode(&u)
    
    if err != nil {
        slog.Error("error decoding addUser request body", "err", err)
        http.Error(w, "bad request body", http.StatusBadRequest)
        return
    }

    err = s.userManager.AddUser(u.FirstName, u.LastName, u.Email)
    if err != nil {
        http.Error(w, fmt.Sprintf("error adding user: %v\n", err), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)

}

func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
    contentType := r.Header.Get("Content-Type")
    if contentType != "application/json" {
        http.Error(w, fmt.Sprintf("unsupported Content-Type header: %q", contentType), http.StatusUnsupportedMediaType)
        return
    }
    requestBody := http.MaxBytesReader(w, r.Body, 1048576)

    decoder := json.NewDecoder(requestBody)
    decoder.DisallowUnknownFields()

    var u UserData
    err := decoder.Decode(&u)
    
    if err != nil {
        slog.Error("error decoding getUser request body", "err", err)
        http.Error(w, "bad request body", http.StatusBadRequest)
        return
    }

    foundUser, err := s.userManager.GetUserByName(u.FirstName, u.LastName)
    if err != nil {
        if errors.Is(err, users.ErrNoResultFound) {
            http.Error(w, "no users found", http.StatusNotFound)
        } else {
            slog.Error("error retrieving user", "err", err)
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }

    converted := convertUserToUserData(foundUser)
    marshalled, err := json.Marshal(*converted)
    if err != nil {
        slog.Error("error marshaling getUser reposne", "err", err)
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    _, err = w.Write(marshalled)
    if err != nil {
        slog.Error("error writing getUser response", "err", err)
    }

    w.WriteHeader(http.StatusCreated)
}
