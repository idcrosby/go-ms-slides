package catalogue

import (
	"errors"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")
var ErrDBConnection = errors.New("database connection error")

var baseQuery = "SELECT sock.sock_id AS id, sock.name, sock.description, sock.price, sock.count, sock.image_url_1, sock.image_url_2, GROUP_CONCAT(tag.name) AS tag_name FROM sock JOIN sock_tag ON sock.sock_id=sock_tag.sock_id JOIN tag ON sock_tag.tag_id=tag.tag_id"

type Sock struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	ImageURL    []string `json:"imageUrl" db:"-"`
	Price       float32  `json:"price" db:"price"`
}

// START OMIT
type Service interface {
	List(tags []string, order string, pageNum, pageSize int) ([]Sock, error) // GET /catalogue
	Get(id string) (Sock, error)                                             // GET /catalogue/{id}
}

func NewCatalogueService(db *sqlx.DB, logger log.Logger) Service {
	return &catalogueService{
		db:     db,
		logger: logger,
	}
}

type catalogueService struct {
	db     *sqlx.DB
	logger log.Logger
}

// END OMIT
// START 2-OMIT
func (s *catalogueService) List(order string, pageNum, pageSize int) ([]Sock, error) {
	var socks []Sock
	query := baseQuery

	var args []interface{}

	query += " GROUP BY id"

	if order != "" {
		query += " ORDER BY ?"
		args = append(args, order)
	}

	query += ";"

	err := s.db.Select(&socks, query, args...)
	if err != nil {
		s.logger.Log("database error", err)
		return []Sock{}, ErrDBConnection
	}

	socks = cut(socks, pageNum, pageSize)

	return socks, nil
}

// END 2-OMIT
func (s *catalogueService) Get(id string) (Sock, error) {
	query := baseQuery + " WHERE sock.sock_id =? GROUP BY sock.sock_id;"

	var sock Sock
	err := s.db.Get(&sock, query, id)
	if err != nil {
		s.logger.Log("database error", err)
		return Sock{}, ErrNotFound
	}

	sock.ImageURL = []string{sock.ImageURL_1, sock.ImageURL_2}
	sock.Tags = strings.Split(sock.TagString, ",")

	return sock, nil
}

func cut(socks []Sock, pageNum, pageSize int) []Sock {
	if pageNum == 0 || pageSize == 0 {
		return []Sock{} // pageNum is 1-indexed
	}
	start := (pageNum * pageSize) - pageSize
	if start > len(socks) {
		return []Sock{}
	}
	end := (pageNum * pageSize)
	if end > len(socks) {
		end = len(socks)
	}
	return socks[start:end]
}
