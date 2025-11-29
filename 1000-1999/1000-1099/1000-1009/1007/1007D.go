package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const nax = 200*1000 + 105

// Pair represents an interval on the preorder numbering.
type Pair struct {
	first  int
	second int
}

type Sat struct {
	n    int
	ile  int
	imp  [][]int
	vis  []bool
	val  []int
	sort []int
}

func (s *Sat) Add() int {
	s.n += 2
	s.ile += 2
	s.imp = append(s.imp, nil, nil)
	s.vis = append(s.vis, false, false)
	s.val = append(s.val, -1, -1)
	s.sort = append(s.sort, 0, 0)
	return s.n - 2
}

func (s *Sat) Or(a, b int) {
	s.imp[a^1] = append(s.imp[a^1], b)
	s.imp[b^1] = append(s.imp[b^1], a)
}

func (s *Sat) Impl(a, b int) {
	s.Or(a^1, b)
}

func (s *Sat) Nie(a int) {
	s.Or(a^1, a^1)
}

func (s *Sat) Tak(a int) {
	s.Or(a, a)
}

func (s *Sat) Nand(a, b int) {
	s.Or(a^1, b^1)
}

func (s *Sat) dfs(x int) {
	s.vis[x] = true
	for _, i := range s.imp[x^1] {
		if !s.vis[i^1] {
			s.dfs(i ^ 1)
		}
	}
	s.ile--
	s.sort[s.ile] = x
}

func (s *Sat) dfsMark(x int) {
	s.vis[x] = false
	if s.val[x] == -1 {
		if s.val[x^1] == -1 {
			s.val[x] = 1
		} else {
			s.val[x] = 0
		}
	}
	for _, i := range s.imp[x] {
		if s.vis[i] {
			s.dfsMark(i)
		}
	}
}

func (s *Sat) Run() bool {
	for i := 0; i < s.n; i++ {
		if !s.vis[i] {
			s.dfs(i)
		}
	}
	for _, i := range s.sort {
		if s.vis[i] {
			s.dfsMark(i)
		}
	}
	for i := 0; i < s.n; i++ {
		if s.val[i] != 0 {
			for _, x := range s.imp[i] {
				if s.val[x] == 0 {
					return false
				}
			}
		}
	}
	return true
}

// Solver encapsulates per-test data.
type Solver struct {
	drz    [][]int
	roz    []int
	jump   []int
	pre    []int
	post   []int
	fad    []int
	prel   int
	n2     int
	przedz [][]int
	nizej  [][]int
	sat    *Sat
	zmm    []int
}

func (s *Solver) removeParent(w, o int) {
	for i := 0; i < len(s.drz[w]); i++ {
		if s.drz[w][i] == o {
			s.drz[w][i], s.drz[w][len(s.drz[w])-1] = s.drz[w][len(s.drz[w])-1], s.drz[w][i]
			s.drz[w] = s.drz[w][:len(s.drz[w])-1]
			i--
			continue
		}
		s.removeParent(s.drz[w][i], w)
	}
}

func (s *Solver) dfsRoz(v int) {
	s.roz[v] = 1
	for i := 0; i < len(s.drz[v]); i++ {
		child := s.drz[v][i]
		s.fad[child] = v
		s.dfsRoz(child)
		s.roz[v] += s.roz[child]
		if s.roz[child] > s.roz[s.drz[v][0]] {
			s.drz[v][0], s.drz[v][i] = s.drz[v][i], s.drz[v][0]
		}
	}
}

func (s *Solver) dfsPre(v int) {
	if s.jump[v] == 0 {
		s.jump[v] = v
	}
	s.prel++
	s.pre[v] = s.prel
	if len(s.drz[v]) > 0 {
		s.jump[s.drz[v][0]] = s.jump[v]
	}
	for _, child := range s.drz[v] {
		s.dfsPre(child)
	}
	s.post[v] = s.prel
}

func (s *Solver) lca(v, u int) int {
	for s.jump[v] != s.jump[u] {
		if s.pre[v] < s.pre[u] {
			v, u = u, v
		}
		v = s.fad[s.jump[v]]
	}
	if s.pre[v] < s.pre[u] {
		return v
	}
	return u
}

func (s *Solver) pathUp(v, u int) []Pair {
	var ret []Pair
	for s.jump[v] != s.jump[u] {
		ret = append(ret, Pair{first: s.pre[s.jump[v]], second: s.pre[v]})
		v = s.fad[s.jump[v]]
	}
	ret = append(ret, Pair{first: s.pre[u], second: s.pre[v]})
	return ret
}

func (s *Solver) getPath(v, u int) []Pair {
	w := s.lca(v, u)
	ret := s.pathUp(v, w)
	pom := s.pathUp(u, w)
	for i := range ret {
		ret[i].first, ret[i].second = ret[i].second, ret[i].first
	}
	for len(pom) > 0 {
		ret = append(ret, pom[len(pom)-1])
		pom = pom[:len(pom)-1]
	}
	return ret
}

func (s *Solver) wrzuc(w, p, k, a, b, zm int) {
	if k < a || b < p {
		return
	}
	if a <= p && k <= b {
		if len(s.przedz[w]) == 0 || s.przedz[w][len(s.przedz[w])-1] != zm {
			s.przedz[w] = append(s.przedz[w], zm)
		}
		return
	}
	if len(s.nizej[w]) == 0 || s.nizej[w][len(s.nizej[w])-1] != zm {
		s.nizej[w] = append(s.nizej[w], zm)
	}
	mid := (p + k) / 2
	s.wrzuc(w*2, p, mid, a, b, zm)
	s.wrzuc(w*2+1, mid+1, k, a, b, zm)
}

func (s *Solver) dodaj(zm, a, b int) {
	l := s.pre[s.lca(a, b)]
	for _, p := range s.getPath(a, b) {
		if p.first > p.second {
			p.first, p.second = p.second, p.first
		}
		if p.first == l {
			p.first++
		}
		if p.second == l {
			p.second--
		}
		if p.first <= p.second {
			s.wrzuc(1, 0, s.n2-1, p.first, p.second, zm)
		}
	}
}

func (s *Solver) ogarnij(zm, nizejx []int) {
	if len(zm) == 0 {
		return
	}
	nor := s.sat.Add()
	for _, x := range nizejx {
		s.sat.Impl(x, nor)
	}
	last := nor
	for _, z := range zm {
		s.sat.Nand(z, last)
		nl := s.sat.Add()
		s.sat.Impl(last, nl)
		s.sat.Impl(z, nl)
		last = nl
	}
}

func solveCase(n int, edges [][2]int, queries [][4]int) (string, error) {
	s := &Solver{
		drz:  make([][]int, n+1),
		roz:  make([]int, n+1),
		jump: make([]int, n+1),
		pre:  make([]int, n+1),
		post: make([]int, n+1),
		fad:  make([]int, n+1),
		sat:  &Sat{},
	}
	for _, e := range edges {
		a, b := e[0], e[1]
		s.drz[a] = append(s.drz[a], b)
		s.drz[b] = append(s.drz[b], a)
	}

	s.removeParent(1, -1)
	s.dfsRoz(1)
	s.dfsPre(1)

	s.n2 = 1
	for s.n2 <= n {
		s.n2 *= 2
	}
	s.przedz = make([][]int, s.n2*2)
	s.nizej = make([][]int, s.n2*2)

	for _, q := range queries {
		a, b, c, d := q[0], q[1], q[2], q[3]
		zm := s.sat.Add()
		s.zmm = append(s.zmm, zm)
		s.dodaj(zm, a, b)
		s.dodaj(zm^1, c, d)
	}

	for i := 1; i < s.n2*2; i++ {
		s.ogarnij(s.przedz[i], s.nizej[i])
	}

	if !s.sat.Run() {
		return "NO\n", nil
	}
	var buf bytes.Buffer
	buf.WriteString("YES\n")
	for _, zm := range s.zmm {
		if s.sat.val[zm] != 0 {
			buf.WriteString("1\n")
		} else {
			buf.WriteString("2\n")
		}
	}
	return buf.String(), nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	n, err := readInt(in)
	if err != nil {
		return
	}
	m, err := readInt(in)
	if err != nil {
		return
	}
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		a, _ := readInt(in)
		b, _ := readInt(in)
		edges[i] = [2]int{a, b}
	}
	queries := make([][4]int, m)
	for i := 0; i < m; i++ {
		a, _ := readInt(in)
		b, _ := readInt(in)
		c, _ := readInt(in)
		d, _ := readInt(in)
		queries[i] = [4]int{a, b, c, d}
	}
	out, _ := solveCase(n, edges, queries)
	fmt.Print(out)
}

func readInt(r *bufio.Reader) (int, error) {
	var x int
	sign := 1
	c, err := r.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = r.ReadByte()
	}
	if err != nil {
		return 0, err
	}
	if c == '-' {
		sign = -1
		c, err = r.ReadByte()
		if err != nil {
			return 0, err
		}
	}
	for c >= '0' && c <= '9' {
		x = x*10 + int(c-'0')
		c, err = r.ReadByte()
		if err != nil {
			break
		}
	}
	return x * sign, nil
}
