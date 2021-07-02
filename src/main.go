package main

import (
	"fmt"
	"github.com/kbats183/LibraryApplicationGo/src/model"
	"net/http"
	"os"
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

	port := "80"

	if otherPort, otherPortExists := os.LookupEnv("PORT"); otherPortExists {
		port = otherPort
	}

	fmt.Println("Server has been started on port", port)

	serverError := http.ListenAndServe(":"+port, nil)
	if serverError != nil {
		fmt.Println("Error:", serverError)
	}
}
