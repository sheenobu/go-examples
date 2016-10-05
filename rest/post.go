package rest

// A Post is a blog post written by a user
type Post struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	User  User   `json:"user"`
}
