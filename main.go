package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	database "github.com/mishaRomanov/learn-fiber/db"
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
	router := mux.NewRouter()

	//create a server
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  300 * time.Millisecond,
		WriteTimeout: 300 * time.Millisecond,
		Handler:      router,
	}

	//open a database
	db, err := database.OpenDb()
	if err != nil {
		logrus.Errorf("!WARNING!: Error while opening database! %v\n", err)
	}
	defer db.Close()
	logrus.Infoln("Database connect successful!")

	//Handles /about request
	router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logrus.Infoln("An attempt to acces /about page with different method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}
		logrus.Infof("New %s request", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`This small app is CRUD to-do list-type application.
Send a POST-request to create a new task: /tasks/add 
and monitor it by visiting /tasks`))
	})

	//handler that adds a new task into a database (POST REQUEST)
	router.HandleFunc("/tasks/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf("Method %s not allowed on that endpoint", r.Method)))
			return
		}
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
			logrus.Errorf("Error while parsing the request body! %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		sqlStatement := `INSERT INTO tasks (Description) VALUES ($1);`
		result, err := db.Exec(sqlStatement, newTask.Desc)
		if err != nil {
			logrus.Errorf(":Error while inserting values into database! %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data written. New task added"))
		logrus.Infof("%v", result)
	})

	//here we create a handler for DELETE request
	router.HandleFunc("/tasks/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		query := fmt.Sprintf(`SELECT id FROM tasks WHERE id = $1`)
		rows, err := db.Exec(query, id)
		if err != nil {
			logrus.Infof("Error while performing SQL query: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error. Try again"))
			return
		}
		res, _ := rows.RowsAffected()
		if res == 0 {
			logrus.Infof("Can't seem to find a task with id %s ", id)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("A task with such ID doesn't exist"))
			return
		}
		deleteQuery := fmt.Sprintf(`DELETE FROM tasks WHERE id = $1`)
		_, err = db.Exec(deleteQuery, id)
		if err != nil {
			logrus.Infof("Error while performing SQL query: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal error"))
			return
		}
		logrus.Infoln("Delete successful")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//here we start a server
	logrus.Fatal(server.ListenAndServe())

}
