package domain

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateCustomerRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}

type SyncResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Deleted int `json:"deleted"`
}
