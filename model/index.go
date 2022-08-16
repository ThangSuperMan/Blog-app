package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	Id          int
	Username    string
	Password    string
	avatar_name string
	created_at  string
	updated_at  string
}

func createTables(db *sql.DB) {
	var statement string = `
		create table if not exists "users" (
		"id_user"     integer not null primary key autoincrement,
		"username"	  text not null,
		"password"	  text not null,
    "avatar_name" text not null,
		"created_at"  text not null,
    "updated_at"  text
	);

    create table if not exists "sessions" (
    "id_session" integer not null primary key autoincrement,
    "session_token" text not null,
    "id_user" integer not null,
    foreign key (id_user) references users (id_user)
    )
    `

	_, err := db.Exec(statement)
	if err != nil {
		fmt.Println("Error of syntax sql here: ", err)
	}
}

func InitModel() {
	fmt.Println("InitModel")
	sqliteDatabase, _ := sql.Open("sqlite3", "./my_database.db")

	defer sqliteDatabase.Close()
	createTables(sqliteDatabase)
}

func ConnectDatabase() *sql.DB {
	sqliteDatabase, _ := sql.Open("sqlite3", "./my_database.db")
	return sqliteDatabase
}

func ReadSessionCookie(db *sql.DB, sessionTokenParams string) (string, int) {
	fmt.Println("model ReadSession")
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

func AddSession(db *sql.DB, sessionToken string, idUser int) {
	fmt.Println("model AddSession")
	statment := `insert into sessions(session_token, id_user) values (?, ?)`

	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(sessionToken, idUser)
	fmt.Println("result when add session record: ", result)
}

func GetInfoUser(db *sql.DB, usernameParams string, passwordParams string) (string, string) {
	fmt.Println("validateUser")
	var username string
	var password string
	statement := `select username, password from users where username=?`

	row := db.QueryRow(statement, usernameParams, passwordParams)

	err := row.Scan(&username, &password)
	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
		return "", ""
	}

	return username, password
}

func GetIdUser(db *sql.DB, username string) int {
	fmt.Println("GetIdUser")
	var idUser int
	statement := `select id_user from users where username=?`

	row := db.QueryRow(statement, username)

	err := row.Scan(&idUser)
	if err != nil {
		fmt.Println("Error when scan info user: ", err)
		return -1
	}

	return idUser
}

func AddUser(db *sql.DB, username string, password string, avatarName string, createdAt string, updatedAt string) {
	fmt.Println("model AddUser")
	statment := `insert into users(username,  password, avatar_name, created_at) values (?, ?, ?, ?)`

	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(username, password, avatarName, createdAt)
	fmt.Println("result", result)
}
