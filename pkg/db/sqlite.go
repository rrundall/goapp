package db

import (
	"fmt"
	"goapp/config"
	"goapp/pkg/model"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

type SqliteStorage struct{ db *sqlx.DB }
type PageList struct {
	OrderBy string
	Limit   int
	OffSet  int
}
type Storage interface {
	ListBooks(p *PageList) ([]model.Book, error)
	GetBooks(str string) ([]model.Book, error)
	InsertBooks(str string) (int64, error)
	UpdateBooks(bk *model.Book) (int64, error)
	PatchBooks(str string) (int64, error)
	DeleteBooks(str string) (int64, error)
}

// OpenSqliteStorage to initialize the sqlite database
func OpenSqliteStorage() *SqliteStorage {
	return &SqliteStorage{db: OpenDB()}
}

// OpenDB will access to sqlite database that locate in local file path
func OpenDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite", config.DBFile)
	if err != nil {
		log.Fatal().Err(err).Msg(config.DBConnectErrMsg)
	}
	return db
}

// CloseDB to close the database connection
func (s SqliteStorage) CloseDB() {
	if err := s.db.Close(); err != nil {
		log.Fatal().Err(err).Msg(config.DBCloseErrMsg)
	}
}

// ListBooks will return all books with order by, page id and page size configuration that passing through
func (s SqliteStorage) ListBooks(p *PageList) ([]model.Book, error) {
	query := fmt.Sprintf("SELECT * FROM book ORDER BY %v LIMIT %v OFFSET %v", p.OrderBy, p.Limit, p.OffSet)
	rows, err := s.db.Queryx(query)
	log.Debug().Msgf("ListBooks: %s", query)
	if err != nil {
		return nil, err
	}

	bks := []model.Book{}
	defer rows.Close()
	for rows.Next() {
		var bk model.Book
		if err = rows.StructScan(&bk); err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	log.Debug().Msgf("%v", bks)
	return bks, nil
}

// GetBooks will return all books that match with condition string that passing through
func (s SqliteStorage) GetBooks(str string) ([]model.Book, error) {
	query := fmt.Sprintf("SELECT * FROM book WHERE %s", str)
	rows, err := s.db.Queryx(query)
	log.Debug().Msgf("FindAllBooks: %s", query)
	if err != nil {
		return nil, err
	}

	bks := []model.Book{}
	defer rows.Close()
	for rows.Next() {
		var bk model.Book
		if err = rows.StructScan(&bk); err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	log.Debug().Msgf("%v", bks)
	return bks, nil
}

// InsertBooks is able to insert single/multiple books depends on string that passing through
// and will return the number of rows that inserted.
func (s SqliteStorage) InsertBooks(str string) (int64, error) {
	query := fmt.Sprintf("INSERT INTO book (isbn, title, author_name, author_surname, published, publisher) "+
		"VALUES %s", str)
	result, err := s.db.Exec(query)
	log.Debug().Msgf("InsertBooks: %s", query)

	var rowsAffected int64
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}
	log.Debug().Msgf("RowsAffected: %d", rowsAffected)
	return rowsAffected, nil
}

// UpdateBooks will update single book and all the book fields are required
// it will return number of book that is updated and return 0 if no book update
func (s SqliteStorage) UpdateBooks(bk *model.Book) (int64, error) {
	query := fmt.Sprintf("UPDATE book SET isbn = '%v', title = '%v', author_name = '%v', author_surname = '%v', "+
		"published = '%v', publisher = '%v' WHERE book_id = '%v'", bk.ISBN, bk.Title, bk.AuthorName, bk.AuthorSurname,
		bk.Published, bk.Publisher, bk.ID)
	result, err := s.db.Exec(query)
	log.Debug().Msgf("UpdateBooks: %s", query)

	var rowsAffected int64
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}
	log.Debug().Msgf("RowsAffected: %d", rowsAffected)
	return rowsAffected, nil
}

// PatchBooks will patch single book and only book_id that is required,
// other field that is empty or not define will be ignored
// it will return number of book that is updated and return 0 if no book update
func (s SqliteStorage) PatchBooks(str string) (int64, error) {
	query := fmt.Sprintf("UPDATE book SET %s", str)
	result, err := s.db.Exec(query)
	log.Debug().Msgf("UpdateBooks: %s", query)

	var rowsAffected int64
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}
	log.Debug().Msgf("RowsAffected: %d", rowsAffected)
	return rowsAffected, nil
}

// DeleteBooks will delete a book id is matched and return number of book that is deleted and return 0 if no book delete
func (s SqliteStorage) DeleteBooks(str string) (int64, error) {
	query := fmt.Sprintf("DELETE from book WHERE book_id = '%s'", str)
	result, err := s.db.Exec(query)
	log.Debug().Msgf("DeleteBooks: %s", query)

	var rowsAffected int64
	if err != nil {
		return rowsAffected, err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return rowsAffected, err
	}
	log.Debug().Msgf("RowsAffected: %d", rowsAffected)
	return rowsAffected, nil
}
