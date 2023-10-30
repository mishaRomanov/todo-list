package main

import (
	"encoding/json"
	database "github.com/mishaRomanov/learn-fiber/db"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// UpdateTask is used to unmrashall all jsons sent to PATCH tasks
type Task struct {
	Desc string `json:"desc"`
}

func main() {

	logrus.Infoln(`|------------Starting a server!------------|`)
	defer func() {
		if r := recover(); r != nil {
			logrus.Println("Recovered after a panic:\n", r)
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
		logrus.Errorf("!WARNING!: Error while opening database! %v\n", err)
	}
	defer db.Close()
	logrus.Infoln("Database connect successful!")

	//Handles /about request
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("New %s request", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`This small app is CRUD to-do list-type application.
Send a POST-request to create a new task: /tasks/new 
and monitor it by visiting /tasks`))
	})

	//handler that adds a new task into a database
	http.HandleFunc("/tasks/add", func(w http.ResponseWriter, r *http.Request) {
		newTask := &Task{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("%s:Error while reading the request body! %v\n", time.Now().Format(time.RFC822), err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		err = json.Unmarshal(body, newTask)
		if err != nil {
			logrus.Errorf(":Error while parsing the request body! %v\n", err)
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
			logrus.Errorf(":Error while inserting values into database! %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data written. New task added"))
		logrus.Info(result.RowsAffected())
	})

	//here we start a server

	logrus.Fatal(server.ListenAndServe())

}
