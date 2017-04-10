package util

import "regexp"

func DecodeNamedSubmatches(matches [][]string, re *regexp.Regexp) [][]string {
	groupNames := re.SubexpNames()
	results := make([][]string, 10)

	for _, group := range matches {
		for nameIdx, name := range groupNames[1:] {
			if group[nameIdx+1] != "" {
				results = append(results, []string{name, group[nameIdx+1]})
			}
		}
	}

	return results
}
