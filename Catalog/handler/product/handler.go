package product

import (
	"encoding/json"
	"exercises/Catalog/model"
	"exercises/Catalog/service"
	"fmt"
	"net/http"
	"strconv"
)

type prodHandler struct {
	prodService service.Product
}

func New(prodService service.Product) prodHandler {
	return prodHandler{
		prodService: prodService,
	}
}
func (c prodHandler) Get(w http.ResponseWriter,r *http.Request){
	query:=r.URL.Query()
	if len(query)==0 {
		c.GetAll(w,r)
		return
	}
	c.GetById(w,r)
}

func (c prodHandler) GetById(w http.ResponseWriter, r *http.Request) {
	res := r.URL.Query()["id"][0]
	key, err := strconv.Atoi(res)
	var respErr model.ErrResp
	if err != nil {
		respErr.RespCode=http.StatusBadRequest
		respErr.ErrorResp="Invalid ID"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	prodDet, err := c.prodService.GetById(key)
	if err != nil {
		respErr.RespCode=404
		respErr.ErrorResp=fmt.Sprintf("Product details with matching product id=%v not found",key)
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonRes)
		return
	}
	jsonRes, _ := json.Marshal(prodDet)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
func (c prodHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	prodDet, err := c.prodService.GetAll()
	var respErr model.ErrResp
	if err != nil {
		respErr.RespCode=404
		respErr.ErrorResp="Requested product details not found"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonRes)
		return
	}
	jsonRes, _ := json.Marshal(prodDet)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (c prodHandler) Create(w http.ResponseWriter, r *http.Request) {
	var pD model.Prod
	var respErr model.ErrResp
	body := r.Body
	err := json.NewDecoder(body).Decode(&pD)
	fmt.Println(err)
	if err != nil {
		respErr.RespCode=404
		respErr.ErrorResp="New Product not added to database"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	res, err := c.prodService.Create(pD.Name, pD.BrandDetails.Brand)
	if err != nil {
		respErr.RespCode=404
		respErr.ErrorResp="New Product not added to database"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	jsonRes, _ := json.Marshal(res)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)

}
func (c prodHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pD model.Prod
	var respErr model.ErrResp
	queryParams,ok:=r.URL.Query()["id"]
	if !ok{
		respErr.RespCode=400
		respErr.ErrorResp="Missing product ID"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	key, err := strconv.Atoi(queryParams[0])
	if err != nil {
		respErr.RespCode=400
		respErr.ErrorResp="Invalid ID"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	body := r.Body
	err = json.NewDecoder(body).Decode(&pD)
	if err != nil {
		respErr.RespCode=404
		respErr.ErrorResp="Product not updated in the database"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	res, err := c.prodService.Update(key,pD.Name, pD.BrandDetails.Brand)
	if err != nil {
		respErr.RespCode=404
		respErr.ErrorResp="Product not updated in the database"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	jsonRes, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}


func (c prodHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var respErr model.ErrResp
	queryParams,ok:=r.URL.Query()["id"]
	//fmt.Println(ok)
	if !ok{
		respErr.RespCode=400
		respErr.ErrorResp="Missing product ID"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	key, err := strconv.Atoi(queryParams[0])
	if err != nil {
		respErr.RespCode=400
		respErr.ErrorResp="Invalid ID"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	resErr := c.prodService.Delete(key)

	if resErr != nil {
		respErr.RespCode=404
		respErr.ErrorResp="Product not deleted in the database"
		jsonRes, _ := json.Marshal(respErr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
