package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Line struct {
	m int64
	b int64
}

type Node struct {
	ln    Line
	left  *Node
	right *Node
}

func eval(ln Line, x int64) int64 {
	return ln.m*x + ln.b
}

func insert(node *Node, l, r int64, ln Line) *Node {
	if node == nil {
		return &Node{ln: ln}
	}
	mid := (l + r) >> 1
	if eval(ln, mid) > eval(node.ln, mid) {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if eval(ln, l) > eval(node.ln, l) {
		node.left = insert(node.left, l, mid, ln)
	} else if eval(ln, r) > eval(node.ln, r) {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *Node, l, r, x int64) int64 {
	if node == nil {
		return math.MinInt64
	}
	res := eval(node.ln, x)
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if x <= mid {
		if v := query(node.left, l, mid, x); v > res {
			res = v
		}
	} else {
		if v := query(node.right, mid+1, r, x); v > res {
			res = v
		}
	}
	return res
}

type LiChao struct {
	root *Node
	l    int64
	r    int64
}

func NewLiChao(l, r int64) *LiChao {
	return &LiChao{l: l, r: r}
}

func (lc *LiChao) Insert(ln Line) {
	lc.root = insert(lc.root, lc.l, lc.r, ln)
}

func (lc *LiChao) Query(x int64) int64 {
	return query(lc.root, lc.l, lc.r, x)
}

func solveF(n int, a []int64) int64 {
	S := make([]int64, n+1)
	T := make([]int64, n+1)
	minS, maxS := int64(0), int64(0)
	for i := 1; i <= n; i++ {
		S[i] = S[i-1] + a[i]
		if S[i] < minS {
			minS = S[i]
		}
		if S[i] > maxS {
			maxS = S[i]
		}
		T[i] = T[i-1] + int64(i)*a[i]
	}
	lc := NewLiChao(minS, maxS)
	lc.Insert(Line{m: 0, b: 0})
	ans := int64(0)
	for r := 1; r <= n; r++ {
		val := lc.Query(S[r])
		cand := T[r] + val
		if cand > ans {
			ans = cand
		}
		line := Line{m: -int64(r), b: int64(r)*S[r] - T[r]}
		lc.Insert(line)
	}
	return ans
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(15) + 1
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			a[i] = int64(rand.Intn(10) - 5)
		}
		var sb strings.Builder
		fmt.Fprintln(&sb, n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, a[i])
		}
		sb.WriteByte('\n')
		expected := solveF(n, a)
		var exp strings.Builder
		fmt.Fprintln(&exp, expected)
		output, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(output)
		want := strings.TrimSpace(exp.String())
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
