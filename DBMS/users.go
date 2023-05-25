package DBMS

import (
	"crypto/sha1"
	"encoding/hex"
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
	rows, err := PostgreSQL.Query(`SELECT * FROM users WHERE token = $1`, token)
	if err != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	var trash string
	var trashDate time.Time
	var date time.Time
	err = rows.Scan(&trash, &trash, &trash, &trash, &trash, &trash, &date, &trashDate, &trash)
	if err != nil {
		return false
	}
	if time.Since(date) > (time.Hour * 72) {
		return false
	}
	return true
}

func RegisterCustomer(email string, name string, surname string, gender string, password string) (string, error) {
	checkConnection()
	token := GenerateToken(password)
	_, err := PostgreSQL.Query(`INSERT INTO users(email, name, surname, password, gender, token, tokendate, registrationdate, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`, email, name, surname, gender, token, time.Now(), time.Now(), "customer")
	if err != nil {
		return "", err
	}
	return token, nil
}

func LoginCustomer(email string, password string) (string, error) {
	checkConnection()

	rows, err := PostgreSQL.Query(`SELECT * FROM users WHERE email = $1 AND password = $2`, email, password)
	if err != nil {
		return "", nil
	}
	if !rows.Next() {
		return "", nil
	}

	token := GenerateToken(password)
	_, err = PostgreSQL.Query(`UPDATE users SET token = $1 WHERE email = $2 AND password = $3`, token, email, SHA1(password))
	if err != nil {
		return "", err
	}
	return token, nil
}
