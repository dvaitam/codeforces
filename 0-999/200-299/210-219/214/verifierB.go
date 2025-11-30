package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var testcases = []string{
	"14 0 1 1 5 2 4 4 9 3 9 0 9 2 6",
	"11 6 8 5 8 7 8 4 0 0 5 7",
	"15 5 6 6 8 2 8 2 3 3 0 2 5 2 2 8",
	"9 5 8 8 2 7 6 8 5 9",
	"6 5 7 2 6 7 8",
	"4 7 4 7 8",
	"9 5 7 7 5 9 8 7 7 3",
	"6 2 9 4 7 4 4",
	"13 8 8 8 8 9 9 6 4 3 7 8 5 9",
	"15 1 5 0 3 1 0 9 0 4 9 3 1 8 2 4",
	"4 3 0 6 0",
	"1 5",
	"6 2 3 0 1 1 1",
	"1 0",
	"12 0 5 4 2 2 2 8 0 6 9 0 3",
	"3 0 0 5",
	"10 1 4 5 7 0 4 7 8 9 0",
	"15 4 6 9 2 7 3 1 5 1 0 7 2 8 9 6",
	"8 8 5 2 5 4 4 9 6",
	"11 0 8 2 0 4 0 2 2 2 1 7",
	"11 3 8 0 3 3 7 1 4 1 9 3",
	"10 9 5 4 6 4 8 0 2 0 6",
	"7 2 1 8 1 3 1 1",
	"1 2",
	"13 3 1 3 0 8 7 7 4 8 6 3 3 6",
	"7 8 0 9 9 0 6 8",
	"10 2 1 7 5 0 8 1 9 5 4",
	"12 5 4 0 6 1 1 4 3 0 7 0 6",
	"11 7 7 3 9 9 1 0 4 0 5 4",
	"15 1 3 7 3 1 9 5 6 7 2 5 6 1 4 1",
	"2 1 9",
	"14 5 6 3 1 0 9 7 0 7 4 5 7 2 5",
	"5 7 8 7 6 7",
	"14 4 6 3 2 7 9 4 8 6 1 9 9 1 1",
	"6 2 8 2 6 1 1",
	"15 0 2 4 6 3 5 7 2 8 4 1 2 8 6 1",
	"6 8 3 8 4 2 2",
	"8 3 6 5 9 2 7 7 0",
	"13 9 6 2 6 8 0 7 4 6 4 6 7 5",
	"9 5 1 3 8 9 3 6 6 0",
	"6 7 8 7 2 1 0",
	"7 3 9 9 6 3 1 6",
	"14 8 3 4 9 9 3 7 9 2 0 9 6 7 4",
	"9 9 2 7 3 1 5 0 7 8",
	"14 1 9 7 5 7 4 8 7 0 1 9 5 2 6",
	"5 2 0 2 7 6",
	"8 4 2 0 4 8 7 0 5",
	"1 8",
	"14 6 9 7 3 4 7 2 7 8 4 1 4 5 4",
	"6 4 6 8 1 8 3",
	"7 9 8 2 8 1 4 0",
	"4 7 8 3 8",
	"5 0 1 1 6 5",
	"4 5 5 1 5",
	"8 5 2 7 7 4 7 2 7",
	"11 3 4 5 2 1 3 7 3 5 2 5",
	"3 2 3 4",
	"13 8 6 6 5 4 9 8 9 5 6 4 8 9",
	"11 1 5 4 6 7 2 4 5 7 7 1",
	"15 2 5 6 2 0 1 5 2 5 1 6 0 8 5 3",
	"14 9 6 8 4 7 2 5 5 3 7 1 2 3 5",
	"5 2 6 5 4 1",
	"6 3 3 3 9 0 5",
	"6 9 0 2 2 1 6",
	"8 4 2 5 8 9 1 5 9",
	"7 3 0 6 7 7 9 5",
	"9 9 9 1 9 8 8 7 6 7",
	"3 6 6 8",
	"8 0 1 7 9 2 1 8 2",
	"2 6 4",
	"8 0 4 1 5 3 2 0 2",
	"7 1 5 7 0 7 3 1",
	"8 2 8 0 2 8 8 0 0",
	"4 8 0 8 5",
	"11 8 3 2 5 7 0 2 8 1 3 1",
	"8 3 0 9 3 6 5 9 6",
	"15 8 8 2 8 1 2 3 2 6 3 4 5 6 2 6",
	"3 6 5 4",
	"13 1 8 1 7 4 4 8 7 4 3 6 2 8",
	"11 1 0 9 8 3 3 3 6 9 0 2",
	"11 0 4 7 8 0 3 2 9 5 0 3",
	"2 2 8",
	"15 2 1 7 4 3 2 5 4 8 9 1 6 6 0 7",
	"5 1 4 0 3 6",
	"6 4 8 6 9 8 3",
	"7 2 2 7 7 5 6 7",
	"10 4 9 3 9 7 7 3 7 9 5",
	"5 1 2 5 9 7",
	"4 9 9 2 4",
	"14 3 8 4 1 0 0 3 5 0 5 8 4 5 7",
	"2 6 7",
	"14 0 4 9 9 2 3 2 2 9 6 1 9 7 4",
	"11 1 7 7 3 2 9 4 3 3 9 5",
	"10 9 6 8 6 3 3 8 0 4 3",
	"3 9 6 6",
	"2 7 6",
	"7 7 6 4 3 3 3 0",
	"9 8 1 9 8 0 0 6 6 6",
	"4 8 4 1 5",
	"9 5 8 7 9 1 7 3 4 0",
}

func referenceSolve(input string) (string, error) {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return "", fmt.Errorf("parse n: %w", err)
	}
	freq := make([]int, 10)
	for i := 0; i < n; i++ {
		var d int
		if _, err := fmt.Fscan(reader, &d); err != nil {
			return "", fmt.Errorf("parse digit %d: %w", i+1, err)
		}
		if d >= 0 && d <= 9 {
			freq[d]++
		}
	}
	if freq[0] == 0 {
		return "-1", nil
	}
	sum := 0
	for d, c := range freq {
		sum += d * c
	}
	mod := sum % 3
	mod1 := make([]int, 0, n)
	mod2 := make([]int, 0, n)
	for d := 1; d <= 9; d++ {
		for i := 0; i < freq[d]; i++ {
			if d%3 == 1 {
				mod1 = append(mod1, d)
			} else if d%3 == 2 {
				mod2 = append(mod2, d)
			}
		}
	}
	remove := func(list []int, cnt int) bool {
		if len(list) < cnt {
			return false
		}
		for i := 0; i < cnt; i++ {
			freq[list[i]]--
		}
		return true
	}
	switch mod {
	case 1:
		if !remove(mod1, 1) {
			if !remove(mod2, 2) {
				return "-1", nil
			}
		}
	case 2:
		if !remove(mod2, 1) {
			if !remove(mod1, 2) {
				return "-1", nil
			}
		}
	}
	if freq[0] == 0 {
		return "-1", nil
	}
	nonZero := 0
	for d := 1; d <= 9; d++ {
		nonZero += freq[d]
	}
	if nonZero == 0 {
		return "0", nil
	}
	var sb strings.Builder
	for d := 9; d >= 0; d-- {
		for i := 0; i < freq[d]; i++ {
			sb.WriteByte(byte('0' + d))
		}
	}
	return sb.String(), nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	idx := 0
	for _, tc := range testcases {
		line := strings.TrimSpace(tc)
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Fprintf(os.Stderr, "case %d: empty line\n", idx)
			os.Exit(1)
		}
		n := 0
		if _, err := fmt.Sscan(parts[0], &n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse n: %v\n", idx, err)
			os.Exit(1)
		}
		if len(parts)-1 != n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d digits got %d\n", idx, n, len(parts)-1)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts[1:], " "))
		expected, err := referenceSolve(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference solve failed: %v\n", idx, err)
			os.Exit(1)
		}
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
