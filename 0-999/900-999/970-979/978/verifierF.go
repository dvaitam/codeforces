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

type testF struct {
	n       int
	k       int
	ratings []int
	pairs   [][2]int
}

func genTestsF() []testF {
	rand.Seed(47)
	tests := make([]testF, 100)
	for i := range tests {
		n := rand.Intn(8) + 2
		k := rand.Intn(n * (n - 1) / 2)
		ratings := make([]int, n)
		for j := range ratings {
			ratings[j] = rand.Intn(100)
		}
		pairs := make([][2]int, k)
		used := make(map[[2]int]bool)
		for j := 0; j < k; j++ {
			for {
				a := rand.Intn(n)
				b := rand.Intn(n)
				if a == b {
					continue
				}
				if a > b {
					a, b = b, a
				}
				p := [2]int{a, b}
				if !used[p] {
					used[p] = true
					pairs[j] = [2]int{a + 1, b + 1}
					break
				}
			}
		}
		tests[i] = testF{n: n, k: k, ratings: ratings, pairs: pairs}
	}
	return tests
}

func solveF(tc testF) []int {
	sorted := make([]int, len(tc.ratings))
	copy(sorted, tc.ratings)
	sort.Ints(sorted)
	ans := make([]int, tc.n)
	for i, r := range tc.ratings {
		ans[i] = sort.SearchInts(sorted, r)
	}
	for _, p := range tc.pairs {
		u := p[0] - 1
		v := p[1] - 1
		if tc.ratings[u] > tc.ratings[v] {
			ans[u]--
		} else if tc.ratings[v] > tc.ratings[u] {
			ans[v]--
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.ratings {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		for _, p := range tc.pairs {
			fmt.Fprintf(&input, "%d %d\n", p[0], p[1])
		}
	}

	expected := make([][]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveF(tc)
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
		tc := tests[idx]
		for j := 0; j < tc.n; j++ {
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
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", idx+1)
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
