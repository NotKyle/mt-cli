package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"

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

func GetAllTasks(output bool) []Task {

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

	if output {
		for _, task := range tasks {
			color.Cyan("Task Details:")
			color.Green("Task ID: %d\n", task.ID)
			color.Green("Task Name: %s\n", task.Name)
			fmt.Println("\n")
		}

		return nil
	}

	return tasks
}

func GetTaskDetailsByIndex(index int) {
	tasks := GetAllTasks(false)

	if index < 0 || index >= len(tasks) {
		color.Red("Invalid task index\n")
		os.Exit(1)
	}

	if len(tasks) == 0 {
		color.Red("No tasks found\n")
		os.Exit(1)
	}

	taskComments := GetTaskDetails(tasks[index].ID)

	color.Cyan("Task Details:")
	color.Green("Task ID: %d\n", tasks[index].ID)
	color.Green("Task Name: %s\n", tasks[index].Name)
	color.Green("Task Section: %s\n", tasks[index].SectionName)

	color.Cyan("\nTask Comments:")
	for _, comment := range taskComments {
		color.Green("%s\n", comment.Text)
	}
}

func GetTaskDetails(taskToGet int) []TaskComment {
	requestURL := fmt.Sprintf("%s/tasks/%d/comments",
		os.Getenv("MEISTERTASK_API_BASE_URL"),
		taskToGet)

	resBody, err := MakeRequest(requestURL)
	if err != nil {
		color.Red("Error making request: %s\n", err)
		os.Exit(1)
	}

	var taskComments []TaskComment
	err = json.Unmarshal(resBody, &taskComments)
	if err != nil {
		color.Red("Error unmarshalling response: %s\n", err)
		os.Exit(1)
	}

	return taskComments
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

	taskToGet := flag.Int(
		"task",
		-1,
		"The task to fetch more details for",
	)

	flag.Parse()

	if *taskToGet != -1 {
		GetTaskDetailsByIndex(*taskToGet)
	} else {
		GetAllTasks(true)
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
