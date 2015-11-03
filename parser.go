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
	Rule            Rule
	Attributes      []Attribute
	Submatches      []ParseMatch
	Hypotheses      []string
	HypothesisCount uint
}

const TerminalSeparators = " ,."

func parseTerminal(pattern string) (string, string) {
	var term string
	found, unescape := false, false
	for i := 0; i < len(pattern); i++ {
		escaped := i > 0 && pattern[i-1] == '\\'
		if escaped {
			unescape = true
		}
		if pattern[i] == '{' && !escaped {
			term, pattern = pattern[:i], pattern[i:]
			found = true
			break
		}
	}

	if !found {
		term, pattern = pattern, ""
	}

	if unescape {
		term = strings.Replace(term, "\\{", "{", -1)
	}

	return term, pattern
}

func parseNonterminal(pattern string) (string, []Attribute, string) {
	attrs := []Attribute{}

	if len(pattern) == 0 || pattern[0] != '{' {
		return "", attrs, pattern
	}

	var body string
	if i := strings.Index(pattern, "}"); i != -1 {
		body, pattern = pattern[1:i], pattern[i+1:]
	} else {
		body, pattern = pattern, ""
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

	return nterm, attrs, pattern
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
	if !ok || len(attr2.Value) == 0 || attr2.Value == attr.Value {
		return true
	}

	*hypo = fmt.Sprintf("%s = %s", attr.Name, attr2.Value)
	return false
}

func (match *ParseMatch) checkAttrToken(attr Attribute, hypo *string) bool {
	if len(attr.Value) == 0 {
		return true
	}

	for i := 0; i < len(match.Attributes); i++ {
		a := &match.Attributes[i]
		if a.Name == attr.Name && a.token == attr.token {
			if attr.export {
				a.export = true
			}

			if len(a.Value) > 0 {
				if a.Value == attr.Value {
					return true
				}
			} else {
				a.Value = attr.Value
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

	match.HypothesisCount += submatch.HypothesisCount
	if match.HypothesisCount > hypotheses_limit {
		return false
	}

	var hypo string
	for _, a := range attrs {
		attr, ok := submatch.findAttr(a.Name)
		if !ok {
			attr.Name = a.Name
		}

		if len(a.token) > 0 {
			a.Value = attr.Value
		}

		if (len(a.Value) > 0 && !submatch.checkAttrValue(a, &hypo)) ||
			(len(a.token) > 0 && !match.checkAttrToken(a, &hypo)) {
			match.HypothesisCount++
			if match.HypothesisCount > hypotheses_limit {
				return false
			}
			hypo = fmt.Sprintf(
				"%d: %s", len(match.Submatches)+1, hypo)
			match.Hypotheses = append(match.Hypotheses, hypo)
		}

		if a.export && len(a.token) == 0 {
			match.Attributes = append(match.Attributes, attr)
		}
	}

	match.Text += submatch.Text
	match.Submatches = append(match.Submatches, submatch)
	return true
}

func parsePatternPart(text, pattern string, match ParseMatch,
	hypotheses_limit uint, output chan []ParseMatch) {

	term, pattern := parseTerminal(pattern)
	if !strings.HasPrefix(text, term) {
		output <- nil
		return
	}
	match.Text += term
	text = text[len(term):]

	if len(pattern) == 0 {
		match.cleanExport()
		output <- []ParseMatch{match}
		return
	}

	pred, attrs, pattern := parseNonterminal(pattern)
	output2 := make(chan []ParseMatch)
	count := 0
	for _, m := range Parse(text, pred, hypotheses_limit) {
		match2 := match.clone()
		if match2.checkSubmatch(attrs, m, hypotheses_limit) {
			go parsePatternPart(text[len(m.Text):],
				pattern, match2, hypotheses_limit, output2)
			count++
		}
	}

	matches := []ParseMatch{}
	for i := 0; i < count; i++ {
		matches = append(matches, <-output2...)
	}

	output <- matches
}

func parsePattern(text string, match ParseMatch,
	hypotheses_limit uint) []ParseMatch {

	pattern := match.Rule.Pattern
	if strings.HasPrefix(pattern, "@") {
		_, match.Attributes, pattern = parseNonterminal(pattern[1:])
	}

	output := make(chan []ParseMatch)
	go parsePatternPart(text, pattern, match, hypotheses_limit, output)
	return <-output
}

func Parse(text, nonterminal string, hypotheses_limit uint) []ParseMatch {
	matches, ok := FindInCache(text, nonterminal, hypotheses_limit)
	if ok {
		return matches
	}

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
				parsePattern(text, match, hypotheses_limit)...)
		}
	}

	AddToCache(text, nonterminal, hypotheses_limit, matches)
	return matches
}
