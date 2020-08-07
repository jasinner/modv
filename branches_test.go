package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

var moduleGraph *ModGraph

//TODO fix me so this is an actual setup method
func TestMain(m *testing.M) {
	reader, err := os.Open("testdata/testmod.txt")
	if err != nil {
		log.Fatal(err)
	}

	moduleGraph = NewModuleGraph(reader)
	moduleGraph.Parse()
	m.Run()
}

func TestParse(t *testing.T) {
	//TODO change this back to a series of decending branches
	doTest("testdata/expected_branches.gob", t)
}

func TestFilterRoot(t *testing.T) {
	err := moduleGraph.Filter(newModule("github.com/poloxue/testmod"))
	if err == nil {
		t.Logf("Expected a root module error")
		t.Fail()
	}
}

func TestFilterDirectDep(t *testing.T) {
	moduleGraph.Filter(newModule("rsc.io/sampler@v1.3.1"))
	doTest("testdata/expected_direct.gob", t)
}

func TestFilterNestedDep(t *testing.T) {
	moduleGraph.Filter(newModule("golang/x/fiction@v0.1.1"))
	doTest("testdata/expectedNestedDep.gob", t)
}

func TestFilterShort(t *testing.T) {
	moduleGraph.FilterShort(newModule("golang/x/fiction@v0.1.1"))
	doTest("testdata/expectedNestedShort.gob", t)
}

func TestFilterShortDirect(t *testing.T) {
	moduleGraph.FilterShort(newModule("rsc.io/sampler@v1.3.1"))
	fmt.Println(moduleGraph.branches)
}

func doTest(expected string, t *testing.T) {
	expectedBranches := make(map[Module][]Module, 0)
	err := Load(expected, &expectedBranches)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expectedBranches, moduleGraph.branches) {
		t.Logf("expected\n%v", expectedBranches)
		t.FailNow()
	}
}
