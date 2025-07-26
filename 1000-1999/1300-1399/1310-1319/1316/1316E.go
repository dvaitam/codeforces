package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "math/bits"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, p, k int
   fmt.Fscan(in, &n, &p, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   type person struct {
       a int64
       s []int64
   }
   people := make([]person, n)
   for i := 0; i < n; i++ {
       people[i].a = a[i]
       people[i].s = make([]int64, p)
       for j := 0; j < p; j++ {
           fmt.Fscan(in, &people[i].s[j])
       }
   }
   // sort by audience strength descending
   sort.Slice(people, func(i, j int) bool {
       return people[i].a > people[j].a
   })

   fullMask := 1<<p
   // precompute popcounts
   pop := make([]int, fullMask)
   for m := 1; m < fullMask; m++ {
       pop[m] = pop[m>>1] + (m & 1)
   }
   const negInf = -1 << 60
   dp := make([]int64, fullMask)
   next := make([]int64, fullMask)
   for m := 1; m < fullMask; m++ {
       dp[m] = negInf
   }
   // DP over people
   for i := 0; i < n; i++ {
       // copy dp to next
       copy(next, dp)
       pi := &people[i]
       // audience option
       for m := 0; m < fullMask; m++ {
           used := i - pop[m]
           if used < k {
               val := dp[m] + pi.a
               if val > next[m] {
                   next[m] = val
               }
           }
       }
       // position assignment
       for m := 0; m < fullMask; m++ {
           if dp[m] == negInf {
               continue
           }
           for j := 0; j < p; j++ {
               bit := 1 << j
               if m & bit == 0 {
                   nm := m | bit
                   val := dp[m] + pi.s[j]
                   if val > next[nm] {
                       next[nm] = val
                   }
               }
           }
       }
       // swap dp and next
       dp, next = next, dp
   }
   // answer is dp[fullMask-1]
   fmt.Fprintln(out, dp[fullMask-1])
}
