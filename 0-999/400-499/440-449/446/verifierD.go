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

type Matrix struct {
	n  [][]float64
	sz int
}

func NewMatrix(sz int) *Matrix {
	n := make([][]float64, sz)
	for i := range n {
		n[i] = make([]float64, sz)
	}
	return &Matrix{n: n, sz: sz}
}

func (m *Matrix) Multiply(o *Matrix) *Matrix {
	sz := m.sz
	res := NewMatrix(sz)
	for i := 0; i < sz; i++ {
		for k := 0; k < sz; k++ {
			mik := m.n[i][k]
			if mik != 0 {
				for j := 0; j < sz; j++ {
					res.n[i][j] += mik * o.n[k][j]
				}
			}
		}
	}
	return res
}

func (m *Matrix) Pow(exp int64) *Matrix {
	sz := m.sz
	res := NewMatrix(sz)
	for i := 0; i < sz; i++ {
		res.n[i][i] = 1
	}
	base := m
	for exp > 0 {
		if exp&1 == 1 {
			res = res.Multiply(base)
		}
		base = base.Multiply(base)
		exp >>= 1
	}
	return res
}

func expected(n, m int, k int64, a []int, edges [][2]int) float64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	comp := make([]int, n+1)
	W := 0
	for i := 1; i <= n; i++ {
		if a[i] == 0 && comp[i] == 0 {
			W++
			queue := []int{i}
			comp[i] = W
			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]
				for _, v := range adj[u] {
					if a[v] == 0 && comp[v] == 0 {
						comp[v] = W
						queue = append(queue, v)
					}
				}
			}
		}
	}
	B := 0
	for i := 1; i <= n; i++ {
		if a[i] == 1 {
			B++
			comp[i] = B
		}
	}
	cnt := make([][]int, B)
	for i := 0; i < B; i++ {
		cnt[i] = make([]int, W)
	}
	Cnt := make([][]int, B)
	for i := 0; i < B; i++ {
		Cnt[i] = make([]int, B)
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		if a[u] == 1 && a[v] == 1 {
			ui, vi := comp[u]-1, comp[v]-1
			Cnt[ui][vi]++
			Cnt[vi][ui]++
		} else if a[u] == 1 && a[v] == 0 {
			ui, vj := comp[u]-1, comp[v]-1
			cnt[ui][vj]++
		} else if a[u] == 0 && a[v] == 1 {
			ui, vj := comp[v]-1, comp[u]-1
			cnt[ui][vj]++
		}
	}
	s1 := make([]int, W)
	for j := 0; j < W; j++ {
		sum := 0
		for i := 0; i < B; i++ {
			sum += cnt[i][j]
		}
		s1[j] = sum
	}
	s2 := make([]int, B)
	for i := 0; i < B; i++ {
		sum := 0
		for j := 0; j < B; j++ {
			sum += Cnt[j][i]
		}
		for j := 0; j < W; j++ {
			sum += cnt[i][j]
		}
		s2[i] = sum
	}
	p := NewMatrix(B)
	for i := 0; i < B; i++ {
		if s2[i] == 0 {
			continue
		}
		for j := 0; j < B; j++ {
			p.n[i][j] = float64(Cnt[i][j]) / float64(s2[i])
			extra := 0.0
			for k2 := 0; k2 < W; k2++ {
				if s1[k2] == 0 {
					continue
				}
				extra += (float64(cnt[i][k2]) / float64(s2[i])) * (float64(cnt[j][k2]) / float64(s1[k2]))
			}
			p.n[i][j] += extra
		}
	}
	exp := k - 2
	pPow := p.Pow(exp)
	j0 := comp[1] - 1
	ans := 0.0
	if j0 >= 0 && j0 < W {
		for i := 0; i < B; i++ {
			if s1[j0] > 0 {
				ans += (float64(cnt[i][j0]) / float64(s1[j0])) * pPow.n[j0][B-1]
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(4) + 2
	maxM := n * (n - 1) / 2
	m := n - 1 + rng.Intn(maxM-(n-1)+1)
	// create connected graph via tree
	edges := make([][2]int, 0, m)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{u, i})
	}
	used := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		used[[2]int{e[0], e[1]}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a := u
		b := v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, [2]int{u, v})
	}
	k := int64(rng.Intn(4) + 2)
	a := make([]int, n+1)
	trapCnt := 1
	for i := 2; i < n; i++ {
		if rng.Intn(2) == 0 && trapCnt < 3 {
			a[i] = 1
			trapCnt++
		}
	}
	a[n] = 1
	a[1] = 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), k))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	ans := expected(n, len(edges), k, a, edges)
	return sb.String(), ans
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if diff := got - exp; diff < -1e-4 || diff > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d: expected %.6f got %.6f\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
