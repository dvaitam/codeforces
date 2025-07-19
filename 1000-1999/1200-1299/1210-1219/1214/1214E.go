package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   type pair struct{ d, u, v int }
   p := make([]pair, n)
   for i := 0; i < n; i++ {
       var di int
       fmt.Fscan(in, &di)
       p[i] = pair{di, 2*i + 1, 2*i + 2}
   }
   sort.Slice(p, func(i, j int) bool { return p[i].d > p[j].d })

   maxSize := 2*n + 5
   is := make([]int, maxSize)
   to := make([][]int, maxSize)
   t := n
   // assign primary nodes
   for i := 1; i <= n; i++ {
       is[i] = p[i-1].u
   }
   // build chain and extra edges
   for i := 1; i <= n; i++ {
       di := p[i-1].d
       y := p[i-1].v
       nx := i + di - 1
       if nx == t {
           t++
           is[t] = y
       } else {
           to[nx] = append(to[nx], y)
       }
   }
   // output edges on main chain
   for i := 1; i < t; i++ {
       fmt.Fprintln(out, is[i], is[i+1])
   }
   // output extra edges
   for i := 1; i <= t; i++ {
       for _, j := range to[i] {
           fmt.Fprintln(out, is[i], j)
       }
   }
}
