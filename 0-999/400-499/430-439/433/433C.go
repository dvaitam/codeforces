package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // neighbors lists for each page
   neighbors := make([][]int, n+1)
   var total int64
   for i := 1; i < m; i++ {
       u, v := a[i-1], a[i]
       d := u - v
       if d < 0 {
           d = -d
       }
       total += int64(d)
       // record adjacency
       neighbors[u] = append(neighbors[u], v)
       neighbors[v] = append(neighbors[v], u)
   }
   var bestGain int64 = 0
   // evaluate for each page
   for p := 1; p <= n; p++ {
       nb := neighbors[p]
       if len(nb) == 0 {
           continue
       }
       // current cost contributions for p
       var curr int64
       for _, t := range nb {
           diff := p - t
           if diff < 0 {
               diff = -diff
           }
           curr += int64(diff)
       }
       // find best target x = median of neighbors
       sort.Ints(nb)
       k := len(nb)
       med := nb[k/2]
       var newCost int64
       for _, t := range nb {
           diff := med - t
           if diff < 0 {
               diff = -diff
           }
           newCost += int64(diff)
       }
       gain := curr - newCost
       if gain > bestGain {
           bestGain = gain
       }
   }
   result := total - bestGain
   fmt.Fprintln(writer, result)
}
