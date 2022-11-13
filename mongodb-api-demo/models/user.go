package models

type Address struct {
    Street     string `json:"street" bson:"street"`
    City       string `json:"city" bson:"city"`
    PostalCode int    `json:"postal_code" bson:"postal_code"`
}

type User struct {
    Name    string  `json:"name" bson:"user_name"`
    Age     int     `json:"age" bson:"user_age"`
    Address Address `json:"address" bson:"user_address"`
}
