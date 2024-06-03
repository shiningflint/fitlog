package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminUsersIndex(t *testing.T) {
	db, err := sql.Open("sqlite3", "app_test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close() // Close the connection when done
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	w := httptest.NewRecorder()
	usersIndexHandler(w, req, db)
	res := w.Result()

	if res.Status != "200 OK" {
		t.Errorf("expected response status to be 200, got %v", res.Status)
	}
}
