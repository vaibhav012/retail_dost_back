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
	"crypto/sha256"
    "encoding/base64"
	"bytes"
)

func CreateStore(w http.ResponseWriter, req *http.Request) {
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
	var store = result["store"].(map[string]interface{})

	name := store["name"].(string)
	address := store["address"].(string)
	store_phone := store["phone"].(string)
	// email := store["email"].(string)
	title := store["title"].(string)
	store_type := store["type"].(string)
	categories := store["categories"].(string)

	var query = "SELECT * FROM users WHERE phone='" + phone + "'";
	results := vvdatabase.DBCon.QueryRow(query)
    var user User
    err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			query := fmt.Sprintf("INSERT INTO stores (`name`, `owner`, `address`, `phone`, `title`, `type`, `categories`) VALUES ('%s',%d,'%s','%s','%s','%s','%s')", name, user.Id, address, store_phone, title, store_type, categories)
			fmt.Println(query)
			insert, err := vvdatabase.DBCon.Query(query)
			// resultsInsert, err := vvdatabase.DBCon.Exec(name, user.Id, address, store_phone, email, title, store_type, categories)
			if err != nil {
				panic(err)
			}
			fmt.Println(insert)
			response.ResponseCode = 0
			response.Response = "success"
			response.ResponseMessage = "Store Added"
			json.NewEncoder(w).Encode(response)
	}
}

func GetStores(w http.ResponseWriter, req *http.Request){
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
			var stores []Store
			var query = fmt.Sprintf("SELECT * FROM stores WHERE owner=%d", user.Id)
			fmt.Println(query);
			results, err := vvdatabase.DBCon.Query(query)
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
				fmt.Printf("len=%d cap=%d %v\n", len(stores), cap(stores), stores)
			}
			response.ResponseCode = 0
			response.Response = "success"
			storesString, err := json.Marshal(stores)
			if err != nil {
				panic(err)
			}
			response.ResponseMessage = string(storesString)
			response.ResponseData = stores
			json.NewEncoder(w).Encode(response)
	}
}

func StoreProfile(w http.ResponseWriter, req *http.Request){
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
			var storeProducts []StoreProduct = make([]StoreProduct, 0)
			var query = fmt.Sprintf("SELECT store_products.id, store_products.MRP,  store_products.SKU,  store_products.description,  store_products.images,  store_products.price, product.Id, product.Name, product.Category FROM store_products JOIN product ON store_products.product = product.id WHERE store=%d", store)
			fmt.Println(query);
			results, err := vvdatabase.DBCon.Query(query)
			if err != nil {
				panic(err.Error())
			}
			for results.Next() {
				var storeProduct StoreProduct
				err = results.Scan(&storeProduct.Id, &storeProduct.MRP, &storeProduct.SKU, &storeProduct.Description, &storeProduct.Images, &storeProduct.Price, &storeProduct.Product.Id, &storeProduct.Product.Name, &storeProduct.Product.Category)
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
 			response.ResponseMessage = "storeProfile"
			response.ResponseData = storeProducts
			json.NewEncoder(w).Encode(response)
	}
}

func CreateStoreProduct(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	newStr := buf.String()
	var result map[string]interface{}
	json.Unmarshal([]byte(newStr), &result)

	phone := result["phone"].(string)
	pass := result["pass"].(string)

	err = vvdatabase.DBCon.Ping()
	if err != nil {
		fmt.Println(err)
	}

	var query = "SELECT * FROM users WHERE phone=" + phone + " AND password='" + pass + "'";
	results := vvdatabase.DBCon.QueryRow(query)

    var user User
    err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
			hasher := sha256.New()
		    hasher.Write([]byte(phone))
		    token := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			query = "INSERT INTO sessions_user (token, user) VALUES (?,?)"
			resultsInsert, err := vvdatabase.DBCon.Exec(query, token, user.Id)
			if err != nil {
		        panic(err)
		    }
			fmt.Println(resultsInsert)
			user.Token = token
			response.ResponseCode = 0
			response.Response = "success"
			userString, err := json.Marshal(user)
			if err != nil {
		        panic(err)
		    }
			response.ResponseMessage = string(userString)
			response.ResponseData = user
			json.NewEncoder(w).Encode(response)
		default:
			var response Response
			response.ResponseCode = 1
			response.Response = "failure"
			response.ResponseMessage = "Invalid Credentials"
			json.NewEncoder(w).Encode(response)
	}
}
