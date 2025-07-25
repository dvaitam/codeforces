package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Frog struct {
	x   int
	r   int
	cnt int
	id  int
}

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

func insertNode(t *Node, key int, val int) *Node {
	if t == nil {
		return &Node{key: key, vals: []int{val}, pr: key*1103515245 + 12345}
	}
	if key == t.key {
		t.vals = append(t.vals, val)
	} else if key < t.key {
		t.left = insertNode(t.left, key, val)
		if t.left.pr > t.pr {
			t = rotateRight(t)
		}
	} else {
		t.right = insertNode(t.right, key, val)
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

func expectedF(n, m int, frogs []Frog, mosqP, mosqB []int) ([]int, []int) {
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
		p := mosqP[k]
		b := mosqB[k]
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
		root = insertNode(root, p, b)
	}
	ansCnt := make([]int, len(frogs))
	ansLen := make([]int, len(frogs))
	for _, f := range frogs {
		ansCnt[f.id] = f.cnt
		ansLen[f.id] = f.r - f.x
	}
	return ansCnt, ansLen
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(parts[pos])
		pos++
		m, _ := strconv.Atoi(parts[pos])
		pos++
		if len(parts) < pos+2*n+2*m {
			fmt.Printf("test %d invalid length\n", idx)
			os.Exit(1)
		}
		frogs := make([]Frog, n)
		for i := 0; i < n; i++ {
			x, _ := strconv.Atoi(parts[pos])
			pos++
			t, _ := strconv.Atoi(parts[pos])
			pos++
			frogs[i] = Frog{x: x, r: x + t, cnt: 0, id: i}
		}
		mosqP := make([]int, m)
		mosqB := make([]int, m)
		for i := 0; i < m; i++ {
			p, _ := strconv.Atoi(parts[pos])
			pos++
			b, _ := strconv.Atoi(parts[pos])
			pos++
			mosqP[i] = p
			mosqB[i] = b
		}
		expCnt, expLen := expectedF(n, m, frogs, mosqP, mosqB)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d", n, m)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, " %d %d", frogs[i].x, frogs[i].r-frogs[i].x)
		}
		for i := 0; i < m; i++ {
			fmt.Fprintf(&buf, " %d %d", mosqP[i], mosqB[i])
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outParts := strings.Fields(strings.TrimSpace(string(out)))
		if len(outParts) < 2*n {
			fmt.Printf("Test %d invalid output length\n", idx)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			gotC, _ := strconv.Atoi(outParts[2*i])
			gotL, _ := strconv.Atoi(outParts[2*i+1])
			if gotC != expCnt[i] || gotL != expLen[i] {
				fmt.Printf("Test %d failed for frog %d: expected %d %d got %d %d\n", idx, i+1, expCnt[i], expLen[i], gotC, gotL)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
