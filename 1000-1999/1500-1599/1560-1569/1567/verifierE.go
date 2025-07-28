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
	first int
	last  int
	pref  int
	suff  int
	len   int
	ans   int64
}

func makeNode(val int) Node {
	return Node{first: val, last: val, pref: 1, suff: 1, len: 1, ans: 1}
}

func merge(a, b Node) Node {
	if a.len == 0 {
		return b
	}
	if b.len == 0 {
		return a
	}
	res := Node{}
	res.len = a.len + b.len
	res.first = a.first
	res.last = b.last
	res.pref = a.pref
	if a.pref == a.len && a.last <= b.first {
		res.pref = a.len + b.pref
	}
	res.suff = b.suff
	if b.suff == b.len && a.last <= b.first {
		res.suff = b.len + a.suff
	}
	res.ans = a.ans + b.ans
	if a.last <= b.first {
		res.ans += int64(a.suff) * int64(b.pref)
	}
	return res
}

func build(tree []Node, arr []int, idx, l, r int) {
	if l == r {
		tree[idx] = makeNode(arr[l])
		return
	}
	mid := (l + r) / 2
	build(tree, arr, idx*2, l, mid)
	build(tree, arr, idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func update(tree []Node, arr []int, idx, l, r, pos, val int) {
	if l == r {
		arr[pos] = val
		tree[idx] = makeNode(val)
		return
	}
	mid := (l + r) / 2
	if pos <= mid {
		update(tree, arr, idx*2, l, mid, pos, val)
	} else {
		update(tree, arr, idx*2+1, mid+1, r, pos, val)
	}
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func query(tree []Node, idx, l, r, ql, qr int) Node {
	if ql <= l && r <= qr {
		return tree[idx]
	}
	mid := (l + r) / 2
	if qr <= mid {
		return query(tree, idx*2, l, mid, ql, qr)
	} else if ql > mid {
		return query(tree, idx*2+1, mid+1, r, ql, qr)
	}
	left := query(tree, idx*2, l, mid, ql, mid)
	right := query(tree, idx*2+1, mid+1, r, mid+1, qr)
	return merge(left, right)
}

func solveCase(n int, arr []int, ops [][3]int) []string {
	tree := make([]Node, 4*n)
	build(tree, arr, 1, 0, n-1)
	var res []string
	for _, op := range ops {
		t, x, y := op[0], op[1], op[2]
		if t == 1 {
			update(tree, arr, 1, 0, n-1, x-1, y)
		} else {
			q := query(tree, 1, 0, n-1, x-1, y-1)
			res = append(res, fmt.Sprint(q.ans))
		}
	}
	return res
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) (int, []int, [][3]int) {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(10) + 1
	}
	ops := make([][3]int, q)
	hasQuery := false
	for i := 0; i < q; i++ {
		t := rng.Intn(2) + 1
		if t == 1 {
			x := rng.Intn(n) + 1
			y := rng.Intn(10) + 1
			ops[i] = [3]int{1, x, y}
		} else {
			l := rng.Intn(n)
			r := rng.Intn(n-l) + l
			ops[i] = [3]int{2, l + 1, r + 1}
			hasQuery = true
		}
	}
	if !hasQuery {
		ops[0][0] = 2
		ops[0][1] = 1
		ops[0][2] = n
	}
	return n, arr, ops
}

func buildInput(n int, arr []int, ops [][3]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(ops)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", op[0], op[1], op[2]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr, ops := genTest(rng)
		input := buildInput(n, arr, ops)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		expected := solveCase(n, append([]int(nil), arr...), ops)
		expectedStr := strings.Join(expected, "\n")
		if out != expectedStr {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expectedStr, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
