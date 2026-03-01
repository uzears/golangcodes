package portfolio

import "time"

type StockResponse struct {
	ID          string    `json:"id"`
	StockName   string    `json:"stock_name"`
	TargetPrice float64   `json:"target_price"`
	StopLoss    float64   `json:"stop_loss"`
	Term        string    `json:"term"`
	CreatedAt   time.Time `json:"created_at"`
}
