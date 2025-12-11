package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"task-tracker-cli/internal/models"
)

const file = "tasks.json"

type FileStorage struct {
	File string
}

func (fs *FileStorage) LoadTask() ([]models.Tasks, error) {

	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Tasks{}, nil
		}
	}

	var tasks []models.Tasks
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return []models.Tasks{}, nil
	}

	return tasks, nil
}

func (fs *FileStorage) SaveTasks(tasks []models.Tasks) error {

	for _, n := range tasks {
		fmt.Printf("save task :%d: :%s:", n.Id, n.Status)

	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) UpdateTasks(id int, status string) error {
	fmt.Println("Update Task")
	tasks, err := f.LoadTask()
	if err != nil {
		return err
	}

	for i, n := range tasks {
		if n.Id == id {
			tasks[i].Status = "done"
		}
	}
	err = f.SaveTasks(tasks)
	if err != nil {
		return err
	}

	return nil
}

/*func (f *FileStorage) DeleteTaskById(id int)*/
