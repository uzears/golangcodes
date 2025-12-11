package main

import (
	"fmt"
	"task-tracker-cli/internal/services"
	"task-tracker-cli/internal/storage"
)

func main() {

	//fmt.Printf("Start time :%d:"time.Now())
	fmt.Println("Task tracker cli\nChoose Option")
	/* by using file */
	fileStore := &storage.FileStorage{
		File: "tasks.json",
	}
	actn := &services.TaskServiceImpl{
		Store: fileStore,
	}

	err := actn.ActnSelector()
	if err != nil {
		return
	}

}
