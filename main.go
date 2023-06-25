package main

import (
	"go-reporting/task"
	"os"
)

func main() {
	taskName := "gsheets-contract"
	if envTaskName := os.Getenv("TASK"); envTaskName != "" {
		taskName = envTaskName
	}
	if taskName == "gdocs-contract" {
		task.GdocsExport()
	} else if taskName == "gsheets-contract" {
		task.GsheetsExport()
	} else {
		task.GdocsExport()
	}
}
