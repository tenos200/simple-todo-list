package main

import (
	"fmt"
	"log"
	"os"
	"todoList/dbhandler"
	"todoList/taskhandler"
	"todoList/tui"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

const filePath = "./test.db"

func main() {

	// if we can't open the file then we should create the schema for the db.
	file, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println(openErr)
		createDbFile()
		dbhandler.CreateSchema(filePath)
	}

	// If it already exist, ensure to close.
	defer file.Close()

	// Start the task
	if false {
		taskhandler.TodoListRunner(filePath)
	} else {
		dbOutputCache := dbhandler.GetDbRows(filePath)
		p := tea.NewProgram(tui.InitialModel(dbOutputCache),
			tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Theres been an error %v", err)
		}
	}
	os.Exit(0)
}

// createDbFile creates a database on file should it not exist.
func createDbFile() {
	_, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Something went wrong when creating file", "err", err)
	}
}
