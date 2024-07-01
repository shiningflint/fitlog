package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAdminUsersIndex(t *testing.T) {
	os.Setenv("APP_ENV", "testing")

	db = GetConnection()

	tx, err := db.Begin()
	if err != nil {
		return
	}

	if _, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "test1", "test1@banana.com"); err != nil {
		return
	}
	if _, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "test2", "test2@banana.com"); err != nil {
		return
	}
	switch err {
	case nil:
		err = tx.Commit()
	default:
		tx.Rollback()
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	w := httptest.NewRecorder()
	usersIndexHandler(w, req)
	res := w.Result()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	items := doc.Find("tbody > tr")
	fmt.Printf("%T\n", items)

	if res.Status != "200 OK" {
		t.Errorf("expected response status to be 200, got %v", res.Status)
	}

	if items.Length() != 3 {
		t.Errorf("expected rows to be 3 including header, got %v", items.Length())
	}

	dx, err := db.Begin()
	if err != nil {
		return
	}

	if _, err = dx.Exec("DELETE FROM users"); err != nil {
		return
	}
	switch err {
	case nil:
		err = dx.Commit()
	default:
		dx.Rollback()
	}
}

func TestMe(t *testing.T) {
	// db = GetConnection()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	meHandler(w, req)
	res := w.Result()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	title := doc.Find("html")

	titleHtml, err := title.Html()
	if err != nil {
		panic(err)
	}
	fmt.Println(titleHtml)
}
