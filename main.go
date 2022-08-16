package main

import (
	"Blog/controller"
	"Blog/helper"
	"Blog/model"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"text/template"
	"time"
)

// Suppose this is our data :)
var users = map[string]string{
	// username : password
	"thang": "1",
	"ngoc":  "1",
}

type AccessToken struct {
	IsSignedIn bool
	Username   string
}

type session struct {
	username string
	expiry   time.Time
}

var sessions = map[string]session{}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func Signin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call me Signin")
	if r.Method == http.MethodGet {
		fmt.Println(http.MethodGet)
		tpl, err := template.ParseGlob("./templates/*.html")
		if err != nil {
			fmt.Println("Error when render Signin page: ", err)
			return
		}

		tpl.ExecuteTemplate(w, "signin.html", nil)
	} else if r.Method == http.MethodPost {
		fmt.Println(http.MethodPost)
		e := r.ParseForm()
		if e != nil {
			log.Fatal(e)
		}

		creds := Credentials{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		expectedPassword, ok := users[creds.Username]

		if !ok || expectedPassword != creds.Password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionsToken := uuid.NewString()
		fmt.Println("sessionsToken: ", sessionsToken)
		expiresAt := time.Now().Add(120 * time.Second)

		//          Key
		sessions[sessionsToken] = session{
			username: creds.Username,
			expiry:   expiresAt,
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "my_cookie",
			Value:   sessionsToken,
			Expires: expiresAt,
		})

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome page")
	cookie, err := r.Cookie("my_cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionsToken := cookie.Value
	userSession, exists := sessions[sessionsToken]
	fmt.Println("userSession: ", userSession)
	fmt.Println("exist: ", exists)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionsToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "<p>Welcome "+userSession.username+"</p>")
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("my_cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := cookie.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[newSessionToken] = session{
		username: userSession.username,
		expiry:   expiresAt,
	}

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout")
	cookie, err := r.Cookie("my_cookie")
	fmt.Println(cookie)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := cookie.Value
	delete(sessions, sessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:    "my_cookie",
		Value:   "",
		Expires: time.Now(),
	})

	fmt.Fprintln(w, "<p>Logout successfully!</p>")
}

func GetSessionCookie(cookieName string, r *http.Request) string {
  fmt.Println("GetCookie")
  cookie, err := r.Cookie(cookieName) 

  if err != nil {
    fmt.Println("Error when trying to get the cookie: ", err)
  }
  sessionToken := cookie.Value 
  return  sessionToken
}


func About(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseGlob("./templates/*.html")
	helper.HaltOn(err)

	var isAuthenticate bool = CheckSessionCookieExists(w, r)

	if isAuthenticate {
		cookie, _ := r.Cookie("my_cookie")
		sessionsToken := cookie.Value
		accessToken := AccessToken{
			IsSignedIn: true,
			Username:   sessions[sessionsToken].username,
		}

		tpl.ExecuteTemplate(w, "about.html", accessToken)
	} else {
		tpl.ExecuteTemplate(w, "about.html", nil)
	}
}

// func Home(w http.ResponseWriter, r *http.Request) {
// 	tpl, err := template.ParseGlob("./templates/*.html")
// 	helper.HaltOn(err)
//
// 	var isAuthenticated bool = CheckSessionCookieExists(w, r)
//
//   if isAuthenticated {
//     fmt.Println("isAuthenticated: ", isAuthenticated)
//     cookieName := "my_cookie"
//     var sessionToken string = GetSessionCookie(cookieName, r)
//     fmt.Println("sessionToken: ", sessionToken )
//   }
//
//   tpl.ExecuteTemplate(w, "index.html", nil)
//   
//
// 	// if isAuthenticate {
// 	// 	cookie, _ := r.Cookie("my_cookie")
// 	// 	sessionsToken := cookie.Value
//  //    // db := model.ConnectDatabase()
//  //    // model.GetInfoUser(db)
// 	// 	accessToken := AccessToken{
// 	// 		IsSignedIn: true,
// 	// 		Username:   sessions[sessionsToken].username,
// 	// 	}
//  //
// 	// 	tpl.ExecuteTemplate(w, "index.html", accessToken)
// 	// } else {
// 	// 	tpl.ExecuteTemplate(w, "index.html", nil)
// 	// }
// }

func handler() {
	http.HandleFunc("/", controller.RenderHomePage)
	http.HandleFunc("/about", About)
	http.HandleFunc("/signin", controller.HandleSignIn)
	http.HandleFunc("/signup", controller.HandlerSignup)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/logout", Logout)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	fileName := "thang"
	fileName += ".jpg"
	fmt.Println("fileName: ", fileName)
	model.InitModel()
	handler()
	http.ListenAndServe(":3000", nil)
}
