package main

var nonterminals = map[string][]string{
	"sentence": {
		"{place_adverb} {part_of_speech=verb number=@2} {extended_objs case=@1 number=@2} {extended_participle case=@1 number=@2}.",
	},
	"place_adverb": {
		"В {extended_objs case=prepositional}",
	},
	"extended_objs": {
		"{extended_obj !case !number}",
		"@{number=plural}{extended_obj !case=@1}, {extended_objs case=@1}",
		"@{number=plural}{extended_obj !case=@1} и {extended_objs case=@1}",
	},
	"extended_obj": {
		"{part_of_speech=noun !case !number}",
		"{part_of_speech=adjective case=@1 number=@2} {part_of_speech=noun !case=@1 !number=@2}",
		"{part_of_speech=adjective case=@1 number=@2} {part_of_speech=noun !case=@1 !number=@2} {extended_obj case=genetive}",
	},
	"extended_participle": {
		"{part_of_speech=participle !case !number} {extended_objs case=instrumental}",
	},
}

var terminals = map[string][]Attribute{
	"больничном": {
		{Name: "case", Values: []string{"prepositional"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"adjective"}},
	},
	"дворе": {
		{Name: "case", Values: []string{"prepositional"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
	"стоит": {
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
	"небольшой": {
		{Name: "case", Values: []string{"nominative"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"adjective"}},
	},
	"флигель": {
		{Name: "case", Values: []string{"nominative"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
	"окруженный": {
		{Name: "case", Values: []string{"nominative"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"participle"}},
	},
	"целым": {
		{Name: "case", Values: []string{"instrumental"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"adjective"}},
	},
	"лесом": {
		{Name: "case", Values: []string{"instrumental"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
	"репейника": {
		{Name: "case", Values: []string{"genetive"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
	"крапивы": {
		{Name: "case", Values: []string{"genetive"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
	"дикой": {
		{Name: "case", Values: []string{"genetive"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"adjective"}},
	},
	"конопли": {
		{Name: "case", Values: []string{"genetive"}},
		{Name: "number", Values: []string{"single"}},
		{Name: "part_of_speech", Values: []string{"noun"}},
	},
}

func FindNonterminalRules(nonterminal string) []string {
	return nonterminals[nonterminal]
}

func FindTerminalAttrs(terminal string) []Attribute {
	return terminals[terminal]
}
