package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   k := make([]int, n)
   totalK := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &k[i])
       totalK += k[i]
   }
   // sale days per type
   sales := make([][]int, n)
   for j := 0; j < m; j++ {
       var d, t int
       fmt.Fscan(reader, &d, &t)
       if t >= 1 && t <= n {
           sales[t-1] = append(sales[t-1], d)
       }
   }
   for i := 0; i < n; i++ {
       sort.Ints(sales[i])
   }
   // binary search minimal day
   lo, hi := 0, 3000
   for lo < hi {
       mid := (lo + hi) / 2
       if canFinish(mid, k, sales, totalK) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(lo)
}

// canFinish checks if all transactions can be bought by day D
func canFinish(D int, k []int, sales [][]int, totalK int) bool {
   // count cheap buys from sale events
   cheap := 0
   for i, ki := range k {
       if ki <= 0 {
           continue
       }
       // count number of sale events <= D
       ci := sort.SearchInts(sales[i], D+1)
       if ci > ki {
           ci = ki
       }
       cheap += ci
   }
   expensive := totalK - cheap
   // total cost: cheap*1 + expensive*2
   cost := cheap + expensive*2
   return cost <= D
}
