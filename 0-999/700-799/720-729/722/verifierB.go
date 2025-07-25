package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n       int
	pattern []int
	lines   []string
}

func genTestsB() []testCaseB {
	rand.Seed(43)
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		pat := make([]int, n)
		lines := make([]string, n)
		for j := 0; j < n; j++ {
			pat[j] = rand.Intn(6)
			words := rand.Intn(3) + 1
			var sb strings.Builder
			for w := 0; w < words; w++ {
				if w > 0 {
					sb.WriteByte(' ')
				}
				l := rand.Intn(5) + 1
				for k := 0; k < l; k++ {
					ch := byte('a' + rand.Intn(26))
					sb.WriteByte(ch)
				}
			}
			lines[j] = sb.String()
		}
		tests[i] = testCaseB{n, pat, lines}
	}
	return tests
}

func countVowels(s string) int {
	cnt := 0
	for _, r := range s {
		switch r {
		case 'a', 'e', 'i', 'o', 'u', 'y':
			cnt++
		}
	}
	return cnt
}

func solveB(tc testCaseB) string {
	for i, line := range tc.lines {
		c := countVowels(line)
		if c != tc.pattern[i] {
			return "NO"
		}
	}
	return "YES"
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j, v := range tc.pattern {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for _, line := range tc.lines {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp := solveB(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
