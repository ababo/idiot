package main

import (
	"fmt"
	"strings"
)

type Attribute struct {
	Name   string
	Value  string
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

const TerminalSeparators = " ,."

func parseTerminal(rule string) (string, string) {
	var term string
	found, unescape := false, false
	for i := 0; i < len(rule); i++ {
		escaped := i > 0 && rule[i-1] == '\\'
		if escaped {
			unescape = true
		}
		if rule[i] == '{' && !escaped {
			term, rule = rule[:i], rule[i:]
			found = true
			break
		}
	}

	if !found {
		term, rule = rule, ""
	}

	if unescape {
		term = strings.Replace(term, "\\{", "{", -1)
	}

	return term, rule
}

func parseNonterminal(rule string) (string, []Attribute, string) {
	attrs := []Attribute{}

	if len(rule) == 0 || rule[0] != '{' {
		return "", attrs, rule
	}

	var body string
	if i := strings.Index(rule, "}"); i != -1 {
		body, rule = rule[1:i], rule[i+1:]
	} else {
		body, rule = rule, ""
	}

	var nterm string
	for i, a := range strings.Split(body, " ") {
		if len(a) == 0 {
			continue
		}

		var attr Attribute
		split := strings.Split(a, "=")

		if len(split[0]) == 0 {
			continue
		}
		rvalue := len(split) > 1 && len(split[1]) > 0

		if split[0][0] == '!' {
			attr.Name = split[0][1:]
			attr.export = true
		} else if i > 0 || rvalue {
			attr.Name = split[0]
		} else {
			nterm = split[0]
			continue
		}

		if rvalue {
			if split[1][0] == '@' {
				attr.token = split[1][1:]
			} else {
				attr.Value = split[1]
			}
		}

		attrs = append(attrs, attr)
	}

	return nterm, attrs, rule
}

func (match *ParseMatch) clone() ParseMatch {
	nmatch := *match

	nmatch.Attributes = make([]Attribute, len(match.Attributes))
	copy(nmatch.Attributes, match.Attributes)

	nmatch.Submatches = make([]ParseMatch, len(match.Submatches))
	copy(nmatch.Submatches, match.Submatches)

	nmatch.Hypotheses = make([]string, len(match.Hypotheses))
	copy(nmatch.Hypotheses, match.Hypotheses)

	return nmatch
}

func (match *ParseMatch) findAttr(name string) (Attribute, bool) {
	for _, a := range match.Attributes {
		if a.Name == name {
			return a, true
		}
	}
	return Attribute{}, false
}

func (match *ParseMatch) checkAttrValue(attr Attribute, hypo *string) bool {
	attr2, ok := match.findAttr(attr.Name)
	if ok {
		if attr.Value == attr2.Value {
			return true
		}
	}

	*hypo = fmt.Sprintf("%s = %s", attr.Name, attr2.Value)
	return false
}

func (match *ParseMatch) checkAttrToken(attr Attribute, hypo *string) bool {
	for i := 0; i < len(match.Attributes); i++ {
		a := &match.Attributes[i]
		if a.Name == attr.Name && a.token == attr.token {
			if attr.export {
				a.export = true
			}

			if a.Value == attr.Value {
				return true
			}

			*hypo = fmt.Sprintf("%s = %s", attr.Name, a.Value)
			return false
		}
	}

	match.Attributes = append(match.Attributes, attr)
	return true
}

func (match *ParseMatch) cleanExport() {
	attrs := []Attribute{}

	for _, a := range match.Attributes {
		if len(a.token) == 0 || a.export {
			a.token = ""
			attrs = append(attrs, a)
		}
	}

	match.Attributes = attrs
}

func (match *ParseMatch) checkSubmatch(attrs []Attribute,
	submatch ParseMatch, hypotheses_limit uint) bool {

	var hypo string
	for _, a := range attrs {
		attr, _ := submatch.findAttr(a.Name)

		if len(a.token) > 0 {
			a.Value = attr.Value
		}

		if (len(a.Value) > 0 && !submatch.checkAttrValue(a, &hypo)) ||
			(len(a.token) > 0 && !match.checkAttrToken(a, &hypo)) {
			if match.HypothesisCount+submatch.HypothesisCount >=
				hypotheses_limit {
				return false
			}
			match.HypothesisCount++
			hypo = fmt.Sprintf(
				"%d: %s", len(match.Submatches)+1, hypo)
			match.Hypotheses = append(match.Hypotheses, hypo)
		}

		if a.export && len(a.token) == 0 {
			match.Attributes = append(match.Attributes, attr)
		}
	}

	match.Text += submatch.Text
	match.HypothesisCount += submatch.HypothesisCount
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
		match2 := match.clone()
		if match2.checkSubmatch(attrs, m, hypotheses_limit) {
			go parseRulePart(text[len(m.Text):],
				rule, match2, hypotheses_limit, output2)
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

	matches := []ParseMatch{}

	if len(nonterminal) == 0 {
		var pref, sep string
		if i := strings.IndexAny(text, TerminalSeparators); i != -1 {
			pref, sep = text[:i], text[i:i+1]
		} else {
			pref, sep = text, ""
		}

		for _, m := range FindTerminals(pref, sep) {
			if strings.HasPrefix(text, m.Text) {
				matches = append(matches, m)
			}
		}
	} else {
		for _, r := range FindNonterminalRules(nonterminal) {
			match := ParseMatch{Nonterminal: nonterminal, Rule: r}
			matches = append(matches,
				parseRule(text, match, hypotheses_limit)...)
		}
	}

	return matches
}
