package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(a, b []int) string {
	pos, neg := false, false
	for i := range a {
		if b[i] > a[i] && !pos {
			return "NO\n"
		}
		if b[i] < a[i] && !neg {
			return "NO\n"
		}
		if a[i] == 1 {
			pos = true
		}
		if a[i] == -1 {
			neg = true
		}
	}
	return "YES\n"
}

func buildCase(a, b []int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	n := len(a)
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solveCase(a, b)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(3) - 1 // -1..1
		b[i] = rng.Intn(7) - 3 // -3..3
	}
	return buildCase(a, b)
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
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", tc.expected, got)
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
		buildCase([]int{0}, []int{0}),
		buildCase([]int{1, -1}, []int{1, -2}),
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
