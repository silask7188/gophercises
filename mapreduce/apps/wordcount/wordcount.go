package wordcount

import (
	"silask7188/mapreduce/mr"
	"strconv"
	"strings"
	"unicode"
)

func Map(filename string, contents string) mr.KeyValues {
	words := splitAndCleanString(contents)

	var kv mr.KeyValues
	for _, w := range words {
		kv = append(kv, mr.KeyValue{Key: w, Value: "1"})
	}
	return kv
}

func Reduce(key string, values []string) string {
	return strconv.Itoa(len(values))
}

func splitAndCleanString(s string) []string {
	w := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return w
}
