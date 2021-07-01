package main

import (
	"./model"
	"./templater"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

var pageAlerts = map[string]*templater.PageAlert{
	"Login.Incorrect":          {Status: "danger", Content: "Incorrect login or password!"},
	"Logout.NoLogin":           {Status: "danger", Content: "You haven't authorization on site!"},
	"Logout.Successful":        {Status: "success", Content: "Successful. You have been logout!"},
	"Reg.IncorrectData":        {Status: "danger", Content: "Login, name, password mast be not empty string"},
	"Reg.LoginAlreadyExist":    {Status: "danger", Content: "Unfortunately, user with this login already exist. Choose another login please!"},
	"Reg.Problem":              {Status: "danger", Content: "Unfortunately, we can't create this user!"},
	"Reg.Successful":           {Status: "success", Content: "Successful. You have been create new user!"},
	"Reg.AlreadyHaveLogin":     {Status: "danger", Content: "You already have account. Please do logout if you want create new account."},
	"GetBook.InvalidBookId":    {Status: "danger", Content: "Invalid book identifier!"},
	"GetBook.InvalidBook":      {Status: "danger", Content: "Unfortunately, you can't get this book!"},
	"GetBook.Successful":       {Status: "success", Content: "Successful. Now you have this book!"},
	"SubmitBook.InvalidBookId": {Status: "danger", Content: "Invalid book identifier!"},
	"SubmitBook.InvalidBook":   {Status: "danger", Content: "Unfortunately, you can't submit this book!"},
	"SubmitBook.Successful":    {Status: "success", Content: "Successful. You returned this book!"},
	"Book.InvalidBook":         {Status: "danger", Content: "Invalid book identifier!"},
}

func getAlert(r *http.Request) *templater.PageAlert {
	msg := r.FormValue("msg")
	if alert, inMap := pageAlerts[msg]; inMap {
		return alert
	}
	err := r.FormValue("err")
	if alert, inMap := pageAlerts[err]; inMap {
		return alert
	}
	return nil
}

func checkAuthentication(r *http.Request) *model.User {
	userIdC, errUserId := r.Cookie("UserId")
	userPasswordC, errUserPassword := r.Cookie("UserPassword")
	if errUserId != nil || errUserPassword != nil {
		return nil
	}

	userId, errUerIdParse := strconv.ParseInt(userIdC.Value, 10, 64)
	if errUerIdParse != nil {
		return nil
	}
	userPassword := userPasswordC.Value

	user := model.GetUserByIdAndPassword(userId, userPassword)
	return user
}

func checkAuthenticationAndRedirect(w http.ResponseWriter, r *http.Request) *model.User {
	user := checkAuthentication(r)

	if user == nil {
		http.Redirect(w, r, "./login", 302)
	}
	return user
}

func pageHome(w http.ResponseWriter, r *http.Request) {
	templater.ExecuteTemplate("home", "Home", w, nil, getAlert(r))
}

func pageLogin(w http.ResponseWriter, r *http.Request) {
	type UserLoginForm struct {
		UserLogin string
	}

	switch r.Method {
	case "GET":
		templater.ExecuteTemplate("login", "Login", w, UserLoginForm{}, getAlert(r))
	case "POST":
		userLogin := r.FormValue("userLogin")
		passHashBytes := sha1.Sum([]byte(r.FormValue("userPassword")))
		passHashString := hex.EncodeToString(passHashBytes[:])
		user := model.GetUserByLoginAndPassword(userLogin, passHashString)
		if user.Id == 0 {
			templater.ExecuteTemplate("login", "Login", w, UserLoginForm{userLogin},
				pageAlerts["Login.Incorrect"])
		} else {
			setCookie(w, "UserId", strconv.FormatInt(user.Id, 10), time.Hour)
			setCookie(w, "UserPassword", passHashString, time.Hour)
			http.Redirect(w, r, "./profile", 302)
		}
	}
}

func pageLogout(w http.ResponseWriter, r *http.Request) {
	if user := checkAuthentication(r); user == nil {
		http.Redirect(w, r, "./?err=Logout.NoLogin", 302)
	} else {
		setCookie(w, "UserId", "", time.Hour)
		setCookie(w, "UserPassword", "", time.Hour)
		http.Redirect(w, r, "./?msg=Logout.Successful", 302)
	}
}

func pageRegistration(w http.ResponseWriter, r *http.Request) {
	if user := checkAuthentication(r); user != nil {
		http.Redirect(w, r, "./?err=Reg.AlreadyHaveLogin", 302)
	}

	type UserRegForm struct {
		UserLogin string
		UserName  string
	}

	switch r.Method {
	case "GET":
		templater.ExecuteTemplate("registration", "Registration", w, UserRegForm{}, nil)
	case "POST":
		userLogin := r.FormValue("userLogin")
		userPassword := r.FormValue("userPassword")
		passHashBytes := sha1.Sum([]byte(userPassword))
		passHashString := hex.EncodeToString(passHashBytes[:])
		userName := r.FormValue("userName")

		formData := UserRegForm{userLogin, userName}

		if userLogin == "" || userPassword == "" || userName == "" {
			templater.ExecuteTemplate("registration", "Registration", w, formData, pageAlerts["Reg.IncorrectData"])
			return
		}

		if free := model.CheckUserLoginFree(userLogin); !free {
			templater.ExecuteTemplate("registration", "Registration", w, formData, pageAlerts["Reg.LoginAlreadyExist"])
			return
		}

		ok := model.CreateUser(userLogin, userName, passHashString)
		if ok {
			http.Redirect(w, r, "./login?msg=Reg.Successful", 302)
		} else {
			templater.ExecuteTemplate("registration", "Registration", w, formData, pageAlerts["Reg.Problem"])
		}
	}
}

func pageProfile(w http.ResponseWriter, r *http.Request) {

	user := checkAuthenticationAndRedirect(w, r)

	if user == nil {
		return
	}
	userBooks := model.GetUserBooks(user.Id)
	templater.ExecuteTemplate("profile", "Profile", w, struct {
		User      model.User
		UserBooks []model.Book
	}{*user, userBooks}, getAlert(r))
}

func pageBooks(w http.ResponseWriter, r *http.Request) {
	books := model.GetBooks()
	templater.ExecuteTemplate("books", "Books", w, books, getAlert(r))
}

func pageBook(w http.ResponseWriter, r *http.Request) {
	bookIdStr := r.FormValue("bookId")

	if bookId, errBookId := strconv.ParseInt(bookIdStr, 10, 64); errBookId != nil {
		templater.ExecuteTemplate("book", "Book", w, nil, pageAlerts["Book.InvalidBook"])
	} else {
		book := model.GetBook(bookId)
		var userId int64
		if user := checkAuthentication(r); user != nil {
			userId = user.Id
		}
		if book.Id != bookId {
			templater.ExecuteTemplate("book", "Book", w, nil, pageAlerts["Book.InvalidBook"])
		} else {
			templater.ExecuteTemplate("book", "Book", w, struct {
				UserId int64
				Book   *model.Book
			}{userId, book}, nil)
		}
	}
}

func pageGetBookForUser(w http.ResponseWriter, r *http.Request) {
	user := checkAuthenticationAndRedirect(w, r)
	if user == nil {
		return
	}

	bookIdStr := r.FormValue("bookId")
	bookId, errBookId := strconv.ParseInt(bookIdStr, 10, 64)

	if errBookId != nil {
		http.Redirect(w, r, "./books?err=GetBook.InvalidBookId", 302)
	}

	isOk := model.SetBookForUser(user.Id, bookId)
	if !isOk {
		http.Redirect(w, r, "./books?err=GetBook.InvalidBook", 302)
	} else {
		http.Redirect(w, r, "./books?msg=GetBook.Successful", 302)
	}
}

func pageSubmitBookForUser(w http.ResponseWriter, r *http.Request) {
	var user = checkAuthenticationAndRedirect(w, r)
	if user == nil {
		return
	}

	r.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")

	bookIdStr := r.FormValue("bookId")
	bookId, errBookId := strconv.ParseInt(bookIdStr, 10, 64)

	if errBookId != nil {
		http.Redirect(w, r, "./profile?err=SubmitBook.InvalidBookId", 302)
	}

	isOk := model.SetBookBackForUser(user.Id, bookId)
	if !isOk {
		http.Redirect(w, r, "./profile?err=SubmitBook.InvalidBook", 302)
	} else {
		http.Redirect(w, r, "./profile?msg=SubmitBook.Successful", 302)
	}
}

func setCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}
