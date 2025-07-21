package main

import (
   "bufio"
   "fmt"
   "math"
   "math/bits"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]int64, n)
   for i := 0; i < n; i++ {
       var v int64
       fmt.Fscan(reader, &v)
       s[i] = v
   }
   var m int
   fmt.Fscan(reader, &m)
   // actions: m lines of action and team
   isPick := make([]bool, m)
   team := make([]int, m)
   for i := 0; i < m; i++ {
       var ac string
       var t int
       fmt.Fscan(reader, &ac, &t)
       if len(ac) > 0 && ac[0] == 'p' {
           isPick[i] = true
       }
       team[i] = t
   }
   // take top m heroes by strength
   sort.Slice(s, func(i, j int) bool { return s[i] > s[j] })
   if m > n {
       m = n
       isPick = isPick[:m]
       team = team[:m]
   }
   a := make([]int64, m)
   for i := 0; i < m; i++ {
       a[i] = s[i]
   }
   // dp over bitmask of used heroes
   total := 1 << m
   dp := make([]int64, total)
   // iterate masks from high to low
   for mask := total - 1; mask >= 0; mask-- {
       k := bits.OnesCount(uint(mask))
       if k >= m {
           dp[mask] = 0
           continue
       }
       actPick := isPick[k]
       actTeam := team[k]
       if actTeam == 1 {
           // maximize
           best := int64(math.MinInt64 / 2)
           for i := 0; i < m; i++ {
               if mask&(1<<i) != 0 {
                   continue
               }
               next := mask | (1 << i)
               val := dp[next]
               if actPick {
                   val += a[i]
               }
               if val > best {
                   best = val
               }
           }
           dp[mask] = best
       } else {
           // team 2: minimize
           best := int64(math.MaxInt64 / 2)
           for i := 0; i < m; i++ {
               if mask&(1<<i) != 0 {
                   continue
               }
               next := mask | (1 << i)
               val := dp[next]
               if actPick {
                   val -= a[i]
               }
               if val < best {
                   best = val
               }
           }
           dp[mask] = best
       }
   }
   // result
   fmt.Println(dp[0])
}
