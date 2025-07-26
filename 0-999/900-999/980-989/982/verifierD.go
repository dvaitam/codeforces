package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

type pair struct {
	val int
	idx int
}

func solve(arr []int) string {
	n := len(arr)
	ps := make([]pair, n)
	for i := 0; i < n; i++ {
		ps[i] = pair{arr[i], i}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].val < ps[j].val })
	parent := make([]int, n)
	sz := make([]int, n)
	active := make([]bool, n)
	cnt := map[int]int{}
	segments := 0
	bestSeg := 0
	bestK := 0

	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		ra := find(a)
		rb := find(b)
		if ra == rb {
			return
		}
		la := sz[ra]
		lb := sz[rb]
		cnt[la]--
		if cnt[la] == 0 {
			delete(cnt, la)
		}
		cnt[lb]--
		if cnt[lb] == 0 {
			delete(cnt, lb)
		}
		segments--
		if la < lb {
			ra, rb = rb, ra
			la, lb = lb, la
		}
		parent[rb] = ra
		sz[ra] = la + lb
		cnt[la+lb]++
	}

	for _, p := range ps {
		i := p.idx
		active[i] = true
		parent[i] = i
		sz[i] = 1
		cnt[1]++
		segments++
		if i > 0 && active[i-1] {
			union(i, i-1)
		}
		if i+1 < n && active[i+1] {
			union(i, i+1)
		}
		if len(cnt) == 1 {
			length := 0
			for k := range cnt {
				length = k
			}
			if cnt[length] == segments {
				k := p.val + 1
				if segments > bestSeg || (segments == bestSeg && k < bestK) {
					bestSeg = segments
					bestK = k
				}
			}
		}
	}

	return fmt.Sprintf("%d", bestK)
}

func generateTests() []testCase {
	rand.Seed(4)
	var tests []testCase
	tests = append(tests, testCase{input: "1\n5\n", output: "6"})
	for len(tests) < 120 {
		n := rand.Intn(10) + 1
		arr := rand.Perm(n*3 + 5)[:n]
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v+1)
		}
		b.WriteString("\n")
		tests = append(tests, testCase{input: b.String(), output: solve(arr)})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
