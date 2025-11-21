package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
	if actual != expect {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int{9, 1, 6}, []int{1, 2, 3}),
		makeCase([]int{1}, []int{1}),
		makeCase([]int{1000000000}, []int{999999999}),
	}
	for i := 0; i < 200; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		a := randomArray(n)
		b := randomArray(m)
		tests = append(tests, makeCase(a, b))
	}
	return tests
}

func randomArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000) + 1
	}
	return arr
}

func makeCase(a, b []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", len(a), len(b))
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
	expect := outcome(a, b)
	return testCase{
		input:  sb.String(),
		expect: expect,
	}
}

func outcome(a, b []int) string {
	sumA := 0
	for _, v := range a {
		sumA += v
	}
	sumB := 0
	for _, v := range b {
		sumB += v
	}
	if sumA > sumB {
		return "Tsondu"
	} else if sumA < sumB {
		return "Tenzing"
	}
	return "Draw"
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
