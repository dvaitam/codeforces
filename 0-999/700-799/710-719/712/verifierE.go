package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

type node struct {
	r1, r2 float64
}

func merge(a, b node) node {
	if a.r1 < 0 {
		return b
	}
	if b.r1 < 0 {
		return a
	}
	r1 := a.r1 * b.r1 / (1 - a.r2*(1-b.r1))
	r2 := b.r2 + ((1 - b.r2) * a.r2 * b.r1 / (1 - (1-b.r1)*a.r2))
	return node{r1, r2}
}

func solveE(n int, p []float64, ops []op) []float64 {
	size := n
	t := make([]node, 2*size)
	for i := 0; i < n; i++ {
		t[size+i] = node{p[i], p[i]}
	}
	for i := size - 1; i >= 1; i-- {
		t[i] = merge(t[2*i], t[2*i+1])
	}
	var res []float64
	for _, op := range ops {
		if op.typ == 1 {
			prob := float64(op.a) / float64(op.b)
			pos := size + op.idx - 1
			t[pos] = node{prob, prob}
			for pos >>= 1; pos > 0; pos >>= 1 {
				t[pos] = merge(t[2*pos], t[2*pos+1])
			}
		} else {
			l := op.l + size - 1
			r := op.r + size
			left := node{-1, 0}
			right := node{-1, 0}
			for l < r {
				if l&1 == 1 {
					left = merge(left, t[l])
					l++
				}
				if r&1 == 1 {
					r--
					right = merge(t[r], right)
				}
				l >>= 1
				r >>= 1
			}
			val := merge(left, right)
			res = append(res, val.r1)
		}
	}
	return res
}

type op struct {
	typ    int
	idx, a int
	b      int
	l, r   int
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(5))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		q := rng.Intn(4) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		probs := make([]float64, n)
		for j := 0; j < n; j++ {
			a := rng.Intn(9) + 1
			b := rng.Intn(9) + 1
			probs[j] = float64(a) / float64(b)
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
		var ops []op
		var expOut strings.Builder
		for j := 0; j < q; j++ {
			typ := rng.Intn(2) + 1
			if typ == 1 {
				idx := rng.Intn(n) + 1
				a := rng.Intn(9) + 1
				b := rng.Intn(9) + 1
				ops = append(ops, op{typ: 1, idx: idx, a: a, b: b})
				sb.WriteString(fmt.Sprintf("1 %d %d %d\n", idx, a, b))
			} else {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				ops = append(ops, op{typ: 2, l: l, r: r})
				sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
			}
		}
		results := solveE(n, probs, ops)
		for _, val := range results {
			expOut.WriteString(fmt.Sprintf("%.9f\n", val))
		}
		tests[i] = testCase{in: sb.String(), out: expOut.String()}
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(tc.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
