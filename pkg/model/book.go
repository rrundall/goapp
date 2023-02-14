package model

// Book is main struct for most of the handlers which is no require flag set
type Book struct {
	ID            int    `json:"book_id" db:"book_id"`
	ISBN          string `json:"isbn" db:"isbn"`
	Title         string `json:"title" db:"title"`
	AuthorName    string `json:"author_name" db:"author_name"`
	AuthorSurname string `json:"author_surname" db:"author_surname"`
	Published     string `json:"published" db:"published"`
	Publisher     string `json:"publisher" db:"publisher"`
}

// PatchBook for handler that need to have required flag set like patch api
type PatchBook struct {
	ID            int    `json:"book_id" db:"book_id" binding:"required"`
	ISBN          string `json:"isbn,omitempty" db:"isbn"`
	Title         string `json:"title,omitempty" db:"title"`
	AuthorName    string `json:"author_name,omitempty" db:"author_name"`
	AuthorSurname string `json:"author_surname,omitempty" db:"author_surname"`
	Published     string `json:"published,omitempty" db:"published"`
	Publisher     string `json:"publisher,omitempty" db:"publisher"`
}

// ListBookRequest to define the order by filed, page id and page size to list the books
type ListBookRequest struct {
	OrderBy  string `form:"order_by,default=book_id" binding:"omitempty"`
	PageID   int    `form:"page_id,default=1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size,default=25" binding:"omitempty,min=5,max=1000"`
}
