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

type Node struct {
	imp, id, prio, sz, maxImp int
	left, right               *Node
}

func sz(n *Node) int {
	if n == nil {
		return 0
	}
	return n.sz
}

func upd(n *Node) {
	if n == nil {
		return
	}
	n.sz = 1 + sz(n.left) + sz(n.right)
	n.maxImp = n.imp
	if n.left != nil && n.left.maxImp > n.maxImp {
		n.maxImp = n.left.maxImp
	}
	if n.right != nil && n.right.maxImp > n.maxImp {
		n.maxImp = n.right.maxImp
	}
}

func merge(l, r *Node) *Node {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.prio > r.prio {
		l.right = merge(l.right, r)
		upd(l)
		return l
	}
	r.left = merge(l, r.left)
	upd(r)
	return r
}

func split(t *Node, k int) (l, r *Node) {
	if t == nil {
		return nil, nil
	}
	if sz(t.left) >= k {
		l0, r0 := split(t.left, k)
		t.left = r0
		upd(t)
		return l0, t
	}
	l0, r0 := split(t.right, k-sz(t.left)-1)
	t.right = l0
	upd(t)
	return t, r0
}

func findLastGE(n *Node, value int) int {
	if n == nil || n.maxImp < value {
		return 0
	}
	if n.right != nil && n.right.maxImp >= value {
		idx := findLastGE(n.right, value)
		if idx > 0 {
			return sz(n.left) + 1 + idx
		}
	}
	if n.imp >= value {
		return sz(n.left) + 1
	}
	return findLastGE(n.left, value)
}

func solveCase(a, c []int) string {
	n := len(a)
	var root *Node
	for i := 0; i < n; i++ {
		imp := a[i]
		cnt := c[i]
		cur := sz(root)
		lo := 1
		if cur+1-cnt > 1 {
			lo = cur + 1 - cnt
		}
		left, right := split(root, lo-1)
		posInRight := findLastGE(right, imp)
		var pos int
		if posInRight > 0 {
			pos = (lo - 1) + posInRight
		} else {
			pos = 0
		}
		p := lo
		if pos+1 > lo {
			p = pos + 1
		}
		root = merge(left, right)
		l2, r2 := split(root, p-1)
		node := &Node{imp: imp, id: i + 1, prio: rand.Int(), sz: 1, maxImp: imp}
		root = merge(merge(l2, node), r2)
	}
	var out strings.Builder
	var dfs func(n *Node)
	dfs = func(n *Node) {
		if n == nil {
			return
		}
		dfs(n.left)
		out.WriteString(fmt.Sprintf("%d ", n.id))
		dfs(n.right)
	}
	dfs(root)
	return strings.TrimSpace(out.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	perm := rng.Perm(n)
	a := make([]int, n)
	for i, v := range perm {
		a[i] = v + 1
	}
	cVals := make([]int, n)
	for i := 0; i < n; i++ {
		cVals[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", a[i], cVals[i]))
	}
	exp := solveCase(a, cVals)
	return sb.String(), exp
}

func runCase(bin, input, expected string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
