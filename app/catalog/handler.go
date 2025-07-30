package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/services"
	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	service *services.ProductService
}

func NewCatalogHandler(s *services.ProductService) *CatalogHandler {
	return &CatalogHandler{
		service: s,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Параметры пагинации и фильтрации
	q := r.URL.Query()
	offset := 0
	limit := 10
	category := q.Get("category")
	priceLessThan := q.Get("priceLessThan")
	if v := q.Get("offset"); v != "" {
		fmt.Sscanf(v, "%d", &offset)
	}
	if v := q.Get("limit"); v != "" {
		fmt.Sscanf(v, "%d", &limit)
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	var pricePtr *float64
	if priceLessThan != "" {
		var f float64
		fmt.Sscanf(priceLessThan, "%f", &f)
		pricePtr = &f
	}
	products, total, err := h.service.GetCatalog(category, pricePtr, offset, limit)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	respProducts := make([]Product, len(products))
	for i, p := range products {
		respProducts[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}
	response := Response{
		Products: respProducts,
		Total:    total,
	}
	api.OKResponse(w, response)
}
