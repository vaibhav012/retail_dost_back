package api_functions

import (
	"database/sql"
)

type Response struct {
	ResponseCode   int    `json:"responseCode"`
    Response string `json:"response"`
    ResponseMessage string `json:"responseMessage"`
    ResponseData interface{} `json:"responseData"`
}

type SearchResponse struct {
	Stores   interface{}    `json:"stores"`
    Products interface{} `json:"products"`
    ProductCategories interface{} `json:"categories"`
}

type User struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
    Phone string `json:"phone"`
    Email string `json:"email"`
    City string `json:"city"`
    Password string `json:"password"`
	Active string `json:"active"`
	Token string `json:"token"`
}

type Product struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
    Category string `json:"category"`
    MRP string `json:"MRP"`
	SKU string `json:"SKU"`
    Images string `json:"images"`
	MarketPrice string `json:"market_price"`
	Description string `json:"description"`
	Price string `json:"price"`
}

type ProductCategory struct {
    Id   int    `json:"id"`
    Category string `json:"category"`
    ParentCategory sql.NullString `json:"parent_category"`
    HeadCategory sql.NullString `json:"head_category"`
}

type Store struct {
    Owner User `json:"owner"`
    Id   int    `json:"id"`
    Name string `json:"name"`
	OwnerRef int `json:"owner_ref"`
    Address string `json:"address"`
    Phone string `json:"phone"`
    Email string `json:"email"`
	Title string `json:"title"`
	Type string `json:"type"`
	Categories string `json:"categories"`
	CurrentPlan string `json:"current_plan"`
}

type StoreProduct struct {
    Id int `json:"id"`
	Name string  `json:"name"`
	Category int `json:"category"`
	Description string `json:"description"`
	MRP string `json:"MRP"`
    SKU string `json:"SKU"`
    Price string `json:"price"`
    Images string `json:"images"`
	Store Store `json:"store"`
	Product Product `json:"product"`
	StoreName string `json:"storeName"`
	StoreTitle string `json:"storeTitle"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Address string `json:"address"`
}
