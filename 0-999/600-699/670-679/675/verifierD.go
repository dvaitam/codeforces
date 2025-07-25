package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Node struct {
	key   int
	prio  int32
	left  *Node
	right *Node
}

func rotateRight(t *Node) *Node {
	l := t.left
	t.left = l.right
	l.right = t
	return l
}

func rotateLeft(t *Node) *Node {
	r := t.right
	t.right = r.left
	r.left = t
	return r
}

func insert(t *Node, key int) *Node {
	if t == nil {
		return &Node{key: key, prio: rand.Int31()}
	}
	if key < t.key {
		t.left = insert(t.left, key)
		if t.left.prio < t.prio {
			t = rotateRight(t)
		}
	} else {
		t.right = insert(t.right, key)
		if t.right.prio < t.prio {
			t = rotateLeft(t)
		}
	}
	return t
}

func predecessor(t *Node, key int) *Node {
	var res *Node
	for t != nil {
		if key > t.key {
			res = t
			t = t.right
		} else {
			t = t.left
		}
	}
	return res
}

func successor(t *Node, key int) *Node {
	var res *Node
	for t != nil {
		if key < t.key {
			res = t
			t = t.left
		} else {
			t = t.right
		}
	}
	return res
}

func expected(a []int) []string {
	rand.Seed(time.Now().UnixNano())
	n := len(a)
	idx := make(map[int]int, n)
	root := &Node{key: a[0], prio: rand.Int31()}
	idx[a[0]] = 0
	res := make([]string, n-1)
	for i := 1; i < n; i++ {
		x := a[i]
		p := predecessor(root, x)
		s := successor(root, x)
		var parent int
		if p == nil {
			parent = s.key
		} else if s == nil {
			parent = p.key
		} else if idx[p.key] > idx[s.key] {
			parent = p.key
		} else {
			parent = s.key
		}
		res[i-1] = fmt.Sprintf("%d", parent)
		idx[x] = i
		root = insert(root, x)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expectedOut := make([][]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		expectedOut[i] = expected(arr)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n%s", err, out)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		for j := 0; j < len(expectedOut[i]); j++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for test %d\n", i+1)
				os.Exit(1)
			}
			got := outScan.Text()
			if got != expectedOut[i][j] {
				fmt.Printf("test %d failed: expected %s got %s\n", i+1, expectedOut[i][j], got)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
