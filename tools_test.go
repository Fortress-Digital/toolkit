package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Errorf("TestTools.RandomString returned wrong length: got %v want %v", len(s), 10)
	}
}
