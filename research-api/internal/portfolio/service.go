package portfolio

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateStock(ctx context.Context, userID, stockName string, targetPrice, stopLoss float64, term StockTerm) (*Stock, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if stockName == "" {
		return nil, errors.New("stock name is required")
	}
	if targetPrice <= 0 {
		return nil, errors.New("target price must be greater than 0")
	}
	if stopLoss <= 0 {
		return nil, errors.New("stop loss must be greater than 0")
	}
	if stopLoss >= targetPrice {
		return nil, errors.New("stop loss must be less than target price")
	}
	if !isValidTerm(term) {
		return nil, errors.New("term must be one of: short_term, mid_term, long_term")
	}

	now := time.Now().UTC()
	stock := &Stock{
		ID:          newID(),
		UserID:      userID,
		StockName:   stockName,
		TargetPrice: targetPrice,
		StopLoss:    stopLoss,
		Term:        term,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateStock(ctx, stock); err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *service) ListStocksByUser(ctx context.Context, userID string) ([]Stock, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	return s.repo.ListStocksByUser(ctx, userID)
}

func isValidTerm(term StockTerm) bool {
	switch term {
	case TermShort, TermMid, TermLong:
		return true
	default:
		return false
	}
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(b)
}
