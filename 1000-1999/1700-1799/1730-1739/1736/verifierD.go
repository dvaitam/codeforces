package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesD.txt. Each line is "n binaryString".
const testcasesRaw = `2 1011
2 0001
3 001100
3 001100
3 110110
2 0110
3 011011
3 110011
1 00
4 11100010
4 10111010
1 01
2 0101
4 00101101
3 101010
4 10110101
1 10
3 110101
2 0000
2 0010
3 001101
2 0111
1 01
4 01010001
3 000110
3 010101
4 01001111
4 01111010
1 10
3 001111
2 0000
3 011000
2 1110
1 00
4 10010011
1 01
4 01101100
1 11
4 10011100
3 000111
3 000100
4 10011100
1 01
4 11110111
2 1000
3 101001
2 1110
3 000100
2 1000
4 01101011
2 1111
3 111010
4 11100011
3 011010
2 0010
1 01
3 000000
3 001010
2 1001
1 01
2 0110
4 00111110
4 01010100
2 0000
1 10
4 11000111
4 00111100
3 001110
2 0100
2 0101
1 00
2 1000
4 00101111
1 01
4 01101100
1 11
3 111101
3 001110
3 000010
4 01010000
1 01
3 000000
3 011100
2 0000
3 100011
4 11110000
3 010101
4 10001000
1 00
3 110110
2 1111
1 10
4 01001100
4 00001001
4 00000101
3 101111
1 01
2 1100
2 0000
3 000010`

type testCase struct {
	n int
	s string
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d expected 2 fields got %d", idx+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		s := fields[1]
		cases = append(cases, testCase{n: n, s: s})
	}
	return cases, nil
}

// Embedded solver logic from 1736D.go (current reference outputs -1 for every case).
func solve(tc testCase) string {
	return "-1"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		want := solve(tc)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		input.WriteString(tc.s)
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
