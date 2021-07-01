package main

import (
	"./model"
	"fmt"
	"net/http"
)

func main() {
	if err := model.InitDB(); err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", pageHome)
	http.HandleFunc("/login", pageLogin)
	http.HandleFunc("/logout", pageLogout)
	http.HandleFunc("/registration", pageRegistration)
	http.HandleFunc("/book", pageBook)
	http.HandleFunc("/books", pageBooks)
	http.HandleFunc("/profile", pageProfile)
	http.HandleFunc("/getBook", pageGetBookForUser)
	http.HandleFunc("/submitBook", pageSubmitBookForUser)

	fmt.Println("Server has been started on port 80")

	serverError := http.ListenAndServe(":80", nil)
	if serverError != nil {
		fmt.Println("Error:", serverError)
	}
}
