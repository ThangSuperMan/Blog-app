package controller

import (
	"Blog/helper"
	"Blog/model"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	IsSignedIn  bool
	Username    string
	ProfileName string
	AvatarName  string
}

// One hour
var LIVETIME_COOKIE int = 1

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

	fmt.Println("Have cookie, now delete that.")
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
}

func RenderTemplate(templateName string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderTemplate")
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)

	var cookieExists bool = CheckSessionCookieExists(w, r)
	cookieName := "my_cookie"

	if cookieExists {
		var sessionToken string = GetSessionTokenCookie(cookieName, r)
		if sessionToken != "" {
			db := model.ConnectDatabase()
			var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
			var user model.User
			user = model.GetInfoUser(db, idUser)
			data := AccessToken{
				IsSignedIn:  true,
				Username:    user.Username,
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

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderHomePage")
	templateName := "index.html"
	RenderTemplate(templateName, w, r)
	// tpl, err := template.ParseGlob("./templates/*.html")
	// helper.HaltOn(err)
	//
	// var cookieExists bool = CheckSessionCookieExists(w, r)
	// cookieName := "my_cookie"
	//
	// if cookieExists {
	// 	var sessionToken string = GetSessionTokenCookie(cookieName, r)
	// 	if sessionToken != "" {
	// 		db := model.ConnectDatabase()
	// 		var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
	// 		var user model.User
	// 		user = model.GetInfoUser(db, idUser)
	// 		data := AccessToken{
	// 			IsSignedIn:  true,
	// 			Username:    user.Username,
	// 			ProfileName: user.Profile_name,
	// 			AvatarName:  user.Avatar_name,
	// 		}
	//
	// 		tpl.ExecuteTemplate(w, "index.html", data)
	// 	} else {
	// 		tpl.ExecuteTemplate(w, "index.html", nil)
	// 	}
	// } else {
	// 	tpl.ExecuteTemplate(w, "index.html", nil)
	// }
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
			// username, avatarName := model.GetUsernameAndAvatarNameOfUsersTable(db, idUser)

			var user model.User
			user = model.GetInfoUser(db, idUser)
			fmt.Println("user: ", user)

			data := AccessToken{
				IsSignedIn:  true,
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
	// tpl, err := template.ParseGlob("./templates/*.html")
	// helper.HaltOn(err)
	// var cookieExists bool = CheckSessionCookieExists(w, r)
	//
	// cookieName := "my_cookie"
	//
	// if cookieExists {
	// 	var sessionToken string = GetSessionTokenCookie(cookieName, r)
	// 	if sessionToken != "" {
	// 		db := model.ConnectDatabase()
	// 		var idUser int = model.GetIdUserFromSessionsTable(db, sessionToken)
	// 		// username, avatarName := model.GetUsernameAndAvatarNameOfUsersTable(db, idUser)
	//
	// 		var user model.User
	// 		user = model.GetInfoUser(db, idUser)
	// 		fmt.Println("user: ", user)
	//
	// 		data := AccessToken{
	// 			IsSignedIn:  true,
	// 			Username:    user.Username,
	// 			ProfileName: user.Profile_name,
	// 			AvatarName:  user.Avatar_name,
	// 		}
	//
	// 		tpl.ExecuteTemplate(w, "about.html", data)
	// 	} else {
	// 		tpl.ExecuteTemplate(w, "about.html", nil)
	// 	}
	// } else {
	// 	tpl.ExecuteTemplate(w, "about.html", nil)
	// }
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

	fmt.Println("username: ", username)
	fmt.Println("password: ", password)
	fmt.Println("avatarProfile: ", profileName)
	fmt.Println("confirmPassword: ", confirmPassword)
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

	fmt.Println("r.Method: ", r.Method)
	e := r.ParseForm()
	if e != nil {
		log.Fatal(e)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	db := model.ConnectDatabase()
	var idUser int = model.GetIdUser(db, username)

	model.GetInfoUser(db, idUser)

	usernameModel, passwordModel := model.GetUsernameAndPasswordOfUser(db, username)
	fmt.Println("usernameModel: ", usernameModel)
	fmt.Println("passwordModel: ", passwordModel)

	if usernameModel != "" && passwordModel != "" {
		fmt.Println("if statement")
		if username == usernameModel && password == passwordModel {
			var idUser int = model.GetIdUser(db, username)
			var sessionToken string = uuid.NewString()

			expiry := time.Now().Add(1 * time.Hour)
			createdAt := time.Now().Format("2006-01-02 15:04:05")
			fmt.Println("username == usernameModel && password == passwordModel")
			fmt.Println("createdAt: ", createdAt)
			fmt.Println("expriesAt: ", expiry)
			model.AddSession(db, sessionToken, expiry, createdAt, idUser)

			http.SetCookie(w, &http.Cookie{
				Name:    "my_cookie",
				Value:   sessionToken,
				Expires: expiry,
			})

			fmt.Fprintln(w, "<p>Successfully signin, now have fun and enjoy.</p>")
		} else if username == usernameModel && password != passwordModel {
			fmt.Fprintln(w, "<p>Username exist but the password was wrong, please make sure that you was typing a correct one.</p>")
		}
	} else {
		fmt.Println("else case")
		fmt.Println("Sorry but the username you just typed is not exist, please make sure you was typing the correct username.")
		fmt.Fprintln(w, "<p>Sorry but the username you just typed is not exist, please make sure you was typing the correct username.</p>")
	}
}
