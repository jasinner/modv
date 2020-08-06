package branches

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

var reader io.Reader

//TODO move to utils
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

//TODO move to utils
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

func TestParse(t *testing.T) {
	reader, err := os.Open("testdata/testmod.txt")
	if err != nil {
		log.Fatal(err)
	}

	moduleGraph := NewModuleGraph(reader)
	moduleGraph.Parse()

	expectedBranches := make(map[Module][]Module, 0)

	err = Load("testdata/expected_branches.gob", &expectedBranches)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expectedBranches, moduleGraph.branches) {
		t.Logf("expected\n%v", expectedBranches)
		t.FailNow()
	}
}
