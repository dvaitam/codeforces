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
	a []int
}

type node struct {
	ch  [2]int
	cnt int
}

var nodes []node
var tot int

func insert(x int) {
	idx := 1
	nodes[idx].cnt++
	for b := 30; b >= 0; b-- {
		bit := (x >> b) & 1
		nxt := nodes[idx].ch[bit]
		if nxt == 0 {
			tot++
			nodes = append(nodes, node{})
			nxt = tot
			nodes[idx].ch[bit] = nxt
		}
		idx = nxt
		nodes[idx].cnt++
	}
}

func f(idx, b int) int {
	if idx == 0 || nodes[idx].cnt == 0 {
		return 0
	}
	if nodes[idx].cnt == 1 || b < 0 {
		return 1
	}
	l := nodes[idx].ch[0]
	r := nodes[idx].ch[1]
	if l == 0 {
		return f(r, b-1)
	}
	if r == 0 {
		return f(l, b-1)
	}
	left := f(l, b-1)
	right := f(r, b-1)
	if left > right {
		return left + 1
	}
	return right + 1
}

func expected(a []int) int {
	nodes = make([]node, 2)
	tot = 1
	for _, x := range a {
		insert(x)
	}
	best := f(1, 30)
	if best < 2 {
		best = 2
	}
	return len(a) - best
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	tests := []testCase{
		{a: []int{0, 1}},
		{a: []int{0, 1, 5, 2}},
		{a: []int{1, 2, 4, 8}},
		{a: []int{7, 3, 5, 1}},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := r.Intn(10) + 2
		m := make(map[int]bool)
		arr := make([]int, 0, n)
		for len(arr) < n {
			v := r.Intn(1024)
			if !m[v] {
				m[v] = true
				arr = append(arr, v)
			}
		}
		tests = append(tests, testCase{a: arr})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(t.a)))
		for j, v := range t.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		want := expected(append([]int(nil), t.a...))
		out, err := runBin(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output\n", i+1)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
