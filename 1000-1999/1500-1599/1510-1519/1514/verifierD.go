package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	val int
	cnt int
}

var a []int
var tree []Node

func merge(left, right Node) Node {
	if left.val == right.val {
		return Node{left.val, left.cnt + right.cnt}
	}
	if left.cnt > right.cnt {
		return Node{left.val, left.cnt - right.cnt}
	}
	if right.cnt > left.cnt {
		return Node{right.val, right.cnt - left.cnt}
	}
	return Node{0, 0}
}

func build(idx, l, r int) {
	if l == r {
		tree[idx] = Node{a[l], 1}
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func query(idx, l, r, L, R int) Node {
	if L <= l && r <= R {
		return tree[idx]
	}
	mid := (l + r) / 2
	if R <= mid {
		return query(idx*2, l, mid, L, R)
	}
	if L > mid {
		return query(idx*2+1, mid+1, r, L, R)
	}
	left := query(idx*2, l, mid, L, R)
	right := query(idx*2+1, mid+1, r, L, R)
	return merge(left, right)
}

func solveCase(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}
	idx := 0
	n, _ := strconv.Atoi(fields[idx])
	idx++
	q, _ := strconv.Atoi(fields[idx])
	idx++
	a = make([]int, n+1)
	pos := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		v, _ := strconv.Atoi(fields[idx])
		idx++
		a[i] = v
		pos[v] = append(pos[v], i)
	}
	tree = make([]Node, 4*(n+1))
	build(1, 1, n)
	var outputs []string
	for ; q > 0; q-- {
		l, _ := strconv.Atoi(fields[idx])
		idx++
		r, _ := strconv.Atoi(fields[idx])
		idx++
		cand := query(1, 1, n, l, r).val
		freq := 0
		if cand != 0 {
			arr := pos[cand]
			left := sort.SearchInts(arr, l)
			right := sort.SearchInts(arr, r+1)
			freq = right - left
		}
		length := r - l + 1
		ans := 1
		if tmp := 2*freq - length; tmp > 1 {
			ans = tmp
		}
		outputs = append(outputs, fmt.Sprintf("%d", ans))
	}
	return strings.Join(outputs, "\n")
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idxCase := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idxCase++
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idxCase, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idxCase, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idxCase)
}
