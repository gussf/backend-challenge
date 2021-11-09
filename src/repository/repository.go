package repository

import (
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found in repository")
	ErrNoGiftFound     = errors.New("no gift was found in repository")
)

type ProductDAO struct {
	Id          int
	Title       string
	Description string
	Amount      int
	Is_gift     bool
}

type Repository interface {
	Find(id int) (ProductDAO, error)
	FindGift() (ProductDAO, error)
}
