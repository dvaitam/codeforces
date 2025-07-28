package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type testCaseG struct {
	n int
	a []int
}

func genTestsG() []testCaseG {
	rng := rand.New(rand.NewSource(48))
	tests := make([]testCaseG, 100)
	for i := range tests {
		n := rng.Intn(4) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(6)
		}
		tests[i] = testCaseG{n, a}
	}
	return tests
}

func mex(arr []int) int {
	exists := make(map[int]bool)
	for _, v := range arr {
		exists[v] = true
	}
	m := 0
	for exists[m] {
		m++
	}
	return m
}

func encode(arr []int) string {
	b := make([]byte, 0, len(arr)*4)
	for i, v := range arr {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, []byte(fmt.Sprintf("%d", v))...)
	}
	return string(b)
}

func f(arr []int) int {
	sort.Ints(arr)
	best := -1
	type node []int
	q := []node{append([]int{}, arr...)}
	seen := map[string]bool{}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if len(cur) == 1 {
			if cur[0] > best {
				best = cur[0]
			}
			continue
		}
		key := encode(cur)
		if seen[key] {
			continue
		}
		seen[key] = true
		n := len(cur)
		// enumerate subsets by mask
		for mask := 1; mask < (1 << n); mask++ {
			subset := []int{}
			rest := []int{}
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					subset = append(subset, cur[i])
				} else {
					rest = append(rest, cur[i])
				}
			}
			m := mex(subset)
			next := append(rest, m)
			sort.Ints(next)
			q = append(q, next)
		}
	}
	return best
}

func solveG(tc testCaseG) []int {
	res := make([]int, tc.n)
	for i := 1; i <= tc.n; i++ {
		res[i-1] = f(tc.a[:i])
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsG()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	expected := make([][]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveG(tc)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for idx, exp := range expected {
		for j := 0; j < len(exp); j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", idx+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
				os.Exit(1)
			}
			if val != exp[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d element %d: expected %d got %d\n", idx+1, j+1, exp[j], val)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
