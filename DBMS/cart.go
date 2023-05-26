package DBMS

import (
	"GoSoft/graph/model"
	"errors"
)

func CartGet(token string) ([]*model.CartItem, error) {
	checkConnection()
	rows, err := PostgreSQL.Query(`
SELECT cart.productid, name, description, photo, file, price, count
FROM public.cart INNER JOIN store ON store.productid = cart.productid
	WHERE userid = (SELECT userid FROM users
		WHERE token = $1);
`, token)
	if err != nil {
		return nil, err
	}

	var cart []*model.CartItem
	for rows.Next() {
		p := new(model.CartItem)
		p.Product = new(model.Product)
		err := rows.Scan(&p.Product.ID, &p.Product.Name, &p.Product.Description, &p.Product.Photo, &p.Product.File, &p.Product.Price, &p.Count)
		if err != nil {
			return nil, err
		}
		cart = append(cart, p)
	}
	return cart, nil
}

func CartAdd(token string, productid int, count int) error {
	checkConnection()
	_, err := PostgreSQL.Query(`
INSERT INTO cart (userid, productid, count)
SELECT userid, $1, $2
FROM users 
WHERE token = $3
ON CONFLICT (userid, productid) DO UPDATE SET count = cart.count + $2;
`, productid, count, token)
	if err != nil {
		return err
	}
	return nil
}

func CartRemove(token string, productid int, count int) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
WITH updated AS (
  UPDATE cart
  SET count = count - $3
  WHERE productid = $2
    AND userid = (SELECT userid FROM users WHERE token = $1)
  RETURNING *
)
DELETE FROM cart
WHERE productid = $2 AND count <= 0;`, token, productid, count)
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("can't remove product")
	}

	if err != nil {
		return err
	}
	return nil
}

func CartGetItem(token string, productid int) (*model.CartItem, error) {
	checkConnection()
	rows, err := PostgreSQL.Query(`
SELECT cart.productid, name, description, photo, file, price, count
FROM public.cart INNER JOIN store ON store.productid = cart.productid
	WHERE userid = (SELECT userid FROM users
		WHERE token = $1) AND productid = $2;
`, token, productid)
	if err != nil {
		return nil, err
	}

	var p *model.CartItem
	for rows.Next() {
		p = new(model.CartItem)
		err := rows.Scan(&p.Product.ID, &p.Product.Name, &p.Product.Description, &p.Product.Photo, &p.Product.File, &p.Product.Price, &p.Count)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
