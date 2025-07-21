package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   b := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // Special case: single node
   if n <= 1 {
       fmt.Println(1)
       return
   }
   // Number of children positions
   m := n - 1
   // breakChild[i] = true if between child i and i+1 must be different parent (b[i+2] < b[i+1])
   breakChild := make([]bool, m)
   for i := 1; i < m; i++ {
       // child i is b[i+1], next child is b[i+2]
       if b[i+2] < b[i+1] {
           breakChild[i] = true
       }
   }
   // dp arrays: prev and curr
   // dpPrev[k] = number of ways for first i children with parent index p_i = k
   dpPrev := make([]int, m+2)
   dpPrev[1] = 1
   // iterate children positions
   for i := 1; i < m; i++ {
       // compute prefix sums of dpPrev
       pre := make([]int, i+2)
       pre[0] = 0
       for k := 1; k <= i; k++ {
           pre[k] = pre[k-1] + dpPrev[k]
           if pre[k] >= MOD {
               pre[k] -= MOD
           }
       }
       // dpCurr for i+1 children, k from 1..i+1
       dpCurr := make([]int, m+2)
       for k := 1; k <= i+1; k++ {
           if breakChild[i] {
               // strict: previous p < k
               if k-1 >= 1 {
                   dpCurr[k] = pre[k-1]
               }
           } else {
               // non-decreasing
               dpCurr[k] = pre[min(k, i)]
           }
       }
       dpPrev = dpCurr
   }
   // sum ways for last child position m
   var ans int
   for k := 1; k <= m; k++ {
       ans += dpPrev[k]
       if ans >= MOD {
           ans -= MOD
       }
   }
   fmt.Println(ans)
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
