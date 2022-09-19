package main

import (
	"backend/globals"
	"backend/rabbit_mq"
	"fmt"
)

func main() {
	// fmt.Println("lol")
	// c := &mysql_client.Client{}
	// // fmt.Println(c)
	// c.Init()
	// fmt.Println("ADD")
	// c.Add(
	// 	mysql_client.Mysql_Entry{
	// 		PersonID:  1,
	// 		LastName:  "LOL",
	// 		FirstName: "LOL",
	// 		Address:   "abv",
	// 		City:      "iop",
	// 		Age:       23,
	// 	},
	// )
	// fmt.Println("GET")
	// c.Get(1)
	// fmt.Println(c.Get(1))
	// fmt.Println("Update")
	// c.Update(
	// 	mysql_client.Mysql_Entry{
	// PersonID:  1,
	// LastName:  "LAL",
	// FirstName: "LAL",
	// Address:   "abv",
	// City:      "iYp",
	// Age:       23,
	// 	},
	// )
	// fmt.Println("GET")
	// fmt.Println(c.Get(1))
	// fmt.Println("delete")
	// c.Remove(1)
	err := globals.MYSQL.Init()
	if err != nil {
		fmt.Println("mysql", err)
	}
	globals.Redis.Init()
	var RMQ rabbit_mq.RMQ
	RMQ.Init("test1")
	RMQ.Consume()
	// for {
	// 	fmt.Println("lol")
	// }
}
