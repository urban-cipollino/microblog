package request

type CreateUserRequest struct {
	Name     string `json:"name"  validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
