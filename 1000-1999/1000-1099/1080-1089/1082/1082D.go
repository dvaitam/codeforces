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
   d := make([]int, n+1)
   cnt := 0
   tot := 0
   type node struct{ deg, idx int }
   nodes := make([]node, n)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &d[i])
       nodes[i-1] = node{d[i], i}
       if d[i] == 1 {
           cnt++
       } else {
           tot += d[i]
       }
   }
   // All ones
   if cnt == n {
       if n <= 2 {
           // diameter = n-1, edges = n-1
           fmt.Fprintf(out, "YES %d\n%d\n", n-1, n-1)
           if n > 1 {
               fmt.Fprintln(out, "1 2")
           }
       } else {
           fmt.Fprintln(out, "NO")
       }
       return
   }
   // subtract internal chain degrees
   heavyCount := n - cnt
   tot -= (heavyCount - 1) * 2
   if tot < cnt {
       fmt.Fprintln(out, "NO")
       return
   }
   // sort by degree ascending
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].deg < nodes[j].deg
   })
   used := make([]int, n+1)
   var ans [][2]int
   ins := func(u, v int) {
       ans = append(ans, [2]int{u, v})
       used[u]++
       used[v]++
   }
   // build main chain of heavy nodes
   for i := cnt; i < n-1; i++ {
       ins(nodes[i].idx, nodes[i+1].idx)
   }
   // attach first two leaves
   if cnt > 0 {
       ins(nodes[0].idx, nodes[cnt].idx)
   }
   if cnt > 1 {
       ins(nodes[1].idx, nodes[n-1].idx)
   }
   // attach remaining leaves
   j := n - 1
   for i := 2; i < cnt; i++ {
       // find next heavy with available degree
       for j > cnt && used[nodes[j].idx] >= nodes[j].deg {
           j--
       }
       ins(nodes[i].idx, nodes[j].idx)
   }
   // calculate diameter
   diameter := heavyCount - 1
   // if extra branches increase diameter
   // last heavy end
   if used[nodes[n-1].idx] >= 2 {
       diameter++
   }
   // first heavy start
   if used[nodes[cnt].idx] >= 2 {
       diameter++
   }
   // output
   fmt.Fprintf(out, "YES %d\n", diameter)
   fmt.Fprintln(out, len(ans))
   for _, e := range ans {
       fmt.Fprintf(out, "%d %d\n", e[0], e[1])
   }
}
