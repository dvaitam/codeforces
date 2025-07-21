package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   w := make([]int, n)
   h := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &w[i], &h[i])
   }
   k := n / 2
   // collect candidate heights
   candMap := make(map[int]struct{})
   for i := 0; i < n; i++ {
       candMap[w[i]] = struct{}{}
       candMap[h[i]] = struct{}{}
   }
   cand := make([]int, 0, len(candMap))
   for hh := range candMap {
       cand = append(cand, hh)
   }
   sort.Ints(cand)
   const INF64 = int64(4e18)
   best := INF64
   for _, H := range cand {
       // evaluate for height H
       m := 0 // mandatory lying
       ok := true
       baseW := int64(0)
       deltas := make([]int64, 0)
       for i := 0; i < n; i++ {
           wi, hi := w[i], h[i]
           wi64, hi64 := int64(wi), int64(hi)
           if wi <= H && hi <= H {
               // optional: standing gives width wi, lying gives hi
               baseW += wi64
               deltas = append(deltas, hi64-wi64)
           } else if wi <= H {
               // only lying (height=wi)
               m++
               baseW += hi64 // lying width = hi
           } else if hi <= H {
               // only standing
               baseW += wi64
           } else {
               ok = false
               break
           }
       }
       if !ok || m > k {
           continue
       }
       // choose optional to lie: at most k-m, pick those with negative deltas
       need := k - m
       // sort deltas
       sort.Slice(deltas, func(i, j int) bool { return deltas[i] < deltas[j] })
       // apply negative deltas up to need
       for i, d := range deltas {
           if i >= need {
               break
           }
           if d < 0 {
               baseW += d
           } else {
               break
           }
       }
       area := baseW * int64(H)
       if area < best {
           best = area
       }
   }
   // output
   if best == INF64 {
       fmt.Println(0)
   } else {
       fmt.Println(best)
   }
}
