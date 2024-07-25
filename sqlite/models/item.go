package models

type Item struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type CreateItem struct {
	Description string `json:"description"`
}
