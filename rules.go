package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Rule struct {
	Pattern     string   `yaml:"pat"`
	Equivalents []string `yaml:"equiv"`
}

var nonterminals map[string][]Rule

func InitRules(rules_filename string) error {
	FinalizeRules()

	data, err := ioutil.ReadFile(rules_filename)
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(data, &nonterminals)
	if err != nil {
		nonterminals = nil
		return err
	}

	return nil
}

func FinalizeRules() {
	nonterminals = nil
}

func FindNonterminalRules(nonterminal string) []Rule {
	return nonterminals[nonterminal]
}
