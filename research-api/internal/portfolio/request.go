package portfolio

type CreateStockRequest struct {
	StockName   string  `json:"stock_name"`
	TargetPrice float64 `json:"target_price"`
	StopLoss    float64 `json:"stop_loss"`
	Term        string  `json:"term"`
}
