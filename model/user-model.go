package model

import (
	"Blog/structs"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func GetAllSmallInfoUsers(db *sql.DB) []structs.SmallInfoUser {
	fmt.Println("GetAllSmallInfoUsers")
	var idUser int
	var profileName string
	var avatarName string
	statment := `select id_user, profile_name, avatar_name from users`
	rows, err := db.Query(statment)

	if err != nil {
		fmt.Println("Error when get all user: ", err)
	}

	users := make([]structs.SmallInfoUser, 0)
	for rows.Next() {
		err := rows.Scan(&idUser, &profileName, &avatarName)
		if err != nil {
			fmt.Println("Error when scan data here: ", err)
		}
		user := structs.SmallInfoUser{
			Id_user:      idUser,
			Profile_name: profileName,
			Avatar_name:  avatarName,
		}

		users = append(users, user)
	}

	return users
}

func GetInfoUser(db *sql.DB, idUserParams int) structs.User {
	fmt.Println("GetInfoUser")
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
		return structs.User{}
	}

	user := structs.User{
		Id_user:      idUser,
		Username:     username,
		Password:     password,
		Profile_name: profileName,
		Avatar_name:  avatarName,
	}

	return user
}

func UpdatePassword(db *sql.DB, idUser int, newPassword string) {
	fmt.Println("UpdatePassword")
	statment := `update users set password=? where id_user=?`
	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(newPassword, idUser)
	fmt.Println("result update user's password record: ", result)
}

func UpdateProfileName(db *sql.DB, idUser int, newProfileUserName string) {
	fmt.Println("EditProfileNameOFUser")
	statment := `update users set profile_name=? where id_user=?`
	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(newProfileUserName, idUser)
	fmt.Println("result update  user's profile name record: ", result)
}

func GetUsernameAndAvatarNameOfUsersTable(db *sql.DB, idUser int) (string, string) {
	fmt.Println("GetUsernameAndAvatarNameOfUsersTable")
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
	fmt.Println("GetUsernameAndPasswordOfUser")
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
	fmt.Println("GetIdUserFromSessionsTable")
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
