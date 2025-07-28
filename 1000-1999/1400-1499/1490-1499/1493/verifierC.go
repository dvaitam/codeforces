package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	n int
	k int
	s string
}

func generateTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 100)
	edge := []Test{{1, 1, "a"}, {2, 1, "ab"}, {3, 3, "abc"}, {4, 2, "aaaa"}, {5, 5, "abcde"}}
	tests = append(tests, edge...)
	letters := []rune("abcde")
	for len(tests) < 100 {
		n := rand.Intn(50) + 1
		k := rand.Intn(n) + 1
		runes := make([]rune, n)
		for i := 0; i < n; i++ {
			runes[i] = letters[rand.Intn(len(letters))]
		}
		tests = append(tests, Test{n, k, string(runes)})
	}
	return tests
}

func solve(n, k int, s string) string {
	if n%k != 0 {
		return "-1"
	}
	freq := make([]int, 26)
	for i := 0; i < n; i++ {
		freq[s[i]-'a']++
	}
	good := true
	for i := 0; i < 26; i++ {
		if freq[i]%k != 0 {
			good = false
		}
	}
	if good {
		return s
	}
	for i := n - 1; i >= 0; i-- {
		cur := int(s[i] - 'a')
		freq[cur]--
		for c := cur + 1; c < 26; c++ {
			freq[c]++
			rem := n - i - 1
			needSum := 0
			for j := 0; j < 26; j++ {
				needSum += (k - freq[j]%k) % k
			}
			if needSum <= rem && (rem-needSum)%k == 0 {
				prefix := s[:i] + string(rune('a'+c))
				rest := make([]byte, 0, rem)
				for t := 0; t < rem-needSum; t++ {
					rest = append(rest, 'a')
				}
				for j := 0; j < 26; j++ {
					cnt := (k - freq[j]%k) % k
					for t := 0; t < cnt; t++ {
						rest = append(rest, byte('a'+j))
					}
				}
				return prefix + string(rest)
			}
			freq[c]--
		}
	}
	return "-1"
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := generateTests()
	var in strings.Builder
	fmt.Fprintln(&in, len(tests))
	for _, t := range tests {
		fmt.Fprintf(&in, "%d %d\n%s\n", t.n, t.k, t.s)
	}
	expectedParts := make([]string, len(tests))
	for i, t := range tests {
		expectedParts[i] = solve(t.n, t.k, t.s) + "\n"
	}
	expect := strings.Join(expectedParts, "")

	got, err := run(binary, in.String())
	if err != nil {
		fmt.Printf("runtime error: %v\noutput:\n%s", err, got)
		os.Exit(1)
	}
	got = strings.ReplaceAll(strings.TrimSpace(got), "\r\n", "\n")
	expect = strings.ReplaceAll(strings.TrimSpace(expect), "\r\n", "\n")
	if got != expect {
		fmt.Println("wrong answer")
		fmt.Println("input:")
		fmt.Print(in.String())
		fmt.Println("expected:")
		fmt.Print(expect)
		fmt.Println("got:")
		fmt.Print(got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
	time.Sleep(0)
}
