package pgsender

import "testing"

func Test_enc(t *testing.T) {
	authHeader := Encode("login", "password")
	if authHeader != "Basic bG9naW46cGFzc3dvcmQ=" {
		t.Errorf("Wrong header, got: %s, want: %s.", authHeader, "Basic bG9naW46cGFzc3dvcmQ=")
	}
}
