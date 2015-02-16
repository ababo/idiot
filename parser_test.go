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
	Matches         json.RawMessage
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
		if err := json.Unmarshal(data, &case_); err != nil {
			t.Error(err)
		}

		matches := Parse(
			case_.Text, case_.Nonterminal, case_.HypothesesLimit)

		json_ := getParseMatchesJson(matches)
		if json_ == strings.Trim(string(case_.Matches), " \t\r\n") {
			t.Logf("parse matches for case '%s' are good\n",
				f.Name())
		} else {
			t.Errorf("parse matches for case '%s' "+
				"defers from expected", f.Name())

			if testing.Verbose() {
				t.Logf("actual matches:\n%s\n", json_)
			}
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
