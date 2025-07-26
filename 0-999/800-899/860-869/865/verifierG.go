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

const mod int64 = 1000000007

func polyMul(a, b []int64, d int, c []int) []int64 {
	tmp := make([]int64, len(a)+len(b)-1)
	for i, av := range a {
		if av == 0 {
			continue
		}
		for j, bv := range b {
			if bv == 0 {
				continue
			}
			tmp[i+j] = (tmp[i+j] + av*bv) % mod
		}
	}
	for len(tmp) > d {
		k := len(tmp) - 1
		coeff := tmp[k] % mod
		tmp = tmp[:k]
		if coeff != 0 {
			for _, cj := range c {
				idx := k - cj
				tmp[idx] = (tmp[idx] + coeff) % mod
			}
		}
	}
	if len(tmp) < d {
		out := make([]int64, d)
		copy(out, tmp)
		tmp = out
	}
	if len(tmp) > d {
		tmp = tmp[:d]
	}
	return tmp
}

func polyPow(base []int64, exp int64, d int, c []int) []int64 {
	res := []int64{1}
	b := base
	e := exp
	for e > 0 {
		if e&1 == 1 {
			res = polyMul(res, b, d, c)
		}
		b = polyMul(b, b, d, c)
		e >>= 1
	}
	if len(res) < d {
		out := make([]int64, d)
		copy(out, res)
		res = out
	}
	if len(res) > d {
		res = res[:d]
	}
	return res
}

func polyPowX(exp int64, d int, c []int) []int64 {
	base := []int64{0, 1}
	base = polyMul([]int64{1}, base, d, c)
	return polyPow(base, exp, d, c)
}

func solve(F, B int, N int64, p []int64, c []int) int64 {
	d := 0
	for _, ci := range c {
		if ci > d {
			d = ci
		}
	}
	P := make([]int64, d)
	for _, pi := range p {
		poly := polyPowX(pi, d, c)
		for i := 0; i < d; i++ {
			P[i] = (P[i] + poly[i]) % mod
		}
	}
	res := polyPow(P, N, d, c)
	Bseq := make([]int64, d)
	Bseq[0] = 1
	for t := 1; t < d; t++ {
		var val int64
		for _, cj := range c {
			if t-cj >= 0 {
				val = (val + Bseq[t-cj]) % mod
			}
		}
		Bseq[t] = val
	}
	var ans int64
	for i := 0; i < d; i++ {
		ans = (ans + res[i]*Bseq[i]) % mod
	}
	return ans % mod
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

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		F := rng.Intn(3) + 1
		B := rng.Intn(3) + 1
		N := int64(rng.Intn(5) + 1)
		p := make([]int64, F)
		for i := 0; i < F; i++ {
			p[i] = int64(rng.Intn(5) + 1)
		}
		c := make([]int, B)
		for i := 0; i < B; i++ {
			c[i] = rng.Intn(5) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", F, B, N)
		for i := 0; i < F; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", p[i])
		}
		sb.WriteByte('\n')
		for i := 0; i < B; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", c[i])
		}
		sb.WriteByte('\n')
		expected := fmt.Sprintf("%d", solve(F, B, N, p, c))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\n---\ngot:\n%s\n", t+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
