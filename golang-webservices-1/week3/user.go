package main

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:company`
	Country  string   `json:country`
	Email    string   `json:"email"`
	Job      string   `json:job`
	Name     string   `json:"name"`
	Phone    string   `json:phone`
}
