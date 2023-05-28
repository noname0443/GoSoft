package DBMS

import (
	"GoSoft/graph/model"
	"errors"
)

func StoreAdd(product model.NewProduct) error {
	CheckConnection()
	result, err := PostgreSQL.Exec(`
INSERT INTO public.store(
	name, description, photo, file, price, subscriptiontype, company)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`, product.Name, product.Description, product.Photo, product.File, product.Price, product.Subscriptiontype, product.Company)
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
	CheckConnection()
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
	CheckConnection()
	result, err := PostgreSQL.Exec(`
UPDATE store
	SET name=$1, description=$2, photo=$3, file=$4, price=$5, company=$6
	WHERE productid = $7;`, product.Name, product.Description, product.Photo, product.File, product.Price, product.Company, productid)
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
