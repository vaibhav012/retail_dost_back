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
	"strconv"
)

func CreateProduct(w http.ResponseWriter, req *http.Request) {
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
	store := int(result["store"].(float64))
	var product = result["product"].(map[string]interface{})

	name := product["name"].(string)
	category, err := strconv.Atoi(product["category"].(string))
	MRP := int(product["MRP"].(float64))
	images := "none";

	// images := product["images"].(string)
	description := product["description"].(string)
	SKU := product["SKU"].(string)
	price := int(product["price"].(float64))
	market_price := price

	var query = "SELECT * FROM users WHERE phone='" + phone + "'";
	results := vvdatabase.DBCon.QueryRow(query)
    var user User
    err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			query, err := vvdatabase.DBCon.Prepare("INSERT INTO product (`name`, `category`, `MRP`, `SKU`, `images`, `market_price`, `description`) VALUES (?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			insert, err := query.Exec(name, category, MRP, SKU, images, market_price, description)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			productId, err := insert.LastInsertId()

			query, err = vvdatabase.DBCon.Prepare("INSERT INTO store_products (`store`, `product`, `SKU`, `images`, `price`, `MRP`) VALUES (?, ?, ?, ?, ?, ?)") // ? = placeholder
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}
			insert, err = query.Exec(store, productId, SKU, images, price, MRP)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			response.ResponseCode = 0
			response.Response = "success"
			response.ResponseMessage = "Product Added"
			json.NewEncoder(w).Encode(response)
	}
}

func GetProductsInStore(w http.ResponseWriter, req *http.Request){
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
	store := int(result["store"].(float64))

	var query = "SELECT * FROM users WHERE phone='" + phone + "'";
	results := vvdatabase.DBCon.QueryRow(query)
	var user User
	err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			var storeProducts []StoreProduct
			var query = fmt.Sprintf("SELECT store_products.id, product.name, product.category, product.description, store_products.MRP,  store_products.SKU,  store_products.images,  store_products.price FROM store_products JOIN product ON store_products.product = product.id WHERE store=%d", store)
			fmt.Println(query);
			results, err := vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var storeProduct StoreProduct
				err = results.Scan(&storeProduct.Id, &storeProduct.Name, &storeProduct.Category, &storeProduct.Description, &storeProduct.MRP, &storeProduct.SKU, &storeProduct.Images, &storeProduct.Price)
				if err != nil {
					panic(err.Error())
				}
				storeProducts = append(storeProducts, storeProduct)
			}
			response.ResponseCode = 0
			response.Response = "success"
			storesString, err := json.Marshal(storeProducts)
			if err != nil {
				panic(err)
			}
			response.ResponseMessage = string(storesString)
			response.ResponseData = storeProducts
			json.NewEncoder(w).Encode(response)
	}
}

func ProductProfile(w http.ResponseWriter, req *http.Request){
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
	product := int(result["product"].(float64))

	var query = "SELECT * FROM users WHERE phone='" + phone + "'";
	results := vvdatabase.DBCon.QueryRow(query)
	var user User
	err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			var storeProducts []StoreProduct = make([]StoreProduct, 0)
			var query = fmt.Sprintf("SELECT store_products.id, store_products.MRP,  store_products.SKU,  store_products.description,  store_products.images,  store_products.price, store_products.store, stores.name, stores.title, stores.email, stores.phone, stores.address FROM store_products JOIN stores ON store_products.store = stores.id WHERE product=%d", product)
			fmt.Println(query);
			results, err := vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var storeProduct StoreProduct
				err = results.Scan(&storeProduct.Id, &storeProduct.MRP, &storeProduct.SKU, &storeProduct.Description, &storeProduct.Images, &storeProduct.Price, &storeProduct.Store.Id, &storeProduct.Store.Name, &storeProduct.Store.Title, &storeProduct.Store.Email, &storeProduct.Store.Phone, &storeProduct.Store.Address)
				if err != nil {
					panic(err.Error())
				}
				storeProducts = append(storeProducts, storeProduct)
			}
			response.ResponseCode = 0
			response.Response = "success"
			// storesString, err := json.Marshal(storeProducts)
			if err != nil {
				panic(err)
			}
 			response.ResponseMessage = "productProfile"
			response.ResponseData = storeProducts
			json.NewEncoder(w).Encode(response)
	}
}
