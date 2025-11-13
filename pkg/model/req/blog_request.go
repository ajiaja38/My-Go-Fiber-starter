package req

type CreateBlogDto struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}
