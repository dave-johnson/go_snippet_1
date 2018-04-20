package util

import (
	"fmt"
	"io"
	"os"

	"errors"
)

func cp(f, t string) error {
	if f == "" {
		return errors.New("from file can not be blank")
	}
	if t == "" {
		return errors.New("to file can not be blank")
	}

	in, err := os.Open(f)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(t)
	if err != nil {
		return err
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()

	fmt.Println("copy was successful")
	return nil
}
