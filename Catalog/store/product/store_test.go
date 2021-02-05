package product

import (
	"errors"
	"exercises/Catalog/model"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestProdStore_GetAll(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	pS := New(fDB)
	rows := sqlmock.NewRows([]string{"id", "name", "brand"}).
		AddRow(1, "Black Sneakers", 2).
		AddRow(2, "Cricket Shoes", 3).
		AddRow(3, "Sandals", 1)
	//mock.ExpectBegin()
	mock.ExpectQuery("select P.id,P.name,P.brand from product as P").WillReturnRows(rows)
	res, err := pS.GetAll()
	//fmt.Println(res)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	tc := []model.Prod{{Id: 1,Name: "Black Sneakers",BrandDetails: model.Brand{2, ""}},
		{Id: 2,Name: "Cricket Shoes",BrandDetails: model.Brand{3, ""}},
		{Id: 3,Name: "Sandals",BrandDetails: model.Brand{1, ""}}}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//TestCase-2
	mock.ExpectQuery("select P.id,P.name,P.brand from product as P").WillReturnError(errors.New("Products Not Found"))
	res, err = pS.GetAll()
	//fmt.Println(res)
	expErr := errors.New("Products Not Found")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	tc = []model.Prod{}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//TestCase-3
	rows = sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Black Sneakers").
		AddRow(2, "Cricket Shoes").
		AddRow(3, "Sandals")
	//mock.ExpectBegin()
	mock.ExpectQuery("select P.id,P.name,P.brand from product as P").WillReturnRows(rows)
	res, err = pS.GetAll()
	//fmt.Println(res)
	expErr=errors.New("sql: expected 2 destination arguments in Scan, not 3")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v",expErr, err)
	}
	tc = []model.Prod{}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}
}

func TestProdStore_GetById(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	//bS := prodStore{DB: fDB}
	pS:=New(fDB)
	rows := sqlmock.NewRows([]string{"id", "name", "brand"}).AddRow(1, "Black Sneakers", 2)
	//mock.ExpectBegin()
	mock.ExpectQuery("select P.id,P.name,P.brand from product as P*").WithArgs(1).WillReturnRows(rows)
	res, err := pS.GetById(1)
	//fmt.Println(res)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	tc := model.Prod{1, "Black Sneakers", model.Brand{2, ""}}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//TestCase-2
	mock.ExpectQuery("select P.id,P.name,P.brand from product as P*").WithArgs(5).WillReturnError(errors.New("Product Id Not Found"))
	res, err = pS.GetById(5)
	//fmt.Println(res)
	expErr := errors.New("Product Id Not Found")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	tc = model.Prod{0, "", model.Brand{0, ""}}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

	//Test Case-3
	rows = sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Black Sneakers")
	//mock.ExpectBegin()
	mock.ExpectQuery("select P.id,P.name,P.brand from product as P*").WithArgs(1).WillReturnRows(rows)
	res, err = pS.GetById(1)
	expErr=errors.New("sql: expected 2 destination arguments in Scan, not 3")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v",expErr, err)
	}
	tc = model.Prod{}
	if !reflect.DeepEqual(res, tc) {
		t.Errorf("Expected %v error but got %v", res, tc)
	}

}

func TestProdStore_Create(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	pS := New(fDB)

	mock.ExpectExec("insert into product").WithArgs("Sandals", 3).WillReturnResult(sqlmock.NewResult(5, 1))
	res1, err := pS.Create("Sandals", 3)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 5) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//TestCase-2
	mock.ExpectExec("insert into product").WithArgs("Sandals", 3).WillReturnError(errors.New("Product not inserted into database"))
	res1, err = pS.Create("Sandals", 3)
	expErr := errors.New("Product not inserted into database")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//Test Case-3
	mock.ExpectExec("insert into product").WithArgs("Sandals", 3).WillReturnResult(sqlmock.NewResult(0, 0))
	res1, err = pS.Create("Sandals", 3)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}


}
func TestProdStore_Update(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	//pS := prodStore{DB: fDB}
	pS:=New(fDB)
	mock.ExpectExec("update product").WithArgs("Sandals", 3,5).WillReturnResult(sqlmock.NewResult(5, 1))
	res1, err := pS.Update(5,"Sandals", 3)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 5) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//TestCase-2
	mock.ExpectExec("update product").WithArgs("Sandals", 1,10).WillReturnError(errors.New("Product not updated in database"))
	res1, err = pS.Update(10,"Sandals", 1)
	expErr := errors.New("Product not updated in database")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//Test Case-3
	mock.ExpectExec("update product").WithArgs("Sandals", 3,5).WillReturnResult(sqlmock.NewResult(0, 1))
	res1, err = pS.Update(5,"Sandals", 3)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 0, res1)
	}
}
func TestProdStore_Delete(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a stub database connection", err)
	}
	defer fDB.Close()
	pS:=New(fDB)

	mock.ExpectExec("delete from product").WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))
	res1, err := pS.Delete(5)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 1) {
		t.Errorf("Expected last inserted id as %v but got %v", 1, res1)
	}

	//TestCase-2
	mock.ExpectExec("delete from product").WithArgs(10).WillReturnError(errors.New("Product not deleted in database"))
	res1, err = pS.Delete(10)
	expErr := errors.New("Product not deleted in database")
	if !reflect.DeepEqual(err, expErr) {
		t.Errorf("Expected %v error but got %v", expErr, err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 5, res1)
	}

	//TestCase-3
	mock.ExpectExec("delete from product").WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 0))
	res1, err = pS.Delete(5)
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil error but got %v", err)
	}
	if !reflect.DeepEqual(res1, 0) {
		t.Errorf("Expected last inserted id as %v but got %v", 0, res1)
	}

}
