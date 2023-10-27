package main

import (
	"log"
	"net/http"
	"os"
)

var TasksTracker = map[int]*Task{}

type Task struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

// UpdateTask is used to unmrashall all jsons sent to PATCH tasks
type UpdateTask struct {
	Desc string `json:"desc"`
}

func main() {
	//create a new logger
	var logger log.Logger
	file, err := os.Open("logs")
	if err != nil {
		log.Print("Error while opening a file: %v", err)
	}
	//here we defer the closing of the file
	defer file.Close()

	//we define the output of our logger
	logger.SetOutput(file)

	//here we start a server
	logger.Fatal(http.ListenAndServe(":8080", nil))

}