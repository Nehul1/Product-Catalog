package brand

import (
	"database/sql"
	"exercises/Catalog/model"
	"exercises/Catalog/store"
	_ "github.com/go-sql-driver/mysql"
)

type brandStore struct {
	DB *sql.DB
}

func New(db *sql.DB) store.Brand {
	return brandStore{
		DB: db,
	}
}
func (bS brandStore) CheckBrand(bName string) (int, error) {
	var b model.Brand
	res, err := bS.DB.Query("select id from brand where brand=?", bName)
	if err != nil {
		return 0, err
		//return 0, errors.New("Brand Name not Found")
	}

	defer res.Close()
	rescount := 0
	for res.Next() {
		rescount = rescount + 1
		err := res.Scan(&b.Id)
		if err != nil {
			return 0, err
		}
	}
	if rescount == 0 {
		return 0, err
	}
	return b.Id, nil
}

func (bS brandStore) Create(bName string) (int, error) {
	res, err := bS.DB.Exec("insert into brand(brand) values(?)", bName)
	if err != nil {
		return 0, err
		//return 0, errors.New("Brand not created")
	}
	lastId, err := res.LastInsertId()
	if err != nil || lastId==0{
		return 0, err
	}
	return int(lastId), nil
}
func (bS brandStore) GetById(key int) (model.Brand, error) {
	var b model.Brand
	var res *sql.Rows
	var err error
	res, err = bS.DB.Query("select B.id,B.brand from brand as B where id=?", key)
	if err != nil {
		return model.Brand{}, err
		//return model.Brand{}, errors.New("Brand Id not Found")
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&b.Id, &b.Brand)
		if err != nil {
			return model.Brand{}, err
		}
	}
	return b, nil
}

func (bS brandStore) GetAll() ([]model.Brand, error) {
	var b model.Brand
	var res *sql.Rows
	var err error
	brandDet:=make([]model.Brand,0)
	res, err = bS.DB.Query("select B.id,B.brand from brand as B")
	if err != nil {
		return []model.Brand{}, err
		//return []model.Brand{}, errors.New("Brands not Found")
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&b.Id, &b.Brand)
		if err != nil {
			return []model.Brand{}, err
		}
		brandDet=append(brandDet,b)
	}
	return brandDet, nil
}
