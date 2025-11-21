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

const refPath = "2000-2999/2100-2199/2130-2139/2136/2136A.go"

type testCase struct {
	a, b, c, d int
}

func main() {
	if len(os.Args) != 2 {
		if len(os.Args) == 3 && os.Args[1] == "--" {
			os.Args = []string{os.Args[0], os.Args[2]}
		} else {
			fmt.Println("usage: go run verifierA.go /path/to/binary")
			os.Exit(1)
		}
	}
	bin := os.Args[1]

	tests := buildTests()
	input := renderInput(tests)

	expRaw, err := runBinary(refPath, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}
	actRaw, err := runBinary(bin, input)
	if err != nil {
		fmt.Printf("runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	exp := normalize(expRaw, len(tests))
	act := normalize(actRaw, len(tests))
	if len(exp) != len(tests) {
		fmt.Printf("reference produced %d lines, expected %d\n", len(exp), len(tests))
		os.Exit(1)
	}
	if len(act) != len(tests) {
		fmt.Printf("output has %d lines, expected %d\ninput:\n%s\nactual:\n%s\n", len(act), len(tests), input, actRaw)
		os.Exit(1)
	}
	for i := range tests {
		if exp[i] != act[i] {
			fmt.Printf("case %d mismatch: expected %s got %s\ninput case: %d %d %d %d\nfull input:\n%s\n", i+1, exp[i], act[i], tests[i].a, tests[i].b, tests[i].c, tests[i].d, input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func buildTests() []testCase {
	rand.Seed(time.Now().UnixNano())
	var tests []testCase

	// Exhaustive small scores
	for a := 0; a <= 3; a++ {
		for b := 0; b <= 3; b++ {
			for c := a; c <= a+3; c++ {
				for d := b; d <= b+3; d++ {
					tests = append(tests, testCase{a, b, c, d})
				}
			}
		}
	}

	// Edge cases
	edges := []testCase{
		{0, 0, 0, 0},
		{0, 0, 100, 0},
		{0, 0, 0, 100},
		{0, 0, 100, 100},
		{2, 2, 2, 2},
		{2, 2, 3, 4},
		{1, 4, 1, 4},
		{1, 4, 4, 4},
		{0, 3, 0, 3},
	}
	tests = append(tests, edges...)

	// Random cases within constraints
	for i := 0; i < 200; i++ {
		a := rand.Intn(101)
		b := rand.Intn(101)
		c := a + rand.Intn(101-a)
		d := b + rand.Intn(101-b)
		tests = append(tests, testCase{a, b, c, d})
	}

	return tests
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func normalize(out string, t int) []string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	res := make([]string, 0, t)
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		res = append(res, strings.ToUpper(ln))
	}
	return res
}
