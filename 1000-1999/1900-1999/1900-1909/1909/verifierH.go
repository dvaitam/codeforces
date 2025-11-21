package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
	n     int
	perm  []int
}

type operation struct {
	l, r int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("failed to read number of operations: %v", err)
	}
	if k < 0 || k > 1_000_000 {
		return fmt.Errorf("invalid number of operations %d", k)
	}
	ops := make([]operation, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &ops[i].l, &ops[i].r); err != nil {
			return fmt.Errorf("failed to read operation %d: %v", i+1, err)
		}
		if ops[i].l < 1 || ops[i].r > tc.n || ops[i].l >= ops[i].r {
			return fmt.Errorf("operation %d out of range", i+1)
		}
		if (ops[i].r-ops[i].l+1)%2 != 0 {
			return fmt.Errorf("operation %d length not even", i+1)
		}
	}
	arr := append([]int(nil), tc.perm...)
	for _, op := range ops {
		for i := op.l - 1; i < op.r; i += 2 {
			arr[i], arr[i+1] = arr[i+1], arr[i]
		}
	}
	for i := 0; i < tc.n; i++ {
		if arr[i] != i+1 {
			return fmt.Errorf("array not sorted after operations")
		}
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int{2, 5, 4, 1, 3}),
		makeCase([]int{1, 2, 3, 4}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 2
		perm := randPerm(n)
		tests = append(tests, makeCase(perm))
	}
	return tests
}

func makeCase(perm []int) testCase {
	var sb strings.Builder
	n := len(perm)
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{
		input: sb.String(),
		n:     n,
		perm:  perm,
	}
}

func randPerm(n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rand.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
