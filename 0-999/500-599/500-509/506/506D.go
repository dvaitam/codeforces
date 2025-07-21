package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type edge struct{ u, v int }

// DSU (Union-Find)
type dsu struct {
   p []int
}

func newDSU(n int) *dsu {
   d := &dsu{p: make([]int, n+1)}
   for i := 1; i <= n; i++ {
       d.p[i] = i
   }
   return d
}

func (d *dsu) find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.find(d.p[x])
   }
   return d.p[x]
}

func (d *dsu) unite(a, b int) {
   ra := d.find(a)
   rb := d.find(b)
   if ra != rb {
       d.p[rb] = ra
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   // colors indexed 1..m
   cols := make([][]edge, m+1)
   for i := 0; i < m; i++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       cols[c] = append(cols[c], edge{a, b})
   }
   var q int
   fmt.Fscan(in, &q)
   qs := make([][2]int, q)
   // map query pair to indices
   qmap := make(map[[2]int][]int, q)
   for i := 0; i < q; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       if u > v {
           u, v = v, u
       }
       qs[i] = [2]int{u, v}
       key := [2]int{u, v}
       qmap[key] = append(qmap[key], i)
   }
   ans := make([]int, q)

   // threshold for heavy colors
   const B = 316
   heavy := make([]int, 0, m/B+1)

   // separate heavy and light colors
   for c := 1; c <= m; c++ {
       if len(cols[c]) > B {
           heavy = append(heavy, c)
       }
   }
   // process heavy colors
   for _, c := range heavy {
       d := newDSU(n)
       for _, e := range cols[c] {
           d.unite(e.u, e.v)
       }
       for i, qr := range qs {
           if d.find(qr[0]) == d.find(qr[1]) {
               ans[i]++
           }
       }
   }
   // process light colors
   // for temporary DSU and components
   for c := 1; c <= m; c++ {
       if len(cols[c]) == 0 || len(cols[c]) > B {
           continue
       }
       // gather vertices
       var verts []int
       for _, e := range cols[c] {
           verts = append(verts, e.u, e.v)
       }
       sort.Ints(verts)
       // unique
       j := 0
       for i := 1; i < len(verts); i++ {
           if verts[i] != verts[j] {
               j++
               verts[j] = verts[i]
           }
       }
       verts = verts[:j+1]
       // map vertex to index
       idx := make(map[int]int, len(verts))
       for i, v := range verts {
           idx[v] = i
       }
       d := &dsu{p: make([]int, len(verts))}
       for i := range d.p {
           d.p[i] = i
       }
       // unite edges
       for _, e := range cols[c] {
           ui := idx[e.u]
           vi := idx[e.v]
           d.unite(ui, vi)
       }
       // group by root
       comps := make(map[int][]int)
       for _, v := range verts {
           r := d.find(idx[v])
           comps[r] = append(comps[r], v)
       }
       // for each component, enumerate pairs
       for _, comp := range comps {
           L := len(comp)
           if L < 2 {
               continue
           }
           // pairs
           for i := 0; i < L; i++ {
               u := comp[i]
               for j := i + 1; j < L; j++ {
                   v := comp[j]
                   key := [2]int{u, v}
                   if idxs, ok := qmap[key]; ok {
                       for _, qi := range idxs {
                           ans[qi]++
                       }
                   }
               }
           }
       }
   }
   // output
   for i := 0; i < q; i++ {
       fmt.Fprintln(out, ans[i])
   }
}
