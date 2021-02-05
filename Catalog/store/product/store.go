package product

import (
	"database/sql"
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/store"
	_ "github.com/go-sql-driver/mysql"
)

type prodStore struct {
	DB *sql.DB
}

func New(db *sql.DB) store.Product {
	return &prodStore{
		DB: db,
	}
}
func (pS prodStore) GetById(key int) (model.Prod, error) {
	var p model.Prod
	var res *sql.Rows
	var err error
	res, err = pS.DB.Query("select P.id,P.name,P.brand from product as P where id=?", key)
	if err != nil {
		//return model.Prod{}, errors.New("Product Id Not Found")
		return model.Prod{}, err
	}
	defer res.Close()
	resCount:=0
	for res.Next() {
		err := res.Scan(&p.Id, &p.Name, &p.BrandDetails.Id)
		if err != nil {
			return model.Prod{}, err
		}
		resCount++
	}
	if resCount==0{
		return p,errors.New("Product Not Found")
	}
	return p, nil
}
func (pS prodStore) GetAll() ([]model.Prod, error) {
	var p model.Prod
	var res *sql.Rows
	var err error
	prodDetails := make([]model.Prod, 0)
	res, err = pS.DB.Query("select P.id,P.name,P.brand from product as P")
	if err != nil {
		return []model.Prod{}, err
		//return []model.Prod{}, errors.New("Products Not Found")
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&p.Id, &p.Name, &p.BrandDetails.Id)
		if err != nil {
			return []model.Prod{}, err
		}
		prodDetails = append(prodDetails, p)
	}
	return prodDetails, nil
}
func (pS prodStore) Create(name string, brandId int) (int, error) {
	res, err := pS.DB.Exec("insert into product(name,brand) values(?,?)", name, brandId)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil || lastId <= 0 {
		return 0, err
	}
	return int(lastId), nil
}

func (pS prodStore) Update(key int,name string, brandId int) (int, error) {
	res, err := pS.DB.Exec("update product set name=?,brand=? where id=?", name, brandId,key)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil || lastId <= 0 {
		return 0, err
	}
	return int(lastId), nil
}

func (pS prodStore) Delete(key int) (int, error) {
	res, err := pS.DB.Exec("delete from product where id=?", key)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows != 1 {
		return 0, err
	}
	return int(rows), nil
}