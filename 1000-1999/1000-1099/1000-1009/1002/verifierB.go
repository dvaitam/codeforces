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

const testcasesRaw = `3 101
4 1100
8 10101100
1 1
8 10110001
5 00001
5 00110
2 10
2 10
1 1
7 1011101
2 01
1 1
6 000100
4 0000
5 10000
8 01110001
6 111100
2 10
7 0001000
5 11000
7 0001110
7 1100110
7 1011101
5 10001
2 01
7 0110011
3 110
8 01011001
2 10
7 1010110
2 10
1 1
7 1101110
2 00
5 11000
2 11
5 00011
5 11100
2 01
7 1011110
8 11101101
1 0
1 1
8 10001010
2 01
5 00100
7 1110111
3 101
3 001
6 111100
8 00100001
1 0
4 1011
7 1011111
4 1100
1 0
6 010001
4 1010
6 001111
5 00010
4 1001
5 11111
4 0100
8 11011110
2 00
2 01
4 1011
1 0
4 1001
6 110101
8 11100000
7 0010100
7 1010100
2 11
4 1110
3 011
1 0
7 0000010
4 0011
3 100
6 111000
8 11010101
5 10001
4 1101
6 010100
6 010000
1 0
6 001100
8 10110111
1 1
5 10011
3 000
2 10
3 000
3 010
8 00111001
2 11
3 101
3 111`

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

func expected(bits string) string {
	for _, c := range bits {
		if c == '1' {
			return "1"
		}
	}
	return "0"
}

func checkCase(bin string, n int, bits string) error {
	input := fmt.Sprintf("%d\n%s\n", n, bits)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	exp := expected(bits)
	if out != exp {
		return fmt.Errorf("expected %s got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
