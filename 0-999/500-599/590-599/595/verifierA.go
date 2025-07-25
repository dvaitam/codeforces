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

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	input  string
	expect int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			a := rng.Intn(2)
			b := rng.Intn(2)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d %d", a, b))
			if a == 1 || b == 1 {
				cnt++
			}
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expect: cnt}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// simple deterministic case
	cases = append(cases, testCase{input: "1 1\n0 0\n", expect: 0})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected single integer output, got %q\ninput:\n%s", i+1, out, tc.input)
			os.Exit(1)
		}
		val, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid integer output\ninput:\n%soutput:\n%s", i+1, tc.input, out)
			os.Exit(1)
		}
		if val != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d, got %d\ninput:\n%s", i+1, tc.expect, val, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
