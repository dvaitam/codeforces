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
	"time"
)

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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(n, k int, arr []int64) int64 {
	skip := n/2 + 1
	idx := len(arr) - (n / 2) - 1
	var sum int64
	for i := 0; i < k; i++ {
		sum += arr[idx]
		idx -= skip
	}
	return sum
}

func genCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	total := n * k
	arr := make([]int64, total)
	for i := 0; i < total; i++ {
		arr[i] = rng.Int63n(1000)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < total; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(n, k, arr)
}

func parseOutput(out string) (int64, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	v, err := strconv.ParseInt(strings.Fields(out)[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		input string
		exp   int64
	}
	cases := []test{}
	// simple deterministic cases
	cases = append(cases, func() test {
		arr := []int64{1, 2, 3}
		input := "1\n1 3\n1 2 3\n"
		return test{input, expected(1, 3, arr)}
	}())
	cases = append(cases, func() test {
		arr := []int64{5, 5, 5, 5}
		input := "1\n2 2\n5 5 5 5\n"
		return test{input, expected(2, 2, arr)}
	}())
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, test{in, exp})
	}
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %v\noutput:%s\n", i+1, err, out)
			os.Exit(1)
		}
		if got != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, tc.exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
