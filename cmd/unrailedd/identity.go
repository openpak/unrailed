package main

// Who holds a ticket, asked of nx-baas (/internal/switch/identity) rather than worked out
// here: it also knows whether the account is banned, which a token signature cannot say.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	errTokenRejected = errors.New("token rejected")
	errAccountBanned = errors.New("account banned")
)

func isRefusal(err error) bool {
	return errors.Is(err, errTokenRejected) || errors.Is(err, errAccountBanned)
}

type identityChecker interface {
	check(idToken string) error
}

type identityClient struct {
	url, key string
	http     *http.Client
}

func newIdentityClient() identityClient {
	u := os.Getenv("NX_INTERNAL_URL")
	if u == "" {
		u = "http://127.0.0.1:20070" // nx-baas game/internal API
	}
	return identityClient{url: u, key: os.Getenv("NX_INTERNAL_KEY"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c identityClient) check(tok string) error {
	if tok == "" {
		return errTokenRejected
	}
	body, _ := json.Marshal(map[string]string{"id_token": tok})
	req, err := http.NewRequest(http.MethodPost, c.url+"/internal/switch/identity", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Key", c.key)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errTokenRejected
	case http.StatusForbidden:
		return errAccountBanned
	}
	return fmt.Errorf("nx-baas identity: %s", resp.Status)
}
