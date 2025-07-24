package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int64
   var n int
   if _, err := fmt.Fscan(in, &m); err != nil {
       return
   }
   fmt.Fscan(in, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // compute gaps between consecutive owls
   g := make([]int64, n)
   for i := 0; i < n; i++ {
       j := i + 1
       if j == n {
           j = 0
       }
       diff := a[j] - a[i] - 1
       if j == 0 {
           diff = a[0] + m - a[i] - 1
       }
       if diff < 0 {
           diff += m
       }
       g[i] = diff
   }
   // precompute minBoth (ceil((g+1)/2)) and minSingle (g)
   minBoth := make([]int64, n)
   for i := 0; i < n; i++ {
       if g[i] <= 0 {
           minBoth[i] = 0
       } else {
           d := g[i] + 1
           minBoth[i] = (d + 1) / 2
       }
   }
   // function to check if time T is sufficient
   can := func(T int64) bool {
       // ctype: 0 ANY, 1 OR, 2 AND
       c := make([]uint8, n)
       for i := 0; i < n; i++ {
           if g[i] > 0 {
               if T < minBoth[i] {
                   return false
               }
               if T < g[i] {
                   c[i] = 2 // AND
               } else {
                   c[i] = 1 // OR
               }
           } else {
               c[i] = 0 // ANY
           }
       }
       // try both initial x[0] = 0 or 1
       for start := 0; start < 2; start++ {
           curr0 := false
           curr1 := false
           if start == 0 {
               curr0 = true
           } else {
               curr1 = true
           }
           // propagate
           for i := 0; i < n; i++ {
               next0, next1 := false, false
               switch c[i] {
               case 0: // ANY
                   if curr0 {
                       next0 = true
                       next1 = true
                   }
                   if curr1 {
                       next0 = true
                       next1 = true
                   }
               case 1: // OR: allow [0,0], [1,0], [1,1]
                   if curr0 {
                       // from 0, only to 0
                       next0 = true
                   }
                   if curr1 {
                       // from 1, to 0 or 1
                       next0 = true
                       next1 = true
                   }
               case 2: // AND: only [1,0]
                   if curr1 {
                       next0 = true
                   }
               }
               curr0, curr1 = next0, next1
               if !curr0 && !curr1 {
                   break
               }
           }
           // must end with x[n] == x[0] == start
           if (start == 0 && curr0) || (start == 1 && curr1) {
               return true
           }
       }
       return false
   }
   // binary search on T
   var lo, hi int64 = 0, m
   for lo < hi {
       mid := (lo + hi) / 2
       if can(mid) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(lo)
}
