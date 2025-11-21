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
	if strings.TrimSpace(expect) != strings.TrimSpace(actual) {
		return fmt.Errorf("expected %q but got %q", expect, actual)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int{1, 1, 2, 2}),
		makeCase([]int{2, 1, 2, 1}),
		makeCase([]int{3, 3, 1, 1, 2, 2}),
	}
	for i := 0; i < 50; i++ {
		n := rand.Intn(5) + 2
		arr := make([]int, 2*n)
		for j := 0; j < n; j++ {
			arr[2*j] = j + 1
			arr[2*j+1] = j + 1
		}
		rand.Shuffle(2*n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		tests = append(tests, makeCase(arr))
	}
	return tests
}

func makeCase(arr []int) testCase {
	n := len(arr) / 2
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: strings.TrimSpace(sb.String()[len(fmt.Sprintf("1\n%d\n", n)):]),
	}
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
