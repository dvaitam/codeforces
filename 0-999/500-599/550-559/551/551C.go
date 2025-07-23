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
   a := make([]int64, n)
   var sum int64
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       sum += a[i]
   }
   // binary search on time t
   var lo, hi int64 = 0, sum + int64(n) + 5
   for lo+1 < hi {
       mid := (lo + hi) / 2
       if can(a, n, m, mid) {
           hi = mid
       } else {
           lo = mid
       }
   }
   // hi is minimal time
   fmt.Println(hi)
}

// can all boxes be removed in time t?
func can(a []int64, n, m int, t int64) bool {
   // copy remaining boxes
   rem := make([]int64, n)
   copy(rem, a)
   cur := 0
   for s := 0; s < m && cur < n; s++ {
       op := t
       // initial move to pile 1
       if op <= 0 {
           break
       }
       op--
       if op <= 0 {
           continue
       }
       // move from pile1 to current pile index cur+1: cost cur
       if int64(cur) > op {
           continue
       }
       op -= int64(cur)
       // now at pile cur
       for op > 0 && cur < n {
           if rem[cur] <= op {
               // remove all remaining at cur
               op -= rem[cur]
               rem[cur] = 0
               cur++
               // move to next pile if any
               if cur < n {
                   if op <= 0 {
                       break
                   }
                   op--
               }
           } else {
               // remove part
               rem[cur] -= op
               op = 0
               break
           }
       }
   }
   return cur >= n
}
