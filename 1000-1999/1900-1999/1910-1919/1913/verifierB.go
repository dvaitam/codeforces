package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// solution logic from 1913B.go
func solveCase(s string) string {
	total0, total1 := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			total0++
		} else {
			total1++
		}
	}
	pref0, pref1, best := 0, 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			pref0++
		} else {
			pref1++
		}
		if pref0 <= total1 && pref1 <= total0 {
			best = i + 1
		}
	}
	return fmt.Sprintf("%d", len(s)-best)
}

const testcasesData = `
10011011110101101
1010011001
1101001100011011
111000100110
0011
11110011
101000010011
00110011011
0000
01011111111
1011000101
0001001
00101110101
1101
00110010100
10111101011
11110111011
110
01110101
01111111001
010010
11110111111
1000000101
110101
001110101010
1011011100101
10000000011
101101011
110000011110
11101111111
110111001111
101110111
0101101
11111011110
1100000011010
10100011110
101101011000011
111011010100111
101000011100100
111101001000
10000111111
111101
001011101101
110010010
100101
010100110011
100000101
100011
100001110010
11101
0111001
01110
10101110110
100101100010
11010010011
1000011111111
101010
111100101
001101
111111100111
010101101011010
01001
1111111010
01110
0101
00001011100010
1110001
111111100001
10010000
01011
1110010
00110111
01110011111011
0110100
00000101
00011
0111111111011
11110010011
10111110000
1010011
0001011000
001001
010
110011
00010100110
111010001001
1001100
01010000
100000000
110010011
0110110
01100
010
11001
000000
111111111010101
01111000
1101
010111110000
1000111
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := fmt.Sprintf("1\n%s\n", line)
		want := solveCase(line)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != want {
			fmt.Printf("test %d failed\ninput: %s\nexpected: %s\ngot: %s\n", idx, line, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
