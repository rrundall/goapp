package api

import (
	"errors"
	"fmt"
	"goapp/config"
	"goapp/pkg/db"
	"goapp/pkg/model"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "goapp/docs"
)

type Server struct {
	address *http.Server
	router  *gin.Engine
	db      db.Storage
}

// GetServer to initialize the api server and database
func GetServer(addr string, h *gin.Engine, d db.Storage) *Server {
	log.Debug().Msgf("GetServer info -> %s, %v, %v", addr, h, d)
	return &Server{
		address: &http.Server{
			Addr:    addr,
			Handler: h,
		},
		router: h,
		db:     d,
	}
}

// StartServer to start up the services
func (s *Server) StartServer() error {
	docs.SwaggerInfo.BasePath = "/v1"
	v1 := s.router.Group("/v1")
	{
		v1.GET("/", s.homePageRequest)
		v1.GET("/books", s.listBooksRequest)
		v1.POST("/books/search", s.searchBooksRequest)
		v1.POST("/books/get", s.getBooksRequest)
		v1.POST("/books", s.insertBooksRequest)
		v1.PUT("/books", s.updateBooksRequest)
		v1.PATCH("/books", s.patchBooksRequest)
		v1.DELETE("/books/:id", s.deleteBooksRequest)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return s.address.ListenAndServe()
}

// homePageRequest for accessing to home page
func (s *Server) homePageRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": config.HomepageMsg})
}

// listBooksRequest godoc
//
//	@Summary		Get Books
//	@Description	For listing books per page.
//	@Description	By default will order by book_id and displays 1000 books in a page.
//	@Tags			books
//	@Produce		json
//	@Param			order_by	query	string	false	"Order by field"	default(book_id)
//	@Param			page_id		query	int		false	"Page number"		default(1)	minimum(1)
//	@Param			page_size	query	int		false	"Results per page"	default(25)	minimum(5)	maximum(1000)
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/books [get]
func (s *Server) listBooksRequest(c *gin.Context) {
	// Default page list configuration
	var list *model.ListBookRequest
	if err := c.ShouldBindQuery(&list); err != nil {
		log.Error().Msgf("listBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": config.BadRequestErrMsg})
		return
	}

	p := &db.PageList{
		OrderBy: list.OrderBy,
		Limit:   list.PageSize,
		OffSet:  (list.PageID - 1) * list.PageSize,
	}
	bks, err := s.db.ListBooks(p)
	if err != nil {
		log.Error().Msgf("allBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		c.JSON(http.StatusOK, bks)
	}
}

// searchBooksRequest godoc
//
//	@Summary		Search Books
//	@Description	For searching books with OR criteria and using LIKE %string% pattern.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			body	body	model.Book	false	"Fields Required: At least one. Empty fields will be ignored"
//	@Success		200
//	@Failure		415
//	@Failure		500
//	@Router			/books/search [post]
func (s *Server) searchBooksRequest(c *gin.Context) {
	if !ValidateContentType(c) {
		return
	}

	var bk *model.Book
	if err := c.ShouldBindJSON(&bk); err != nil {
		log.Error().Msgf("%s: %s", config.InvalidDataErrMsg, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
		return
	}

	v := reflect.Indirect(reflect.ValueOf(bk))
	t := v.Type()
	var str []string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != "" && !(t.Field(i).Type.String() == "int" && v.Field(i).Interface() == 0) {
			str = append(str, fmt.Sprintf("%s LIKE '%%%v%%'", t.Field(i).Tag.Get("db"), v.Field(i).Interface()))
		}
	}
	if !WarnEmptyData(c, str) {
		return
	}

	bks, err := s.db.GetBooks(strings.Join(str, " OR "))
	if err != nil {
		log.Error().Msgf("searchBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		c.JSON(http.StatusOK, bks)
	}
}

// getBooksRequest godoc
//
//	@Summary		Find Matching Books
//	@Description	For searching books with AND criteria and using WHERE column = string pattern.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			body	body	model.Book	false	"Fields Required: At least one. Empty fields will be ignored"
//	@Success		200
//	@Failure		415
//	@Failure		500
//	@Router			/books/get [post]
func (s *Server) getBooksRequest(c *gin.Context) {
	if !ValidateContentType(c) {
		return
	}

	var bk *model.Book
	if err := c.ShouldBindJSON(&bk); err != nil {
		log.Error().Msgf("%s: %s", config.InvalidDataErrMsg, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
		return
	}

	v := reflect.Indirect(reflect.ValueOf(bk))
	t := v.Type()
	var str []string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != "" && !(t.Field(i).Type.String() == "int" && v.Field(i).Interface() == 0) {
			str = append(str, fmt.Sprintf("%s = '%v'", t.Field(i).Tag.Get("db"), v.Field(i).Interface()))
		}
	}

	if !WarnEmptyData(c, str) {
		return
	}

	bks, err := s.db.GetBooks(strings.Join(str, " AND "))
	if err != nil {
		log.Error().Msgf("getBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		c.JSON(http.StatusOK, bks)
	}
}

// insertBooksRequest godoc
//
//	@Summary		Insert Books
//	@Description	For inserting single/multiple books.
//	@Description	Will return number of rows that are inserted, if there is no row updated, will return no data update with 0 row affected.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			body	body	[]model.Book	true	"Fields Required: ALL except book_id. Fields cannot be empty. Unique fields: isbn. If book_id is included it will be ignored."
//	@Success		200
//	@Failure		415
//	@Failure		500
//	@Router			/books [post]
func (s *Server) insertBooksRequest(c *gin.Context) {
	if !ValidateContentType(c) {
		return
	}
	var bks []*model.Book
	if err := c.ShouldBindJSON(&bks); err != nil {
		log.Error().Msgf("%s: %s", config.InvalidDataErrMsg, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
		return
	}

	var str []string
	for i := 0; i < len(bks); i++ {
		v := reflect.Indirect(reflect.ValueOf(bks[i]))
		t := v.Type()
		for j := 0; j < v.NumField(); j++ {
			if t.Field(j).Type.String() == "string" && v.Field(j).Interface() == "" {
				log.Error().Msgf("%s %s", t.Field(j).Tag.Get("json"),
					errors.New(config.DataCouldNotBeEmptyErrMsg))
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s %s",
					t.Field(j).Tag.Get("json"), config.DataCouldNotBeEmptyErrMsg)})
				return
			}
		}
		str = append(str, fmt.Sprintf("('%v', '%v', '%v', '%v', '%v', '%v')",
			bks[i].ISBN, bks[i].Title, bks[i].AuthorName, bks[i].AuthorSurname, bks[i].Published, bks[i].Publisher))
	}

	rowsAffected, err := s.db.InsertBooks(strings.Join(str, ","))
	if err != nil {
		log.Error().Msgf("insertBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		ValidateRowsAffected(c, rowsAffected, config.AddSuccessMsg)
	}
}

// updateBooksRequest godoc
//
//	@Summary		Update Book by book_id
//	@Description	For updating a book by book_id.
//	@Description	Will return number of row that is updated, if there is no row updated, will return no data update with 0 row affected.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			body	body	model.Book	true	"Fields Required: ALL. Fields cannot be empty. Unique fields: isbn."
//	@Success		200
//	@Failure		415
//	@Failure		500
//	@Router			/books [put]
func (s *Server) updateBooksRequest(c *gin.Context) {
	if !ValidateContentType(c) {
		return
	}
	var bk *model.Book
	if err := c.ShouldBindJSON(&bk); err != nil {
		log.Error().Msgf("%s: %s", config.InvalidDataErrMsg, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
		return
	}

	v := reflect.Indirect(reflect.ValueOf(bk))
	t := v.Type()
	for j := 0; j < v.NumField(); j++ {
		if (t.Field(j).Type.String() == "string" && v.Field(j).Interface() == "") ||
			(t.Field(j).Type.String() == "int" && v.Field(j).Interface() == 0) {
			log.Error().Msgf("%s %s", t.Field(j).Tag.Get("json"),
				errors.New(config.DataCouldNotBeEmptyErrMsg))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s %s",
				t.Field(j).Tag.Get("json"), config.DataCouldNotBeEmptyErrMsg)})
			return
		}
	}

	rowsAffected, err := s.db.UpdateBooks(bk)
	if err != nil {
		log.Error().Msgf("updateBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		ValidateRowsAffected(c, rowsAffected, config.UpdateSuccessMsg)
	}
}

// patchBooksRequest godoc
//
//	@Summary		Update Book by book_id
//	@Description	For updating a book by book_id.
//	@Description	Will return number of row that is updated, if there is no row updated, will return no data update with 0 row affected.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			body	body	model.Book	true	"Fields Required: book_id. Empty fields will be ignored. Unique fields: isbn."
//	@Success		200
//	@Failure		415
//	@Failure		500
//	@Router			/books [patch]
func (s *Server) patchBooksRequest(c *gin.Context) {
	if !ValidateContentType(c) {
		return
	}
	var bk *model.PatchBook
	if err := c.ShouldBindJSON(&bk); err != nil {
		log.Error().Msgf("%s: %s", config.InvalidDataErrMsg, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
		return
	}

	v := reflect.Indirect(reflect.ValueOf(bk))
	t := v.Type()
	var str []string
	var emptyFields []string
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Type.String() == "int" && v.Field(i).Interface() == 0 {
			log.Error().Msgf("%s %s", t.Field(i).Tag.Get("json"),
				errors.New(config.DataCouldNotBeEmptyErrMsg))
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s %s",
				t.Field(i).Tag.Get("json"), config.DataCouldNotBeEmptyErrMsg)})
			return
		} else if t.Field(i).Type.String() == "string" && v.Field(i).Interface() == "" {
			emptyFields = append(emptyFields, t.Field(i).Tag.Get("json"))
		} else {
			str = append(str, fmt.Sprintf("%s = '%v'", t.Field(i).Tag.Get("db"), v.Field(i).Interface()))
		}
	}

	rowsAffected, err := s.db.PatchBooks(fmt.Sprintf("%s WHERE book_id = %v",
		strings.Join(str, ","), bk.ID))
	if err != nil {
		log.Error().Msgf("patchBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		if msg := WarnFieldsCannotBeEmpty(emptyFields); msg != "" {
			c.JSON(http.StatusOK, gin.H{"message": config.UpdateSuccessMsg, "rows_affected": rowsAffected,
				"warning": msg})
		} else {
			ValidateRowsAffected(c, rowsAffected, config.UpdateSuccessMsg)
		}
	}
}

// deleteBooksRequest godoc
//
//	@Summary		Delete Book
//	@Description	For deleting book by id.
//	@Description	Header is required for content-type.
//	@Description	Will return number of row that is deleted, if there is no row deleted, will return no data update with 0 row affected.
//	@Tags			books
//	@Produce		json
//	@Param			id	path	string	true	"The book_id to be deleted."
//	@Success		200
//	@Failure		500
//	@Router			/books/{id} [delete]
func (s *Server) deleteBooksRequest(c *gin.Context) {
	rowsAffected, err := s.db.DeleteBooks(c.Param("id"))
	if err != nil {
		log.Error().Msgf("deleteBooksRequest failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.DBOperationErrMsg})
	} else {
		ValidateRowsAffected(c, rowsAffected, config.DeleteSuccessMsg)
	}
}
