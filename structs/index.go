package structs

type User struct {
	Id_user      int
	Username     string
	Password     string
	Profile_name string
	Avatar_name  string
	Updated_at   string
	Created_at   string
}

type Blog struct {
	Id_blog    int
	Title      string
	Body       string
	Updated_at string
	Created_at string
	Image_name string
	Id_comment int
	Id_user    int
}

type SmallInfoUser struct {
	Id_user      int
	Profile_name string
	Avatar_name  string
}

type AuthorOfTheLastestBlog struct {
	Id_user      int
	Profile_name string
	Avatar_name  string
}

type AccessToken struct {
	Is_signed_in bool
	Id_user      int
	Username     string
	Password     string
	Profile_name string
	Avatar_name  string
	// Just for render blogs with homepage template
  Blogs                              []Blog
  Blog_detail                        Blog
	Lastest_blog                       Blog
	Author_of_the_lastest_blog         AuthorOfTheLastestBlog
  Author_of_the_current_blog_detail  User
	Small_info_user_own_blogs          []SmallInfoUser
}
