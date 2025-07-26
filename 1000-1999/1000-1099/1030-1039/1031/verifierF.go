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
)

const N = 1000000

var spf []int
var dist [][]int
var patMap map[string]int
var M int

func buildData() {
	spf = make([]int, N+1)
	for i := 2; i <= N; i++ {
		if spf[i] == 0 {
			for j := i; j <= N; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	d := make([]int, N+1)
	d[1] = 1
	for i := 2; i <= N; i++ {
		p := spf[i]
		cnt := 0
		x := i
		for x%p == 0 {
			x /= p
			cnt++
		}
		d[i] = d[x] * (cnt + 1)
	}
	maxD := 0
	for i := 1; i <= N; i++ {
		if d[i] > maxD {
			maxD = d[i]
		}
	}
	M = maxD
	type pat struct {
		exps []int
		D    int
	}
	patterns := make([]pat, 0)
	patMap = make(map[string]int)
	patterns = append(patterns, pat{exps: []int{}, D: 1})
	patMap[""] = 0
	var record func([]int, int) int
	record = func(exps []int, D int) int {
		var key string
		for i, e := range exps {
			if i > 0 {
				key += ","
			}
			key += strconv.Itoa(e)
		}
		if id, ok := patMap[key]; ok {
			return id
		}
		id := len(patterns)
		c := make([]int, len(exps))
		copy(c, exps)
		patterns = append(patterns, pat{exps: c, D: D})
		patMap[key] = id
		return id
	}
	var dfs func(cur []int, lastE, curP int)
	dfs = func(cur []int, lastE, curP int) {
		for e := lastE; e >= 1; e-- {
			np := curP * (e + 1)
			if np > M {
				continue
			}
			nxt := append(cur, e)
			id := record(nxt, np)
			_ = id
			dfs(nxt, e, np)
		}
	}
	dfs([]int{}, M-1, 1)
	P := len(patterns)
	nbr := make([][]int, P)
	for i, p := range patterns {
		exps := p.exps
		for j, e := range exps {
			var nxt []int
			if e > 1 {
				nxt = make([]int, len(exps))
				copy(nxt, exps)
				nxt[j] = e - 1
			} else {
				nxt = append([]int{}, exps[:j]...)
				nxt = append(nxt, exps[j+1:]...)
			}
			sort.Slice(nxt, func(a, b int) bool { return nxt[a] > nxt[b] })
			var key string
			for k, ee := range nxt {
				if k > 0 {
					key += ","
				}
				key += strconv.Itoa(ee)
			}
			if id, ok := patMap[key]; ok {
				nbr[i] = append(nbr[i], id)
			}
		}
		for j, e := range exps {
			nxt := make([]int, len(exps))
			copy(nxt, exps)
			nxt[j] = e + 1
			prod := 1
			for _, ee := range nxt {
				prod *= (ee + 1)
				if prod > M {
					break
				}
			}
			if prod <= M {
				sort.Slice(nxt, func(a, b int) bool { return nxt[a] > nxt[b] })
				var key string
				for k, ee := range nxt {
					if k > 0 {
						key += ","
					}
					key += strconv.Itoa(ee)
				}
				if id, ok := patMap[key]; ok {
					nbr[i] = append(nbr[i], id)
				}
			}
		}
		if p.D*2 <= M {
			nxt := append([]int{}, exps...)
			nxt = append(nxt, 1)
			sort.Slice(nxt, func(a, b int) bool { return nxt[a] > nxt[b] })
			var key string
			for k, ee := range nxt {
				if k > 0 {
					key += ","
				}
				key += strconv.Itoa(ee)
			}
			if id, ok := patMap[key]; ok {
				nbr[i] = append(nbr[i], id)
			}
		}
		if len(nbr[i]) > 1 {
			sort.Ints(nbr[i])
			u := nbr[i][:1]
			for _, v := range nbr[i][1:] {
				if v != u[len(u)-1] {
					u = append(u, v)
				}
			}
			nbr[i] = u
		}
	}
	INF := 1 << 30
	dist = make([][]int, M+1)
	for D := 1; D <= M; D++ {
		dist[D] = make([]int, P)
		for i := 0; i < P; i++ {
			dist[D][i] = INF
		}
		var q []int
		for i, p := range patterns {
			if p.D == D {
				dist[D][i] = 0
				q = append(q, i)
			}
		}
		for qi := 0; qi < len(q); qi++ {
			u := q[qi]
			du := dist[D][u]
			for _, v := range nbr[u] {
				if dist[D][v] > du+1 {
					dist[D][v] = du + 1
					q = append(q, v)
				}
			}
		}
	}
}

func sigToPat(n int) int {
	if n == 1 {
		return patMap[""]
	}
	var exps []int
	for n > 1 {
		p := spf[n]
		cnt := 0
		for n%p == 0 {
			n /= p
			cnt++
		}
		exps = append(exps, cnt)
	}
	sort.Slice(exps, func(i, j int) bool { return exps[i] > exps[j] })
	var key string
	for i, e := range exps {
		if i > 0 {
			key += ","
		}
		key += strconv.Itoa(e)
	}
	return patMap[key]
}

func solvePair(a, b int) int {
	INF := 1 << 30
	pa := sigToPat(a)
	pb := sigToPat(b)
	ans := INF
	for D := 1; D <= M; D++ {
		da := dist[D][pa]
		db := dist[D][pb]
		if da+db < ans {
			ans = da + db
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	answers := make([]int, t)
	for i := 0; i < t; i++ {
		a := rng.Intn(N) + 1
		b := rng.Intn(N) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		answers[i] = solvePair(a, b)
	}
	var expect bytes.Buffer
	for i, ans := range answers {
		if i > 0 {
			expect.WriteByte('\n')
		}
		fmt.Fprintf(&expect, "%d", ans)
	}
	expect.WriteByte('\n')
	return sb.String(), expect.String()
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	buildData()
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
