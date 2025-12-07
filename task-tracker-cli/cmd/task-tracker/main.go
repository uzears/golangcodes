package main

import (
	"fmt"
	"task-tracker-cli/internal/services"
)

func main() {

	fmt.Println("Task tracker cli\nChoose Option")
	actn := &services.TaskService{}
	err := actn.ActnSelector()
	if err != nil {
		return
	}

}
