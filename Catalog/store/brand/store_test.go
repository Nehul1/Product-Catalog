package brand

import (
	"errors"
	"exercises/Catalog/model"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestBrandStore_GetById(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	bS := New(fDB)
	rows := sqlmock.NewRows([]string{"id", "brand"}).AddRow(2, "Puma")
	//mock.ExpectBegin()
	mock.ExpectQuery("select B.id,B.brand from brand as B*").WithArgs(2).WillReturnRows(rows)
	res, err := bS.GetById(2)
	//fmt.Println(res)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	tc := model.Brand{Id: 2,Brand: "Puma"}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//TestCase-2
	//mock.ExpectBegin()
	mock.ExpectQuery("select B.id,B.brand from brand as B*").WithArgs(5).WillReturnError(errors.New("Brand Id not Found"))
	res, err = bS.GetById(5)
	//fmt.Println(res)
	expErr := errors.New("Brand Id not Found")
	//if err.Error()!=expErr.Error(){
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	tc = model.Brand{Id: 0,Brand: ""}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//Test Case-3
	rows = sqlmock.NewRows([]string{"id"}).AddRow(2)
	//mock.ExpectBegin()
	mock.ExpectQuery("select B.id,B.brand from brand as B*").WithArgs(2).WillReturnRows(rows)
	res, err = bS.GetById(2)
	expErr=errors.New("sql: expected 1 destination arguments in Scan, not 2")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr,err)
	}
	tc = model.Brand{}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}
}
func TestBrandStore_GetAll(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	bS := New(fDB)
	rows := sqlmock.NewRows([]string{"id", "brand"}).
		AddRow(2, "Puma").
		AddRow(1,"Nike").
		AddRow(3,"Adidas")
	//mock.ExpectBegin()
	mock.ExpectQuery("select B.id,B.brand from brand").WillReturnRows(rows)
	res, err := bS.GetAll()
	//fmt.Println(res)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	tc := []model.Brand{{Id: 2,Brand: "Puma"},{Id: 1,Brand: "Nike"},{Id: 3,Brand: "Adidas"}}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//TestCase-2
	//mock.ExpectBegin()
	mock.ExpectQuery("select B.id,B.brand from brand").WillReturnError(errors.New("Brands not Found"))
	res, err = bS.GetAll()
	//fmt.Println(res)
	expErr := errors.New("Brands not Found")
	//if err.Error()!=expErr.Error(){
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	tc = []model.Brand{}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}
}

func TestBrandStore_CheckBrand(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	bS := New(fDB)

	//TestCase1
	rows := sqlmock.NewRows([]string{"id"}).AddRow(2)
	//mock.ExpectBegin()
	mock.ExpectQuery("select id from brand*").WithArgs("Puma").WillReturnRows(rows)
	res, err := bS.CheckBrand("Puma")
	//fmt.Println(res)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected %v error but got %v", nil, err)
	}
	tc := 2
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v but got %v", res, tc)
	}

	//TestCase-2
	//mock.ExpectBegin()
	mock.ExpectQuery("select id from brand*").WithArgs("Hello").WillReturnError(errors.New("Brand Name not Found"))
	res, err = bS.CheckBrand("Hello")
	//fmt.Println(res)
	expErr := errors.New("Brand Name not Found")
	//if err.Error()!=expErr.Error(){
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	tc = 0
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//TestCase-3
	rows = sqlmock.NewRows([]string{"id","name"}).AddRow(1,"Raj")
	//mock.ExpectBegin()
	mock.ExpectQuery("select id from brand*").WithArgs("Puma").WillReturnRows(rows)
	res, err = bS.CheckBrand("Puma")
	expErr=errors.New("sql: expected 2 destination arguments in Scan, not 1")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	tc = 0
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v but got %v", res, tc)
	}

}
func TestBrandStore_Create(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	bS := New(fDB)

	mock.ExpectExec("insert into brand").WithArgs("Asics").WillReturnResult(sqlmock.NewResult(6, 1))
	res1, err := bS.Create("Asics")
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 6) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//Test Case 2
	mock.ExpectExec("insert into brand").WithArgs("Asics").WillReturnError(errors.New("Brand not created"))
	res1, err = bS.Create("Asics")
	if !reflect.DeepEqual(err, errors.New("Brand not created")) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//TestCase-3
	mock.ExpectExec("insert into brand").WithArgs("Asics").WillReturnResult(sqlmock.NewResult(0, 0))
	res1, err = bS.Create("Asics")
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 0, res1)
	}

}
