package models

type Item struct {
	ID          uint
	Name        string
	Price       uint
	Quantity    uint
	Description string
	SoldOut     bool
}
