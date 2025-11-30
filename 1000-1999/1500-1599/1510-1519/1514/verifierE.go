package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `4 0010 1001 0101 1000
4 0100 0011 1000 1010
5 01101 00100 00011 11000 01010
5 00100 10111 00000 10101 10100
6 001001 100000 010011 111000 110101 010100
4 0100 0000 1101 1100
6 011010 000100 010010 101011 010001 111000
5 01101 00001 01010 11001 00100
5 00100 10010 01010 10001 11100
5 00100 10101 00001 11100 10010
4 0000 1010 1001 1100
5 01111 00111 00010 00001 00100
5 01100 00101 00010 11000 10110
5 01100 00101 00001 11101 10000
4 0001 1011 1000 0010
5 01010 00101 10000 01101 10100
5 01100 00000 01010 11000 11110
5 00101 10101 00001 11101 00000
6 000001 101111 100010 101001 100101 001000
4 0100 0001 1101 1000
4 0100 0011 1000 1010
4 0100 0011 1000 1010
5 01001 00001 11011 11000 00010
6 000011 100001 110101 110011 011000 000010
5 00011 10101 10010 01000 00110
6 011000 001111 000000 101000 101100 101110
4 0000 1010 1001 1100
6 011001 000110 010100 100001 101100 011010
4 0110 0011 0001 1000
6 001110 101001 000010 011001 010100 101010
6 000000 100010 110100 110000 101101 111100
4 0100 0001 1100 1010
5 01011 00110 10011 00001 01000
5 00111 10111 00001 00101 00000
5 01011 00011 11001 00101 00000
6 001100 101011 000000 011000 101101 101100
5 01010 00101 10011 01001 10000
6 011110 001000 000111 010000 010100 110110
6 001100 100011 010011 011000 100100 100110
4 0110 0011 0001 1000
6 010101 001110 100111 000001 100101 010000
5 01000 00011 11000 10101 10100
6 000100 101110 100000 001000 101100 111110
6 010011 000001 110011 111011 010000 000010
6 010101 000111 110111 000000 100101 000100
6 011111 000100 010011 001011 010001 010000
4 0100 0011 1001 1000
5 00000 10001 11011 11001 10000
5 01010 00100 10001 01101 11000
4 0011 1011 0001 0000
4 0000 1010 1000 1110
4 0100 0000 1101 1100
5 01010 00110 10001 00101 11000
5 01111 00011 01000 00101 00100
6 011110 001000 000100 010000 011101 111100
6 010100 001111 100111 000001 100101 100000
6 011010 000010 010011 111000 000101 110100
4 0001 1010 1001 0100
5 01011 00111 10011 00001 00000
5 00110 10111 00011 00001 10000
5 01100 00011 01011 10001 10000
4 0000 1001 1100 1010
6 001010 101100 000010 101010 010000 111110
5 01000 00010 11011 10000 11010
6 011010 001111 000111 100000 000100 100110
4 0000 1010 1001 1100
4 0110 0010 0001 1100
5 01100 00011 01000 10101 10100
5 01000 00101 10010 11000 10110
6 011101 001100 000000 001000 111100 011110
6 010010 000011 110001 111001 001101 100000
6 010100 000010 110100 010000 101101 111100
4 0001 1000 1101 0100
5 00001 10001 11011 11000 00010
6 001111 100100 010010 001011 010001 011000
5 01110 00111 00011 00001 10000
6 011100 000001 010010 011001 110101 101000
5 01011 00110 10000 00101 01100
6 011000 001111 000100 100000 101100 101110
4 0001 1000 1100 0110
6 001111 100111 010110 000011 000000 001010
6 011001 001100 000001 101010 111000 010110
6 010010 001110 100111 100000 000101 110100
6 010000 000100 110001 101010 111000 110110
6 000100 101011 100010 011001 100101 101000
5 01101 00100 00000 11100 01110
5 00111 10110 00000 00101 01100
4 0010 1000 0101 1100
5 00001 10011 11000 10100 00110
4 0000 1011 1001 1000
5 01101 00011 01000 10100 00110
5 01111 00000 01010 01000 01110
6 001100 100000 010110 010001 110101 111000
6 011000 001000 000010 111010 110000 111110
5 01010 00000 11001 01100 11010
6 000011 101010 100000 111001 001100 011010
4 0000 1010 1000 1110
4 0101 0001 1101 0000
5 01100 00001 01011 11000 10010
5 01001 00010 11000 10101 01100`

type testCase struct {
	n    int
	grid []string
}

func solveCase(tc testCase) string {
	n := tc.n
	g := make([][]bool, n)
	for i := 0; i < n; i++ {
		g[i] = make([]bool, n)
		for j := 0; j < n && j < len(tc.grid[i]); j++ {
			if tc.grid[i][j] == '1' {
				g[i][j] = true
			}
		}
	}
	reach := make([][]bool, n)
	for i := 0; i < n; i++ {
		reach[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		queue := []int{i}
		reach[i][i] = true
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for u := 0; u < n; u++ {
				if g[v][u] && !reach[i][u] {
					reach[i][u] = true
					queue = append(queue, u)
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("3\n")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if reach[i][j] {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d: expected %d rows got %d", idx+1, n, len(fields)-1)
		}
		grid := make([]string, n)
		copy(grid, fields[1:])
		res = append(res, testCase{n: n, grid: grid})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
