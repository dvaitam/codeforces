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
	key   int
	left  *Node
	right *Node
	size  int
	sum   int64
}

func update(t *Node) {
	if t == nil {
		return
	}
	ls, lsum := 0, int64(0)
	if t.left != nil {
		ls = t.left.size
		lsum = t.left.sum
	}
	rs, rsum := 0, int64(0)
	if t.right != nil {
		rs = t.right.size
		rsum = t.right.sum
	}
	t.size = ls + rs + 1
	t.sum = int64(t.size) + lsum + rsum
}

func split(t *Node, key int) (*Node, *Node) {
	if t == nil {
		return nil, nil
	}
	if t.key <= key {
		var r *Node
		t.right, r = split(t.right, key)
		update(t)
		return t, r
	}
	var l *Node
	l, t.left = split(t.left, key)
	update(t)
	return l, t
}

func newNode(key int) *Node {
	return &Node{key: key, size: 1, sum: 1}
}

func solveE(a []int) []int64 {
	n := len(a)
	pos := make([]int, n+1)
	for i, v := range a {
		pos[v] = i
	}
	var root *Node
	ans := make([]int64, n)
	for i := 1; i <= n; i++ {
		p := pos[i]
		left, right := split(root, p-1)
		node := newNode(p)
		node.left = left
		node.right = right
		update(node)
		root = node
		ans[i-1] = root.sum
	}
	return ans
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(20) + 1
		perm := rand.Perm(n)
		for i := range perm {
			perm[i]++
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveE(perm)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", t, n, len(lines), input)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			want := fmt.Sprintf("%d", expected[i])
			if strings.TrimSpace(lines[i]) != want {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected %s got %s\ninput:\n%s", t, i+1, want, lines[i], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
