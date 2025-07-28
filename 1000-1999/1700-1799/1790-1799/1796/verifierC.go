package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	l int64
	r int64
}

const mod int64 = 998244353

func solve(l, r int64) (int64, int64) {
	k := int64(0)
	v := l
	for v <= r {
		k++
		v *= 2
	}
	pow2 := int64(1) << (k - 1)
	count1 := r/pow2 - l + 1
	if count1 < 0 {
		count1 = 0
	}
	count2 := int64(0)
	if k >= 2 {
		pow3 := int64(3) * (int64(1) << (k - 2))
		c := r/pow3 - l + 1
		if c > 0 {
			count2 = c * (k - 1)
		}
	}
	total := (count1 + count2) % mod
	return k, total
}

func genTests() []TestCase {
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		l := int64(i*10 + 1)
		r := l + int64((i%3+1)*5)
		tests = append(tests, TestCase{l, r})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.l, tc.r)
	}

	expectedLines := make([]string, len(tests))
	for i, tc := range tests {
		k, cnt := solve(tc.l, tc.r)
		expectedLines[i] = fmt.Sprintf("%d %d", k, cnt)
	}

	out, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		fmt.Print(out)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(expectedLines) {
		fmt.Printf("wrong number of lines: got %d want %d\n", len(lines), len(expectedLines))
		os.Exit(1)
	}
	for i, got := range lines {
		got = strings.TrimSpace(got)
		if got != expectedLines[i] {
			fmt.Printf("test %d failed: input (%d %d) expected %s got %s\n", i+1, tests[i].l, tests[i].r, expectedLines[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
