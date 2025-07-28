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

type Edge struct {
	to int
	w  int64
}

type Seg struct {
	zone  int
	start int64
	end   int64
}

type solverH struct {
	n     int
	g     [][]Edge
	k     int
	zones []int
	pass  []int64
	fine  []int64
	T     int64

	dist      []int64
	tin       []int
	tout      []int
	timer     int
	zoneNodes [][]int
	counts    [][]int64
}

func (s *solverH) addEdge(u, v int, w int64) {
	s.g[u] = append(s.g[u], Edge{v, w})
	s.g[v] = append(s.g[v], Edge{u, w})
}

func (s *solverH) dfsDist(u, p int, d int64) {
	s.timer++
	s.tin[u] = s.timer
	s.dist[u] = d
	s.zoneNodes[s.zones[u]] = append(s.zoneNodes[s.zones[u]], u)
	for _, e := range s.g[u] {
		if e.to == p {
			continue
		}
		s.dfsDist(e.to, u, d+e.w)
	}
	s.tout[u] = s.timer + 1
}

func (s *solverH) countScans(dv, start, end int64) int64 {
	L := dv - end
	if L < s.T {
		L = s.T
	}
	R := dv - start
	if R <= L {
		return 0
	}
	return (R-1)/s.T - (L-1)/s.T
}

func (s *solverH) dfsCounts(u, p int, segs []Seg) {
	dv := s.dist[u]
	for _, seg := range segs {
		c := s.countScans(dv, seg.start, seg.end)
		s.counts[u][seg.zone] += c
	}
	for _, e := range s.g[u] {
		if e.to == p {
			continue
		}
		zc := s.zones[e.to]
		if len(segs) > 0 && segs[len(segs)-1].zone == zc {
			oldEnd := segs[len(segs)-1].end
			segs[len(segs)-1].end = s.dist[e.to]
			s.dfsCounts(e.to, u, segs)
			segs[len(segs)-1].end = oldEnd
		} else {
			seg := Seg{zone: zc, start: segs[len(segs)-1].end, end: s.dist[e.to]}
			segs = append(segs, seg)
			s.dfsCounts(e.to, u, segs)
			segs = segs[:len(segs)-1]
		}
	}
}

type queryH struct {
	tp  int
	ch  string
	val int64
	u   int
}

func (s *solverH) solve(queries []queryH) []string {
	s.dist = make([]int64, s.n+1)
	s.tin = make([]int, s.n+1)
	s.tout = make([]int, s.n+1)
	s.zoneNodes = make([][]int, s.k)
	s.timer = 0
	s.dfsDist(1, 0, 0)

	s.counts = make([][]int64, s.n+1)
	for i := 0; i <= s.n; i++ {
		s.counts[i] = make([]int64, s.k)
	}
	segs := []Seg{{zone: s.zones[1], start: 0, end: 0}}
	s.dfsCounts(1, 0, segs)

	res := []string{}
	for _, q := range queries {
		switch q.tp {
		case 1:
			idx := int(q.ch[0] - 'A')
			s.pass[idx] = q.val
		case 2:
			idx := int(q.ch[0] - 'A')
			s.fine[idx] = q.val
		case 3:
			u := q.u
			z := s.zones[u]
			nodes := s.zoneNodes[z]
			l := sort.Search(len(nodes), func(i int) bool { return s.tin[nodes[i]] >= s.tin[u] })
			r := sort.Search(len(nodes), func(i int) bool { return s.tin[nodes[i]] >= s.tout[u] })
			best := int64(1<<63 - 1)
			for i := l; i < r; i++ {
				v := nodes[i]
				cost := int64(0)
				for j := 0; j < s.k; j++ {
					if j == z {
						continue
					}
					cnt := s.counts[v][j]
					p := s.pass[j]
					f := s.fine[j]
					if p < f*cnt {
						cost += p
					} else {
						cost += f * cnt
					}
				}
				if cost < best {
					best = cost
				}
			}
			if best == int64(1<<63-1) {
				best = 0
			}
			res = append(res, fmt.Sprintf("%d", best))
		}
	}
	return res
}

func solveCaseH(n int, edges [][3]int, k int, zoneStr string, pass, fine []int64, T int64, queries []queryH) []string {
	s := &solverH{n: n, k: k, T: T}
	s.g = make([][]Edge, n+1)
	for _, e := range edges {
		s.addEdge(e[0], e[1], int64(e[2]))
	}
	s.zones = make([]int, n+1)
	for i := 1; i <= n; i++ {
		s.zones[i] = int(zoneStr[i-1] - 'A')
	}
	s.pass = append([]int64(nil), pass...)
	s.fine = append([]int64(nil), fine...)
	return s.solve(queries)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseH(rng *rand.Rand) (string, []string) {
	n := rng.Intn(4) + 1
	edges := make([][3]int, n-1)
	for i := 0; i < n-1; i++ {
		u := i + 1
		v := i + 2
		w := rng.Intn(5) + 1
		edges[i] = [3]int{u, v, w}
	}
	k := rng.Intn(3) + 1
	zones := make([]byte, n)
	zones[0] = 'A'
	for i := 1; i < n; i++ {
		maxZone := int(zones[i-1]-'A') + rng.Intn(2)
		if maxZone >= k {
			maxZone = k - 1
		}
		zones[i] = byte('A' + maxZone)
	}
	pass := make([]int64, k)
	fine := make([]int64, k)
	for i := 0; i < k; i++ {
		pass[i] = int64(rng.Intn(5) + 1)
		fine[i] = int64(rng.Intn(5) + 1)
	}
	T := int64(rng.Intn(5) + 1)
	qn := rng.Intn(3) + 1
	queries := make([]queryH, qn)
	input := fmt.Sprintf("1\n%d\n", n)
	for _, e := range edges {
		input += fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2])
	}
	input += fmt.Sprintf("%d\n%s\n", k, string(zones))
	for i := 0; i < k; i++ {
		input += fmt.Sprintf("%d ", pass[i])
	}
	input = strings.TrimSpace(input) + "\n"
	for i := 0; i < k; i++ {
		input += fmt.Sprintf("%d ", fine[i])
	}
	input = strings.TrimSpace(input) + "\n"
	input += fmt.Sprintf("%d\n%d\n", T, qn)
	for i := 0; i < qn; i++ {
		tp := rng.Intn(3) + 1
		if tp == 1 {
			ch := string('A' + byte(rng.Intn(k)))
			val := int64(rng.Intn(5) + 1)
			queries[i] = queryH{tp: 1, ch: ch, val: val}
			input += fmt.Sprintf("1 %s %d\n", ch, val)
		} else if tp == 2 {
			ch := string('A' + byte(rng.Intn(k)))
			val := int64(rng.Intn(5) + 1)
			queries[i] = queryH{tp: 2, ch: ch, val: val}
			input += fmt.Sprintf("2 %s %d\n", ch, val)
		} else {
			u := rng.Intn(n) + 1
			queries[i] = queryH{tp: 3, u: u}
			input += fmt.Sprintf("3 %d\n", u)
		}
	}
	expLines := solveCaseH(n, edges, k, string(zones), pass, fine, T, queries)
	return input, expLines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseH(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(exp), len(gotLines), in)
			os.Exit(1)
		}
		for j := range exp {
			if strings.TrimSpace(gotLines[j]) != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp[j], gotLines[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
