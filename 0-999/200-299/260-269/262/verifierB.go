package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(input string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) < 2 {
		return "", fmt.Errorf("not enough fields")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) != 2+n {
		return "", fmt.Errorf("expected %d numbers, got %d", n, len(fields)-2)
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return "", err
		}
		a[i] = v
	}

	for i := 0; i < n && k > 0; i++ {
		if a[i] < 0 {
			a[i] = -a[i]
			k--
		}
	}

	sum := int64(0)
	minAbs := int64(1 << 60)
	zero := false
	for _, v := range a {
		if v == 0 {
			zero = true
		}
		if int64(abs(v)) < minAbs {
			minAbs = int64(abs(v))
		}
		sum += int64(v)
	}
	if k > 0 && k%2 == 1 && !zero {
		sum -= 2 * minAbs
	}
	return fmt.Sprint(sum), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var testcases = []string{
	"3 9 -8 -7 -2",
	"8 7 -10 -7 -4 2 2 5 5 10",
	"7 9 -10 -7 -3 -2 0 4 8",
	"1 0 -10",
	"9 0 -10 -4 -3 2 3 4 5 6 7",
	"4 5 -3 -3 -1 4",
	"1 6 7",
	"2 2 -1 10",
	"2 5 3 6",
	"9 10 -9 -4 -1 -1 2 5 6 8 8",
	"8 3 -8 -5 1 1 2 3 4 7",
	"9 1 -10 -9 -5 -1 1 2 5 5 6",
	"10 9 -10 -5 -5 -4 -3 2 6 7 8 10",
	"9 3 -2 1 1 2 4 6 7 8 9",
	"1 6 6",
	"3 8 -4 3 7",
	"1 7 1",
	"10 8 -10 -4 1 1 3 3 5 6 7 7",
	"10 9 -10 -5 -5 -3 0 4 7 8 9 10",
	"2 8 -9 -2",
	"2 1 -10 4",
	"1 4 -3",
	"5 1 -8 -5 -1 1 9",
	"3 2 -5 -2 6",
	"5 10 -1 0 4 5 5",
	"2 0 -1 2",
	"6 6 -7 -4 -4 -2 -2 6",
	"10 6 -10 -10 -9 -6 -5 -3 2 3 4 6",
	"9 3 -10 -3 2 4 6 6 8 10 10",
	"6 10 -9 -6 -4 -1 3 10",
	"1 4 -8",
	"2 4 -5 -1",
	"7 9 -10 -9 -6 -4 -2 7 8",
	"10 7 -9 -7 -5 -4 -4 1 2 6 8 9",
	"7 9 -7 -4 -1 2 5 5 6",
	"1 5 9",
	"7 4 -10 -6 -5 -4 0 0 8",
	"7 3 -7 -2 1 2 5 7 7",
	"9 3 -9 -8 -8 -6 -5 -5 -4 -2 7",
	"6 9 -7 -2 0 0 1 6",
	"5 3 -6 5 7 8 9",
	"2 5 -9 3",
	"2 6 -6 -6",
	"6 1 -8 2 7 8 8 9",
	"4 9 -8 -2 -1 1",
	"10 8 -10 -10 -9 -8 -7 -7 -2 -1 4 9",
	"7 1 -9 -7 -5 -4 -3 3 8",
	"8 2 -7 -5 -3 -1 2 3 7 7",
	"5 7 -7 -4 0 0 10",
	"1 0 -10",
	"5 9 0 0 2 2 4",
	"2 1 0 9",
	"8 1 -5 -4 -2 -2 1 5 7 9",
	"9 3 -8 -8 -8 -4 -3 -2 -1 1 4",
	"10 10 -9 -5 -3 -1 -1 0 0 0 2 8",
	"4 5 -7 7 8 9",
	"10 1 -10 -8 -8 -8 -3 -3 -3 -2 2 7",
	"1 10 -10",
	"5 5 -7 -6 5 5 6",
	"6 1 -6 -6 -5 -5 0 6",
	"5 1 -6 -4 -1 6 9",
	"3 8 -9 0 9",
	"9 3 -9 -8 -5 -5 -3 -2 -1 3 7",
	"8 6 -10 -2 2 4 4 7 7 7",
	"6 2 -10 -2 3 5 8 10",
	"1 0 1",
	"10 2 -6 -6 -5 -2 -2 2 2 8 8 9",
	"2 3 -10 5",
	"3 8 0 6 10",
	"8 10 -3 -3 -3 0 3 5 5 10",
	"6 8 -9 -3 -2 9 10 10",
	"2 8 1 10",
	"3 8 -4 -1 -1",
	"5 8 -8 -5 1 4 9",
	"2 9 6 8",
	"7 2 -9 -6 -4 -2 3 5 8",
	"7 10 -9 -5 1 2 6 6 7",
	"2 4 -7 10",
	"5 1 -8 -6 -3 4 9",
	"7 6 -6 -5 0 2 4 5 9",
	"4 1 3 3 7 9",
	"2 10 -2 -1",
	"4 6 -10 -4 6 7",
	"8 9 -10 -10 -5 -4 -3 -2 9 10",
	"5 2 -4 -2 -1 7 8",
	"5 10 -5 1 4 5 7",
	"7 1 -10 -7 -4 -4 -1 2 8",
	"2 9 -10 7",
	"5 10 -8 -6 1 6 10",
	"10 4 -10 -7 0 1 1 3 4 4 6 6",
	"5 8 -7 0 2 5 8",
	"7 6 -10 -4 -2 6 7 9 10",
	"4 7 -1 3 6 9",
	"3 7 -4 6 9",
	"6 8 -10 0 2 2 3 8",
	"10 9 -10 -8 -3 -1 3 5 10 10 10 10",
	"3 10 -5 -2 2",
	"2 9 -10 1",
	"5 6 -6 -2 -1 4 7",
	"8 2 -9 -7 -2 3 4 6 6 8",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range testcases {
		input := strings.TrimSpace(tc) + "\n"

		expected, err := solve(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}
