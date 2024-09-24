package main

import (
    "encoding/json"
    "net/http"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
)

type Todo struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Task string `json:"task"`
}

var db *gorm.DB

func init() {
    var err error
    db, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect to the database")
    }
    db.AutoMigrate(&Todo{})
}

func getTodos(w http.ResponseWriter, r *http.Request) {
    var todos []Todo
    db.Find(&todos)
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    json.NewDecoder(r.Body).Decode(&todo)
    db.Create(&todo)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func main() {
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getTodos(w, r)
        case http.MethodPost:
            createTodo(w, r)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    log.Println("Server is running on port 3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

