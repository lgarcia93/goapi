package repository

import (
	"fitgoapi/database"
	"fitgoapi/model"
)

// ICityRepository ...
type ICityRepository interface {
	FetchCityByCode(code string) ([]model.City, error)
	FetchCities() ([]model.City, error)
}

// CityRepository ...
type CityRepository struct {
}

func (r CityRepository) FetchCityByCode(code string) ([]model.City, error) {
	db := database.Connection

	var city []model.City

	db.Query("")

	return city, nil
}

func (r CityRepository) FetchCities() ([]model.City, error) {

	db := database.Connection

	var cities []model.City

	results, err := db.Query("SELECT code, name, zip_code, uf FROM city")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var city model.City
		err = results.Scan(&city.Code, &city.Name, &city.ZipCode, &city.UF)
		if err != nil {
			panic(err.Error())
		}
		cities = append(cities, city)
	}

	return cities, err
}
