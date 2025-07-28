package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   type Comp struct{ idx, k, t int }
   var group1, group2 []Comp
   comps := make([]Comp, n)
   for i := 0; i < n; i++ {
       var k, t int
       fmt.Fscan(reader, &k, &t)
       c := Comp{i, k, t}
       comps[i] = c
       if t == 1 {
           group1 = append(group1, c)
       } else {
           group2 = append(group2, c)
       }
   }
   ans := make([]int, n)
   used := make([]bool, m+1)
   // schedule t=2: parity-based
   M1 := (m + 1) / 2
   M2 := m / 2
   type interval struct{ l, r int }
   par := map[int][]interval{1: {{1, M1}}, 2: {{1, M2}}}
   sort.Slice(group2, func(i, j int) bool { return group2[i].k > group2[j].k })
   for _, c := range group2 {
       placed := false
       for p := 1; p <= 2; p++ {
           ivs := par[p]
           for idx, iv := range ivs {
               if iv.r-iv.l+1 >= c.k {
                   j := iv.l
                   start := p + 2*(j-1)
                   ans[c.idx] = start
                   // mark used days
                   for d := 0; d < c.k; d++ {
                       day := p + 2*(j-1+d)
                       used[day] = true
                   }
                   // update interval
                   if iv.l+c.k <= iv.r {
                       par[p][idx].l = iv.l + c.k
                   } else {
                       par[p] = append(par[p][:idx], par[p][idx+1:]...)
                   }
                   placed = true
                   break
               }
           }
           if placed {
               break
           }
       }
       if !placed {
           fmt.Println(-1)
           return
       }
   }
   // build free intervals for t=1
   type ivl struct{ l, r int }
   var free []ivl
   i := 1
   for i <= m {
       if used[i] {
           i++
           continue
       }
       l := i
       for i <= m && !used[i] {
           i++
       }
       free = append(free, ivl{l, i - 1})
   }
   sort.Slice(group1, func(i, j int) bool { return group1[i].k > group1[j].k })
   for _, c := range group1 {
       placed := false
       for idx, iv := range free {
           if iv.r-iv.l+1 >= c.k {
               start := iv.l
               ans[c.idx] = start
               if iv.l+c.k <= iv.r {
                   free[idx].l = iv.l + c.k
               } else {
                   free = append(free[:idx], free[idx+1:]...)
               }
               placed = true
               break
           }
       }
       if !placed {
           fmt.Println(-1)
           return
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, x := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       out.WriteString(strconv.Itoa(x))
   }
   out.WriteByte('\n')
}
