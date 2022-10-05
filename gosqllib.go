package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"strings"
)

var db *sql.DB

type SqlManager interface {
	OpenConnection()
	InsertInto()
	ExecuteQuery()
}

type Record struct {
	IP     string `db:"ip"`
	Status string `db:"status"`
}

type Records []Record

// Object
type GoSqlManager struct {
	myDb     *sql.DB
	error    *error
	server   string
	database string
}

// Create a connection
func (sqlManager GoSqlManager) OpenConnection(server string, database string) {
	sqlManager.server = server
	sqlManager.database = database

	connString := fmt.Sprintf("server=%s;port=1433;trusted_connection=yes;", server)
	thisDb, err := sql.Open("mssql", connString)

	sqlManager.myDb = thisDb
	db = thisDb
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connected!\n")
}

// Insert records
// db.QueryContext(ctx, `select * from t where ID = @ID and Name = @p2;`, sql.Named("ID", 6), "Bob")
func (sqlManager GoSqlManager) InsertInto(table string, values map[string]string) {
	columns := ""
	valuesStr := ""
	for k, v := range values {
		columns += "[" + k + "]" + ", "
		valuesStr += v + "', '"
	}
	columns = columns[:len(columns)-2]
	valuesStr = valuesStr[:len(valuesStr)-3]

	query := "INSERT INTO " + table + " (" + columns + ") VALUES ('" + valuesStr + ");"
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	for rows.Next() {
		var id int
		var txt string
		err := rows.Scan(&id, &txt)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, txt)
	}
}

// Retrieve query results
func (sqlManager GoSqlManager) ExecuteQuery(query string) (*sql.Rows, error) {

	rows, err2 := db.Query(query)
	if err2 != nil {
		log.Fatal("Error executing query: " + err2.Error())
		return rows, err2
	} else {
		log.Println("Query executed successfully.")
	}
	return rows, err2
}

func (sqlManager GoSqlManager) BulkInsert(cs Records) (sql.Result, error) {
	var (
		placeholders []string
		vals         []interface{}
	)

	for index, contact := range cs {
		placeholders = append(placeholders, fmt.Sprintf("(?, ?)"))
		index = index*2 + 1
		index = index*2 + 2
		vals = append(vals, contact.IP, contact.Status)
	}

	insertStatement := fmt.Sprintf("INSERT INTO savvytest.dbo.HaProxy_Test(IP,Status) VALUES %s", strings.Join(placeholders, ","))

	stmt, err := db.Prepare(insertStatement)
	rows, err := stmt.Exec(vals...)
	if err != nil {
		return rows, err
	}

	return rows, err
}
