package main

import (
	"go-reporting/task"
	"os"
)

func main() {
	taskName := "gdocs-contract"

	if envTaskName := os.Getenv("TASK"); envTaskName != "" {
		taskName = envTaskName
	}

	if taskName == "gdocs-contract" {
		task.GdocsExport()
	} else if taskName == "gsheet-contract" {
		println(taskName)
		// task.GdocsExport()
	} else {
		task.GdocsExport()
	}
}
