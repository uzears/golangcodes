package services

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task-tracker-cli/internal/models"
	"task-tracker-cli/internal/storage"
)

type TaskService struct{}

func (s *TaskService) ActnSelector() error {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Choose Option")
	fmt.Println("A. Add task")
	fmt.Println("L. List all task")
	fmt.Println("M. Mark as done")
	fmt.Println("D. Delete task")
	fmt.Println("S. List by status")
	fmt.Println("E. Exit")
	fmt.Println("")

	taskType, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Unable to read ")
		return err
	}
	taskType = strings.TrimSpace(taskType)
	taskType = strings.ToLower(taskType)

	switch taskType {
	case "a":
		fmt.Println("Inside Add Task")
		add := TaskService{}
		err := add.AddTask()
		if err != nil {
			fmt.Println("Error while adding task :", err)
			return err
		}
		return nil

	case "l":
		fmt.Println("Inside List Task")
		display := TaskService{}
		err := display.DisplayTask()
		if err != nil {
			return err
		}
		return nil

	case "d":
		fmt.Println("Inside Delete Task")
		delete := TaskService{}
		err := delete.DeleteTask()
		if err != nil {
			return err
		}

	case "e":
		fmt.Println("Exiting ...")
		return nil

	default:
		fmt.Println("Wrong Choice")
		return nil
	}

	return nil
}

func (s *TaskService) AddTask() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the task")

	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)

	data, err := storage.LoadTask()
	if err != nil {
		return err
	}

	newTask := models.Tasks{
		Id:          len(data) + 1,
		Description: desc,
		Status:      "pending",
	}

	data = append(data, newTask)

	err = storage.SaveTasks(data)
	if err != nil {
		return err
	}
	othActn := TaskService{}
	err = othActn.OtherActn()
	if err != nil {
		return err
	}

	fmt.Println("Successfully Added")
	return nil
}

func (s *TaskService) DisplayTask() error {
	fmt.Println("Displaying Task :")
	tasks, err := storage.LoadTask()
	if err != nil {
		return err
	}

	fmt.Println("Tasks :")
	for _, n := range tasks {
		fmt.Printf("ID :%d: Description :%s: Status :%s:\n", n.Id, n.Description, n.Status)
	}
	othActn := TaskService{}
	err = othActn.OtherActn()
	if err != nil {
		return err
	}

	return nil
}

func (d *TaskService) DeleteTask() error {

	var newData []models.Tasks

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the task Id you want to delete")
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)
	deleteId, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid Id")
		return err
	}

	deleteData, err := storage.LoadTask()
	if err != nil {
		return err
	}

	for _, i := range deleteData {
		if i.Id == deleteId {
			continue
		}
		newData = append(newData, i)
	}

	err = storage.SaveTasks(newData)
	if err != nil {
		return err
	}

	fmt.Println("Task Delete Successfully ")

	othActn := TaskService{}
	err = othActn.OtherActn()
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) OtherActn() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Do you want to Repeat Action?")
	actn, _ := reader.ReadString('\n')
	actn = strings.TrimSpace(actn)
	actn = strings.ToLower(actn)

	if actn == "y" {
		newActn := TaskService{}
		err := newActn.ActnSelector()
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	return nil
}
