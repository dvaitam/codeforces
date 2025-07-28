package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(n, q int, prices []int, queries [][2]int) []int64 {
	arr := append([]int(nil), prices...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(arr[i])
	}
	res := make([]int64, q)
	for idx, qu := range queries {
		x := qu[0]
		y := qu[1]
		res[idx] = prefix[x] - prefix[x-y]
	}
	return res
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	prices := make([]int, n)
	for i := range prices {
		prices[i] = rng.Intn(100) + 1
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(x) + 1
		queries[i] = [2]int{x, y}
	}
	ans := solveCase(n, q, prices, queries)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", prices[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", queries[i][0], queries[i][1]))
	}
	var out strings.Builder
	for i, v := range ans {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(fmt.Sprintf("%d", v))
	}
	return testCase{input: sb.String(), expected: out.String()}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\n----got----\n%s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
