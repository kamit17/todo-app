package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&Todo{})
    return db
}

func TestGetTodos(t *testing.T) {
    db = setupDB()
    req, err := http.NewRequest("GET", "/todos", nil)
    if err != nil {
        t.Fatal(err)
    }
    w := httptest.NewRecorder()
    getTodos(w, req)

    res := w.Result()
    if res.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK; got %v", res.Status)
    }
}

func TestCreateTodo(t *testing.T) {
    db = setupDB()
    payload := `{"task": "Test Task"}`
    req, err := http.NewRequest("POST", "/todos", bytes.NewBufferString(payload))
    if err != nil {
        t.Fatal(err)
    }
    w := httptest.NewRecorder()
    createTodo(w, req)

    res := w.Result()
    if res.StatusCode != http.StatusCreated {
        t.Errorf("Expected status Created; got %v", res.Status)
    }
}

