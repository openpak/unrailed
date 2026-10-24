package main

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeIDs map[string]error

func (f fakeIDs) check(tok string) error {
	if err, ok := f[tok]; ok {
		return err
	}
	return errTokenRejected
}

func get(t *testing.T, s *server, target string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
	return w
}

func TestSubscriptionGate(t *testing.T) {
	s := &server{ids: fakeIDs{
		"good":   nil,
		"banned": errAccountBanned,
		"outage": errors.New("connection refused"),
	}}
	for tok, want := range map[string]string{
		"good":   "OK",
		"banned": noOnlineService,
		"forged": noOnlineService,
		"":       noOnlineService,
		"outage": "OK", // an nx-baas outage must not take online away
	} {
		w := get(t, s, "/nintendoHasSubscription?ticket="+tok)
		if w.Code != 200 || w.Body.String() != want {
			t.Errorf("ticket %q: got %d %q, want 200 %q", tok, w.Code, w.Body.String(), want)
		}
	}
}

func TestKickListIsEmpty(t *testing.T) {
	s := &server{ids: fakeIDs{}}
	w := get(t, s, "/getKick?userid=u&ticket=t&buildNumber=1")
	if w.Code != 200 || w.Body.Len() != 0 {
		t.Fatalf("getKick: %d %q", w.Code, w.Body.String())
	}
	if w := get(t, s, "/notifyKick?id=1&name=x&userid=u&ticket=t&buildNumber=1"); w.Code != 200 {
		t.Fatalf("notifyKick: %d", w.Code)
	}
	if w := get(t, s, "/nope"); w.Code != 404 {
		t.Fatalf("unknown path: %d", w.Code)
	}
}

func TestUploadAccepted(t *testing.T) {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.WriteField("userid", "u")
	mw.WriteField("ticket", "t")
	fw, _ := mw.CreateFormFile("avatar", "a.jpg")
	fw.Write([]byte{0xff, 0xd8, 0xff})
	mw.Close()
	req := httptest.NewRequest(http.MethodPost, "/uploadAvatar", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	(&server{ids: fakeIDs{}}).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("uploadAvatar: %d", w.Code)
	}
}
