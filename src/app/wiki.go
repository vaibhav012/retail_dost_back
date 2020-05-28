package main

import (
	"fmt"
	"log"
    "net/http"
    "github.com/gorilla/mux"
	"encoding/json"
	"api_functions"
	"vvdatabase"
	// "io/ioutil"
	// "html/template"
)

type Article struct {
    Title string `json:"Title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
}

type Tag struct {
    id   int    `json:"id"`
    name string `json:"name"`
    address string `json:"address"`
    contact string `json:"contact"`
}

var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/appdata", returnAllArticles)
    myRouter.HandleFunc("/login", api_functions.Login)
    myRouter.HandleFunc("/logout", api_functions.Login)
    myRouter.HandleFunc("/store/create", api_functions.CreateStore)
    myRouter.HandleFunc("/store/get", api_functions.GetStores)
    myRouter.HandleFunc("/account/profile", api_functions.GetProfile)
    myRouter.HandleFunc("/product/add", api_functions.CreateProduct)
    myRouter.HandleFunc("/product/get", api_functions.Login)
    myRouter.HandleFunc("/product/profile", api_functions.ProductProfile)
	myRouter.HandleFunc("/store/profile", api_functions.StoreProfile)
	myRouter.HandleFunc("/store/product/get", api_functions.GetProductsInStore)
    myRouter.HandleFunc("/category/add", api_functions.CreateProduct)
    myRouter.HandleFunc("/category/get", api_functions.GetCategory)
	myRouter.HandleFunc("/search", api_functions.Search)
    log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)
}

func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []Article{
        Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
    }

	vvdatabase.ConnectDatabase()

    handleRequests()
}
