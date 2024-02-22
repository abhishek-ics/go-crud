package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
	"fmt"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
	"github.com/abhishek-ics/go-mysql-api/user"
)

var db *sql.DB

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    db = connectDB()
    defer db.Close()

    router := mux.NewRouter()
    router.HandleFunc("/users", user.GetUsers(db)).Methods("GET")
    router.HandleFunc("/users/{id}", user.GetUser(db)).Methods("GET")
    router.HandleFunc("/users", user.CreateUser(db)).Methods("POST")
    router.HandleFunc("/users/{id}", user.UpdateUser(db)).Methods("PUT")
    router.HandleFunc("/users/{id}", user.DeleteUser(db)).Methods("DELETE")

    port := ":8080"
    log.Printf("Server is listening on port %s", port)
    log.Fatal(http.ListenAndServe(port, router))
}

func connectDB() *sql.DB {
    dbUsername := os.Getenv("DB_USERNAME")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)

    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Fatal("Error connecting to database")
    }
    return db
}
