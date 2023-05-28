package DBMS

import (
	"GoSoft/graph/model"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"
)

func SHA1(password string) string {
	hash := sha1.Sum([]byte((password)))
	return hex.EncodeToString(hash[:])
}

func GenerateToken(password string) string {
	return SHA1(password + time.Now().String())
}

func ValidateToken(token string) bool {
	CheckConnection()
	rows, err := PostgreSQL.Query(`SELECT tokendate FROM users WHERE token = $1`, token)
	if err != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	var date time.Time
	err = rows.Scan(&date)
	if err != nil {
		return false
	}
	err = rows.Close()
	if err != nil {
		return false
	}
	if time.Since(date) > (time.Hour * 72) {
		return false
	}
	return true
}

func ValidatePrivileges(token string, role string) bool {
	CheckConnection()
	rows, err := PostgreSQL.Query(`SELECT tokendate FROM users WHERE token = $1 AND role = $2`, token, role)
	if err != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	var date time.Time
	err = rows.Scan(&date)
	if err != nil {
		return false
	}
	err = rows.Close()
	if err != nil {
		return false
	}
	if time.Since(date) > (time.Hour * 72) {
		return false
	}
	return true
}

func RegisterCustomer(email string, name string, surname string, gender string, password string) (string, error) {
	CheckConnection()
	encpassword := SHA1(password)
	token := GenerateToken(password)
	result, err := PostgreSQL.Exec(`
INSERT INTO users(email, name, surname, password, gender, token, tokendate, registrationdate, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`, email, name, surname, encpassword, gender, token, time.Now(), time.Now(), "customer")
	if err != nil {
		return "", err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("can't register this user")
	}
	return token, nil
}

func LoginCustomer(email string, password string) (string, error) {
	CheckConnection()

	rows, err := PostgreSQL.Query(`SELECT * FROM users WHERE email = $1 AND password = $2`, email, SHA1(password))
	if err != nil {
		return "", nil
	}
	if !rows.Next() {
		return "", nil
	}
	err = rows.Close()
	if err != nil {
		return "", nil
	}

	token := GenerateToken(password)
	result, err := PostgreSQL.Exec(`UPDATE users SET token = $1, tokendate =  $2 WHERE email = $3 AND password = $4`, token, time.Now(), email, SHA1(password))
	if err != nil {
		return "", err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("user not found")
	}
	return token, nil
}

func CheckUsersProduct(email string, password string, productid int) (bool, error) {
	CheckConnection()
	rows, err := PostgreSQL.Query(`
SELECT userid, datetime, productid, price, orderid, paid, count, subscriptiontype
	FROM public.purchase
	WHERE userid = (SELECT userid FROM users
					WHERE email = $1 AND password = $2)
					AND productid = $3;`, email, SHA1(password), productid)
	if err != nil {
		return false, nil
	}
	if !rows.Next() {
		return false, nil
	}
	err = rows.Close()
	if err != nil {
		return false, nil
	}
	return true, nil
}

func GetProfile(token string) (*model.User, error) {
	CheckConnection()
	rows, err := PostgreSQL.Query(`SELECT userid, email, name, surname, gender, registrationdate, role FROM users WHERE token = $1`, token)
	if err != nil {
		return nil, err
	}

	var user *model.User
	if rows.Next() {
		user = new(model.User)
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Surname, &user.Gender, &user.Date, &user.Role)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("user not found")
	}
	err = rows.Close()
	if err != nil {
		return nil, nil
	}
	return user, nil
}

func UpdateProfile(token string, email string, name string, surname string, gender string, password string) error {
	CheckConnection()
	result, err := PostgreSQL.Exec(`
UPDATE users 
SET 
  email = CASE 
             WHEN length($1) > 0 THEN $2
             ELSE email
           END,
  name = CASE 
             WHEN length($2) > 0 THEN $2
             ELSE name
           END,
  surname = CASE 
             WHEN length($3) > 0 THEN $3
             ELSE surname
           END,
  gender = CASE 
             WHEN length($4) > 0 THEN $4
             ELSE gender
           END,
  password = CASE 
             WHEN length($5) > 0 THEN $5
             ELSE password
           END
  WHERE token = $6;
`, email, name, surname, gender, password, token)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not found")
	}
	return nil
}
