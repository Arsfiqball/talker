package exco_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arsfiqball/talker/exco"
)

func TestHttpGetCheck(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fakeServer := httptest.NewServer(handler)
	defer fakeServer.Close()

	t.Run("success", func(t *testing.T) {
		err := exco.HttpGetCheck(fakeServer.URL, 100*time.Millisecond)(context.Background())
		if err != nil {
			t.Error("expected no error, got", err)
		}
	})
}
