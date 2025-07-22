package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	output string
}

var primes = []int{
	13, 17, 31, 37, 71, 73, 79, 97, 107, 113, 149, 157, 167, 179, 199, 311, 337, 347, 359, 389,
	701, 709, 733, 739, 743, 751, 761, 769, 907, 937, 941, 953, 967, 971, 983, 991, 1009, 1021,
	1031, 1033, 1061, 1069, 1091, 1097, 1103, 1109, 1151, 1153, 1181, 1193, 1201, 1213, 1217,
	1223, 1229, 1231, 1237, 1249, 1259, 1279, 1283, 1301, 1321, 1381, 1399, 1409, 1429, 1439,
	1453, 1471, 1487, 1499, 1511, 1523, 1559, 1583, 1597, 1601, 1619, 1657, 1669, 1723, 1733,
	1741, 1753, 1789, 1811, 1831, 1847, 1867, 1879, 1901, 1913, 1933, 1949, 1979, 3011, 3019,
	3023, 3049,
}

var tests []testCase

func init() {
	for i, p := range primes {
		tests = append(tests, testCase{
			input:  fmt.Sprintf("%d\n", i+1),
			output: fmt.Sprintf("%d", p),
		})
	}
}

func runTest(binary string, t testCase, idx int) error {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(t.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test %d runtime error: %v\n%s", idx, err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(t.output)
	if got != want {
		return fmt.Errorf("test %d failed: expected %q got %q", idx, want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, t := range tests {
		if err := runTest(bin, t, i+1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
