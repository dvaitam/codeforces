package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type node struct {
	val         int
	parent      int
	left, right *node
	idx         int
}

func buildBST(vals []int) (*node, []*node) {
	sort.Ints(vals)
	var nodes []*node
	var build func(l, r int) *node
	build = func(l, r int) *node {
		if l > r {
			return nil
		}
		m := (l + r) / 2
		nd := &node{val: vals[m]}
		nodes = append(nodes, nd)
		nd.left = build(l, m-1)
		nd.right = build(m+1, r)
		if nd.left != nil {
			nd.left.parent = len(nodes) // placeholder
		}
		if nd.right != nil {
			nd.right.parent = len(nodes)
		}
		return nd
	}
	root := build(0, len(vals)-1)
	// assign indices BFS
	queue := []*node{root}
	idx := 1
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		n.idx = idx
		idx++
		if n.left != nil {
			queue = append(queue, n.left)
		}
		if n.right != nil {
			queue = append(queue, n.right)
		}
	}
	// fix parent indices
	queue = []*node{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if n.left != nil {
			n.left.parent = n.idx
			queue = append(queue, n.left)
		}
		if n.right != nil {
			n.right.parent = n.idx
			queue = append(queue, n.right)
		}
	}
	// collect nodes in index order
	ordered := make([]*node, idx-1)
	queue = []*node{root}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		ordered[n.idx-1] = n
		if n.left != nil {
			queue = append(queue, n.left)
		}
		if n.right != nil {
			queue = append(queue, n.right)
		}
	}
	return root, ordered
}

func generateCase(rng *rand.Rand) string {
	h := rng.Intn(3) + 1 // height 1..3 => 3,7,15 nodes
	n := (1 << (h + 1)) - 1
	vals := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; {
		v := rng.Intn(1000) + 1
		if !used[v] {
			used[v] = true
			vals[i] = v
			i++
		}
	}
	_, nodes := buildBST(vals)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, nd := range nodes {
		p := nd.parent
		if p == 0 {
			p = -1
		}
		fmt.Fprintf(&sb, "%d %d\n", p, nd.val)
	}
	q := rng.Intn(5) + 1
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		for {
			v := rng.Intn(1000) + 1001
			if !used[v] {
				fmt.Fprintf(&sb, "%d\n", v)
				break
			}
		}
	}
	return sb.String()
}

func buildOracle() (string, error) {
	exe := "oracleC"
	cmd := exec.Command("go", "build", "-o", exe, "./0-999/0-99/80-89/85/85C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
