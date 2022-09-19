package handler

import (
	"backend/globals"
	"encoding/json"
	"fmt"
	"time"
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

func WaitForLock(key string) {
	rkey := fmt.Sprintf("%s-lock", key)
	_, err := globals.Redis.Get(rkey)
	for err == nil {
		time.Sleep(1 * time.Second)
		_, err = globals.Redis.Get(rkey)
	}
}

func CreateLock(key string) error {
	rkey := fmt.Sprintf("%s-lock", key)
	err := globals.Redis.Add(rkey, "locked")
	if err != nil {
		return err
	}
	return nil
}

func DeleteLock(key string) {
	rkey := fmt.Sprintf("%s-lock", key)
	globals.Redis.Remove(rkey)
}

func HandleRequest(d []byte) error {
	var entry JobEntry
	json.Unmarshal(d, &entry)

	switch entry.Action {
	case "get":
		data, err := globals.MYSQL.Get(entry.PersonID) // Mysql get the data
		if err != nil {
			return err
		}
		rKey := fmt.Sprintf("%s-row", entry.PersonID) // Caching data key
		globals.Redis.Add(rKey, data)                 // Caching to redis
		globals.Redis.Add(
			entry.JobID,
			ActionStatus{"completed", rKey},
		)
	case "add":
		err := globals.MYSQL.Add(d) // Mysql Add
		if err != nil {
			return err
		}
		rKey := fmt.Sprintf("%s-row", entry.PersonID) // Caching data key
		globals.Redis.Add(rKey, entry)                // Caching to redis
		globals.Redis.Add(
			entry.JobID,
			ActionStatus{"completed", rKey},
		)
	case "delete":
		WaitForLock(entry.PersonID)       // Wait for the lock
		err := CreateLock(entry.PersonID) // Create Lock
		if err != nil {
			return err
		}
		rKey := fmt.Sprintf("%s-row", entry.PersonID) // Cached key
		globals.Redis.Remove(rKey)                    // Remove the cached key
		err = globals.MYSQL.Remove(entry.PersonID)    // Mysql Remove
		if err != nil {
			return err
		}
		DeleteLock(entry.PersonID) // Remove lock
		globals.Redis.Add(
			entry.JobID,
			ActionStatus{"completed", ""},
		)
	case "update":
		WaitForLock(entry.PersonID)       // Wait for the lock
		err := CreateLock(entry.PersonID) // Create Lock
		if err != nil {
			return err
		}
		rKey := fmt.Sprintf("%s-row", entry.PersonID) // Cached key
		globals.Redis.Remove(rKey)                    // Remove the cached key
		globals.MYSQL.Update(d)                       // Mysql Update
		DeleteLock(entry.PersonID)                    // Remove lock
		globals.Redis.Add(
			entry.JobID,
			ActionStatus{"completed", ""},
		)
	}
	return nil
}
