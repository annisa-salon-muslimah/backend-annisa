package input

type InputBlog struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
}
