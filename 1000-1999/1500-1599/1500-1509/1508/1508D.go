package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Point struct { x, y int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   p := make([]Point, n)
   to := make([]int, n)
   ord := make([]int, 0, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i].x, &p[i].y, &to[i])
       to[i]--
       if to[i] != i {
           ord = append(ord, i)
       }
   }
   if len(ord) == 0 {
       fmt.Fprintln(out, 0)
       return
   }
   // find minimal point
   mn := ord[0]
   for _, v := range ord[1:] {
       if p[v].x < p[mn].x || (p[v].x == p[mn].x && p[v].y < p[mn].y) {
           mn = v
       }
   }
   // rotate ord so mn is first
   idx := 0
   for i, v := range ord {
       if v == mn {
           idx = i
           break
       }
   }
   ord = append(ord[idx:], ord[:idx]...)
   // sort by polar angle around mn
   sort.Slice(ord[1:], func(i, j int) bool {
       a := ord[1+i]
       b := ord[1+j]
       dx1, dy1 := p[a].x-p[mn].x, p[a].y-p[mn].y
       dx2, dy2 := p[b].x-p[mn].x, p[b].y-p[mn].y
       return int64(dx1)*int64(dy2) - int64(dy1)*int64(dx2) > 0
   })
   // operations and DSU
   ops := make([][2]int, 0, n)
   makeOp := func(i, j int) {
       to[i], to[j] = to[j], to[i]
       ops = append(ops, [2]int{i, j})
   }
   par := make([]int, n)
   for i := range par {
       par[i] = i
   }
   var find func(int) int
   find = func(v int) int {
       if par[v] != v {
           par[v] = find(par[v])
       }
       return par[v]
   }
   used := make([]bool, n)
   for i := 0; i < n; i++ {
       if used[i] {
           continue
       }
       for j := i; !used[j]; j = to[j] {
           par[find(j)] = find(i)
           used[j] = true
       }
   }
   for i := 1; i < len(ord)-1; i++ {
       u := ord[i]
       v := ord[i+1]
       r1, r2 := find(u), find(v)
       if r1 == r2 {
           continue
       }
       makeOp(u, v)
       par[r1] = r2
   }
   // finalize cycle from mn
   for i := range used {
       used[i] = false
   }
   path := make([]int, 0, n)
   for i := mn; !used[i]; i = to[i] {
       path = append(path, i)
       used[i] = true
   }
   for j := 1; j < len(path); j++ {
       makeOp(mn, path[j])
   }
   // output
   fmt.Fprintln(out, len(ops))
   for _, pr := range ops {
       fmt.Fprintln(out, pr[0]+1, pr[1]+1)
   }
}
