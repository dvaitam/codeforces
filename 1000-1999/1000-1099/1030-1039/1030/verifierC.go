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
	expected string
}

func solveTicket(s string) string {
	n := len(s)
	for i := 0; i < n-1; i++ {
		sum := 0
		for k := 0; k <= i; k++ {
			sum += int(s[k] - '0')
		}

		possible := true
		curr := 0
		for j := i + 1; j < n; j++ {
			curr += int(s[j] - '0')
			if curr == sum {
				curr = 0
			} else if curr > sum {
				possible = false
				break
			}
		}

		if possible && curr == 0 {
			return "YES\n"
		}
	}
	return "NO\n"
}

func buildCase(s string) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(s)))
	sb.WriteByte('\n')
	sb.WriteString(s)
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solveTicket(s)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	digits := make([]byte, n)
	for i := range digits {
		digits[i] = byte(rng.Intn(10)) + '0'
	}
	return buildCase(string(digits))
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.String() != tc.expected {
		return fmt.Errorf("expected %q got %q", tc.expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase("350178"),
		buildCase("99"),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
