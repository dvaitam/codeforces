package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   B   = 256
   n, q int
   BN   int
   v    []int
   cnt  int
   ans  []pair
)

type pair struct { a, b int }

func query(a, b int) int {
   if a == -1 {
       return b
   }
   if b == -1 {
       return a
   }
   ans = append(ans, pair{a, b})
   id := cnt
   cnt++
   return id
}

// Comp stores prefix counts of values in [low, up)
type Comp struct { to []int }

func NewComp(pos []int, low, up int) *Comp {
   m := len(pos)
   c := &Comp{to: make([]int, m+1)}
   for i, p := range pos {
       c.to[i+1] = c.to[i]
       x := v[p]
       if x >= low && x < up {
           c.to[i+1]++
       }
   }
   return c
}

func (c *Comp) down(l, r int) (int, int) {
   return c.to[l], c.to[r]
}

// gen builds merge results for positions in value range [low, up)
func gen(low, up int, pos []int) [][]int {
   m := len(pos)
   res := make([][]int, m+1)
   for i := 0; i <= m; i++ {
       row := make([]int, m+1)
       for j := range row {
           row[j] = -1
       }
       res[i] = row
   }
   if m == 0 {
       return res
   }
   if m == 1 {
       res[0][1] = pos[0]
       return res
   }
   mid := (low + up) / 2
   lcomp := NewComp(pos, low, mid)
   ucomp := NewComp(pos, mid, up)
   var lpos, upos []int
   for _, p := range pos {
       if v[p] < mid {
           lpos = append(lpos, p)
       } else {
           upos = append(upos, p)
       }
   }
   lres := gen(low, mid, lpos)
   ures := gen(mid, up, upos)
   for i := 0; i < m; i++ {
       for j := i + 1; j <= m; j++ {
           ll, lr := lcomp.down(i, j)
           ul, ur := ucomp.down(i, j)
           res[i][j] = query(lres[ll][lr], ures[ul][ur])
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &q)
   v = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &v[i])
       v[i]--
   }
   BN = (n + B - 1) / B
   poss := make([][]int, BN)
   for i := 0; i < n; i++ {
       bi := v[i] / B
       poss[bi] = append(poss[bi], i)
   }
   idx := make([]int, n)
   for i := 0; i < n; i++ {
       idx[i] = i
   }
   comps := make([]*Comp, BN)
   bks := make([][][]int, BN)
   for i := 0; i < BN; i++ {
       bks[i] = gen(i*B, i*B+B, poss[i])
       comps[i] = NewComp(idx, i*B, i*B+B)
   }
   qans := make([]int, q)
   for i := range qans {
       qans[i] = -1
   }
   for ti := 0; ti < q; ti++ {
       var l, r int
       fmt.Fscan(in, &l, &r)
       l--
       for i := 0; i < BN; i++ {
           ll, rr := comps[i].down(l, r)
           qans[ti] = query(qans[ti], bks[i][ll][rr])
       }
   }
   // output
   fmt.Fprintln(out, cnt)
   for _, p := range ans {
       fmt.Fprintln(out, p.a+1, p.b+1)
   }
   for i, x := range qans {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprintf(out, "%d", x+1)
   }
   out.WriteByte('\n')
}
