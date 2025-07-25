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
	in  string
	out string
}

func solveCase(s string) string {
	pos := strings.IndexByte(s, '.')
	if pos == -1 {
		pos = len(s)
	}
	digits := strings.ReplaceAll(s, ".", "")
	first := 0
	for first < len(digits) && digits[first] == '0' {
		first++
	}
	last := len(digits) - 1
	for last >= 0 && digits[last] == '0' {
		last--
	}
	digits = digits[first : last+1]
	exponent := pos - (first + 1)
	var sb strings.Builder
	sb.WriteByte(digits[0])
	if len(digits) > 1 {
		sb.WriteByte('.')
		sb.WriteString(digits[1:])
	}
	if exponent != 0 {
		sb.WriteByte('E')
		sb.WriteString(strconv.Itoa(exponent))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildCase(s string) testCase {
	return testCase{in: s + "\n", out: solveCase(s)}
}

func randomCase(rng *rand.Rand) testCase {
	length := rng.Intn(25) + 1
	digits := make([]byte, length)
	for i := range digits {
		digits[i] = byte('0' + rng.Intn(10))
	}
	// ensure at least one non-zero
	digits[rng.Intn(length)] = byte('1' + rng.Intn(9))
	if rng.Intn(2) == 0 && length > 1 {
		// insert decimal point
		pos := rng.Intn(length-1) + 1
		s := string(digits[:pos]) + "." + string(digits[pos:])
		return buildCase(s)
	}
	return buildCase(string(digits))
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.out), strings.TrimSpace(got))
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
	var cases []testCase
	cases = append(cases, buildCase("100"))
	cases = append(cases, buildCase("1"))
	cases = append(cases, buildCase("0.1"))
	cases = append(cases, buildCase("0000123.45000"))
	cases = append(cases, buildCase("9.87654321"))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
