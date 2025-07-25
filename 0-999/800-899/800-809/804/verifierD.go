package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const INF = 1000000000

var (
	adj  [][]int
	num  []int
	cool []bool
	tot  []int
	sr   [][]int
	ssr  [][]int64
	down []int
	ds   []int
	mem  map[int64]float64
)

func dfs(v, c int) (int, int) {
	num[v] = c
	cnt, sum := 1, len(adj[v])
	for _, w := range adj[v] {
		if num[w] != -1 {
			continue
		}
		a, b := dfs(w, c)
		cnt += a
		sum += b
	}
	return cnt, sum
}

func dfs2(v, pr int) {
	down[v] = 0
	for _, w := range adj[v] {
		if w == pr {
			continue
		}
		dfs2(w, v)
		if down[w]+1 > down[v] {
			down[v] = down[w] + 1
		}
	}
}

func dfs3(v, pr, sof, h int, cur *[]int) {
	ds[v] = down[v]
	if sof+h > ds[v] {
		ds[v] = sof + h
	}
	if h > ds[v] {
		ds[v] = h
	}
	*cur = append(*cur, ds[v])
	mx, mx2 := -INF, -INF
	for _, w := range adj[v] {
		if w == pr {
			continue
		}
		if down[w] > mx {
			mx2 = mx
			mx = down[w]
		} else if down[w] > mx2 {
			mx2 = down[w]
		}
	}
	for _, w := range adj[v] {
		if w == pr {
			continue
		}
		nx := sof
		if down[w] != mx {
			if mx-h+1 > nx {
				nx = mx - h + 1
			}
		} else {
			if mx2-h+1 > nx {
				nx = mx2 - h + 1
			}
		}
		dfs3(w, v, nx, h+1, cur)
	}
}

func prepare(s, c int) {
	dfs2(s, -1)
	cur := make([]int, 0, 16)
	dfs3(s, -1, -INF, 0, &cur)
	maxd := 0
	for _, d := range cur {
		if d > maxd {
			maxd = d
		}
	}
	tot[c] = len(cur)
	sr[c] = make([]int, maxd+1)
	ssr[c] = make([]int64, maxd+1)
	for _, d := range cur {
		sr[c][d]++
	}
	for i := maxd - 1; i >= 0; i-- {
		sr[c][i] += sr[c][i+1]
	}
	for i := maxd; i >= 0; i-- {
		ssr[c][i] = int64(sr[c][i])
		if i < maxd {
			ssr[c][i] += ssr[c][i+1]
		}
	}
}

func solveComp(a, b int) float64 {
	if a > b {
		a, b = b, a
	}
	key := int64(a)*1234567 + int64(a^b)
	if v, ok := mem[key]; ok {
		return v
	}
	if len(sr[a]) > len(sr[b]) {
		a, b = b, a
	}
	da := len(sr[a]) - 1
	db := len(sr[b]) - 1
	r := da
	if db > r {
		r = db
	}
	var sum float64
	for i := 0; i < len(sr[a]); i++ {
		need := r - i
		if need < 0 {
			need = 0
		}
		cc := sr[a][i]
		if i+1 < len(sr[a]) {
			cc -= sr[a][i+1]
		}
		if cc == 0 {
			continue
		}
		if need >= len(sr[b]) {
			sum += float64(tot[b]) * float64(cc) * float64(r)
		} else {
			cnt := sr[b][need]
			lf := tot[b] - cnt
			sum += float64(lf) * float64(cc) * float64(r)
			d0 := need + i + 1
			sum += float64(d0) * float64(cc) * float64(cnt)
			if need+1 < len(sr[b]) {
				sum += float64(cc) * float64(ssr[b][need+1])
			}
		}
	}
	res := sum / float64(tot[a]*tot[b])
	mem[key] = res
	return res
}

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	adj = make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	num = make([]int, n)
	for i := range num {
		num[i] = -1
	}
	cool = make([]bool, n)
	tot = make([]int, n)
	sr = make([][]int, n)
	ssr = make([][]int64, n)
	down = make([]int, n)
	ds = make([]int, n)
	mem = make(map[int64]float64)
	cid := 0
	for i := 0; i < n; i++ {
		if num[i] != -1 {
			continue
		}
		cnt, sum := dfs(i, cid)
		if sum == cnt*2-2 {
			cool[cid] = true
			prepare(i, cid)
		}
		cid++
	}
	var sb strings.Builder
	for i := 0; i < k; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		va := num[a]
		vb := num[b]
		if va == vb || !cool[va] || !cool[vb] {
			sb.WriteString("-1\n")
		} else {
			ans := solveComp(va, vb)
			sb.WriteString(fmt.Sprintf("%.9f\n", ans))
		}
	}
	return strings.TrimSpace(sb.String())
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type edge struct{ u, v int }

func genTest() string {
	r := rand.New(rand.NewSource(rand.Int63()))
	n := r.Intn(4) + 2
	m := n - 1
	k := r.Intn(n) + 1
	lines := []string{fmt.Sprintf("%d %d %d", n, m, k)}
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		lines = append(lines, fmt.Sprintf("%d %d", i, p))
	}
	for i := 0; i < k; i++ {
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		lines = append(lines, fmt.Sprintf("%d %d", a, b))
	}
	return strings.Join(lines, "\n") + "\n"
}

func generateTests() []string {
	tests := make([]string, 100)
	rand.Seed(4)
	for i := 0; i < 100; i++ {
		tests[i] = genTest()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveD(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed. expected %s got %s\ninput:\n%s", i+1, expected, got, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
