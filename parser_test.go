package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"sort"
	"strings"
	"testing"
)

type parseCase struct {
	Text            string
	Nonterminal     string
	HypothesesLimit uint
	Matches         []ParseMatch
}

type byContent []ParseMatch

func (ms byContent) Len() int      { return len(ms) }
func (ms byContent) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }
func (ms byContent) Less(i, j int) bool {
	return fmt.Sprintf("%v", ms[i]) < fmt.Sprintf("%v", ms[j])
}

func TestParser(t *testing.T) {
	if err := initParser(); err != nil {
		t.Fatal(err)
	}
	defer finalizeParser()

	dir := path.Join(getRootDir(), "parser_test")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		filename := path.Join(dir, f.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		var case_ parseCase
		if err := json.Unmarshal(data, &case_); err != nil {
			t.Error(err)
		}

		case_.Matches = Parse(
			case_.Text, case_.Nonterminal, case_.HypothesesLimit)
		sort.Sort(byContent(case_.Matches))

		json_, _ := json.MarshalIndent(&case_, "", " ")
		if string(json_) == strings.Trim(string(data), " \t\r\n") {
			t.Logf("case '%s' is passed", f.Name())
		} else {
			t.Errorf("case '%s' is FAILED", f.Name())

			filename2 := path.Join(dir, f.Name()+".actual")
			ioutil.WriteFile(filename2, json_, 0664)

			if testing.Verbose() {
				cmd := exec.Command("diff", filename, filename2)

				var out bytes.Buffer
				cmd.Stdout = &out

				if err := cmd.Run(); err != nil {
					fmt.Printf("case '%s' diff:\n%s",
						f.Name(), out.String())
				}
			}
		}
	}
}
