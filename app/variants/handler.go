package variants

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type ProductDetailService interface {
	GetProductDetail(code string) (*models.Product, error)
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Variant struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type ProductDetail struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category Category  `json:"category"`
	Variants []Variant `json:"variants"`
}

type Handler struct {
	service ProductDetailService
}

func NewHandler(s ProductDetailService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("id")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "missing id")
		return
	}

	product, err := h.service.GetProductDetail(code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.ErrorResponse(w, http.StatusNotFound, "product not found")
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	variants := make([]Variant, len(product.Variants))
	for i, v := range product.Variants {
		variants[i] = Variant{
			ID:    v.ID,
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price.InexactFloat64(),
		}
	}

	resp := ProductDetail{
		Code:  product.Code,
		Price: product.Price.InexactFloat64(),
		Category: Category{
			Code: product.Category.Code,
			Name: product.Category.Name,
		},
		Variants: variants,
	}

	api.OKResponse(w, resp)
}