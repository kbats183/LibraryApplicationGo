package tests

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"testing"
)

func getAndParse(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HTTP status code error: %d %s", res.StatusCode, res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func TestHome(t *testing.T) {
	document, err := getAndParse("http://localhost/home")
	if err != nil {
		t.Error(err)
		return
	}

	pageHome := document.Find("#PageName").Nodes
	if len(pageHome) != 1 {
		t.Error("Expected #PageName in home page")
	}

	if pageHome[0].FirstChild.Data != "Library" {
		t.Errorf("Expected page header Library, but found %v", pageHome[0].FirstChild.Data)
	}
}

func TestBooks(t *testing.T) {
	document, err := getAndParse("http://localhost/books")
	if err != nil {
		t.Error(err)
		return
	}

	pageHome := document.Find("#PageName").Nodes
	if len(pageHome) != 1 {
		t.Error("Expected #PageName in home page")
	}

	if pageHome[0].FirstChild.Data != "Books" {
		t.Errorf("Expected page header Books, but found %v", pageHome[0].FirstChild.Data)
	}

	bookNames := document.Find("#booksList .bookLine .bookLineName").Nodes
	bookAuthors := document.Find("#booksList .bookLine .bookLineAuthor").Nodes

	expectedBooks := [][2]string{
		{"1984", "Джордж Оруэлл"},
		{"Война и мир", "Лев Толстой"},
		{"Лолита", "Владимир Набоков"},
		{"Улисс", "Джеймс Джойс"},
		{"Шум и ярость", "Уильям Фолкнер"}}

	if len(bookNames) != len(expectedBooks) || len(bookAuthors) != len(expectedBooks) {
		t.Errorf("On page Books expected %d, but found %d bookNames and %d bookAuthors", len(expectedBooks), len(bookNames), len(bookAuthors))
		return
	}

	for i, book := range expectedBooks {
		bookName := bookNames[i].FirstChild.Data
		bookAuthor := bookAuthors[i].FirstChild.Data
		if bookName != book[0] || bookAuthor != book[1] {
			t.Errorf("On page Books expected %v(%v), but found %v(%v)", book[0], book[1], bookName, bookAuthor)
		}
	}
}

func main() {
	fmt.Println("hello world")
}
