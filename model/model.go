package model

import "time"

type PlayerInfo struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Team         string `json:"team"`
	Nationality  string `json:"nationality"`
	Height       string `json:"height"`
	Born         string `json:"born"`
	Shirt_number string `json:"shirt_number"`
}

type Player struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerNation struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

type TeamInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Coach     string `json:"coach"`
	President string `json:"president"`
	Website   string `json:"website"`
}

type Team struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TimeOfCreation struct {
	Id   string    `json:"id"`
	Date time.Time `json:"date"`
}

type Statistic struct {
	Country_name      string `json:"country_name"`
	Number_of_players int    `json:"number_of_players"`
}

//type ResponseData struct {
//	Page           int    `json:"page"`
//	PageSize       int    `json:"page_size"`
//	TotalItems     int    `json:"total_items"`
//	OrderBy        string `json:"order_by"`
//	OrderDirection string `json:"order_direction"`
//	FilterBy       string `json:"filter_by"`
//	FilterValue    string `json:"filter_value"`
//	FilterOperator string `json:"filter_operator"`
//	Data           []Team `json:"data"`
//}

type ResponseData struct {
	Page           int         `json:"page"`
	PageSize       int         `json:"page_size"`
	TotalItems     int         `json:"total_items"`
	OrderBy        string      `json:"order_by"`
	OrderDirection string      `json:"order_direction"`
	FilterBy       string      `json:"filter_by"`
	FilterValue    string      `json:"filter_value"`
	FilterOperator string      `json:"filter_operator"`
	Data           interface{} `json:"data"`
}

type UrlParams struct {
	Page           int
	PageSize       int
	OrderBy        string
	OrderDirection string
	FilterBy       string
	FilterValue    string
	FilterOperator string
}
