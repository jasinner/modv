package main

import (
	"bufio"
	"encoding/gob"
	"os"
)

// Marshal is a function that marshals the object into an
// io.Writer.
// By default, it uses the GOB encoder
func Save(v interface{}, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	e := gob.NewEncoder(w)

	err = e.Encode(v)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
func Load(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	d := gob.NewDecoder(r)

	err = d.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

//IsDir checks is a path is a directory
func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	mode := fi.Mode()
	if mode.IsDir() {
		return true, nil
	}
	return false, nil
}
