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

func LogOut(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LogOut controller")
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
	db := model.ConnectDatabase()
	model.DeleteSessionCookie(db, sessionToken)

	cookie := &http.Cookie{
		Name:     "my_cookie",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "Successfully logout, hope you will came back, see ya.")
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
				IsSignedIn:  true,
				Username:    user.Username,
				Password:    user.Password,
				ProfileName: user.Profile_name,
				AvatarName:  user.Avatar_name,
			}

			tpl.ExecuteTemplate(w, templateName, data)
		} else {
			tpl.ExecuteTemplate(w, templateName, nil)
		}
	} else {
		tpl.ExecuteTemplate(w, templateName, nil)
	}
}

func RenderAddBlogPage(w http.ResponseWriter, r *http.Request) {
	var cookieExists bool = CheckSessionCookieExists(w, r)
	cookieName := "my_cookie"

	var sessionTokenCurrentInBrowser string = GetSessionTokenCookie(cookieName, r)
	fmt.Println("value of cookie: ", sessionTokenCurrentInBrowser)
  r.ParseForm()
  idUser := r.FormValue("id_user")
  fmt.Println("idUser: ", idUser)
  db := model.ConnectDatabase()
  idUserInteger, _ := strconv.Atoi(idUser)  
  accessTokenInDatabase := model.GetAccessToken(db, idUserInteger)
  fmt.Println("accessToken: ", accessTokenInDatabase )

	// SessionToken != "" => User still did not sign in
	if r.Method == http.MethodPost && cookieExists && sessionTokenCurrentInBrowser != "" {
    if sessionTokenCurrentInBrowser == accessTokenInDatabase {
      fmt.Println("sessionTokenCurrentInBrowser == accessTokenInDatabase")
      fmt.Println("r.Method == http.MethodPost")
      fmt.Println("RenderAddBlogPage new")
      fmt.Fprintln(w, "<p>RenderAddBlogPage</p>")
    }
	} else {
		fmt.Fprintln(w, "<p>sorry, you would like to use this feature you have to sign in first, thank you so much.</p>")
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
			var user structs.User
			var blogs []structs.Blog = model.ReadAllBlogs(db)

			user = model.GetInfoUser(db, idUser)
			fmt.Println("user :", user)
			data := structs.AccessToken{
				IsSignedIn:  true,
				Username:    user.Username,
				Password:    user.Password,
				ProfileName: user.Profile_name,
				AvatarName:  user.Avatar_name,
				Blogs:       blogs,
			}

			tpl.ExecuteTemplate(w, templateName, data)
		} else {
			tpl.ExecuteTemplate(w, templateName, nil)
		}
	} else {
		tpl.ExecuteTemplate(w, templateName, nil)
	}
}

func RenderProfilePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderProfilePage")
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)
	var cookieExists bool = CheckSessionCookieExists(w, r)

	cookieName := "my_cookie"

	if cookieExists {
		var sessionToken string = GetSessionTokenCookie(cookieName, r)
		if sessionToken != "" {
			db := model.ConnectDatabase()
			var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
			var user structs.User
			user = model.GetInfoUser(db, idUser)

			data := structs.AccessToken{
				IsSignedIn:  true,
        Id_user: idUser,
				Username:    user.Username,
				ProfileName: user.Profile_name,
				AvatarName:  user.Avatar_name,
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
	fmt.Println("RenderAboutPage")
	templateName := "about.html"
	RenderTemplate(templateName, w, r)
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

	file, fileHeader, e := r.FormFile("avatar_profile")
	helper.HaltOn(e)
	defer file.Close()
	contentType := fileHeader.Header["Content-Type"][0]
	var osFile *os.File
	var err error
	var avatarName string

	if contentType == "image/jpeg" {
		osFile, err = ioutil.TempFile("static/uploads/images", "*.jpg")
		avatarName = strings.TrimLeft(osFile.Name(), "static/uploads/images")
	} else if contentType == "image/png" {
		osFile, err = ioutil.TempFile("static/uploads/images", "*.png")
		avatarName = strings.TrimLeft(osFile.Name(), "static/uploads/images/")
	}

	defer osFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	osFile.Write(fileBytes)

	if password == confirmPassword {
		db := model.ConnectDatabase()
		model.AddUser(db, username, password, profileName, avatarName, createdAt, updatedAt)
		fmt.Fprintln(w, "Successfully signup have fun and ejoy.")
		tpl, err := template.ParseGlob("./templates/*.html")
		helper.HaltOn(err)
		tpl.ExecuteTemplate(w, "signup-successfully.html", nil)
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
