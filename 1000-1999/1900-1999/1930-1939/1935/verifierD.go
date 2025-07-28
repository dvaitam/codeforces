package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	n   int
	c   int
	arr []int
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(5) + 1
		c := n + rand.Intn(20)
		vals := rand.Perm(c + 1)[:n]
		sortInts(vals)
		tests[i] = TestCase{n: n, c: c, arr: vals}
	}
	return tests
}

func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&buf, "%d %d\n", tc.n, tc.c)
		for i, v := range tc.arr {
			if i > 0 {
				fmt.Fprint(&buf, " ")
			}
			fmt.Fprint(&buf, v)
		}
		fmt.Fprintln(&buf)
	}
	return buf.Bytes()
}

func runBinary(path string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func expected(tc TestCase) string {
	s := make(map[int]bool)
	for _, v := range tc.arr {
		s[v] = true
	}
	var count int64
	for x := 0; x <= tc.c; x++ {
		for y := x; y <= tc.c; y++ {
			if s[x+y] || s[y-x] {
				continue
			}
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := genTests()
	input := buildInput(tests)

	expectedOutputs := make([]string, len(tests))
	for i, tc := range tests {
		expectedOutputs[i] = expected(tc)
	}

	out, err := runBinary(binary, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	got := strings.Fields(strings.TrimSpace(string(out)))
	if len(got) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(expectedOutputs), len(got))
		os.Exit(1)
	}
	for i := range expectedOutputs {
		if got[i] != expectedOutputs[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d expected %s got %s\n", i+1, expectedOutputs[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
