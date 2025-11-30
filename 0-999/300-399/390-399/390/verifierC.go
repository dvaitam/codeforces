package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcase corpus from testcasesC.txt.
const rawTestcasesData = `
14 4 1 11100001110011 2 13
6 1 3 101100 2 3 4 6 3 5
9 4 2 011101111 2 9 2 9
6 4 3 010100 2 5 1 4 2 5
11 4 4 01110100111 7 10 1 4 1 8 2 9
10 2 1 0100001000 5 10
8 1 4 00100001 4 5 7 7 7 7 5 5
5 5 3 11100 1 5 1 5 1 5
7 5 2 1010100 2 6 3 7
12 4 2 110111010101 4 11 8 11
5 5 4 10101 1 5 1 5 1 5 1 5
14 2 1 01000101011011 3 4
11 1 3 01010101001 6 8 4 6 1 1
5 2 2 10011 2 5 1 4
7 3 3 1111100 1 3 4 6 1 6
10 1 1 1001100011 9 10
13 1 1 1001110000101 2 3
14 4 3 01110100011011 4 11 1 12 7 10
14 2 2 00100111011100 9 12 1 6
7 4 1 1101010 1 4
13 1 2 1001000010110 12 13 9 10
8 2 4 10100111 1 2 1 6 1 4 6 7
9 2 3 101100110 2 5 4 9 3 8
5 1 1 10111 2 2
13 2 3 0111011110110 1 6 12 13 7 10
5 2 1 10001 3 4
13 4 1 0111111001000 5 8
7 2 3 0001000 1 6 1 6 2 7
5 4 2 01011 2 5 2 5
6 2 1 000011 1 6
11 3 4 10000000101 3 11 9 11 1 6 1 9
10 3 3 1100101111 2 7 3 8 2 10
10 5 4 1110010100 1 10 1 10 1 10 1 10
8 3 4 00001000 1 6 2 7 2 7 6 8
15 1 1 010111101000011 4 6
11 2 3 00101000100 1 6 5 8 4 7
5 3 2 11111 3 5 2 4
7 1 3 1110100 1 3 6 6 3 4
15 1 4 110100000011101 3 5 12 12 11 12 12 14
8 1 2 10010111 1 3 4 6
8 5 1 10001010 1 4
8 2 4 11100010 4 4 2 4 2 7 6 8
12 4 2 111111001111 2 10 1 6
11 4 2 10010000101 3 8 6 9
9 3 3 101100111 2 7 2 9 2 8
5 1 3 11111 1 2 1 3 1 4
10 4 3 1111101110 1 6 3 5 5 9
13 3 3 1010101111111 1 11 1 4 2 6
11 3 4 00110100011 1 6 2 9 2 7 3 5
13 5 3 0101111101011 1 8 8 9 11 12
9 4 3 011111011 4 6 1 3 2 6
10 1 3 0011111010 1 4 1 7 1 8
7 5 1 0101010 6 7
14 4 2 01010010000000 3 8 10 13
7 1 3 0100011 4 4 6 7 6 6
11 2 1 00011100010 10 11
10 1 3 1011001110 2 7 2 8 7 10
11 3 1 01000000100 4 10
14 1 1 00010100100100 12 12
6 3 4 101010 1 6 1 3 4 6 2 5
12 3 3 011110010011 7 10 11 12 4 9
11 1 2 01000111000 5 9 2 7
12 4 4 110111001100 1 5 5 11 8 11 1 7
13 4 4 1001000101110 5 7 1 4 5 7 2 10
14 5 1 01011001111010 10 14
5 5 1 01101 1 3
15 1 3 000111010011011 5 6 2 3 11 14
13 1 4 1100010110111 4 4 11 12 2 10 12 12
11 1 4 00100100011 7 8 4 7 5 7 6 10
10 5 1 0101101001 9 10
8 1 2 01001000 1 7 3 4
13 3 3 0001100100010 6 11 10 11 11 13
5 2 1 11000 3 3
12 1 1 000011011100 3 10
9 3 3 011011010 3 8 1 8 1 8
12 1 2 000010010111 9 10 12 12
15 4 2 001100010100000 3 10 2 8
13 5 2 0110111011101 5 13 10 13
15 2 1 101010111011101 3 12
12 3 4 010110110000 1 5 1 4 9 12 5 11
14 2 3 11011011001100 2 13 2 9 6 7
12 3 3 000001100111 1 6 7 10 5 8
6 2 1 101101 3 3
10 4 4 0011000001 1 10 1 10 1 10 1 10
9 5 1 011010111 2 5
6 1 4 001100 1 6 1 5 4 5 1 5
10 4 4 1110001000 2 4 3 8 7 7 2 6
7 3 3 1010000 4 4 2 2 1 7
14 4 1 00110011001001 3 5
6 2 2 000000 2 4 1 1
6 3 1 011101 3 6
11 3 4 10100110011 7 10 6 7 4 10 5 6
12 1 1 011001110100 4 4
9 2 2 011010110 2 7 7 9
14 4 1 10001010010110 1 3
10 5 1 1111111111 10 10
14 4 3 00000001101101 9 14 1 6 5 6
14 4 2 11101100110011 1 6 6 11
14 3 3 11010110110101 4 10 1 10 6 6
7 1 4 1000000 2 6 3 4 1 6 5 5
10 4 3 0001010000 4 8 1 4 7 8
6 3 3 100001 4 5 2 5 4 4
12 3 2 101001111001 11 12 3 11
8 2 1 10011001 3 5
`

func loadTestcases() []string {
	lines := strings.Split(strings.TrimSpace(rawTestcasesData), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func solveInstance(n, k int, s string, queries [][2]int) []int {
	// Prefix sums per residue mod k for positions that are '1'.
	ps := make([][]int, k)
	for i := 0; i < k; i++ {
		ps[i] = make([]int, n+1)
	}
	// Total ones prefix sum.
	tot := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for r := 0; r < k; r++ {
			ps[r][i] = ps[r][i-1]
		}
		tot[i] = tot[i-1]
		if s[i-1] == '1' {
			tot[i]++
			ps[(i-1)%k][i]++
		}
	}

	ans := make([]int, len(queries))
	for idx, q := range queries {
		l, r := q[0], q[1]
		total1 := tot[r] - tot[l-1]
		targetResidue := (l - 1) % k
		good1 := ps[targetResidue][r] - ps[targetResidue][l-1]
		numTargets := (r - l + 1) / k
		ans[idx] = total1 + numTargets - 2*good1
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcase(line string) (string, string, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return "", "", fmt.Errorf("invalid testcase line")
	}
	n, err1 := strconv.Atoi(fields[0])
	k, err2 := strconv.Atoi(fields[1])
	w, err3 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return "", "", fmt.Errorf("failed to parse header")
	}
	str := fields[3]
	if len(str) != n {
		return "", "", fmt.Errorf("string length mismatch")
	}
	expectedFields := 4 + 2*w
	if len(fields) != expectedFields {
		return "", "", fmt.Errorf("expected %d fields, got %d", expectedFields, len(fields))
	}
	queries := make([][2]int, w)
	for i := 0; i < w; i++ {
		l, _ := strconv.Atoi(fields[4+2*i])
		r, _ := strconv.Atoi(fields[5+2*i])
		queries[i] = [2]int{l, r}
	}
	expVals := solveInstance(n, k, str, queries)
	var inputBuilder strings.Builder
	fmt.Fprintf(&inputBuilder, "%d %d %d\n", n, k, w)
	inputBuilder.WriteString(str)
	inputBuilder.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&inputBuilder, "%d %d\n", q[0], q[1])
	}
	var expBuilder strings.Builder
	for i, v := range expVals {
		if i > 0 {
			expBuilder.WriteByte('\n')
		}
		fmt.Fprintf(&expBuilder, "%d", v)
	}
	return inputBuilder.String(), expBuilder.String(), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := loadTestcases()
	for idx, line := range testcases {
		input, expect, err := parseTestcase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
