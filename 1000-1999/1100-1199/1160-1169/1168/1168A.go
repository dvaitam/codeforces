package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // binary search on number of operations
   lo, hi := -1, m
   for lo+1 < hi {
       mid := (lo + hi) >> 1
       if ok(a, n, m, mid) {
           hi = mid
       } else {
           lo = mid
       }
   }
   fmt.Println(hi)
}

// ok checks if with x operations the array can be made non-decreasing
func ok(a []int, n, m, x int) bool {
   cur := 0
   for i := 0; i < n; i++ {
       ai := a[i]
       if ai + x < m {
           // cannot wrap
           if ai < cur {
               return false
           }
           // set to original
           cur = ai
       } else {
           // wrap possible: reachable [0..wrap] and [ai..m-1]
           wrap := (ai + x) % m
           if wrap < cur {
               // cannot wrap to cur, must use [ai..m-1]
               if ai < cur {
                   return false
               }
               cur = ai
           }
           // else keep cur unchanged
       }
   }
   return true
}
