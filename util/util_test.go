package util

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	err := cp("", "")
	if err.Error() != "from file can not be blank" {
		t.Error("did not get correct response", err)
		t.FailNow()
	}

	from := "src/a.txt"
	err = cp(from, "")
	if err.Error() != "to file can not be blank" {
		t.Error("did not get correct response", err)
		t.FailNow()
	}

	to := "dest/b.txt"
	err = cp(from, to)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("did it work")
	t.FailNow()
}
