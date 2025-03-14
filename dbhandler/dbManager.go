package dbhandler

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DbRow struct {
	Id     int
	Name   string
	Status string
	Date   string
}

func InserIntoDb() {
	fmt.Println("Inserted into database")
}

func UpdateDb() {
	fmt.Println("Updated Database Entry")
}

func DeleteFromDb() {
	fmt.Println("Deleted from database.")
}

func CreateSchema(filePath string) {
	db, openDbErr := sql.Open("sqlite3", filePath)
	if openDbErr != nil {
		fmt.Println(openDbErr)
	}
	_, err := db.Exec(`CREATE TABLE TodoList(
    number_id INTEGER PRIMARY KEY,
    task_name TEXT NOT NULL, 
    task_status TEXT NOT NULL, 
    task_date TEXT NOT NULL);`)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Created schema")
}

// GetDbRows takes a filepath to a database and returns a slice with
// all database rows.
//
// Example:
//
//	dbRow := GetDbRows("database.db")
//	for i := 0; i < 3; i++ {
//	    fmt.Printf("%d ",dbRow[i].Id)
//	}
//
// Output:
//
//	1 2 3
func GetDbRows(filePath string) []DbRow {
	var dbOutput []DbRow
	db, openDbErr := sql.Open("sqlite3", filePath)
	if openDbErr != nil {
		log.Fatal(openDbErr)
	}
	rows, fetchErr := db.Query(`SELECT * FROM TodoList;`)
	if fetchErr != nil {
		log.Fatal(fetchErr)
	}

	for rows.Next() {
		var row DbRow
		rows.Scan(&row.Id, &row.Name, &row.Status, &row.Date)
		dbOutput = append(dbOutput, row)
	}

	defer rows.Close()
	return dbOutput
}
