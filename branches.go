package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

//Module presents a golang module
type Module struct {
	Name    string
	Version string
	IsRoot  bool
}

//NewModule creates a new golang module
func newModule(bytes []byte) Module {
	modParts := strings.Split(string(bytes), "@")
	parent := strings.TrimSpace(modParts[0])
	if len(modParts) < 2 {
		return Module{parent, "", true}
	}
	dependant := strings.TrimSpace(modParts[1])
	return Module{parent, dependant, false}
}

func (m Module) String() string {
	if m.Version == "" {
		return fmt.Sprintf("%s", m.Name)
	}
	return fmt.Sprintf("%s:%s", m.Name, m.Version)
}

//use a pointer to reduce memory for a large number of modules and relations
type relation struct {
	parent    *Module
	dependant *Module
}

func (r relation) String() string {
	return fmt.Sprintf("%v %v", r.parent, r.dependant)
}

//ModGraph represents a golang module graph as a tree
type ModGraph struct {
	Reader   io.Reader
	branches map[Module][]Module
}

//NewModuleGraph creates a new golang module graph from an io.Reader such as os.Stdin
//The module graph is a prepresentation of golang module graph map of leaf to branch
func NewModuleGraph(r io.Reader) *ModGraph {
	return &ModGraph{
		Reader:   r,
		branches: make(map[Module][]Module, 0),
	}
}

func addModule(uniqModules map[Module]bool, bytes []byte) Module {
	next := newModule(bytes)
	if !uniqModules[next] {
		uniqModules[next] = true
	}
	return next
}

//Parse reads 'go mod graph' output into a ModGraph for filtering
func (m *ModGraph) Parse() error {
	bufReader := bufio.NewReader(m.Reader)
	uniqModules := make(map[Module]bool)
	relations := make([]relation, 0)
	for {
		relationBytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line := bytes.Split(relationBytes, []byte(" "))
		parent := addModule(uniqModules, line[0])
		dependant := addModule(uniqModules, line[1])

		relations = append(relations, relation{&parent, &dependant})
	}
	for _, relation := range relations {
		if relation.parent.IsRoot {
			newBranch := []Module{*relation.parent}
			m.branches[*relation.dependant] = newBranch
		} else {
			branch, ok := m.branches[*relation.parent]
			if ok {
				delete(m.branches, *relation.parent)
				m.branches[*relation.dependant] = append(branch, *relation.parent)
			} else {
				return fmt.Errorf("Didn't find branch with leaf: %v", *relation.parent)
			}
		}
	}
	return nil
}
