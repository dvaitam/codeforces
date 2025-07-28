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
	expected string
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(20)
	arr := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20)
		sum += arr[i]
	}
	need := sum - m
	if need < 0 {
		need = 0
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d", need)}
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		// deterministic cases
		{input: "1\n1 0\n5\n", expected: "5"},
		{input: "1\n3 10\n1 2 3\n", expected: "0"},
		{input: "1\n2 1\n1 1\n", expected: "1"},
	}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
