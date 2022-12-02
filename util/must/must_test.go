package must

import (
	"strconv"
	"testing"
)

func TestMust(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()

	if Must(strconv.Atoi("99")) != 99 {
		t.Error("99")
	}

	// Should panic
	Must(strconv.Atoi("AAA"))

}
