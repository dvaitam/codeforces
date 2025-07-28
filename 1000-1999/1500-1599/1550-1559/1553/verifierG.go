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

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), size: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func sieve(n int) []int {
	spf := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i*i <= n {
				for j := i * i; j <= n; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	return spf
}

func factorize(x int, spf []int) []int {
	res := []int{}
	for x > 1 {
		p := spf[x]
		res = append(res, p)
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func uniqueInts(arr []int) []int {
	m := make(map[int]struct{}, len(arr))
	out := make([]int, 0, len(arr))
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var n, q int
	fmt.Fscan(rdr, &n, &q)
	arr := make([]int, n)
	maxVal := 1000001
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &arr[i])
		if arr[i]+1 > maxVal {
			maxVal = arr[i] + 1
		}
	}
	spf := sieve(maxVal)
	d := newDSU(maxVal + 1)
	for _, x := range arr {
		fac := uniqueInts(factorize(x, spf))
		if len(fac) == 0 {
			continue
		}
		base := fac[0]
		for _, p := range fac[1:] {
			d.union(base, p)
		}
	}
	roots := make([]int, n)
	plusRoots := make([][]int, n)
	edges := make(map[[2]int]struct{})
	for idx, x := range arr {
		fac := uniqueInts(factorize(x, spf))
		baseRoot := d.find(fac[0])
		roots[idx] = baseRoot
		unionSet := []int{baseRoot}
		for _, p := range fac[1:] {
			rp := d.find(p)
			if rp != baseRoot {
				unionSet = append(unionSet, rp)
			}
		}
		fac2 := uniqueInts(factorize(x+1, spf))
		tmp := make([]int, 0, len(fac2))
		for _, p := range fac2 {
			rp := d.find(p)
			tmp = append(tmp, rp)
			unionSet = append(unionSet, rp)
		}
		plusRoots[idx] = uniqueInts(tmp)
		unionSet = uniqueInts(unionSet)
		for i := 0; i < len(unionSet); i++ {
			for j := i + 1; j < len(unionSet); j++ {
				a, b := unionSet[i], unionSet[j]
				if a > b {
					a, b = b, a
				}
				edges[[2]int{a, b}] = struct{}{}
			}
		}
	}
	var outputs []string
	for ; q > 0; q-- {
		var s, t int
		fmt.Fscan(rdr, &s, &t)
		s--
		t--
		if roots[s] == roots[t] {
			outputs = append(outputs, "0")
			continue
		}
		checkSetS := append([]int{roots[s]}, plusRoots[s]...)
		checkSetT := append([]int{roots[t]}, plusRoots[t]...)
		found := false
		for _, rs := range checkSetS {
			for _, rt := range checkSetT {
				if rs == rt {
					found = true
					break
				}
				a, b := rs, rt
				if a > b {
					a, b = b, a
				}
				if _, ok := edges[[2]int{a, b}]; ok {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			outputs = append(outputs, "1")
		} else {
			outputs = append(outputs, "2")
		}
	}
	return strings.Join(outputs, "\n")
}

func randPerm(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = rand.Intn(10) + 1
	}
	return p
}

func generateCases() []testCase {
	rand.Seed(7)
	cases := []testCase{}
	fixed := []struct {
		n, q    int
		arr     []int
		queries [][2]int
	}{
		{2, 1, []int{2, 3}, [][2]int{{1, 2}}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", f.n, f.q)
		for i, v := range f.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for _, qv := range f.queries {
			fmt.Fprintf(&sb, "%d %d\n", qv[0], qv[1])
		}
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(4) + 2
		q := rand.Intn(3) + 1
		arr := randPerm(n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, q)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			s := rand.Intn(n) + 1
			t := rand.Intn(n) + 1
			fmt.Fprintf(&sb, "%d %d\n", s, t)
		}
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierG.go <binary>")
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
