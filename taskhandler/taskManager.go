package taskhandler

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"time"
	"todoList/dbhandler"
)

const timeFormat = "2006-01-02"

// TodoListRunner the takes a file path to required database file,
// handles the event loop. Provides the user with 5 choices for input:
//
// Picks:
//  1. Show all tasks
//  2. Add Task
//  3. Mark Task Done
//  4. Exit program
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
			// TODO: something is not right with return here.
			dbOutputCache = markAsDone(dbOutputCache)
			fmt.Println(dbOutputCache)
		case 4:
			// Insert and update all necessary
			dbhandler.UpdateDatabase(filePath, &dbOutputCache)
			os.Exit(0)
		}
	}
}

func addTaskToList(dbOutputCache *[]dbhandler.DbRow) *[]dbhandler.DbRow {
	var rowToAdd dbhandler.DbRow

	// For loop for adding task to list
	for {
		var confirmation bool

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
		formattedDate := fmt.Sprintf("Task: %s\nDate: %s\n",
			rowToAdd.Name,
			rowToAdd.Date)

		huh.NewConfirm().
			Title(formattedDate).
			Affirmative("Yes.").
			Negative("No.").
			Value(&confirmation).
			Run()

		if confirmation {
			*dbOutputCache = append(*dbOutputCache, rowToAdd)
			break
		}

	}
	return dbOutputCache
}

// getInputDate allows user to pick an input date from menu, returns a date string
// formatted as 2006-01-02.
func getInputDate() string {

	// Instantiate variables for huh.NewOptions input values
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
		pickedDate = tomorrow.Format(timeFormat)
	case 2:
		nextWeek := currentDate.Add(time.Hour * 24 * 7)
		pickedDate = nextWeek.Format(timeFormat)
	case 3:

		// Initialise a new scanner and take the date.
		var dateInput string
		huh.NewInput().
			Title("Enter date formatted as '2006-01-02'").
			Prompt("> ").
			Value(&dateInput).
			Run()

		parsedDate, dateErr := time.Parse(timeFormat, dateInput)
		if dateErr != nil {
			fmt.Println("Incorrectly formatted date", "date", dateErr)
		}
		pickedDate = parsedDate.Format(timeFormat)
	}
	return pickedDate
}

// showTasks displays all tasks for tasklist.
func showTasks(tasks *[]dbhandler.DbRow) {

	if len(*tasks) == 0 {
		fmt.Println("There are no tasks on the list.")
	} else {
		for _, v := range *tasks {
			// Show the task if its not complete
			if v.Status != "Complete" {
				fmt.Printf("ID: %d\nName: %s\nStatus: %s\nDate: %s\n",
					v.Id, v.Name, v.Status, v.Date)
			}
		}
	}
}

// TODO: Look back at this function and try to understand what is going on
// with the function passing here. Something with how slices are passed here
// doesn't make sense in my mind.
// TODO: Document this function...
func markAsDone(tasks []dbhandler.DbRow) []dbhandler.DbRow {
	var selectedTasks []*dbhandler.DbRow
	var choices []huh.Option[*dbhandler.DbRow]

	// Convert the tasks into huhNewOptions
	for i := range tasks {
		choices = append(choices, huh.NewOption(tasks[i].Name, &tasks[i]))
	}
	if len(tasks) == 0 {
		fmt.Println("There are no tasks on the list.")
		return tasks
	} else {
		huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[*dbhandler.DbRow]().
					Title("Task").
					Options(choices...).
					Value(&selectedTasks),
			),
		).Run()
	}
	for _, v := range selectedTasks {
		v.Status = "Complete"
	}
	return tasks
}

// showMenu displays a menu with options from 1-5, returns int from user input.
//
// Choices:
//  1. View Tasks
//  2. Add a new Task
//  3. Mark Complete
//  4. Exit
func showMenu() int {
	var userInput int
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("View Tasks", 1),
					huh.NewOption("Add New Task", 2),
					huh.NewOption("Mark Complete", 3),
					huh.NewOption("Exit", 4),
				).
				Value(&userInput),
		),
	).WithHeight(10).Run()
	return userInput
}
