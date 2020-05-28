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

func CreateCategory(w http.ResponseWriter, req *http.Request) {
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
	category := int(product["category"].(float64))
	MRP := int(product["MRP"].(float64))
	images := product["images"].(string)
	market_price := int(product["market_price"].(float64))
	description := product["description"].(string)
	SKU := product["SKU"].(string)
	price := int(product["price"].(float64))

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

func GetCategory(w http.ResponseWriter, req *http.Request){
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

	var query = "SELECT * FROM users WHERE phone='" + phone + "'";
	results := vvdatabase.DBCon.QueryRow(query)
	var user User
	err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			var categories []ProductCategory
			var query = fmt.Sprintf("SELECT * FROM product_categories WHERE 1")
			fmt.Println(query);
			results, err := vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var product_category ProductCategory
				err = results.Scan(&product_category.Id, &product_category.Category, &product_category.ParentCategory, &product_category.HeadCategory)
				if err != nil {
					panic(err.Error())
				}
				categories = append(categories, product_category)
			}
			response.ResponseCode = 0
			response.Response = "success"
			// storesString, err := json.Marshal(categories)
			if err != nil {
				panic(err)
			}
			response.ResponseData = categories
			json.NewEncoder(w).Encode(response)
	}
}
