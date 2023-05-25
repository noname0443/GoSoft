package DBMS

import (
	"GoSoft/graph/model"
)

func GetComments(productid int, from int, count int) ([]*model.Comment, error){
	checkConnection()
	rows, err := PostgreSQL.Query(`SELECT * FROM comment WHERE productid = $1 ORDER BY date LIMIT $2 OFFSET $3`, productid, count, from)
	if err != nil{
		return nil, err
	}
	var comments []*model.Comment
	for rows.Next(){
		c := new(model.Comment)
		err := rows.Scan(&c.ID, &c.Date, &c.Userid, &c.Productid, &c.Content)
		if err != nil{
			return nil, err
		}
		comments = append(comments, c)
	}
	if err != nil{
		return nil, err
	}
	return comments, nil
}