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

const testcasesRaw = `8 10111010
8 10111001
1 1
5 01110
2 10
6 001000
6 111111
7 0000101
1 1
5 00110
3 000
7 1111011
5 11010
1 0
7 1101001
3 011
1 0
4 1010
8 01000101
7 1111101
1 1
1 1
4 1111
1 1
6 011011
3 001
6 000010
6 011011
3 011
1 0
3 010
1 1
3 100
3 101
6 111010
1 0
1 1
2 01
4 1010
8 00100011
8 11100100
3 010
1 1
1 1
7 0010101
6 111100
8 10101000
6 001011
8 10011000
6 101001
1 0
7 1000000
6 011000
5 10010
3 111
3 100
1 1
3 000
6 000000
7 0111001
4 1001
6 011011
4 1011
2 11
6 010101
7 0001011
7 0110011
7 0101101
8 00101110
5 10111
4 1000
6 001110
2 10
3 000
4 0011
5 10110
6 101101
6 111110
3 000
4 0101
2 00
2 11
6 011110
3 000
3 101
7 1001010
7 0110101
2 11
7 1100101
3 100
8 11011001
3 100
1 0
4 0011
2 11
6 000011
5 00111
2 01
7 0110010
8 11100111`

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

func expectedOps(bits string, n int) []string {
	ops := []string{}
	for i := 0; i < n; i++ {
		if bits[i] == '1' {
			ops = append(ops, fmt.Sprintf("CNOT %d %d", i+1, n+1))
		}
	}
	return ops
}

func checkCase(bin string, n int, bits string) error {
	input := fmt.Sprintf("%d\n%s\n", n, bits)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	lines := strings.FieldsFunc(strings.TrimSpace(out), func(r rune) bool { return r == '\n' || r == '\r' })
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("failed to parse op count: %v", err)
	}
	ops := expectedOps(bits, n)
	if k != len(ops) {
		return fmt.Errorf("expected %d operations, got %d", len(ops), k)
	}
	if len(lines)-1 != len(ops) {
		return fmt.Errorf("expected %d lines, got %d", len(ops), len(lines)-1)
	}
	for i, op := range ops {
		if strings.TrimSpace(lines[i+1]) != op {
			return fmt.Errorf("line %d: expected %q got %q", i+2, op, strings.TrimSpace(lines[i+1]))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "bad testcase on line %d\n", idx+1)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad N on line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		bits := fields[1]
		if len(bits) != n {
			fmt.Fprintf(os.Stderr, "bitstring length mismatch on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		if err := checkCase(bin, n, bits); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
