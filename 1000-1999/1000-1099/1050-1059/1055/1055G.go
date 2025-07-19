package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// point represents a 2D point or vector
type point struct { x, y int64 }

func (a point) add(b point) point  { return point{a.x + b.x, a.y + b.y} }
func (a point) sub(b point) point  { return point{a.x - b.x, a.y - b.y} }
func (a point) cross(b point) int64 { return a.x*b.y - a.y*b.x }
func (a point) dot(b point) int64   { return a.x*b.x + a.y*b.y }
func (a point) slen() int64         { return a.dot(a) }
// type returns half-plane classification: 0 for upper, 1 for lower
func (a point) typ() int {
   if a.x < 0 || (a.x == 0 && a.y < 0) {
       return 1
   }
   return 0
}

// comparator for sorting by polar angle
func less(a, b point) bool {
   ta, tb := a.typ(), b.typ()
   if ta != tb {
       return ta < tb
   }
   return a.cross(b) > 0
}

func getAlpha(a, b point) float64 {
   return math.Atan2(float64(a.cross(b)), float64(a.dot(b)))
}

func checkPoint(a point, d int64) bool {
   return a.slen() < d*d
}

func checkSeg(a, b point, d int64) bool {
   // projection tests
   if (point{0, 0}.sub(a)).dot(b.sub(a)) < 0 {
       return checkPoint(a, d)
   }
   if (point{0, 0}.sub(b)).dot(a.sub(b)) < 0 {
       return checkPoint(b, d)
   }
   // distance to line
   x := a.cross(b)
   if x < 0 {
       x = -x
   }
   // compare x^2 < d^2 * |b-a|^2
   lhs := x * x
   rhs := d * d * b.sub(a).slen()
   return lhs < rhs
}

func check(pol []point, d int64) bool {
   n := len(pol)
   for i := 0; i+1 < n; i++ {
       if checkSeg(pol[i], pol[i+1], d) {
           return true
       }
   }
   var alpha float64
   for i := 0; i+1 < n; i++ {
       alpha += getAlpha(pol[i], pol[i+1])
   }
   if alpha > 0.5 {
       return true
   }
   return false
}

const inf = 1000000000

// Flow implements a simple max-flow using augmenting DFS
type Flow struct {
   n    int
   es   []edge
   g    [][]int
   used []bool
}
type edge struct{
   to, cap, flow int
}

func NewFlow(n int) *Flow {
   return &Flow{n: n, es: make([]edge, 0), g: make([][]int, n), used: make([]bool, n)}
}

func (f *Flow) AddEdge(v, u, c int) {
   f.g[v] = append(f.g[v], len(f.es))
   f.es = append(f.es, edge{u, c, 0})
   f.g[u] = append(f.g[u], len(f.es))
   f.es = append(f.es, edge{v, 0, 0})
}

func (f *Flow) dfs(v, t int) bool {
   if v == t {
       return true
   }
   f.used[v] = true
   for _, ei := range f.g[v] {
       e := &f.es[ei]
       if !f.used[e.to] && e.cap > e.flow {
           if f.dfs(e.to, t) {
               e.flow++
               f.es[ei^1].flow--
               return true
           }
       }
   }
   return false
}

func (f *Flow) MaxFlow(s, t int) int {
   var res int
   for {
       for i := range f.used {
           f.used[i] = false
       }
       if !f.dfs(s, t) {
           break
       }
       res++
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for {
       var n int
       var w int64
       if _, err := fmt.Fscan(reader, &n, &w); err != nil {
           return
       }
       // read polygon
       pol := make([]point, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &pol[i].x, &pol[i].y)
       }
       var m int
       fmt.Fscan(reader, &m)
       qs := make([]point, m)
       rs := make([]int64, m)
       for i := 0; i < m; i++ {
           fmt.Fscan(reader, &qs[i].x, &qs[i].y, &rs[i])
       }
       // normalize polygon
       pos := 0
       for i := 1; i < n; i++ {
           if pol[i].x < pol[pos].x || (pol[i].x == pol[pos].x && pol[i].y < pol[pos].y) {
               pos = i
           }
       }
       // rotate
       pol = append(pol[pos:], pol[:pos]...)
       p0 := pol[0]
       for i := 0; i < n; i++ {
           pol[i] = pol[i].sub(p0)
       }
       pol = append(pol, pol[0])
       // adjust width
       var wpol int64
       for i := range pol {
           if pol[i].x > wpol {
               wpol = pol[i].x
           }
       }
       w -= wpol
       // build combined edge-sum polygon
       // reversed polygon for jellyfish shape
       qol := make([]point, n)
       for i := 0; i < n; i++ {
           qol[i] = point{0, 0}.sub(pol[i])
       }
       // rotate qol to minimal lex
       pos = 0
       for i := 1; i < n; i++ {
           if qol[i].x < qol[pos].x || (qol[i].x == qol[pos].x && qol[i].y < qol[pos].y) {
               pos = i
           }
       }
       qol = append(qol[pos:], qol[:pos]...)
       qol = append(qol, qol[0])
       // edges
       es0 := make([]point, n)
       for i := 0; i < n; i++ {
           es0[i] = pol[i+1].sub(pol[i])
       }
       es1 := make([]point, n)
       for i := 0; i < n; i++ {
           es1[i] = qol[i+1].sub(qol[i])
       }
       // merge es0 and es1 by angle
       es := make([]point, 0, len(es0)+len(es1))
       i0, i1 := 0, 0
       for i0 < len(es0) && i1 < len(es1) {
           if less(es0[i0], es1[i1]) {
               es = append(es, es0[i0]); i0++
           } else {
               es = append(es, es1[i1]); i1++
           }
       }
       for i0 < len(es0) { es = append(es, es0[i0]); i0++ }
       for i1 < len(es1) { es = append(es, es1[i1]); i1++ }
       // prefix sums of edges
       psum := make([]point, len(es)+1)
       cur := qol[0]
       psum[0] = cur
       for i := 0; i < len(es); i++ {
           cur = cur.add(es[i])
           psum[i+1] = cur
       }
       // build flow
       s := 2*m
       t := s + 1
       flow := NewFlow(t + 1)
       for i := 0; i < m; i++ {
           flow.AddEdge(2*i, 2*i+1, 1)
           l := qs[i].x - wpol - rs[i]
           r := qs[i].x + rs[i]
           if l < 0 {
               flow.AddEdge(s, 2*i, inf)
           }
           if w < r {
               flow.AddEdge(2*i+1, t, inf)
           }
           for j := i + 1; j < m; j++ {
               // translate polygon
               curPoly := make([]point, len(psum))
               delta := qs[j].sub(qs[i])
               for k := range psum {
                   curPoly[k] = psum[k].add(delta)
               }
               if check(curPoly, rs[i]+rs[j]) {
                   flow.AddEdge(2*i+1, 2*j, inf)
                   flow.AddEdge(2*j+1, 2*i, inf)
               }
           }
       }
       res := flow.MaxFlow(s, t)
       fmt.Fprintln(writer, res)
   }
}
