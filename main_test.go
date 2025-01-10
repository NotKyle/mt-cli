package main

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestTasks(t *testing.T) {
	godotenv.Load()

	tasks := GetAllTasks(false)

	if len(tasks) == 0 {
		t.Errorf("Expected tasks to be greater than 0")
	}
}
