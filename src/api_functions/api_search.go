package api_functions

import (
	"fmt"
	_ "log"
	// "net/http"
    // "github.com/gorilla/mux"
	"encoding/json"
	// "io/ioutil"
	// "html/template"
    "net/http"
	"vvdatabase"
	"bytes"
)

func Search(w http.ResponseWriter, req *http.Request){
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	err = vvdatabase.DBCon.Ping()
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	newStr := buf.String()
	var result map[string]interface{}
	json.Unmarshal([]byte(newStr), &result)
	phone := result["phone"].(string)
	search := result["search"].(string)

	var query = "SELECT * FROM users WHERE phone='" + phone + "'";
	results := vvdatabase.DBCon.QueryRow(query)
	var user User
	err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			// var storeProducts []StoreProduct = make([]StoreProduct, 0)
			var stores []Store = make([]Store, 0)
			var products []Product = make([]Product, 0)
			var productCategories []ProductCategory = make([]ProductCategory, 0)

			var query = "SELECT * FROM product_categories WHERE category LIKE '%" + search + "%'"
			fmt.Println(query);
			results, err := vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var productCategory ProductCategory
				err = results.Scan(&productCategory.Id, &productCategory.Category, &productCategory.ParentCategory, &productCategory.HeadCategory)
				if err != nil {
					panic(err.Error())
				}
				productCategories = append(productCategories, productCategory)
			}

			query = "SELECT * FROM product WHERE name LIKE '%" + search + "%' OR description LIKE '%" + search + "%'"
			fmt.Println(query);
			results, err = vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var product Product
				err = results.Scan(&product.Id, &product.Name, &product.Category, &product.MRP, &product.SKU, &product.Images, &product.MarketPrice, &product.Description)
				if err != nil {
					panic(err.Error())
				}
				products = append(products, product)
			}

			query = "SELECT * FROM stores WHERE name LIKE '%" + search + "%' OR title LIKE '%" + search + "%'"
			fmt.Println(query);
			results, err = vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var store Store
				err = results.Scan(&store.Id, &store.Name, &store.OwnerRef, &store.Address, &store.Phone, &store.Email, &store.Title, &store.Type, &store.Categories, &store.CurrentPlan)
				if err != nil {
					panic(err.Error())
				}
				stores = append(stores, store)
			}

			var searchResponse SearchResponse
			searchResponse.Stores = stores
			searchResponse.Products = products
			searchResponse.ProductCategories = productCategories
			response.ResponseCode = 0
			response.Response = "success"
			// someString, err := json.Marshal(productCategories)
			if err != nil {
				panic(err)
			}
			response.ResponseData = searchResponse
			json.NewEncoder(w).Encode(response)
	}
}
