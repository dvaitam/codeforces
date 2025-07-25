package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type rope struct{ l, r int }

func run(bin, input string) (string, error) {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type segTree struct {
	n int
	t []int
}

func newSegTree(n int) *segTree {
	size := 1
	for size < n+2 {
		size <<= 1
	}
	return &segTree{n: size, t: make([]int, 2*size)}
}

func (st *segTree) update(pos, val int) {
	pos += st.n
	if st.t[pos] >= val {
		return
	}
	st.t[pos] = val
	for pos > 1 {
		pos >>= 1
		if st.t[pos] >= val {
			break
		}
		st.t[pos] = val
	}
}

func (st *segTree) query(l, r int) int {
	if l > r {
		return 0
	}
	l += st.n
	r += st.n
	res := 0
	for l <= r {
		if l&1 == 1 {
			res = max(res, st.t[l])
			l++
		}
		if r&1 == 0 {
			res = max(res, st.t[r])
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func solveF(n int, ropes []rope, queries [][2]int) []int {
	sort.Slice(ropes, func(i, j int) bool { return ropes[i].r < ropes[j].r })
	type qu struct{ x, y, idx int }
	qs := make([]qu, len(queries))
	for i, q := range queries {
		qs[i] = qu{q[0], q[1], i}
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].y < qs[j].y })
	st := newSegTree(n)
	ans := make([]int, len(queries))
	j := 0
	for _, q := range qs {
		for j < len(ropes) && ropes[j].r <= q.y {
			if ropes[j].l <= ropes[j].r {
				st.update(ropes[j].l, ropes[j].r)
			}
			j++
		}
		h := q.x
		if h > q.y {
			ans[q.idx] = q.x
			continue
		}
		for {
			mx := st.query(q.x, h)
			if mx > h {
				h = mx
				if h >= q.y {
					if h > q.y {
						h = q.y
					}
					break
				}
			} else {
				break
			}
		}
		if h > q.y {
			h = q.y
		}
		ans[q.idx] = h
	}
	return ans
}

func genCase(rng *rand.Rand) (int, []rope, [][2]int) {
	n := rng.Intn(10) + 1
	m := rng.Intn(5)
	ropes := make([]rope, m)
	usedR := map[int]bool{}
	for i := 0; i < m; i++ {
		r := rng.Intn(n) + 1
		for usedR[r] {
			r = rng.Intn(n) + 1
		}
		usedR[r] = true
		l := rng.Intn(r) + 1
		ropes[i] = rope{l, r}
	}
	q := rng.Intn(5) + 1
	qs := make([][2]int, q)
	for i := 0; i < q; i++ {
		x := rng.Intn(n) + 1
		y := x + rng.Intn(n-x+1)
		qs[i] = [2]int{x, y}
	}
	return n, ropes, qs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, ropes, qs := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d\n", len(ropes))
		for _, rp := range ropes {
			fmt.Fprintf(&sb, "%d %d\n", rp.l, rp.r)
		}
		fmt.Fprintf(&sb, "%d\n", len(qs))
		for _, q := range qs {
			fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
		}
		expAns := solveF(n, ropes, qs)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		lines := strings.Fields(got)
		if len(lines) != len(expAns) {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong number of lines\n", i+1)
			os.Exit(1)
		}
		for idx, ansStr := range lines {
			if ansStr != fmt.Sprint(expAns[idx]) {
				fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %d got %s\n", i+1, sb.String(), expAns[idx], ansStr)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
