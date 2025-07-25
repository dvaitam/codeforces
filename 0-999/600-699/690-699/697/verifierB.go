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

func convert(s string) string {
	ePos := strings.IndexByte(s, 'e')
	if ePos == -1 {
		return s
	}
	beforeE := s[:ePos]
	afterE := s[ePos+1:]
	var exp int
	fmt.Sscanf(afterE, "%d", &exp)

	dotPos := strings.IndexByte(beforeE, '.')
	if dotPos == -1 {
		dotPos = len(beforeE)
	}
	intPart := beforeE[:dotPos]
	fracPart := ""
	if dotPos < len(beforeE) {
		fracPart = beforeE[dotPos+1:]
	}
	digits := intPart + fracPart
	decIndex := len(intPart) + exp

	if decIndex >= len(digits) {
		out := digits + strings.Repeat("0", decIndex-len(digits))
		out = strings.TrimLeft(out, "0")
		if out == "" {
			out = "0"
		}
		return out
	}

	iPart := digits[:decIndex]
	fPart := digits[decIndex:]
	iPart = strings.TrimLeft(iPart, "0")
	if iPart == "" {
		iPart = "0"
	}
	fPart = strings.TrimRight(fPart, "0")
	if fPart == "" {
		return iPart
	}
	return iPart + "." + fPart
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) testCase {
	a := rng.Intn(10)
	b := 0
	if a != 0 {
		b = rng.Intn(101)
	}
	var d string
	if rng.Intn(4) == 0 {
		d = "0"
	} else {
		length := rng.Intn(10) + 1
		digits := make([]byte, length)
		for i := range digits {
			digits[i] = byte('0' + rng.Intn(10))
		}
		if digits[length-1] == '0' {
			digits[length-1] = byte('1' + rng.Intn(9))
		}
		d = string(digits)
	}
	s := fmt.Sprintf("%d.%se%d", a, d, b)
	return testCase{input: s + "\n", expected: convert(s)}
}

func runCase(bin string, tc testCase) error {
	out, err := runBinary(bin, tc.input)
	if err != nil {
		return err
	}
	if out != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// some fixed cases
	cases = append(cases,
		testCase{"3.1415e0\n", convert("3.1415e0")},
		testCase{"1.0e2\n", convert("1.0e2")},
		testCase{"0.0e0\n", convert("0.0e0")},
		testCase{"9.99e1\n", convert("9.99e1")},
		testCase{"2.5e0\n", convert("2.5e0")},
	)

	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
