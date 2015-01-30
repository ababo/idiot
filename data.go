package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

var db *sql.DB = nil

func InitData(db_filename string) error {
	db2, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		return err
	}

	if db != nil {
		FinalizeData()
	}
	db = db2

	return nil
}

func FinalizeData() {
	if db != nil {
		db.Close()
	}
}

func createTerminalTable() error {
	db.Exec("DROP TABLE terminal")

	if _, err := db.Exec(
		`CREATE TABLE terminal(text TEXT NOT NULL,
		                       pos TEXT NOT NULL,
		                       number TEXT, "case" TEXT)`); err != nil {
		return err
	}

	if _, err := db.Exec(
		"CREATE INDEX terminal_text ON terminal(text)"); err != nil {
		return err
	}

	return nil
}

var terminalAttrValues = map[string][2]string{
	"NOUN": {"pos", "noun"},
	"ADVB": {"pos", "advb"},
	"ADJF": {"pos", "adjf"},
	"ADJS": {"pos", "adjs"},
	"COMP": {"pos", "comp"},
	"VERB": {"pos", "verb"},
	"INFN": {"pos", "infn"},
	"PRTF": {"pos", "prtf"},
	"PRTS": {"pos", "prts"},
	"GRND": {"pos", "grnd"},
	"CONJ": {"pos", "conj"},
	"INTJ": {"pos", "intj"},
	"PRCL": {"pos", "prcl"},
	"PREP": {"pos", "prep"},
	"PRED": {"pos", "pred"},
	"NUMR": {"pos", "numr"},
	"NPRO": {"pos", "npro"},

	"sing": {"number", "sing"},
	"plur": {"number", "plur"},

	"nomn": {"case", "nomn"},
	"gen1": {"case", "gen1"},
	"gen2": {"case", "gen2"},
	"datv": {"case", "datv"},
	"accs": {"case", "accs"},
	"ablt": {"case", "ablt"},
	"loct": {"case", "loct"},
	"loc1": {"case", "loc1"},
	"loc2": {"case", "loc2"},
	"voct": {"case", "voct"},
}

func loadTerminalTable(file *os.File) (skipped []string, err error) {
	scanner := bufio.NewScanner(file)
	skipped2 := map[string]bool{}

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "\t")
		if len(split) < 2 {
			continue
		}

		text := strings.Replace(split[0], "'", "''", -1)
		text = strings.ToLower(text)

		split = strings.FieldsFunc(split[1], func(r rune) bool {
			return r == ',' || r == ' '

		})

		var names, values string
		for _, a := range split {
			if pair, ok := terminalAttrValues[a]; ok {
				names += ",'" + pair[0] + "'"
				values += ",'" + pair[1] + "'"
			} else {
				skipped2[a] = true
			}
		}

		query := fmt.Sprintf(
			"INSERT INTO terminal(text%s) VALUES('%s'%s)",
			names, text, values)
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	skipped = make([]string, 0, len(skipped2))
	for k := range skipped2 {
		skipped = append(skipped, k)
	}

	return skipped, nil
}

func buildTerminalData(txt_filename string) (skipped []string, err error) {
	file, err := os.Open(txt_filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = createTerminalTable(); err != nil {
		return nil, err
	}

	skipped, err = loadTerminalTable(file)
	if err != nil {
		return nil, err
	}

	return skipped, nil
}

func appendTerminalMatches(matches []ParseMatch, rows *sql.Rows) []ParseMatch {
	for rows.Next() {
		var text, pos, number, case_ string
		rows.Scan(&text, &pos, &number, &case_)

		match := ParseMatch{Text: text}
		if len(pos) > 0 {
			match.Attributes = append(
				match.Attributes,
				Attribute{Name: "pos", Value: pos})
		}
		if len(number) > 0 {
			match.Attributes = append(
				match.Attributes,
				Attribute{Name: "number", Value: number})
		}
		if len(case_) > 0 {
			match.Attributes = append(
				match.Attributes,
				Attribute{Name: "case", Value: case_})
		}

		matches = append(matches, match)
	}

	return matches
}

func FindTerminals(prefix, separator string) []ParseMatch {
	if db == nil || len(prefix) == 0 {
		return nil
	}

	matches := []ParseMatch{}

	rows, err := db.Query(
		"SELECT DISTINCT * FROM terminal WHERE text=?", prefix)
	if err != nil {
		return nil
	}
	defer rows.Close()
	matches = appendTerminalMatches(matches, rows)

	if len(separator) == 0 {
		return matches
	}

	from := prefix + separator
	to := from[:len(from)-1] + string(from[len(from)-1]+1)

	rows2, err2 := db.Query(
		"SELECT DISTINCT * FROM terminal WHERE text>=? AND text<?",
		from, to)
	if err2 != nil {
		return nil
	}
	defer rows2.Close()
	matches = appendTerminalMatches(matches, rows2)

	return matches
}

////////

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

func FindNonterminalRules(nonterminal string) []string {
	return nonterminals[nonterminal]
}
