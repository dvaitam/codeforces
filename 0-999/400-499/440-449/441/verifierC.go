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
	n, m, k int
}

func (tc testCase) Input() string {
	return fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k)
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	k := rng.Intn(n*m) + 1
	return testCase{n, m, k}
}

func expected(tc testCase) []string {
	total := tc.n * tc.m
	base := total / tc.k
	extra := total % tc.k
	coords := make([][2]int, total)
	idx := 0
	for i := 0; i < tc.n; i++ {
		if i%2 == 0 {
			for j := 0; j < tc.m; j++ {
				coords[idx] = [2]int{i + 1, j + 1}
				idx++
			}
		} else {
			for j := tc.m - 1; j >= 0; j-- {
				coords[idx] = [2]int{i + 1, j + 1}
				idx++
			}
		}
	}
	var lines []string
	pos := 0
	for t := 1; t <= tc.k; t++ {
		size := base
		if t == tc.k {
			size += extra
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d", size))
		for i := 0; i < size; i++ {
			c := coords[pos]
			sb.WriteString(fmt.Sprintf(" %d %d", c[0], c[1]))
			pos++
		}
		lines = append(lines, sb.String())
	}
	return lines
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	input := tc.Input()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	expectLines := expected(tc)
	if len(gotLines) != len(expectLines) {
		return fmt.Errorf("expected %d lines got %d", len(expectLines), len(gotLines))
	}
	for i := range expectLines {
		if strings.TrimSpace(gotLines[i]) != expectLines[i] {
			return fmt.Errorf("line %d mismatch expected '%s' got '%s'", i+1, expectLines[i], strings.TrimSpace(gotLines[i]))
		}
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

	cases := []testCase{{1, 1, 1}, {2, 3, 2}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
