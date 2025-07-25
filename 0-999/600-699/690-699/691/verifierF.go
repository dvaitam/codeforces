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

func solveCase(n int, arr []int, queries []int) string {
	var sb strings.Builder
	for qi, p := range queries {
		count := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				if arr[i]*arr[j] >= p {
					count++
				}
			}
		}
		if qi > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", count))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func buildCase(n int, arr []int, qs []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(qs)))
	for i, v := range qs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{in: sb.String(), out: solveCase(n, arr, qs)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(100) + 1
	}
	m := rng.Intn(5) + 1
	qs := make([]int, m)
	for i := range qs {
		qs[i] = rng.Intn(200) + 1
	}
	return buildCase(n, arr, qs)
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
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(tc.out)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase(1, []int{1}, []int{1, 2, 3}))
	cases = append(cases, buildCase(2, []int{1, 2}, []int{1}))
	cases = append(cases, buildCase(3, []int{2, 2, 2}, []int{4}))
	cases = append(cases, buildCase(4, []int{1, 2, 3, 4}, []int{4, 5, 6}))
	cases = append(cases, buildCase(5, []int{5, 4, 3, 2, 1}, []int{10}))
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
