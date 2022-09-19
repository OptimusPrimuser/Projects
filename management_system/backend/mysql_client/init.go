package mysql_client

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"fmt"

	// "log"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/ziutek/mymysql/native" // Native engine
	// "github.com/ziutek/mymysql/mysql"
)

type Mysql_Entry struct {
	PersonID  string `json:"PersonID"`
	LastName  string `json:"LastName"`
	FirstName string `json:"FirstName"`
	Address   string `json:"Address"`
	City      string `json:"City"`
	Age       int    `json:"Age"`
}

type Client struct {
	Conn       *sql.DB
	DB_Name    string
	Table_Name string
}

func (c *Client) Init() error {
	user := "root"
	pass := "root"
	host := "localhost"
	port := "3306"
	db := "mgmt"
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		user, pass,
		host, port,
		db,
	)
	var err error
	c.Table_Name = "Persons"
	c.Conn, err = sql.Open(
		"mysql",
		dns,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Add(data []byte) error {
	var entry Mysql_Entry
	json.Unmarshal(data, &entry)
	// colNames := "LastName,FirstName,Address,City,AGE"
	query := fmt.Sprintf(
		"INSERT INTO %s VALUES ('%s','%s','%s','%s','%s',%d);",
		c.Table_Name,
		entry.PersonID,
		entry.LastName,
		entry.FirstName,
		entry.Address,
		entry.City,
		entry.Age,
	)
	_, err := c.Conn.Query(query)

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(key string) ([]Mysql_Entry, error) {
	query := fmt.Sprintf(
		"select * from %s where PersonID='%s';",
		c.Table_Name,
		key,
	)
	res, err := c.Conn.Query(query)
	if err != nil {
		return []Mysql_Entry{}, err
	}
	resOut := make([]Mysql_Entry, 0)
	for res.Next() {

		var entry Mysql_Entry
		err := res.Scan(
			&entry.PersonID,
			&entry.LastName,
			&entry.FirstName,
			&entry.Address,
			&entry.City,
			&entry.Age,
		)

		if err != nil {
			return []Mysql_Entry{}, err
		}
		resOut = append(resOut, entry)
	}
	res.Close()
	return resOut, nil
}

func (c *Client) Remove(key string) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE PersonID='%s';",
		c.Table_Name,
		key,
	)
	_, err := c.Conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Update(data []byte) error {
	var entry Mysql_Entry
	json.Unmarshal(data, &entry)
	// columns := []string{
	// 	"PersonID=%d",
	// 	"LastName=\"%s\"",
	// 	"FirstName=\"%s\"",
	// 	"Address=\"%s\"",
	// 	"City=\"%s\"",
	// 	"AGE=%d",
	// }
	columns := ""
	if entry.Address != "" {
		columns = columns + fmt.Sprintf(" Address='%s',", entry.Address)
	}
	if entry.Age != 0 {
		columns = columns + fmt.Sprintf(" AGE=%d,", entry.Age)
	}
	if entry.City != "" {
		columns = columns + fmt.Sprintf(" City='%s',", entry.City)
	}
	if entry.FirstName != "" {
		columns = columns + fmt.Sprintf(" FirstName='%s',", entry.FirstName)
	}
	if entry.LastName != "" {
		columns = columns + fmt.Sprintf(" LastName='%s',", entry.LastName)
	}
	columns = columns[:len(columns)-1]
	query := fmt.Sprintf(
		"UPDATE %s SET %s where PersonID='%s';",
		c.Table_Name,
		columns,
		entry.PersonID,
	)
	// query = fmt.Sprintf(
	// 	query,
	// 	entry.PersonID,
	// 	entry.LastName,
	// 	entry.FirstName,
	// 	entry.Address,
	// 	entry.City,
	// 	entry.Age,
	// )
	_, err := c.Conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}
