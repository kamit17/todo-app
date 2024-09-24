package main

import (
    "testing"
)

func TestTodoCreation(t *testing.T) {
    todo := Todo{Task: "Test Task"}
    if todo.Task != "Test Task" {
        t.Errorf("Expected 'Test Task', got '%s'", todo.Task)
    }
}

