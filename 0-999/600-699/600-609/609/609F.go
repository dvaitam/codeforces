package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Frog struct {
	x   int
	r   int
	cnt int
	id  int
}

// Treap node to keep mosquitoes ordered by position
type Node struct {
	key         int
	vals        []int
	pr          int
	left, right *Node
}

func rotateLeft(t *Node) *Node {
	r := t.right
	t.right = r.left
	r.left = t
	return r
}

func rotateRight(t *Node) *Node {
	l := t.left
	t.left = l.right
	l.right = t
	return l
}

func insert(t *Node, key int, val int) *Node {
	if t == nil {
		return &Node{key: key, vals: []int{val}, pr: rand.Int()}
	}
	if key == t.key {
		t.vals = append(t.vals, val)
	} else if key < t.key {
		t.left = insert(t.left, key, val)
		if t.left.pr > t.pr {
			t = rotateRight(t)
		}
	} else {
		t.right = insert(t.right, key, val)
		if t.right.pr > t.pr {
			t = rotateLeft(t)
		}
	}
	return t
}

func lowerBound(t *Node, key int) *Node {
	var best *Node
	for t != nil {
		if t.key >= key {
			best = t
			t = t.left
		} else {
			t = t.right
		}
	}
	return best
}

func deleteNode(t *Node, key int) *Node {
	if t == nil {
		return nil
	}
	if key < t.key {
		t.left = deleteNode(t.left, key)
	} else if key > t.key {
		t.right = deleteNode(t.right, key)
	} else {
		if t.left == nil {
			return t.right
		}
		if t.right == nil {
			return t.left
		}
		if t.left.pr > t.right.pr {
			t = rotateRight(t)
			t.right = deleteNode(t.right, key)
		} else {
			t = rotateLeft(t)
			t.left = deleteNode(t.left, key)
		}
	}
	return t
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	frogs := make([]Frog, n)
	for i := 0; i < n; i++ {
		var x, t int
		fmt.Fscan(reader, &x, &t)
		frogs[i] = Frog{x: x, r: x + t, cnt: 0, id: i}
	}

	sort.Slice(frogs, func(i, j int) bool { return frogs[i].x < frogs[j].x })
	xs := make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = frogs[i].x
	}

	parent := make([]int, n+1)
	for i := 0; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(v int) int {
		if v >= n {
			return n
		}
		if parent[v] != v {
			parent[v] = find(parent[v])
		}
		return parent[v]
	}

	var root *Node

	// mergeRight merges frog i with any following frogs whose starting position is within its reach
	var mergeRight func(int)
	mergeRight = func(i int) {
		for {
			j := find(i + 1)
			if j >= n || frogs[i].r < frogs[j].x {
				break
			}
			if frogs[i].r < frogs[j].r {
				frogs[i].r = frogs[j].r
			}
			parent[j] = find(j + 1)
		}
	}

	// processPending consumes any pending mosquitoes within frog i's reach
	var processPending func(int)
	processPending = func(i int) {
		for {
			node := lowerBound(root, frogs[i].x)
			if node == nil || node.key > frogs[i].r {
				break
			}
			val := node.vals[0]
			if len(node.vals) == 1 {
				root = deleteNode(root, node.key)
			} else {
				node.vals = node.vals[1:]
			}
			frogs[i].r += val
			frogs[i].cnt++
			mergeRight(i)
		}
	}

	for k := 0; k < m; k++ {
		var p, b int
		fmt.Fscan(reader, &p, &b)
		idx := sort.Search(len(xs), func(i int) bool { return xs[i] > p }) - 1
		if idx >= 0 {
			i := find(idx)
			if i < n && frogs[i].r >= p {
				frogs[i].r += b
				frogs[i].cnt++
				mergeRight(i)
				processPending(i)
				continue
			}
		}
		root = insert(root, p, b)
	}

	ansCnt := make([]int, n)
	ansLen := make([]int, n)
	for _, f := range frogs {
		ansCnt[f.id] = f.cnt
		ansLen[f.id] = f.r - f.x
	}
	for i := 0; i < n; i++ {
		fmt.Fprintln(writer, ansCnt[i], ansLen[i])
	}
}
