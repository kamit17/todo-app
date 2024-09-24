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

// Handle CORS
func handleCORS(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow specific methods
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow specific headers

    // Handle preflight requests
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content for preflight
        return
    }
}

func getTodos(w http.ResponseWriter, r *http.Request) {
    handleCORS(w, r) // Call CORS handler
    var todos []Todo
    db.Find(&todos)
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    handleCORS(w, r) // Call CORS handler
    var todo Todo
    json.NewDecoder(r.Body).Decode(&todo)
    db.Create(&todo)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func main() {
    // Serve static files (CSS, JS)
    http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend/"))))

    // Handle root route to serve index.html
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "../frontend/index.html") // Serve the HTML file
    })

    // Handle /todos route
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
        handleCORS(w, r) // Call CORS handler for all requests
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

