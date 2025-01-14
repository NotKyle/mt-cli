package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	vars := getVars()
	API_BASE_URL := vars[0]
	LANE_ID := vars[1]

	requestURL := fmt.Sprintf("%s/sections/%s/tasks?assigned_to_me=true&status=open",
		API_BASE_URL,
		LANE_ID)

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

	homeDir, err := os.UserHomeDir()

	lastAccessed := time.Now().Unix()

	tasksFile, err := os.OpenFile(homeDir+"/tasks.json", os.O_RDWR|os.O_CREATE, 0755)
	var asJson []byte

	if err != nil {
		fmt.Println("Tasks file does not exist, creating it now", err)

		_ = os.NewFile(0, homeDir+"/tasks.json")
	}

	tasksFile, err = os.OpenFile(homeDir+"/tasks.json", os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		color.Red("Error creating tasks file %s\n", err)
		os.Exit(1)
	}

	if output {
		color.Cyan("Tasks Found: %d\n", len(tasks))

		asJson = []byte(fmt.Sprintf(`{
			"lastAccessed": "%d",
			"tasks": %s
		}`, lastAccessed, resBody))

		_, err = tasksFile.Write(asJson)

		if err != nil {
			color.Red("Error writing to tasks file %s\n", err)
			os.Exit(1)
		}

		for _, task := range tasks {
			color.Cyan("Task Details:")
			color.Green("Task ID: %d\n", task.ID)
			color.Green("Task Name: %s\n", task.Name)
			fmt.Println()
		}

		return nil
	}

	return tasks
}

func getVars() [3]string {
	// Attempt to load .env only if the file exists
	if _, fileErr := os.Stat(".env"); fileErr == nil {
		err := godotenv.Load()

		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Read variables from environment
	API_BASE_URL := os.Getenv("MEISTERTASK_API_BASE_URL")
	LANE_ID := os.Getenv("LANE_ID")
	API_TOKEN := os.Getenv("MEISTERTASK_API_KEY")

	return [3]string{API_BASE_URL, LANE_ID, API_TOKEN}
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
	vars := getVars()
	API_BASE_URL := vars[0]

	requestURL := fmt.Sprintf("%s/tasks/%d/comments",
		API_BASE_URL,
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
	vars := getVars()
	API_TOKEN := vars[2]

	var bearer = "Bearer " + API_TOKEN

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

	// Get environment variables
	vars := getVars()

	API_BASE_URL := vars[0]
	LANE_ID := vars[1]

	// Ensure critical environment variables are present
	if API_BASE_URL == "" {
		log.Fatal("Error: MEISTERTASK_API_BASE_URL is not set")
	}
	if LANE_ID == "" {
		log.Fatal("Error: LANE_ID is not set")
	}

	taskToGet := flag.Int(
		"get",
		-1,
		"The task to fetch more details for",
	)

	flag.Parse()

	if *taskToGet != -1 {
		index := (*taskToGet - 1)
		GetTaskDetailsByIndex(index)
	} else {
		GetAllTasks(true)
	}
}
