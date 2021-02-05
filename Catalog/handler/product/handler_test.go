package product

import (
	"bytes"
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/service"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProdHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	prodService := service.NewMockProduct(ctrl)
	testcasesById := []struct {
		reqMethod string
		expIn     int
		reqPath   string
		expRespBody    string
		expServOut    model.Prod
		expRespCode   int
		expServErr    error
	}{
		{reqMethod: http.MethodGet, expIn: 1,reqPath: "/product?id=1", expRespBody: `{"id":1,"name":"Shoes","brand":{"name":"Puma"}}`,expServOut: model.Prod{1, "Shoes", model.Brand{3, "Puma"}}, expRespCode: http.StatusOK, expServErr: nil},
		{reqMethod: http.MethodGet,expIn: 5, reqPath: "/product?id=5", expRespBody: `{"responseCode":404,"errorResponse":"Product details with matching product id=5 not found"}`, expServOut: model.Prod{}, expRespCode: http.StatusNotFound, expServErr: errors.New("Id not found")},
		{reqMethod: http.MethodGet, expIn: 3,reqPath: "/product?id=abc", expRespBody: `{"responseCode":400,"errorResponse":"Invalid ID"}`,expServOut: model.Prod{3, "Cricket Shoes", model.Brand{5, ""}}, expRespCode: http.StatusBadRequest, expServErr: errors.New("Brand Id not found")},
	}
	for i, tc := range testcasesById {
		prodService.EXPECT().GetById(tc.expIn).Return(tc.expServOut, tc.expServErr)
		prodHandler := New(prodService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, nil)
		prodHandler.Get(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expRespCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expRespCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expRespBody) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expRespBody, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected nil error but got %v", i, resErr)
		}
	}
	testcases := []struct {
		reqMethod string
		reqPath   string
		expRespBody    string
		expServOut    []model.Prod
		expRespCode   int
		expServErr    error
	}{
		{"GET","/product", `[{"id":1,"name":"Shoes","brand":{"name":"Puma"}}]`, []model.Prod{{1, "Shoes", model.Brand{3, "Puma"}}}, 200, nil},
		{"GET","/product", `{"responseCode":404,"errorResponse":"Requested product details not found"}`, []model.Prod{}, 404, errors.New("Products not found")},

	}
	for i, tc := range testcases {
		prodService.EXPECT().GetAll().Return(tc.expServOut, tc.expServErr)
		prodHandler := New(prodService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, nil)
		prodHandler.Get(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expRespCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expRespCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expRespBody) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expRespBody, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected nil error but got %v", i, resErr)
		}
	}
}
func TestProdHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	prodService := service.NewMockProduct(ctrl)
	testcases := []struct {
		reqMethod  string
		expInName  string
		expInBName string
		reqPath    string
		reqBody    []byte
		expRespBody     string
		expServOut     model.Prod
		expRespCode    int
		expServErr error
	}{
		{reqMethod: http.MethodPost,expInName: "Cricket Shoes",expInBName: "Puma", reqPath: "/product", reqBody: []byte(`{"name":"Cricket Shoes","brand":{"name":"Puma"}}`),expRespBody: `{"id":1,"name":"Cricket Shoes","brand":{"name":"Puma"}}`, expServOut: model.Prod{1, "Cricket Shoes", model.Brand{3, "Puma"}}, expRespCode: http.StatusCreated, expServErr: nil},
		{reqMethod: http.MethodPost, expInName: "Shoes", expInBName: "Nike",reqPath: "/product",reqBody: []byte(`{"name":"Shoes","brand":{"name":"Nike"}}`), expRespBody: `{"responseCode":404,"errorResponse":"New Product not added to database"}`,expServOut: model.Prod{0, "", model.Brand{0, ""}},expRespCode: http.StatusBadRequest, expServErr: errors.New("Product not created")},
		//{"POST","Shoes","Nike","/product",[]byte(`{"names":"Shoes","brands":{"name":"Nike"}}`),"Product not added to database",model.Prod{0,"",model.Brand{0,""}},400,errors.New("Product not created")},

	}
	for i, tc := range testcases {
		prodService.EXPECT().Create(tc.expInName, tc.expInBName).Return(tc.expServOut, tc.expServErr)
		prodHandler := New(prodService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBuffer(tc.reqBody))
		prodHandler.Create(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expRespCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expRespCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expRespBody) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expRespBody, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected %v error but got %v", i, nil, resErr)
		}
	}
}

func TestProdHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	prodService := service.NewMockProduct(ctrl)
	testcases := []struct {
		reqMethod  string
		expInName  string
		expInBName string
		reqPath    string
		reqBody    []byte
		expRespBody     string
		expServOut     model.Prod
		expRespCode    int
		expServErr error
	}{
		{reqMethod: http.MethodPut, expInName: "Cricket Shoes", expInBName: "Puma", reqPath:"/product?id=1", reqBody: []byte(`{"name":"Cricket Shoes","brand":{"name":"Puma"}}`),expRespBody: `{"id":1,"name":"Cricket Shoes","brand":{"name":"Puma"}}`,expServOut: model.Prod{1, "Cricket Shoes", model.Brand{3, "Puma"}},expRespCode: http.StatusOK, expServErr: nil},
		{reqMethod: http.MethodPut,expInName: "Shoes",expInBName: "Nike",reqPath: "/product?id=0",reqBody: []byte(`{"name":"Shoes","brand":{"name":"Nike"}}`),expRespBody: `{"responseCode":404,"errorResponse":"Product not updated in the database"}`, expServOut: model.Prod{0, "", model.Brand{0, ""}}, expRespCode: http.StatusBadRequest, expServErr: errors.New("Product not updated")},
		{reqMethod: http.MethodPut,expInName: "Shoes",expInBName: "Nike",reqPath: "/product?id=abc",reqBody: []byte(`{"names":"Shoes","brands":{"name":"Nike"}}`),expRespBody: `{"responseCode":400,"errorResponse":"Invalid ID"}`,expServOut: model.Prod{0,"",model.Brand{0,""}},expRespCode: http.StatusBadRequest,expServErr: errors.New("Product not created")},
		//{"PUT", "Cricket Shoes", "Puma", "/product?id=1", []byte(`{"names":"Cricket Shoes","brands":{"names":"Puma"}}`), `{"id":1,"name":"Cricket Shoes","brand":{"name":"Puma"}}`, model.Prod{1, "Cricket Shoes", model.Brand{3, "Puma"}}, 200, nil},
		{reqMethod: http.MethodPut,expInName: "Cricket Shoes",expInBName: "Puma", reqPath: "/product",reqBody: []byte(`{"name":"Cricket Shoes","brand":{"name":"Puma"}}`), expRespBody: `{"responseCode":400,"errorResponse":"Missing product ID"}`, expServOut: model.Prod{1, "Cricket Shoes", model.Brand{3, "Puma"}}, expRespCode: http.StatusBadRequest,expServErr: nil},
	}
	for i, tc := range testcases {
		prodService.EXPECT().Update(tc.expServOut.Id,tc.expInName, tc.expInBName).Return(tc.expServOut, tc.expServErr)
		prodHandler := New(prodService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBuffer(tc.reqBody))
		prodHandler.Update(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expRespCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expRespCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expRespBody) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expRespBody, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected %v error but got %v", i, nil, resErr)
		}
	}
}

func TestProdHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	prodService := service.NewMockProduct(ctrl)
	testcases := []struct {
		reqMethod  string
		reqPath    string
		expIn int
		expRespBody string
		expServOutErr   error
		expRespCode    int
	}{
		{reqMethod: http.MethodDelete,reqPath: "/product?id=1",expIn: 1,expRespBody: "",expServOutErr: nil,expRespCode: http.StatusNoContent},
		{reqMethod: http.MethodDelete,reqPath: "/product?id=abc",expIn: 1,expRespBody: `{"responseCode":400,"errorResponse":"Invalid ID"}`,expServOutErr: nil,expRespCode: http.StatusBadRequest},
		{reqMethod: http.MethodDelete,reqPath: "/product?id=5",expIn: 5, expRespBody: `{"responseCode":404,"errorResponse":"Product not deleted in the database"}`, expServOutErr: errors.New("Product not deleted from database"),expRespCode: http.StatusBadRequest},
		{reqMethod: http.MethodDelete,reqPath: "/product",expIn: 1,expRespBody: `{"responseCode":400,"errorResponse":"Missing product ID"}`,expServOutErr: nil,expRespCode: http.StatusBadRequest},
	}
	for i, tc := range testcases {
		prodService.EXPECT().Delete(tc.expIn).Return(tc.expServOutErr)
		prodHandler := New(prodService)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBuffer(nil))
		prodHandler.Delete(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expRespCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expRespCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expRespBody) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expRespBody, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected %v error but got %v", i, nil, resErr)
		}
	}
}