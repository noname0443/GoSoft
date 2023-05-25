package DBMS

import (
	"GoSoft/graph/model"
	"errors"
	_ "errors"
)

func SearchProducts(name string, preparedCategory string, lowerPrice float64, highestPrice float64) ([]*model.Product, error) {
	checkConnection()
	rows, err := PostgreSQL.Query(`SELECT DISTINCT store.productid, name, description, photo, file, price FROM public.store INNER JOIN categories ON store.productid = categories.productid WHERE name LIKE $1 AND price >= $2 AND price <= $3 AND (LENGTH($4) = 0 OR category = $4);`, "%" + name + "%", lowerPrice, highestPrice, preparedCategory)
	if err != nil{
		return nil, err
	}

	var products []*model.Product
	for rows.Next(){
		p := new(model.Product)
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Photo, &p.File, &p.Price)
		if err != nil{
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProduct(id int) (*model.Product, error){
	checkConnection()
	rows, err := PostgreSQL.Query(`SELECT * FROM store WHERE productid = $1`, id)
	if err != nil{
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New("Product not found")
	}
	p := new(model.Product)
	err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Photo, &p.File, &p.Price)
	if err != nil{
		return nil, err
	}
	return p, nil
}