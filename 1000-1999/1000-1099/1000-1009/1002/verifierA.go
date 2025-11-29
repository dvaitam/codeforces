package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n    int
	bits string
}

const testcasesRaw = `7 1101111
8 11001001
3 110
2 11
8 10111011
5 10001
1 1
8 11010000
3 110
2 11
8 10110101
8 10110100
3 101
8 10000001
5 10011
8 11101011
4 1001
2 10
6 101100
3 100
2 10
2 10
2 11
2 11
2 10
1 1
4 1001
4 1001
2 11
2 10
2 11
6 110010
2 11
4 1111
3 100
3 101
5 10100
8 11111100
8 10101001
6 111011
7 1000100
4 1111
1 1
7 1100101
3 111
1 1
1 1
8 11110100
3 101
7 1100000
4 1001
5 10011
2 10
5 11010
6 100100
1 1
4 1111
1 1
8 11110011
1 1
3 101
6 111010
1 1
5 10011
7 1010100
5 10011
6 111001
8 11101011
1 1
5 11001
2 10
2 10
4 1010
2 11
5 11101
2 11
4 1101
4 1100
1 1
8 10011100
4 1000
8 11100101
7 1111100
1 1
6 111001
3 111
8 10000000
5 10110
3 111
7 1000101
7 1110000
5 11101
5 11111
1 1
3 100
5 11001
7 1010011
1 1
8 11100100
1 1`

var testcases = mustParseTestcases(testcasesRaw)

// Embedded reference logic from 1002A.go.
func generateOps(n int, bits string) ([]string, error) {
	if len(bits) != n {
		return nil, fmt.Errorf("bitstring length mismatch: got %d want %d", len(bits), n)
	}
	ops := []string{"H 1"}
	for i := 1; i < n; i++ {
		if bits[i] == '1' {
			ops = append(ops, fmt.Sprintf("CNOT 1 %d", i+1))
		}
	}
	return ops, nil
}

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(raw))
	var res []testcase
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		fields := strings.Fields(text)
		if len(fields) != 2 {
			panic(fmt.Sprintf("bad testcase format on line %d", line))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(fmt.Sprintf("bad N on line %d: %v", line, err))
		}
		bits := fields[1]
		if len(bits) != n {
			panic(fmt.Sprintf("bitstring length mismatch on line %d: got %d want %d", line, len(bits), n))
		}
		res = append(res, testcase{n: n, bits: bits})
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(res) == 0 {
		panic("no embedded testcases")
	}
	return res
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) (int, []string, error) {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	if !scanner.Scan() {
		return 0, nil, fmt.Errorf("no output")
	}
	countLine := strings.TrimSpace(scanner.Text())
	k, err := strconv.Atoi(countLine)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse op count: %v", err)
	}
	var ops []string
	for scanner.Scan() {
		ops = append(ops, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return 0, nil, fmt.Errorf("scanner error: %v", err)
	}
	return k, ops, nil
}

func checkCase(bin string, tc testcase) error {
	expectedOps, err := generateOps(tc.n, tc.bits)
	if err != nil {
		return err
	}
	input := fmt.Sprintf("%d\n%s\n", tc.n, tc.bits)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	k, gotOps, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if k != len(gotOps) {
		return fmt.Errorf("reported %d operations but printed %d lines", k, len(gotOps))
	}
	if len(gotOps) != len(expectedOps) {
		return fmt.Errorf("expected %d operations, got %d", len(expectedOps), len(gotOps))
	}
	for i, op := range expectedOps {
		if gotOps[i] != op {
			return fmt.Errorf("line %d: expected %q got %q", i+2, op, gotOps[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range testcases {
		if err := checkCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
