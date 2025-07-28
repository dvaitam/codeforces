package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type portal struct {
	l int
	r int
	a int
	b int
}

func escape(x int, ps []portal) int {
	l, r := x, x
	changed := true
	for changed {
		changed = false
		for _, p := range ps {
			if p.r < l || p.l > r {
				continue
			}
			if p.a < l {
				l = p.a
				changed = true
			}
			if p.b > r {
				r = p.b
				changed = true
			}
		}
	}
	return r
}

func solveD(ps []portal, xs []int) []int {
	res := make([]int, len(xs))
	for i, x := range xs {
		res[i] = escape(x, ps)
	}
	return res
}

func generateCasesD() []struct {
	ps []portal
	xs []int
} {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]struct {
		ps []portal
		xs []int
	}, 0, 100)
	for len(cases) < 100 {
		n := rng.Intn(4) + 1
		ps := make([]portal, n)
		for i := 0; i < n; i++ {
			a := rng.Intn(20) + 1
			b := a + rng.Intn(5)
			l := a - rng.Intn(5)
			if l < 1 {
				l = 1
			}
			r := b + rng.Intn(5)
			ps[i] = portal{l: l, r: r, a: a, b: b}
		}
		q := rng.Intn(4) + 1
		xs := make([]int, q)
		for i := 0; i < q; i++ {
			xs[i] = rng.Intn(25) + 1
		}
		cases = append(cases, struct {
			ps []portal
			xs []int
		}{ps: ps, xs: xs})
	}
	return cases
}

func runCase(bin string, ps []portal, xs []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(ps)))
	for _, p := range ps {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", p.l, p.r, p.a, p.b))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(xs)))
	for i, x := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(x))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	exp := solveD(ps, xs)
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad int %q", f)
		}
		if v != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesD()
	for i, tc := range cases {
		if err := runCase(bin, tc.ps, tc.xs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
