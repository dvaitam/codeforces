package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func naiveCount(a, b string) int {
	palsA := make(map[string]struct{})
	for i := 0; i < len(a); i++ {
		for j := i + 1; j <= len(a); j++ {
			s := a[i:j]
			if isPalindrome(s) {
				palsA[s] = struct{}{}
			}
		}
	}
	palsB := make(map[string]struct{})
	for i := 0; i < len(b); i++ {
		for j := i + 1; j <= len(b); j++ {
			s := b[i:j]
			if isPalindrome(s) {
				palsB[s] = struct{}{}
			}
		}
	}
	seen := make(map[string]struct{})
	for sa := range palsA {
		for sb := range palsB {
			seen[sa+sb] = struct{}{}
		}
	}
	return len(seen)
}

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 1
	m := r.Intn(5) + 1
	letters := "abc"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[r.Intn(len(letters))])
	}
	a := sb.String()
	sb.Reset()
	for i := 0; i < m; i++ {
		sb.WriteByte(letters[r.Intn(len(letters))])
	}
	b := sb.String()
	return fmt.Sprintf("%s\n%s\n", a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	cases := []string{
		"a\na\n",
		"ab\naba\n",
		"abc\nabc\n",
	}
	for i := 0; i < 97; i++ {
		cases = append(cases, genCase(r))
	}
	for idx, input := range cases {
		scanner := bufio.NewScanner(strings.NewReader(input))
		scanner.Scan()
		A := scanner.Text()
		scanner.Scan()
		B := scanner.Text()
		want := fmt.Sprintf("%d", naiveCount(A, B))
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
