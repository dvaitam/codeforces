package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var testcases = []struct {
	n int
	a []int
}{
	{n: 3, a: []int{5, 1, 3}},
	{n: 2, a: []int{4, 4}},
	{n: 8, a: []int{4, 2, 1, 4, 1, 4, 4, 5}},
	{n: 1, a: []int{4}},
	{n: 5, a: []int{2, 5, 1, 3, 1}},
	{n: 1, a: []int{1}},
	{n: 9, a: []int{1, 4, 2, 4, 1, 5, 2, 4, 4}},
	{n: 9, a: []int{2, 3, 2, 2, 4, 3, 1, 4, 5}},
	{n: 2, a: []int{2, 3}},
	{n: 2, a: []int{3, 5}},
	{n: 7, a: []int{5, 2, 3, 3, 5, 4, 5}},
	{n: 7, a: []int{5, 1, 4, 2, 4, 4, 2}},
	{n: 6, a: []int{5, 3, 1, 4, 5, 1}},
	{n: 3, a: []int{5, 4, 3}},
	{n: 8, a: []int{1, 4, 1, 3, 5, 5, 5, 4}},
	{n: 3, a: []int{2, 5, 2}},
	{n: 1, a: []int{2}},
	{n: 9, a: []int{5, 2, 4, 5, 3, 5, 3, 4, 3}},
	{n: 9, a: []int{5, 1, 4, 5, 2, 5, 5, 2, 4}},
	{n: 1, a: []int{4}},
	{n: 6, a: []int{5, 5, 2, 5, 4, 4}},
	{n: 6, a: []int{4, 3, 1, 5, 5, 5}},
	{n: 10, a: []int{3, 4, 5, 1, 2, 2, 5, 5, 2, 1}},
	{n: 9, a: []int{3, 1, 1, 1, 1, 4, 1, 3, 2}},
	{n: 5, a: []int{1, 5, 2, 3, 3}},
	{n: 2, a: []int{2, 2}},
	{n: 5, a: []int{5, 2, 3, 3, 4}},
	{n: 6, a: []int{4, 4, 1, 1, 3, 4}},
	{n: 6, a: []int{4, 2, 3, 1, 3, 5}},
	{n: 4, a: []int{5, 4, 1, 2}},
	{n: 1, a: []int{4}},
	{n: 3, a: []int{1, 2, 4}},
	{n: 9, a: []int{4, 5, 2, 5, 4, 2, 5, 1, 4}},
	{n: 10, a: []int{3, 4, 1, 3, 2, 2, 1, 3, 1, 1}},
	{n: 5, a: []int{3, 2, 4, 5, 3}},
	{n: 3, a: []int{1, 5, 1}},
	{n: 10, a: []int{2, 5, 4, 2, 5, 5, 1, 4, 2, 3}},
	{n: 2, a: []int{2, 5}},
	{n: 7, a: []int{5, 2, 4, 1, 4, 3, 5}},
	{n: 8, a: []int{1, 3, 5, 4, 3, 1, 2, 2}},
	{n: 6, a: []int{5, 2, 3, 4, 2, 3}},
	{n: 2, a: []int{4, 5}},
	{n: 6, a: []int{5, 4, 5, 2, 1, 1}},
	{n: 2, a: []int{2, 2}},
	{n: 3, a: []int{5, 2, 3}},
	{n: 6, a: []int{5, 5, 3, 3, 3, 3}},
	{n: 2, a: []int{3, 2}},
	{n: 10, a: []int{4, 2, 5, 5, 1, 3, 1, 4, 1, 4}},
	{n: 3, a: []int{2, 3, 1}},
	{n: 10, a: []int{5, 4, 1, 5, 5, 2, 5, 1, 3, 3}},
	{n: 5, a: []int{5, 5, 1, 4, 3}},
	{n: 2, a: []int{1, 3}},
	{n: 1, a: []int{5}},
	{n: 1, a: []int{1}},
	{n: 7, a: []int{1, 1, 2, 2, 5, 4, 2}},
	{n: 2, a: []int{4, 2}},
	{n: 4, a: []int{2, 1, 4, 4}},
	{n: 9, a: []int{3, 5, 3, 4, 3, 1, 2, 3, 1}},
	{n: 1, a: []int{1}},
	{n: 5, a: []int{5, 3, 4, 4, 3}},
	{n: 7, a: []int{1, 1, 3, 5, 4, 1, 3}},
	{n: 4, a: []int{5, 5, 4, 3}},
	{n: 5, a: []int{2, 5, 2, 3, 2}},
	{n: 4, a: []int{3, 1, 3, 1}},
	{n: 8, a: []int{1, 5, 3, 2, 4, 3, 1, 3}},
	{n: 3, a: []int{3, 5, 3}},
	{n: 4, a: []int{3, 1, 5, 5}},
	{n: 10, a: []int{5, 1, 2, 2, 1, 2, 4, 1, 3, 5}},
	{n: 2, a: []int{1, 1}},
	{n: 1, a: []int{3}},
	{n: 6, a: []int{4, 4, 2, 1, 5, 3}},
	{n: 2, a: []int{5, 2}},
	{n: 3, a: []int{2, 2, 3}},
	{n: 5, a: []int{1, 5, 5, 3, 2}},
	{n: 4, a: []int{2, 5, 1, 3}},
	{n: 10, a: []int{5, 2, 2, 3, 4, 5, 2, 1, 2, 3}},
	{n: 2, a: []int{4, 4}},
	{n: 9, a: []int{3, 5, 4, 5, 4, 1, 4, 3, 2}},
	{n: 5, a: []int{4, 1, 4, 5, 1}},
	{n: 1, a: []int{3}},
	{n: 10, a: []int{2, 5, 2, 2, 3, 3, 4, 5, 4, 2}},
	{n: 10, a: []int{1, 2, 4, 1, 2, 5, 3, 5, 4, 2}},
	{n: 4, a: []int{3, 4, 4, 2}},
	{n: 7, a: []int{3, 5, 5, 3, 2, 1, 1}},
	{n: 9, a: []int{3, 2, 5, 2, 3, 3, 3, 5, 3}},
	{n: 3, a: []int{4, 5, 1}},
	{n: 2, a: []int{5, 5}},
	{n: 10, a: []int{4, 2, 2, 3, 4, 2, 5, 1, 4, 4}},
	{n: 6, a: []int{4, 5, 2, 5, 1, 5}},
	{n: 2, a: []int{3, 1}},
	{n: 5, a: []int{1, 2, 5, 1, 4}},
	{n: 4, a: []int{4, 4, 4, 2}},
	{n: 6, a: []int{4, 2, 5, 4, 2, 1}},
	{n: 7, a: []int{5, 5, 4, 1, 3, 3, 2}},
	{n: 7, a: []int{5, 1, 2, 5, 4, 5, 1}},
	{n: 1, a: []int{5}},
	{n: 4, a: []int{3, 2, 2, 3}},
	{n: 3, a: []int{5, 2, 3}},
	{n: 5, a: []int{5, 3, 4, 2, 5}},
	{n: 6, a: []int{4, 4, 1, 2, 5, 4}},
}

const testcasesCount = 100

func solveCase(a []int) string {
	prefix := 0
	for prefix < len(a) && a[prefix] == 1 {
		prefix++
	}
	if prefix == len(a) {
		if len(a)%2 == 1 {
			return "First"
		}
		return "Second"
	}
	if prefix%2 == 0 {
		return "First"
	}
	return "Second"
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	if len(testcases) != testcasesCount {
		fmt.Fprintf(os.Stderr, "unexpected testcase count: got %d want %d\n", len(testcases), testcasesCount)
		os.Exit(1)
	}
	cand := os.Args[1]

	for i, tc := range testcases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		want := solveCase(tc.a)
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
