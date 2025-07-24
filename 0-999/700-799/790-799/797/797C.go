package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	freq := make([]int, 26)
	for i := 0; i < len(s); i++ {
		freq[s[i]-'a']++
	}

	// find smallest index with remaining letters
	minIdx := 0
	for minIdx < 26 && freq[minIdx] == 0 {
		minIdx++
	}

	stack := make([]byte, 0, len(s))
	result := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		c := s[i]
		stack = append(stack, c)
		freq[c-'a']--
		for minIdx < 26 && freq[minIdx] == 0 {
			minIdx++
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if minIdx >= 26 || top <= byte('a'+minIdx) {
				result = append(result, top)
				stack = stack[:len(stack)-1]
			} else {
				break
			}
		}
	}

	for i := len(stack) - 1; i >= 0; i-- {
		result = append(result, stack[i])
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, string(result))
	writer.Flush()
}
