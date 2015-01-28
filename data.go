package main

var nonterminals = map[string][]string{
	"sentence": {
		"{place_adverb} {part_of_speech=verb number=@2} {extended_objs case=@1 number=@2}, {extended_participle case=@1 number=@2}.",
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
		"{part_of_speech=adjective case=@1 number=@2} {part_of_speech=noun !case=@1 !number=@2} {extended_objs case=genetive}",
	},
	"extended_participle": {
		"{part_of_speech=participle !case !number} {extended_objs case=instrumental}",
	},
}

var terminals = map[string][]Attribute{
	"больничном": {
		{Name: "case", Value: "prepositional"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "adjective"},
	},
	"дворе": {
		{Name: "case", Value: "prepositional"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "noun"},
	},
	"стоит": {
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "verb"},
	},
	"небольшой": {
		{Name: "case", Value: "nominative"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "adjective"},
	},
	"флигель": {
		{Name: "case", Value: "nominative"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "noun"},
	},
	"окруженный": {
		{Name: "case", Value: "nominative"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "participle"},
	},
	"целым": {
		{Name: "case", Value: "instrumental"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "adjective"},
	},
	"лесом": {
		{Name: "case", Value: "instrumental"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "noun"},
	},
	"репейника": {
		{Name: "case", Value: "genetive"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "noun"},
	},
	"крапивы": {
		{Name: "case", Value: "genetive"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "noun"},
	},
	"дикой": {
		{Name: "case", Value: "genetive"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "adjective"},
	},
	"конопли": {
		{Name: "case", Value: "genetive"},
		{Name: "number", Value: "single"},
		{Name: "part_of_speech", Value: "noun"},
	},
}

func FindNonterminalRules(nonterminal string) []string {
	return nonterminals[nonterminal]
}

func FindTerminals(prefix string) []ParseMatch {
	return []ParseMatch{{Text: prefix, Attributes: terminals[prefix]}}
}
