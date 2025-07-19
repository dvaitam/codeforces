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
	fmt.Fscan(reader, &n)
	m := 2*n - 2
	inputs := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &inputs[i])
	}
	// identify the two candidates of length n-1
	var s1, s2 string
	for _, s := range inputs {
		if len(s) == n-1 {
			if s1 == "" {
				s1 = s
			} else {
				s2 = s
			}
		}
	}
	cand1 := s1 + string(s2[len(s2)-1])
	cand2 := s2 + string(s1[len(s1)-1])
	res := make([]byte, m)
	if test(cand1, inputs, res) {
		writer.Write(res)
	} else {
		test(cand2, inputs, res)
		writer.Write(res)
	}
}

// test attempts to assign 'S' or 'P' for each input based on candidate
func test(pref string, inputs []string, res []byte) bool {
	used := make(map[string]bool)
	for i, s := range inputs {
		if hasSuffix(pref, s) {
			if !used[s] {
				res[i] = 'S'
				used[s] = true
			} else {
				res[i] = 'P'
			}
		} else if hasPrefix(pref, s) {
			res[i] = 'P'
		} else {
			return false
		}
	}
	return true
}

func hasPrefix(s, t string) bool {
	if len(s) < len(t) {
		return false
	}
	return s[:len(t)] == t
}

func hasSuffix(s, t string) bool {
	if len(s) < len(t) {
		return false
	}
	return s[len(s)-len(t):] == t
}
