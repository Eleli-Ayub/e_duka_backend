package models

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductID          string `gorm:"column:product_id;not null;primary key;unique;" json:"producttid"`
	ProductName        string `gorm:"column:product_name;unique;not null" json:"productname"`
	ProductPrice       string `gorm:"column:product_price;not null" json:"productprice"`
	ProductDescription string `gorm:"column:product_description;" json:"productdescription"`
	UserID             string `gorm:"size:255;not null" json:"userid"`
	MainImage          string `gorm:"not null;" json:"mainimage"`
	ProductStatus      string `gorm:"not null;" json:"productstatus"`
	Quantity           int    `gorm:"default:0" json:"quantity"`
	ProductType        string `gorm:"column:product_type;" json:"producttype"`
	TotalLikes         int    `gorm:"default:0" json:"totallikes"`
	TotalComments      int    `gorm:"default:0" json:"totalcomments"`
	DateAdded          string `gorm:"" json:"dateadded"`
	LastUpdated        string `gorm:"size:255;not null" json:"lastupdated"`
	LatestInteractions string `gorm:"size:255;not null" json:"latestinteractions"`
	TotalInteractions  int    `gorm:"size:255;not null" json:"totalinteractions"`
	TotalBookmarks     int    `gorm:"size:255;not null" json:"totalbookmarks"`
	Brand              string `gorm:"column:brand" json:"brand"`
	Category           string `gorm:"category" json:"category"`
	SubCategory        string `gorm:"column:subcategory" json:"subcategory"`
}

func (product *Product) Save() (*Product, error) {
	err := database.Database.Create(&product).Error
	if err != nil {
		return &Product{}, err
	}
	return product, nil
}
