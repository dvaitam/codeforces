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

type testCase struct {
	input    string
	expected []string
}

func buildDeterministic() []testCase {
	cases := []testCase{}
	// Simple Yes case
	cases = append(cases, testCase{input: "1\n1\n2 3\n", expected: []string{"Yes"}})
	// Simple No case
	cases = append(cases, testCase{input: "1\n1\n1 1\n", expected: []string{"No"}})
	return cases
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	exp := make([]string, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(100) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		even, odd := 0, 0
		for j := 0; j < 2*n; j++ {
			v := rng.Intn(101)
			sb.WriteString(fmt.Sprintf("%d ", v))
			if v%2 == 0 {
				even++
			} else {
				odd++
			}
		}
		sb.WriteString("\n")
		if even == odd {
			exp[i] = "Yes"
		} else {
			exp[i] = "No"
		}
	}
	return testCase{input: sb.String(), expected: exp}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(tc.expected) {
		return fmt.Errorf("expected %d tokens got %d", len(tc.expected), len(fields))
	}
	for i, f := range fields {
		if !strings.EqualFold(f, tc.expected[i]) {
			return fmt.Errorf("expected %s got %s", tc.expected[i], f)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := buildDeterministic()
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
