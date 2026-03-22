package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// ---------- embedded reference solver for 1648D ----------

const refINF int64 = 2e18

type refSegment struct {
	id int
	L  int
	R  int
	k  int64
}

type refNode struct {
	max_W    int64
	lazy_f   int64
	min_f    int64
	max_f    int64
	max_val  int64
	has_lazy bool
}

var refTree []refNode

func refMax(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func refMin(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func refPushup(node int) {
	refTree[node].max_W = refMax(refTree[node*2].max_W, refTree[node*2+1].max_W)
	refTree[node].min_f = refMin(refTree[node*2].min_f, refTree[node*2+1].min_f)
	refTree[node].max_f = refMax(refTree[node*2].max_f, refTree[node*2+1].max_f)
	refTree[node].max_val = refMax(refTree[node*2].max_val, refTree[node*2+1].max_val)
}

func refApplyLazy(node int, V int64) {
	refTree[node].lazy_f = V
	refTree[node].has_lazy = true
	refTree[node].min_f = V
	refTree[node].max_f = V
	if refTree[node].max_W == -refINF {
		refTree[node].max_val = -refINF
	} else {
		refTree[node].max_val = V + refTree[node].max_W
	}
}

func refPushdown(node int) {
	if refTree[node].has_lazy {
		refApplyLazy(node*2, refTree[node].lazy_f)
		refApplyLazy(node*2+1, refTree[node].lazy_f)
		refTree[node].has_lazy = false
	}
}

func refBuild(node, L, R int) {
	refTree[node].max_W = -refINF
	refTree[node].min_f = -refINF
	refTree[node].max_f = -refINF
	refTree[node].max_val = -refINF
	refTree[node].lazy_f = -refINF
	refTree[node].has_lazy = false
	if L == R {
		return
	}
	mid := (L + R) / 2
	refBuild(node*2, L, mid)
	refBuild(node*2+1, mid+1, R)
}

func refUpdateW(node, L, R, pos int, W int64) {
	if L == R {
		refTree[node].max_W = W
		if W == -refINF {
			refTree[node].max_val = -refINF
		} else {
			if refTree[node].min_f == -refINF {
				refTree[node].max_val = -refINF
			} else {
				refTree[node].max_val = refTree[node].min_f + W
			}
		}
		return
	}
	refPushdown(node)
	mid := (L + R) / 2
	if pos <= mid {
		refUpdateW(node*2, L, mid, pos, W)
	} else {
		refUpdateW(node*2+1, mid+1, R, pos, W)
	}
	refPushup(node)
}

func refChmax(node, L, R, ql, qr int, V int64) {
	if ql <= L && R <= qr {
		if refTree[node].min_f >= V {
			return
		}
		if refTree[node].max_f <= V {
			refApplyLazy(node, V)
			return
		}
	}
	refPushdown(node)
	mid := (L + R) / 2
	if ql <= mid {
		refChmax(node*2, L, mid, ql, qr, V)
	}
	if qr > mid {
		refChmax(node*2+1, mid+1, R, ql, qr, V)
	}
	refPushup(node)
}

func refQueryMaxVal(node, L, R, ql, qr int) int64 {
	if ql <= L && R <= qr {
		return refTree[node].max_val
	}
	refPushdown(node)
	mid := (L + R) / 2
	res := -refINF
	if ql <= mid {
		res = refMax(res, refQueryMaxVal(node*2, L, mid, ql, qr))
	}
	if qr > mid {
		res = refMax(res, refQueryMaxVal(node*2+1, mid+1, R, ql, qr))
	}
	return res
}

func refSolve(input string) string {
	reader := bufio.NewReaderSize(strings.NewReader(input), 1024*1024)
	var n, q int
	fmt.Fscan(reader, &n, &q)

	P1 := make([]int64, n+2)
	P2 := make([]int64, n+2)
	P3 := make([]int64, n+2)

	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(reader, &val)
		P1[i] = P1[i-1] + val
	}
	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(reader, &val)
		P2[i] = P2[i-1] + val
	}
	for i := 1; i <= n; i++ {
		var val int64
		fmt.Fscan(reader, &val)
		P3[i] = P3[i-1] + val
	}

	X := make([]int64, n+1)
	for i := 0; i <= n-1; i++ {
		X[i] = P1[i+1] - P2[i]
	}
	Y := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		Y[i] = P2[i] - P3[i-1]
	}

	segments := make([]refSegment, q)
	add := make([][]refSegment, n+2)
	exp := make([][]int, n+2)

	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &segments[i].L, &segments[i].R, &segments[i].k)
		segments[i].id = i
		add[segments[i].L] = append(add[segments[i].L], segments[i])
		exp[segments[i].R] = append(exp[segments[i].R], i)
	}

	for i := 1; i <= n; i++ {
		sort.Slice(add[i], func(a, b int) bool {
			return add[i][a].k < add[i][b].k
		})
	}

	ptr := make([]int, n+2)
	deleted := make([]bool, q)

	getMaxW := func(L int) int64 {
		for ptr[L] < len(add[L]) && deleted[add[L][ptr[L]].id] {
			ptr[L]++
		}
		if ptr[L] < len(add[L]) {
			return -add[L][ptr[L]].k
		}
		return -refINF
	}

	refTree = make([]refNode, 4*n+4)
	refBuild(1, 1, n)

	dp := make([]int64, n+2)
	for i := range dp {
		dp[i] = -refINF
	}

	ans := -refINF

	for i := 1; i <= n; i++ {
		for _, id := range exp[i-1] {
			deleted[id] = true
			L := segments[id].L
			refUpdateW(1, 1, n, L, getMaxW(L))
		}

		V := X[i-1]
		if dp[i-1] > V {
			V = dp[i-1]
		}
		refChmax(1, 1, n, 1, i, V)

		refUpdateW(1, 1, n, i, getMaxW(i))

		dp[i] = refQueryMaxVal(1, 1, n, 1, i)
		if dp[i] != -refINF {
			ans = refMax(ans, dp[i]+Y[i])
		}
	}

	ans += P3[n]
	return fmt.Sprintf("%d\n", ans)
}

// ---------- verifier harness ----------

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	q := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for row := 0; row < 3; row++ {
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		k := rng.Intn(5)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
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
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		tests = append(tests, genCase(rng))
	}

	for i, tc := range tests {
		exp := refSolve(tc)
		got, err := runExe(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
