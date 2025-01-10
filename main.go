package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Task struct {
	ID int
	Name,
	SectionName string
}

type TaskComment struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func GetAllTasks() {

	requestURL := fmt.Sprintf("%s/sections/%s/tasks?assigned_to_me=true&status=open",
		os.Getenv("MEISTERTASK_API_BASE_URL"),
		os.Getenv("LANE_ID"))

	resBody, err := MakeRequest(requestURL)

	if err != nil {
		color.Red("Error making request %s\n", err)
		os.Exit(1)
	}

	// Convert response to struct
	var tasks []Task
	err = json.Unmarshal(resBody, &tasks)

	if err != nil {
		color.Red("Error unmarshalling response %s\n", err)
		os.Exit(1)
	}

	for _, task := range tasks {
		color.Cyan("Task Details:")
		color.Green("Task ID: %d\n", task.ID)
		color.Green("Task Name: %s\n", task.Name)
		fmt.Println("\n")
	}
}

func GetTaskDetails(taskToGet string) {
	requestURL := fmt.Sprintf("%s/tasks/%s/comments",
		os.Getenv("MEISTERTASK_API_BASE_URL"),
		taskToGet)

	resBody, err := MakeRequest(requestURL)
	if err != nil {
		color.Red("Error making request %s\n", err)
		os.Exit(1)
	}

	var taskComments []TaskComment

	err = json.Unmarshal(resBody, &taskComments)
	if err != nil {
		color.Red("Error unmarshalling response %s\n", err)
		os.Exit(1)
	}

	color.Cyan("Task Comment Details (%s):", taskToGet)
	for _, taskComment := range taskComments {
		color.Green("Task Comment: %s\n", taskComment.Text)
		fmt.Println("\n")
	}
}

func MakeRequest(url string) ([]byte, error) {

	var bearer = "Bearer " + os.Getenv("MEISTERTASK_API_KEY")

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Go-HTTP-Client")

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	taskToGet := flag.String(
		"task",
		"",
		"The task to fetch more details for",
	)

	flag.Parse()

	if *taskToGet != "" {
		GetTaskDetails(*taskToGet)
	} else {
		GetAllTasks()
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
