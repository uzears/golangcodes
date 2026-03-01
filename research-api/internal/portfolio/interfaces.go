package portfolio

import "context"

type Service interface {
	CreateStock(ctx context.Context, userID, stockName string, targetPrice, stopLoss float64, term StockTerm) (*Stock, error)
	ListStocksByUser(ctx context.Context, userID string) ([]Stock, error)
}

type Repository interface {
	CreateStock(ctx context.Context, stock *Stock) error
	ListStocksByUser(ctx context.Context, userID string) ([]Stock, error)
}
