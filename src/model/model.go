package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
)

type User struct {
	Id       int64
	Login    string
	Name     string
	password string
}

type Book struct {
	Id          int64
	Name        string
	Author      string
	Description string
	HolderId    int64
}

var connStr string

func InitDB() error {
	dbConfig, err := ioutil.ReadFile("db.config")
	if err != nil {
		return err
	}

	connStr = string(dbConfig)
	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	rowsUsers, err := db.Query("select count(*) from users")
	if err != nil {
		return err
	}
	defer rowsUsers.Close()

	rowsBooks, err := db.Query("select count(*) from books")
	if err != nil {
		return err
	}
	defer rowsBooks.Close()

	return nil
}

func GetUserByLoginAndPassword(login string, password string) *User {
	u := User{}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return &u
	}
	defer db.Close()

	rows, err := db.Query("select id, login, name, password from users where login = $1 and password = $2", login, password)
	if err != nil {
		fmt.Println("DB query error:", err)
		return &u
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Login, &u.Name, &u.password)
		if err != nil {
			fmt.Println("DB rows scan error:", err)
			continue
		}
	}
	return &u
}

func CheckUserLoginFree(login string) bool {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return true
	}
	defer db.Close()

	rows, err := db.Query("select count(*) from users where login = $1", login)
	if err != nil {
		fmt.Println("DB query error:", err)
		return true
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			fmt.Println("DB rows scan error:", err)
			continue
		}
	}
	return count == 0
}

func GetUserByIdAndPassword(id int64, password string) *User {
	u := User{}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return &u
	}
	defer db.Close()

	rows, err := db.Query("select id, login, name, password from users where id = $1 and password = $2", id, password)
	if err != nil {
		fmt.Println("DB query error:", err)
		return &u
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Login, &u.Name, &u.password)
		if err != nil {
			fmt.Println("DB rows scan error:", err)
			continue
		}
	}
	return &u
}

func CreateUser(login, name, password string) bool {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return false
	}
	defer db.Close()

	result, err := db.Exec("insert into users(login, name, password)  values ($1, $2, $3)", login, name, password)
	if err != nil {
		fmt.Println("DB update error:", err)
		return false
	}
	rowsAffected, err := result.RowsAffected()
	return err == nil && rowsAffected == 1
}

func GetBooks() []Book {
	var books []Book
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return books
	}
	defer db.Close()

	rows, err := db.Query("select id, name, author, description, holderId from books order by name asc")
	if err != nil {
		fmt.Println("DB query error:", err)
		return books
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Name, &book.Author, &book.Description, &book.HolderId)
		if err != nil {
			fmt.Println("DB rows scan error:", err)
			continue
		}
		books = append(books, book)
	}
	return books
}

func GetBook(bookId int64) *Book {
	var book Book
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return &book
	}
	defer db.Close()

	rows, err := db.Query("select id, name, author, description, holderId from books where id = $1", bookId)
	if err != nil {
		fmt.Println("DB query error:", err)
		return &book
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.Id, &book.Name, &book.Author, &book.Description, &book.HolderId)
		if err != nil {
			fmt.Println("DB rows scan error:", err)
			continue
		}
	}
	return &book
}

func GetUserBooks(userId int64) []Book {
	var books []Book
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return books
	}
	defer db.Close()

	rows, err := db.Query("select id, name, author, description, holderId from books where holderId = $1 order by name asc", userId)
	if err != nil {
		fmt.Println("DB query error:", err)
		return books
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Name, &book.Author, &book.Description, &book.HolderId)
		if err != nil {
			fmt.Println("DB rows scan error:", err)
			continue
		}
		books = append(books, book)
	}
	return books
}

func SetBookForUser(userId, bookId int64) bool {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return false
	}
	defer db.Close()

	result, err := db.Exec("update Books set holderId = $1 where id = $2 and holderId = 0", userId, bookId)
	if err != nil {
		fmt.Println("DB update error:", err)
		return false
	}
	rowsAffected, err := result.RowsAffected()
	return err == nil && rowsAffected == 1
}

func SetBookBackForUser(userId, bookId int64) bool {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB connection error:", err)
		return false
	}
	defer db.Close()

	result, err := db.Exec("update Books set holderId = 0 where id = $2 and holderId = $1", userId, bookId)
	if err != nil {
		fmt.Println("DB update error:", err)
		return false
	}
	rowsAffected, err := result.RowsAffected()
	return err == nil && rowsAffected == 1
}
