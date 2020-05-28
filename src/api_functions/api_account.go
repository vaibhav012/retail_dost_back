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

func Login(w http.ResponseWriter, req *http.Request) {
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

func GetProfile(w http.ResponseWriter, req *http.Request) {
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
	token := result["token"].(string)

	err = vvdatabase.DBCon.Ping()
	if err != nil {
		fmt.Println(err)
	}

	var query = "SELECT users.* FROM sessions_user JOIN users ON sessions_user.user = users.id WHERE users.phone=" + phone + " AND token='" + token + "'";
	results := vvdatabase.DBCon.QueryRow(query)

    var user User
    err = results.Scan(&user.Id, &user.Name, &user.Phone, &user.Email, &user.City, &user.Password, &user.Active)
	switch err {
		case nil:
			var response Response
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
			// response.ResponseMessage = "Invalid Account"
			response.ResponseMessage = err.Error()
			json.NewEncoder(w).Encode(response)
	}
}
