package taskhandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"todoList/dbhandler"
)

// TodoListRunner the entry point for the todo list, handles the event loop
// for the todo list. Takes a file path to the required database file.
func TodoListRunner(filePath string) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Fetch relevant entries from db, and show menu
		dbOutput := dbhandler.GetDbRows(filePath)
		fmt.Println(len(dbOutput))
		showMenu()
		scanner.Scan()
		text := scanner.Text()

		// Switch statement for initialising different menu options.
		switch text {
		case "1":
			addTaskToList()
		case "2":
			fmt.Println(len(dbOutput))
			showTasks(&dbOutput)
		case "3":
			markAsDone()
		case "4":
			deleteTask()
		case "5":
			os.Exit(0)
		default:
			fmt.Println("Error: Invalid input")
		}

	}
}

func addTaskToList() {
	fmt.Println("adding tasks.")
}

// TODO: fix so that the output is correct.
func showTasks(tasks *[]dbhandler.DbRow) {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("| %-25s | %-10s | %-10s |\n", "Task Name", "Status", "Date")
	fmt.Println(strings.Repeat("=", 50))

	for _, task := range *tasks {
		fmt.Printf("| %-25s | %-10s | %-10s |\n", task.Name, task.Status, task.Date)
	}

	fmt.Println(strings.Repeat("=", 50))
}

func markAsDone() {
	fmt.Println("mark task as done")
}

func deleteTask() {
	fmt.Println("delete tasks")
}

// showMenu shows the output menu options for the todo list.
func showMenu() {
	menu := `====================================
          📌 TO-DO LIST 📌         
====================================
[1] ➜ Add a new task
[2] ➜ View tasks
[3] ➜ Mark task as done ✅
[4] ➜ Delete a task ❌
[5] ➜ Exit 🚪
------------------------------------
Enter your choice: `
	fmt.Println(menu)
}
