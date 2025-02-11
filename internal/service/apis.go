package service

import "github.com/Kontentski/develops-today-task/internal/api/cat"

// APIs provides a collection of API interfaces.
type APIs struct {
	CatAPI
}

type CatAPI interface {
	GetBreeds() ([]cat.Breed, error)
}
