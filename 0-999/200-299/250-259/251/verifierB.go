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

type testCaseB struct {
	n   int
	k   int
	q   []int
	s   []int
	ans string
}

func identity(n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = i
	}
	return a
}

func compose(p, q []int) []int {
	n := len(p)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = p[q[i]]
	}
	return r
}

func inverse(p []int) []int {
	n := len(p)
	inv := make([]int, n)
	for i, v := range p {
		inv[v] = i
	}
	return inv
}

func equalPerm(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func minSteps(q, qinv [][]int, target []int, k int) int {
	type node struct{ exp, dist int }
	offset := k
	vis := make([]bool, 2*k+1)
	qnodes := []node{{0, 0}}
	vis[offset] = true
	for len(qnodes) > 0 {
		cur := qnodes[0]
		qnodes = qnodes[1:]
		var perm []int
		if cur.exp >= 0 {
			perm = q[cur.exp]
		} else {
			perm = qinv[-cur.exp]
		}
		if equalPerm(perm, target) {
			return cur.dist
		}
		if cur.dist == k {
			continue
		}
		for _, next := range []int{cur.exp + 1, cur.exp - 1} {
			if next < -k || next > k {
				continue
			}
			idx := next + offset
			if !vis[idx] {
				vis[idx] = true
				qnodes = append(qnodes, node{next, cur.dist + 1})
			}
		}
	}
	return -1
}

func computeB(n, k int, qperm, s []int) string {
	qpow := make([][]int, k+1)
	qpow[0] = identity(n)
	for i := 1; i <= k; i++ {
		qpow[i] = compose(qpow[i-1], qperm)
	}
	inv := inverse(qperm)
	invpow := make([][]int, k+1)
	invpow[0] = identity(n)
	for i := 1; i <= k; i++ {
		invpow[i] = compose(invpow[i-1], inv)
	}
	d := minSteps(qpow, invpow, s, k)
	if d == k {
		return "YES"
	}
	return "NO"
}

func genCaseB() testCaseB {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5) + 1
	k := rand.Intn(10) + 1
	q := rand.Perm(n)
	s := rand.Perm(n)
	ans := computeB(n, k, q, s)
	// convert to 1-indexed for input
	q1 := make([]int, n)
	s1 := make([]int, n)
	for i := 0; i < n; i++ {
		q1[i] = q[i] + 1
		s1[i] = s[i] + 1
	}
	return testCaseB{n, k, q1, s1, ans}
}

func buildInputB(cs []testCaseB) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintf(&sb, "%d %d\n", c.n, c.k)
		for i, v := range c.q {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i, v := range c.s {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]testCaseB, 100)
	for i := range cases {
		cases[i] = genCaseB()
	}
	input := buildInputB(cases)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cases) {
		fmt.Printf("expected %d lines got %d\n", len(cases), len(outputs))
		os.Exit(1)
	}
	for i, s := range outputs {
		if s != cases[i].ans {
			fmt.Printf("mismatch on case %d: expected %s got %s\n", i+1, cases[i].ans, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
