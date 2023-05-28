package DBMS

import (
	"GoSoft/graph/model"
	"errors"
	"time"
)

func CartGet(token string) ([]*model.CartItem, error) {
	CheckConnection()
	rows, err := PostgreSQL.Query(`
SELECT cart.productid, name, description, photo, file, price, subscriptiontype, count
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
		err := rows.Scan(&p.Product.ID, &p.Product.Name, &p.Product.Description, &p.Product.Photo, &p.Product.File, &p.Product.Price, &p.Product.Subscriptiontype, &p.Count)
		if err != nil {
			return nil, err
		}
		cart = append(cart, p)
	}
	return cart, nil
}

func CartAdd(token string, productid int, count int) error {
	if count <= 0 {
		return errors.New("count less or equal then zero")
	}
	CheckConnection()
	err := removeOldSoftware()
	if err != nil {
		return err
	}
	rows, err := PostgreSQL.Query(`SELECT * FROM purchase WHERE userid = (SELECT userid FROM users WHERE token = $1) AND productid = $2 AND paid = true;`, token, productid)
	if rows.Next() {
		return errors.New("software already bought")
	}

	_, err = PostgreSQL.Query(`
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
	CheckConnection()
	result, err := PostgreSQL.Exec(`
UPDATE cart SET count = count - $3 WHERE productid = $2 AND userid = (SELECT userid FROM users WHERE token = $1);`, token, productid, count)
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
	CheckConnection()
	rows, err := PostgreSQL.Query(`
SELECT cart.productid, name, description, photo, file, price, subscriptiontype, company, count
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
		err := rows.Scan(&p.Product.ID, &p.Product.Name, &p.Product.Description, &p.Product.Photo, &p.Product.File, &p.Product.Price, &p.Product.Subscriptiontype, &p.Product.Company, &p.Count)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func CartPurchase(token string, orderID string) error {
	CheckConnection()
	result, err := PostgreSQL.Exec(`
INSERT INTO purchase (productid, userid, datetime, subscriptiontype, price, count, paid, orderid)
	(SELECT productid, userid, $2, (SELECT subscriptiontype FROM store WHERE cart.productid = productid),
	        (SELECT price FROM store WHERE cart.productid = productid), count, FALSE, $3
		FROM cart WHERE userid = (SELECT userid FROM users WHERE token = $1));`, token, time.Now(), orderID)
	if err != nil {
		return err
	}
	_, err = PostgreSQL.Exec(`
DELETE FROM cart WHERE userid = (SELECT userid FROM users WHERE token = $1);`, token)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("can't purchase items in the cart")
	}

	if err != nil {
		return err
	}
	return nil
}

func PurchasedSoftware(token string) ([]*model.Purchase, error) {
	CheckConnection()
	err := removeOldSoftware()
	if err != nil {
		return nil, err
	}
	rows, err := PostgreSQL.Query(`
SELECT purchase.productid, name, description, photo, file, purchase.price, count, purchase.subscriptiontype, company, datetime FROM purchase INNER JOIN store
ON purchase.productid = store.productid WHERE userid = (SELECT userid FROM users WHERE token = $1) AND paid = TRUE ORDER BY datetime DESC;
`, token)
	if err != nil {
		return nil, err
	}

	var products []*model.Purchase
	for rows.Next() {
		item := new(model.Purchase)
		item.Product = new(model.Product)
		err := rows.Scan(&item.Product.ID, &item.Product.Name, &item.Product.Description, &item.Product.Photo, &item.Product.File, &item.Product.Price, &item.Count, &item.Product.Subscriptiontype, &item.Product.Company, &item.Date)
		if err != nil {
			return nil, err
		}
		products = append(products, item)
	}
	return products, nil
}

func CartMakePaid(token string, orderID string) error {
	CheckConnection()
	_, err := PostgreSQL.Exec(`
UPDATE purchase
	SET paid=TRUE
	WHERE orderid = $2 AND userid = (SELECT userid FROM users WHERE token = $1);
`, token, orderID)
	if err != nil {
		return err
	}
	return nil
}

func removeOldSoftware() error {
	CheckConnection()
	_, err := PostgreSQL.Query(`
DELETE FROM purchase
WHERE (now() > (datetime + INTERVAL '1 month' * count) AND subscriptiontype = 'month')
OR (now() > (datetime + INTERVAL '1 year' * count) AND subscriptiontype = 'year');`)
	if err != nil {
		return err
	}
	return nil
}

func CheckFilePermission(filepath string, token string) (bool, error) {
	CheckConnection()
	err := removeOldSoftware()
	if err != nil {
		return false, err
	}
	rows, err := PostgreSQL.Query(`
SELECT productid FROM store
	WHERE file = $1 AND productid IN
		(SELECT productid FROM purchase WHERE userid = 
			(SELECT userid FROM users WHERE token = $2));`, filepath, token)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		return false, err
	}
	return true, nil
}
