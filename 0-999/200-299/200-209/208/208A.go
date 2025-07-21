package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Read the entire input line
	input, err := reader.ReadString('\n')
	if err != nil && len(input) == 0 {
		return
	}
	s := strings.TrimSpace(input)
	// Split by "WUB" and filter out empty tokens
	parts := strings.Split(s, "WUB")
	var words []string
	for _, p := range parts {
		if p != "" {
			words = append(words, p)
		}
	}
	// Join words with a single space
	if len(words) > 0 {
		fmt.Println(strings.Join(words, " "))
	}
}
