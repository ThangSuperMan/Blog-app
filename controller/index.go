package controller

import (
	"Blog/helper"
	"Blog/model"
	"Blog/structs"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// One hour
var LIVETIME_COOKIE int = 1

func HandleSessionsTokenCookieStillNotLogOut(username string) {
	db := model.ConnectDatabase()
	model.DeleteSessionCookie(db, username)
}

func GetSessionTokenCookie(cookieName string, r *http.Request) string {
	cookie, err := r.Cookie(cookieName)

	if err != nil {
		fmt.Println("Error when trying to get the cookie: ", err)
	}
	sessionToken := cookie.Value
	return sessionToken
}

func CheckSessionCookieExists(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("my_cookie")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	// Session cookie exist
	return true
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Lougout handler")
	_, err := r.Cookie("my_cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookieName := "my_cookie"
	var sessionToken string = GetSessionTokenCookie(cookieName, r)
	fmt.Println("sessionToken: ", sessionToken)
	db := model.ConnectDatabase()
	model.DeleteSessionCookie(db, sessionToken)

	cookie := &http.Cookie{
		Name:     "my_cookie",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	fmt.Println("here we go")
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)
	templateName := "logout-successfully.html"
	err = tpl.ExecuteTemplate(w, templateName, nil)
	helper.HaltOn(err)
}

func RenderTemplate(templateName string, w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)

	var cookieExists bool = CheckSessionCookieExists(w, r)
	cookieName := "my_cookie"

	if cookieExists {
		var sessionToken string = GetSessionTokenCookie(cookieName, r)
		if sessionToken != "" {
			db := model.ConnectDatabase()
			defer db.Close()
			var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
			var user structs.User
			user = model.GetInfoUser(db, idUser)
			data := structs.AccessToken{
				Is_signed_in: true,
				Id_user:      user.Id_user,
				Username:     user.Username,
				Password:     user.Password,
				Profile_name: user.Profile_name,
				Avatar_name:  user.Avatar_name,
			}

			tpl.ExecuteTemplate(w, templateName, data)
		} else {
			tpl.ExecuteTemplate(w, templateName, nil)
		}
	} else {
		tpl.ExecuteTemplate(w, templateName, nil)
	}
}

func HandleAddSingleBlog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleAddSingleBlog")
	if r.Method == http.MethodGet {
		fmt.Println("http.MethodGet: ")
		return
	}

	fmt.Println("http.MethodPost")
	r.ParseForm()
	r.ParseMultipartForm(10)
	title := r.FormValue("title")
	body := r.FormValue("content")
	idUser := r.FormValue("id_user")
	idUserTypeInt, _ := strconv.Atoi(idUser)
	nameInputfile := "image_blog"
	locationUpload := "static/uploads/images/blogs/"
	var nameBlogImage string = UploadFile(nameInputfile, locationUpload, r)
	createdAt := time.Now().Format("2006-01-02 15:04:05")

	db := model.ConnectDatabase()
	model.AddBlog(db, title, body, nameBlogImage, createdAt, idUserTypeInt)
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)
	tpl.ExecuteTemplate(w, "add-blog-successfully.html", nil)
}

func RenderAddBlogPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderAddBlogPage")
	if r.Method == http.MethodGet {
		fmt.Println("Get method")
		templateName := "add-blog.html"
		RenderTemplate(templateName, w, r)
		// return
	}

	// Post Method
	// For adding the blog
	var cookieExists bool = CheckSessionCookieExists(w, r)
	if cookieExists {
		fmt.Println("is cookieExists  :>> ", cookieExists)
	}
	fmt.Println("is cookieExists  :>> ", cookieExists)

	cookieName := "my_cookie"
	var sessionTokenCurrentInBrowser string = GetSessionTokenCookie(cookieName, r)
	r.ParseForm()
	idUser := r.FormValue("id_user")
	db := model.ConnectDatabase()
	idUserInteger, _ := strconv.Atoi(idUser)
	accessTokenInDatabase := model.GetAccessToken(db, idUserInteger)
	fmt.Println("accessToken: ", accessTokenInDatabase)

	// SessionToken != "" => User still did not sign in
	if r.Method == http.MethodPost && cookieExists && sessionTokenCurrentInBrowser != "" {
		if sessionTokenCurrentInBrowser == accessTokenInDatabase {
			templateName := "add-blog.html"
			RenderTemplate(templateName, w, r)
		}
	} else {
		fmt.Fprintln(w, "<p>sorry, you would like to use this feature you have to sign in first, thank you so much.</p>")
	}
}

func RenderBlogDetailPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderDetailBlogPage")

	if r.Method == http.MethodPost {
		return
	}

	templateName := "blog-detail.html"
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)

	idBlog := r.URL.Query().Get("id")
	idBlogTypeInt, _ := strconv.Atoi(idBlog)
	db := model.ConnectDatabase()
	var blogDetail structs.Blog
	blogDetail = model.ReadOneBlog(db, idBlogTypeInt)
	idUser := blogDetail.Id_user
	var user structs.User
	user = model.GetInfoUser(db, idUser)

	var cookieExists bool = CheckSessionCookieExists(w, r)
	cookieName := "my_cookie"

	if cookieExists {
		var sessionToken string = GetSessionTokenCookie(cookieName, r)
		if sessionToken != "" {
			var user structs.User
			var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
			user = model.GetInfoUser(db, idUser)
			data := structs.AccessToken{
				Is_signed_in:                      true,
				Username:                          user.Username,
				Password:                          user.Password,
				Profile_name:                      user.Profile_name,
				Avatar_name:                       user.Avatar_name,
				Blog_detail:                       blogDetail,
				Author_of_the_current_blog_detail: user,
			}

			fmt.Println("blogDetail.Created_at: ", blogDetail.Created_at)

			tpl.ExecuteTemplate(w, templateName, data)
		} else {
			data := structs.AccessToken{
				Is_signed_in:                      false,
				Blog_detail:                       blogDetail,
				Author_of_the_current_blog_detail: user,
			}

			fmt.Println("blogDetail.Created_at: ", blogDetail.Created_at)

			tpl.ExecuteTemplate(w, templateName, data)
		}
	} else {
		data := structs.AccessToken{
			Is_signed_in:                      false,
			Blog_detail:                       blogDetail,
			Author_of_the_current_blog_detail: user,
		}

		fmt.Println("blogDetail.Created_at: ", blogDetail.Created_at)
		tpl.ExecuteTemplate(w, templateName, data)
	}
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderHomePage")
	templateName := "index.html"
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)

	var cookieExists bool = CheckSessionCookieExists(w, r)
	cookieName := "my_cookie"

	if cookieExists {
		var sessionToken string = GetSessionTokenCookie(cookieName, r)
		if sessionToken != "" {
			db := model.ConnectDatabase()
			defer db.Close()
			var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
			var blogs []structs.Blog = model.ReadAllBlogs(db)

			fmt.Println("blogs :>> ", blogs)

			var lastestBlog structs.Blog
			lastestBlog = model.ReadTheLastestBlog(db)
			var authorOwnTheLastestBlog structs.AuthorOfTheLastestBlog
			OwnerTheLastestBlog := model.GetInfoUser(db, lastestBlog.Id_user)
			authorOwnTheLastestBlog = structs.AuthorOfTheLastestBlog{
				Id_user:      OwnerTheLastestBlog.Id_user,
				Profile_name: OwnerTheLastestBlog.Profile_name,
				Avatar_name:  OwnerTheLastestBlog.Avatar_name,
			}

			var user structs.User
			user = model.GetInfoUser(db, idUser)
			var smallInfoUsersOwnBlogs []structs.SmallInfoUser = model.GetAllSmallInfoUsers(db)

			fmt.Println()
			data := structs.AccessToken{
				Is_signed_in:               true,
				Username:                   user.Username,
				Password:                   user.Password,
				Profile_name:               user.Profile_name,
				Avatar_name:                user.Avatar_name,
				Blogs:                      blogs,
				Lastest_blog:               lastestBlog,
				Author_of_the_lastest_blog: authorOwnTheLastestBlog,
				Small_info_user_own_blogs:  smallInfoUsersOwnBlogs,
			}

			fmt.Println("data :>> ", data)

			tpl.ExecuteTemplate(w, templateName, data)
		} else {
			db := model.ConnectDatabase()
			var lastestBlog structs.Blog
			lastestBlog = model.ReadTheLastestBlog(db)
			var blogs []structs.Blog = model.ReadAllBlogs(db)
			var smallInfoUsersOwnBlogs []structs.SmallInfoUser = model.GetAllSmallInfoUsers(db)
			var authorOwnTheLastestBlog structs.AuthorOfTheLastestBlog
			OwnerTheLastestBlog := model.GetInfoUser(db, lastestBlog.Id_user)
			fmt.Println("OwnerTheLastestBlog: ", OwnerTheLastestBlog)

			authorOwnTheLastestBlog = structs.AuthorOfTheLastestBlog{
				Id_user:      OwnerTheLastestBlog.Id_user,
				Profile_name: OwnerTheLastestBlog.Profile_name,
				Avatar_name:  OwnerTheLastestBlog.Avatar_name,
			}
			lastestBlog = model.ReadTheLastestBlog(db)
			data := structs.AccessToken{
				Blogs:                      blogs,
				Lastest_blog:               lastestBlog,
				Author_of_the_lastest_blog: authorOwnTheLastestBlog,
				Small_info_user_own_blogs:  smallInfoUsersOwnBlogs,
			}

			tpl.ExecuteTemplate(w, templateName, data)
		}
	} else {
		db := model.ConnectDatabase()
		var lastestBlog structs.Blog
		lastestBlog = model.ReadTheLastestBlog(db)
		var blogs []structs.Blog = model.ReadAllBlogs(db)
		var smallInfoUsersOwnBlogs []structs.SmallInfoUser = model.GetAllSmallInfoUsers(db)
		var authorOwnTheLastestBlog structs.AuthorOfTheLastestBlog
		OwnerTheLastestBlog := model.GetInfoUser(db, lastestBlog.Id_user)
		authorOwnTheLastestBlog = structs.AuthorOfTheLastestBlog{
			Id_user:      OwnerTheLastestBlog.Id_user,
			Profile_name: OwnerTheLastestBlog.Profile_name,
			Avatar_name:  OwnerTheLastestBlog.Avatar_name,
		}
		data := structs.AccessToken{
			Blogs:                      blogs,
			Small_info_user_own_blogs:  smallInfoUsersOwnBlogs,
			Author_of_the_lastest_blog: authorOwnTheLastestBlog,
			Lastest_blog:               lastestBlog,
		}

		tpl.ExecuteTemplate(w, templateName, data)
	}
}

func RenderProfilePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderProfilePage")
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)
	var cookieExists bool = CheckSessionCookieExists(w, r)

	cookieName := "my_cookie"

	if cookieExists {
		fmt.Println("cookieExists: ", cookieExists)
		var sessionToken string = GetSessionTokenCookie(cookieName, r)
		if sessionToken != "" {
			db := model.ConnectDatabase()
			var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
			var user structs.User
			user = model.GetInfoUser(db, idUser)

			data := structs.AccessToken{
				Is_signed_in: true,
				Id_user:      idUser,
				Username:     user.Username,
				Profile_name: user.Profile_name,
				Avatar_name:  user.Avatar_name,
			}

			tpl.ExecuteTemplate(w, "profile.html", data)
		} else {
			tpl.ExecuteTemplate(w, "profile.html", nil)
		}
	} else {
		tpl.ExecuteTemplate(w, "profile.html", nil)
	}
}

func RenderAboutPage(w http.ResponseWriter, r *http.Request) {
	templateName := "about.html"
	RenderTemplate(templateName, w, r)
}

func UploadFile(nameInputFile string, locationUpload string, r *http.Request) string {
	fmt.Println("UploadFile")
	file, fileHeader, e := r.FormFile(nameInputFile)
	helper.HaltOn(e)
	defer file.Close()
	contentType := fileHeader.Header["Content-Type"][0]
	var osFile *os.File
	var err error
	var imageName string

	if contentType == "image/jpeg" {
		osFile, err = ioutil.TempFile(locationUpload, "*.jpg")
		imageName = strings.TrimLeft(osFile.Name(), locationUpload)
	} else if contentType == "image/png" {
		osFile, err = ioutil.TempFile(locationUpload, "*.png")
		imageName = strings.TrimLeft(osFile.Name(), locationUpload)
	}

	defer osFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	osFile.Write(fileBytes)

	return imageName
}

func HandlerSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerSignup")

	if r.Method == http.MethodGet {
		tpl, e := template.ParseGlob("./templates/*.html")

		helper.HaltOn(e)
		tpl.ExecuteTemplate(w, "signup.html", nil)
		return
	}

	// Post method
	r.ParseForm()
	r.ParseMultipartForm(10)
	username := r.FormValue("username")
	password := r.FormValue("password")
	profileName := r.FormValue("profile_name")
	confirmPassword := r.FormValue("confirm_password")
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := ""

	locationUpload := "static/uploads/images/users/"
	avatarName := UploadFile("avatar_profile", locationUpload, r)

	if password == confirmPassword {
		db := model.ConnectDatabase()
		model.AddUser(db, username, password, profileName, avatarName, createdAt, updatedAt)
		tpl, err := template.ParseGlob("./templates/*.html")
		helper.HaltOn(err)
		fmt.Println(tpl.ExecuteTemplate(w, "signup-successfully.html", nil))
	}
}

func HandleEditProfileName(w http.ResponseWriter, sessionToken string, profileNewName string, password string) {
	fmt.Println("EditProfileName")
	db := model.ConnectDatabase()
	var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
	username, _ := model.GetUsernameAndAvatarNameOfUsersTable(db, idUser)
	_, passwordModel := model.GetUsernameAndPasswordOfUser(db, username)

	if password == passwordModel {
		// Allow user edit info
		fmt.Println("Allow user edit info")
		model.UpdateProfileName(db, idUser, profileNewName)

		tpl, err := template.ParseGlob("./templates/*.html")
		helper.HaltOn(err)
		tpl.ExecuteTemplate(w, "edit-profile-successfully.html", nil)
	} else {
		fmt.Println("Dont allow user edit info")
		fmt.Fprintln(w, "<p>Sorry but your password did not correct, please make sure that you typed a correct one.</p>")
	}
}

func HandleEditPassword(w http.ResponseWriter, sessionToken string, currentPassword string, newPassword string, confirmNewPassword string) {
	fmt.Println("HandleEditPassword")
	db := model.ConnectDatabase()
	var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
	username, _ := model.GetUsernameAndAvatarNameOfUsersTable(db, idUser)
	_, passwordModel := model.GetUsernameAndPasswordOfUser(db, username)

	if currentPassword != passwordModel {
		fmt.Println("currentPassword != passwordModel")
		fmt.Fprintln(w, "<p>Sorry but your current password did not correct, please make sure that you typed a correct one.</p>")
	} else if newPassword != confirmNewPassword {
		fmt.Println("newPassword != confirmNewPassword")
		fmt.Fprintln(w, "<p>Sorry but your new password and confirm new password did not match with each other, please make sure that you typed a similar one.</p>")
	} else {
		fmt.Println("every thing was totally correct.")
		model.UpdatePassword(db, idUser, newPassword)
		tpl, err := template.ParseGlob("./templates/*.html")
		helper.HaltOn(err)
		tpl.ExecuteTemplate(w, "edit-profile-successfully.html", nil)
	}
}

func HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleEditProfile")
	if r.Method == http.MethodPost {
		r.ParseForm()
		editProfileName := r.FormValue("edit_profile_name")
		editPassword := r.FormValue("edit_password")
		// Get user's id based on the access token cookie in the client browser
		cookie, err := r.Cookie("my_cookie")

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var sessionToken string = cookie.Value
		if editProfileName != "" {
			fmt.Println("editProfileName  ")
			profileNewName := r.FormValue("profile_new_name")
			password := r.FormValue("password")
			HandleEditProfileName(w, sessionToken, profileNewName, password)
		} else if editPassword != "" {
			fmt.Println("editPassword")
			currentPassword := r.FormValue("current_password")
			newPassword := r.FormValue("new_password")
			confirmNewPassword := r.FormValue("confirm_new_password")
			HandleEditPassword(w, sessionToken, currentPassword, newPassword, confirmNewPassword)
		}
	}
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleSignIn")
	if r.Method == http.MethodGet {
		tpl, err := template.ParseGlob("./templates/*.html")
		helper.HaltOn(err)
		tpl.ExecuteTemplate(w, "signin.html", nil)
		return
	}

	e := r.ParseForm()
	if e != nil {
		log.Fatal(e)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	db := model.ConnectDatabase()
	var idUser int = model.GetIdUser(db, username)
	// Delete old sessions token cookie
	model.DeleteAllSessionsCookieRelatedToUser(db, idUser)
	model.GetInfoUser(db, idUser)
	usernameModel, passwordModel := model.GetUsernameAndPasswordOfUser(db, username)

	if usernameModel != "" && passwordModel != "" {
		fmt.Println("if statement")
		if username == usernameModel && password == passwordModel {
			var idUser int = model.GetIdUser(db, username)
			var sessionToken string = uuid.NewString()

			expiry := time.Now().Add(1 * time.Hour)
			createdAt := time.Now().Format("2006-01-02 15:04:05")
			model.AddSession(db, sessionToken, expiry, createdAt, idUser)
			http.SetCookie(w, &http.Cookie{
				Name:    "my_cookie",
				Value:   sessionToken,
				Expires: expiry,
			})

			tpl, err := template.ParseGlob("./templates/*.html")
			helper.HaltOn(err)
			tpl.ExecuteTemplate(w, "signin-successfully.html", nil)
		} else if username == usernameModel && password != passwordModel {
			fmt.Fprintln(w, "<p>Username exist but the password was wrong, please make sure that you was typing a correct one.</p>")
		}
	} else {
		fmt.Fprintln(w, "<p>Sorry but the username you just typed is not exist, please make sure you was typing the correct username.</p>")
	}
}
