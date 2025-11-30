package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []string{
	"9 37 98",
	"5 17 16",
	"32 49 58",
	"31 42 49",
	"14 7 63",
	"2 25 56",
	"39 49 99",
	"1 45 58",
	"18 47 30",
	"38 7 41",
	"2 2 4",
	"42 35 2",
	"25 44 28",
	"28 47 4",
	"34 15 98",
	"29 32 71",
	"15 23 30",
	"44 15 98",
	"30 19 3",
	"27 36 83",
	"7 12 81",
	"47 19 16",
	"48 22 93",
	"46 33 55",
	"33 43 25",
	"20 19 76",
	"32 33 51",
	"38 3 62",
	"16 48 52",
	"27 43 23",
	"24 36 90",
	"50 44 95",
	"24 6 57",
	"43 33 14",
	"50 11 67",
	"26 24 63",
	"47 2 61",
	"3 20 91",
	"40 38 75",
	"26 42 22",
	"11 33 30",
	"1 50 26",
	"35 36 30",
	"26 33 45",
	"37 23 59",
	"18 43 71",
	"39 47 1",
	"25 48 66",
	"9 34 100",
	"36 14 55",
	"4 31 47",
	"37 36 26",
	"33 27 63",
	"23 27 45",
	"1 35 70",
	"40 40 43",
	"30 39 4",
	"15 41 23",
	"36 38 24",
	"6 36 33",
	"3 44 10",
	"6 2 58",
	"1 49 97",
	"18 16 35",
	"8 40 24",
	"23 19 9",
	"11 11 33",
	"34 11 85",
	"18 42 92",
	"19 30 90",
	"21 32 61",
	"8 2 40",
	"25 22 54",
	"13 17 14",
	"17 47 66",
	"14 39 56",
	"2 15 3",
	"26 10 5",
	"47 11 58",
	"46 33 87",
	"28 35 29",
	"41 45 67",
	"29 15 68",
	"42 2 51",
	"44 37 42",
	"43 41 55",
	"4 48 39",
	"9 14 7",
	"20 5 10",
	"20 20 96",
	"11 27 73",
	"17 9 2",
	"36 3 76",
	"14 37 59",
	"11 50 91",
	"40 33 5",
	"25 13 45",
	"7 14 74",
	"44 28 76",
	"13 32 14",
}

// countBlack returns number of initially black squares after removing k layers.
func countBlack(n, m int, k int64) int64 {
	nn := int64(n) - 2*k
	mm := int64(m) - 2*k
	if nn <= 0 || mm <= 0 {
		return 0
	}
	area := nn * mm
	return (area + 1) / 2
}

func referenceSolve(n, m int, x int64) string {
	k := x - 1
	b0 := countBlack(n, m, k)
	b1 := countBlack(n, m, k+1)
	return strconv.FormatInt(b0-b1, 10)
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

func parseLine(line string) (int, int, int64, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("expected 3 numbers, got %d", len(parts))
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse n: %w", err)
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse m: %w", err)
	}
	x, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("parse x: %w", err)
	}
	return n, m, x, nil
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
		n, m, x, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n%d\n", n, m, x)
		expected := referenceSolve(n, m, x)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
