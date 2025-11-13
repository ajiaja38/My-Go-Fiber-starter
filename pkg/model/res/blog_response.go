package res

type FindAllBlogsResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Image     string `json:"image"`
	UserId    string `json:"userId"`
	Owner     string `json:"owner"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
