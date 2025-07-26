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

const MOD = 1000000007

type Fenwick struct {
	n int
	t []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, t: make([]int64, n+1)}
}

func (f *Fenwick) Add(i int, v int64) {
	for x := i; x <= f.n; x += x & -x {
		f.t[x] += v
	}
}

func (f *Fenwick) Sum(i int) int64 {
	var s int64
	for x := i; x > 0; x -= x & -x {
		s += f.t[x]
	}
	return s
}

func (f *Fenwick) LowerBound(target int64) int {
	idx := 0
	var sum int64
	bit := 1
	for bit<<1 <= f.n {
		bit <<= 1
	}
	for step := bit; step > 0; step >>= 1 {
		nxt := idx + step
		if nxt <= f.n && sum+f.t[nxt] < target {
			sum += f.t[nxt]
			idx = nxt
		}
	}
	return idx + 1
}

type FenwickMod struct {
	n int
	t []int64
}

func NewFenwickMod(n int) *FenwickMod {
	return &FenwickMod{n: n, t: make([]int64, n+1)}
}

func (f *FenwickMod) Add(i int, v int64) {
	v %= MOD
	if v < 0 {
		v += MOD
	}
	for x := i; x <= f.n; x += x & -x {
		f.t[x] += v
		if f.t[x] >= MOD {
			f.t[x] -= MOD
		}
	}
}

func (f *FenwickMod) Sum(i int) int64 {
	var s int64
	for x := i; x > 0; x -= x & -x {
		s += f.t[x]
		if s >= MOD {
			s -= MOD
		}
	}
	return s
}

func solve(n, q int, a, w []int64, queries [][2]int64) string {
	B := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		B[i] = a[i] - int64(i)
	}
	bitW := NewFenwick(n)
	bitWB := NewFenwickMod(n)
	for i := 1; i <= n; i++ {
		bitW.Add(i, w[i])
		bi := B[i] % MOD
		if bi < 0 {
			bi += MOD
		}
		bitWB.Add(i, w[i]%MOD*bi%MOD)
	}
	var out strings.Builder
	for _, qu := range queries {
		x, y := qu[0], qu[1]
		if x < 0 {
			id := int(-x)
			neww := y
			delta := neww - w[id]
			w[id] = neww
			bitW.Add(id, delta)
			bi := B[id] % MOD
			if bi < 0 {
				bi += MOD
			}
			bitWB.Add(id, delta%MOD*bi%MOD)
		} else {
			l := int(x)
			r := int(y)
			tot := bitW.Sum(r) - bitW.Sum(l-1)
			half := (tot + 1) / 2
			base := bitW.Sum(l - 1)
			m := bitW.LowerBound(base + half)
			if m > r {
				m = r
			}
			sumWl := bitW.Sum(m) - bitW.Sum(l-1)
			sumWr := bitW.Sum(r) - bitW.Sum(m)
			sumWB_l := (bitWB.Sum(m) - bitWB.Sum(l-1) + MOD) % MOD
			sumWB_r := (bitWB.Sum(r) - bitWB.Sum(m) + MOD) % MOD
			med := B[m] % MOD
			if med < 0 {
				med += MOD
			}
			swl := sumWl % MOD
			swr := sumWr % MOD
			c1 := (med*swl%MOD - sumWB_l + MOD) % MOD
			c2 := (sumWB_r - med*swr%MOD + MOD) % MOD
			res := (c1 + c2) % MOD
			out.WriteString(fmt.Sprintf("%d\n", res))
		}
	}
	return out.String()
}

func buildCase(n int, q int, a, w []int64, queries [][2]int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(a[i], 10))
	}
	sb.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(w[i], 10))
	}
	sb.WriteByte('\n')
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	return testCase{input: sb.String(), expected: solve(n, q, a, w, queries)}
}

type testCase struct {
	input    string
	expected string
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	a := make([]int64, n+1)
	w := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = int64(rng.Intn(20))
		w[i] = int64(rng.Intn(10) + 1)
	}
	queries := make([][2]int64, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			id := int64(rng.Intn(n) + 1)
			nw := int64(rng.Intn(10) + 1)
			queries[i] = [2]int64{-id, nw}
		} else {
			l := int64(rng.Intn(n) + 1)
			r := int64(rng.Intn(n-int(l)+1) + int(l))
			queries[i] = [2]int64{l, r}
		}
	}
	return buildCase(n, q, a, w, queries)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.String() != tc.expected {
		return fmt.Errorf("expected %q got %q", tc.expected, out.String())
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
	a := []int64{0, 1}
	w := []int64{0, 1}
	cases := []testCase{
		buildCase(1, 1, a, w, [][2]int64{{1, 1}}),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
