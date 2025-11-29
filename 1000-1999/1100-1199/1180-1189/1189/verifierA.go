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

const testcasesRaw = `28 1101111110010010100110111011
28 1000101101000001001101101011
3 110
10 0000110000
30 001100111110101100010010110000
2 00
4 0010
12 000000100010
27 100011100100101110000011010
1 1
22 1111100101010011110111
27 000100011101100101011001111
2 10
18 000111000000001100
16 1001101010010000
22 1110111100110001111010
2 10
5 11101
4 1001
6 011111
4 0111
27 101011011001000000100111101
3 101
19 0101000100111000000
15 110010111111000
11 11001011100
17 00000101100111000
24 101111000011101111110000
10 1001101001
14 01110010000101
7 1100110
22 1011000111011000110000
25 1100110010010011101110100
12 011000100000
22 0011000010111000111111
6 001010
5 01000
10 1100010001
10 1011100110
13 1011101110001
4 0110
25 1100110110101011001010110
22 1011001001111011100001
10 1000011100
18 011111100001110111
24 101111011010001110001010
4 0110
5 10011
23 11011101010001111110010
19 0100001000101111011
16 1101100001010001
22 0101010011111000100100
9 111111001
30 001110111100000010101100010011
11 10101111100
4 0010
7 1010011
4 1010
24 001101110001100100000100
24 001101001111000111010101
26 11000101101101010010100000
8 10011001
14 01101110111001
13 0000010000000
14 01001110010110
15 010111110111010
15 101110010110111
7 0101001
28 0001111111100001010110011000
22 0011111011111010001110
18 100100110001010101
28 0001011111110101010111101101
11 01100011000
23 01010110110011000010010
12 000101111101
18 000010010101010010
27 001111110010000100101100101
2 11
20 11110011010100010010
11 11100110001
12 010010011000
6 001011
2 00
21 110010011101000100001
17 00000010111001010
5 11011
8 11010110
27 111010101100010111011001110
19 1011011001011101101
9 101000100
14 11001000000001
12 110110110110
29 11111110000000101000011101111
26 00000010111001010101101010
9 111100101
28 0100111011001010000000110111
28 0000111001110011011001011110
27 011101111111011010100110011
15 101110100111010
19 0101010001000010000`

func isGood(s string) bool {
	zeros := strings.Count(s, "0")
	ones := len(s) - zeros
	return zeros != ones
}

func expected(n int, s string) (int, []string) {
	zeros := strings.Count(s, "0")
	ones := n - zeros
	if n%2 == 1 || zeros != ones {
		return 1, []string{s}
	}
	return 2, []string{s[:1], s[1:]}
}

func run(bin, input string) (string, error) {
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
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid test case format on line %d\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid n on line %d\n", idx)
			os.Exit(1)
		}
		s := parts[1]
		input := fmt.Sprintf("%d\n%s\n", n, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "test %d: output too short\n", idx)
			os.Exit(1)
		}
		k, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid k\n", idx)
			os.Exit(1)
		}
		if len(fields) != 1+k {
			fmt.Fprintf(os.Stderr, "test %d: expected %d substrings got %d\n", idx, k, len(fields)-1)
			os.Exit(1)
		}
		substrs := fields[1 : 1+k]
		if strings.Join(substrs, "") != s {
			fmt.Fprintf(os.Stderr, "test %d: concatenation mismatch\n", idx)
			os.Exit(1)
		}
		for _, sub := range substrs {
			if !isGood(sub) {
				fmt.Fprintf(os.Stderr, "test %d: substring %s not good\n", idx, sub)
				os.Exit(1)
			}
		}
		expK, _ := expected(n, s)
		if k != expK {
			fmt.Fprintf(os.Stderr, "test %d: expected k=%d got %d\n", idx, expK, k)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
