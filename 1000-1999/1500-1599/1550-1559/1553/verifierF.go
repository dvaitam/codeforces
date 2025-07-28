package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick { return &Fenwick{n: n, bit: make([]int64, n+2)} }
func (f *Fenwick) Add(idx int, delta int64) {
	idx++
	for idx < len(f.bit) {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}
func (f *Fenwick) Sum(idx int) int64 {
	idx++
	res := int64(0)
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}
func (f *Fenwick) RangeSum(l, r int) int64 {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var n int
	fmt.Fscan(rdr, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &a[i])
	}
	const MAX = 300000
	cnt := NewFenwick(MAX + 2)
	sum := NewFenwick(MAX + 2)
	divf := NewFenwick(MAX + 2)
	ans := make([]int64, n)
	var acc int64
	for i, x := range a {
		var part1 int64
		for m := 0; m <= MAX; m += x {
			r := m + x - 1
			if r > MAX {
				r = MAX
			}
			c := cnt.RangeSum(m, r)
			s := sum.RangeSum(m, r)
			part1 += s - c*int64(m)
		}
		tot := int64(i)
		part2 := int64(x)*tot - divf.Sum(x)
		acc += part1 + part2
		ans[i] = acc
		cnt.Add(x, 1)
		sum.Add(x, int64(x))
		for m := x; m <= MAX; m += x {
			divf.Add(m, int64(x))
		}
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func generateCases() []testCase {
	rand.Seed(6)
	cases := []testCase{}
	fixed := [][]int{{1}, {1, 2}, {3, 3, 3}}
	for _, arr := range fixed {
		var sb strings.Builder
		n := len(arr)
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			v := rand.Intn(50) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierF.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
