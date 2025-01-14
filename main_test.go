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

func Test_main(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

