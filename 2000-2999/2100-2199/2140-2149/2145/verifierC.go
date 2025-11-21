package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.Atoi(expect)
	val, err := strconv.Atoi(actual)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase("bbbab"),
		makeCase("bbaaba"),
		makeCase("aaaa"),
		makeCase("aabbaaabbaab"),
		makeCase("aabaa"),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(10) + 2
		s := randomString(n)
		tests = append(tests, makeCase(s))
	}
	return tests
}

func randomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			sb.WriteByte('a')
		} else {
			sb.WriteByte('b')
		}
	}
	return sb.String()
}

func makeCase(s string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n%s\n", len(s), s)
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(s)),
	}
}

func solveReference(s string) int {
	diff := 0
	for _, ch := range s {
		if ch == 'a' {
			diff++
		} else {
			diff--
		}
	}
	if diff == 0 {
		return 0
	}
	prefix := 0
	seen := map[int]int{0: 0}
	best := len(s) + 1
	for i, ch := range s {
		if ch == 'a' {
			prefix++
		} else {
			prefix--
		}
		if idx, ok := seen[prefix-diff]; ok {
			if length := i + 1 - idx; length < best {
				best = length
			}
		}
		if _, ok := seen[prefix]; !ok {
			seen[prefix] = i + 1
		}
	}
	if best == len(s)+1 {
		return -1
	}
	return best
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
