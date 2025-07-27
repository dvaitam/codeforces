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

func solveCase(arr []int) int {
	n := len(arr)
	i := n - 1
	for i > 0 && arr[i-1] >= arr[i] {
		i--
	}
	for i > 0 && arr[i-1] <= arr[i] {
		i--
	}
	return i
}

func buildCase(a []int) testCase {
	n := len(a)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for idx, v := range a {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	res := solveCase(a)
	return testCase{n: n, a: a, input: sb.String(), expected: fmt.Sprintf("%d\n", res)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(50)
	}
	return buildCase(a)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase([]int{1}),
		buildCase([]int{5, 4, 3, 2, 1}),
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
