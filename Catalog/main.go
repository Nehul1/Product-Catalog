package main

import (
	"database/sql"
	product3 "exercises/Catalog/handler/product"
	product2 "exercises/Catalog/service/product"
	"exercises/Catalog/store/brand"
	"exercises/Catalog/store/product"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "nehul:9618181838@tcp(127.0.0.1)/catalogue")
	defer db.Close()
	if err != nil {
		log.Println(err)
	}
	pS := product.New(db)
	bS := brand.New(db)
	catServ := product2.New(pS, bS)
	catHandle := product3.New(catServ)
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/product", catHandle.Get).Methods(http.MethodGet)
	myRouter.HandleFunc("/product", catHandle.Create).Methods(http.MethodPost)
	myRouter.HandleFunc("/product", catHandle.Update).Methods(http.MethodPut)
	myRouter.HandleFunc("/product", catHandle.Delete).Methods(http.MethodDelete)
	err = http.ListenAndServe(":8080", myRouter)
	if err != nil {
		log.Println(err)
	}
}
