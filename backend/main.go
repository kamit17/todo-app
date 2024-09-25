package main

import (
    "encoding/json"
    "net/http"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
    "os"
)

type Todo struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Task string `json:"task"`
}

var db *gorm.DB
var logFile *os.File

type LogMessage struct {
    Level   string `json:"level"`
    Message string `json:"message"`
}

func logInfo(message string) {
    logMessage := LogMessage{Level: "INFO", Message: message}
    jsonLog, _ := json.Marshal(logMessage)
    log.Println(string(jsonLog))
}

func logError(message string) {
    logMessage := LogMessage{Level: "ERROR", Message: message}
    jsonLog, _ := json.Marshal(logMessage)
    log.Println(string(jsonLog))
}

func init() {
    var err error

    // Set up logging to file
    os.MkdirAll("logs", os.ModePerm) // Create logs directory if it doesn't exist
    logFile, err = os.OpenFile("logs/application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Error opening log file:", err)
    }
    log.SetOutput(logFile) // Redirect log output to the file

    // Initialize database
    db, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
    if err != nil {
        logError("failed to connect to the database")
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
    logInfo("Fetched todos") // Structured logging
    json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    handleCORS(w, r) // Call CORS handler
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        logError("Failed to decode todo: " + err.Error())
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    db.Create(&todo)
    logInfo("Created todo") // Structured logging
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func main() {
    // Serve static files (CSS, JS)
    http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("../frontend/"))))

    // Handle root route to serve index.html
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "../frontend/index.html") // Serve the HTML file
        logInfo("Served index.html") // Structured logging
    })

    // Handle /todos route
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
        handleCORS(w, r) // Call CORS handler for all requests
        logInfo("Received " + r.Method + " request for " + r.URL.Path) // Structured logging
        switch r.Method {
        case http.MethodGet:
            getTodos(w, r)
        case http.MethodPost:
            createTodo(w, r)
        default:
            logError("Method not allowed: " + r.Method) // Structured logging
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    logInfo("Server is running on port 3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

