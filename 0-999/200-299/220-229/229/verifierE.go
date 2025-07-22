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
	"time"
)

type testCaseE struct {
	n    int
	m    int
	k    []int
	vals [][]int
	ans  float64
}

type pair struct{ val, grp int }

func solveE(tc testCaseE) float64 {
	var v []pair
	for i := 0; i < tc.m; i++ {
		for _, x := range tc.vals[i] {
			v = append(v, pair{x, i})
		}
	}
	sort.Slice(v, func(i, j int) bool { return v[i].val > v[j].val })
	threshold := v[tc.n-1].val
	a := make([]int, tc.m)
	kRem := make([]int, tc.m)
	copy(kRem, tc.k)
	p := 1.0
	done := 0
	for i := 0; i < len(v) && v[i].val != threshold; i++ {
		cur := v[i].grp
		den := kRem[cur]
		if den < 1 {
			den = 1
		}
		p *= float64(a[cur]+1) / float64(den)
		a[cur]++
		kRem[cur]--
		done++
	}
	var ids []int
	for i := 0; i < len(v) && v[i].val >= threshold; i++ {
		if v[i].val == threshold {
			ids = append(ids, v[i].grp)
		}
	}
	s := len(ids)
	tmp := make([]float64, s)
	for j := 0; j < s; j++ {
		cur := ids[j]
		den := kRem[cur]
		if den < 1 {
			den = 1
		}
		tmp[j] = float64(a[cur]+1) / float64(den)
	}
	f := make([][]float64, s+1)
	for i := range f {
		f[i] = make([]float64, s+1)
	}
	f[0][0] = p
	for j := 0; j < s; j++ {
		for i := 0; i <= j; i++ {
			cur := f[i][j]
			if cur == 0 {
				continue
			}
			f[i+1][j+1] += cur * tmp[j] * float64(i+1) / float64(j+1)
			f[i][j+1] += cur * float64(j-i+1) / float64(j+1)
		}
	}
	return f[tc.n-done][s]
}

func genCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	k := make([]int, m)
	vals := make([][]int, m)
	total := 0
	for i := 0; i < m; i++ {
		k[i] = rng.Intn(3) + 1
		vals[i] = make([]int, k[i])
		for j := 0; j < k[i]; j++ {
			vals[i][j] = rng.Intn(20) + 1
		}
		sort.Ints(vals[i])
		sort.Slice(vals[i], func(a, b int) bool { return vals[i][a] > vals[i][b] })
		total += k[i]
	}
	if n > total {
		n = total
	}
	tc := testCaseE{n: n, m: m, k: k, vals: vals}
	tc.ans = solveE(tc)
	return tc
}

func runCaseE(bin string, tc testCaseE) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.m; i++ {
		fmt.Fprintf(&sb, "%d", tc.k[i])
		for _, v := range tc.vals[i] {
			fmt.Fprintf(&sb, " %d", v)
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
	str := strings.TrimSpace(out.String())
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("invalid float output: %v", err)
	}
	diff := val - tc.ans
	if diff < 0 {
		diff = -diff
	}
	if diff > 1e-6 {
		return fmt.Errorf("expected %.10f got %.10f", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
