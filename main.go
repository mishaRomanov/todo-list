package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

var requestID int //this variable is used to keep track of each request
// in case something goes bad we can track it in our server logs

// TasksTracker is a storage for our tasks
var TasksTracker = map[int]*Task{}

type Task struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
	Time time.Time
}

// UpdateTask is used to unmrashall all jsons sent to PATCH tasks
type UpdateTask struct {
	Desc string `json:"desc"`
}

// About func returns information about our app
func About(c *fiber.Ctx) error {
	requestID++
	logrus.Infof("New About request number %d at:%s", requestID, time.Now().Format(time.RFC822))
	return c.SendString(`This small app is CRUD to-do list-type application.
Send a POST-request to create a new task: /tasks/new 
and monitor it by visiting /tasks`)
}

func main() {
	// start a server
	webApp := fiber.New()

	//all handlers are stored here
	webApp.Get("/about", About)

	webApp.Get("/", About)

	//func that returns specific task
	webApp.Get("/tasks/:id", func(ctx *fiber.Ctx) error {
		requestID++
		param := ctx.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logrus.Infof("Error while converting string to int in request number %d", requestID)
			return ctx.Status(fiber.StatusConflict).SendString(fmt.Sprintf("Error:%v", err))
		}
		logrus.Infof("Succesfull request number %d", requestID)
		return ctx.Status(fiber.StatusOK).JSON(TasksTracker[id])

	})

	//POSTS a new task
	webApp.Post("/tasks/new", func(ctx *fiber.Ctx) error {
		requestID++
		newTask := Task{}
		err := ctx.BodyParser(&newTask)
		if err != nil {
			logrus.Infof("Error while unmarshalling json to struct at request number %d: %v", requestID, err)
			return ctx.Status(fiber.StatusUnprocessableEntity).SendString("Error while creating a new task!")
		}
		//adding a new task
		TasksTracker[newTask.Id] = &newTask
		logrus.Infof("New task added: request number %d", requestID)
		return ctx.Status(fiber.StatusOK).SendString(fmt.Sprintf("A new task created! Unique task id is:%d", newTask.Id))
	})

	//PATCH already existing task
	webApp.Patch("/tasks/:id", func(ctx *fiber.Ctx) error {
		requestID++
		param := ctx.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logrus.Infof("Error while converting string to int in request number %d", requestID)
		}
		new := UpdateTask{}
		err = ctx.BodyParser(&new)
		if err != nil {
			logrus.Infof("Error while unmarshalling json to struct at request number %d: %v", requestID, err)
		}
		old := TasksTracker[id].Desc

		TasksTracker[id].Desc = new.Desc
		logrus.Infof("Updated task number %d with request %d", id, requestID)
		return ctx.Status(fiber.StatusOK).SendString(fmt.Sprintf(`Updated your task from "%s" to "%s"`, old, TasksTracker[id].Desc))

	})

	logrus.Fatal(webApp.Listen(":8080"))
}
