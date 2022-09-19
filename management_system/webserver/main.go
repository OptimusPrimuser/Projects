package main

import (
	"webserver/globals"
	"webserver/rabbit_mq"
	"webserver/redis_store"
	"webserver/ws"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Entry struct {
	PersonID  string `json:"PersonID"`
	LastName  string `json:"LastName"`
	FirstName string `json:"FirstName"`
	Address   string `json:"Address"`
	City      string `json:"City"`
	Age       int    `json:"Age"`
}

func main() {
	globals.RMQ = rabbit_mq.RMQ{}
	globals.RMQ.Init("test1")
	defer globals.RMQ.Close()
	globals.Redis = redis_store.Redis{}
	globals.Redis.Init()
	defer globals.Redis.Close()
	// globals.RMQ.Publish(
	// 	ws.JobEntry{
	// 		PersonID:  1,
	// 		LastName:  "LOL",
	// 		FirstName: "LOL",
	// 		Address:   "abv",
	// 		City:      "iop",
	// 		Age:       23,
	// 		JobID:     "5",
	// 		Action:    "lol",
	// 	},
	// )
	// temp, x := globals.Redis.Get("4")
	// fmt.Println(temp)
	// fmt.Println("err", x)
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/up", ws.HealthCheck)
	e.GET("/get", ws.Get)
	e.GET("/jobStatus", ws.CheckJob)
	e.GET("/jobSolution", ws.GetSolution)
	e.POST("/add", ws.Add)
	e.DELETE("/delete", ws.Delete)
	e.PUT("/update", ws.Update)
	e.Logger.Fatal(e.Start(":1323"))
}
