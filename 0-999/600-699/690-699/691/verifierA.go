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
	in  string
	out string
}

func solveCase(n int, arr []int) string {
	zeros := 0
	for _, v := range arr {
		if v == 0 {
			zeros++
		}
	}
	if n == 1 {
		if zeros == 0 {
			return "YES\n"
		}
		return "NO\n"
	}
	if zeros == 1 {
		return "YES\n"
	}
	return "NO\n"
}

func buildCase(arr []int) testCase {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{in: sb.String(), out: solveCase(n, arr)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(2)
	}
	return buildCase(arr)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.out), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// deterministic cases
	cases = append(cases, buildCase([]int{0}))
	cases = append(cases, buildCase([]int{1}))
	cases = append(cases, buildCase([]int{0, 0}))
	cases = append(cases, buildCase([]int{0, 1}))
	cases = append(cases, buildCase([]int{1, 1}))
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
