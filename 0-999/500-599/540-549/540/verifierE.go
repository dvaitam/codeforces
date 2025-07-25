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

type pair struct{ a, b int64 }

// BIT Fenwick tree
type BIT struct {
	n int
	t []int64
}

func NewBIT(n int) *BIT { return &BIT{n, make([]int64, n+1)} }
func (b *BIT) Add(i int, v int64) {
	for i <= b.n {
		b.t[i] += v
		i += i & -i
	}
}
func (b *BIT) Sum(i int) int64 {
	var s int64
	for i > 0 {
		s += b.t[i]
		i -= i & -i
	}
	return s
}

func solveCase(pairs []pair) int64 {
	all := make([]int64, 0, len(pairs)*2)
	for _, p := range pairs {
		all = append(all, p.a, p.b)
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	uniq := all[:0]
	for i, v := range all {
		if i == 0 || v != all[i-1] {
			uniq = append(uniq, v)
		}
	}
	all = uniq
	k := len(all)
	pos2idx := make(map[int64]int, k)
	for i, v := range all {
		pos2idx[v] = i
	}
	M := make([]int, k)
	for i := 0; i < k; i++ {
		M[i] = i
	}
	for _, p := range pairs {
		ia := pos2idx[p.a]
		ib := pos2idx[p.b]
		M[ia], M[ib] = M[ib], M[ia]
	}
	bit := NewBIT(k)
	var invSS int64
	for i := 0; i < k; i++ {
		v := M[i]
		seen := int64(i)
		sum := bit.Sum(v + 1)
		invSS += seen - sum
		bit.Add(v+1, 1)
	}
	var sumA, sumB int64
	for i := 0; i < k; i++ {
		pi := all[i]
		vi := all[M[i]]
		if vi > pi+1 {
			a := pi + 1
			b := vi - 1
			tot := b - a + 1
			l := sort.Search(k, func(j int) bool { return all[j] >= a })
			r := sort.Search(k, func(j int) bool { return all[j] > b }) - 1
			var cnt int64
			if r >= l {
				cnt = int64(r - l + 1)
			}
			sumA += tot - cnt
		}
		if vi < pi-1 {
			a := vi + 1
			b := pi - 1
			tot := b - a + 1
			l := sort.Search(k, func(j int) bool { return all[j] >= a })
			r := sort.Search(k, func(j int) bool { return all[j] > b }) - 1
			var cnt int64
			if r >= l {
				cnt = int64(r - l + 1)
			}
			sumB += tot - cnt
		}
	}
	return invSS + sumA + sumB
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	pairs := make([]pair, n)
	for i := 0; i < n; i++ {
		a := int64(rng.Intn(200) + 1)
		b := int64(rng.Intn(200) + 1)
		if a == b {
			if a == 1 {
				b = 2
			} else {
				b = a - 1
			}
		}
		pairs[i] = pair{a, b}
	}
	var bld strings.Builder
	fmt.Fprintf(&bld, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&bld, "%d %d\n", pairs[i].a, pairs[i].b)
	}
	ans := solveCase(pairs)
	return bld.String(), fmt.Sprintf("%d", ans)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, expect := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, expect, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
