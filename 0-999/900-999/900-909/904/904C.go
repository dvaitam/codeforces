package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	possible := make([]bool, 26)
	for i := 0; i < 26; i++ {
		possible[i] = true
	}
	remaining := 26
	determined := false
	excessive := 0

	for i := 0; i < n; i++ {
		var typ, word string
		fmt.Fscan(reader, &typ, &word)

		if determined && i != n-1 {
			if typ != "." {
				excessive++
			}
		}

		switch typ {
		case ".":
			for j := 0; j < len(word); j++ {
				idx := int(word[j] - 'a')
				if possible[idx] {
					possible[idx] = false
					remaining--
				}
			}
		case "!":
			present := make([]bool, 26)
			for j := 0; j < len(word); j++ {
				present[int(word[j]-'a')] = true
			}
			for j := 0; j < 26; j++ {
				if possible[j] && !present[j] {
					possible[j] = false
					remaining--
				}
			}
		case "?":
			if i != n-1 {
				idx := int(word[0] - 'a')
				if possible[idx] {
					possible[idx] = false
					remaining--
				}
			}
		}

		if !determined && remaining == 1 {
			determined = true
		}
	}

	fmt.Fprintln(writer, excessive)
}
