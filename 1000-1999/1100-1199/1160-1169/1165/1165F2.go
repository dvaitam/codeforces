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
   offers := make([][]int, n)
   maxDay := 0
   for j := 0; j < m; j++ {
       var d, t int
       fmt.Fscan(reader, &d, &t)
       if d > maxDay {
           maxDay = d
       }
       // convert t to 0-based
       t--
       offers[t] = append(offers[t], d)
   }
   for i := 0; i < n; i++ {
       sort.Ints(offers[i])
   }
   // binary search on day
   // lower bound l=0 (impossible), r = max(maxDay, 2*totalK)
   left, right := 0, max(maxDay, 2*totalK)
   for left+1 < right {
       mid := (left + right) / 2
       if can(mid, k, offers, totalK) {
           right = mid
       } else {
           left = mid
       }
   }
   // right is first possible day
   fmt.Println(right)
}

// can we finish by day D?
func can(D int, k []int, offers [][]int, totalK int) bool {
   // count number of cheap buys available
   cheap := 0
   // for each type, c = number of offers <= D, at most k[i]
   for i := 0; i < len(k); i++ {
       if k[i] == 0 || len(offers[i]) == 0 {
           continue
       }
       // find count of days <= D
       c := sort.SearchInts(offers[i], D+1)
       if c > k[i] {
           c = k[i]
       }
       cheap += c
   }
   // total cost = full price cost - savings: 2*totalK - cheap
   need := 2*totalK - cheap
   return need <= D
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
