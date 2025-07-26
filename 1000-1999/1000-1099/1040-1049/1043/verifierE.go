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

type pair struct {
	diff int64
	idx  int
}

func expected(a, b []int64, forbid [][2]int) string {
	n := len(a)
	tab := make([]pair, n)
	for i := 0; i < n; i++ {
		tab[i] = pair{b[i] - a[i], i}
	}
	sort.Slice(tab, func(i, j int) bool { return tab[i].diff < tab[j].diff })
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + b[tab[i].idx]
	}
	suf := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suf[i] = suf[i+1] + a[tab[i].idx]
	}
	wynik := make([]int64, n)
	for i := 0; i < n; i++ {
		idx := tab[i].idx
		wynik[idx] = b[idx]*int64(n-1-i) + a[idx]*int64(i)
		wynik[idx] += pref[i]
		wynik[idx] += suf[i+1]
	}
	for _, p := range forbid {
		c, d := p[0], p[1]
		delta := a[c] + b[d]
		if a[d]+b[c] < delta {
			delta = a[d] + b[c]
		}
		wynik[c] -= delta
		wynik[d] -= delta
	}
	var sb strings.Builder
	for i, v := range wynik {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input, want string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %q got %q", want, got)
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

	type test struct {
		a, b   []int64
		forbid [][2]int
	}

	tests := []test{
		{a: []int64{1, 2}, b: []int64{3, 4}, forbid: nil},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		a := make([]int64, n)
		b := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = int64(rng.Intn(20) - 10)
			b[j] = int64(rng.Intn(20) - 10)
			if a[j] == b[j] && rng.Intn(2) == 0 {
				b[j]++
			}
		}
		m := rng.Intn(n)
		forb := make([][2]int, 0, m)
		used := make(map[[2]int]bool)
		for j := 0; j < m; j++ {
			c := rng.Intn(n)
			d := rng.Intn(n)
			if c == d {
				if c+1 < n {
					d = c + 1
				} else {
					d = c - 1
				}
			}
			p := [2]int{c, d}
			if used[p] {
				continue
			}
			used[p] = true
			forb = append(forb, p)
		}
		tests = append(tests, test{a: a, b: b, forbid: forb})
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), len(tc.forbid)))
		for i := 0; i < len(tc.a); i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tc.a[i], tc.b[i]))
		}
		for _, p := range tc.forbid {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0]+1, p[1]+1))
		}
		want := expected(tc.a, tc.b, tc.forbid)
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
