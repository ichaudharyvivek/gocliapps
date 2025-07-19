package main

import (
	"flag"
	"fmt"
	"gotask/tasks"
	"os"
	"text/tabwriter"
)

var (
	gotask  = flag.NewFlagSet("gotask", flag.ExitOnError)
	addCmd  = flag.NewFlagSet("add", flag.ExitOnError)
	listCmd = flag.NewFlagSet("list", flag.ExitOnError)
	delCmd  = flag.NewFlagSet("delete", flag.ExitOnError)
	doneCmd = flag.NewFlagSet("done", flag.ExitOnError)
)

func main() {
	gotask.Usage = usage
	if len(os.Args) < 2 {
		fmt.Println("Expected gotask command")
		os.Exit(1)
	}

	service := tasks.NewTaskService("tasks.json")
	switch os.Args[1] {
	case "gotask":
		gotask.Parse(os.Args[2:])
		handleGoTask(service, gotask.Args())
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}

}

func handleGoTask(service *tasks.TaskService, args []string) {
	if len(args) < 1 {
		fmt.Println("Expected subcommands. Eg. add, list, delete")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		task := addCmd.String("task", "", "The description of todo")
		addCmd.Parse(args[1:])

		if *task == "" {
			fmt.Println("Error: task description cannot be empty")
			os.Exit(1)
		}

		t, err := service.AddTask(*task)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Task added successfully: %d\n", t.ID)

	case "list":
		listCmd.Parse(args[1:])
		tasks, err := service.ListTasks()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer w.Flush()

		fmt.Fprintf(w, "%-5s\t%-20s\t%-10s\n", "ID", "Description", "Status")
		fmt.Fprintf(w, "%s\t%s\t%s\n", "-----", "--------------------", "----------")

		for _, task := range tasks {
			fmt.Fprintf(w, "%-5d\t%-20s\t%-10s\n", task.ID, task.Description, task.Status)
		}

	case "delete":
		id := delCmd.Int("id", -1, "Id of todo to delete")
		delCmd.Parse(args[1:])

		if *id == -1 {
			fmt.Println("Please provide a valid id")
			os.Exit(1)
		}

		if err := service.DeleteTask(*id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Task deleted successfully: %d\n", *id)

	case "done":
		id := doneCmd.Int("id", -1, "Id of todo to mark as complete")
		doneCmd.Parse(args[1:])

		if *id == -1 {
			fmt.Println("Please provide a valid id")
			os.Exit(1)
		}

		if err := service.CompleteTask(*id); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Task marked as completed: %d\n", *id)

	case "help":
		gotask.Usage()

	default:
		fmt.Println("Invalid sub-command")
	}

}

func usage() {
	gotask.PrintDefaults()
	fmt.Print(cmd)
}

const cmd = `
Usage:
	gotask [command] [flag]

Available Commands:
  	add       Add a new task to the task list
  	list      List all tasks
  	done      Mark a task as completed
  	delete    Delete a task from the task list
  	help      Display help window

Available Flags:
	--task    Todo task
	--id      Todo id    
`
