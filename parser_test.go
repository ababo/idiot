package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

type parseCase struct {
	Text            string
	Nonterminal     string
	HypothesesLimit uint
	Matches         []ParseMatch
}

func areEqualStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, s := range slice1 {
		if s != slice2[i] {
			return false
		}
	}

	return true
}

func areEqualAttributeSlices(slice1, slice2 []Attribute) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, s := range slice1 {
		if s.Name != slice2[i].Name || s.Value != slice2[i].Value {
			return false
		}
	}

	return true
}

func areEqualParseMatches(matches1, matches2 []ParseMatch) bool {
	if len(matches1) != len(matches2) {
		return false
	}

	removed := make(map[int]bool, len(matches2))
	for _, m1 := range matches1 {
		found := false
		for j, m2 := range matches2 {
			if _, ok := removed[j]; ok {
				continue
			}

			if m1.Text == m2.Text &&
				m1.Nonterminal == m2.Nonterminal &&
				m1.Rule == m2.Rule &&
				areEqualAttributeSlices(
					m1.Attributes, m2.Attributes) &&
				areEqualParseMatches(
					m1.Submatches, m2.Submatches) &&
				areEqualStringSlices(
					m1.Hypotheses, m2.Hypotheses) &&
				m1.HypothesisCount == m2.HypothesisCount {
				removed[j], found = true, true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func TestParser(t *testing.T) {
	dir := path.Join(getRootDir(), "parser_test")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		data, err := ioutil.ReadFile(path.Join(dir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}

		var case_ parseCase
		err = json.Unmarshal(data, &case_)
		if err != nil {
			t.Error(err)
			continue
		}

		matches := Parse(
			case_.Text, case_.Nonterminal, case_.HypothesesLimit)
		if areEqualParseMatches(matches, case_.Matches) {
			t.Logf("parse matches for case '%s' are good\n",
				f.Name())
		} else {
			t.Errorf("parse matches for case '%s' "+
				"defers from expected", f.Name())

			case_.Matches = matches
			actual, _ := json.Marshal(case_)
			t.Logf("case '%s' with actual matches:\n%s\n",
				f.Name(), actual)
		}
	}
}

func TestMain(m *testing.M) {
	dir := getRootDir()

	err := InitMorph(path.Join(dir, "morph.bin"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer FinalizeMorph()

	err = InitRules(path.Join(dir, "rules.yaml"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer FinalizeRules()

	InitCache(256)
	defer FinalizeCache()

	os.Exit(m.Run())
}
