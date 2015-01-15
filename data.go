package main

var predicates = map[string][]string{
	"text": {
		"{sentence} {text}",
		"{sentence}",
	},
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

var attributes = map[string]map[string][]string{
	"больничном": {
		"case":           {"prepositional"},
		"number":         {"single"},
		"part_of_speech": {"adjective"},
	},
	"дворе": {
		"case":           {"prepositional"},
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
	"стоит": {
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
	"небольшой": {
		"case":           {"nominative"},
		"number":         {"single"},
		"part_of_speech": {"adjective"},
	},
	"флигель": {
		"case":           {"nominative"},
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
	"окруженный": {
		"case":           {"nominative"},
		"number":         {"single"},
		"part_of_speech": {"participle"},
	},
	"целым": {
		"case":           {"instrumental"},
		"number":         {"single"},
		"part_of_speech": {"adjective"},
	},
	"лесом": {
		"case":           {"instrumental"},
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
	"репейника": {
		"case":           {"genetive"},
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
	"крапивы": {
		"case":           {"genetive"},
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
	"дикой": {
		"case":           {"genetive"},
		"number":         {"single"},
		"part_of_speech": {"adjective"},
	},
	"конопли": {
		"case":           {"genetive"},
		"number":         {"single"},
		"part_of_speech": {"noun"},
	},
}

func FindPredicateRules(predicate string) []string {
	return predicates[predicate]
}

func FindTerminalAttr(terminal string, attr string) ([]string, bool) {
	attrs, ok := attributes[terminal]
	if !ok {
		return nil, false
	}
	val, ok := attrs[attr]
	return val, ok
}
