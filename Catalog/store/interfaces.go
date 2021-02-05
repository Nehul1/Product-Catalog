package store

import "exercises/Catalog/model"

type Brand interface {
	GetById(int) (model.Brand, error)
	GetAll()([]model.Brand,error)
	Create(string) (int, error)
	CheckBrand(string) (int, error)
}

type Product interface {
	GetById(int) (model.Prod, error)
	GetAll()([]model.Prod,error)
	Create(string, int) (int, error)
	Update(int,string,int)(int,error)
	Delete(int)(int,error)
}
