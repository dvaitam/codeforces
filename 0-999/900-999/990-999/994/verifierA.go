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
	n int
	m int
	x []int
	y []int
}

func expected(tc testCase) string {
	mark := make([]bool, 10)
	for _, v := range tc.y {
		if v >= 0 && v < 10 {
			mark[v] = true
		}
	}
	var out []int
	for _, v := range tc.x {
		if v >= 0 && v < 10 && mark[v] {
			out = append(out, v)
		}
	}
	var sb strings.Builder
	for i, v := range out {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.x {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.y {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	exp := expected(tc)
	gFields := strings.Fields(strings.TrimSpace(got))
	eFields := strings.Fields(strings.TrimSpace(exp))
	if len(gFields) != len(eFields) {
		return fmt.Errorf("expected %d numbers got %d", len(eFields), len(gFields))
	}
	for i := range gFields {
		if gFields[i] != eFields[i] {
			return fmt.Errorf("expected %v got %v", eFields, gFields)
		}
	}
	return nil
}

func uniqueDigits(rng *rand.Rand, cnt int) []int {
	digits := rng.Perm(10)
	return digits[:cnt]
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	x := uniqueDigits(rng, n)
	y := uniqueDigits(rng, m)
	return testCase{n: n, m: m, x: x, y: y}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, testCase{n: 3, m: 3, x: []int{7, 1, 2}, y: []int{2, 1, 7}})
	cases = append(cases, testCase{n: 4, m: 4, x: []int{1, 0, 2, 3}, y: []int{0, 1, 7, 9}})
	cases = append(cases, testCase{n: 1, m: 1, x: []int{5}, y: []int{3}})
	cases = append(cases, testCase{n: 10, m: 10, x: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, y: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
