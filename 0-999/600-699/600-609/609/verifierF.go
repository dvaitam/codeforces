package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD = 1000000007

var rawTestcases = []string{
	"4 5 4 2 7 8 24 0 45 0 29 6 44 8 19 9 43 2 20 1",
	"1 5 35 10 48 7 11 5 59 6 31 5 21 7",
	"2 8 24 0 36 10 9 4 48 2 46 8 40 4 50 2 45 5 7 2 15 10",
	"2 7 17 2 25 7 52 6 2 8 36 4 29 2 44 8 45 3 43 7",
	"2 7 4 8 39 0 22 6 50 6 39 1 23 3 57 3 34 0 31 10",
	"4 1 4 0 9 0 29 9 40 10 28 9",
	"1 8 22 0 51 9 44 8 50 9 1 3 18 7 36 2 50 4 51 10",
	"1 8 7 9 13 8 23 3 21 3 50 7 11 3 26 7 15 1 18 1",
	"1 1 35 8 50 3",
	"2 8 10 4 14 4 32 5 22 2 7 6 16 9 44 7 38 7 33 9 9 4",
	"2 7 12 8 35 9 9 5 13 0 35 0 10 8 37 4 14 4 23 1",
	"5 8 17 3 30 3 37 9 38 8 45 3 53 7 51 5 35 4 48 4 57 8 58 5 37 2 30 4",
	"3 4 15 5 20 4 31 3 13 0 11 6 55 4 44 3",
	"3 2 24 7 30 0 44 6 1 9 13 3",
	"4 1 1 6 5 4 7 6 27 1 41 4",
	"1 2 37 2 7 2 26 7",
	"3 4 1 5 13 3 39 3 14 3 26 1 50 3 22 8",
	"5 2 11 3 26 6 29 1 40 8 47 10 5 8 53 10",
	"3 8 16 5 29 3 46 2 25 4 51 7 56 8 23 9 26 2 9 0 23 5 13 2",
	"2 5 15 8 21 10 39 1 55 7 6 3 10 2 17 4",
	"1 5 26 6 36 1 46 4 24 5 36 10 3 10",
	"3 4 32 8 34 5 47 8 18 4 5 10 30 9 20 8",
	"2 4 1 3 15 0 56 10 27 9 58 0 40 8",
	"5 3 7 9 20 10 30 8 37 4 46 1 44 8 41 2 5 1",
	"2 2 33 2 48 8 44 2 26 10",
	"4 2 17 8 20 5 28 9 47 4 46 8 52 2",
	"4 4 7 2 20 6 23 9 41 4 8 1 35 3 30 9 0 8",
	"1 7 5 10 13 2 39 4 11 2 9 5 2 8 33 8 46 0",
	"4 8 6 8 19 2 27 5 36 2 48 7 7 2 20 2 40 5 60 0 60 2 58 3 35 8",
	"1 2 9 0 52 2 21 8",
	"5 4 8 9 17 3 20 6 27 1 31 7 57 6 38 9 27 10 29 6",
	"4 2 0 10 25 7 26 10 49 4 57 7 46 5",
	"1 2 14 3 1 0 47 6",
	"4 5 7 1 24 4 29 10 42 10 33 1 59 9 22 5 17 6 4 7",
	"3 5 1 2 22 10 47 8 30 0 53 9 41 0 11 9 33 3",
	"1 4 11 4 49 1 16 3 5 4 47 1",
	"2 1 3 6 18 8 25 10",
	"4 8 25 0 29 3 41 4 48 9 3 0 22 5 2 9 28 3 9 9 2 5 18 6 14 1",
	"4 6 11 10 39 2 47 10 49 8 37 3 6 4 53 8 45 9 53 10 58 2",
	"5 7 4 7 19 2 22 9 28 8 41 9 59 7 5 5 20 2 45 8 31 9 20 5 11 0",
	"1 4 18 7 43 4 54 2 31 5 48 8",
	"3 8 12 0 29 0 41 4 60 1 47 7 37 9 1 6 36 9 47 2 14 6 58 10",
	"2 5 31 10 46 8 19 9 35 0 28 4 60 4 15 6",
	"2 4 2 1 30 1 7 4 21 9 17 7 27 3",
	"5 1 7 4 11 7 32 2 35 4 42 4 11 4",
	"2 1 42 10 45 4 43 8",
	"2 8 30 2 35 6 20 1 27 7 3 1 12 8 49 9 46 6 17 1 4 8",
	"4 8 21 7 25 3 31 0 42 4 8 9 43 2 55 8 59 7 59 2 36 9 25 4 6 9",
	"5 7 12 6 14 2 27 8 29 1 33 0 11 8 14 0 8 2 32 1 46 3 16 2 16 8",
	"3 2 8 3 21 3 31 7 29 5 0 10",
	"2 6 34 5 46 9 47 8 44 5 22 9 12 3 44 3 59 10",
	"3 7 18 2 23 7 30 6 17 1 48 6 31 8 53 4 24 2 0 5 33 6",
	"4 6 5 7 6 5 9 2 17 2 46 5 52 1 18 10 17 0 13 1 32 2",
	"1 7 35 8 11 2 52 10 1 0 51 1 17 9 11 10 9 6",
	"4 1 3 4 8 10 11 1 35 10 25 8",
	"3 4 22 10 41 10 44 6 10 2 51 6 50 9 47 0",
	"2 4 6 7 27 0 17 2 31 9 56 1 27 3",
	"1 3 1 7 56 8 13 2 41 2",
	"2 8 5 1 22 8 4 6 23 9 43 1 11 6 22 4 2 7 50 3 3 5",
	"3 5 14 10 36 2 49 5 3 0 17 4 10 4 11 7 41 5",
	"2 6 2 10 47 9 48 5 44 10 46 9 51 7 25 10 52 7",
	"4 1 3 3 7 5 22 9 36 5 22 0",
	"3 8 11 10 13 6 40 3 8 7 24 7 2 10 18 8 10 2 15 7 47 7 27 5",
	"4 6 25 7 26 7 34 1 35 4 29 2 20 5 19 6 36 4 10 3 19 10",
	"1 7 48 9 18 6 21 3 11 5 57 6 55 10 60 0 41 3",
	"4 6 0 9 5 8 21 9 43 3 21 9 59 7 8 10 11 9 27 3 58 7",
	"1 1 28 4 44 5",
	"2 8 12 7 27 8 59 10 5 9 46 5 56 2 41 8 31 6 41 9 4 9",
	"2 4 22 0 27 0 26 3 36 3 55 7 2 2",
	"1 2 19 1 24 3 53 6",
	"1 2 7 2 3 0 44 8",
	"3 2 3 6 11 9 48 9 7 0 39 9",
	"3 5 19 1 30 5 42 3 49 6 49 5 46 10 14 9 34 7",
	"2 8 2 7 48 1 54 10 53 0 15 2 39 4 51 10 30 7 7 8 26 3",
	"4 2 2 6 15 7 23 7 24 0 18 0 10 8",
	"1 2 7 5 46 4 2 2",
	"4 4 16 4 34 10 39 3 48 2 52 6 3 3 15 9 49 10",
	"4 5 22 0 44 0 46 0 49 0 27 3 49 1 36 9 10 2 26 1",
	"2 8 38 7 45 10 58 2 18 5 30 9 60 1 47 6 41 0 58 2 27 9",
	"2 3 0 0 26 1 12 5 26 2 10 5",
	"5 4 11 9 19 2 28 10 38 6 46 1 16 10 33 5 9 10 26 1",
	"1 4 14 9 22 2 7 1 1 7 6 5",
	"2 6 1 7 16 7 59 5 59 5 0 10 18 5 6 8 25 1",
	"3 1 36 6 46 2 47 4 33 9",
	"1 7 25 2 11 6 13 6 60 7 40 1 47 5 52 7 7 7",
	"1 8 17 9 25 5 49 5 17 7 59 6 19 0 10 6 56 3 2 2",
	"5 4 3 3 5 7 34 9 35 5 46 3 24 5 3 3 20 7 60 10",
	"5 6 16 5 20 3 21 3 25 10 36 9 58 0 33 10 58 2 49 10 48 10 49 1",
	"3 2 16 10 29 5 46 3 9 4 29 7",
	"2 3 4 0 13 7 16 3 53 7 30 3",
	"1 3 17 10 3 7 1 4 36 3",
	"3 1 4 10 17 8 40 4 40 0",
	"3 2 13 2 16 6 35 0 33 1 10 0",
	"1 1 48 10 32 7",
	"5 3 6 8 11 6 22 4 27 3 46 7 48 2 16 8 45 5",
	"3 3 30 3 33 7 37 8 17 3 17 3 26 9",
	"3 5 17 0 32 0 45 3 4 1 35 0 15 8 29 2 56 5",
	"3 6 7 9 13 1 18 7 25 9 24 1 11 8 4 1 54 1 2 0",
	"4 4 8 8 26 3 29 1 38 4 58 7 7 10 5 5 7 6",
	"1 5 3 2 43 9 4 9 29 6 15 5 26 10",
}

type testCase struct {
	n, m  int
	frogs []frog
	foods []pair
	input string
}

type frog struct {
	x int
	t int
}

type pair struct {
	p int
	b int
}

func parseCases() []testCase {
	var cases []testCase
	for _, line := range rawTestcases {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		m, _ := strconv.Atoi(fields[pos])
		pos++
		frogs := make([]frog, n)
		for i := 0; i < n; i++ {
			x, _ := strconv.Atoi(fields[pos])
			pos++
			t, _ := strconv.Atoi(fields[pos])
			pos++
			frogs[i] = frog{x: x, t: t}
		}
		foods := make([]pair, m)
		for i := 0; i < m; i++ {
			p, _ := strconv.Atoi(fields[pos])
			pos++
			b, _ := strconv.Atoi(fields[pos])
			pos++
			foods[i] = pair{p: p, b: b}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d", n, m)
		for _, fr := range frogs {
			fmt.Fprintf(&sb, " %d %d", fr.x, fr.t)
		}
		for _, f := range foods {
			fmt.Fprintf(&sb, " %d %d", f.p, f.b)
		}
		sb.WriteByte('\n')
		cases = append(cases, testCase{n: n, m: m, frogs: frogs, foods: foods, input: sb.String()})
	}
	return cases
}

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

func expected(n, m int, frogsData []frog, foods []pair) ([]int, []int) {
	frogs := make([]Frog, n)
	for i, fr := range frogsData {
		frogs[i] = Frog{x: fr.x, r: fr.x + fr.t, id: i}
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
	for _, food := range foods {
		p := food.p
		b := food.b
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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for idx, tc := range cases {
		expCnt, expLen := expected(tc.n, tc.m, tc.frogs, tc.foods)
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("Test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		parts := strings.Fields(got)
		if len(parts) < 2*tc.n {
			fmt.Printf("Test %d invalid output length\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			c, _ := strconv.Atoi(parts[2*i])
			l, _ := strconv.Atoi(parts[2*i+1])
			if c != expCnt[i] || l != expLen[i] {
				fmt.Printf("Test %d failed for frog %d: expected %d %d got %d %d\n", idx+1, i+1, expCnt[i], expLen[i], c, l)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
