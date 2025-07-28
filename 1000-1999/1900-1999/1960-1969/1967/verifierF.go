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

const infVal = int(1e9)

type QuerySet struct {
	q  int
	xs []int
}

type TestCaseF struct {
	n       int
	p       []int
	queries []QuerySet
}

func genCaseF(rng *rand.Rand) TestCaseF {
	n := rng.Intn(5) + 1
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	qs := make([]QuerySet, n)
	for q := 1; q <= n; q++ {
		k := rng.Intn(3)
		xs := make([]int, k)
		for i := range xs {
			xs[i] = rng.Intn(q) + 1
		}
		qs[q-1] = QuerySet{q: q, xs: xs}
	}
	return TestCaseF{n: n, p: p, queries: qs}
}

func expectedF(tc TestCaseF) []int {
	ans := []int{}
	seq := make([]int, 0, tc.n)
	for _, qs := range tc.queries {
		q := qs.q
		xs := qs.xs
		seq = seq[:0]
		for _, v := range tc.p {
			if v <= q {
				seq = append(seq, v)
			}
		}
		m := len(seq)
		pre := make([]int, m)
		stack := make([]int, 0, m)
		for i := 0; i < m; i++ {
			for len(stack) > 0 && seq[stack[len(stack)-1]] <= seq[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				pre[i] = -infVal
			} else {
				pre[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		nxt := make([]int, m)
		stack = stack[:0]
		for i := m - 1; i >= 0; i-- {
			for len(stack) > 0 && seq[stack[len(stack)-1]] <= seq[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				nxt[i] = infVal
			} else {
				nxt[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		diff := make([]int, m)
		for i := 0; i < m; i++ {
			if pre[i] <= -infVal/2 || nxt[i] >= infVal/2 {
				diff[i] = infVal
			} else {
				diff[i] = nxt[i] - pre[i]
			}
		}
		for _, x := range xs {
			sum := 0
			for i := 0; i < m; i++ {
				d := diff[i]
				if d > x {
					sum += x
				} else {
					sum += d
				}
			}
			ans = append(ans, sum)
		}
	}
	return ans
}

func runCaseF(bin string, tc TestCaseF, expect []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, qs := range tc.queries {
		fmt.Fprintf(&sb, "%d\n", len(qs.xs))
		for i, x := range qs.xs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", x)
		}
		sb.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outs) != len(expect) {
		return fmt.Errorf("expected %d numbers got %d", len(expect), len(outs))
	}
	for i, s := range outs {
		var v int
		fmt.Sscan(s, &v)
		if v != expect[i] {
			return fmt.Errorf("index %d expected %d got %d", i, expect[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseF(rng)
		exp := expectedF(tc)
		if err := runCaseF(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d p=%v queries=%v\n", i+1, err, tc.n, tc.p, tc.queries)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
