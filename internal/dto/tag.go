package dto

import "github.com/Vatsal-S-Patel/Bloggy/models"

type AddTagRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}

type TagBlogs struct {
	Name  string         `json:"name"`
	Blogs []*models.Blog `json:"blogs"`
}
