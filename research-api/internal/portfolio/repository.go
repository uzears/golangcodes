package portfolio

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	r := &repository{db: db}
	r.ensureSchema(context.Background())
	return r
}

func (r *repository) ensureSchema(ctx context.Context) {
	query := `
		CREATE TABLE IF NOT EXISTS portfolio_stocks (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			stock_name TEXT NOT NULL,
			target_price DOUBLE PRECISION NOT NULL,
			stop_loss DOUBLE PRECISION NOT NULL,
			term TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`
	_, _ = r.db.ExecContext(ctx, query)
}

func (r *repository) CreateStock(ctx context.Context, stock *Stock) error {
	query := `
		INSERT INTO portfolio_stocks
			(id, user_id, stock_name, target_price, stop_loss, term, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		stock.ID,
		stock.UserID,
		stock.StockName,
		stock.TargetPrice,
		stock.StopLoss,
		stock.Term,
		stock.CreatedAt,
		stock.UpdatedAt,
	)

	return err
}

func (r *repository) ListStocksByUser(ctx context.Context, userID string) ([]Stock, error) {
	query := `
		SELECT id, user_id, stock_name, target_price, stop_loss, term, created_at, updated_at
		FROM portfolio_stocks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stocks := make([]Stock, 0)
	for rows.Next() {
		var stock Stock
		if err := rows.Scan(
			&stock.ID,
			&stock.UserID,
			&stock.StockName,
			&stock.TargetPrice,
			&stock.StopLoss,
			&stock.Term,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}
