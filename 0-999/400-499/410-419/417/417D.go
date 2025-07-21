package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var b int64
   if _, err := fmt.Fscan(in, &n, &m, &b); err != nil {
       return
   }
   type Friend struct {
       k   int64
       x   int64
       mask int
   }
   friends := make([]Friend, n)
   for i := 0; i < n; i++ {
       var xi, ki int64
       var mi int
       fmt.Fscan(in, &xi, &ki, &mi)
       mask := 0
       for j := 0; j < mi; j++ {
           var p int
           fmt.Fscan(in, &p)
           if p > 0 && p <= m {
               mask |= 1 << (p - 1)
           }
       }
       friends[i] = Friend{k: ki, x: xi, mask: mask}
   }
   // sort by required monitors k ascending
   sort.Slice(friends, func(i, j int) bool {
       return friends[i].k < friends[j].k
   })
   full := (1 << m) - 1
   INF := int64(9e18)
   dp := make([]int64, 1<<m)
   for i := 1; i <= full; i++ {
       dp[i] = INF
   }
   dp[0] = 0
   ans := INF
   // process friends grouped by k
   for i := 0; i < n; {
       curK := friends[i].k
       j := i
       for j < n && friends[j].k == curK {
           f := friends[j]
           // update dp with this friend
           xi := f.x
           fm := f.mask
           // iterate masks in reverse to avoid reuse
           for mask := full; ; mask-- {
               prev := dp[mask]
               if prev < INF {
                   nm := mask | fm
                   cost := prev + xi
                   if cost < dp[nm] {
                       dp[nm] = cost
                   }
               }
               if mask == 0 {
                   break
               }
           }
           j++
       }
       // after adding all friends with curK, check if full covered
       if dp[full] < INF {
           total := dp[full] + curK*b
           if total < ans {
               ans = total
           }
       }
       i = j
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
