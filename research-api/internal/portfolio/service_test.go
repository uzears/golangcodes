package portfolio

import (
	"context"
	"testing"
)

type fakeRepo struct {
	createStockErr error
	listErr        error
	created        *Stock
	list           []Stock
}

func (r *fakeRepo) CreateStock(ctx context.Context, stock *Stock) error {
	r.created = stock
	return r.createStockErr
}

func (r *fakeRepo) ListStocksByUser(ctx context.Context, userID string) ([]Stock, error) {
	return r.list, r.listErr
}

func TestCreateStockSuccess(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	stock, err := svc.CreateStock(context.Background(), "u1", "AAPL", 210.5, 180.0, TermMid)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stock == nil {
		t.Fatalf("expected stock to be returned")
	}
	if repo.created == nil {
		t.Fatalf("expected repo create to be called")
	}
	if stock.StockName != "AAPL" {
		t.Fatalf("unexpected stock name: %s", stock.StockName)
	}
}

func TestCreateStockInvalidTerm(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	if _, err := svc.CreateStock(context.Background(), "u1", "AAPL", 200, 150, StockTerm("weekly")); err == nil {
		t.Fatalf("expected invalid term error")
	}
}

func TestCreateStockInvalidRisk(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	if _, err := svc.CreateStock(context.Background(), "u1", "AAPL", 100, 120, TermShort); err == nil {
		t.Fatalf("expected stop loss validation error")
	}
}
