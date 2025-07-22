package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Node struct {
	left, right *Node
	sz          int
	label       int
	prio        int
}

func upd(n *Node) {
	n.sz = 1
	if n.left != nil {
		n.sz += n.left.sz
	}
	if n.right != nil {
		n.sz += n.right.sz
	}
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.prio < b.prio {
		a.right = merge(a.right, b)
		upd(a)
		return a
	}
	b.left = merge(a, b.left)
	upd(b)
	return b
}

func split(n *Node, k int) (a, b *Node) {
	if n == nil {
		return nil, nil
	}
	lsz := 0
	if n.left != nil {
		lsz = n.left.sz
	}
	if k <= lsz {
		a, n.left = split(n.left, k)
		upd(n)
		return a, n
	}
	n.right, b = split(n.right, k-lsz-1)
	upd(n)
	return n, b
}

func solveD(n, m int, ops [][2]int) string {
	rand.Seed(42)
	var root *Node
	for i := 0; i < n; i++ {
		nd := &Node{prio: rand.Int()}
		nd.sz = 1
		root = merge(root, nd)
	}
	assigned := make([]bool, n+1)
	for i := 0; i < m; i++ {
		x := ops[i][0]
		y := ops[i][1]
		if y < 1 || y > root.sz {
			return "-1"
		}
		t1, t2 := split(root, y-1)
		tmid, t3 := split(t2, 1)
		if tmid == nil {
			return "-1"
		}
		if tmid.label != 0 && tmid.label != x {
			return "-1"
		}
		if tmid.label == 0 {
			if assigned[x] {
				return "-1"
			}
			tmid.label = x
			assigned[x] = true
		}
		tmp := merge(t1, t3)
		root = merge(tmid, tmp)
	}
	res := make([]int, 0, n)
	var stack []*Node
	cur := root
	curLabel := 1
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.left
		}
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.label == 0 {
			for curLabel <= n && assigned[curLabel] {
				curLabel++
			}
			node.label = curLabel
			assigned[curLabel] = true
		}
		res = append(res, node.label)
		cur = node.right
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	m := rng.Intn(n) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	arr := make([]int, len(perm))
	copy(arr, perm)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		y := rng.Intn(n) + 1
		x := arr[y-1]
		fmt.Fprintf(&sb, "%d %d\n", x, y)
		arr = append([]int{x}, append(arr[:y-1], arr[y:]...)...)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		scanner := bufio.NewScanner(strings.NewReader(tc))
		scanner.Split(bufio.ScanWords)
		var fields []string
		for scanner.Scan() {
			fields = append(fields, scanner.Text())
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		ops := make([][2]int, m)
		idx := 2
		for j := 0; j < m; j++ {
			x, _ := strconv.Atoi(fields[idx])
			y, _ := strconv.Atoi(fields[idx+1])
			ops[j] = [2]int{x, y}
			idx += 2
		}
		expect := solveD(n, m, ops)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
