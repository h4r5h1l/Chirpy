package main

import (
	"strings"
)

func replace_words(input string) string {
	bad := []string{"kerfuffle", "sharbert", "fornax"}
	parts := strings.Split(input, " ")
	for i, word := range parts {
		for _, bad_word := range bad {
			if strings.ToLower(word) == bad_word {
				parts[i] = "****"
			}
		}
	}
	return strings.Join(parts, " ")
}
