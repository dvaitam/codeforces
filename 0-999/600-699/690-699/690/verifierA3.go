package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type scenario struct {
	n    int64
	r    int64
	seen []int64
}

type testCase struct {
	input  string
	cases  []scenario
	expect []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA3.go /path/to/binary")
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
		answerLines := strings.Fields(out)
		if len(answerLines) != len(tc.cases) {
			fmt.Printf("test %d failed: expected %d answers but got %d\ninput:\n%s\noutput:\n%s\n", idx+1, len(tc.cases), len(answerLines), tc.input, out)
			os.Exit(1)
		}
		guesses := make([]int64, len(answerLines))
		for i, tok := range answerLines {
			val, ok := parseInt64(tok)
			if !ok {
				fmt.Printf("test %d failed: invalid integer %q\ninput:\n%s\noutput:\n%s\n", idx+1, tok, tc.input, out)
				os.Exit(1)
			}
			guesses[i] = val
		}
		if err := checkStrategy(tc.cases, guesses); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func parseInt64(s string) (int64, bool) {
	var sign int64 = 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i++
	}
	var val int64
	for ; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			return 0, false
		}
		val = val*10 + int64(ch-'0')
	}
	return val * sign, true
}

func checkStrategy(cases []scenario, guesses []int64) error {
	for i, sc := range cases {
		if guesses[i] < 1 || guesses[i] > sc.n {
			return fmt.Errorf("scenario %d: guess %d out of range [1,%d]", i+1, guesses[i], sc.n)
		}
	}
	if flawlesslyWins(cases, guesses) {
		return nil
	}
	return fmt.Errorf("strategy does not guarantee a correct guess")
}

func flawlesslyWins(cases []scenario, guesses []int64) bool {
	// Simulate the core logic: reconstruct full assignment for each T by inserting candidate self-number
	for idx, sc := range cases {
		n := sc.n
		full := make([]int64, 0, n)
		insertPos := sc.r - 1
		cur := int64(0)
		for i := int64(0); i < n; i++ {
			if i == insertPos {
				full = append(full, guesses[idx])
			} else {
				full = append(full, sc.seen[cur])
				cur++
			}
		}
		if !verifyAssignment(full) {
			return false
		}
	}
	return true
}

func verifyAssignment(nums []int64) bool {
	n := int64(len(nums))
	tot := int64(0)
	for _, v := range nums {
		tot += v
	}
	mod := tot % n
	for i, val := range nums {
		r := int64(i+1) % n
		if r == 0 {
			r = n
		}
		delta := (r - mod + n) % n
		expected := delta
		if expected == 0 {
			expected = n
		}
		if val == expected {
			return true
		}
	}
	return false
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

func genTests() []testCase {
	rand.Seed(42)
	var tests []testCase
	tests = append(tests, buildCase([]scenario{
		{n: 2, r: 1, seen: []int64{1}},
		{n: 2, r: 2, seen: []int64{2}},
	}))
	for t := 0; t < 200; t++ {
		cases := make([]scenario, 0)
		T := rand.Intn(10) + 1
		for i := 0; i < T; i++ {
			n := int64(rand.Intn(5) + 2)
			r := int64(rand.Intn(int(n)) + 1)
			vals := make([]int64, 0, n-1)
			for k := int64(0); k < n-1; k++ {
				vals = append(vals, int64(rand.Intn(int(n))+1))
			}
			cases = append(cases, scenario{n: n, r: r, seen: vals})
		}
		tests = append(tests, buildCase(cases))
	}
	return tests
}

func buildCase(cases []scenario) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, sc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", sc.n, sc.r))
		for i, v := range sc.seen {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	expect := make([]int64, len(cases))
	refInput := sb.String()
	refOut, _ := runRef(refInput)
	fields := strings.Fields(refOut)
	for i, tok := range fields {
		val, _ := parseInt64(tok)
		expect[i] = val
	}
	return testCase{
		input:  refInput,
		cases:  cases,
		expect: expect,
	}
}

func runRef(input string) (string, error) {
	cmd := exec.Command("go", "run", "690A3.go")
	cmd.Dir = "."
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
