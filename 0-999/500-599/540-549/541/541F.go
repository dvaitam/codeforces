package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, T int
   if _, err := fmt.Fscan(in, &n, &T); err != nil {
       return
   }
   // buckets[d] holds qi for tasks with max depth L_i = d
   buckets := make([][]int, T)
   for i := 0; i < n; i++ {
       var ti, qi int
       fmt.Fscan(in, &ti, &qi)
       // compute allowed depth
       L := T - ti
       if L < 0 {
           continue
       }
       if L >= 0 && L < T {
           buckets[L] = append(buckets[L], qi)
       }
   }
   // sort buckets descending and build prefix sums
   prefs := make([][]int, T)
   for d := 0; d < T; d++ {
       b := buckets[d]
       if len(b) == 0 {
           prefs[d] = []int{0}
           continue
       }
       sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
       pref := make([]int, len(b)+1)
       pref[0] = 0
       for i, v := range b {
           pref[i+1] = pref[i] + v
       }
       prefs[d] = pref
   }
   // dp arrays: prev and next, size up to n free slots
   // prev[s] = max profit at current depth with s free slots
   const NEG_INF = -1_000_000_000
   prev := make([]int, n+1)
   for i := range prev {
       prev[i] = NEG_INF
   }
   // start with one slot at depth 0
   prev[1] = 0
   // iterate depths
   for d := 0; d < T; d++ {
       next := make([]int, n+1)
       for i := range next {
           next[i] = NEG_INF
       }
       pref := prefs[d]
       m := len(pref) - 1
       for s := 0; s <= n; s++ {
           base := prev[s]
           if base < 0 {
               continue
           }
           // try taking x tasks at this depth
           maxTake := m
           if s < maxTake {
               maxTake = s
           }
           for x := 0; x <= maxTake; x++ {
               profit := base + pref[x]
               newSlots := (s - x) * 2
               if newSlots > n {
                   newSlots = n
               }
               if profit > next[newSlots] {
                   next[newSlots] = profit
               }
           }
       }
       prev = next
   }
   // answer is max over prev
   ans := 0
   for _, v := range prev {
       if v > ans {
           ans = v
       }
   }
   fmt.Println(ans)
}
