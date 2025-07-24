package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   var totalSlices int64
   maxA := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       totalSlices += int64(a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   if totalSlices < k {
       fmt.Println(-1)
       return
   }
   lo, hi, ans := 1, maxA, 1
   for lo <= hi {
       mid := (lo + hi) / 2
       if can(a, k, mid, totalSlices) {
           ans = mid
           lo = mid + 1
       } else {
           hi = mid - 1
       }
   }
   fmt.Println(ans)
}

// can checks if it's possible to get at least k parts of size >= m
func can(a []int, k int64, m int, totalSlices int64) bool {
   if m == 1 {
       return totalSlices >= k
   }
   var cntTotal int64
   // stack for manual DFS
   stack := make([]int, 0, 32)
   for _, v := range a {
       if v < m {
           continue
       }
       // count pieces from this tangerine
       cnt := int64(0)
       stack = stack[:0]
       stack = append(stack, v)
       for len(stack) > 0 {
           x := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           if x < m {
               continue
           }
           if x < 2*m {
               cnt++
               if cntTotal+cnt >= k {
                   return true
               }
           } else {
               // split into two parts: floor and ceil
               half := x / 2
               stack = append(stack, half)
               stack = append(stack, x-half)
           }
       }
       cntTotal += cnt
       if cntTotal >= k {
           return true
       }
   }
   return false
}
