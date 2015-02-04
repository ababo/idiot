package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var nonterminals map[string][]string

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

func FindNonterminalRules(nonterminal string) []string {
	return nonterminals[nonterminal]
}
