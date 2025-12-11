package storage

import (
	"task-tracker-cli/internal/models"
)

type Storage interface {
	LoadTask() ([]models.Tasks, error)
	SaveTasks([]models.Tasks) error
	UpdateTasks(int, string) error
	//DeleteTask()
}
