package product

import "C"
import (
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/service"
	"exercises/Catalog/store"
	"fmt"
	"log"
)

type prodService struct {
	prodStore  store.Product
	brandStore store.Brand
}

//type category struct {
//	prodStore  store.Product
//	brandStore store.Brand
//}

func New(prodStore store.Product, brandStore store.Brand) service.Product {
	return prodService{
		prodStore:  prodStore,
		brandStore: brandStore,
	}
}


func (S prodService) GetById(id int) (model.Prod, error) {
	prodDet, err := S.prodStore.GetById(id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return model.Prod{}, errors.New("Id not found")
	}
	brandDet, err := S.brandStore.GetById(prodDet.BrandDetails.Id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return prodDet, errors.New("Brand Id not found")
	}
	prodDet.BrandDetails.Brand = brandDet.Brand
	return prodDet, nil
}
func (S prodService) GetAll() ([]model.Prod, error) {
	prodDet, err := S.prodStore.GetAll()
	log.Println(err)
	if err != nil {
		log.Println(err)
		return []model.Prod{}, errors.New("Products not found")
	}
	for i,p:=range prodDet {
		brandDet, err := S.brandStore.GetById(p.BrandDetails.Id)
		log.Println(err)
		if err != nil {
			log.Println(err)
			return prodDet, errors.New("Brand Id not found")
		}
		prodDet[i].BrandDetails.Brand = brandDet.Brand
	}
	return prodDet, nil
}

func (S prodService) Create(empName, bName string) (model.Prod, error) {
	bId, err := S.brandStore.CheckBrand(bName)
	if err != nil || bId == 0 {
		res, err := S.brandStore.Create(bName)
		if err != nil {
			return model.Prod{}, errors.New("Brand not created")
		}
		bId = res
	}
	res, err := S.prodStore.Create(empName, bId)
	if err != nil {
		return model.Prod{}, errors.New("Product not created")
	}
	prodDet, err := S.prodStore.GetById(res)
	if err != nil {
		return model.Prod{}, errors.New("Created Product not found")
	}
	brandDet, err := S.brandStore.GetById(prodDet.BrandDetails.Id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return prodDet, errors.New("Brand not found")
	}
	prodDet.BrandDetails.Brand = brandDet.Brand
	return prodDet, nil
}
func (S prodService) Update(empId int, empName, bName string) (model.Prod, error) {
	bId, err := S.brandStore.CheckBrand(bName)
	if err != nil || bId == 0 {
		res, err := S.brandStore.Create(bName)
		if err != nil {
			return model.Prod{}, err
		}
		bId = res
	}
	res, err := S.prodStore.Update(empId,empName, bId)
	if err != nil {
		return model.Prod{}, errors.New("Product not updated")
	}
	fmt.Println(res)
	prodDet, err := S.prodStore.GetById(empId)
	if err != nil {
		return model.Prod{}, errors.New("Updated Product not found")
	}
	brandDet, err := S.brandStore.GetById(prodDet.BrandDetails.Id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return prodDet, errors.New("Updated Brand not found")
	}
	prodDet.BrandDetails.Brand = brandDet.Brand
	return prodDet, nil
}
func (S prodService) Delete(empId int) error {
	res, err := S.prodStore.Delete(empId)
	if err != nil || res != 1 {
		return  errors.New("Product not deleted from database")
	}
	return nil
}
