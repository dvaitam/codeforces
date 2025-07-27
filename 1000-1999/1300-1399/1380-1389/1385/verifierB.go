package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n        int
	a        []int
	input    string
	expected string
}

func solveCase(n int, arr []int) string {
	seen := make([]bool, n+1)
	res := make([]int, 0, n)
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			res = append(res, v)
		}
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildCase(p []int) testCase {
	n := len(p)
	// generate random merge of two copies of p
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	arr := make([]int, 0, 2*n)
	i, j := 0, 0
	for i < n || j < n {
		if i == n {
			arr = append(arr, p[j])
			j++
		} else if j == n {
			arr = append(arr, p[i])
			i++
		} else if rng.Intn(2) == 0 {
			arr = append(arr, p[i])
			i++
		} else {
			arr = append(arr, p[j])
			j++
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for idx, v := range arr {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	exp := solveCase(n, arr)
	return testCase{n: n, a: arr, input: sb.String(), expected: exp}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	return buildCase(perm)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.expected)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase([]int{1}),
		buildCase([]int{2, 1}),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
