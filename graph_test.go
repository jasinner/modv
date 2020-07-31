package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var Reader io.Reader

func TestMain(m *testing.M) {
	var err error
	Reader, err = os.Open("testdata/testmod.txt")
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func TestEmptyTarget(t *testing.T) {
	DoTest(t, "testdata/expected_no_target.txt", "")
}

func TestQuote(t *testing.T) {
	DoTest(t, "testdata/expected_quote.txt", "rsc.io/quote/v3")
}

func DoTest(t *testing.T, expected string, target string) {
	expectedBytes, err := ioutil.ReadFile(expected)
	if err != nil {
		log.Fatal(err)
	}
	moduleGraph := NewModuleGraph(Reader)
	moduleGraph.Parse()
	var results bytes.Buffer
	w := bufio.NewWriter(&results)
	moduleGraph.Render(w)
	moduleGraph.Filter(target)
	w.Flush()
	if !bytes.Equal(results.Bytes(), expectedBytes) {
		t.Logf("expected: %v", string(expectedBytes))
		t.Errorf("result: %v", string(results.Bytes()))
	}
}
