package model

import (
	"Blog/structs"
	"database/sql"
	"fmt"
)

func AddBlog(db *sql.DB, title string, body string, imageName, createdAt string, idUser int) {
	fmt.Println("AddBlog model")
	statment := `insert into blogs(title, body, image_name, created_at, id_user) values (?, ?, ?, ?, ?)`
	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(title, body, imageName, createdAt, idUser)
	fmt.Println("result add blog record: ", result)
}

func ReadAllBlogs(db *sql.DB) []structs.Blog {
	fmt.Println("ReadAllBlogs")
	var idBlog int
	var title string
	var body string
	var imageName string
	var createdAt string
	var updatedAt sql.NullString
	var idComment sql.NullInt64
	var idUser int
	statement := `select * from blogs`
	rows, err := db.Query(statement)

	if err != nil {
		fmt.Println("Error when read all blogs here: ", err)
	}

	blogs := make([]structs.Blog, 0)
	for rows.Next() {
		err := rows.Scan(&idBlog, &title, &body, &imageName, &createdAt, &updatedAt, &idComment, &idUser)
		if err != nil {
			fmt.Println("Error when scan data here: ", err)
		}

		blog := structs.Blog{
			Id_blog:    idBlog,
			Title:      title,
			Body:       body,
			Image_name: imageName,
			Created_at: createdAt,
			Id_user:    idUser,
		}

		blogs = append(blogs, blog)
	}

	return blogs
}
