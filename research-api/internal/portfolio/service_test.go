package portfolio

import (
	"context"
	"testing"

	"github.com/uzears/golangcodes/research-api/internal/platform/logger"
)

type fakeRepo struct {
	createStockErr error
	listErr        error
	created        *Stock
	list           []Stock
}

type fakeLogger struct {
	debugCalls int
}

func (l *fakeLogger) Info(msg string, fields ...any) {}
func (l *fakeLogger) Error(msg string, fields ...any) {}
func (l *fakeLogger) Warn(msg string, fields ...any) {}
func (l *fakeLogger) With(fields ...any) logger.Logger { return l }
func (l *fakeLogger) Debug(msg string, fields ...any)  { l.debugCalls++ }

func (r *fakeRepo) CreateStock(ctx context.Context, stock *Stock) error {
	r.created = stock
	return r.createStockErr
}

func (r *fakeRepo) ListStocksByUser(ctx context.Context, userID string) ([]Stock, error) {
	return r.list, r.listErr
}

func TestCreateStockSuccess(t *testing.T) {
	repo := &fakeRepo{}
	log := &fakeLogger{}
	svc := NewService(repo, log)

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
	if log.debugCalls == 0 {
		t.Fatalf("expected debug log for generated id")
	}
}

func TestCreateStockInvalidTerm(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo, &fakeLogger{})

	if _, err := svc.CreateStock(context.Background(), "u1", "AAPL", 200, 150, StockTerm("weekly")); err == nil {
		t.Fatalf("expected invalid term error")
	}
}

func TestCreateStockInvalidRisk(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo, &fakeLogger{})

	if _, err := svc.CreateStock(context.Background(), "u1", "AAPL", 100, 120, TermShort); err == nil {
		t.Fatalf("expected stop loss validation error")
	}
}
