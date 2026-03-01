package portfolio

import "time"

type StockTerm string

const (
	TermShort StockTerm = "short_term"
	TermMid   StockTerm = "mid_term"
	TermLong  StockTerm = "long_term"
)

type Stock struct {
	ID          string
	UserID      string
	StockName   string
	TargetPrice float64
	StopLoss    float64
	Term        StockTerm
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
