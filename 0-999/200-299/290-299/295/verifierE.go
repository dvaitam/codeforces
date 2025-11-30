package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []string{
	"1 -2 3 2 -2 -1 2 6 9 1 1 5",
	"2 3 -2 4 1 1 1 1 2 -5 1 2 5 2 1 7",
	"4 6 4 7 5 2 1 3 1 1 2 -4",
	"2 9 4 2 2 0 4 1 2 -1",
	"4 0 6 -3 -9 3 2 2 7 1 3 -3 1 2 -4",
	"3 -4 -9 -10 4 2 -8 8 1 1 5 2 -6 8 2 -8 -1",
	"3 -6 4 1 1 2 -7 -7",
	"1 -1 3 2 -10 -5 2 2 2 1 1 1",
	"3 5 3 10 4 1 2 -2 2 -4 10 2 -8 -6 2 0 1",
	"3 4 0 -3 2 2 -6 0 1 2 -4",
	"4 0 -5 -8 10 2 2 -7 1 1 4 -1",
	"3 -10 -4 -1 2 1 3 2 1 1 3",
	"3 -7 -9 2 1 2 4 5",
	"2 -6 -9 1 1 1 1",
	"1 -1 1 1 1 -1",
	"4 -5 6 4 -2 3 1 4 -1 1 2 -4 1 3 5",
	"4 -8 0 6 3 3 2 -9 -6 2 -5 5 1 2 -4",
	"2 4 3 1 2 -7 -3",
	"3 -6 -4 -3 3 2 0 8 2 -2 6 2 -10 9",
	"2 -2 -4 1 2 -5 4",
	"1 -9 4 2 0 3 2 -3 0 2 -6 6 1 1 -3",
	"1 -8 1 1 1 -5",
	"3 -1 -4 -8 1 1 3 -3",
	"3 2 6 1 4 1 2 4 1 3 -5 2 -10 3 1 1 -4",
	"2 0 -2 3 2 0 6 2 0 5 2 -1 6",
	"4 7 0 2 9 4 2 -10 -4 2 -1 2 2 -6 -5 2 -5 -4",
	"2 -8 -2 2 1 2 -4 2 -3 3",
	"2 -7 -9 1 1 1 4",
	"2 -6 8 2 2 -7 1 1 2 2",
	"2 0 8 2 1 2 4 2 -7 -6",
	"4 -7 -4 -6 -5 3 2 -5 0 2 0 7 1 1 -5",
	"1 -4 3 2 -3 0 2 -7 8 2 -8 5",
	"4 -1 -3 2 -9 3 2 -8 -5 2 -10 4 1 2 3",
	"3 4 -8 2 3 2 -10 -1 2 -10 4 1 1 -1",
	"3 -1 9 -6 4 2 -10 9 2 -6 3 1 2 -3 1 1 1",
	"4 7 -6 3 0 2 2 -5 8 1 1 -2",
	"1 -4 4 1 1 1 2 -2 1 2 4 6 1 1 2",
	"4 -3 5 -10 6 4 2 -4 8 2 -5 -1 1 3 -2 1 3 1",
	"2 -1 -2 2 2 -10 -8 1 2 -4",
	"1 8 2 1 1 2 1 1 -3",
	"1 4 3 1 1 2 1 1 5 1 1 -3",
	"4 -4 6 -1 -8 2 2 -10 7 2 -1 10",
	"2 8 -1 1 2 -8 1",
	"3 -2 1 -7 3 1 2 -5 1 3 4 2 -10 3",
	"2 3 8 4 2 -10 -5 2 -8 7 1 1 -4 1 2 1",
	"2 -3 7 1 1 1 1",
	"4 -1 -5 2 -10 3 2 -4 6 2 -3 9 1 4 2",
	"1 -4 3 2 -5 8 1 1 -5 1 1 -1",
	"3 0 -1 -6 2 2 -9 -2 2 -7 8",
	"4 -5 2 -3 3 4 1 4 4 2 -9 10 1 2 3 2 -3 5",
	"1 7 1 1 1 -1",
	"1 0 1 1 1 3",
	"4 -2 -10 1 -6 1 1 4 5",
	"3 3 -9 -3 3 1 2 5 1 3 2 2 -9 -1",
	"3 -7 1 -1 2 2 -3 -3 1 1 5",
	"4 -7 7 -2 -5 3 2 -10 7 1 4 5 1 3 -2",
	"3 -8 -7 6 4 2 -6 -5 1 2 -3 2 -3 -1 1 1 2",
	"2 1 3 3 2 -7 -6 2 1 10 1 1 1",
	"4 0 9 -1 5 2 2 -9 3 1 4 -3",
	"1 -7 4 1 1 5 1 1 5 2 -10 7 1 1 1",
	"3 -3 10 -5 1 1 3 4",
	"2 -1 1 4 1 2 -3 1 2 -3 1 2 -3 1 2 4",
	"2 -7 7 1 2 6 9",
	"4 -1 0 -8 -4 1 2 -2 6",
	"3 -3 8 0 3 1 2 -3 1 3 -4 2 -2 0",
	"3 4 -10 0 3 1 2 1 1 1 1 2 2 9",
	"1 7 1 2 -10 3",
	"2 -8 -6 1 2 -4 10",
	"3 -3 8 -8 2 2 0 10 2 -10 -8",
	"4 -2 2 -4 -5 4 1 2 -5 1 1 -4 2 -9 -3 1 4 4",
	"3 10 7 0 4 1 1 -2 1 1 -3 1 1 4 1 3 4",
	"4 4 10 2 -10 3 1 2 4 2 -1 7 2 6 9",
	"1 4 2 2 -2 7 2 -7 0",
	"3 6 -2 4 4 1 3 -1 2 -8 3 2 -10 -7 2 -7 8",
	"2 8 -10 1 1 2 2",
	"1 6 1 2 -3 -2",
	"3 -2 -4 2 4 2 -1 0 1 3 4 2 -7 0 1 1 1",
	"4 -2 1 10 -3 2 2 -2 4 2 -8 -7",
	"2 7 -9 3 2 0 1 1 2 5 2 -5 -1",
	"4 -9 7 6 -1 3 2 9 10 2 -8 1 1 4 4",
	"2 1 9 3 2 -10 5 2 -5 -3 2 -7 -3",
	"1 2 4 2 9 10 1 1 -5 1 1 -1 1 1 2",
	"1 2 1 2 -7 -6",
	"1 -10 3 2 6 8 2 -8 0 2 0 8",
	"3 0 -3 -9 2 2 -10 6 2 -9 0",
	"2 -7 7 2 1 2 -3 1 2 3",
	"4 3 -5 4 -2 2 1 3 5 2 0 10",
	"3 -3 2 -10 1 1 1 4",
	"3 8 -5 -10 1 1 3 -3",
	"1 -6 2 1 1 1 1 1 -1",
	"3 8 -10 -3 1 2 0 1",
	"1 8 3 2 -8 3 1 1 4 2 -6 -4",
	"2 9 -6 1 2 8 8",
	"1 9 2 2 -2 4 2 -5 4",
	"4 3 1 -7 0 2 2 -5 0 2 -8 1",
	"3 -2 1 -6 4 1 1 -3 2 3 9 1 1 -3 1 2 1",
	"1 -9 4 2 -4 9 2 -10 9 1 1 4 2 5 9",
	"4 2 10 -9 -1 2 1 3 -4 2 -5 8",
	"4 -3 -6 -2 -10 4 2 1 5 1 3 4 1 4 -3 2 -4 1",
	"1 -1 4 1 1 -2 2 -6 -2 1 1 1 2 1 5",
}

// Treap node storing key, subtree size, sum of keys, and sum of pairwise distances
type Node struct {
	key   int64
	pri   int
	left  *Node
	right *Node
	sz    int
	sumX  int64
	sumP  int64
}

func update(n *Node) {
	if n == nil {
		return
	}
	n.sz = 1
	n.sumX = n.key
	n.sumP = 0
	l, r := n.left, n.right
	if l != nil {
		n.sz += l.sz
		n.sumX += l.sumX
		n.sumP += l.sumP
	}
	if r != nil {
		n.sz += r.sz
		n.sumX += r.sumX
		n.sumP += r.sumP
	}
	if l != nil {
		n.sumP += int64(l.sz)*n.key - l.sumX
	}
	if r != nil {
		n.sumP += r.sumX - int64(r.sz)*n.key
	}
	if l != nil && r != nil {
		n.sumP += int64(l.sz)*r.sumX - int64(r.sz)*l.sumX
	}
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pri > b.pri {
		a.right = merge(a.right, b)
		update(a)
		return a
	}
	b.left = merge(a, b.left)
	update(b)
	return b
}

func split(n *Node, key int64) (l, r *Node) {
	if n == nil {
		return nil, nil
	}
	if n.key <= key {
		var rr *Node
		n.right, rr = split(n.right, key)
		update(n)
		return n, rr
	}
	var ll *Node
	ll, n.left = split(n.left, key)
	update(n)
	return ll, n
}

func insert(root *Node, key int64, rnd *randState) *Node {
	newNode := &Node{key: key, pri: rnd.next(), sz: 1, sumX: key, sumP: 0}
	var a, b *Node
	a, b = split(root, key)
	return merge(merge(a, newNode), b)
}

func erase(root *Node, key int64) *Node {
	var a, b, c *Node
	a, b = split(root, key-1)
	c, b = split(b, key)
	_ = c
	return merge(a, b)
}

func rangeSumP(root *Node, lkey, rkey int64) (ans int64, rootOut *Node) {
	var a, b, c, d *Node
	a, b = split(root, lkey-1)
	c, d = split(b, rkey)
	if c != nil {
		ans = c.sumP
	}
	rootOut = merge(a, merge(c, d))
	return
}

type randState struct{ x uint32 }

func (r *randState) next() int {
	r.x = r.x*1664525 + 1013904223
	return int(r.x & 0x7fffffff)
}

func referenceSolve(fields []string) (string, error) {
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid test case")
	}
	p := 0
	n, err := strconv.Atoi(fields[p])
	if err != nil {
		return "", fmt.Errorf("parse n: %w", err)
	}
	p++
	if len(fields) < p+n+1 {
		return "", fmt.Errorf("not enough coordinates")
	}
	pos := make([]int64, n)
	for i := 0; i < n; i++ {
		val, err := strconv.ParseInt(fields[p+i], 10, 64)
		if err != nil {
			return "", fmt.Errorf("parse coord %d: %w", i+1, err)
		}
		pos[i] = val
	}
	p += n
	if p >= len(fields) {
		return "", fmt.Errorf("missing m")
	}
	m, err := strconv.Atoi(fields[p])
	if err != nil {
		return "", fmt.Errorf("parse m: %w", err)
	}
	p++
	rnd := &randState{x: 1}
	var root *Node
	for _, v := range pos {
		root = insert(root, v, rnd)
	}
	var outputs []string
	for i := 0; i < m; i++ {
		if p >= len(fields) {
			return "", fmt.Errorf("missing query %d", i+1)
		}
		t, err := strconv.Atoi(fields[p])
		if err != nil {
			return "", fmt.Errorf("parse query type %d: %w", i+1, err)
		}
		p++
		if t == 1 {
			if p+1 >= len(fields) {
				return "", fmt.Errorf("missing data for update %d", i+1)
			}
			idx, _ := strconv.Atoi(fields[p])
			delta, _ := strconv.ParseInt(fields[p+1], 10, 64)
			p += 2
			old := pos[idx-1]
			newVal := old + delta
			root = erase(root, old)
			root = insert(root, newVal, rnd)
			pos[idx-1] = newVal
		} else {
			if p+1 >= len(fields) {
				return "", fmt.Errorf("missing data for query %d", i+1)
			}
			l, _ := strconv.ParseInt(fields[p], 10, 64)
			r, _ := strconv.ParseInt(fields[p+1], 10, 64)
			p += 2
			ans, newRoot := rangeSumP(root, l, r)
			root = newRoot
			outputs = append(outputs, strconv.FormatInt(ans, 10))
		}
	}
	return strings.Join(outputs, "\n"), nil
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	idx := 0
	for _, tc := range testcases {
		line := strings.TrimSpace(tc)
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		expected, err := referenceSolve(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		// rebuild input lines
		p := 0
		n, _ := strconv.Atoi(fields[p])
		p++
		coords := strings.Join(fields[p:p+n], " ")
		p += n
		m, _ := strconv.Atoi(fields[p])
		p++
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		input.WriteString(coords + "\n")
		fmt.Fprintf(&input, "%d\n", m)
		qFields := fields[p:]
		qp := 0
		for i := 0; i < m; i++ {
			t, _ := strconv.Atoi(qFields[qp])
			qp++
			if t == 1 {
				a := qFields[qp]
				b := qFields[qp+1]
				qp += 2
				input.WriteString(fmt.Sprintf("1 %s %s\n", a, b))
			} else {
				a := qFields[qp]
				b := qFields[qp+1]
				qp += 2
				input.WriteString(fmt.Sprintf("2 %s %s\n", a, b))
			}
		}

		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected:\n%s\ngot:\n%s\n", idx, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
