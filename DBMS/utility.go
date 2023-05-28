package DBMS

import (
	"GoSoft/Utility"
	"database/sql"
	"fmt"
	"log"
)

var PostgreSQL *sql.DB

func CheckConnection() {
	var err error
	var configs map[string]string

	if PostgreSQL != nil {
		err = PostgreSQL.Ping()
	}
	for err != nil || PostgreSQL == nil {
		log.Println("Reconnect to DB")
		err, configs = Utility.GetConfig("config.ini") // TODO: Get from env
		if err != nil {
			log.Fatal("Can't read config file.")
		}
		DBMSString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", configs["user"], configs["password"], configs["dbname"], configs["sslmode"])
		PostgreSQL, err = sql.Open("postgres", DBMSString)
		if err != nil {
			log.Fatal(err)
		}
	}
}
