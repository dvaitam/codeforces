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
)

const M = 998244353

type Node struct {
	val  int64
	prod int64
}

func merge(l, r Node) Node {
	newVal := (l.val*r.prod + r.val) % M
	newProd := (l.prod * r.prod) % M
	return Node{newVal, newProd}
}

func power(a, b int64) int64 {
	res := int64(1)
	a %= M
	for b > 0 {
		if b&1 == 1 {
			res = (res * a) % M
		}
		a = (a * a) % M
		b >>= 1
	}
	return res
}

func modInverse(x int64) int64 { return power(x, M-2) }

type BIT struct {
	size int
	tree []int
}

func newBIT(n int) *BIT {
	return &BIT{size: n, tree: make([]int, n+1)}
}

func (b *BIT) add(idx, delta int) {
	for ; idx <= b.size; idx += idx & -idx {
		b.tree[idx] += delta
	}
}

func (b *BIT) prefixSum(idx int) int {
	if idx < 1 {
		return 0
	}
	sum := 0
	for ; idx > 0; idx -= idx & -idx {
		sum += b.tree[idx]
	}
	return sum
}

func (b *BIT) findKth(k int) int {
	if k <= 0 {
		return 0
	}
	pos, sum := 0, 0
	logN := 0
	for (1 << (logN + 1)) <= b.size {
		logN++
	}
	for i := logN; i >= 0; i-- {
		p := 1 << i
		if pos+p <= b.size && sum+b.tree[pos+p] < k {
			sum += b.tree[pos+p]
			pos += p
		}
	}
	return pos + 1
}

func solveCase(input string) (string, error) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	nextInt := func() int {
		if !sc.Scan() {
			return 0
		}
		v, _ := strconv.Atoi(sc.Text())
		return v
	}

	n := nextInt()
	q := nextInt()
	p := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		p[i] = int64(nextInt())
	}

	invP := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		invP[i] = (100 * modInverse(p[i])) % M
	}

	seg := make([]Node, 4*(n+1))
	var build func(int, int, int)
	build = func(v, l, r int) {
		if l == r {
			seg[v] = Node{val: invP[l], prod: invP[l]}
		} else {
			m := (l + r) / 2
			build(v*2, l, m)
			build(v*2+1, m+1, r)
			seg[v] = merge(seg[v*2], seg[v*2+1])
		}
	}

	var query func(int, int, int, int, int) Node
	query = func(v, l, r, L, R int) Node {
		if L > R {
			return Node{0, 1}
		}
		if L == l && R == r {
			return seg[v]
		}
		m := (l + r) / 2
		if R <= m {
			return query(v*2, l, m, L, R)
		}
		if L > m {
			return query(v*2+1, m+1, r, L, R)
		}
		a := query(v*2, l, m, L, m)
		b := query(v*2+1, m+1, r, m+1, R)
		return merge(a, b)
	}

	build(1, 1, n)
	bit := newBIT(n)
	isCP := map[int]bool{1: true}
	bit.add(1, 1)
	total := query(1, 1, n, 1, n).val

	var out bytes.Buffer
	for i := 0; i < q; i++ {
		u := nextInt()
		var cp, cn int
		isU := isCP[u]
		numCp := bit.prefixSum(n)
		if isU {
			rankU := bit.prefixSum(u)
			cp = bit.findKth(rankU - 1)
			nextRank := rankU + 1
			if nextRank > numCp {
				cn = n + 1
			} else {
				cn = bit.findKth(nextRank)
			}
		} else {
			rankBefore := bit.prefixSum(u - 1)
			cp = bit.findKth(rankBefore)
			nextRank := rankBefore + 1
			if nextRank > numCp {
				cn = n + 1
			} else {
				cn = bit.findKth(nextRank)
			}
		}
		val := query(1, 1, n, cp, u-1)
		prod := query(1, 1, n, u, cn-1)
		delta := (val.val * (1 - prod.prod + M)) % M
		if isU {
			total = (total - delta + M) % M
			bit.add(u, -1)
			isCP[u] = false
		} else {
			total = (total + delta) % M
			bit.add(u, 1)
			isCP[u] = true
		}
		fmt.Fprintf(&out, "%d", total)
		if i+1 < q {
			out.WriteByte('\n')
		}
	}
	return out.String(), nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2 // 2..9
	q := rng.Intn(8) + 1 // 1..8
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(100)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d\n", rng.Intn(n-1)+2)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(1264))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		want, err := solveCase(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
