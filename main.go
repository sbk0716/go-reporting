package main

import (
	"go-reporting/task"
	"os"
)

func main() {
	taskName := "gsheets-bank-statement-export"
	if envTaskName := os.Getenv("TASK"); envTaskName != "" {
		taskName = envTaskName
	}
	if taskName == "gdocs-contract-export" {
		task.GdocsContractExport()
	} else if taskName == "gsheets-contract-export" {
		task.GsheetsContractExport()
	} else if taskName == "gsheets-bank-statement-export" {
		task.GsheetsBankStatementExport()
	} else {
		task.GdocsContractExport()
	}
}
