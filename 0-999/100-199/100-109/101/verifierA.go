package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// testCase represents a single verifier test
type testCase struct {
	s string
	k int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := []testCase{
		{s: "aaaaa", k: 4},
		{s: "abacaba", k: 3},
		{s: "abc", k: 0},
		{s: "abcdefgh", k: 10},
		{s: "ababcd", k: 2},
	}
	for i, t := range tests {
		input := fmt.Sprintf("%s\n%d\n", t.s, t.k)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\nstderr: %s\n", i+1, err, stderr.String())
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(lines) < 2 {
			fmt.Printf("test %d: expected two output lines, got %d\n", i+1, len(lines))
			os.Exit(1)
		}
		reportedAns := strings.TrimSpace(lines[0])
		reportedStr := strings.TrimSpace(lines[1])

		expectedAns, errExp := expectedAnswer(t.s, t.k)
		if errExp != nil {
			fmt.Printf("test %d: internal error: %v\n", i+1, errExp)
			os.Exit(1)
		}
		if reportedAns != fmt.Sprint(expectedAns) {
			fmt.Printf("test %d: expected answer %d, got %s\n", i+1, expectedAns, reportedAns)
			os.Exit(1)
		}
		if !isSubsequence(t.s, reportedStr) {
			fmt.Printf("test %d: output string is not a subsequence\n", i+1)
			os.Exit(1)
		}
		if len(t.s)-len(reportedStr) > t.k {
			fmt.Printf("test %d: removed more than k characters\n", i+1)
			os.Exit(1)
		}
		if distinctCount(reportedStr) != expectedAns {
			fmt.Printf("test %d: output string has wrong number of distinct characters\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

// expectedAnswer computes the minimal number of distinct characters
func expectedAnswer(s string, k int) (int, error) {
	n := len(s)
	if k >= n {
		return 0, nil
	}
	freq := make([]int, 256)
	for i := 0; i < n; i++ {
		freq[s[i]]++
	}
	// sort frequencies descending
	counts := make([]int, 0, 256)
	for _, f := range freq {
		if f > 0 {
			counts = append(counts, f)
		}
	}
	for i := 0; i < len(counts); i++ {
		for j := i + 1; j < len(counts); j++ {
			if counts[j] > counts[i] {
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}
	remain := n - k
	ans := 0
	for i := 0; i < len(counts) && remain > 0; i++ {
		remain -= counts[i]
		ans++
	}
	if remain > 0 {
		return 0, fmt.Errorf("invalid k")
	}
	if ans < 0 {
		ans = 0
	}
	return ans, nil
}

func isSubsequence(s, sub string) bool {
	j := 0
	for i := 0; i < len(s) && j < len(sub); i++ {
		if s[i] == sub[j] {
			j++
		}
	}
	return j == len(sub)
}

func distinctCount(s string) int {
	seen := make(map[byte]struct{})
	for i := 0; i < len(s); i++ {
		seen[s[i]] = struct{}{}
	}
	return len(seen)
}
