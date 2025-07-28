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
	n   int
	m   int
	arr []int
}

func solve(tc testCase) int64 {
	vals := append([]int(nil), tc.arr...)
	sort.Ints(vals)
	n, m := tc.n, tc.m
	k := n * m
	mn := vals[0]
	mn2 := vals[1]
	mx := vals[k-1]
	mx2 := vals[k-2]
	s1 := int64(n-1)*int64(mx-mn2) + int64(n)*int64(m-1)*int64(mx-mn)
	s2 := int64(m-1)*int64(mx-mn2) + int64(m)*int64(n-1)*int64(mx-mn)
	s3 := int64(n-1)*int64(mx2-mn) + int64(n)*int64(m-1)*int64(mx-mn)
	s4 := int64(m-1)*int64(mx2-mn) + int64(m)*int64(n-1)*int64(mx-mn)
	ans := s1
	if s2 > ans {
		ans = s2
	}
	if s3 > ans {
		ans = s3
	}
	if s4 > ans {
		ans = s4
	}
	return ans
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(43))
	tests := []testCase{
		{2, 2, []int{4, 1, 1, 3}},
		{2, 2, []int{0, 0, 0, 0}},
	}
	for len(tests) < 100 {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + 2
		arr := make([]int, n*m)
		for i := range arr {
			arr[i] = rng.Intn(201) - 100
		}
		tests = append(tests, testCase{n, m, arr})
	}
	return tests
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := fmt.Sprintf("%d", solve(tc))
		output, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if output != expected {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
