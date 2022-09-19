package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"webserver/globals"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JobEntry struct {
	PersonID  string `json:"PersonID"`
	LastName  string `json:"LastName"`
	FirstName string `json:"FirstName"`
	Address   string `json:"Address"`
	City      string `json:"City"`
	Age       int    `json:"Age"`
	JobID     string `json:"jobID"`
	Action    string `json:"action"`
}

type ActionStatus struct {
	Status      string `json:"status"`
	SolutionKey string `json:"solutionKey"`
}

func generateUniqueID() string {
	return uuid.New().String()
}

func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "UP")
}

func Get(c echo.Context) error {
	var entry JobEntry
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	err = json.Unmarshal(body, &entry)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	val, err := globals.Redis.Get(entry.PersonID + "-row")
	if err == nil {
		return c.String(http.StatusOK, val)
	}
	entry.JobID = generateUniqueID()
	entry.Action = "get"
	globals.RMQ.Publish(entry)
	globals.Redis.Add(entry.JobID, ActionStatus{"pending", ""})
	return c.String(http.StatusOK, entry.JobID)
}

func Add(c echo.Context) error {
	var entry JobEntry
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	err = json.Unmarshal(body, &entry)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	entry.JobID = generateUniqueID()
	entry.Action = "add"
	entry.PersonID = generateUniqueID()
	globals.RMQ.Publish(entry)
	globals.Redis.Add(entry.JobID, ActionStatus{"pending", ""})
	return c.String(http.StatusOK, entry.JobID)
}

func Update(c echo.Context) error {
	var entry JobEntry
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	err = json.Unmarshal(body, &entry)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	entry.JobID = generateUniqueID()
	entry.Action = "update"
	globals.RMQ.Publish(entry)
	globals.Redis.Add(entry.JobID, ActionStatus{"pending", ""})
	return c.String(http.StatusOK, entry.JobID)
}

func Delete(c echo.Context) error {
	var entry JobEntry
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	err = json.Unmarshal(body, &entry)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	entry.JobID = generateUniqueID()
	entry.Action = "delete"
	globals.RMQ.Publish(entry)
	globals.Redis.Add(entry.JobID, ActionStatus{"pending", ""})
	return c.String(http.StatusOK, entry.JobID)
}

func CheckJob(c echo.Context) error {
	var val JobEntry
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	err = json.Unmarshal(body, &val)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	retVal, err := globals.Redis.Get(val.JobID)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	return c.String(http.StatusOK, retVal)
}

func GetSolution(c echo.Context) error {
	var val map[string]string
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("error", err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	err = json.Unmarshal(body, &val)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	retVal, err := globals.Redis.Get(val["key"])
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error)
	}
	return c.String(http.StatusOK, retVal)
}
