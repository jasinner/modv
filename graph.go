package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/elliotchance/orderedmap"
)

var graphTemplate = `digraph {
{{ $mods := .mods}}
{{- if eq .direction "horizontal" -}}
rankdir=LR;
{{ end -}}
node [shape=box];
{{ range $mod := $mods.Keys -}}
{{ $mods.GetOrDefault $mod "" }} [label="{{ $mod }}"];
{{ end -}}

{{- range $dep := .dependencies -}}
{{ $dep.Print $mods }}
{{ end -}}
}`

func NewModule(s string) string {
	modParts := strings.Split(s, "@")
	if len(modParts) == 0 {
		errors.New("Delimeter @ not found in module")
	}
	return modParts[0]
}

type Dependency struct {
	Module string
	Next   *string
}

func NewDependency(module string, next *string) Dependency {
	return Dependency{module, next}
}

func (d Dependency) Print(mods orderedmap.OrderedMap) string {
	mod := mods.GetOrDefault(d.Module, "")
	next := mods.GetOrDefault(*d.Next, "")
	return fmt.Sprintf("%d -> %d", mod, next)
}

type ModuleGraph struct {
	Reader io.Reader

	Mods         *orderedmap.OrderedMap
	Dependencies []Dependency
}

func NewModuleGraph(r io.Reader) *ModuleGraph {
	return &ModuleGraph{
		Reader: r,

		Mods:         orderedmap.NewOrderedMap(),
		Dependencies: make([]Dependency, 0),
	}
}

func (m *ModuleGraph) Parse() error {
	bufReader := bufio.NewReader(m.Reader)

	serialID := 1
	for {
		relationBytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		relation := bytes.Split(relationBytes, []byte(" "))
		mod, depMod := strings.TrimSpace(string(relation[0])), strings.TrimSpace(string(relation[1]))

		module := NewModule(mod)
		depModule := NewModule(depMod)

		modID, ok := m.Mods.Get(module)
		if !ok {
			modID = serialID
			m.Mods.Set(module, modID)
			serialID++
		}

		depModID, ok := m.Mods.Get(depModule)
		if !ok {
			depModID = serialID
			m.Mods.Set(depModule, depModID)
			serialID++
		}
		dependency := NewDependency(module, &depModule)
		m.Dependencies = append(m.Dependencies, dependency)
	}
}

func (m *ModuleGraph) Filter(target string) error {
	fmt.Printf("target: %v", target)

	filteredModules := orderedmap.NewOrderedMap()

	_, ok := m.Mods.Get(target)
	if !ok {
		fmt.Println("couldn't find target")
		return nil
	}
	fmt.Println("found target")

	topMod := m.Mods.Front()
	filteredModules.Set(topMod, 1)
	for mod := topMod; mod != nil; mod = mod.Next() {
		if mod.Key == target {
			filteredModules.Set(mod.Key, mod.Value)
		}
	}
	m.Mods = filteredModules
	//	for mod, index := range m.Mods {
	//		parts := strings.Split(mod, "\n")
	//		if parts[0] == target {
	//			fmt.Printf("target index: %v\n", m.Mods[mod])
	//			filteredModules[mod] = index
	//		}
	//	}

	//fmt.Printf("filtered modules: %v", filteredModules)
	//}
	return nil
}

func (m *ModuleGraph) Render(w io.Writer) error {
	templ, err := template.New("graph").Parse(graphTemplate)

	if err != nil {
		return fmt.Errorf("templ.Parse: %v", err)
	}

	var direction string
	if len(m.Dependencies) > 15 {
		direction = "horizontal"
	}

	if err := templ.Execute(w, map[string]interface{}{
		"mods":         m.Mods,
		"dependencies": m.Dependencies,
		"direction":    direction,
		"next":         0,
	}); err != nil {
		return fmt.Errorf("templ.Execute: %v", err)
	}

	return nil
}
