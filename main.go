package main

import (
	"encoding/json"
	"fmt"
	database "github.com/mishaRomanov/learn-fiber/db"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// UpdateTask is used to unmrashall all jsons sent to PATCH tasks
type Task struct {
	Desc string `json:"desc"`
}

func main() {

	log.Println(`|------------Starting a server!------------|`)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered after a panic:\n", r)
		}
	}()

	//create a server
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  300 * time.Millisecond,
		WriteTimeout: 300 * time.Millisecond,
	}
	//open a database
	db, err := database.OpenDb()
	if err != nil {
		log.Fatalf("Error while opening database! %v\n", err)
	}
	defer db.Close()
	log.Println("Database connect successful!")

	//create a new logger
	var logger log.Logger
	file, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error while opening a file: %v", err)
	}
	//here we defer the closing of the file
	defer file.Close()

	//we define the output of our logger
	//this is very simple
	logger.SetOutput(file)

	//Handles /about request
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s: New %s request", time.Now().Format(time.RFC822), r.Method)
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`This small app is CRUD to-do list-type application.
Send a POST-request to create a new task: /tasks/new 
and monitor it by visiting /tasks`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Method not allowed"))
		}
	})

	//handler that adds a new task into a database
	http.HandleFunc("/tasks/add", func(w http.ResponseWriter, r *http.Request) {
		newTask := &Task{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s:Error while reading the request body! %v\n", time.Now().Format(time.RFC822), err)
			logger.Printf("Error while reading the request body! %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		err = json.Unmarshal(body, newTask)
		if err != nil {
			log.Printf("%s:Error while parsing the request body! %v\n", time.Now().Format(time.RFC822), err)
			logger.Printf("Error while parsing the request body! %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		if newTask.Desc == "" {
			panic("Empty description!")
		}
		sqlStatement := `INSERT INTO tasks (Description)
VALUES ($1);`
		result, err := db.Exec(sqlStatement, newTask.Desc)
		if err != nil {
			log.Printf("%s:Error while inserting values into database! %v\n", time.Now().Format(time.RFC822), err)
			logger.Printf("Error while inserting values into database! %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data written. New task added"))
		log.Println(result.RowsAffected())
	})

	//here we start a server

	logger.Fatal(server.ListenAndServe())

}
