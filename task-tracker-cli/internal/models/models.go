package models

type Tasks struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
