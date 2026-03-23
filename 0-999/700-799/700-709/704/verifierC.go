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

const mod = 1000000007

func combine(t0, t1, a, b int64) (int64, int64) {
	n0 := (t0*a + t1*b) % mod
	n1 := (t0*b + t1*a) % mod
	return n0, n1
}

func solveC(n, m int, clauses [][]int) string {
	p := make([]int, m+1)
	deg := make([]int, m+1)
	adj := make([][]int, m+1)
	u := make([]int, 0, n)
	v := make([]int, 0, n)
	c := 0

	addEdge := func(a, b int) {
		id := len(u)
		u = append(u, a)
		v = append(v, b)
		adj[a] = append(adj[a], id)
		deg[a]++
		adj[b] = append(adj[b], id)
		deg[b]++
	}

	for _, cl := range clauses {
		k := len(cl)
		if k == 1 {
			x := cl[0]
			s := 0
			if x < 0 {
				s = 1
				x = -x
			}
			p[x] ^= 1
			c ^= s
		} else {
			a := cl[0]
			b := cl[1]
			s := 0
			t := 0
			if a < 0 {
				s = 1
				a = -a
			}
			if b < 0 {
				t = 1
				b = -b
			}
			if a == b {
				if s == t {
					p[a] ^= 1
					c ^= s
				} else {
					c ^= 1
				}
			} else {
				addEdge(a, b)
				p[a] ^= 1 ^ t
				p[b] ^= 1 ^ s
				c ^= s | t
			}
		}
	}

	usedEdge := make([]bool, len(u))
	visited := make([]bool, m+1)

	getPath := func(start int) []int {
		verts := []int{start}
		visited[start] = true
		cur := start
		for {
			nextEdge := -1
			for _, e := range adj[cur] {
				if !usedEdge[e] {
					nextEdge = e
					break
				}
			}
			if nextEdge == -1 {
				break
			}
			usedEdge[nextEdge] = true
			nxt := u[nextEdge]
			if nxt == cur {
				nxt = v[nextEdge]
			}
			visited[nxt] = true
			verts = append(verts, nxt)
			cur = nxt
		}
		return verts
	}

	getCycle := func(start int) []int {
		verts := []int{start}
		visited[start] = true
		cur := start
		for {
			nextEdge := -1
			for _, e := range adj[cur] {
				if !usedEdge[e] {
					nextEdge = e
					break
				}
			}
			if nextEdge == -1 {
				break
			}
			usedEdge[nextEdge] = true
			nxt := u[nextEdge]
			if nxt == cur {
				nxt = v[nextEdge]
			}
			if nxt != start {
				visited[nxt] = true
				verts = append(verts, nxt)
			}
			cur = nxt
		}
		return verts
	}

	countPath := func(verts []int) (int64, int64) {
		var dp, ndp [2][2]int64
		p0 := p[verts[0]]
		dp[0][0] = 1
		dp[1][p0] = 1
		for i := 1; i < len(verts); i++ {
			ndp = [2][2]int64{}
			pi := p[verts[i]]
			for prev := 0; prev < 2; prev++ {
				for par := 0; par < 2; par++ {
					val := dp[prev][par]
					if val == 0 {
						continue
					}
					ndp[0][par] += val
					if ndp[0][par] >= mod {
						ndp[0][par] -= mod
					}
					npar := par ^ pi ^ prev
					ndp[1][npar] += val
					if ndp[1][npar] >= mod {
						ndp[1][npar] -= mod
					}
				}
			}
			dp = ndp
		}
		a := dp[0][0] + dp[1][0]
		if a >= mod {
			a -= mod
		}
		b := dp[0][1] + dp[1][1]
		if b >= mod {
			b -= mod
		}
		return a, b
	}

	countCycle := func(verts []int) (int64, int64) {
		var ans [2]int64
		start := verts[0]
		for x0 := 0; x0 < 2; x0++ {
			var dp, ndp [2][2]int64
			dp[x0][p[start]&x0] = 1
			for i := 1; i < len(verts); i++ {
				ndp = [2][2]int64{}
				pi := p[verts[i]]
				for prev := 0; prev < 2; prev++ {
					for par := 0; par < 2; par++ {
						val := dp[prev][par]
						if val == 0 {
							continue
						}
						ndp[0][par] += val
						if ndp[0][par] >= mod {
							ndp[0][par] -= mod
						}
						npar := par ^ pi ^ prev
						ndp[1][npar] += val
						if ndp[1][npar] >= mod {
							ndp[1][npar] -= mod
						}
					}
				}
				dp = ndp
			}
			for last := 0; last < 2; last++ {
				for par := 0; par < 2; par++ {
					val := dp[last][par]
					if val == 0 {
						continue
					}
					totalPar := par ^ (last & x0)
					ans[totalPar] += val
					if ans[totalPar] >= mod {
						ans[totalPar] -= mod
					}
				}
			}
		}
		return ans[0], ans[1]
	}

	total0, total1 := int64(1), int64(0)

	for i := 1; i <= m; i++ {
		if visited[i] {
			continue
		}
		if deg[i] == 0 {
			visited[i] = true
			var a, b int64
			if p[i] == 0 {
				a, b = 2, 0
			} else {
				a, b = 1, 1
			}
			total0, total1 = combine(total0, total1, a, b)
		} else if deg[i] == 1 {
			verts := getPath(i)
			a, b := countPath(verts)
			total0, total1 = combine(total0, total1, a, b)
		}
	}

	for i := 1; i <= m; i++ {
		if !visited[i] {
			verts := getCycle(i)
			a, b := countCycle(verts)
			total0, total1 = combine(total0, total1, a, b)
		}
	}

	target := 1 ^ c
	if target == 0 {
		return fmt.Sprintf("%d", total0)
	}
	return fmt.Sprintf("%d", total1)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(4) + 1
	counts := make([]int, m+1)
	clauses := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		k := rng.Intn(2) + 1
		clause := make([]int, k)
		for j := 0; j < k; j++ {
			var vv int
			attempts := 0
			for {
				attempts++
				vv = rng.Intn(m) + 1
				if counts[vv] < 2 || attempts > 10 {
					break
				}
			}
			counts[vv]++
			if rng.Intn(2) == 0 {
				clause[j] = vv
			} else {
				clause[j] = -vv
			}
		}
		clauses[i] = clause
		sb.WriteString(fmt.Sprintf("%d", k))
		for _, lit := range clause {
			sb.WriteString(fmt.Sprintf(" %d", lit))
		}
		sb.WriteByte('\n')
	}
	expected := solveC(n, m, clauses)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
