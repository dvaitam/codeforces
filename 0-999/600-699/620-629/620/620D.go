package main

import (
   "fmt"
   "sort"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   var n, m int
   fmt.Scan(&n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&a[i])
   }
   fmt.Scan(&m)
   b := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Scan(&b[i])
   }
   // initial sums
   var sa, sb int64
   for _, v := range a {
       sa += v
   }
   for _, v := range b {
       sb += v
   }
   best := abs64(sa - sb)
   bestCnt := 0
   bestIdx := make([]int, 4)
   // single swap
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           diff := abs64(sa - sb - 2*a[i] + 2*b[j])
           if diff < best {
               best = diff
               bestCnt = 1
               bestIdx[0] = i
               bestIdx[1] = j
           }
       }
   }
   // double swap
   if n >= 2 && m >= 2 {
       ca := n * (n - 1) / 2
       cb := m * (m - 1) / 2
       asum := make([]int64, ca)
       bsum := make([]int64, cb)
       idx := 0
       for i := 0; i < n; i++ {
           for j := i + 1; j < n; j++ {
               asum[idx] = a[i] + a[j]
               idx++
           }
       }
       idx = 0
       for i := 0; i < m; i++ {
           for j := i + 1; j < m; j++ {
               bsum[idx] = b[i] + b[j]
               idx++
           }
       }
       sort.Slice(asum, func(i, j int) bool { return asum[i] < asum[j] })
       sort.Slice(bsum, func(i, j int) bool { return bsum[i] < bsum[j] })
       ptr := 0
       base := sa - sb
       for i := 0; i < ca; i++ {
           // move ptr to best for this asum[i]
           for ptr < cb-1 && abs64(base-2*asum[i]+2*bsum[ptr]) >= abs64(base-2*asum[i]+2*bsum[ptr+1]) {
               ptr++
           }
           diff := abs64(base - 2*asum[i] + 2*bsum[ptr])
           if diff < best {
               best = diff
               bestCnt = 2
               bestIdx[0] = i
               bestIdx[1] = ptr
           }
       }
       if bestCnt == 2 {
           // find original pairs for a
           sumA := asum[bestIdx[0]]
           found := false
           for i := 0; i < n && !found; i++ {
               for j := i + 1; j < n; j++ {
                   if a[i]+a[j] == sumA {
                       bestIdx[0] = i
                       bestIdx[2] = j
                       found = true
                       break
                   }
               }
           }
           // find original pairs for b
           sumB := bsum[bestIdx[1]]
           found = false
           for i := 0; i < m && !found; i++ {
               for j := i + 1; j < m; j++ {
                   if b[i]+b[j] == sumB {
                       bestIdx[1] = i
                       bestIdx[3] = j
                       found = true
                       break
                   }
               }
           }
       }
   }
   // output
   fmt.Println(best)
   fmt.Println(bestCnt)
   for i := 0; i < bestCnt; i++ {
       ai := bestIdx[2*i] + 1
       bi := bestIdx[2*i+1] + 1
       fmt.Println(ai, bi)
   }
}
