package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

var graphTemplate = `digraph {
{{ $mods := .mods}}
{{- if eq .direction "horizontal" -}}
rankdir=LR;
{{ end -}}
node [shape=box];
{{ range $mod, $modId := .mods -}}
{{ $modId }} [label="{{ $mod.Print }}"];
{{ end -}}

{{- range $dep := .dependencies -}}
{{ $dep.Print $mods }}
{{ end -}}
}
`

type Module struct {
	Name    string
	Version string
}

func NewModule(s string) Module {
	modParts := strings.Split(s, "@")
	if len(modParts) > 1 {
		return Module{modParts[0], modParts[1]}
	}
	return Module{modParts[0], ""}
}

func (m Module) Print() string {
	return fmt.Sprintf("%v\n%v", m.Name, m.Version)
}

type Dependency struct {
	Module Module
	Next   *Module
}

func NewDependency(module Module, next *Module) Dependency {
	return Dependency{module, next}
}

func (d Dependency) Print(mods map[Module]int) string {
	return fmt.Sprintf("%d -> %d", mods[d.Module], mods[*d.Next])
}

type ModuleGraph struct {
	Reader io.Reader

	Mods         map[Module]int
	Dependencies []Dependency
}

func NewModuleGraph(r io.Reader) *ModuleGraph {
	return &ModuleGraph{
		Reader: r,

		Mods:         make(map[Module]int),
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

		modID, ok := m.Mods[module]
		if !ok {
			modID = serialID
			m.Mods[module] = modID
			serialID++
		}

		depModID, ok := m.Mods[depModule]
		if !ok {
			depModID = serialID
			m.Mods[depModule] = depModID
			serialID += 1
		}
		dependency := NewDependency(module, &depModule)
		m.Dependencies = append(m.Dependencies, dependency)
	}
}

func (m *ModuleGraph) Render(w io.Writer, target string) error {
	templ, err := template.New("graph").Parse(graphTemplate)
	//_, err := template.New("graph").Parse(graphTemplate)
	if err != nil {
		return fmt.Errorf("templ.Parse: %v", err)
	}

	//fmt.Printf("target: %v", target)

	//if target != "" {
	//	filteredModules := make(map[string]int)

	//	for mod, index := range m.Mods {
	//		parts := strings.Split(mod, "\n")
	//		if parts[0] == target {
	//			fmt.Printf("target index: %v\n", m.Mods[mod])
	//			filteredModules[mod] = index
	//		}
	//	}

	//fmt.Printf("filtered modules: %v", filteredModules)
	//}

	//Need to loop through again, or add this to the template
	//mod = strings.Replace(mod, "@", "\n", 1)
	//depMod = strings.Replace(depMod, "@", "\n", 1)

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
