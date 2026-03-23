package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
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
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	exp, _ := strconv.ParseInt(expect, 10, 64)
	if val != exp {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTestCase([]int{0, 0}, []int{1, 1}),
		makeTestCase([]int{1, 3, 5}, []int{2, 4, 6}),
		makeTestCase([]int{10, 20, 30}, []int{10, 20, 30}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(8) + 2
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(50)
			b[j] = rand.Intn(50)
		}
		tests = append(tests, makeTestCase(a, b))
	}
	return tests
}

func makeTestCase(a, b []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(a))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(a, b)),
	}
}

func solveReference(a, b []int) int64 {
	n := len(a)
	A := make([]int64, n)
	B := make([]int64, n)
	for i := 0; i < n; i++ {
		A[i] = int64(a[i])
		B[i] = int64(b[i])
	}
	sort.Slice(A, func(i, j int) bool { return A[i] < A[j] })
	sort.Slice(B, func(i, j int) bool { return B[i] < B[j] })

	var sumW int64
	for i := 0; i < n; i++ {
		diff := A[i] - B[i]
		if diff < 0 {
			sumW -= diff
		} else {
			sumW += diff
		}
	}

	prefB := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefB[i+1] = prefB[i] + B[i]
	}

	var sumC int64
	k := 0
	for i := 0; i < n; i++ {
		for k < n && B[k] <= A[i] {
			k++
		}
		sumC += int64(k)*A[i] - prefB[k]
		sumC += (prefB[n] - prefB[k]) - int64(n-k)*A[i]
	}

	return sumC - int64(n-1)*sumW
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
