package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"goOrderAPI/logger"
	"googlemaps.github.io/maps"
	"os"
)

var validate *validator.Validate
var fieldValidate *validator.Validate

type Order struct {
	Base
	Origin      pq.StringArray `gorm:"not null;type:varchar(100)[]" validate:"len=2,required,dive,gt=0,required"`
	Destination pq.StringArray `gorm:"not null;type:varchar(100)[]" validate:"len=2,required,dive,gt=0,required"`
	Status      string         `json:"status"`
	Distance    int            `json:"distance"`
}

func (order *Order) Validate() error {
	validate = validator.New()
	if err := validate.Struct(order); err != nil {
		return err
	}
	if err := validate.Var(order.Origin[0], "latitude"); err != nil {
		return errors.New("invalid origin latitude")
	}
	if err := validate.Var(order.Origin[1], "longitude"); err != nil {
		return errors.New("invalid origin longitude")
	}
	if err := validate.Var(order.Destination[0], "latitude"); err != nil {
		return errors.New("invalid destination latitude")
	}
	if err := validate.Var(order.Destination[1], "longitude"); err != nil {
		return errors.New("invalid destination longitude")
	}
	return nil
}

func (order *Order) Create() (interface{}, error) {
	if err := GetDB().Create(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (order *Order) Update(data map[string]interface{}) (interface{}, error) {
	if err := GetDB().Model(order).Updates(data).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (order *Order) Delete() error {
	return nil
}
func (order *Order) Get(forUpdate bool) (interface{}, error) {
	fetchedOrder := &Order{}
	db := GetDB()
	if forUpdate {
		db = db.Set("gorm:query_option", "FOR UPDATE")
	}
	if err := db.Where(order).First(&fetchedOrder).Error; err != nil {
		return nil, err
	}
	return fetchedOrder, nil
}

func (order *Order) All(offset int, limit int) (interface{}, error) {
	orders := []*Order{}
	if err := GetDB().Offset((offset - 1) * limit).Limit(limit).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (order *Order) BeforeSave() error {
	googleAPIKey := os.Getenv("googleAPIKey")
	c, err := maps.NewClient(maps.WithAPIKey(googleAPIKey))
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Cannot use google API")
		return err
	}

	r := &maps.DistanceMatrixRequest{
		Origins:      []string{fmt.Sprintf("%s,%s", order.Origin[0], order.Origin[1])},
		Destinations: []string{fmt.Sprintf("%s,%s", order.Destination[0], order.Destination[1])},
		Units:        maps.UnitsMetric,
	}

	distance, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Unable to get distance")
		return err
	}

	if distance.Rows[0].Elements[0].Status != "OK" {
		return errors.New("Cannot calculate distance, invalid combination of coordinates")
	}
	order.Distance = distance.Rows[0].Elements[0].Distance.Meters
	order.Status = "UNASSIGNED"
	return nil
}

func (order *Order) Filter() (interface{}, error) {
	return nil, nil
}
