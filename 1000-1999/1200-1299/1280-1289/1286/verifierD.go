package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

type matrix struct {
	a [2][2]int64
}

func mul(A, B matrix) matrix {
	var C matrix
	for i := 0; i < 2; i++ {
		for k := 0; k < 2; k++ {
			var s int64
			for j := 0; j < 2; j++ {
				s = (s + A.a[i][j]*B.a[j][k]) % mod
			}
			C.a[i][k] = s
		}
	}
	return C
}

func identityMatrix() matrix {
	return matrix{[2][2]int64{{1, 0}, {0, 1}}}
}

type segTree struct {
	n    int
	tree []matrix
}

func newSegTree(arr []matrix) *segTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &segTree{n, make([]matrix, 2*n)}
	for i := 0; i < len(arr); i++ {
		st.tree[n+i] = arr[i]
	}
	for i := n - 1; i > 0; i-- {
		left := st.tree[i<<1]
		right := st.tree[i<<1|1]
		if left.a == [2][2]int64{} {
			st.tree[i] = right
		} else if right.a == [2][2]int64{} {
			st.tree[i] = left
		} else {
			st.tree[i] = mul(left, right)
		}
	}
	return st
}

func (st *segTree) update(pos int, val matrix) {
	idx := st.n + pos
	st.tree[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		left := st.tree[idx<<1]
		right := st.tree[idx<<1|1]
		if left.a == [2][2]int64{} {
			st.tree[idx] = right
		} else if right.a == [2][2]int64{} {
			st.tree[idx] = left
		} else {
			st.tree[idx] = mul(left, right)
		}
	}
}

func (st *segTree) product() matrix {
	if len(st.tree) == 0 {
		return identityMatrix()
	}
	return st.tree[1]
}

type frac struct {
	num int64
	den int64
}

func less(a, b frac) bool  { return a.num*b.den < b.num*a.den }
func equal(a, b frac) bool { return a.num*b.den == b.num*a.den }

type event struct {
	f   frac
	idx int
	r   int
	c   int
}

func solveD(n int, x []int64, v []int64, pIn []int64) string {
	inv100 := modPow(100, mod-2)
	p := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = pIn[i] % mod * inv100 % mod
	}
	mats := make([]matrix, max(0, n-1))
	for i := 0; i+1 < n; i++ {
		q := (1 - p[i+1] + mod) % mod
		mats[i] = matrix{[2][2]int64{{q, p[i+1]}, {q, p[i+1]}}}
	}
	st := newSegTree(mats)
	probNoCollision := func() int64 {
		prod := st.product()
		q1 := (1 - p[0] + mod) % mod
		res0 := (q1*prod.a[0][0] + p[0]*prod.a[1][0]) % mod
		res1 := (q1*prod.a[0][1] + p[0]*prod.a[1][1]) % mod
		return (res0 + res1) % mod
	}
	events := make([]event, 0)
	for i := 0; i+1 < n; i++ {
		d := x[i+1] - x[i]
		if v[i+1] > v[i] {
			events = append(events, event{frac{d, v[i+1] - v[i]}, i, 0, 0})
		}
		events = append(events, event{frac{d, v[i] + v[i+1]}, i, 1, 0})
		if v[i] > v[i+1] {
			events = append(events, event{frac{d, v[i] - v[i+1]}, i, 1, 1})
		}
	}
	sort.Slice(events, func(i, j int) bool { return less(events[i].f, events[j].f) })
	ans := int64(0)
	prevProb := int64(1)
	i := 0
	for i < len(events) {
		j := i
		f := events[i].f
		for j < len(events) && equal(events[j].f, f) {
			j++
		}
		timeMod := f.num % mod * modPow(f.den%mod, mod-2) % mod
		for k := i; k < j; k++ {
			e := events[k]
			m := mats[e.idx]
			if m.a[e.r][e.c] != 0 {
				m.a[e.r][e.c] = 0
				mats[e.idx] = m
				st.update(e.idx, m)
			}
		}
		newProb := probNoCollision()
		delta := (prevProb - newProb) % mod
		if delta < 0 {
			delta += mod
		}
		ans = (ans + delta*timeMod) % mod
		prevProb = newProb
		i = j
	}
	return fmt.Sprint(ans)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(5))
	tests := make([]string, 100)
	for i := range tests {
		n := rng.Intn(5) + 2
		x := make([]int64, n)
		cur := int64(0)
		for j := 0; j < n; j++ {
			cur += int64(rng.Intn(5) + 1)
			x[j] = cur
		}
		v := make([]int64, n)
		p := make([]int64, n)
		for j := 0; j < n; j++ {
			v[j] = int64(rng.Intn(5) + 1)
			p[j] = int64(rng.Intn(101))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d %d %d\n", x[j], v[j], p[j])
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, input := range tests {
		lines := strings.Fields(input)
		n := atoi(lines[0])
		x := make([]int64, n)
		v := make([]int64, n)
		p := make([]int64, n)
		idx := 1
		for j := 0; j < n; j++ {
			x[j] = int64(atoi(lines[idx]))
			v[j] = int64(atoi(lines[idx+1]))
			p[j] = int64(atoi(lines[idx+2]))
			idx += 3
		}
		exp := solveD(n, x, v, p)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
