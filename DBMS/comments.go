package DBMS

import (
	"GoSoft/graph/model"
	"errors"
	"time"
)

func GetComments(productid int) ([]*model.ExtendedComment, error) {
	checkConnection()
	rows, err := PostgreSQL.Query(`
SELECT commentid, users.name, users.surname, users.role, date, productid, content
	FROM public.comment INNER JOIN users ON users.userid = comment.userid WHERE productid = $1 ORDER BY DATE DESC;
`,
		productid)
	if err != nil {
		return nil, err
	}
	var comments []*model.ExtendedComment
	for rows.Next() {
		c := new(model.ExtendedComment)
		err := rows.Scan(&c.ID, &c.Name, &c.Surname, &c.Role, &c.Date, &c.Productid, &c.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func AddComment(token string, productid int, content string) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
INSERT INTO comment(
	date, userid, productid, content)
	VALUES ($1, (SELECT userid FROM users WHERE token = $2), $3, $4);`, time.Now(), token, productid, content)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("can't add comment")
	}
	return nil
}

func RemoveComment(token string, commentid int) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
DELETE FROM comment
WHERE userid = (SELECT userid FROM users WHERE token = $1)
AND commentid = $2;`, token, commentid)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func UpdateComment(token string, commentid int, content string) error {
	checkConnection()
	result, err := PostgreSQL.Exec(`
UPDATE comment SET content = $1
WHERE userid = (SELECT userid FROM users WHERE token = $2)
AND commentid = $3;`, content, token, commentid)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("comment not found")
	}
	return nil
}
