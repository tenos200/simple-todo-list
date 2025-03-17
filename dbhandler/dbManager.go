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

// UpdateDatabase updates the database with the tasks list collected from the
// taskManager. All the values that have an id of 0 will be inserted into the
// database as these did not exists previously, all the rest of the rows will
// update the status row. After this we carry out a query which removes all
// rows that has the updated status.
func UpdateDatabase(filePath string, tasks *[]DbRow) {

	db, openDbErr := sql.Open("sqlite3", filePath)

	if openDbErr != nil {
		fmt.Println(openDbErr)
	}
	// for loop to update all the task information
	for _, v := range *tasks {
		// if the id is 0 we know its a new entry otherwise we just update the row
		if v.Id == 0 {
			_, insertErr := db.Exec(`INSERT INTO todolist(task_name, task_status, task_date)
            VALUES(?, ?, ?);`, v.Name, v.Status, v.Date)
			if insertErr != nil {
				fmt.Println(insertErr)
			}
		} else {
			// If an entry has been moved to complete we need to change this.
			_, updateErr := db.Exec(`UPDATE todolist SET task_status = ?
            WHERE number_id = ?;`, v.Status, v.Id)
			if updateErr != nil {
				fmt.Println(updateErr)
			}
		}
	}
	// Finally removes all the tasks that have been completed
	_, deleteErr := db.Exec(
		`DELETE FROM todolist WHERE task_status = 'Complete';`)
	if deleteErr != nil {
		fmt.Println(deleteErr)
	}
}

// CreateSchema takes a file path for the database and creates the required
// table on that database.
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
