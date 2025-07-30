package variants

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/mytheresa/go-hiring-challenge/services"
	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductDetail struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Variants []Variant `json:"variants"`
}

type Handler struct {
	service *services.ProductService
}

func NewHandler(s *services.ProductService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["id"]
	product, err := h.service.GetProductDetail(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}
	variants := make([]Variant, len(product.Variants))
	for i, v := range product.Variants {
		price := v.Price.InexactFloat64()
		if price == 0 {
			price = product.Price.InexactFloat64()
		}
		variants[i] = Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price,
		}
	}
	resp := ProductDetail{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: product.Category.Name,
		Variants: variants,
	}
	api.OKResponse(w, resp)
} 