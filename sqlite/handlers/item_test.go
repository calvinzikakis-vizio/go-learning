package handlers

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetItemHandler(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	t.Run("GetItem_Id=1", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "description"}).
			AddRow(1, "hello there!")

		mock.ExpectQuery("SELECT id, description FROM items WHERE id = ?").WithArgs(1).WillReturnRows(row)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/item?id=1", nil)

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		GetItemHandler(w, r, db)

		assert.Equal(t, http.StatusOK, w.Code)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("GetItem_Id=100", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, description FROM items WHERE id = ?").WithArgs(100).WillReturnError(sql.ErrNoRows)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/item?id=100", nil)

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		GetItemHandler(w, r, db)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "sql: no rows in result set\n", w.Body.String())
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("GetItem_Id=invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/item?id=invalid", nil)

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		GetItemHandler(w, r, db)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "id parameter must be an integer\n", w.Body.String())
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("GetItem_Id=empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/item", nil)

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		GetItemHandler(w, r, db)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "id parameter is required\n", w.Body.String())
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetItemsHandler(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	t.Run("GetItems", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "description"}).
			AddRow(1, "hello there!").
			AddRow(2, "general kenobi")

		mock.ExpectQuery("SELECT id, description FROM items").WillReturnRows(rows)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/items", nil)

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		GetItemsHandler(w, r, db)

		assert.Equal(t, http.StatusOK, w.Code)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("GetItems_NoRows", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, description FROM items").WillReturnError(sql.ErrNoRows)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/items", nil)

		ctx := context.WithValue(r.Context(), "db", db)
		r = r.WithContext(ctx)

		GetItemsHandler(w, r, db)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "sql: no rows in result set\n", w.Body.String())
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestUpdateItemHandler(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	t.Run("UpdateItem_Id=1", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "description"}).
			AddRow(1, "hello there!")
		mock.ExpectQuery("SELECT id, description FROM items WHERE id = ?").WithArgs(1).WillReturnRows(row)
		mock.ExpectExec("UPDATE items SET description = ? WHERE id = ?").WithArgs("Updated description!", 1).WillReturnResult(sqlmock.NewResult(1, 1))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/item?id=1", nil)
		r = r.WithContext(context.WithValue(r.Context(), "db", db))

		UpdateItemHandler(w, r, db)

		assert.Equal(t, http.StatusOK, w.Code)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
