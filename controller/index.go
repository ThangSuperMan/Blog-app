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

// 7200 miliseconds
var LIVETIME_COOKIE int = 7200

func CheckSessionCookieExists(w http.ResponseWriter, r *http.Request) bool {
  fmt.Println("AutoAdaptedAccountAuthenticated")
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

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)

	var isAuthenticated bool = CheckSessionCookieExists(w, r)

  if isAuthenticated {
    fmt.Println("isAuthenticated: ", isAuthenticated)
    cookieName := "my_cookie"
    var sessionToken string = GetSessionCookie(cookieName, r)
    fmt.Println("sessionToken: ", sessionToken )
  }

  tpl.ExecuteTemplate(w, "index.html", nil)
  

	// if isAuthenticate {
	// 	cookie, _ := r.Cookie("my_cookie")
	// 	sessionsToken := cookie.Value
 //    // db := model.ConnectDatabase()
 //    // model.GetInfoUser(db)
	// 	accessToken := AccessToken{
	// 		IsSignedIn: true,
	// 		Username:   sessions[sessionsToken].username,
	// 	}
 //
	// 	tpl.ExecuteTemplate(w, "index.html", accessToken)
	// } else {
	// 	tpl.ExecuteTemplate(w, "index.html", nil)
	// }
}

func RenderSignUpPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RenderSignUpPage")
	tpl, err := template.ParseGlob("./templates/*.html")

	helper.HaltOn(err)
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

func HandlerSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerSignup")

	if r.Method == http.MethodGet {
		tpl, e := template.ParseGlob("./templates/*.html")

		helper.HaltOn(e)
		tpl.ExecuteTemplate(w, "signup.html", nil)
		return
	}

	r.ParseForm()
	r.ParseMultipartForm(10)

	username := r.FormValue("username")
	password := r.FormValue("password")
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
		osFile, err = ioutil.TempFile("uploads/images", "*.jpg")
		avatarName = strings.TrimLeft(osFile.Name(), "uploads/images/")
	} else if contentType == "image/png" {
		osFile, err = ioutil.TempFile("uploads/images", "*.png")
		avatarName = strings.TrimLeft(osFile.Name(), "uploads/images/")
	}

	defer osFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	osFile.Write(fileBytes)

	if password == confirmPassword {
		db := model.ConnectDatabase()
		model.AddUser(db, username, password, avatarName, createdAt, updatedAt)
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

	e := r.ParseForm()
	if e != nil {
		log.Fatal(e)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
  db := model.ConnectDatabase()
  usernameModel, passwordModel := model.GetInfoUser(db, username, password)

  if usernameModel != "" && passwordModel != ""{
    if username == usernameModel && password == passwordModel {
      fmt.Println("username == usernameModel && password == passwordModel")
      var idUser int = model.GetIdUser(db, username)  
      fmt.Println("idUser: ", idUser)

      var sessionToken string = uuid.NewString()

      expriesAt := time.Now().Add(time.Duration(LIVETIME_COOKIE) * time.Second)
      model.AddSession(db, sessionToken, idUser)

      http.SetCookie(w, &http.Cookie {
        Name: "my_cookie",
        Value: sessionToken,
        Expires: expriesAt,
      })

      fmt.Fprintln(w, "<p>Successfully signup, now have fun and enjoy.</p>")
    } else {
      fmt.Println("The info user you just type had somthing wrong!")
    }
  }


}
