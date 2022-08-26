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
	// Number_of_likes int
	Id_comment int
	Id_user    int
}

// Just for render blogs with homepage template
type AccessToken struct {
	IsSignedIn  bool
  Id_user     int
	Username    string
	Password    string
	ProfileName string
	AvatarName  string
	Blogs       []Blog
}
