package main

import (
	"Blog/controller"
	"Blog/model"
	"fmt"
	"net/http"
)

func handler() {
	http.HandleFunc("/", controller.RenderHomePage)
	http.HandleFunc("/profile", controller.RenderProfilePage)
	http.HandleFunc("/edit_profile", controller.HandleEditProfile)
  http.HandleFunc("/blog_detail", controller.RenderBlogDetailPage)
	http.HandleFunc("/add_blog", controller.RenderAddBlogPage)
	http.HandleFunc("/add_single_blog", controller.HandleAddSingleBlog)
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
	port := ":3000"
	fmt.Println("Listenning on the port", port)
	http.ListenAndServe(port, nil)
}
