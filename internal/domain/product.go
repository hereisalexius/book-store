package domain

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProductRequest struct {
	Name  string  `json:"name"  binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price" binding:"omitempty,gt=0"`
}

type PatchProductRequest struct {
	Name  *string  `json:"name"  binding:"omitempty,min=1"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
}
