package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n, k int64) (int64, int64) {
	if k == 0 || k == n {
		return 0, 0
	}
	minGood := int64(1)
	maxGood := n - k
	if maxGood > 2*k {
		maxGood = 2 * k
	}
	return minGood, maxGood
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() [][2]int64 {
	tests := make([][2]int64, 0, 100)
	for i := int64(1); i <= 25; i++ {
		tests = append(tests, [2]int64{i, 0})
		tests = append(tests, [2]int64{i + 25, i + 25})
		tests = append(tests, [2]int64{i + 50, i})
		tests = append(tests, [2]int64{i + 75, i + 50})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		n, k := t[0], t[1]
		expMin, expMax := expected(n, k)
		got, err := run(bin, fmt.Sprintf("%d %d\n", n, k))
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var gMin, gMax int64
		if _, err := fmt.Sscanf(got, "%d %d", &gMin, &gMax); err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", i+1, got)
			os.Exit(1)
		}
		if gMin != expMin || gMax != expMax {
			fmt.Printf("test %d failed: n=%d k=%d expected %d %d got %d %d\n", i+1, n, k, expMin, expMax, gMin, gMax)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
