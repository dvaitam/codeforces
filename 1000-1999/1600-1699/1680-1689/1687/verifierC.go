package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type DSU struct{ parent []int }

func NewDSU(n int) *DSU {
	p := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		p[i] = i
	}
	return &DSU{p}
}
func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}
func (d *DSU) Union(x, y int) {
	x = d.Find(x)
	y = d.Find(y)
	if x != y {
		d.parent[x] = y
	}
}

func expected(n, m int, a, b []int, segs [][2]int) string {
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + a[i-1] - b[i-1]
	}
	if pref[n] != 0 {
		return "NO"
	}
	adj := make([][]int, n+1)
	for idx, s := range segs {
		l := s[0]
		r := s[1]
		adj[l] = append(adj[l], idx)
		adj[r] = append(adj[r], idx)
	}
	dsu := NewDSU(n)
	visited := make([]bool, n+1)
	queue := make([]int, 0)
	for i := 0; i <= n; i++ {
		if pref[i] == 0 {
			visited[i] = true
			queue = append(queue, i)
			dsu.Union(i, i+1)
		}
	}
	cnt := make([]int, m)
	segArr := segs
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for len(adj[v]) > 0 {
			id := adj[v][len(adj[v])-1]
			adj[v] = adj[v][:len(adj[v])-1]
			cnt[id]++
			if cnt[id] == 2 {
				l := segArr[id][0]
				r := segArr[id][1]
				if l > r {
					l, r = r, l
				}
				x := dsu.Find(l + 1)
				for x <= r {
					if !visited[x] {
						visited[x] = true
						queue = append(queue, x)
					}
					dsu.Union(x, x+1)
					x = dsu.Find(x)
				}
			}
		}
	}
	for i := 0; i <= n; i++ {
		if !visited[i] {
			return "NO"
		}
	}
	return "YES"
}

func genTest(rng *rand.Rand) (int, int, []int, []int, [][2]int) {
	n := rng.Intn(8) + 2
	m := rng.Intn(n) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10)
		b[i] = rng.Intn(10)
	}
	segs := make([][2]int, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n)
		r := rng.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		segs[i] = [2]int{l, r}
	}
	return n, m, a, b, segs
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(2))
	for i := 0; i < 100; i++ {
		n, m, a, b, segs := genTest(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for _, s := range segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s[0]+1, s[1]))
		}
		expect := expected(n, m, a, b, segs)
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
