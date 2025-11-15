package res

type FindOwnBlogResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Image     string `json:"image"`
	UserId    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type FindBlogResponse struct {
	FindOwnBlogResponse
	Owner string `json:"owner"`
}
