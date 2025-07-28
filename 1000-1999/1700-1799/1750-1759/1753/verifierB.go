package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func solveCase(n, x int, arr []int) string {
	cnt := make([]int64, x+1)
	for _, v := range arr {
		if v <= x {
			cnt[v]++
		}
	}
	q := cnt[1]
	for i := 2; i <= x; i++ {
		val := q + cnt[i]*int64(i)
		if val%int64(i) != 0 {
			return "No"
		}
		q = val / int64(i)
	}
	return "Yes"
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(43))
	var tests []test
	tests = append(tests, test{input: "1 1\n1\n", expected: "Yes"})
	tests = append(tests, test{input: "2 2\n1 2\n", expected: solveCase(2, 2, []int{1, 2})})
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		x := rng.Intn(5) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(x) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		tests = append(tests, test{input: sb.String(), expected: solveCase(n, x, arr)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
