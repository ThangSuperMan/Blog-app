package model

import (
	"Blog/structs"
	"database/sql"
	"fmt"
)

func ReadOneBlog(db *sql.DB, idBlog int) structs.Blog {
	fmt.Println("ReadOneBlog")
	var title string
	var body string
	var imageName string
	var idUser int
	var createdAt string

	statement := `select title, body, image_name, created_at, id_user from blogs where id_blog = ?`
	row := db.QueryRow(statement, idBlog)
	err := row.Scan(&title, &body, &imageName, &createdAt, &idUser)
	if err != nil {
		fmt.Println("Error when scan blog record: ", err)
		return structs.Blog{}
	}

	blog := structs.Blog{
		Title:      title,
		Body:       body,
		Image_name: imageName,
		Created_at: createdAt,
		Id_user:    idUser,
	}

	return blog
}

func ReadTheLastestBlog(db *sql.DB) structs.Blog {
	fmt.Println("ReadTheLastestBlog")
	var idBlog int
	var title string
	var body string
	var imageName string
	var createdAt string
	var updatedAt sql.NullString
	// var idComment sql.NullInt64
	var amountOfLikes int
	var idUser int
	statement := `select id_blog, title, body, image_name, created_at, updated_at, amount_of_likes, id_user  
                from blogs 
                order by id_blog desc
                limit 1;`
	row := db.QueryRow(statement)
	err := row.Scan(&idBlog, &title, &body, &imageName, &createdAt, &updatedAt, &amountOfLikes, &idUser)
	fmt.Println("amoutOfLikes: ", amountOfLikes)
	fmt.Println("id_user: ", idUser)
	if err != nil {
		fmt.Println("Error when scan the lastest blog: ", err)
	}

	blog := structs.Blog{
		Id_blog:    idBlog,
		Title:      title,
		Body:       body,
		Image_name: imageName,
		Created_at: createdAt,
		Id_user:    idUser,
	}

	fmt.Println("blog: ", blog)

	return blog
}

func AddBlog(db *sql.DB, title string, body string, imageName, createdAt string, idUser int) {
	fmt.Println("AddBlog model")
	var amountOfLikes int = 0
	statment := `insert into 
               blogs(title, body, image_name, created_at, amount_of_likes, id_user) 
               values (?, ?, ?, ?, ?, ?)`
	stmt, _ := db.Prepare(statment)
	result, _ := stmt.Exec(title, body, imageName, createdAt, amountOfLikes, idUser)
	fmt.Println("result add blog record: ", result)
}

func readTheLastIdBlogInTableBlog(db *sql.DB) int {
	fmt.Println("readTheLastIdBlogInTableBlog")
	var idBlog int
	statement := `select id_blog 
                from blogs
                order by id_blog desc 
                limit 1`
	row := db.QueryRow(statement)
	err := row.Scan(&idBlog)
	if err != nil {
		fmt.Println("Error when scan the lastest blog: ", err)
	}

	return idBlog
}

// Read all blogs except the first one cause it is the lastest blog
func ReadAllBlogs(db *sql.DB) []structs.Blog {
	fmt.Println("ReadAllBlogs mode")
	lastestIdBlog := readTheLastIdBlogInTableBlog(db)
	var idBlog int
	var title string
	var body string
	var imageName string
	var createdAt string
	var updatedAt sql.NullString
	// var idComment sql.NullInt64
	var amountOfLikes int
	var idUser int
	blogs := make([]structs.Blog, 0)

	if lastestIdBlog == 1 {
		fmt.Println("lastestIdBlog == 1")
	} else {
		lastestIdBlog = lastestIdBlog - 1
		statement := `select * 
    from blogs
    where id_blog between 1 and $1
    order by id_blog desc`

		rows, err := db.Query(statement, lastestIdBlog)

		if err != nil {
			fmt.Println("Error when read all blogs here: ", err)
		}

		for rows.Next() {
			// err := rows.Scan(&idBlog, &title, &body, &imageName, &createdAt, &updatedAt, &idComment, &idUser, &amountOfLikes)
			err := rows.Scan(&idBlog, &title, &body, &imageName, &createdAt, &updatedAt, &idUser, &amountOfLikes)
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
	}

	return blogs
}
