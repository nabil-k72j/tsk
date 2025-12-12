package main

import (
	"fmt"
	"os"
)

var TasksPath string

func main() {

	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error getting home path", err)
		os.Exit(1)
	}

	TasksPath = homePath + "/tasks.csv"

	tasks, err := GetTasks(TasksPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading tasks", err)
	}

	Run(tasks)
}
