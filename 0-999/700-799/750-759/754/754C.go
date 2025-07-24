package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func mentions(text, user string) bool {
	for i := 0; ; {
		idx := strings.Index(text[i:], user)
		if idx == -1 {
			return false
		}
		idx += i
		beforeOK := idx == 0 || !isWordChar(rune(text[idx-1]))
		afterIdx := idx + len(user)
		afterOK := afterIdx == len(text) || !isWordChar(rune(text[afterIdx]))
		if beforeOK && afterOK {
			return true
		}
		i = idx + 1
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		names := make([]string, n)
		nameIndex := make(map[string]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &names[i])
			nameIndex[names[i]] = i
		}
		var m int
		fmt.Fscan(reader, &m)
		reader.ReadString('\n')

		prefixes := make([]string, m)
		texts := make([]string, m)
		for i := 0; i < m; i++ {
			line, _ := reader.ReadString('\n')
			line = strings.TrimRight(line, "\r\n")
			idx := strings.IndexByte(line, ':')
			if idx == -1 {
				prefixes[i] = line
				texts[i] = ""
			} else {
				prefixes[i] = line[:idx]
				texts[i] = line[idx+1:]
			}
		}

		allowed := make([][]int, m)
		for i := 0; i < m; i++ {
			if prefixes[i] != "?" {
				if id, ok := nameIndex[prefixes[i]]; ok {
					allowed[i] = []int{id}
				} else {
					allowed[i] = nil
				}
			} else {
				cand := []int{}
				for j, nm := range names {
					if !mentions(texts[i], nm) {
						cand = append(cand, j)
					}
				}
				allowed[i] = cand
			}
		}

		dp := make([][]int, m)
		for i := range dp {
			dp[i] = make([]int, n)
			for j := range dp[i] {
				dp[i][j] = -1
			}
		}
		if len(allowed[0]) == 0 {
			fmt.Fprintln(writer, "Impossible")
			continue
		}
		for _, u := range allowed[0] {
			dp[0][u] = -2
		}
		for i := 1; i < m; i++ {
			for _, u := range allowed[i] {
				for v := 0; v < n; v++ {
					if dp[i-1][v] != -1 && v != u {
						dp[i][u] = v
						break
					}
				}
			}
		}
		endUser := -1
		for u := 0; u < n; u++ {
			if dp[m-1][u] != -1 {
				endUser = u
				break
			}
		}
		if endUser == -1 {
			fmt.Fprintln(writer, "Impossible")
			continue
		}
		result := make([]int, m)
		cur := endUser
		for i := m - 1; i >= 0; i-- {
			result[i] = cur
			cur = dp[i][cur]
		}
		for i := 0; i < m; i++ {
			fmt.Fprintf(writer, "%s:%s\n", names[result[i]], texts[i])
		}
	}
}
