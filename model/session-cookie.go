package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func GetAccessToken(db *sql.DB, idUser int) string {
	statement := `select session_token from sessions where id_user=?`
	var sessionToken string

	row := db.QueryRow(statement, idUser)

	err := row.Scan(&sessionToken)
	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
    return ""
	}

	return sessionToken
}

func DeleteAllSessionsCookieRelatedToUser(db *sql.DB, idUser int) {
	fmt.Println("DeleteAllSessionsCookieRelatedToUser")
	statement := `delete from sessions where id_user=?`
	stmt, _ := db.Prepare(statement)
	result, _ := stmt.Exec(idUser)
	fmt.Println("Result when delete all sessions cookie related inside the sesions table: ", result)
}

func DeleteSessionCookie(db *sql.DB, sessionToken string) {
	fmt.Println("DeleteSessionCookie: ", sessionToken)
	statement := `delete from sessions where session_token=?`
	stmt, _ := db.Prepare(statement)
	result, _ := stmt.Exec(sessionToken)
	fmt.Println("Result when delete session cookie inside the sesions table: ", result)
}

func ReadSessionCookie(db *sql.DB, sessionTokenParams string) (string, int) {
	statement := `select session_token, id_user from sessions where session_token=?`
	var sessionToken string
	var idUser int

	row := db.QueryRow(statement, sessionTokenParams)

	err := row.Scan(&sessionToken, &idUser)
	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
		return "", -1
	}

	return sessionToken, idUser
}

func AddSession(db *sql.DB, sessionToken string, expriy time.Time, createdAt string, idUser int) {
	statement := `insert into sessions(session_token, expiry, created_at, id_user) values (?, ?, ?, ?)`
	stmt, _ := db.Prepare(statement)
	result, _ := stmt.Exec(sessionToken, expriy, createdAt, idUser)
	fmt.Println("result when add session record: ", result)
}
