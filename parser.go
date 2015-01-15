package main

type RuleMatch struct {
	Text       string
	Predicate  string
	Rule       string
	Submatches []RuleMatch
	Export     map[string]string
}

func Parse(predicate string, text string) []RuleMatch {
	return nil
}
