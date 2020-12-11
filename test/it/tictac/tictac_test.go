// +build integration

package tictac

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func Test_Tic(t *testing.T) {
	defer func() {
		if err := Truncate(a.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	req, err := http.NewRequest(http.MethodPost, "/tic", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusNoContent, fmt.Sprintf("status: excpet 204, got %d", w.Code))
}

func Test_Tic_Tac(t *testing.T) {
	defer func() {
		if err := Truncate(a.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, "/tac", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	a.handler.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, fmt.Sprintf("status: excpet 200, got %d", w.Code))
}
