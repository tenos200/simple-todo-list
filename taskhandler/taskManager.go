package taskhandler

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"strings"
	"time"
	"todoList/dbhandler"
)

// TodoListRunner the takes a file path to required database file,
// handles the event loop. Provides the user with 5 choices for input:
//
// Picks:
//  1. Add a task
//  2. Show all tasks
//  3. Show all tasks
//  4. Mark a task as done
//  5. Exit program
func TodoListRunner(filePath string) {

	// Fetch relevant entries from db, and show menu
	dbOutputCache := dbhandler.GetDbRows(filePath)
	for {
		choice := showMenu()

		// Switch statement for initialising different menu options.
		switch choice {
		case 1:
			showTasks(&dbOutputCache)
		case 2:
			addTaskToList(&dbOutputCache)
		case 3:
			markAsDone(&dbOutputCache)
		case 4:
			deleteTask()
		case 5:
			os.Exit(0)
		default:
			fmt.Println("Error: Invalid input")
		}

	}
}

func addTaskToList(dbOutputCache *[]dbhandler.DbRow) *[]dbhandler.DbRow {
	var rowToAdd dbhandler.DbRow

	// For loop for adding task to list
	for {
		var confirmation int

		// Add name to task.
		huh.NewInput().
			Title("Enter task name.").
			Prompt("> ").
			Value(&rowToAdd.Name).
			Run()

		// Get the date for the task
		rowToAdd.Date = getInputDate()

		// Set default status to imcomplete
		rowToAdd.Status = "Incomplete"
		formattedDate := fmt.Sprintf("Task: %s\nDate: %s\n[Y/N]",
			rowToAdd.Name,
			rowToAdd.Date)

		huh.NewSelect[int]().
			Title(formattedDate).
			Options(
				huh.NewOption("Y", 1),
				huh.NewOption("N", 2),
			).
			Value(&confirmation).
			Run()

		if confirmation == 1 {
			*dbOutputCache = append(*dbOutputCache, rowToAdd)
			break
		}

	}
	return dbOutputCache
}

func getInputDate() string {

	var pickedDate string
	var switchSelect int
	currentDate := time.Now()
	huh.NewSelect[int]().
		Title("Input a date for task completion.").
		Options(
			huh.NewOption("Tomorrow.", 1),
			huh.NewOption("In 7 days.", 2),
			huh.NewOption("Enter own date.", 3),
		).
		Value(&switchSelect).
		Run()

		// Switch statement to pick the correct date
	switch switchSelect {
	case 1:
		tomorrow := currentDate.Add(24 * time.Hour)
		pickedDate = tomorrow.Format("2006-01-02")
	case 2:
		nextWeek := currentDate.Add(time.Hour * 24 * 7)
		pickedDate = nextWeek.Format("2006-01-02")
	case 3:
		// Initialise a new scanner and take the date.
		var dateInput string
		huh.NewInput().
			Title("Enter date formatted as '2006-01-02'").
			Prompt("> ").
			Value(&dateInput).
			Run()

		parsedDate, dateErr := time.Parse("2006-01-02", dateInput)
		if dateErr != nil {
			fmt.Println("Incorrectly formatted date", "date", dateErr)
		}
		pickedDate = parsedDate.Format("2006-01-02")
	}
	return pickedDate
}

// showTasks displays all tasks for tasklist.
func showTasks(tasks *[]dbhandler.DbRow) {

	if len(*tasks) == 0 {
		fmt.Println("No tasks have been added.")
	} else {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("| %-25s | %-10s | %-10s |\n", "Task Name", "Status", "Date")
		fmt.Println(strings.Repeat("=", 50))

		for _, task := range *tasks {
			fmt.Printf("| %-25s | %-10s | %-10s |\n",
				task.Name, task.Status, task.Date)
		}

		fmt.Println(strings.Repeat("=", 50))
	}
}

func markAsDone(tasks *[]dbhandler.DbRow) {
}

func deleteTask() {
	fmt.Println("delete tasks")
}

func showMenu() int {
	var userInput int
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("View Tasks", 1),
					huh.NewOption("Add New Task", 2),
					huh.NewOption("Mark Complete", 3),
					huh.NewOption("Delete Task", 4),
					huh.NewOption("Exit", 5),
				).Value(&userInput),
		),
	).Run()
	return userInput
}
