package request

type CreatePostRequest struct {
	Text string `json:"text"  validate:"required,min=1,max=500"`
}
