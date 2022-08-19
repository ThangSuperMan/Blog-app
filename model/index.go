package model

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id_user      int
	Username     string
	Password     string
	Profile_name string
	Avatar_name  string
	Created_at   string
	Updated_at   string
}

func createTables(db *sql.DB) {
	var statement string = `
		create table if not exists "users" (
		"id_user"     integer not null primary key autoincrement,
    "username"	    text not null,
		"password"	    text not null,
		"profile_name"  text not null,
    "avatar_name"   text not null,
		"created_at"    text not null,
    "updated_at"    text
	);

    create table if not exists "sessions" (
    "id_session"    integer not null primary key autoincrement,
    "session_token" text not null,
    "created_at"    text not null,
    "expiry"        datetime not null, 
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
	sqliteDatabase, _ := sql.Open("sqlite3", "./my_database.db")

	defer sqliteDatabase.Close()
	createTables(sqliteDatabase)
}

func ConnectDatabase() *sql.DB {
	sqliteDatabase, _ := sql.Open("sqlite3", "./my_database.db")
	return sqliteDatabase
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

func GetInfoUser(db *sql.DB, idUserParams int) User {
	var idUser int
	var username string
	var password string
	var profileName string
	var avatarName string

	statement := `select id_user, username, password, profile_name, avatar_name from users where id_user=?`
	row := db.QueryRow(statement, idUserParams)
	err := row.Scan(&idUser, &username, &password, &profileName, &avatarName)

	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
		return User{}
	}

	user := User{
		Id_user:      idUser,
		Username:     username,
		Password:     password,
		Profile_name: profileName,
		Avatar_name:  avatarName,
	}

	return user
}

func GetUsernameAndAvatarNameOfUsersTable(db *sql.DB, idUser int) (string, string) {
	var username string
	var avatarName string
	statement := `select username, avatar_name from users where id_user=?`
	row := db.QueryRow(statement, idUser)
	err := row.Scan(&username, &avatarName)

	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
		return "", ""
	}

	return username, avatarName
}

func GetUsernameAndPasswordOfUser(db *sql.DB, usernameParams string) (string, string) {
	var username string
	var password string
	statement := `select username, password from users where username=?`
	row := db.QueryRow(statement, usernameParams)
	err := row.Scan(&username, &password)

	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
		return "", ""
	}

	return username, password
}

func GetIdUserFromSessionsTable(db *sql.DB, sessionToken string) int {
	var idUser int
	statement := `select id_user from sessions where session_token=?`
	row := db.QueryRow(statement, sessionToken)
	err := row.Scan(&idUser)

	if err != nil {
		fmt.Println("Error when scan info user record: ", err)
	}

	return idUser
}

func GetIdUser(db *sql.DB, username string) int {
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

func AddUser(db *sql.DB, username string, password string, profileName string, avatarName string, createdAt string, updatedAt string) {
	statment := `insert into users(username,  password, profile_name, avatar_name, created_at) values (?, ?, ?, ?, ?)`

	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(username, password, profileName, avatarName, createdAt)
	fmt.Println("result add user record: ", result)
}
