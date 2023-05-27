package DBMS

import (
	"GoSoft/graph/model"
	"errors"
	_ "errors"
)

func SearchProducts(name string, preparedCategory string, lowerPrice float64, highestPrice float64) ([]*model.Product, error) {
	checkConnection()
	rows, err := PostgreSQL.Query(`
SELECT store.productid, name, description, photo, file, price, subscriptiontype FROM public.store
WHERE (LENGTH($4) = 0 OR productid in
	   (SELECT productid FROM categories WHERE category = $4)) AND
	   name LIKE $1 AND price >= $2 AND price <= $3;
`, "%"+name+"%", lowerPrice, highestPrice, preparedCategory)
	if err != nil {
		return nil, err
	}

	var products []*model.Product
	for rows.Next() {
		p := new(model.Product)
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Photo, &p.File, &p.Price, &p.Subscriptiontype)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProduct(id int) (*model.Product, error) {
	checkConnection()
	rows, err := PostgreSQL.Query(`SELECT * FROM store WHERE productid = $1`, id)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New("Product not found")
	}
	p := new(model.Product)
	err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Photo, &p.File, &p.Price, &p.Subscriptiontype)
	if err != nil {
		return nil, err
	}
	return p, nil
}
