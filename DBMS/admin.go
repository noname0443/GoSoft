package DBMS

import (
	"GoSoft/graph/model"
	"errors"
)

func StoreAdd(product model.NewProduct) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
INSERT INTO public.store(
	name, description, photo, file, price)
	VALUES ($1, $2, $3, $4, $5);`, product.Name, product.Description, product.Photo, product.File, product.Price)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("can't add product")
	}
	return nil
}

func StoreRemove(productid int) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
DELETE FROM store
WHERE productid = $1;`, productid)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("can't remove product")
	}
	return nil
}

func StoreUpdate(productid int, product model.NewProduct) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
UPDATE store
	SET name=$1, description=$2, photo=$3, file=$4, price=$5
	WHERE productid = $6;`, product.Name, product.Description, product.Photo, product.File, product.Price, productid)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("can't update product")
	}
	return nil
}
