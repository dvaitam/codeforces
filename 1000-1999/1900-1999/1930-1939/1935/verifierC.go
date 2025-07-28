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

type Message struct {
	a int64
	b int64
}

type TestCase struct {
	n    int
	L    int64
	msgs []Message
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(6) + 1
		L := int64(rand.Intn(50) + 1)
		msgs := make([]Message, n)
		for j := 0; j < n; j++ {
			a := int64(rand.Intn(10) + 1)
			b := int64(rand.Intn(10) + 1)
			msgs[j] = Message{a: a, b: b}
		}
		tests[i] = TestCase{n: n, L: L, msgs: msgs}
	}
	return tests
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&buf, "%d %d\n", tc.n, tc.L)
		for _, m := range tc.msgs {
			fmt.Fprintf(&buf, "%d %d\n", m.a, m.b)
		}
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
	n := tc.n
	L := tc.L
	msgs := tc.msgs
	maxSize := 0
	for mask := 1; mask < (1 << n); mask++ {
		sumA := int64(0)
		minB := int64(1 << 60)
		maxB := int64(0)
		sz := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sumA += msgs[i].a
				if msgs[i].b < minB {
					minB = msgs[i].b
				}
				if msgs[i].b > maxB {
					maxB = msgs[i].b
				}
				sz++
			}
		}
		cost := sumA
		if sz > 1 {
			cost += maxB - minB
		}
		if cost <= L && sz > maxSize {
			maxSize = sz
		}
	}
	if 0 <= int(L) && maxSize == 0 {
		return "0"
	}
	return fmt.Sprintf("%d", maxSize)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected %s got %s\n", i+1, expectedOutputs[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
