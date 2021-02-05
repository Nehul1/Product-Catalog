package product

import (
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/store"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestProdService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := store.NewMockProduct(ctrl)
	bS := store.NewMockBrand(ctrl)
	testcases := []struct {
		expIn   int
		pOut    model.Prod
		bOut    model.Brand
		expOut  model.Prod
		expPerr error
		expBerr error
		expErr  error
	}{
		{expIn: 1, pOut: model.Prod{Id: 1, Name: "Shoes",BrandDetails: model.Brand{3, ""}}, bOut: model.Brand{3, "Puma"}, expOut: model.Prod{
			1, "Shoes", model.Brand{3, "Puma"}},expPerr: nil,expBerr: nil,expErr: error(nil)},
		{expIn: 3, pOut: model.Prod{Id: 3, Name: "Cricket Shoes",BrandDetails: model.Brand{5, ""}}, bOut: model.Brand{},expOut: model.Prod{
			3, "Cricket Shoes", model.Brand{5, ""}}, expPerr: nil,expBerr: errors.New("Brand Id not found"), expErr: errors.New("Brand Id not found")},
		{expIn: 5, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{},expPerr: errors.New("Id not found"), expBerr: errors.New("Id not found"),expErr: errors.New("Id not found")},
	}
	for i, tc := range testcases {
		pS.EXPECT().GetById(tc.expIn).Return(tc.pOut, tc.expPerr)
		if tc.expPerr == nil {
			bS.EXPECT().GetById(tc.pOut.BrandDetails.Id).Return(tc.bOut, tc.expBerr)
		}
		prodService := New(pS, bS)
		res, err := prodService.GetById(tc.expIn)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestProdService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := store.NewMockProduct(ctrl)
	bS := store.NewMockBrand(ctrl)
	testcases := []struct {
		pOut    []model.Prod
		bOut    []model.Brand
		expOut  []model.Prod
		expPerr error
		expBerr error
		expErr  error
	}{
		{ pOut: []model.Prod{{Id: 1, Name: "Shoes",BrandDetails: model.Brand{Id: 3,Brand: ""}},{Id: 2,Name: "Cricket Shoes", BrandDetails: model.Brand{Id: 2, Brand: ""}}}, bOut: []model.Brand{{3, "Puma"},{2,"Nike"}}, expOut: []model.Prod{{
			1, "Shoes", model.Brand{3, "Puma"}},{2, "Cricket Shoes", model.Brand{2, "Nike"}}}, expPerr: nil,expBerr: nil, expErr: nil},
		{pOut: []model.Prod{}, bOut: []model.Brand{},expOut: []model.Prod{},expPerr:  errors.New("Products not found"),expBerr: errors.New("Brand Id not Found"),expErr: errors.New("Products not found")},
		{ pOut: []model.Prod{{1, "Shoes", model.Brand{3, ""}}}, bOut: []model.Brand{{}},expOut: []model.Prod{{
			1, "Shoes", model.Brand{3, ""}}}, expPerr: nil,expBerr: errors.New("Brand Id not Found"),  expErr: errors.New("Brand Id not found")},

	}
	for i, tc := range testcases {
		pS.EXPECT().GetAll().Return(tc.pOut, tc.expPerr)
		if tc.expPerr == nil {
			for i,v:=range tc.pOut {
				bS.EXPECT().GetById(v.BrandDetails.Id).Return(tc.bOut[i], tc.expBerr)
			}
		}
		prodService := New(pS, bS)
		res, err := prodService.GetAll()
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestProdService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := store.NewMockProduct(ctrl)
	bS := store.NewMockBrand(ctrl)
	testcases := []struct {
		expCheckBrand     model.Brand
		expCheckBrandErr  error
		expCreateBrand    model.Brand
		expBrandCreateErr error
		expProdCreate     model.Prod
		expProdCreateErr  error
		pOut              model.Prod
		bOut              model.Brand
		expOut            model.Prod
		expPerr           error
		expBerr           error
		expErr            error
	}{
		{expProdCreate: model.Prod{1, "Shoes", model.Brand{3, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{3, "Puma"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Puma"}, expBrandCreateErr: nil, pOut: model.Prod{1, "Shoes", model.Brand{3, ""}}, bOut: model.Brand{3, "Puma"}, expOut: model.Prod{
			1, "Shoes", model.Brand{3, "Puma"}}, expPerr: nil, expBerr: nil, expErr: nil},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: nil, pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, "Asics"}}, bOut: model.Brand{}, expOut: model.Prod{
			2, "Cricket Shoes", model.Brand{5, "Asics"}}, expPerr: nil, expBerr: errors.New("Brand Id not found"), expErr: errors.New("Brand not found")},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: errors.New("Brand not created"), pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: nil, expBerr: errors.New("Brand Id not found"),expErr: errors.New("Brand not created")},
		{expProdCreate: model.Prod{3, "Slippers", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Crocs"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Crocs"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Created Product not found")},
		{expProdCreate: model.Prod{0, "Sandals", model.Brand{3, ""}}, expProdCreateErr: errors.New("Employee not inserted"), expCheckBrand: model.Brand{3, "Mochi"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Mochi"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Product not created")},
	}

	for i, tc := range testcases {
		bS.EXPECT().CheckBrand(tc.expCheckBrand.Brand).Return(tc.expCheckBrand.Id, tc.expCheckBrandErr)
		if tc.expCheckBrand.Id == 0 {
			bS.EXPECT().Create(tc.expCreateBrand.Brand).Return(tc.expCreateBrand.Id, tc.expBrandCreateErr)
		}
		pS.EXPECT().Create(tc.expProdCreate.Name, tc.expProdCreate.BrandDetails.Id).Return(tc.expProdCreate.Id, tc.expProdCreateErr)
		if tc.expProdCreateErr == nil {
			pS.EXPECT().GetById(tc.expProdCreate.Id).Return(tc.pOut, tc.expPerr)
			if tc.expPerr == nil {
				bS.EXPECT().GetById(tc.expProdCreate.BrandDetails.Id).Return(tc.bOut, tc.expBerr)
			}
		}
		prodService := New(pS, bS)
		res, err := prodService.Create(tc.expProdCreate.Name, tc.expCheckBrand.Brand)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestProdService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := store.NewMockProduct(ctrl)
	bS := store.NewMockBrand(ctrl)
	testcases := []struct {
		expCheckBrand     model.Brand
		expCheckBrandErr  error
		expCreateBrand    model.Brand
		expBrandCreateErr error
		expProdCreate     model.Prod
		expProdCreateErr  error
		pOut              model.Prod
		bOut              model.Brand
		expOut            model.Prod
		expPerr           error
		expBerr           error
		expErr            error
	}{
		{expProdCreate: model.Prod{1, "Shoes", model.Brand{3, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{3, "Puma"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Puma"}, expBrandCreateErr: nil, pOut: model.Prod{1, "Shoes", model.Brand{3, ""}}, bOut: model.Brand{3, "Puma"}, expOut: model.Prod{
			1, "Shoes", model.Brand{3, "Puma"}}, expPerr: nil, expBerr: nil, expErr: nil},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: nil, pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, "Asics"}}, bOut: model.Brand{}, expOut: model.Prod{
			2, "Cricket Shoes", model.Brand{5, "Asics"}}, expPerr: nil, expBerr: errors.New("Brand Id not found"), expErr: errors.New("Updated Brand not found")},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: errors.New("Brand not created"), pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: nil, expBerr: errors.New("Brand Id not found"),expErr: errors.New("Brand not created")},
		{expProdCreate: model.Prod{3, "Slippers", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Crocs"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Crocs"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Updated Product not found")},
		{expProdCreate: model.Prod{0, "Sandals", model.Brand{3, ""}}, expProdCreateErr: errors.New("Employee not inserted"), expCheckBrand: model.Brand{3, "Mochi"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Mochi"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Product not updated")},
	}

	for i, tc := range testcases {
		bS.EXPECT().CheckBrand(tc.expCheckBrand.Brand).Return(tc.expCheckBrand.Id, tc.expCheckBrandErr)
		if tc.expCheckBrand.Id == 0 {
			bS.EXPECT().Create(tc.expCreateBrand.Brand).Return(tc.expCreateBrand.Id, tc.expBrandCreateErr)
		}
		pS.EXPECT().Update(tc.expProdCreate.Id,tc.expProdCreate.Name, tc.expProdCreate.BrandDetails.Id).Return(tc.expProdCreate.Id, tc.expProdCreateErr)
		if tc.expProdCreateErr == nil {
			pS.EXPECT().GetById(tc.expProdCreate.Id).Return(tc.pOut, tc.expPerr)
			if tc.expPerr == nil {
				bS.EXPECT().GetById(tc.expProdCreate.BrandDetails.Id).Return(tc.bOut, tc.expBerr)
			}
		}
		prodService := New(pS, bS)
		res, err := prodService.Update(tc.expProdCreate.Id,tc.expProdCreate.Name, tc.expCheckBrand.Brand)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestProdService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := store.NewMockProduct(ctrl)
	bS := store.NewMockBrand(ctrl)
	testcases := []struct {
		expIn int
		expErr            error
		expRet int
		expRetErr error
	}{
		{expIn: 1,expErr: nil,expRet: 1,expRetErr: nil},
		{expIn: 10,expErr: errors.New("Product not deleted from database"),expRet: 0,expRetErr: errors.New("Product not deleted from database")},
		{expIn: 5,expErr: errors.New("Product not deleted from database"),expRet: 0,expRetErr: nil},
	}

	for i, tc := range testcases {
		pS.EXPECT().Delete(tc.expIn).Return(tc.expRet,tc.expRetErr)
		prodService := New(pS, bS)
		resErr := prodService.Delete(tc.expIn)
		if !reflect.DeepEqual(resErr, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, resErr)
		}
		}
}
