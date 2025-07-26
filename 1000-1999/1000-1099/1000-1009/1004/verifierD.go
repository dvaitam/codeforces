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

type pos struct{ n, m, x, y int }

func runCandidate(bin, input string) (string, error) {
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

func counts(n, m, x, y int) []int {
	t := n * m
	c := make([]int, t)
	idx := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			c[idx] = abs(i-x) + abs(j-y)
			idx++
		}
	}
	sort.Ints(c)
	return c
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func parseOutput(out string) (pos, error) {
	fields := strings.Fields(out)
	if len(fields) != 4 {
		return pos{}, fmt.Errorf("expected 4 numbers")
	}
	n, err1 := strconv.Atoi(fields[0])
	m, err2 := strconv.Atoi(fields[1])
	x, err3 := strconv.Atoi(fields[2])
	y, err4 := strconv.Atoi(fields[3])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return pos{}, fmt.Errorf("parse error")
	}
	return pos{n, m, x, y}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct{ n, m, x, y int }
	var cases []test
	cases = append(cases, test{1, 1, 1, 1})
	cases = append(cases, test{2, 2, 1, 2})
	for i := 0; i < 98; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		x := rng.Intn(n) + 1
		y := rng.Intn(m) + 1
		cases = append(cases, test{n, m, x, y})
	}

	for idx, tc := range cases {
		arr := counts(tc.n, tc.m, tc.x, tc.y)
		t := tc.n * tc.m
		input := fmt.Sprintf("%d\n", t)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		p, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if p.n <= 0 || p.m <= 0 || p.n*p.m != t || p.x < 1 || p.x > p.n || p.y < 1 || p.y > p.m {
			fmt.Fprintf(os.Stderr, "case %d invalid values\n", idx+1)
			os.Exit(1)
		}
		exp := counts(p.n, p.m, p.x, p.y)
		if !equal(exp, arr) {
			fmt.Fprintf(os.Stderr, "case %d failed: matrix doesn't match counts\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
