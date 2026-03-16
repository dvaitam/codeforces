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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func expected(n int, s string) int {
	zeros := 0
	for _, c := range s {
		if c == '0' {
			zeros++
		}
	}
	ones := n - zeros
	if zeros > ones {
		return zeros - ones
	}
	return ones - zeros
}

const testcasesARaw = `13 1011111100100
20 10100110111011100010
20 11010000010011011010
20 11011010000110000001
17 10011111010110001
4 0101
14 00000000001010
2 00
6 010001
20 01000111001001011100
2 00
11 10100111111
5 01010
18 100111101110001000
15 110110010101100
16 1110100001110000
17 00001100110011010
17 10010000111011110
7 1100011
11 10100100111
18 010100100111110011
11 10101101100
19 1000000100111101010
9 010100010
4 1110
2 00
4 0111
18 001011111100011100
9 011100000
5 01011
6 011100
17 01011110000111011
10 1110000110
2 11
5 10011
2 11
15 001000010101100
20 11010110001110110001
9 000011001
19 1001001001110111010
4 1011
6 001000
1 0
4 0110
20 00010111000111111000
18 101000100011100010
3 011
18 101110011011011101
10 1000100110
14 10011011010101
20 10010101101011001001
13 1101110000111
4 0001
10 1000111111
19 0000111011110111101
10 0100011100
3 101
2 00
9 100100111
13 0111010100011
13 1110010010000
12 000101111011
16 1101100001010001
8 10101001
11 11100010010
7 1111111
8 01001110
14 11100000010101
18 100010011110101111
9 000001001
1 1
2 01
9 010100011
8 11100011
19 0010000010000110100
17 11110001110101011
13 0001011011010
12 001010000001
3 011
5 01101
10 0111011100
14 10000010000000
14 01001110010110
15 010111110111010
15 101110010110111
7 0101001
1 0
19 0111111110000101011
20 00110000011111011111
7 1000111
2 10
1 1
17 00110001010101000
15 011111110101010
13 1110110110110
2 01
12 000010101101
14 00110000100101
19 0001011111010000100
15 010101001000111
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("bad test format on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("bad n on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		s := fields[1]
		input := fmt.Sprintf("%d\n%s\n", n, s)
		exp := expected(n, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Printf("test %d: cannot parse output: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\ninput:\n%s", idx, exp, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
