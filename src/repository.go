package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
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

type InMemoryRepository struct {
	Products []ProductDAO
}

func NewInMemoryRepository(jsonFilePath string) (InMemoryRepository, error) {

	ret := InMemoryRepository{}

	content, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return ret, errors.New("error opening file " + jsonFilePath + ": " + err.Error())
	}

	err = json.Unmarshal(content, &ret.Products)
	if err != nil {
		return ret, errors.New("error unmarshalling json: " + err.Error())
	}

	return ret, nil
}

// Uses binary search to find the product in the already sorted products.json repository
func (m InMemoryRepository) Find(id int) (ProductDAO, error) {
	ret := sort.Search(len(m.Products), func(i int) bool { return m.Products[i].Id >= id })
	if ret < len(m.Products) && m.Products[ret].Id == id {
		return m.Products[ret], nil
	}
	return ProductDAO{}, ErrProductNotFound
}

func (m InMemoryRepository) FindGift() (ProductDAO, error) {
	for _, p := range m.Products {
		if p.Is_gift {
			return p, nil
		}
	}

	return ProductDAO{}, ErrNoGiftFound
}

type Repository interface {
	Find(id int) (ProductDAO, error)
	FindGift() (ProductDAO, error)
}
