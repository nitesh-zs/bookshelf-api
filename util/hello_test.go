package util

import "testing"

func TestHello(t *testing.T) {
	if hello() != "hello" {
		t.Error("Testcase for hello() failed")
	}
}
