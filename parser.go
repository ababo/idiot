package main

import (
	"strings"
)

type Attribute struct {
	Name   string
	Values []string
	token  string
	export bool
}

type ParseMatch struct {
	Text            string
	Nonterminal     string
	Rule            string
	Attributes      []Attribute
	Submatches      []ParseMatch
	Hypotheses      []string
	HypothesisCount uint
}

const TerminalSeparators = " "

func parseTerminal(rule string) (string, string) {
	return "", rule
}

func parseNonterminal(rule string) (string, []Attribute, string) {
	return "", nil, rule
}

func (match *ParseMatch) clone() ParseMatch {
	return *match
}

func (match *ParseMatch) checkAttrValue(attr *Attribute, hypo *string) bool {
	return false
}

func (match *ParseMatch) checkAttrToken(attr *Attribute, hypo *string) bool {
	return false
}

func (match *ParseMatch) cleanExport() {

}

func (match *ParseMatch) checkSubmatch(attrs []Attribute,
	submatch ParseMatch, hypotheses_limit uint) bool {

	var hypo string
	for _, a := range attrs {
		if (len(a.Values) > 0 && !submatch.checkAttrValue(&a, &hypo)) ||
			(len(a.token) > 0 && !match.checkAttrToken(&a, &hypo)) {
			if match.HypothesisCount+submatch.HypothesisCount >=
				hypotheses_limit {
				return false
			}
			match.HypothesisCount++
			match.Hypotheses = append(match.Hypotheses, hypo)
		}

		if a.export {
			match.Attributes = append(match.Attributes, a)
		}
	}

	match.Text += submatch.Text
	match.Submatches = append(match.Submatches, submatch)
	return true
}

func parseRulePart(text string, rule string, match ParseMatch,
	hypotheses_limit uint, output chan []ParseMatch) {

	term, rule := parseTerminal(rule)
	if !strings.HasPrefix(text, term) {
		output <- nil
		return
	}
	match.Text += term
	text = text[len(term):]

	if len(rule) == 0 {
		match.cleanExport()
		output <- []ParseMatch{match}
		return
	}

	pred, attrs, rule := parseNonterminal(rule)
	output2 := make(chan []ParseMatch)
	count := 0
	for _, m := range Parse(text, pred, hypotheses_limit) {
		match := match.clone()
		if match.checkSubmatch(attrs, m, hypotheses_limit) {
			go parseRulePart(text[len(match.Text):],
				rule, match, hypotheses_limit, output2)
			count++
		}
	}

	matches := []ParseMatch{}
	for i := 0; i < count; i++ {
		matches = append(matches, <-output2...)
	}

	output <- matches
}

func parseRule(text string, match ParseMatch,
	hypotheses_limit uint) []ParseMatch {

	rule := match.Rule
	if strings.HasPrefix(rule, "@") {
		_, match.Attributes, rule = parseNonterminal(rule[1:])
	}

	output := make(chan []ParseMatch)
	go parseRulePart(text, rule, match, hypotheses_limit, output)
	return <-output
}

func Parse(text string, nonterminal string,
	hypotheses_limit uint) []ParseMatch {

	if len(nonterminal) == 0 {
		if i := strings.IndexAny(text, TerminalSeparators); i != -1 {
			text = text[:i]
		}
		return []ParseMatch{{Attributes: FindTerminalAttrs(text)}}
	}

	matches := []ParseMatch{}
	for _, r := range FindNonterminalRules(nonterminal) {
		match := ParseMatch{Nonterminal: nonterminal, Rule: r}
		matches = append(matches,
			parseRule(text, match, hypotheses_limit)...)
	}
	return matches
}
