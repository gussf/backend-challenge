package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ProductDAO struct {
	Id          int
	Title       string
	Description string
	Amount      int
	IsGift      bool
}

type InMemoryRepository struct {
	Products []ProductDAO
}

func NewInMemoryRepository(jsonFilePath string) (InMemoryRepository, error) {

	ret := InMemoryRepository{}

	content, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return ret, errors.New("Error opening file " + jsonFilePath + ": " + err.Error())
	}

	err = json.Unmarshal(content, &ret.Products)
	if err != nil {
		return ret, errors.New("Error unmarshalling json: " + err.Error())
	}

	return ret, nil
}

func (m InMemoryRepository) Find(id int) (ProductDAO, error) {
	return ProductDAO{}, nil
}

type Repository interface {
	Find(id int) (ProductDAO, error)
}
