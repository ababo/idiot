package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	dir := getRootDir()
	data, err := ioutil.ReadFile(path.Join(dir, "test/cases.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	var cases map[string]string
	err = yaml.Unmarshal(data, &cases)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range cases {
		data, err := ioutil.ReadFile(path.Join(dir, "test", k+".yaml"))
		if err != nil {
			t.Error(err)
		}
		expected := strings.Trim(string(data), " \t\n")

		data, _ = yaml.Marshal(Parse(v, "sentence", 0))
		actual := strings.Trim(string(data), " \t\n")

		if actual != expected {
			if testing.Verbose() {
				fmt.Printf("case %s actual result:\n%s\n",
					k, actual)
			}

			t.Errorf("parse result for case "+
				"%s defers from expected", k)
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
