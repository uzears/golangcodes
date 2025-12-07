package storage

import (
	"os"
	"task-tracker-cli/internal/models"

	"encoding/json"
)

func SaveTasks(tasks []models.Tasks) error {

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
