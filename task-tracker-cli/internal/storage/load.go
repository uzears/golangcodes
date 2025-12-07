package storage

import (
	"encoding/json"
	"os"
	"task-tracker-cli/internal/models"
)

const file = "tasks.json"

func LoadTask() ([]models.Tasks, error) {

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
