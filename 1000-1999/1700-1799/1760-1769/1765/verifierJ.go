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
	combined := make([]int, 0, 2*n)
	combined = append(combined, a...)
	combined = append(combined, b...)
	sort.Ints(combined)
	m1 := combined[n-1]
	m2 := combined[n]
	cost1 := int64(0)
	cost2 := int64(0)
	for i := 0; i < n; i++ {
		cost1 += absDiff(a[i], m1)
		cost2 += absDiff(a[i], m2)
	}
	for i := 0; i < n; i++ {
		cost1 += absDiff(b[i], m1)
		cost2 += absDiff(b[i], m2)
	}
	if cost2 < cost1 {
		return cost2
	}
	return cost1
}

func absDiff(x, y int) int64 {
	if x >= y {
		return int64(x - y)
	}
	return int64(y - x)
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
