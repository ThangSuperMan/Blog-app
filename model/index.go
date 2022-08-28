package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func createTables(db *sql.DB) {
	// PRAGMA foreign_keys = on;
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
    foreign key (id_user)
      references users (id_user)
        on update cascade
        on delete cascade
    );

    create table if not exists "blogs" (
    "id_blog"         integer not null primary key autoincrement,
    "title"           text not null,
    "body"            text not null,
    "image_name"      text not null,
    "created_at"      text not null,
    "updated_at"      text,
    "id_comment"      integer, 
    "id_user"         integer,
    foreign key (id_user) 
      references users (id_user) 
        on update cascade
        on delete cascade
    );

    create table if not exists "comments" (
    "id_comment" interger 
    "body"       text not null,    
    "created_at" text not null,
    "updated_at" text
    );
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
	sqliteDatabase, _ := sql.Open("sqlite3", "./my_database.db?_foreign_keys=on")
	return sqliteDatabase
}
