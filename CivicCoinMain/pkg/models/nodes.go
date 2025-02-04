package models

import "time"

type Node struct {
	Id          string    `json:"id"`
	Addr        string    `json:"ip_address"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}
