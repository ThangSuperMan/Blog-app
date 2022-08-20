package main

import (
	"Blog/controller"
	"Blog/model"
	"fmt"
	"net/http"
	"time"
)

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

	// sessionToken := cookie.Value
	// delete(sessions, sessionToken)
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
	return sessionToken
}

func Render404PageNotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Render404PageNotFound")
	fmt.Fprintln(w, "<p>Render404PageNotFound</p>")
}

func handler() {
	http.HandleFunc("/", controller.RenderHomePage)
	http.HandleFunc("/profile", controller.RenderProfilePage)
	http.HandleFunc("/about", controller.RenderAboutPage)
	http.HandleFunc("/signin", controller.HandleSignIn)
	http.HandleFunc("/signup", controller.HandlerSignup)
	http.HandleFunc("/logout", controller.LogOut)
	http.HandleFunc("/*", Render404PageNotFound)
	// http.NotFound

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	model.InitModel()
	handler()
	http.ListenAndServe(":3002", nil)
}
