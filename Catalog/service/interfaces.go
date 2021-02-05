package service

import "exercises/Catalog/model"

type Product interface {
	GetById(int) (model.Prod, error)
	GetAll()([]model.Prod,error)
	Create(string, string) (model.Prod, error)
	Update(int,string,string)(model.Prod,error)
	Delete(int)error
}
