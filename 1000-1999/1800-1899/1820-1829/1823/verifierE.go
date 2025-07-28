package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n     int
	l     int
	r     int
	edges [][2]int
}

func mex(set map[int]bool) int {
	for i := 0; ; i++ {
		if !set[i] {
			return i
		}
	}
}

func expected(tc Test) string {
	n := tc.n
	l, r := tc.l, tc.r
	adj := make([][]int, n)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n)
	cycles := []int{}
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		cur := i
		prev := -1
		length := 0
		for {
			visited[cur] = true
			length++
			next := adj[cur][0]
			if next == prev {
				if len(adj[cur]) > 1 {
					next = adj[cur][1]
				}
			}
			prev, cur = cur, next
			if cur == i {
				break
			}
		}
		cycles = append(cycles, length)
	}
	maxLen := 0
	for _, v := range cycles {
		if v > maxLen {
			maxLen = v
		}
	}
	if maxLen < l {
		return "Bob"
	}
	grundy := make([]int, maxLen+1)
	for length := 1; length <= maxLen; length++ {
		reachable := make(map[int]bool)
		for k := l; k <= r && k <= length; k++ {
			for start := 0; start <= length-k; start++ {
				left := start
				right := length - k - start
				val := grundy[left] ^ grundy[right]
				reachable[val] = true
			}
		}
		grundy[length] = mex(reachable)
	}
	total := 0
	for _, c := range cycles {
		reachable := make(map[int]bool)
		if c >= l && c <= r {
			reachable[0] = true
		}
		for k := l; k <= r && k < c; k++ {
			reachable[grundy[c-k]] = true
		}
		g := mex(reachable)
		total ^= g
	}
	if total != 0 {
		return "Alice"
	}
	return "Bob"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(5)
	const cases = 100
	tests := make([]Test, cases)
	for i := range tests {
		cycles := rand.Intn(3) + 1
		total := 0
		lengths := make([]int, cycles)
		for j := 0; j < cycles; j++ {
			lengths[j] = rand.Intn(4) + 1
			total += lengths[j]
		}
		l := rand.Intn(3) + 1
		r := l + rand.Intn(3)
		edges := make([][2]int, 0, total)
		idx := 0
		for _, lenC := range lengths {
			for j := 0; j < lenC; j++ {
				u := idx + j
				v := idx + (j+1)%lenC
				edges = append(edges, [2]int{u, v})
			}
			idx += lenC
		}
		tests[i] = Test{n: total, l: l, r: r, edges: edges}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.l, tc.r)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0]+1, e[1]+1)
		}
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running binary:", err)
		fmt.Print(out.String())
		return
	}

	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	for idx, tc := range tests {
		var ans string
		if _, err := fmt.Fscan(reader, &ans); err != nil {
			fmt.Printf("test %d: failed to read output\n", idx+1)
			return
		}
		exp := expected(tc)
		if ans != exp {
			fmt.Printf("test %d: expected %s got %s\n", idx+1, exp, ans)
			return
		}
	}
	fmt.Printf("verified %d test cases\n", len(tests))
}
