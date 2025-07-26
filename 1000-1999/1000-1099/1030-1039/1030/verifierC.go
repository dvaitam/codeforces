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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int(s[i] - '0')
	}
	for i := 1; i < n; i++ {
		target := 0
		for j := 0; j < i; j++ {
			target += a[j]
		}
		curr := 0
		cnt := 0
		ok := true
		for j := i; j < n; j++ {
			curr += a[j]
			if curr > target {
				ok = false
				break
			}
			if curr == target {
				curr = 0
				cnt++
			}
		}
		if ok && curr == 0 && cnt >= 1 {
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
