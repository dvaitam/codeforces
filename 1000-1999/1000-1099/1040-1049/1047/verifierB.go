package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected int
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	maxv := 0
	for i := 0; i < n; i++ {
		x := rng.Intn(1_000_000_000) + 1
		y := rng.Intn(1_000_000_000) + 1
		if x+y > maxv {
			maxv = x + y
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return testCase{input: sb.String(), expected: maxv}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// deterministic edge cases
	cases = append(cases, testCase{input: "1\n1 1\n", expected: 2})
	cases = append(cases, testCase{input: "1\n1000000000 1000000000\n", expected: 2000000000})

	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "case %d: expected single integer got %q\n", i+1, out)
			os.Exit(1)
		}
		val, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse integer\n", i+1)
			os.Exit(1)
		}
		if val != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, tc.expected, val, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
