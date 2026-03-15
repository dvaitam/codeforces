package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// solveG is the correct embedded solver for CF 1373/G.
// Uses a segment tree approach with toggle-based pawn placement.
func solveG(input string) string {
	fields := strings.Fields(input)
	idx := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(fields[idx])
		idx++
		return v
	}

	n := nextInt()
	k := nextInt()
	m := nextInt()

	type stNode struct {
		maxA int
		lazy int
		maxV int
	}

	N := 2 * n
	tree := make([]stNode, 4*N+1)

	var build func(node, l, r int)
	build = func(node, l, r int) {
		if l == r {
			tree[node].maxA = l
			tree[node].maxV = 0
			return
		}
		mid := (l + r) / 2
		build(node*2, l, mid)
		build(node*2+1, mid+1, r)
		tree[node].maxA = maxInt(tree[node*2].maxA, tree[node*2+1].maxA)
	}

	pushDown := func(node int) {
		if tree[node].lazy != 0 {
			lz := tree[node].lazy
			tree[node*2].maxA += lz
			tree[node*2].lazy += lz
			tree[node*2+1].maxA += lz
			tree[node*2+1].lazy += lz
			tree[node].lazy = 0
		}
	}

	var updateRange func(node, l, r, ql, qr, val int)
	updateRange = func(node, l, r, ql, qr, val int) {
		if ql <= l && r <= qr {
			tree[node].maxA += val
			tree[node].lazy += val
			return
		}
		pushDown(node)
		mid := (l + r) / 2
		if ql <= mid {
			updateRange(node*2, l, mid, ql, qr, val)
		}
		if qr > mid {
			updateRange(node*2+1, mid+1, r, ql, qr, val)
		}
		tree[node].maxA = maxInt(tree[node*2].maxA, tree[node*2+1].maxA)
	}

	var updateMaxV func(node, l, r, i, val int)
	updateMaxV = func(node, l, r, i, val int) {
		if l == r {
			tree[node].maxV = val
			return
		}
		mid := (l + r) / 2
		if i <= mid {
			updateMaxV(node*2, l, mid, i, val)
		} else {
			updateMaxV(node*2+1, mid+1, r, i, val)
		}
		tree[node].maxV = maxInt(tree[node*2].maxV, tree[node*2+1].maxV)
	}

	var queryMaxA func(node, l, r, ql, qr int) int
	queryMaxA = func(node, l, r, ql, qr int) int {
		if ql <= l && r <= qr {
			return tree[node].maxA
		}
		pushDown(node)
		mid := (l + r) / 2
		res := -1
		if ql <= mid {
			res = queryMaxA(node*2, l, mid, ql, qr)
		}
		if qr > mid {
			res = maxInt(res, queryMaxA(node*2+1, mid+1, r, ql, qr))
		}
		return res
	}

	build(1, 1, N)

	c := make([]int, N+1)
	pawns := make(map[int64]bool)

	var results []string
	for i := 0; i < m; i++ {
		x := nextInt()
		y := nextInt()

		v := y + absInt(x-k)
		key := int64(x)*300000 + int64(y)

		if pawns[key] {
			delete(pawns, key)
			c[v]--
			if c[v] == 0 {
				updateMaxV(1, 1, N, v, 0)
			}
			updateRange(1, 1, N, 1, v, -1)
		} else {
			pawns[key] = true
			c[v]++
			if c[v] == 1 {
				updateMaxV(1, 1, N, v, v)
			}
			updateRange(1, 1, N, 1, v, 1)
		}

		M := tree[1].maxV
		if M == 0 {
			results = append(results, "0")
		} else {
			ans := queryMaxA(1, 1, N, 1, M) - 1 - n
			if ans < 0 {
				ans = 0
			}
			results = append(results, strconv.Itoa(ans))
		}
	}

	return strings.Join(results, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(n) + 1
		m := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, m))
		for i := 0; i < m; i++ {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		input := sb.String()
		expect := solveG(input)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", t+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
