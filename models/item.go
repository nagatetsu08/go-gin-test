package models

import "gorm.io/gorm"

type Item struct {
	// gorm.Modelにカーソルを当てると、内部で管理しているパラメータがみえる。（IDとかCreatedAtとか）
	// このように構造体in構造体でモデルを定義できる
	gorm.Model
	Name        string `gorm:"not null"` // タグでDBの制約をつけることができる
	Price       uint   `gorm:"not null"`
	Quantity    uint
	Description string
	SoldOut     bool `gorm:"not null;default:false"` //複数定義するときはセミコロンで区切る。ただし、スペースとかいれてはいけない
	UserId      uint `gorm:"not null"`
}
