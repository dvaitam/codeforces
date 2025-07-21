package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // Read planet data
   names := make([]string, n)
   a := make([][]int, n)
   b := make([][]int, n)
   c := make([][]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &names[i])
       a[i] = make([]int, m)
       b[i] = make([]int, m)
       c[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j], &b[i][j], &c[i][j])
       }
   }
   best := 0
   // Try all buy planet i and sell planet h
   for i := 0; i < n; i++ {
       for h := 0; h < n; h++ {
           if i == h {
               continue
           }
           // Collect profitable items
           profits := make([]int, 0)
           for j := 0; j < m; j++ {
               p := b[h][j] - a[i][j]
               if p > 0 {
                   for cnt := 0; cnt < c[i][j]; cnt++ {
                       profits = append(profits, p)
                   }
               }
           }
           if len(profits) == 0 {
               continue
           }
           sort.Sort(sort.Reverse(sort.IntSlice(profits)))
           sum := 0
           limit := k
           if len(profits) < k {
               limit = len(profits)
           }
           for idx := 0; idx < limit; idx++ {
               sum += profits[idx]
           }
           if sum > best {
               best = sum
           }
       }
   }
   fmt.Println(best)
}
