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

type TaskServiceImpl struct {
	Store storage.Storage
}

/* not of use
type TaskService interface {
	ActnSelector() error
	AddTask() error
	DeleteTask() error
	ListTasks() error
}*/

func (s *TaskServiceImpl) ActnSelector() error {

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
		err := s.AddTask()
		if err != nil {
			fmt.Println("Error while adding task :", err)
			return err
		}
		return nil

	case "l":
		fmt.Println("Inside List Task")
		err := s.DisplayTask()
		if err != nil {
			return err
		}
		return nil

	case "d":
		fmt.Println("Inside Delete Task")
		err := s.DeleteTask()
		if err != nil {
			return err
		}

	case "e":
		fmt.Println("Exiting ...")
		return nil

	case "m":
		err := s.MarkTaskDone()
		if err != nil {
			return err
		}

	default:
		fmt.Println("Wrong Choice")
		return nil
	}

	return nil
}

func (s *TaskServiceImpl) AddTask() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the task")

	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)

	data, err := s.Store.LoadTask()
	if err != nil {
		return err
	}

	newTask := models.Tasks{
		Id:          len(data) + 1,
		Description: desc,
		Status:      "pending",
	}

	data = append(data, newTask)

	err = s.Store.SaveTasks(data)
	if err != nil {
		return err
	}
	othActn := TaskServiceImpl{}
	err = othActn.OtherActn()
	if err != nil {
		return err
	}

	fmt.Println("Successfully Added")
	return nil
}

func (s *TaskServiceImpl) DisplayTask() error {
	fmt.Println("Displaying Task :")
	tasks, err := s.Store.LoadTask()
	if err != nil {
		return err
	}

	fmt.Println("Tasks :")
	for _, n := range tasks {
		fmt.Printf("ID :%d: Description :%s: Status :%s:\n", n.Id, n.Description, n.Status)
	}
	err = s.OtherActn()
	if err != nil {
		return err
	}

	return nil
}

func (d *TaskServiceImpl) DeleteTask() error {

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

	deleteData, err := d.Store.LoadTask()
	if err != nil {
		return err
	}

	for _, i := range deleteData {
		if i.Id == deleteId {
			continue
		}
		newData = append(newData, i)
	}

	err = d.Store.SaveTasks(newData)
	if err != nil {
		return err
	}

	fmt.Println("Task Delete Successfully ")

	err = d.OtherActn()
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskServiceImpl) MarkTaskDone() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the task id to mark as done")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	markId, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid Id")
		return err
	}
	err = s.Store.UpdateTasks(markId, "done")
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskServiceImpl) OtherActn() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Do you want to Repeat Action?")
	actn, _ := reader.ReadString('\n')
	actn = strings.TrimSpace(actn)
	actn = strings.ToLower(actn)

	if actn == "y" {
		err := s.ActnSelector()
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	return nil
}
