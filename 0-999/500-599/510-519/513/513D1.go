package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, c int
   if _, err := fmt.Fscan(reader, &n, &c); err != nil {
       return
   }
   // constraints for each a: list of (b, dir)
   leftCons := make([][]int, n+2)
   rightCons := make([][]int, n+2)
   for i := 0; i < c; i++ {
       var a, b int
       var dir string
       fmt.Fscan(reader, &a, &b, &dir)
       if b <= a {
           fmt.Println("IMPOSSIBLE")
           return
       }
       if dir == "LEFT" {
           leftCons[a] = append(leftCons[a], b)
       } else {
           rightCons[a] = append(rightCons[a], b)
       }
   }
   // dp[l][r]: -1 unknown, 0 false, 1 true
   dp := make([][]int, n+2)
   choice := make([][]int, n+2)
   for i := 0; i <= n+1; i++ {
       dp[i] = make([]int, n+2)
       choice[i] = make([]int, n+2)
       for j := range dp[i] {
           dp[i][j] = -1
       }
   }
   var solve func(l, r int) bool
   solve = func(l, r int) bool {
       if l > r {
           return true
       }
       if dp[l][r] != -1 {
           return dp[l][r] == 1
       }
       // constraints at root l must have b in [l+1..r]
       // check b bounds
       for _, b := range leftCons[l] {
           if b > r {
               dp[l][r] = 0
               return false
           }
       }
       for _, b := range rightCons[l] {
           if b > r {
               dp[l][r] = 0
               return false
           }
       }
       // determine k bounds
       total := r - l
       lower := 0
       upper := total
       if len(leftCons[l]) > 0 {
           maxb := l + 1
           for _, b := range leftCons[l] {
               if b > maxb {
                   maxb = b
               }
           }
           lower = maxb - l
           if lower < 1 {
               lower = 1
           }
       }
       if len(rightCons[l]) > 0 {
           minb := r
           for _, b := range rightCons[l] {
               if b < minb {
                   minb = b
               }
           }
           // b >= l+k+1 => k <= b-l-1
           ub := minb - l - 1
           if ub > total {
               ub = total
           }
           // need right subtree at least size 1
           if ub > total-1 {
               ub = total - 1
           }
           upper = ub
       }
       if lower > upper {
           dp[l][r] = 0
           return false
       }
       // try k
       for k := lower; k <= upper; k++ {
           if solve(l+1, l+k) && solve(l+k+1, r) {
               choice[l][r] = k
               dp[l][r] = 1
               return true
           }
       }
       dp[l][r] = 0
       return false
   }
   if !solve(1, n) {
       fmt.Println("IMPOSSIBLE")
       return
   }
   // reconstruct in-order
   out := make([]int, 0, n)
   var build func(l, r int)
   build = func(l, r int) {
       if l > r {
           return
       }
       k := choice[l][r]
       build(l+1, l+k)
       out = append(out, l)
       build(l+k+1, r)
   }
   build(1, n)
   w := bufio.NewWriter(os.Stdout)
   for i, v := range out {
       if i > 0 {
           w.WriteByte(' ')
       }
       fmt.Fprint(w, v)
   }
   w.WriteByte('\n')
   w.Flush()
}
