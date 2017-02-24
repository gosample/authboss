package authboss

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

func TestAuthBossInit(t *testing.T) {
	t.Parallel()

	ab := New()
	ab.LogWriter = ioutil.Discard
	ab.ViewLoader = mockRenderLoader{}
	err := ab.Init()
	if err != nil {
		t.Error("Unexpected error:", err)
	}
}

func TestAuthBossCurrentUser(t *testing.T) {
	t.Parallel()

	ab := New()
	ab.LogWriter = ioutil.Discard
	ab.StoreLoader = mockStoreLoader{"joe": mockUser{Email: "john@john.com", Password: "lies"}}
	ab.ViewLoader = mockRenderLoader{}
	ab.SessionStoreMaker = newMockClientStoreMaker(mockClientStore{SessionKey: "joe"})
	ab.CookieStoreMaker = newMockClientStoreMaker(mockClientStore{})

	if err := ab.Init(); err != nil {
		t.Error("Unexpected error:", err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "localhost", nil)

	userStruct := ab.CurrentUserP(rec, req)
	us := userStruct.(mockStoredUser)

	if us.Email != "john@john.com" || us.Password != "lies" {
		t.Error("Wrong user found!")
	}
}

func TestAuthBossCurrentUserCallbacks(t *testing.T) {
	t.Parallel()

	ab := New()
	ab.LogWriter = ioutil.Discard
	ab.StoreLoader = mockStoreLoader{"joe": mockUser{Email: "john@john.com", Password: "lies"}}
	ab.ViewLoader = mockRenderLoader{}
	ab.SessionStoreMaker = newMockClientStoreMaker(mockClientStore{SessionKey: "joe"})
	ab.CookieStoreMaker = newMockClientStoreMaker(mockClientStore{})

	if err := ab.Init(); err != nil {
		t.Error("Unexpected error:", err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "localhost", nil)

	afterGetUser := errors.New("afterGetUser")
	beforeGetUser := errors.New("beforeGetUser")
	beforeGetUserSession := errors.New("beforeGetUserSession")

	ab.Callbacks.After(EventGetUser, func(context.Context) error {
		return afterGetUser
	})
	if _, err := ab.CurrentUser(rec, req); err != afterGetUser {
		t.Error("Want:", afterGetUser, "Got:", err)
	}

	ab.Callbacks.Before(EventGetUser, func(context.Context) (Interrupt, error) {
		return InterruptNone, beforeGetUser
	})
	if _, err := ab.CurrentUser(rec, req); err != beforeGetUser {
		t.Error("Want:", beforeGetUser, "Got:", err)
	}

	ab.Callbacks.Before(EventGetUserSession, func(context.Context) (Interrupt, error) {
		return InterruptNone, beforeGetUserSession
	})
	if _, err := ab.CurrentUser(rec, req); err != beforeGetUserSession {
		t.Error("Want:", beforeGetUserSession, "Got:", err)
	}
}

func TestAuthbossUpdatePassword(t *testing.T) {
	t.Skip("TODO(aarondl): Implement")
	/*
		t.Parallel()

		ab := New()
		session := mockClientStore{}
		cookies := mockClientStore{}
		ab.SessionStoreMaker = newMockClientStoreMaker(session)
		ab.CookieStoreMaker = newMockClientStoreMaker(cookies)

		called := false
		ab.Callbacks.After(EventPasswordReset, func(ctx context.Context) error {
			called = true
			return nil
		})

		user1 := struct {
			Password string
		}{}
		user2 := struct {
			Password sql.NullString
		}{}

		r, _ := http.NewRequest("GET", "http://localhost", nil)

		called = false
		err := ab.UpdatePassword(nil, r, "newpassword", &user1, func() error { return nil })
		if err != nil {
			t.Error(err)
		}

		if len(user1.Password) == 0 {
			t.Error("Password not updated")
		}
		if !called {
			t.Error("Callbacks should have been called.")
		}

		called = false
		err = ab.UpdatePassword(nil, r, "newpassword", &user2, func() error { return nil })
		if err != nil {
			t.Error(err)
		}

		if !user2.Password.Valid || len(user2.Password.String) == 0 {
			t.Error("Password not updated")
		}
		if !called {
			t.Error("Callbacks should have been called.")
		}

		called = false
		oldPassword := user1.Password
		err = ab.UpdatePassword(nil, r, "", &user1, func() error { return nil })
		if err != nil {
			t.Error(err)
		}

		if user1.Password != oldPassword {
			t.Error("Password not updated")
		}
		if called {
			t.Error("Callbacks should not have been called")
		}
	*/
}

func TestAuthbossUpdatePasswordFail(t *testing.T) {
	t.Skip("TODO(aarondl): Implement")
	/*
		t.Parallel()

		ab := New()

		user1 := struct {
			Password string
		}{}

		anErr := errors.New("anError")
		err := ab.UpdatePassword(nil, nil, "update", &user1, func() error { return anErr })
		if err != anErr {
			t.Error("Expected an specific error:", err)
		}
	*/
}
