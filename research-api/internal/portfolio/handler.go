package portfolio

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uzears/golangcodes/research-api/internal/platform/logger"
)

type Handler struct {
	service Service
	log     logger.Logger
}

func NewHandler(service Service, log logger.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) CreateStock(c *gin.Context) {
	log := c.MustGet("logger").(logger.Logger)
	log.Debug("create portfolio stock request received")

	userID, exists := c.Get("user_id")
	if !exists {
		log.Debug("create portfolio stock missing user_id in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid create stock request", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	stock, err := h.service.CreateStock(
		c.Request.Context(),
		userID.(string),
		strings.TrimSpace(req.StockName),
		req.TargetPrice,
		req.StopLoss,
		StockTerm(strings.TrimSpace(req.Term)),
	)
	if err != nil {
		log.Warn("failed to create portfolio stock", "err", err, "user_id", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Debug("portfolio stock created", "user_id", userID, "stock_id", stock.ID)
	c.JSON(http.StatusCreated, toResponse(*stock))
}

func (h *Handler) ListStocks(c *gin.Context) {
	log := c.MustGet("logger").(logger.Logger)
	log.Debug("list portfolio stocks request received")

	userID, exists := c.Get("user_id")
	if !exists {
		log.Debug("list portfolio stocks missing user_id in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stocks, err := h.service.ListStocksByUser(c.Request.Context(), userID.(string))
	if err != nil {
		log.Error("failed to list portfolio stocks", "err", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch portfolio"})
		return
	}

	resp := make([]StockResponse, 0, len(stocks))
	for _, s := range stocks {
		resp = append(resp, toResponse(s))
	}

	log.Debug("portfolio stocks fetched", "user_id", userID, "count", len(resp))
	c.JSON(http.StatusOK, gin.H{"items": resp})
}

func toResponse(stock Stock) StockResponse {
	return StockResponse{
		ID:          stock.ID,
		StockName:   stock.StockName,
		TargetPrice: stock.TargetPrice,
		StopLoss:    stock.StopLoss,
		Term:        string(stock.Term),
		CreatedAt:   stock.CreatedAt,
	}
}
