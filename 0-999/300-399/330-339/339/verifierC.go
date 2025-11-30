package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `0100100110 1
0111001001 10
1100001011 2
1101110100 3
1101111001 10
0000100111 2
1110010010 5
1010101000 8
0010000101 6
0110011010 2
1110001001 4
1101100110 5
0110001011 9
1000101110 5
1110110001 7
0001000101 9
1001111110 3
1000010010 3
0100000101 4
0101101100 2
0000111111 5
0000000101 3
1001111100 8
0111110111 6
1111001010 9
0111111001 8
0000100001 10
1000010111 6
0100001101 4
1001010001 10
0001110010 6
1111101001 9
1110100001 5
0111010000 3
0111010010 1
0100101110 8
0000011110 7
0000100001 1
1010010111 8
1001111110 4
0110000010 2
1111111100 7
1100010110 9
0001010111 8
1001100101 5
1001011110 9
1101110010 9
1101010010 9
0110010010 5
1101110000 10
0101111110 2
0101111110 2
0001000001 2
1001100000 4
1110100010 4
0110101000 4
0100101001 6
1010011111 4
1100100010 7
1000111101 7
1000111010 7
0010100000 2
1100001000 9
0011110110 1
0100000110 10
1111100110 4
0100110100 9
1101001111 1
0101001100 10
1100001100 5
1010010000 8
1110000100 4
0000001011 10
1100010111 7
1010000001 6
0110110101 5
0110101011 1
1110001011 4
1111011000 6
1001110110 2
0111011111 10
0010001101 2
1010110100 1
0111001110 8
1000010101 2
0100001001 5
0101000101 4
1101101010 5
1111111000 8
1011100010 4
1100111110 4
0010011001 1
1100111100 2
1111011101 5
1011111111 10
0101111001 4
1010010000 4
0011011110 3
1000111110 10
1110010110 7`

var (
	allowed [11]bool
	seq     []int
	m       int
)

// dfs mirrors the search in 339C.go.
func dfs(idx, diff int) bool {
	if idx > m {
		return true
	}
	for v := 10; v > diff; v-- {
		if allowed[v] && seq[idx-1] != v {
			seq[idx] = v
			if dfs(idx+1, v-diff) {
				return true
			}
		}
	}
	return false
}

func solveCase(s string, mVal int) (string, error) {
	if len(s) != 10 {
		return "", fmt.Errorf("bad input string length")
	}
	for i := 0; i < 10; i++ {
		allowed[i+1] = s[i] == '1'
	}
	m = mVal
	seq = make([]int, m+1)
	if !dfs(1, 0) {
		return "NO", nil
	}
	var out strings.Builder
	out.WriteString("YES\n")
	for i := 1; i <= m; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.Itoa(seq[i]))
	}
	return out.String(), nil
}

func parseTestcases() ([][2]string, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([][2]string, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		cases = append(cases, [2]string{fields[0], fields[1]})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		mVal, err := strconv.Atoi(tc[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse m: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect, err := solveCase(tc[0], mVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: solve error: %v\n", idx+1, err)
			os.Exit(1)
		}
		input := tc[0] + "\n" + tc[1] + "\n"
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
