package main

import (
	"Blog/controller"
	"Blog/model"
	"net/http"
)

func handler() {
	http.HandleFunc("/", controller.RenderHomePage)
	http.HandleFunc("/profile", controller.RenderProfilePage)
	http.HandleFunc("/about", controller.RenderAboutPage)
	http.HandleFunc("/signin", controller.HandleSignIn)
	http.HandleFunc("/signup", controller.HandlerSignup)
	http.HandleFunc("/logout", controller.LogOut)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	model.InitModel()
	handler()
	http.ListenAndServe(":3000", nil)
}
