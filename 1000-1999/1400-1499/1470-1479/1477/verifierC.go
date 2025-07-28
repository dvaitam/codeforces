package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseC struct {
	n  int
	xs []int64
	ys []int64
}

func generateTests() []testCaseC {
	rand.Seed(42)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(8) + 3 // 3..10
		xs := make([]int64, n)
		ys := make([]int64, n)
		used := map[[2]int64]bool{}
		for j := 0; j < n; j++ {
			for {
				x := rand.Int63n(21) - 10
				y := rand.Int63n(21) - 10
				key := [2]int64{x, y}
				if !used[key] {
					used[key] = true
					xs[j] = x
					ys[j] = y
					break
				}
			}
		}
		tests[i] = testCaseC{n: n, xs: xs, ys: ys}
	}
	return tests
}

func checkPermutation(t testCaseC, perm []int) bool {
	if len(perm) != t.n {
		return false
	}
	seen := make([]bool, t.n)
	for _, v := range perm {
		if v < 1 || v > t.n || seen[v-1] {
			return false
		}
		seen[v-1] = true
	}
	for i := 0; i+2 < t.n; i++ {
		a := perm[i] - 1
		b := perm[i+1] - 1
		c := perm[i+2] - 1
		dx1 := t.xs[a] - t.xs[b]
		dy1 := t.ys[a] - t.ys[b]
		dx2 := t.xs[c] - t.xs[b]
		dy2 := t.ys[c] - t.ys[b]
		if dx1*dx2+dy1*dy2 <= 0 {
			return false
		}
	}
	return true
}

func buildInput(t testCaseC) string {
	var b strings.Builder
	fmt.Fprintln(&b, t.n)
	for i := 0; i < t.n; i++ {
		fmt.Fprintf(&b, "%d %d\n", t.xs[i], t.ys[i])
	}
	return b.String()
}

func parseOutput(out string) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) == 1 && fields[0] == "-1" {
		return nil, fmt.Errorf("got -1")
	}
	perm := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		perm[i] = v
	}
	return perm, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := buildInput(t)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "execution failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		perm, err := parseOutput(out.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if !checkPermutation(t, perm) {
			fmt.Fprintf(os.Stderr, "test %d failed: output is not a valid permutation\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
