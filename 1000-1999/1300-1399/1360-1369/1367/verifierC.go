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

var testcases = []struct {
	n int
	k int
	s string
}{
	{n: 2, k: 1, s: "01"},
	{n: 6, k: 3, s: "100011"},
	{n: 17, k: 3, s: "11001111100000010"},
	{n: 5, k: 5, s: "10111"},
	{n: 19, k: 3, s: "1101101111111110101"},
	{n: 16, k: 3, s: "1110110100000100"},
	{n: 17, k: 2, s: "10001001100000000"},
	{n: 1, k: 3, s: "1"},
	{n: 5, k: 2, s: "00100"},
	{n: 5, k: 1, s: "01011"},
	{n: 16, k: 1, s: "1101101001001011"},
	{n: 17, k: 3, s: "01111000100000100"},
	{n: 8, k: 2, s: "10100111"},
	{n: 9, k: 5, s: "000110000"},
	{n: 4, k: 1, s: "0000"},
	{n: 7, k: 1, s: "1111001"},
	{n: 14, k: 5, s: "00100110011110"},
	{n: 14, k: 1, s: "01001011100010"},
	{n: 12, k: 3, s: "001001110110"},
	{n: 9, k: 1, s: "001100010"},
	{n: 16, k: 3, s: "1101111111100111"},
	{n: 3, k: 5, s: "001"},
	{n: 6, k: 5, s: "010000"},
	{n: 10, k: 4, s: "0110100101"},
	{n: 17, k: 2, s: "10010110110101011"},
	{n: 13, k: 3, s: "1111000110111"},
	{n: 6, k: 1, s: "010100"},
	{n: 13, k: 5, s: "0101001110100"},
	{n: 12, k: 1, s: "101111100101"},
	{n: 9, k: 2, s: "001111001"},
	{n: 18, k: 4, s: "010110110110111111"},
	{n: 17, k: 1, s: "01001001010001101"},
	{n: 12, k: 1, s: "111011110101"},
	{n: 11, k: 2, s: "00101010001"},
	{n: 18, k: 4, s: "111111011110111100"},
	{n: 11, k: 4, s: "00010101010"},
	{n: 20, k: 4, s: "11011010001101110100"},
	{n: 8, k: 5, s: "01100001"},
	{n: 15, k: 3, s: "010110011110111"},
	{n: 6, k: 4, s: "110010"},
	{n: 4, k: 5, s: "0011"},
	{n: 15, k: 1, s: "101000010110100"},
	{n: 16, k: 2, s: "0000001001100000"},
	{n: 15, k: 2, s: "001110000010111"},
	{n: 5, k: 4, s: "01110"},
	{n: 18, k: 1, s: "111110100000010001"},
	{n: 16, k: 5, s: "0001000000110011"},
	{n: 17, k: 5, s: "01101101001111010"},
	{n: 6, k: 4, s: "111110"},
	{n: 19, k: 4, s: "1011100110010100001"},
	{n: 2, k: 3, s: "11"},
	{n: 15, k: 1, s: "110100001011011"},
	{n: 8, k: 2, s: "10011100"},
	{n: 18, k: 1, s: "100110111111000000"},
	{n: 2, k: 4, s: "11"},
	{n: 8, k: 5, s: "10111010"},
	{n: 9, k: 1, s: "010000100"},
	{n: 20, k: 2, s: "10001001110110101111"},
	{n: 10, k: 5, s: "0001101001"},
	{n: 14, k: 4, s: "00000101000010"},
	{n: 12, k: 1, s: "000000101001"},
	{n: 14, k: 2, s: "01011001110101"},
	{n: 5, k: 1, s: "10011"},
	{n: 15, k: 4, s: "101110101011011"},
	{n: 9, k: 2, s: "111100111"},
	{n: 14, k: 1, s: "01011101111000"},
	{n: 6, k: 4, s: "110111"},
	{n: 19, k: 4, s: "1110100100011110101"},
	{n: 13, k: 5, s: "0100101100011"},
	{n: 16, k: 5, s: "0100101111110010"},
	{n: 12, k: 3, s: "000010110001"},
	{n: 16, k: 3, s: "0001000000100110"},
	{n: 1, k: 3, s: "0"},
	{n: 10, k: 1, s: "1110010010"},
	{n: 9, k: 2, s: "101100010"},
	{n: 6, k: 5, s: "110011"},
	{n: 2, k: 2, s: "01"},
	{n: 10, k: 1, s: "0000101010"},
	{n: 17, k: 5, s: "01010100100001111"},
	{n: 10, k: 3, s: "0000001010"},
	{n: 3, k: 2, s: "010"},
	{n: 6, k: 4, s: "010111"},
	{n: 18, k: 3, s: "000010101000001000"},
	{n: 15, k: 4, s: "101001010011011"},
	{n: 8, k: 1, s: "00001001"},
	{n: 14, k: 4, s: "00100001001100"},
	{n: 9, k: 3, s: "101111111"},
	{n: 13, k: 2, s: "1111100111011"},
	{n: 6, k: 2, s: "100010"},
	{n: 17, k: 2, s: "10011001000101100"},
	{n: 13, k: 2, s: "0011100100101"},
	{n: 10, k: 1, s: "1101001111"},
	{n: 17, k: 3, s: "01010011011000100"},
	{n: 1, k: 1, s: "1"},
	{n: 7, k: 5, s: "1110000"},
	{n: 1, k: 5, s: "0"},
	{n: 5, k: 1, s: "11000"},
	{n: 12, k: 3, s: "010100001010"},
	{n: 8, k: 5, s: "00000000"},
	{n: 19, k: 2, s: "0110000011000110001"},
}

const testcasesCount = 100

func solveCase(n, k int, s string) int {
	right := make([]int, n)
	next := n * 2
	for i := n - 1; i >= 0; i-- {
		if s[i] == '1' {
			next = i
		}
		right[i] = next
	}
	last := -n * 2
	ans := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			last = i
			continue
		}
		if i-last > k && right[i]-i > k {
			ans++
			last = i
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	if len(testcases) != testcasesCount {
		fmt.Fprintf(os.Stderr, "unexpected testcase count: got %d want %d\n", len(testcases), testcasesCount)
		os.Exit(1)
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(testcases))
	for _, tc := range testcases {
		fmt.Fprintf(&input, "%d %d\n%s\n", tc.n, tc.k, tc.s)
	}

	expected := make([]int, len(testcases))
	for i, tc := range testcases {
		expected[i] = solveCase(tc.n, tc.k, tc.s)
	}

	cmd := exec.Command(os.Args[1])
	cmd.Stdin = strings.NewReader(input.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\nstderr: %s\n", err, out)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < len(testcases); i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(outScan.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
