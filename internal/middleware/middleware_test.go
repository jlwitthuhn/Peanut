// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"peanut/internal/keynames/contextkeys"
	"testing"
	"time"
)

func TestWrapHandler_AppliesMiddlewareInOrder(t *testing.T) {
	var order []string

	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	makeMiddleware := func(name string) MiddlewareFunc {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+"-before")
				next.ServeHTTP(w, r)
				order = append(order, name+"-after")
			})
		}
	}

	wrapped := WrapHandler(base, makeMiddleware("first"), makeMiddleware("second"), makeMiddleware("third"))
	wrapped.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	expected := []string{"first-before", "second-before", "third-before", "handler", "third-after", "second-after", "first-after"}
	if len(order) != len(expected) {
		t.Fatalf("got %v, want %v", order, expected)
	}
	for i := range expected {
		if order[i] != expected[i] {
			t.Fatalf("position %d: got %q, want %q\nfull order: %v", i, order[i], expected[i], order)
		}
	}
}

func TestWrapHandler_NoMiddleware(t *testing.T) {
	called := false
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	wrapped := WrapHandler(base)
	wrapped.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	if !called {
		t.Fatal("base handler was not called")
	}
}

func TestSecurityHeaders(t *testing.T) {
	var innerHeaders http.Header
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		innerHeaders = w.Header().Clone()
	})

	handler := SecurityHeaders()(next)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	if got := innerHeaders.Get("Content-Security-Policy"); got != "frame-ancestors 'none';" {
		t.Errorf("Content-Security-Policy = %q, want %q", got, "frame-ancestors 'none';")
	}
	if got := innerHeaders.Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("X-Frame-Options = %q, want %q", got, "DENY")
	}
}

func TestRequestId_SetsContextValue(t *testing.T) {
	var gotId string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotId, _ = r.Context().Value(contextkeys.RequestId).(string)
	})

	handler := RequestId()(next)
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	if gotId == "" {
		t.Fatal("request ID was not set in context")
	}
	if len(gotId) != 8 {
		t.Errorf("request ID length = %d, want 8 hex chars", len(gotId))
	}
	if _, err := hex.DecodeString(gotId); err != nil {
		t.Errorf("request ID %q is not valid hex: %v", gotId, err)
	}
}

func TestRequestId_UniquePerRequest(t *testing.T) {
	ids := make(map[string]bool)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(contextkeys.RequestId).(string)
		ids[id] = true
	})

	handler := RequestId()(next)
	for i := 0; i < 100; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	}

	if len(ids) != 100 {
		t.Errorf("got %d unique IDs out of 100 requests", len(ids))
	}
}

func TestRequestTimer_SetsContextValue(t *testing.T) {
	var gotTime time.Time
	var ok bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTime, ok = r.Context().Value(contextkeys.RequestTimerBegin).(time.Time)
	})

	before := time.Now()
	handler := RequestTimer()(next)
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	after := time.Now()

	if !ok {
		t.Fatal("request timer was not set in context")
	}
	if gotTime.Before(before) || gotTime.After(after) {
		t.Errorf("timer %v not between %v and %v", gotTime, before, after)
	}
}

func TestCheckPermissions_BlocksWhenPermMissing(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called when permission is missing")
	})

	perms := map[string]struct{}{
		"Admin/Gui/View": {},
	}
	ctx := context.WithValue(context.Background(), contextkeys.UserPerms, perms)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	handler := CheckPermissions("Admin/Gui/View", "Admin/Forums/Edit")(next)
	rec := httptest.NewRecorder()

	func() {
		defer func() { recover() }()
		handler.ServeHTTP(rec, req)
	}()

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestCheckPermissions_NoRequiredPerms(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	perms := map[string]struct{}{}
	ctx := context.WithValue(context.Background(), contextkeys.UserPerms, perms)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	handler := CheckPermissions()(next)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if !called {
		t.Fatal("next handler was not called when no permissions are required")
	}
}
