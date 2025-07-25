package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// form represents a linear function f(t) = a*t + b
type form struct {
   a, b int64
}
func (f form) eval(t int64) int64 { return f.a*t + f.b }
func addForm(x, y form) form { return form{x.a + y.a, x.b + y.b} }
func subForm(x, y form) form { return form{x.a - y.a, x.b - y.b} }
func cross(x, y form) int64 { return x.a*y.b - x.b*y.a }
// cross3 returns cross(b-a, c-a)
func cross3(a, b, c form) int64 { return cross(subForm(b, a), subForm(c, a)) }

// Edge of tree: to node and weight form
type edge struct {
   to int
   f  form
}

// solver for centroid decomposition and convex hull
type solver struct {
   adj     [][]edge
   m       int64
   n       int
   sts     []int
   removed []bool
   cands   []form
}

// calculate subtree sizes
func (s *solver) calcSts(v, p int) {
   s.sts[v] = 1
   for _, e := range s.adj[v] {
       if !s.removed[e.to] && e.to != p {
           s.calcSts(e.to, v)
           s.sts[v] += s.sts[e.to]
       }
   }
}
// find centroid
func (s *solver) getCentroid(v, p, total int) int {
   for _, e := range s.adj[v] {
       if !s.removed[e.to] && e.to != p && s.sts[e.to]*2 > total {
           return s.getCentroid(e.to, v, total)
       }
   }
   return v
}
// collect forms from v to leaves
func (s *solver) getForms(v, p int, f form, out *[]form) {
   leaf := true
   for _, e := range s.adj[v] {
       if !s.removed[e.to] && e.to != p {
           leaf = false
           s.getForms(e.to, v, addForm(f, e.f), out)
       }
   }
   if leaf {
       *out = append(*out, f)
   }
}
// keep only maximal forms: build hull
func keepMaximal(forms *[]form) {
   a := *forms
   sort.Slice(a, func(i, j int) bool {
       if a[i].a != a[j].a {
           return a[i].a < a[j].a
       }
       return a[i].b > a[j].b
   })
   // filter increasing a
   w := make([]form, 0, len(a))
   for i := range a {
       if len(w) == 0 || a[i].a > w[len(w)-1].a {
           w = append(w, a[i])
       }
   }
   // filter decreasing b
   w2 := make([]form, 0, len(w))
   for i := range w {
       for len(w2) > 0 && w[i].b >= w2[len(w2)-1].b {
           w2 = w2[:len(w2)-1]
       }
       w2 = append(w2, w[i])
   }
   // convex hull
   res := make([]form, 0, len(w2))
   for _, f := range w2 {
       for len(res) >= 2 && cross3(f, res[len(res)-2], res[len(res)-1]) >= 0 {
           res = res[:len(res)-1]
       }
       res = append(res, f)
   }
   *forms = res
}
// formRange associates form with its valid t range [l,r]
type formRange struct {
   f      form
   l, r   int64
}
// create formRanges from hull forms
func (s *solver) formAndRanges(forms []form) []formRange {
   k := len(forms)
   fr := make([]formRange, k)
   for i := 0; i < k; i++ {
       fr[i].f = forms[i]
       if i == 0 {
           fr[i].l = 0
       } else {
           fr[i].l = fr[i-1].r
       }
       if i+1 < k {
           // solve forms[i].eval(t) >= forms[i+1].eval(t)
           num := forms[i].b - forms[i+1].b
           den := forms[i+1].a - forms[i].a
           t := num / den
           if t > s.m-1 {
               t = s.m - 1
           }
           fr[i].r = t
       } else {
           fr[i].r = s.m - 1
       }
   }
   return fr
}
// add pair sums between two formRange lists
func (s *solver) addPairs(a, b []formRange) {
   j := 0
   for i := range a {
       for j < len(b) && b[j].r < a[i].l {
           j++
       }
       for k := j; k < len(b) && b[k].l <= a[i].r; k++ {
           s.cands = append(s.cands, addForm(a[i].f, b[k].f))
       }
   }
}
// process centroid c
func (s *solver) solveC(c int) {
   var formsList [][]form
   for _, e := range s.adj[c] {
       if s.removed[e.to] {
           continue
       }
       var forms []form
       s.getForms(e.to, c, e.f, &forms)
       formsList = append(formsList, forms)
   }
   var rangesList [][]formRange
   for _, forms := range formsList {
       keepMaximal(&forms)
       fr := s.formAndRanges(forms)
       for _, f := range forms {
           s.cands = append(s.cands, f)
       }
       rangesList = append(rangesList, fr)
   }
   // pairs between subtrees
   for i := 0; i < len(rangesList); i++ {
       for j := i + 1; j < len(rangesList); j++ {
           s.addPairs(rangesList[i], rangesList[j])
       }
   }
}
// recursive solve
func (s *solver) solve(v int) {
   s.calcSts(v, -1)
   c := s.getCentroid(v, -1, s.sts[v])
   s.solveC(c)
   s.removed[c] = true
   for _, e := range s.adj[c] {
       if !s.removed[e.to] {
           s.solve(e.to)
       }
   }
}
// ensure degree <=3 by expanding
func expand(v, p int, adj [][]edge, adj2 *[][]edge) int {
   newV := len(*adj2)
   *adj2 = append(*adj2, nil)
   var children []struct{ node int; f form }
   for _, e := range adj[v] {
       if e.to == p {
           continue
       }
       u := expand(e.to, v, adj, adj2)
       children = append(children, struct{ node int; f form }{u, e.f})
   }
   // split if more than 2 children
   for len(children) > 2 {
       ca := children[len(children)-1]
       cb := children[len(children)-2]
       children = children[:len(children)-2]
       join := len(*adj2)
       *adj2 = append(*adj2, nil)
       // connect ca
       (*adj2)[join] = append((*adj2)[join], edge{ca.node, ca.f})
       (*adj2)[ca.node] = append((*adj2)[ca.node], edge{join, ca.f})
       // connect cb
       (*adj2)[join] = append((*adj2)[join], edge{cb.node, cb.f})
       (*adj2)[cb.node] = append((*adj2)[cb.node], edge{join, cb.f})
       // add to children
       children = append(children, struct{ node int; f form }{join, form{0, 0}})
   }
   // connect remaining to newV
   for _, ch := range children {
       *adj2[newV] = append(*adj2[newV], edge{ch.node, ch.f})
       *adj2[ch.node] = append(*adj2[ch.node], edge{newV, ch.f})
   }
   return newV
}
func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   var m int64
   fmt.Fscan(in, &n, &m)
   adj := make([][]edge, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       var a, b int64
       fmt.Fscan(in, &u, &v, &a, &b)
       u--, v--
       f := form{a, b}
       adj[u] = append(adj[u], edge{v, f})
       adj[v] = append(adj[v], edge{u, f})
   }
   // expand to adj2
   var adj2 [][]edge
   expand(0, -1, adj, &adj2)
   s := &solver{
       adj:     adj2,
       m:       m,
       n:       len(adj2),
       sts:     make([]int, len(adj2)),
       removed: make([]bool, len(adj2)),
   }
   s.solve(0)
   // add zero form
   s.cands = append(s.cands, form{0, 0})
   keepMaximal(&s.cands)
   // output
   idx := 0
   for t := int64(0); t < m; t++ {
       for idx+1 < len(s.cands) && s.cands[idx+1].eval(t) > s.cands[idx].eval(t) {
           idx++
       }
       if t > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, s.cands[idx].eval(t))
   }
   out.WriteByte('\n')
