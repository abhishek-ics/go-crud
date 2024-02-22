package user

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}

func GetUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var users []User
        rows, err := db.Query("SELECT * FROM users")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        for rows.Next() {
            var user User
            err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
            if err != nil {
                log.Fatal(err)
            }
            users = append(users, user)
        }

        respondWithJSON(w, http.StatusOK, users)
    }
}

func GetUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        id := params["id"]

        var user User
        err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
        if err != nil {
            if err == sql.ErrNoRows {
                respondWithError(w, http.StatusNotFound, "User not found")
                return
            }
            respondWithError(w, http.StatusInternalServerError, "Error retrieving user")
            log.Println(err)
            return
        }

        respondWithJSON(w, http.StatusOK, user)
    }
}

func CreateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user User
        json.NewDecoder(r.Body).Decode(&user)

        result, err := db.Exec("INSERT INTO users (first_name, last_name, email) VALUES (?, ?, ?)", user.FirstName, user.LastName, user.Email)
        if err != nil {
            log.Fatal(err)
        }

        lastInsertID, err := result.LastInsertId()
        if err != nil {
            log.Fatal(err)
        }

        user.ID = int(lastInsertID)

        respondWithJSON(w, http.StatusCreated, user)
    }
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        id := params["id"]

        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            respondWithError(w, http.StatusBadRequest, "Invalid request payload")
            return
        }

		exists, err := userExistsByID(db, id)
        if err != nil {
            log.Fatal(err)
        }
        if !exists {
            respondWithError(w, http.StatusNotFound, "User not found")
            return
        }

        // Update the user
        _, err = db.Exec("UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?", user.FirstName, user.LastName, user.Email, id)
        if err != nil {
            log.Fatal(err)
        }

        user.ID, _ = strconv.Atoi(id)

        respondWithJSON(w, http.StatusOK, user)
    }
}


func DeleteUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userID := vars["id"]

		exists, err := userExistsByID(db, userID)
        if err != nil {
            log.Fatal(err)
        }
        if !exists {
            respondWithError(w, http.StatusNotFound, "User not found")
            return
        }

        // Delete the user
        err = deleteUserByID(db, userID)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Error deleting user")
            return
        }

        respondWithJSON(w, http.StatusNoContent, "Delete user successfully")
    }
}


func deleteUserByID(db *sql.DB, userID string) error {
	 // Check if the user exists
	 
    _, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
    return err
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        http.Error(w, "JSON encoding error", http.StatusInternalServerError)
        return
    }
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, ErrorResponse{Message: message})
}

type ErrorResponse struct {
    Message string `json:"message"`
}

func userExistsByID(db *sql.DB, userID string) (bool, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
