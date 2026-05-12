package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

var groupPriority = []string{
	"Соискатель",
	"Работодатель",
	"Общедоступные методы",
	"Приложение",
}

var slugByGroupName = map[string]string{
	"Работодатель":         "employer",
	"Соискатель":           "applicant",
	"Общедоступные методы": "public",
	"Приложение":           "app",
}

type tagGroup struct {
	name string
	tags map[string]struct{}
}

func main() {
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}
	apiDir := filepath.Join(rootDir, "api")
	inPath := filepath.Join(apiDir, "openapi.yml")
	raw, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", inPath, err)
		os.Exit(1)
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		fmt.Fprintf(os.Stderr, "yaml parse: %v\n", err)
		os.Exit(1)
	}
	pruneBareAllOfObjectTail(doc)
	squashAllOfKeepFirstInComponents(doc)
	rawGroups, ok := doc["x-tagGroups"].([]interface{})
	if !ok || len(rawGroups) == 0 {
		fmt.Fprintln(os.Stderr, "missing x-tagGroups")
		os.Exit(1)
	}
	var groups []tagGroup
	for _, g := range rawGroups {
		gm, _ := g.(map[string]interface{})
		if gm == nil {
			continue
		}
		name, _ := gm["name"].(string)
		if name == "" {
			continue
		}
		tagSet := make(map[string]struct{})
		for _, t := range sliceStr(gm["tags"]) {
			tagSet[t] = struct{}{}
		}
		groups = append(groups, tagGroup{name: name, tags: tagSet})
	}
	paths, _ := doc["paths"].(map[string]interface{})
	if paths == nil {
		fmt.Fprintln(os.Stderr, "missing paths")
		os.Exit(1)
	}
	prioIdx := make(map[string]int)
	for i, n := range groupPriority {
		prioIdx[n] = i
	}
	outPaths := map[string]map[string]interface{}{
		"employer":  {},
		"applicant": {},
		"public":    {},
		"app":       {},
	}
	httpMethods := map[string]bool{
		"get": true, "put": true, "post": true, "patch": true,
		"delete": true, "options": true, "head": true, "trace": true,
	}
	for pathKey, pathVal := range paths {
		pathItem, ok := pathVal.(map[string]interface{})
		if !ok {
			continue
		}
		for method, opVal := range pathItem {
			if !httpMethods[method] {
				continue
			}
			op, ok := opVal.(map[string]interface{})
			if !ok {
				continue
			}
			win := pickGroup(sliceStr(op["tags"]), groups, prioIdx)
			if win == "" {
				win = "Общедоступные методы"
			}
			slug := slugByGroupName[win]
			if slug == "" {
				slug = "employer"
			}
			if outPaths[slug][pathKey] == nil {
				outPaths[slug][pathKey] = map[string]interface{}{}
			}
			pi := outPaths[slug][pathKey].(map[string]interface{})
			pi[method] = deepCopy(opVal).(map[string]interface{})
			if p, ok := pathItem["parameters"]; ok {
				pi["parameters"] = deepCopy(p)
			}
			if s, ok := pathItem["servers"]; ok {
				pi["servers"] = deepCopy(s)
			}
		}
	}
	for slug, pm := range outPaths {
		out := cloneRoot(doc)
		delete(out, "x-tagGroups")
		out["paths"] = pm
		data, err := yaml.Marshal(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal %s: %v\n", slug, err)
			os.Exit(1)
		}
		outPath := filepath.Join(apiDir, fmt.Sprintf("openapi.%s.yaml", slug))
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
			os.Exit(1)
		}
		fmt.Println(outPath)
	}
}

func cloneRoot(doc map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(doc)+1)
	for k, v := range doc {
		if k == "paths" {
			continue
		}
		out[k] = deepCopy(v)
	}
	return out
}

func deepCopy(v interface{}) interface{} {
	switch x := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(x))
		for k, v := range x {
			m[k] = deepCopy(v)
		}
		return m
	case []interface{}:
		s := make([]interface{}, len(x))
		for i, v := range x {
			s[i] = deepCopy(v)
		}
		return s
	default:
		return v
	}
}

func sliceStr(v interface{}) []string {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, x := range arr {
		s, _ := x.(string)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func pickGroup(opTags []string, groups []tagGroup, prioIdx map[string]int) string {
	seen := make(map[string]struct{})
	for _, tag := range opTags {
		for _, g := range groups {
			if _, ok := g.tags[tag]; ok {
				seen[g.name] = struct{}{}
			}
		}
	}
	if len(seen) == 0 {
		return ""
	}
	candidates := make([]string, 0, len(seen))
	for n := range seen {
		candidates = append(candidates, n)
	}
	sort.Slice(candidates, func(i, j int) bool {
		a, okA := prioIdx[candidates[i]]
		b, okB := prioIdx[candidates[j]]
		if !okA {
			a = len(groupPriority) + 1
		}
		if !okB {
			b = len(groupPriority) + 1
		}
		if a != b {
			return a < b
		}
		return candidates[i] < candidates[j]
	})
	return candidates[0]
}

func squashAllOfKeepFirstInComponents(doc map[string]interface{}) {
	comp, ok := doc["components"].(map[string]interface{})
	if !ok {
		return
	}
	schemas, ok := comp["schemas"].(map[string]interface{})
	if !ok {
		return
	}
	for _, sch := range schemas {
		squashAllOfKeepFirstWalk(sch)
	}
}

func squashAllOfKeepFirstWalk(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		squashAllOfKeepFirstMap(x)
		for _, child := range x {
			squashAllOfKeepFirstWalk(child)
		}
	case []interface{}:
		for _, item := range x {
			squashAllOfKeepFirstWalk(item)
		}
	}
}

func squashAllOfKeepFirstMap(x map[string]interface{}) {
	ao, ok := x["allOf"].([]interface{})
	if !ok || len(ao) == 0 {
		return
	}
	first, ok := ao[0].(map[string]interface{})
	if !ok {
		return
	}
	extraDesc := x["description"]
	delete(x, "allOf")
	delete(x, "type")
	for k, v := range first {
		x[k] = deepCopy(v)
	}
	if extraDesc != nil {
		if _, has := x["description"]; !has {
			x["description"] = extraDesc
		}
	}
}

func pruneBareAllOfObjectTail(doc map[string]interface{}) {
	comp, ok := doc["components"].(map[string]interface{})
	if !ok {
		return
	}
	schemas, ok := comp["schemas"].(map[string]interface{})
	if !ok {
		return
	}
	for _, sch := range schemas {
		pruneBareAllOfWalk(sch)
	}
}

func pruneBareAllOfWalk(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		if ao, ok := x["allOf"].([]interface{}); ok && len(ao) == 2 {
			first, fok := ao[0].(map[string]interface{})
			second, sok := ao[1].(map[string]interface{})
			if fok && sok && first["$ref"] != nil && bareObjectPlaceholder(second) {
				delete(x, "allOf")
				for k, val := range first {
					x[k] = deepCopy(val)
				}
			}
		}
		for _, child := range x {
			pruneBareAllOfWalk(child)
		}
	case []interface{}:
		for _, item := range x {
			pruneBareAllOfWalk(item)
		}
	}
}

func bareObjectPlaceholder(m map[string]interface{}) bool {
	if len(m) != 1 {
		return false
	}
	if m["type"] != "object" {
		return false
	}
	return true
}
