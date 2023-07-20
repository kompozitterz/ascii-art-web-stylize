package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*"))

	style := http.FileServer(http.Dir("./templates"))
	http.Handle("/", style)
	http.HandleFunc("/asciiart", posthandler)
	// http.HandleFunc("/down", download)
	fmt.Println("Le serveur démarre sur le  port: 8083 * (http://localhost:8083)")
	// http.ListenAndServe(":8082", nil)
	http.ListenAndServe(":8083", nil)
}
