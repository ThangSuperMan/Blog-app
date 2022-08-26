package model

import (
	"Blog/structs"
	"database/sql"
	"fmt"
)

func ReadAllBlogs(db *sql.DB) []structs.Blog {
	fmt.Println("ReadAllBlogs")
	var id int
	var title string
	var body string
	var createdAt string
  var updatedAt string
  var idComment string
	var idUser int
	fmt.Println("ReadAllBlogs")
	statement := `select * from blogs`
	rows, err := db.Query(statement)

	if err != nil {
		fmt.Println("Error when read all blogs here: ", err)
	}

	blogs := make([]structs.Blog, 0)
	for rows.Next() {
		rows.Scan(&id, &title, &body, &createdAt, &updatedAt, &idUser, &idComment)

		blog := structs.Blog{
			Id_blog:    id,
			Title:      title,
			Body:       body,
			Created_at: createdAt,
			Id_user:    idUser,
		}

		blogs = append(blogs, blog)
	}

	return blogs
}
