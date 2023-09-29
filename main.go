package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

// TasksTracker is a storage for our tasks
var TasksTracker = map[int]*Task{}

type Task struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

// UpdateTask is used to unmrashall all jsons sent to PATCH tasks
type UpdateTask struct {
	Desc string `json:"desc"`
}

// About func returns information about our app
func About(c *fiber.Ctx) error {
	logrus.Infof("New About request  at:%s", time.Now().Format(time.RFC822))
	return c.SendString(`This small app is CRUD to-do list-type application.
Send a POST-request to create a new task: /tasks/new 
and monitor it by visiting /tasks`)
}

func main() {
	os.Create("logs.log")
	file, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	// start a server
	webApp := fiber.New()

	webApp.Use(requestid.New())
	webApp.Use(logger.New(logger.Config{
		Format: "[${locals:requestid}]: ${method} ${path} - ${status}\n",
		Output: file,
	}))

	//all handlers are stored here
	webApp.Get("/about", About)

	webApp.Get("/", About)

	//Returns all tasks
	webApp.Get("/tasks/all", func(ctx *fiber.Ctx) error {
		if len(TasksTracker) == 0 {
			ctx.Status(fiber.StatusOK).SendString("There are currently no tasks at all!")
		}
		return ctx.Status(fiber.StatusOK).JSON(TasksTracker)
	})

	//func that returns specific task
	webApp.Get("/tasks/:id", func(ctx *fiber.Ctx) error {
		param := ctx.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logrus.Infof("Error while converting string to int")
			return ctx.Status(fiber.StatusConflict).SendString(fmt.Sprintf(`Invalid id "%s"`, param))
		}
		logrus.Infof("Succesfull request")
		return ctx.Status(fiber.StatusOK).JSON(TasksTracker[id])

	})

	//POSTS a new task
	webApp.Post("/tasks/new", func(ctx *fiber.Ctx) error {
		newTask := Task{}
		err := ctx.BodyParser(&newTask)
		if err != nil {
			logrus.Infof("Error while unmarshalling json to struct at request %v", err)
			return ctx.Status(fiber.StatusUnprocessableEntity).SendString("Error while creating a new task!")
		}
		//adding a new task
		TasksTracker[newTask.Id] = &newTask
		logrus.Infof("New task added")
		return ctx.Status(fiber.StatusOK).SendString(fmt.Sprintf("A new task created! Unique task id is:%d", newTask.Id))
	})

	//PATCH already existing task
	webApp.Patch("/tasks/:id", func(ctx *fiber.Ctx) error {
		param := ctx.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logrus.Infof("Error while converting string to int")
		}
		new := UpdateTask{}
		err = ctx.BodyParser(&new)
		if err != nil {
			logrus.Infof("Error while unmarshalling json to struct %v", err)
		}
		old := TasksTracker[id].Desc

		TasksTracker[id].Desc = new.Desc
		logrus.Infof("Updated task number %d", id)
		return ctx.Status(fiber.StatusOK).SendString(fmt.Sprintf(`Updated your task from "%s" to "%s"`, old, TasksTracker[id].Desc))

	})
	//DELETE a task by id
	webApp.Delete("/tasks/:id", func(ctx *fiber.Ctx) error {
		param := ctx.Params("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logrus.Infof("Error while converting string to int")
		}
		delete(TasksTracker, id)
		logrus.Infof("Deleted task number %d", id)

		return ctx.Status(fiber.StatusOK).SendString(fmt.Sprintf("Deleted task number %d", id))
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
