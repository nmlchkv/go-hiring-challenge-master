package catalog

import (
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CatalogService interface {
	GetCatalog(categoryCode string, priceLessThan *float64, offset, limit int) ([]models.Product, int64, error)
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProductSummary struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type Response struct {
	Items []ProductSummary `json:"items"`
	Total int64            `json:"total"`
}

type CatalogHandler struct {
	service CatalogService
}

func NewCatalogHandler(s CatalogService) *CatalogHandler {
	return &CatalogHandler{
		service: s,
	}
}

func parseQueryParams(r *http.Request) (offset, limit int, category string, priceLessThan *float64, errMsg string, status int) {
	q := r.URL.Query()

	offset = 0
	limit = 10

	if v := q.Get("offset"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 0 {
			return 0, 0, "", nil, "invalid offset", http.StatusBadRequest
		}
		offset = parsed
	}

	if v := q.Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, "", nil, "invalid limit", http.StatusBadRequest
		}
		if parsed < 1 || parsed > 100 {
			return 0, 0, "", nil, "limit must be in [1; 100]", http.StatusBadRequest
		}
		limit = parsed
	}

	category = q.Get("category")

	if v := q.Get("priceLessThan"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, 0, "", nil, "invalid priceLessThan", http.StatusBadRequest
		}
		if f <= 0 {
			return 0, 0, "", nil, "priceLessThan must be > 0", http.StatusBadRequest
		}
		priceLessThan = &f
	}

	return offset, limit, category, priceLessThan, "", 0
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	offset, limit, category, pricePtr, errMsg, status := parseQueryParams(r)
	if errMsg != "" {
		api.ErrorResponse(w, status, errMsg)
		return
	}

	products, total, err := h.service.GetCatalog(category, pricePtr, offset, limit)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respProducts := make([]ProductSummary, len(products))
	for i, p := range products {
		respProducts[i] = ProductSummary{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	response := Response{
		Items: respProducts,
		Total: total,
	}

	api.OKResponse(w, response)
}
